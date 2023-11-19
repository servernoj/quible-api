package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/quible-backend/auth-service/service"
	"gitlab.com/quible-backend/lib/misc"
)

var UserFields = []string{"id", "username", "email", "phone", "full_name"}

// @Description	Register a new user.
// @Tags			user,public
// @Accept			json
// @Produce		json
// @Param			request	body		service.UserRegisterDTO	true	"User registration information"
// @Success		201		{object}	UserResponse
// @Failure		400		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router			/register [post]
func UserRegister(c *gin.Context) {
	userService := getUserServiceFromContext(c)
	var userRegisterDTO service.UserRegisterDTO
	if err := c.ShouldBindJSON(&userRegisterDTO); err != nil {
		log.Printf("invalid request body: %q", err)
		SendError(c, http.StatusBadRequest, 100)
		return
	}

	if foundUser, _ := userService.GetUserByEmail(userRegisterDTO.Email); foundUser != nil {
		SendError(c, http.StatusBadRequest, 8)
		return
	}
	if foundUser, _ := userService.GetUserByUsername(userRegisterDTO.Username); foundUser != nil {
		SendError(c, http.StatusBadRequest, 5)
		return
	}
	createdUser, err := userService.CreateUser(&userRegisterDTO)
	if err != nil {
		log.Printf("unable to register user: %q", err)
		SendError(c, http.StatusInternalServerError, 4)
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
		SendError(c, http.StatusUnauthorized, 1)
		return
	}
	foundUser, _ := userService.GetUserByEmail(userLoginDTO.Email)
	if foundUser == nil {
		log.Printf("user with given email not found: %q", userLoginDTO.Email)
		SendError(c, http.StatusUnauthorized, 1)
		return
	}
	if err := userService.ValidatePassword(foundUser.HashedPassword, userLoginDTO.Password); err != nil {
		log.Printf("invalid password: %+v", userLoginDTO)
		SendError(c, http.StatusUnauthorized, 2)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": generateToken(foundUser),
	})
}

// @Summary		Get user
// @Description	Get user profile associated with token
// @Tags			user,private
// @Produce		json
// @Success		200	{object}	UserResponse
// @Failure		401	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router			/user [get]
func GetUser(c *gin.Context) {
	user := getUserFromContext(c)
	c.JSON(
		http.StatusCreated,
		misc.PickFields(user, UserFields...),
	)
}
