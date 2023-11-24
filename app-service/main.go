package main

import (
	_ "embed"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gitlab.com/quible-backend/app-service/controller"
	"gitlab.com/quible-backend/lib/env"
	"gitlab.com/quible-backend/lib/misc"
	"gitlab.com/quible-backend/lib/store"
)

//	@title			Quible app-service
//	@description	Authentication and authorization service of Quible.io
//	@version		0.1
//	@host			www.quible.io
//	@BasePath		/api/v1

const DefaultPort = 8021

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