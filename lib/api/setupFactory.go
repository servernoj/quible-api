package api

import (
	"encoding/json"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/rs/zerolog/log"
)

type ErrorReporter interface {
	NewError(int, string, ...error) huma.StatusError
}

type OpsRegistrar interface {
	Register(huma.API, VersionConfig)
}

type ServiceAPI interface {
	ErrorReporter
	email.EmailSender
	OpsRegistrar
}

type PostInit func(ServiceAPI, *gin.Engine, VersionConfig, ...WithOption) huma.API

func overrideHumaNewError(implValue ErrorReporter) {
	huma.NewError = func(status int, message string, errs ...error) huma.StatusError {
		if status == 0 {
			return &ErrorResponse{}
		}
		if len(errs) > 0 {
			b, _ := json.MarshalIndent(errs, "", "  ")
			log.Error().Msgf("Validation error(s): %s", b)
		}
		return implValue.NewError(status, message, errs...)
	}
}

func GetPostInit(title string, description string) PostInit {
	return func(serviceAPI ServiceAPI, router *gin.Engine, vc VersionConfig, withOptions ...WithOption) huma.API {
		// 1. Initialize config with version-prefixed fields
		config := vc.GetConfig(title, description)
		// 2. Create API instance
		api := humagin.New(router, config)
		// 3. Override default error reporting facility
		overrideHumaNewError(serviceAPI)
		// 4. Register all version-specific endpoints
		serviceAPI.Register(api, vc)
		// 5. Register all optional [shared] endpoints
		for _, option := range withOptions {
			option(api, vc)
		}
		return api
	}
}
