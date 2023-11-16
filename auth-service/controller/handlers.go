package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/quible-backend/auth-service/service"
	"gitlab.com/quible-backend/lib/misc"
)

func RegisterUser(c *gin.Context) {
	var user *service.UserRegisterDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		SendError(c, http.StatusBadRequest, 100)
		return
	}
	userService := getUserService(c)

	log.Println(userService, user)

	if foundUser, _ := userService.GetUserByEmail(user.Email); foundUser != nil {
		SendError(c, http.StatusBadRequest, 8)
		return
	}
	if foundUser, _ := userService.GetUserByUsername(user.Username); foundUser != nil {
		SendError(c, http.StatusBadRequest, 5)
		return
	}
	createdUser, err := userService.CreateUser(user)
	if err != nil {
		log.Printf("unable to register userL %v", err)
		SendError(c, http.StatusInternalServerError, 4)
		return
	}
	c.JSON(
		http.StatusCreated,
		misc.Pick(createdUser, "id", "username", "email", "phone", "full_name"),
	)
}
