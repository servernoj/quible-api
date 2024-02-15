package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/models"
)

type GetUserInput struct {
	AuthorizationHeaderResolver
}

type GetUserOutput struct {
	Body UserSimplified
}

func (impl *VersionedImpl) RegisterGetUser(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "get-user",
				Summary:     "Get user record",
				Description: "Return user record associated with the provided access token",
				Method:      http.MethodGet,
				Errors: []int{
					http.StatusUnauthorized,
				},
				Tags: []string{"user", "protected"},
				Path: "/user",
			},
		),
		func(ctx context.Context, input *GetUserInput) (*GetUserOutput, error) {
			// 1. Locate user based on access token send via Authorization header
			user, err := models.FindUserG(ctx, input.UserId)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err401_UserNotFound, err)
			}
			// 2. Return simplified version of the user object
			response := &GetUserOutput{
				Body: UserSimplified{
					ID:       user.ID,
					Username: user.Username,
					Email:    user.Email,
					Phone:    user.Phone,
					FullName: user.FullName,
				},
			}
			return response, nil
		},
	)
}
