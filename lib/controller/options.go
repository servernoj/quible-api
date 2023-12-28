package controller

import (
	"bytes"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/auth-service/services/emailService"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/quible-io/quible-api/lib/email/postmark"
	"github.com/quible-io/quible-api/lib/swagger"
)

type Option func(g *gin.RouterGroup)

func WithHealth() Option {
	return func(g *gin.RouterGroup) {
		g.GET("/health", func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		})
	}
}

func WithEmailTester() Option {
	return func(g *gin.RouterGroup) {
		g.POST("/email-tester", func(c *gin.Context) {
			type TestEmailDTO struct {
				Handler string  `json:"handler" binding:"required"`
				Args    []any   `json:"args" binding:"required"`
				Subject *string `json:"subject"`
				To      *string `json:"to"`
			}
			var html bytes.Buffer
			var testEmailDTO TestEmailDTO
			if err := c.ShouldBindJSON(&testEmailDTO); err != nil {
				c.String(http.StatusBadRequest, "%v", err)
				return
			}
			if _, ok := emailService.Handlers[testEmailDTO.Handler]; !ok {
				keys := []string{}
				for key := range emailService.Handlers {
					keys = append(keys, key)
				}
				c.String(http.StatusBadRequest, "invalid handler, supposed to be one of %q", keys)
				return
			}
			fn := reflect.ValueOf(emailService.Handlers[testEmailDTO.Handler])
			// -- args
			args := make([]reflect.Value, len(testEmailDTO.Args))
			for idx, param := range testEmailDTO.Args {
				if fn.Type().In(idx) != reflect.TypeOf(param) {
					c.String(
						http.StatusBadRequest,
						"invalid type of Args[%d]: %q. Got %q, expected %q",
						idx,
						param,
						reflect.TypeOf(param),
						fn.Type().In(idx),
					)
					return
				}
				args[idx] = reflect.ValueOf(param)
			}
			args = append(args, reflect.ValueOf(&html))

			if fn.Type().NumIn() != len(args) {
				c.String(http.StatusBadRequest, "invalid length of args array")
				return
			}
			fn.Call(args)

			if c.Request.URL.Query().Has("debug") {
				c.Data(http.StatusOK, gin.MIMEHTML, html.Bytes())
				return
			}
			if testEmailDTO.To == nil {
				var temp = "contact@quible.tech"
				testEmailDTO.To = &temp
			}
			if testEmailDTO.Subject == nil {
				var temp = "Test"
				testEmailDTO.Subject = &temp
			}

			if err := email.Send(c.Request.Context(), postmark.EmailDTO{
				From:     "no-reply@quible.tech",
				To:       *testEmailDTO.To,
				Subject:  *testEmailDTO.Subject,
				HTMLBody: html.String(),
			}); err != nil {
				c.String(http.StatusFailedDependency, "unable to send activation email: %q", err)
				return
			}
			c.Status(http.StatusAccepted)
		})
	}
}

func WithSwagger(spec string) Option {
	return func(g *gin.RouterGroup) {
		swagger.Register(g, spec, "/docs")
	}
}
