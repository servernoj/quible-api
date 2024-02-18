package api

import (
	"fmt"

	"github.com/danielgtaylor/huma/v2"
)

type VersionConfig struct {
	Tag    string
	SemVer string
}

func (vc VersionConfig) Prefixer(op huma.Operation) huma.Operation {
	op.Path = fmt.Sprintf("/api/%s%s", vc.Tag, op.Path)
	return op
}

func (vc VersionConfig) GetConfig(title, description string) huma.Config {
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
				Description: description,
			},
			Components: &huma.Components{
				Schemas: registry,
			},
			Tags:           []*huma.Tag{},
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
