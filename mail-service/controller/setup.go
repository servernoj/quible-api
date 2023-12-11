package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	a "github.com/quible-io/quible-api/app-service/controller"
	c "github.com/quible-io/quible-api/lib/controller"
	"gitlab.com/quible-backend/mail-service/service"
)

var (
	WithSwagger = c.WithSwagger
	WithHealth  = c.WithHealth
)

// terminator for "protected" group
func terminator(c *gin.Context, fmt string, args ...any) {
	log.Printf(fmt, args...)
	a.ErrorMap.SendError(c, http.StatusInternalServerError, a.Err500_UnknownError)
}

// setup.go

func Setup(g *gin.RouterGroup, client *service.Client, options ...c.Option) {
	// Apply additional options to the router
	for _, option := range options {
		option(g)
	}

	// Create a protected router group
	protected := g.Group("", c.InjectUserIdOrFail(terminator))

	// Setup protected routes
	protected.POST("/send-email", func(c *gin.Context) { SendEmailHandler(c, client) })
}
