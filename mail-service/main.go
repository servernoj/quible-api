package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/lib/env"
	"gitlab.com/quible-backend/mail-service/controller"
)

const DefaultPort = 8003

func main() {
	// set the env
	env.Setup()

	// create the client
	client := controller.NewClient()

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
