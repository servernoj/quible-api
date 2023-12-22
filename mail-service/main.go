package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/lib/env"
	"github.com/quible-io/quible-api/mail-service/controller"
)

const DefaultPort = 8003

// swagger.yaml
var swaggerSpec string

func main() {
	Server()
}

func Server() {
	// -- Environment vars from .env file
	env.Setup()
	// -- HTTP server
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())
	g := r.Group("/api/v1")
	controller.Setup(
		g,
		controller.WithSwagger(swaggerSpec),
		controller.WithHealth(),
	)
	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(DefaultPort)
	}
	log.Printf("starting server on port: %s\n", port)
	log.Fatalf("%v", r.Run(":"+port))
}
