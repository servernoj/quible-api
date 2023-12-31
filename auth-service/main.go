package main

import (
	_ "embed"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/quible-io/quible-api/auth-service/controller"
	"github.com/quible-io/quible-api/auth-service/realtime"
	"github.com/quible-io/quible-api/lib/env"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/store"
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
	Server()
}

func Server() {
	// -- Environment vars from .env file
	env.Setup()
	// -- Custom validators
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		misc.RegisterValidators(validate)
	} else {
		log.Println("unable to attach custom validators")
	}
	// -- Store + ORM
	if err := store.Setup(os.Getenv("ENV_DSN")); err != nil {
		log.Fatalf("unable to setup DB connection: %s", err)
	}
	defer store.Close()
	// -- Ably realtime
	if err := realtime.Setup(); err != nil {
		log.Fatalf("unable to setup Ably SDK: %s", err)
	}
	// -- HTTP server
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(location.Default())
	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "authorization")
	r.Use(cors.New(corsConfig))
	// API group
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
