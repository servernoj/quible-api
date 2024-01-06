package controller

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/quible-io/quible-api/auth-service/realtime"
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
			var host string
			switch os.Getenv("ENV_DEPLOYMENT") {
			case "dev":
				host = "https://auth.dev.quible.io"
			case "prod":
				host = "https://auth.prod.quible.io"
			default:
				host = os.Getenv("ENV_URL_AUTH_SERVICE")
			}
			var html bytes.Buffer
			emailService.UserActivation(
				user.FullName,
				fmt.Sprintf(
					"%s/api/v1/user/activate?token=%s",
					host,
					token.String(),
				),
				&html,
			)
			return email.Send(c.Request.Context(), email.EmailDTO{
				From:     "no-reply@quible.tech",
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
// @Description	Handles click from activation email
// @Tags			user,public
// @Produce		text/plain
// @Param			token	query		string	true	"JWT generated during registration"
// @Success		200	{string}	string
// @Router		/user/activate [get]
func UserActivate(c *gin.Context) {
	us := getUserServiceFromContext(c)
	token := c.Request.URL.Query().Get("token")
	tokenClaims, err := verifyJWT(token, Activate)
	if err != nil {
		log.Printf("unable to verify token: %q", err)
		c.String(http.StatusExpectationFailed, "Unable to verify the request")
		return
	}
	userId := tokenClaims["userId"].(string)
	if err := us.ActivateUser(userId); err != nil {
		log.Printf("unable to activate user: %q", err)
		c.String(http.StatusInternalServerError, "Unable to activate user account")
		return
	}
	c.String(http.StatusOK, "Account has been successfully activated")
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
		log.Printf("provided refresh token has been revoked")
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
// @Success		202		{string}	string
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
	var host string
	switch os.Getenv("ENV_DEPLOYMENT") {
	case "dev":
		host = "https://auth.dev.quible.io"
	case "prod":
		host = "https://auth.prod.quible.io"
	default:
		host = os.Getenv("ENV_URL_AUTH_SERVICE")
	}
	var html bytes.Buffer
	emailService.PasswordReset(
		user.FullName,
		fmt.Sprintf(
			"%s/api/v1/password-reset?token=%s",
			host,
			token.String(),
		),
		&html,
	)
	if err := email.Send(c.Request.Context(), email.EmailDTO{
		From:     "no-reply@quible.tech",
		To:       user.Email,
		Subject:  "Password reset",
		HTMLBody: html.String(),
	}); err != nil {
		log.Printf("unable to send password reset email: %q", err)
		SendError(c, http.StatusFailedDependency, Err424_UnableToSendEmail)
		return
	}
	c.String(http.StatusAccepted, "Password reset request accepted")
}

// @Summary		Render password reset form
// @Description	Render password reset form in response to click on a link from email
// @Tags			user,password
// @Produce		text/html
// @Param			token	query		string	true	"JWT generated while handling password reset request"
// @Success		200	{string}	string
// @Success		417	{string}	string
// @Router		/password-reset [get]
func UserPasswordResetForm(c *gin.Context) {
	us := getUserServiceFromContext(c)
	token := c.Request.URL.Query().Get("token")
	tokenClaims, err := verifyJWT(token, PasswordReset)
	if err != nil {
		log.Printf("unable to verify token: %q", err)
		c.HTML(
			http.StatusExpectationFailed,
			"password.html",
			gin.H{
				"error": "Invalid request",
			},
		)
		return
	}
	userId := tokenClaims["userId"].(string)
	user, err := us.GetUserById(userId)
	if err != nil || user == nil {
		log.Printf("unable to locate requested user: %q", err)
		c.HTML(
			http.StatusExpectationFailed,
			"password.html",
			gin.H{
				"error": "Account not found",
			},
		)
		return
	}
	c.HTML(http.StatusOK, "password.html", gin.H{
		"error": nil,
	})
}

// @Summary		Accepts new password from the rendered web form
// @Description	Validates token provided in query and performs validation of password field, if successful -- updates the password for user identified from JWT
// @Tags			user,password
// @Accept		application/x-www-form-urlencoded
// @Produce		text/html
// @Param			request	body		userService.UserResetPasswordDTO	true	"Password and its confirmation"
// @Param			token	query		string	true	"JWT generated while handling password reset request"
// @Success		200		{string}	string
// @Failure		400		{object}	ErrorResponse
// @Failure		417		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/password-reset [post]
func UserPasswordResetAction(c *gin.Context) {
	us := getUserServiceFromContext(c)
	token := c.Request.URL.Query().Get("token")
	tokenClaims, err := verifyJWT(token, PasswordReset)
	if err != nil {
		log.Printf("unable to verify token: %q", err)
		c.HTML(
			http.StatusExpectationFailed,
			"password.html",
			gin.H{
				"error": "Invalid request",
			},
		)
		return
	}
	userId := tokenClaims["userId"].(string)
	user, err := us.GetUserById(userId)
	if err != nil || user == nil {
		log.Printf("unable to locate requested user: %q", err)
		c.HTML(
			http.StatusExpectationFailed,
			"password.html",
			gin.H{
				"error": "Account not found",
			},
		)
		return
	}
	var userResetPasswordDTO userService.UserResetPasswordDTO
	if err := c.ShouldBindWith(&userResetPasswordDTO, binding.FormPost); err != nil {
		log.Printf("password validation failed: %q", err)
		c.HTML(
			http.StatusExpectationFailed,
			"password.html",
			gin.H{
				"error": "Password(s) don't match or have insufficient complexity",
			},
		)
		return
	}
	// All checks passed
	user.HashedPassword, _ = us.HashPassword(userResetPasswordDTO.Password)
	if err := us.Update(user); err != nil {
		log.Printf("unable to update user with the new password: %q", err)
		c.String(
			http.StatusInternalServerError,
			"Unable to apply store new password",
		)
		return
	}
	c.String(http.StatusOK, "Password has been successfully reset")
}
