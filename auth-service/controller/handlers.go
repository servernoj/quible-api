package controller

import (
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"gitlab.com/quible-backend/auth-service/service"
	"gitlab.com/quible-backend/lib/misc"
)

var UserFields = []string{"id", "username", "email", "phone", "full_name"}

// @Summary		Register
// @Description	Register a new user.
// @Tags			user,public
// @Accept			json
// @Produce		json
// @Param			request	body		service.UserRegisterDTO	true	"User registration information"
// @Success		201		{object}	UserResponse
// @Failure		400		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router			/user [post]
func UserRegister(c *gin.Context) {
	userService := getUserServiceFromContext(c)
	var userRegisterDTO service.UserRegisterDTO
	if err := c.ShouldBindJSON(&userRegisterDTO); err != nil {
		log.Printf("invalid request body: %q", err)
		SendError(c, http.StatusBadRequest, Err400_InvalidRequestBody)
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
// @Accept			json
// @Produce		json
// @Param			request	body		service.UserLoginDTO	true	"User login credentials"
// @Success		200		{object}	TokenResponse
// @Failure		400		{object}	ErrorResponse
// @Failure		401		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router			/login [post]
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
	c.JSON(http.StatusOK, gin.H{
		"token": generateToken(foundUser),
	})
}

// @Summary		Get user
// @Description	Returns user profile associated with the token
// @Tags			user,private
// @Produce		json
// @Success		200	{object}	UserResponse
// @Failure		401	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router			/user [get]
func UserGet(c *gin.Context) {
	user := getUserFromContext(c)
	c.JSON(
		http.StatusCreated,
		misc.PickFields(user, UserFields...),
	)
}

// @Summary		Update user
// @Description	Updates user profile associated with the token
// @Tags			user,private
// @Accept			json
// @Produce		json
// @Param			request	body		service.UserPatchDTO	true	"Partial user object to be used for update"
// @Success		200		{object}	UserResponse
// @Failure		400		{object}	ErrorResponse
// @Failure		401		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router			/user [patch]
func UserPatch(c *gin.Context) {
	var userPatchDTO service.UserPatchDTO
	if err := c.ShouldBindJSON(&userPatchDTO); err != nil {
		errorCode := Err400_InvalidRequestBody
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
		log.Printf("unmet request body contraints: %q", errorFields.GetAllFields())
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
