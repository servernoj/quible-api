package main

import (
	_ "embed"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/quible-backend/auth-service/controller"
	"gitlab.com/quible-backend/lib/env"
	"gitlab.com/quible-backend/lib/store"
)

//	@title			Quible auth-service
//	@description	Authentication and authorization service of Quible.io
//	@version		0.1
//	@host			www.quible.io
//	@BasePath		/api/v1

const DefaultPort = 8001

//go:embed swagger.yaml
var swaggerSpec string

func main() {
	env.Setup()
	// separate the code from the 'main' function.
	// all code that available in main function were not testable
	Server()
}

func Server() {
	// Store + ORM
	if err := store.Setup(os.Getenv("ENV_DSN")); err != nil {
		log.Fatalf("unexpected  error while tried to connect to database: %v\n", err)
	}
	defer store.Close()

	// HTTP server
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
