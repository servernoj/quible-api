package controller

import (
	"github.com/gin-gonic/gin"
	c "gitlab.com/quible-backend/lib/controller"
)

const userContextKey = "user"
const serviceContextKey = "service"

// Add "health" endpoint at /health
var WithHealth = c.WithHealth

// Add "swagger" endpoint at /docs
var WithSwagger = c.WithSwagger

// Setup the controller and all handlers
func Setup(g *gin.RouterGroup, options ...c.Option) {
	for _, option := range options {
		option(g)
	}
	g.GET("docs/errors", GetErrorCodes)
	g.Use(injectUserService)
	// -- Public API
	g.POST("/user", UserRegister)
	g.POST("/user/refresh", UserRefresh)
	g.POST("/login", UserLogin)
	//-- Protected API
	protected := g.Group("", authMiddleware)
	protected.GET("/user", UserGet)
	protected.PATCH("/user", UserPatch)

}
