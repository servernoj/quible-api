package api

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const Title = "Quible app service"

//go:embed serviceDescription.md
var ServiceDescription string

type VersionConfig struct {
	Tag    string
	SemVer string
}

func (vc VersionConfig) Prefixer(op huma.Operation) huma.Operation {
	op.Path = fmt.Sprintf("/api/%s%s", vc.Tag, op.Path)
	return op
}

func (vc VersionConfig) GetConfig(title string) huma.Config {
	prefix := fmt.Sprintf("/api/%s", vc.Tag)
	schemaPrefix := "#/components/schemas/"
	schemasPath := fmt.Sprintf("%s/schemas", prefix)
	registry := huma.NewMapRegistry(schemaPrefix, huma.DefaultSchemaNamer)
	return huma.Config{
		OpenAPI: &huma.OpenAPI{
			OpenAPI: "3.1.0",
			Info: &huma.Info{
				Title:       title,
				Version:     vc.SemVer,
				Description: ServiceDescription,
			},
			Components: &huma.Components{
				Schemas: registry,
			},
			Tags: []*huma.Tag{
				{Name: "public", Description: "API with public access"},
				{Name: "protected", Description: "API requiring authentication"},
				{Name: "service", Description: "API for handling service metadata"},
			},
			OnAddOperation: []huma.AddOpFunc{},
		},
		OpenAPIPath: fmt.Sprintf("%s/docs/openapi", prefix),
		DocsPath:    fmt.Sprintf("%s/docs", prefix),
		SchemasPath: schemasPath,
		Formats: map[string]huma.Format{
			"application/json": huma.DefaultJSONFormat,
			"json":             huma.DefaultJSONFormat,
			"application/cbor": huma.DefaultCBORFormat,
			"cbor":             huma.DefaultCBORFormat,
		},
		DefaultFormat: "application/json",
		Transformers: []huma.Transformer{
			func(ctx huma.Context, status string, v any) (any, error) {
				if statusError, ok := v.(huma.StatusError); ok {
					// override response status to respect the status returned by NewError
					ctx.SetStatus(statusError.GetStatus())
				}
				return v, nil
			},
		},
	}
}

func Setup[Impl ErrorReporter](router *gin.Engine, vc VersionConfig, withOptions ...WithOption) {
	// 1. Initialize config with version-prefixed fields
	config := vc.GetConfig(Title)
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
