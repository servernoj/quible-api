package main

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/quible-backend/auth-service/service"
	"gitlab.com/quible-backend/lib/env"
	"gitlab.com/quible-backend/lib/store"
	"gitlab.com/quible-backend/lib/swagger"
)

//	@title			Quible auth-service
//	@description	Authentication and authorization service of Quible.io
//	@version		0.1
//	@host			www.quible.io
//	@BasePath		/api/v1

const BasePath = "/api/v1"
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
	// prepare gin
	gin.SetMode(gin.ReleaseMode)

	// gin setup
	r := gin.Default()
	r.Use(cors.Default())
	g := r.Group(BasePath)

	if err := store.Init(os.Getenv("ENV_DSN")); err != nil {
		log.Fatalf("unexpected  error while tried to connect to database: %v\n", err)
	}

	userService := service.UserService{
		C: context.Background(),
	}
	user, _ := userService.GetUserById(1)
	log.Printf("%+v\n", user)

	// prepare postgresql database
	// dbPool, err := config.NewDBPool(
	// 	os.Getenv("ENV_DSN"),
	// )

	// log for error if error occur while connecting to the database
	// if err != nil {
	// 	log.Fatalf("unexpected  error while tried to connect to database: %v\n", err)
	// }

	// defer dbPool.Close()

	// // setup api
	// database := user.NewRepository(dbPool)
	// service := user.NewService(database)
	// controller := user.NewController(service)

	// Register User controller routes
	// user.Routes(g, controller)
	// Register Swagger/docs routes
	swagger.Register(g, swaggerSpec, "/docs")
	// Register helper routes
	g.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	// run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(DefaultPort)
	}
	log.Printf("starting server on port: %s\n", port)
	log.Fatalf("%v", r.Run(":"+port))
}
