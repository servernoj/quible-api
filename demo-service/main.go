package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const BasePath = "/api/demo"
const DefaultPort = 8010

func main() {
	// load environment variables based on the value of `APP_ENV`:
	// 1. When it is undefined => from file `.env`
	// 2. When it is defined => from file `.env.${APP_ENV}`
	// Note: content from the file does not override existing env variables, it only adds
	if err := InitializeEnvironment(os.Getenv("APP_ENV")); err != nil {
		log.Fatalln("unable to initialize environment", err)
	}
	Server()
}

func Server() {
	// prepare gin
	gin.SetMode(gin.ReleaseMode)

	// gin setup
	r := gin.Default()
	r.Use(cors.Default())
	g := r.Group(BasePath)

	// Register helper routes
	g.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	// run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(DefaultPort)
		log.Printf("starting server on port: %s\n", port)
	}
	log.Fatalf("%v", r.Run(":"+port))
}
