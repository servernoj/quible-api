package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/auth-service/email"
	"github.com/quible-io/quible-api/auth-service/realtime"
	"github.com/quible-io/quible-api/auth-service/service"
	"github.com/quible-io/quible-api/lib/misc"
)

var UserFields = []string{"id", "username", "email", "phone", "full_name"}
var PublicUserFields = []string{"id", "full_name"}

// @Summary		Register
// @Description	Register a new user.
// @Tags			user,public
// @Accept		json
// @Produce		json
// @Param			request	body		service.UserRegisterDTO	true	"User registration information"
// @Success		201		{object}	UserResponse
// @Failure		400		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/user [post]
func UserRegister(c *gin.Context) {
	userService := getUserServiceFromContext(c)
	var userRegisterDTO service.UserRegisterDTO
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

	if foundUser, _ := userService.GetUserByEmail(userRegisterDTO.Email); foundUser != nil {
		SendError(c, http.StatusBadRequest, Err400_UserWithEmailExists)
		return
	}
	if foundUser, _ := userService.GetUserByUsername(userRegisterDTO.Username); foundUser != nil {
		SendError(c, http.StatusBadRequest, Err400_UserWithUsernameExists)
		return
	}
	createdUser, err := userService.CreateUser(&userRegisterDTO)
	if err != nil {
		log.Printf("unable to register user: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnableToRegister)
		return
	}
	c.JSON(
		http.StatusCreated,
		misc.PickFields(createdUser, UserFields...),
	)
}

// @Summary		Login
// @Description	Login with user credentials to get token
// @Tags			user,public
// @Accept		json
// @Produce		json
// @Param			request	body		service.UserLoginDTO	true	"User login credentials"
// @Success		200		{object}	TokenResponse
// @Failure		400		{object}	ErrorResponse
// @Failure		401		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/login [post]
func UserLogin(c *gin.Context) {
	userService := getUserServiceFromContext(c)

	var userLoginDTO service.UserLoginDTO
	if err := c.ShouldBindJSON(&userLoginDTO); err != nil {
		log.Printf("invalid request body: %q", err)
		SendError(c, http.StatusBadRequest, Err400_InvalidRequestBody)
		return
	}
	foundUser, _ := userService.GetUserByEmail(userLoginDTO.Email)
	if foundUser == nil {
		log.Printf("user with given email not found: %q", userLoginDTO.Email)
		SendError(c, http.StatusUnauthorized, Err401_InvalidCredentials)
		return
	}
	if err := userService.ValidatePassword(foundUser.HashedPassword, userLoginDTO.Password); err != nil {
		log.Printf("invalid password: %+v", userLoginDTO)
		SendError(c, http.StatusUnauthorized, Err401_InvalidCredentials)
		return
	}
	generatedAccessToken, err := generateToken(foundUser, false)
	if err != nil {
		log.Printf("unable to generate access token: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnableToGenerateToken)
		return
	}
	generatedRefreshToken, err := generateToken(foundUser, true)
	if err != nil {
		log.Printf("unable to generate refresh token: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnableToGenerateToken)
		return
	}
	foundUser.Refresh = generatedRefreshToken.ID
	if err := userService.Update(foundUser); err != nil {
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
	userService := getUserServiceFromContext(c)
	user, err := userService.GetUserById(userId)
	if err != nil || user == nil {
		log.Printf("user not found: %q", userId)
		SendError(c, http.StatusNotFound, Err404_UserNotFound)
		return
	}
	imageData := userService.GetUserImage(user)
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
// @Param			request	body		service.UserPatchDTO	true	"Partial user object to be used for update"
// @Success		200		{object}	UserResponse
// @Failure		400		{object}	ErrorResponse
// @Failure		401		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/user [patch]
func UserPatch(c *gin.Context) {
	var userPatchDTO service.UserPatchDTO
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
	userService := getUserServiceFromContext(c)
	if err := userService.Update(user); err != nil {
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
// @Param			request	body		service.UserRefreshDTO	true	"User's refresh token"
// @Success		200		{object}	TokenResponse
// @Failure		400		{object}	ErrorResponse
// @Failure		401		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/user/refresh [post]
func UserRefresh(c *gin.Context) {
	var userRefreshDTO service.UserRefreshDTO
	if err := c.ShouldBindJSON(&userRefreshDTO); err != nil {
		errorFields := misc.ParseValidationError(err)
		log.Printf("unmet request body constraints: %q", errorFields.GetAllFields())
		SendError(c, http.StatusBadRequest, Err400_InvalidRequestBody)
		return
	}

	claims, err := verifyJWT(userRefreshDTO.RefreshToken, true)
	if err != nil {
		log.Printf("invalid refresh token: %q", err)
		SendError(c, http.StatusUnauthorized, Err401_InvalidRefreshToken)
		return
	}

	userId := claims["userId"].(string)
	userService := getUserServiceFromContext(c)
	user, err := userService.GetUserById(userId)
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

	generatedAccessToken, err := generateToken(user, false)
	if err != nil {
		log.Printf("unable to generate access token: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnableToGenerateToken)
		return
	}
	generatedRefreshToken, err := generateToken(user, true)
	if err != nil {
		log.Printf("unable to generate refresh token: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnableToGenerateToken)
		return
	}
	user.Refresh = generatedRefreshToken.ID
	if err := userService.Update(user); err != nil {
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
	imageData := &service.ImageData{
		ContentType:   file.Header.Get("Content-Type"),
		BinaryContent: fileBytes,
	}

	// Save the image data to the user's record using UserService
	userService := getUserServiceFromContext(c)
	if err := userService.UpdateUserProfileImage(user.ID, imageData); err != nil {
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
	userService := getUserServiceFromContext(c)
	user, err := userService.GetUserById(userId)
	if err != nil || user == nil {
		SendError(c, http.StatusNotFound, Err404_UserNotFound)
		return
	}
	imageData := userService.GetUserImage(user)
	if imageData == nil {
		SendError(c, http.StatusNotFound, Err404_UserHasNoImage)
	} else {
		c.Data(http.StatusOK, imageData.ContentType, imageData.BinaryContent)
	}
}

func Activation(c *gin.Context) {
	var html bytes.Buffer
	email.ActivationRender(&html)
	requestBody := map[string]any{
		"from":     "contact@quible.tech",
		"to":       "simonbaev@gmail.com",
		"subject":  "User activation",
		"HtmlBody": html.String(),
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("unable to marshal request body: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
	}
	res, err := http.DefaultClient.Post(
		fmt.Sprintf(
			"%s/api/v1/send",
			os.Getenv("ENV_URL_MAIL_SERVICE"),
		),
		gin.MIMEJSON,
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		log.Printf("unable to send email: %q", err)
		SendError(c, http.StatusFailedDependency, Err424_UnknownError)
	}
	defer res.Body.Close()
	var body map[string]any
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		log.Printf("unable to decode response from mail-service: %q", err)
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
	}
	c.JSON(http.StatusOK, body)
}
