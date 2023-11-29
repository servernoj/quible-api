package controller

import (
	"github.com/gin-gonic/gin"
	c "gitlab.com/quible-backend/lib/controller"
)

// Setup the controller and all handlers
func Setup(g *gin.RouterGroup, options ...c.Option) {
	for _, option := range options {
		option(g)
	}
	g.GET("docs/errors", GetErrorCodes)
	// -- Public API
	//-- Protected API
	protected := g.Group("", c.InjectUserId)
	protected.GET("/test", Test)
}
