package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	c "github.com/quible-io/quible-api/lib/controller"
	"gitlab.com/quible-backend/mail-service/service"
)

var (
	WithSwagger = c.WithSwagger
	WithHealth  = c.WithHealth
)

func SetupRoutes(router *gin.Engine, client *service.Client) {
	router.POST("/send-email", func(c *gin.Context) {
		var email service.Email
		if err := c.BindJSON(&email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		response, err := client.SendEmail(context.Background(), email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error sending email: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"response": response})
	})

	// We could add more routes and use the same structure as app-service/setup
}
