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
	g.GET("/live/token", GetLiveToken)
	//-- Protected API
	protected := g.Group("", c.InjectUserIdOrFail(terminator))
	// BasketAPI
	protected.GET("/games", GetGames)
	protected.GET("/game", GetGameDetails)
	// Chat
	chatPublic := g.Group("/chat")
	chatProtected := protected.Group("/chat")
	chatProtected.POST("groups", CreateChatGroup)
	chatProtected.GET("groups", ListChatGroups)
	chatProtected.DELETE("groups/:chatGroupId", DeleteChatGroup)
	chatProtected.POST("channels", CreateChannel)
	chatProtected.GET("channels", GetMyGroupedChannels)
	chatProtected.GET("channels/list", GetMyChannels)
	chatProtected.POST("channels/:channelId", JoinPublicChannel)
	chatProtected.POST("channels/:channelId/invite", InviteToPrivateChannel)
	chatProtected.DELETE("channels/:channelId", LeaveChannel)
	chatProtected.GET("token", GetChatToken)
	// -- public
	chatPublic.GET("groups/search", SearchPublicChannelsByChatGroupTitle)
	chatPublic.POST("channels/accept", AcceptInvitationToPrivateChannel)
}
