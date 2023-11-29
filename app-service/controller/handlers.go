package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/quible-backend/lib/misc"
	"gitlab.com/quible-backend/lib/models"
)

var UserFields = []string{"id", "username", "email", "phone", "full_name"}

func Test(c *gin.Context) {
	if c.IsAborted() {
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}
	userId, exists := c.Get("userId")
	if !exists {
		SendError(c, http.StatusUnauthorized, Err401_UserIdNotFound)
		return
	}
	user, _ := models.FindUserG(c.Request.Context(), userId.(string))
	if user == nil {
		log.Printf("unable to find user with userId %q", userId)
		SendError(c, http.StatusUnauthorized, Err401_UserNotFound)
		return
	}
	c.JSON(
		http.StatusOK,
		misc.PickFields(user, UserFields...),
	)
}
