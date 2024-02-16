package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	srvAPI "github.com/quible-io/quible-api/app-service-ng/api"
	v1 "github.com/quible-io/quible-api/app-service-ng/api/v1"
	"github.com/quible-io/quible-api/app-service-ng/services/BasketAPI"
	"github.com/quible-io/quible-api/app-service-ng/services/ablyService"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/env"
	"github.com/quible-io/quible-api/lib/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ServiceOptions struct {
	Port int `help:"Port to listen on" short:"p" default:"8002"`
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	Server()
}

func Server() {
	// -- Environment vars from .env file
	env.Setup()
	// -- Store + ORM
	if err := store.Setup(os.Getenv("ENV_DSN")); err != nil {
		log.Error().Msgf("unable to setup DB connection: %s", err)
		os.Exit(1)
	}
	defer store.Close()
	// -- Ably
	if err := ablyService.Setup(); err != nil {
		log.Fatal().Msgf("unable to setup Ably SDK: %s", err)
	}
	// -- Live data BasketAPI
	quit, err := BasketAPI.StartLive()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer func() {
		quit <- struct{}{}
	}()
	// -- Huma CLI
	cli := huma.NewCLI(func(hooks huma.Hooks, options *ServiceOptions) {
		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowCredentials = true
		corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "authorization")
		router.Use(cors.New(corsConfig))
		// http server based on Gin router
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		if port == 0 {
			port = options.Port
		}
		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		}
		// -- V1
		srvAPI.Setup[v1.VersionedImpl](
			router,
			libAPI.VersionConfig{
				Tag:    "v1",
				SemVer: "1.0.0",
			},
			libAPI.WithErrorMap(v1.ErrorMap),
			libAPI.WithVersion(),
			libAPI.WithHealth(),
		)
		// Hooks
		hooks.OnStart(func() {
			log.Info().Msgf("starting server on port: %d", port)
			log.Error().Err(server.ListenAndServe()).Send()
			os.Exit(10)
		})
		hooks.OnStop(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err := server.Shutdown(ctx)
			if err != nil {
				log.Error().Err(err).Send()
			}
		})
	})
	cli.Run()
}
