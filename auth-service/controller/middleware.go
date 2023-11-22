package controller

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/quible-backend/auth-service/service"
	"gitlab.com/quible-backend/lib/models"
)

func getUserServiceFromContext(c *gin.Context) *service.UserService {
	serviceCandidate, ok := c.Get(serviceContextKey)
	if !ok {
		return nil
	}
	userService, ok := serviceCandidate.(*service.UserService)
	if !ok {
		return nil
	}
	return userService
}

func getUserFromContext(c *gin.Context) *models.User {
	userCandidate, ok := c.Get(userContextKey)
	if !ok {
		log.Printf("no user object in context")
		return nil
	}
	user, ok := userCandidate.(*models.User)
	if !ok {
		log.Printf("unable to assert user type")
		return nil
	}
	return user
}

func authMiddleware(c *gin.Context) {
	authToken := strings.TrimSpace(c.GetHeader("Authorization"))
	userService := getUserServiceFromContext(c)
	if userService == nil {
		log.Printf("unable to retrieve user service")
		SendError(c, http.StatusInternalServerError, Err500_UnknownError)
		return
	}
	if authToken == "" {
		SendError(c, http.StatusUnauthorized, Err401_AuthorizationHeaderMissing)
		return
	}
	re, _ := regexp.Compile(`\s+`)
	headerParts := re.Split(authToken, -1)
	if len(headerParts) != 2 {
		log.Printf("authorization header format is invalid, missing space")
		SendError(c, http.StatusUnauthorized, Err401_AuthorizationHeaderInvalid)
		return
	}
	if headerParts[0] != "Bearer" {
		log.Printf("authorization header doesn't carry bearer token")
		SendError(c, http.StatusUnauthorized, Err401_AuthorizationHeaderInvalid)
		return
	}
	token := headerParts[1]
	id, err := verifyJWT(token)
	if err != nil {
		log.Printf("unable to verify token %q: %s", token, err)
		SendError(c, http.StatusUnauthorized, Err401_AuthorizationHeaderInvalid)
		return
	}
	user, err := userService.GetUserById(id)
	if err != nil || user == nil {
		log.Printf("user with id = %q not found", id)
		SendError(c, http.StatusUnauthorized, Err401_UserNotFound)
		return
	}
	c.Set(userContextKey, user)
	c.Next()
}

// Inject user service object for all requests to use
func injectUserService(c *gin.Context) {
	userService := service.UserService{
		// TODO: possibly need to send a context dettached context instead of the original one
		C: c.Request.Context(),
	}
	c.Set(serviceContextKey, &userService)
	c.Next()
}