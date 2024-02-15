package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/lib/models"
)

const userIdContextKey string = "userId"

func getUserFromContext(c *gin.Context) *models.User {
	userId, ok := c.Get(userIdContextKey)
	if !ok {
		log.Printf("no userId in context")
		return nil
	}
	user, err := models.FindUserG(c.Request.Context(), userId.(string))
	if err != nil || user == nil {
		log.Printf("unable to retrieve user object from DB")
		SendError(c, http.StatusUnauthorized, Err401_UserNotFound)
		return nil
	}
	return user
}
