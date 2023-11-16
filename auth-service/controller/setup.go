package controller

import (
	"github.com/gin-gonic/gin"
	c "gitlab.com/quible-backend/lib/controller"
)

// Add "health" endpoint at /health
var WithHealth = c.WithHealth

// Add "swagger" endpoint at /docs
var WithSwagger = c.WithSwagger

// Setup the controller and all handlers
func Setup(g *gin.RouterGroup, options ...c.Option) {
	// Process options
	for _, option := range options {
		option(g)
	}

}
