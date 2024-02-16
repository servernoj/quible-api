package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
)

type AuthTestInput struct {
	AuthorizationHeaderResolver
}

type AuthTestOutput struct {
}

func (impl *VersionedImpl) RegisterAuthTest(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "get-auth-test",
				Summary:     "Test authentication",
				Method:      http.MethodGet,
				Errors: []int{
					http.StatusUnauthorized,
				},
				Tags: []string{"test"},
				Path: "/test",
			},
		),
		func(ctx context.Context, input *AuthTestInput) (*AuthTestOutput, error) {
			return nil, nil
		},
	)
}
