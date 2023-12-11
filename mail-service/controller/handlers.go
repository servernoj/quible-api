package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/quible-backend/mail-service/service"
)

func SendEmailHandler(c *gin.Context, client *service.Client) {
	var email service.Email
	if err := c.BindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	response, err := client.SendEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error sending email: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}
