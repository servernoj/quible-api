package controller

import (
	"github.com/gin-gonic/gin"
	c "github.com/quible-io/quible-api/lib/controller"
)

const userContextKey = "user"
const serviceContextKey = "service"

// Add "health" endpoint at /health
var WithHealth = c.WithHealth

// Add "swagger" endpoint at /docs
var WithSwagger = c.WithSwagger

// Add "email tester" endpoint at /email-tester

// Setup the controller and all handlers
func Setup(g *gin.RouterGroup, options ...c.Option) {
	for _, option := range options {
		option(g)
	}
	g.GET("docs/errors", GetErrorCodes)
	g.Use(injectUserService)
	// -- Public API
	g.POST("/user", UserRegister)
	g.POST("/user/activate", UserActivate)
	g.POST("/user/password-reset", UserPasswordReset)
	g.POST("/user/request-new-password", UserRequestNewPassword)
	g.POST("/user/refresh", UserRefresh)
	g.POST("/login", UserLogin)
	g.GET("/user/:userId/image", UserGetImage)
	g.GET("/rt/read-only-token", AblyTokenDemo)
	//-- Protected API
	protected := g.Group("", authMiddleware)
	protected.GET("/user", UserGet)
	protected.GET("/user/:userId/profile", UserGetById)
	protected.PATCH("/user", UserPatch)
	protected.GET("/rt/token", AblyToken)
	protected.PUT("/user/image", UserUploadImage)
	// -- chat service related
	protected.POST("chat/groups", CreateChatGroup)
	protected.DELETE("chat/groups/:chatGroupId", DeleteChatGroup)
	protected.GET("chat/capabilities", GetCapabilities)
}
