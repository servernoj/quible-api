package api

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type SetupFunc[Impl ErrorReporter] func(router *gin.Engine, vc VersionConfig, withOptions ...WithOption)

func SetupFactory[Impl ErrorReporter](Title, ServiceDescription string) SetupFunc[Impl] {
	return func(router *gin.Engine, vc VersionConfig, withOptions ...WithOption) {
		// 1. Initialize config with version-prefixed fields
		config := vc.GetConfig(Title, ServiceDescription)
		// 2. Create API instance
		api := humagin.New(router, config)
		// 3. Override default error reporting facility
		huma.NewError = func(status int, message string, errs ...error) huma.StatusError {
			if status == 0 {
				// case for https://github.com/danielgtaylor/huma/issues/236
				return &ErrorResponse{}
			}
			var implValue Impl
			if len(errs) > 0 {
				b, _ := json.MarshalIndent(errs, "", "  ")
				log.Error().Msgf("Validation error(s): %s", b)
			}
			return implValue.NewError(status, message, errs...)
		}
		// 4. Register all version-specific endpoints
		var implPtr *Impl
		implType := reflect.TypeOf(implPtr)
		args := []reflect.Value{
			reflect.ValueOf(implPtr),
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
	}
}
