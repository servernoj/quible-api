package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const userIdContextKey = "userId"

func InjectUserId(c *gin.Context) {
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
		log.Printf("unable to send request to auth-service: %q", err)
		c.Abort()
		return
	}
	body := response.Body
	defer body.Close()
	var data map[string]any
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		log.Printf("unable to parse response from auth-service: %q", err)
		c.Abort()
		return
	}
	if response.StatusCode == http.StatusUnauthorized {
		c.AbortWithStatusJSON(response.StatusCode, data)
		return
	}
	if userId, ok := data["id"].(string); ok {
		c.Set(userIdContextKey, userId)
	}

	c.Next()
}
