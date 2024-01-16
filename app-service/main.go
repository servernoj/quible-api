package main

import (
	_ "embed"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/app-service/BasketAPI"
	"github.com/quible-io/quible-api/app-service/controller"
	"github.com/quible-io/quible-api/lib/env"
	"github.com/quible-io/quible-api/lib/store"
)

//	@title			Quible app-service
//	@description	Wrapper to RSC API
//	@version		0.1
//	@host			www.quible.io
//	@BasePath		/api/v1

const DefaultPort = 8002

//go:embed swagger.yaml
var swaggerSpec string

func main() {
	Server()
}

func Server() {
	// -- Environment vars from .env file
	env.Setup()
	// -- Store + ORM
	if err := store.Setup(os.Getenv("ENV_DSN")); err != nil {
		log.Fatalf("unable to setup DB connection: %s", err)
	}
	defer store.Close()
	// -- Live data BasketAPI
	quit, err := BasketAPI.StartLive()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		quit <- struct{}{}
	}()
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
