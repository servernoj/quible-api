package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/mail-service/postmark"
)

// @Summary		Send single email
// @Description	Use 3rd party (Postmark) mail delivery service to send an email
// @Tags			private
// @Accept		json
// @Produce		json
// @Param			request	body		postmark.EmailDTO	true	"Email delivery request"
// @Success		200		{object}	postmark.PostmarkResponse
// @Failure		400		{object}	ErrorResponse
// @Failure		401		{object}	ErrorResponse
// @Failure		424		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/send [post]
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
