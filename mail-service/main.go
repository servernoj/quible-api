package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/lib/env"
	"gitlab.com/quible-backend/mail-service/controller"
	"gitlab.com/quible-backend/mail-service/service"
)

const DefaultPort = 8003

func main() {
	// set the env
	env.Setup()

	// create the client
	serverToken := os.Getenv("ENV_SERVER_TOKEN")
	accountToken := os.Getenv("ENV_ACCOUNT_TOKEN")
	client := &service.Client{
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
		ServerToken:  serverToken,
		AccountToken: accountToken,
		BaseURL:      "https://api.postmarkapp.com",
	}

	// start the gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.Default())

	// set the controller
	controller.SetupRoutes(router, client)

	// we could start other controller.SetupOtherRoutes(router) latter

	// start the service
	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(DefaultPort)
	}
	log.Printf("Starting mail service on port: %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}
