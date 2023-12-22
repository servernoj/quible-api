package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/mail-service/postmark"
)

func SendEmailHandler(c *gin.Context) {

	type DTO = postmark.EmailDTO
	var NewMailer = postmark.NewMailer

	var email DTO
	if err := c.ShouldBindJSON(&email); err != nil {
		ErrorMap.SendError(c, http.StatusBadRequest, Err400_InvalidRequestBody)
		return
	}

	response, err := NewMailer(c.Request.Context()).SendEmail(email)

	if err != nil {
		ErrorMap.SendError(c, http.StatusFailedDependency, Err424_PostmarkSendEmail)
		return
	}

	c.JSON(http.StatusOK, *response)
}
