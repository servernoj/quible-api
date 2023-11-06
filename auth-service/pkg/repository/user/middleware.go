package user

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/quible-backend/auth-service/domain"
)

const userContextKey = "user"

func (h *controller) authenticateUser(c *gin.Context) (*domain.UserResponse, error) {
	authToken := c.GetHeader("Authorization")

	if authToken == "" {
		return nil, errors.New("Authorization header missing")
	}
	headerParts := strings.Split(authToken, " ")
	if len(headerParts) != 2 {
		return nil, errors.New("Authorization header is invalid")
	}
	if headerParts[0] != "Bearer" {
		return nil, errors.New("Authorization header is missing bearer part")
	}

	id, _, err := verifyJWT(headerParts[1])
	if err != nil {
		return nil, errors.New("Failed to verify auth Token.")
	}

	user, err := h.Service.Gets(id)
	if err != nil {
		return nil, errors.New("Failed to get user by email")
	}

	if user == nil {
		return nil, errors.New("User not found")
	}

	return user, nil
}

func (h *controller) AuthMiddleware(c *gin.Context) {
	user, err := h.authenticateUser(c)
	if err != nil {
		h.SendErrorAbort(c, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	c.Set(userContextKey, user)

	c.Next()
}

func (h *controller) currentUser(ctx *gin.Context) *domain.UserResponse {
	user, ok := ctx.Get(userContextKey)
	if !ok {
		// this shouldn't be happen because auth middleware is aborting when there is an Error.
		// this only possible if we try to get current user on non-auth endpoint.
		return nil
	}
	return user.(*domain.UserResponse)
}
