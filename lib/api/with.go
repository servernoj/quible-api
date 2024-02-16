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
	Body []ErrorMapEntry
}

type ErrorMapEntry struct {
	Code        int
	Notation    string
	Description string
}

func WithErrorMap[T ErrorCodeConstraints](errorMap ErrorMap[T]) WithOption {
	return func(api huma.API, vc VersionConfig) {
		huma.Register(
			api,
			vc.Prefixer(
				huma.Operation{
					OperationID: "get-error-map-json",
					Summary:     "List of errors",
					Description: "Returns error codes as JSON object",
					Method:      http.MethodGet,
					Tags:        []string{"service", "public"},
					Path:        "/docs/errors",
				},
			),
			func(ctx context.Context, input *struct{}) (*ErrorMapOutput, error) {
				response := make([]ErrorMapEntry, len(errorMap))
				idx := 0
				for k, v := range errorMap {
					response[idx] = ErrorMapEntry{
						Code:        int(k),
						Notation:    k.String(),
						Description: v,
					}
					idx++
				}
				return &ErrorMapOutput{
					Body: response,
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
					Summary:     "API Version",
					Description: "Returns version of the API",
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
					Summary:       "Health status",
					Description:   "Returns health status, no response body",
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
