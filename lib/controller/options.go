package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func WithSwagger(spec string) Option {
	return func(g *gin.RouterGroup) {
		swagger.Register(g, spec, "/docs")
	}
}
