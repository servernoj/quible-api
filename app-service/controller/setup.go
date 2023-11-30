package controller

import (
	"github.com/gin-gonic/gin"
	c "github.com/quible-io/quible-api/lib/controller"
)

var (
	WithSwagger = c.WithSwagger
	WithHealth  = c.WithHealth
)

// Setup the controller and all handlers
func Setup(g *gin.RouterGroup, options ...c.Option) {
	for _, option := range options {
		option(g)
	}
	g.GET("docs/errors", GetErrorCodes)
	// -- Public API
	//-- Protected API
	protected := g.Group("", c.InjectUserIdOrFail)
	protected.GET("/schedule-season", ScheduleSeason)
	protected.GET("/daily-schedule", DailySchedule)
}
