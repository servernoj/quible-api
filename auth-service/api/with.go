package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type WithOption func(api huma.API, vc VersionConfig)

type RawBody struct {
	ContentType string `header:"content-type"`
	Body        []byte
}

type ErrorMapOutput struct {
	Body any
}

func WithErrorMap(errorMap any) WithOption {
	return func(api huma.API, vc VersionConfig) {
		huma.Register(
			api,
			vc.Prefixer(
				huma.Operation{
					OperationID: "get-error-map-json",
					Summary:     "Returns error codes as JSON object",
					Method:      http.MethodGet,
					Tags:        []string{"service", "public"},
					Path:        "/docs/errors",
				},
			),
			func(ctx context.Context, input *struct{}) (*ErrorMapOutput, error) {
				return &ErrorMapOutput{
					Body: errorMap,
				}, nil
			},
		)
	}
}

func WithVersion() WithOption {
	return func(api huma.API, vc VersionConfig) {
		huma.Register(
			api,
			vc.Prefixer(
				huma.Operation{
					OperationID: "get-version",
					Summary:     "Returns version of the called API",
					Method:      http.MethodGet,
					Tags:        []string{"service", "public"},
					Path:        "/version",
				},
			),
			func(ctx context.Context, input *struct{}) (*RawBody, error) {
				return &RawBody{
					ContentType: "text/plain",
					Body:        []byte(vc.SemVer),
				}, nil
			},
		)
	}
}

func WithHealth() WithOption {
	return func(api huma.API, vc VersionConfig) {
		huma.Register(
			api,
			vc.Prefixer(
				huma.Operation{
					OperationID:   "get-heath",
					Summary:       "Returns health status, no response body",
					Method:        http.MethodGet,
					DefaultStatus: http.StatusOK,
					Tags:          []string{"service", "public"},
					Path:          "/health",
				},
			),
			func(ctx context.Context, input *struct{}) (*struct{}, error) {
				return nil, nil
			},
		)
	}
}
