package main

import (
	_ "embed"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/quible-backend/app-service/controller"
	c "gitlab.com/quible-backend/lib/controller"
	"gitlab.com/quible-backend/lib/env"
)

//	@title			Quible app-service
//	@description	Wrapper to RSC API
//	@version		0.1
//	@host			www.quible.io
//	@BasePath		/api/v1

const DefaultPort = 8001

//go:embed swagger.yaml
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
		c.WithSwagger(swaggerSpec),
		c.WithHealth(),
	)
	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(DefaultPort)
	}
	log.Printf("starting server on port: %s\n", port)
	log.Fatalf("%v", r.Run(":"+port))
}
