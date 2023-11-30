package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	c "github.com/quible-io/quible-api/lib/controller"
)

var (
	WithSwagger = c.WithSwagger
	WithHealth  = c.WithHealth
)

// terminator for "protected" group
func terminator(c *gin.Context, fmt string, args ...any) {
	log.Printf(fmt, args...)
	ErrorMap.SendError(c, http.StatusInternalServerError, Err500_UnknownError)
}

// Setup the controller and all handlers
func Setup(g *gin.RouterGroup, options ...c.Option) {
	for _, option := range options {
		option(g)
	}

	g.GET("docs/errors", ErrorMap.GetErrorCodes)
	// -- Public API
	//-- Protected API
	protected := g.Group("", c.InjectUserIdOrFail(terminator))
	protected.GET("/schedule-season", ScheduleSeason)
	protected.GET("/daily-schedule", DailySchedule)
	protected.GET("/team-info", TeamInfo)
}
