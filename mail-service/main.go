package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/quible-backend/mail-service/controller"
	"github.com/quible-io/quible-api/lib/env"
)

const DefaultPort = 8003

// swagger.yaml
var swaggerSpec string

func main() {
	// Set the environment variables
	env.Setup()

	// Start the Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())
	g := r.Group("/api/v1")

	// Set up the controller with the protected group and client
	controller.Setup(g, controller.WithSwagger(swaggerSpec),
		controller.WithHealth() /*, other necessary options if any*/)

	// Start the service on the specified port
	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(DefaultPort)
	}
	log.Printf("Starting mail service on port: %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}
