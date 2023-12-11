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

// swagger.yaml
var swaggerSpec string

func main() {
	// Set the environment variables
	env.Setup()

	// Create the client
	client := controller.NewClient()

	// Start the Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.Default())

	// Create a router group for the protected routes
	protectedGroup := router.Group("/")

	// Set up the controller with the protected group and client
	controller.Setup(protectedGroup, client, controller.WithSwagger(swaggerSpec),
		controller.WithHealth() /*, other necessary options if any*/)

	// Start the service on the specified port
	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(DefaultPort)
	}
	log.Printf("Starting mail service on port: %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}
