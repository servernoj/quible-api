package api

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/rs/zerolog/log"
)

type ErrorReporter interface {
	NewError(int, string, ...error) huma.StatusError
}

type EmailSenderSetter interface {
	SetEmailSender(email.EmailSender)
}

type ServiceAPI interface {
	ErrorReporter
	EmailSenderSetter
	email.EmailSender
}

type PostInit func(ServiceAPI, *gin.Engine, VersionConfig, ...WithOption) huma.API

func GetPostInit(title string, description string) PostInit {
	return func(serviceAPI ServiceAPI, router *gin.Engine, vc VersionConfig, withOptions ...WithOption) huma.API {
		// 1. Initialize config with version-prefixed fields
		config := vc.GetConfig(title, description)
		// 2. Create API instance
		api := humagin.New(router, config)
		// 3. Override default error reporting facility
		huma.NewError = func(status int, message string, errs ...error) huma.StatusError {
			if status == 0 {
				return &ErrorResponse{}
			}
			if len(errs) > 0 {
				b, _ := json.MarshalIndent(errs, "", "  ")
				log.Error().Msgf("Validation error(s): %s", b)
			}
			return serviceAPI.NewError(status, message, errs...)
		}
		// 4. Register all implementation-specific endpoints
		implType := reflect.TypeOf(serviceAPI)
		args := []reflect.Value{
			reflect.ValueOf(serviceAPI),
			reflect.ValueOf(api),
			reflect.ValueOf(vc),
		}
		for i := 0; i < implType.NumMethod(); i++ {
			m := implType.Method(i)
			if strings.HasPrefix(m.Name, "Register") && len(m.Name) > 8 {
				m.Func.Call(args)
			}
		}
		// 5. Register all optional [shared] endpoints
		for _, option := range withOptions {
			option(api, vc)
		}
		return api
	}
}
