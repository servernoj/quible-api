package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/quible-io/quible-api/lib/email/postmark"
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

	type EmailDTO = postmark.EmailDTO

	var emailDTO EmailDTO
	if err := c.ShouldBindJSON(&emailDTO); err != nil {
		log.Printf("request parsing error %s:", err)
		ErrorMap.SendError(c, http.StatusBadRequest, Err400_InvalidRequestBody)
		return
	}

	if err := email.Send(c.Request.Context(), emailDTO); err != nil {
		log.Printf("email sender error: %q", err)
		ErrorMap.SendError(c, http.StatusFailedDependency, Err424_PostmarkSendEmail)
		return
	}

	c.Status(http.StatusOK)
}
