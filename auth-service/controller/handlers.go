package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/ably/ably-go/ably"
	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/auth-service/realtime"
	"github.com/quible-io/quible-api/auth-service/services/chatService"
	"github.com/quible-io/quible-api/auth-service/services/emailService"
	"github.com/quible-io/quible-api/auth-service/services/userService"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/models"
	"golang.org/x/sync/errgroup"
)

var UserFields = []string{"id", "username", "email", "phone", "full_name"}
var PublicUserFields = []string{"id", "full_name"}

// @Summary		Register
// @Description	Register a new user.
// @Tags			user,public
// @Accept		json
// @Produce		json
// @Param			request	body		userService.UserRegisterDTO	true	"User registration information"
// @Success		201		{object}	UserResponse
// @Failure		400		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/user [post]
func UserRegister(c *gin.Context) {
	us := getUserServiceFromContext(c)
	var userRegisterDTO userService.UserRegisterDTO
	var errorCode ErrorCode
	if err := c.ShouldBindJSON(&userRegisterDTO); err != nil {
		errorCode = Err400_InvalidRequestBody
		errorFields := misc.ParseValidationError(err)
		if errorFields.IsValidationError {
			if errorFields.CheckAll("Email") {
				errorCode = Err400_InvalidEmailFormat
			} else if errorFields.CheckAll("Phone") {
				errorCode = Err400_InvalidPhoneFormat
			}
		} else {
			errorCode = Err400_MalformedJSON
		}
		log.Printf("unmet request body constraints: %q", errorFields.GetAllFields())
		SendError(c, http.StatusBadRequest, errorCode)
		return
	}
	foundUser, _ := us.GetUserByUsernameOrEmail(&userRegisterDTO)
	if foundUser != nil && foundUser.ActivatedAt.Ptr() != nil {
		SendError(c, http.StatusBadRequest, Err400_UserWithEmailOrUsernameExists)
		return
	}
	// branch on whether the user exists or not
	var user *models.User
	if foundUser != nil {
		if err := us.UpdateWith(foundUser, &userRegisterDTO); err != nil {
			log.Printf("unable to update existing user with registration data: %q", err)
			SendError(c, http.StatusInternalServerError, Err500_UnableToRegister)
			return
		}
		user = foundUser
	} else {
		createdUser, err := us.CreateUser(&userRegisterDTO)
		if err != nil {
			log.Printf("unable to register user: %q", err)
			SendError(c, http.StatusInternalServerError, Err500_UnableToRegister)
			return
		}
		user = createdUser
	}
	// send activation email
	g := new(errgroup.Group)
	g.Go(
		func() error {
			token, _ := generateToken(user, Activate)
			var html bytes.Buffer
			emailService.UserActivation(
				user.FullName,
				fmt.Sprintf(
					"%s/forms/activation?token=%s",
					os.Getenv("WEB_CLIENT_URL"),
					token.String(),
				),
				&html,
			)
			return email.Send(c.Request.Context(), email.EmailDTO{
				From:     "no-reply@quible.io",
				To:       user.Email,
				Subject:  "Activate your Quible account",
				HTMLBody: html.String(),
			})
		},
	)
	if err := g.Wait(); err != nil {
		log.Printf("unable to send activation email: %q", err)
		SendError(c, http.StatusFailedDependency, Err424_UnableToSendEmail)
		return
	}
	// response with user profile as data
	c.JSON(
		http.StatusCreated,
		misc.PickFields(user, UserFields...),
	)
}

// @Summary		Activate new user
// @Description	Handles user activation initiated by clicking link from email
// @Tags			user,public
// @Produce		json
// @Param			request	body		userService.TokenDTO	true	"JWT generated during registration"
// @Success   200
// @Failure		400		{object}	ErrorResponse
// @Failure		417		{object}	ErrorResponse
// @Router		/user/activate [post]
func UserActivate(c *gin.Context) {
	us := getUserServiceFromContext(c)
	var userActivateDTO userService.TokenDTO
	if err := c.ShouldBindJSON(&userActivateDTO); err != nil {
		log.Printf("invalid request body: %q", err)
		SendError(c, http.StatusBadRequest, Err400_InvalidOrMalformedToken)
		return
	}
	tokenClaims, err := verifyJWT(userActivateDTO.Token, Activate)
	if err != nil {
		log.Printf("invalid token: %q", err)
		SendError(c, http.StatusExpectationFailed, Err417_InvalidToken)
		return
	}
	userId := tokenClaims["userId"].(string)
	if err := us.ActivateUser(userId); err != nil {
		log.Printf("unable to activate user: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnableToActivateUser)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Login
// @Description	Login with user credentials to get token
// @Tags			user,public
// @Accept		json
// @Produce		json
// @Param			request	body		userService.UserLoginDTO	true	"User login credentials"
// @Success		200		{object}	TokenResponse
// @Failure		400		{object}	ErrorResponse
// @Failure		401		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/login [post]
func UserLogin(c *gin.Context) {
	us := getUserServiceFromContext(c)

	var userLoginDTO userService.UserLoginDTO
	if err := c.ShouldBindJSON(&userLoginDTO); err != nil {
		log.Printf("invalid request body: %q", err)
		SendError(c, http.StatusBadRequest, Err400_InvalidRequestBody)
		return
	}
	foundUser, _ := us.GetUserByEmail(userLoginDTO.Email)
	if foundUser == nil {
		log.Printf("user with given email not found: %q", userLoginDTO.Email)
		SendError(c, http.StatusUnauthorized, Err401_InvalidCredentials)
		return
	}
	if foundUser.ActivatedAt.Ptr() == nil {
		log.Printf("not activated user is attempted to login: %q", userLoginDTO.Email)
		SendError(c, http.StatusUnauthorized, Err401_UserNotActivated)
		return
	}
	if err := us.ValidatePassword(foundUser.HashedPassword, userLoginDTO.Password); err != nil {
		log.Printf("invalid password: %+v", userLoginDTO)
		SendError(c, http.StatusUnauthorized, Err401_InvalidCredentials)
		return
	}

	type TokenJob struct {
		user   *models.User
		action TokenAction
		result *GeneratedToken
	}
	var generatedAccessToken, generatedRefreshToken GeneratedToken
	jobs := map[string]TokenJob{
		"access":  {foundUser, Access, &generatedAccessToken},
		"refresh": {foundUser, Refresh, &generatedRefreshToken},
	}
	g := new(errgroup.Group)
	for name, job := range jobs {
		job, name := job, name
		g.Go(
			func() error {
				generatedToken, err := generateToken(job.user, job.action)
				if err != nil {
					log.Printf("unable to generate %s token: %q", name, err)
					return err
				}
				*job.result = generatedToken
				return nil
			},
		)
	}
	if err := g.Wait(); err != nil {
		SendError(c, http.StatusInternalServerError, Err500_UnableToGenerateToken)
		return
	}

	foundUser.Refresh = generatedRefreshToken.ID
	if err := us.Update(foundUser); err != nil {
		log.Printf("unable to associate refresh token with user: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}

	responseData := TokenResponse{
		AccessToken:  generatedAccessToken.String(),
		RefreshToken: generatedRefreshToken.String(),
	}
	log.Printf(
		"login user %q with refresh signature %q\n",
		userLoginDTO.Email,
		strings.SplitN(responseData.RefreshToken, ".", 3)[2],
	)
	c.JSON(http.StatusOK, responseData)
}

// @Summary		Get user
// @Description	Returns user profile associated with the token
// @Tags			user,private
// @Produce		json
// @Success		200	{object}	UserResponse
// @Failure		401	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/user [get]
func UserGet(c *gin.Context) {
	user := getUserFromContext(c)
	c.JSON(
		http.StatusCreated,
		misc.PickFields(user, UserFields...),
	)
}

// @Summary		Get public user profile by ID
// @Description	Returns user profile corresponding to provided ID
// @Tags			user,private
// @Produce		json
// @Param     userId   path   string  true  "User ID"
// @Success		200	{object}	PublicUserRecord
// @Failure		401	{object}	ErrorResponse
// @Failure		404	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/user/{userId}/profile [get]
func UserGetById(c *gin.Context) {
	userId := c.Param("userId")
	us := getUserServiceFromContext(c)
	user, err := us.GetUserById(userId)
	if err != nil || user == nil {
		log.Printf("user not found: %q", userId)
		SendError(c, http.StatusNotFound, Err404_UserNotFound)
		return
	}
	imageData := us.GetUserImage(user)
	var imageDataURL string
	if imageData != nil {
		imageDataURL = fmt.Sprintf(
			"data:%s;base64,%s",
			imageData.ContentType,
			base64.StdEncoding.EncodeToString(imageData.BinaryContent),
		)
	}
	result := misc.PickFields(user, PublicUserFields...)
	result["image"] = nil
	if len(imageDataURL) > 0 {
		result["image"] = &imageDataURL
	}
	c.JSON(
		http.StatusOK,
		result,
	)
}

// @Summary		Update user
// @Description	Updates user profile associated with the token
// @Tags			user,private
// @Accept		json
// @Produce		json
// @Param			request	body		userService.UserPatchDTO	true	"Partial user object to be used for update"
// @Success		200		{object}	UserResponse
// @Failure		400		{object}	ErrorResponse
// @Failure		401		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/user [patch]
func UserPatch(c *gin.Context) {
	var userPatchDTO userService.UserPatchDTO
	var errorCode ErrorCode
	if err := c.ShouldBindJSON(&userPatchDTO); err != nil {
		errorCode = Err400_InvalidRequestBody
		errorFields := misc.ParseValidationError(err)
		if errorFields.IsValidationError {
			if errorFields.CheckAll("Email") {
				errorCode = Err400_InvalidEmailFormat
			} else if errorFields.CheckAll("Phone") {
				errorCode = Err400_InvalidPhoneFormat
			}
		} else {
			errorCode = Err400_MalformedJSON
		}
		log.Printf("unmet request body constraints: %q", errorFields.GetAllFields())
		SendError(c, http.StatusBadRequest, errorCode)
		return
	}
	user := getUserFromContext(c)
	dtoType := reflect.TypeOf(userPatchDTO)
	dtoValue := reflect.ValueOf(userPatchDTO)
	userValue := reflect.ValueOf(user).Elem()
	for i := 0; i < dtoValue.NumField(); i++ {
		dtoFieldName := dtoType.Field(i).Name
		dtoFieldValue := dtoValue.Field(i).Elem()
		if dtoFieldValue.IsValid() {
			target := userValue.FieldByName(dtoFieldName)
			if target.Kind().String() == "struct" && target.Type().String() == "null.String" {
				target = target.FieldByName("String")
			}
			if target.CanSet() {
				target.SetString(dtoFieldValue.String())
			}
		}
	}
	us := getUserServiceFromContext(c)
	if err := us.Update(user); err != nil {
		log.Printf("unable to update user: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}
	c.JSON(
		http.StatusOK,
		misc.PickFields(user, UserFields...),
	)
}

// @Summary		Refresh access/refresh tokens
// @Description	Login with user credentials to get token
// @Tags			user,public
// @Accept		json
// @Produce		json
// @Param			request	body		userService.UserRefreshDTO	true	"User's refresh token"
// @Success		200		{object}	TokenResponse
// @Failure		400		{object}	ErrorResponse
// @Failure		401		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/user/refresh [post]
func UserRefresh(c *gin.Context) {
	var userRefreshDTO userService.UserRefreshDTO
	if err := c.ShouldBindJSON(&userRefreshDTO); err != nil {
		errorFields := misc.ParseValidationError(err)
		log.Printf("unmet request body constraints: %q", errorFields.GetAllFields())
		SendError(c, http.StatusBadRequest, Err400_InvalidRequestBody)
		return
	}

	claims, err := verifyJWT(userRefreshDTO.RefreshToken, Refresh)
	if err != nil {
		log.Printf("invalid refresh token: %q", err)
		SendError(c, http.StatusUnauthorized, Err401_InvalidRefreshToken)
		return
	}

	userId := claims["userId"].(string)
	us := getUserServiceFromContext(c)
	user, err := us.GetUserById(userId)
	if err != nil {
		log.Printf("unable to retrieve user by id from the refresh token: %q", userId)
		SendError(c, http.StatusUnauthorized, Err401_InvalidRefreshToken)
		return
	}

	if refreshTokenId := claims["jti"].(string); user.Refresh != refreshTokenId {
		log.Printf("provided refresh token %q has been revoked", userRefreshDTO.RefreshToken)
		SendError(c, http.StatusUnauthorized, Err401_InvalidRefreshToken)
		return
	}

	generatedAccessToken, err := generateToken(user, Access)
	if err != nil {
		log.Printf("unable to generate access token: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnableToGenerateToken)
		return
	}
	generatedRefreshToken, err := generateToken(user, Refresh)
	if err != nil {
		log.Printf("unable to generate refresh token: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnableToGenerateToken)
		return
	}
	user.Refresh = generatedRefreshToken.ID
	if err := us.Update(user); err != nil {
		log.Printf("unable to associate refresh token with user: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}
	responseData := TokenResponse{
		AccessToken:  generatedAccessToken.String(),
		RefreshToken: generatedRefreshToken.String(),
	}
	log.Printf(
		"refresh token for %q, new refresh signature %q\n",
		user.Email,
		strings.SplitN(responseData.RefreshToken, ".", 3)[2],
	)
	c.JSON(http.StatusOK, responseData)
}

// @Summary		Get new `TokenRequest` for client Ably SDK
// @Description	Returns user profile corresponding to provided ID
// @Tags			realtime,private
// @Produce		json
// @Success		200	{object}	AblyTokenRequest
// @Failure		400	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/rt/token [get]
func AblyToken(c *gin.Context) {
	user := getUserFromContext(c)
	if c.Request.URL.Query().Has("clientId") && c.Request.URL.Query().Get("clientId") != user.ID {
		log.Printf("provided clientId doesn't match authenticated user")
		SendError(c, http.StatusBadRequest, Err400_InvalidClientId)
		return
	}
	token, err := realtime.GetToken(user.ID)
	if err != nil {
		log.Printf("unable to retrieve ably token for user %q: %q", user.ID, err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}
	c.JSON(http.StatusOK, token)
}

func AblyTokenDemo(c *gin.Context) {
	capabilities, _ := json.Marshal(&map[string][]string{
		"*": {"subscribe", "history"},
	})
	token, err := realtime.CreateTokenRequest(&ably.TokenParams{
		Capability: string(capabilities),
		ClientID:   "nobody",
	})
	if err != nil {
		log.Printf("unable to retrieve demo ably token: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}
	c.JSON(http.StatusOK, token)
}

// UserUploadImage handles the uploading of a user profile image.
// @Summary      Upload Profile Image
// @Description  Uploads a profile image for the current user.
// @Tags         user,private
// @Accept       multipart/form-data
// @Produce      json
// @Param        image formData file true "Profile Image"
// @Success      200
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /user/image [put]
func UserUploadImage(c *gin.Context) {
	user := getUserFromContext(c)
	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("image upload error: %q", err)
		SendError(c, http.StatusBadRequest, Err400_InvalidRequestBody)
		return
	}

	// Restrict image size (example: 1MB)
	if file.Size > 1*1024*1024 {
		SendError(c, http.StatusBadRequest, Err400_FileTooLarge)
		return
	}

	// Open the file
	uploadedFile, err := file.Open()
	if err != nil {
		log.Printf("unable to open uploaded file: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}
	defer uploadedFile.Close()

	// Read the file into a byte slice
	fileBytes, err := io.ReadAll(uploadedFile)
	if err != nil {
		log.Printf("unable to read uploaded file: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}

	// Create an ImageData instance
	imageData := &userService.ImageData{
		ContentType:   file.Header.Get("Content-Type"),
		BinaryContent: fileBytes,
	}

	// Save the image data to the user's record using UserService
	us := getUserServiceFromContext(c)
	if err := us.UpdateUserProfileImage(user.ID, imageData); err != nil {
		log.Printf("unable to update user profile image: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}

	c.Status(http.StatusAccepted)
}

// UserGetImage handles the retrieval of a user profile image.
// @Summary      Get Profile Image
// @Description  Retrieves the profile image of a user.
// @Tags         user,public
// @Produce      image/*
// @Param        userId path string true "User ID"
// @Success      200  {file}  byte[]
// @Failure      404  {object}  ErrorResponse
// @Router       /user/{userId}/image [get]
func UserGetImage(c *gin.Context) {
	userId := c.Param("userId")
	us := getUserServiceFromContext(c)
	user, err := us.GetUserById(userId)
	if err != nil || user == nil {
		SendError(c, http.StatusNotFound, Err404_UserNotFound)
		return
	}
	imageData := us.GetUserImage(user)
	if imageData == nil {
		SendError(c, http.StatusNotFound, Err404_UserHasNoImage)
	} else {
		c.Data(http.StatusOK, imageData.ContentType, imageData.BinaryContent)
	}
}

// @Summary		Request new password (password forgotten flow)
// @Description	Allows users to recover their forgotten password by submitting associated email address
// @Tags			user,public,password
// @Accept		json
// @Produce		json
// @Param			request	body		userService.UserRequestNewPasswordDTO	true	"User's email address"
// @Success		202
// @Failure		400		{object}	ErrorResponse
// @Failure		424		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/user/request-new-password [post]
func UserRequestNewPassword(c *gin.Context) {
	var userRequestNewPasswordDTO userService.UserRequestNewPasswordDTO
	if err := c.ShouldBindJSON(&userRequestNewPasswordDTO); err != nil {
		errorFields := misc.ParseValidationError(err)
		log.Printf("unmet request body constraints: %q", errorFields.GetAllFields())
		SendError(c, http.StatusBadRequest, Err400_InvalidRequestBody)
		return
	}
	us := getUserServiceFromContext(c)
	user, err := us.GetUserByEmail(userRequestNewPasswordDTO.Email)
	if err != nil || user == nil {
		log.Printf("unable to locate user")
		SendError(c, http.StatusBadRequest, Err400_EmailNotRegistered)
		return
	}
	// send password reset email
	token, _ := generateToken(user, PasswordReset)
	var html bytes.Buffer
	emailService.PasswordReset(
		user.FullName,
		fmt.Sprintf(
			"%s/forms/password-reset?token=%s",
			os.Getenv("WEB_CLIENT_URL"),
			token.String(),
		),
		&html,
	)
	if err := email.Send(c.Request.Context(), email.EmailDTO{
		From:     "no-reply@quible.io",
		To:       user.Email,
		Subject:  "Password reset",
		HTMLBody: html.String(),
	}); err != nil {
		log.Printf("unable to send password reset email: %q", err)
		SendError(c, http.StatusFailedDependency, Err424_UnableToSendEmail)
		return
	}
	c.Status(http.StatusAccepted)
}

// @Summary		Resets user's password
// @Description	Validates provided JWT, password, and, if successful -- updates the password for the user identified by JWT
// @Tags			user,password
// @Accept		json
// @Produce		json
// @Param			request	body		userService.UserResetPasswordDTO	true	"Password and its confirmation"
// @Success		200
// @Failure		400		{object}	ErrorResponse
// @Failure		417		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/user/password-reset [post]
func UserPasswordReset(c *gin.Context) {
	var userResetPasswordDTO userService.UserResetPasswordDTO
	if err := c.ShouldBindJSON(&userResetPasswordDTO); err != nil {
		log.Printf("password validation failed: %q", err)
		SendError(c, http.StatusBadRequest, Err400_InvalidRequestBody)
		return
	}
	tokenClaims, err := verifyJWT(userResetPasswordDTO.Token, PasswordReset)
	if err != nil {
		log.Printf("invalid token: %q", err)
		SendError(c, http.StatusExpectationFailed, Err417_InvalidToken)
		return
	}
	userId := tokenClaims["userId"].(string)
	us := getUserServiceFromContext(c)
	user, err := us.GetUserById(userId)
	if err != nil || user == nil {
		log.Printf("unable to associate user: %q", err)
		SendError(c, http.StatusExpectationFailed, Err417_UnableToAssociateUser)
		return
	}
	// All checks passed
	user.HashedPassword, _ = us.HashPassword(userResetPasswordDTO.Password)
	if err := us.Update(user); err != nil {
		log.Printf("unable to update user with the new password: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnableToResetPassword)
		return
	}
	c.Status(http.StatusOK)
}

// -- Chat groups and channels

func CreateChatGroup(c *gin.Context) {
	var dto chatService.CreateChatGroupDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Println(err)
		SendError(c, http.StatusBadRequest, Err400_InvalidRequestBody)
		return
	}
	user := getUserFromContext(c)
	cs := chatService.ChatService{
		C: c.Request.Context(),
	}
	_, err := cs.CreateChatGroup(
		user,
		dto.Name,
		dto.Title,
		dto.Summary,
		dto.IsPrivate,
	)
	if err != nil {
		log.Println(err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}
	c.Status(http.StatusCreated)
}

func DeleteChatGroup(c *gin.Context) {
	id := c.Param("chatGroupId")
	user := getUserFromContext(c)
	cs := chatService.ChatService{
		C: c.Request.Context(),
	}
	if err := cs.DeleteChatGroup(user, id); err != nil {
		log.Println(err)
		if errors.Is(err, chatService.ErrChatGroupNotFound) {
			SendError(c, http.StatusNotFound, Err404_ChatGroupNotFound)
		} else {
			SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		}
		return
	}
	c.Status(http.StatusNoContent)
}
