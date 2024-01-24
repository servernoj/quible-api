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
	SendError(c, http.StatusInternalServerError, Err500_UnknownError)
}

// Setup the controller and all handlers
func Setup(g *gin.RouterGroup, options ...c.Option) {
	for _, option := range options {
		option(g)
	}

	g.GET("docs/errors", ErrorMap.GetErrorCodes)
	// -- Public API
	g.POST("/live-push", LivePush)
	//-- Protected API
	protected := g.Group("", c.InjectUserIdOrFail(terminator))
	// RSC
	protected.GET("/schedule-season", ScheduleSeason)
	protected.GET("/daily-schedule", DailySchedule)
	protected.GET("/team-info", TeamInfo)
	protected.GET("/team-stats", TeamStats)
	protected.GET("/player-info", PlayerInfo)
	protected.GET("/player-stats", PlayerStats)
	protected.GET("/injuries", Injuries)
	protected.GET("/live", LiveFeed)
	// BasketAPI
	protected.GET("/games", GetGames)
	protected.GET("/game", GetGameDetails)
	// Chat
	protected.POST("chat/groups", CreateChatGroup)
	protected.GET("chat/groups", ListChatGroups)
	protected.GET("chat/groups/search", SearchPublicChannelsByChatGroupTitle)
	protected.DELETE("chat/groups/:chatGroupId", DeleteChatGroup)
	protected.POST("chat/channels", CreateChannel)
	protected.POST("chat/channels/:channelId", JoinPublicChannel)
	protected.DELETE("chat/channels/:channelId", LeaveChannel)
	protected.GET("chat/capabilities", GetCapabilities)
}
