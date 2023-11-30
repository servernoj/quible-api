package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const userIdContextKey = "userId"

// this function MUST abort the context
type Terminator func(c *gin.Context, fmt string, args ...any)

func InjectUserIdOrFail(terminator Terminator) gin.HandlerFunc {
	return func(c *gin.Context) {
		request, _ := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf(
				"%s/api/v1/user",
				os.Getenv("ENV_URL_AUTH_SERVICE"),
			),
			http.NoBody,
		)
		request.Header.Add("Authorization", c.Request.Header.Get("Authorization"))
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			terminator(c, "unable to send request to auth-service: %q", err)
			return
		}
		body := response.Body
		defer body.Close()
		var data map[string]any
		if err := json.NewDecoder(body).Decode(&data); err != nil {
			terminator(c, "unable to parse response from auth-service: %q", err)
			return
		}
		if response.StatusCode == http.StatusUnauthorized {
			c.AbortWithStatusJSON(response.StatusCode, data)
			return
		}
		if userId, ok := data["id"].(string); !ok {
			terminator(c, "field `id` is not present in the returned user object")
			return
		} else {
			c.Set(userIdContextKey, userId)
		}

		c.Next()
	}
}
