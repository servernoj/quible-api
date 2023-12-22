package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/mail-service/postmark"
)

func SendEmailHandler(c *gin.Context) {

	type DTO = postmark.EmailDTO
	var NewClient = postmark.NewClient

	var email DTO
	if err := c.ShouldBindJSON(&email); err != nil {
		log.Printf("request parsing error %s:", err)
		ErrorMap.SendError(c, http.StatusBadRequest, Err400_InvalidRequestBody)
		return
	}

	response, err := NewClient(c.Request.Context()).SendEmail(email)

	if err != nil {
		log.Printf("email sender error: %q", err)
		ErrorMap.SendError(c, http.StatusFailedDependency, Err424_PostmarkSendEmail)
		return
	}

	c.JSON(http.StatusOK, *response)
}
