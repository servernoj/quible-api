package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ably/ably-go/ably"
	"github.com/danielgtaylor/huma/v2"
	"github.com/quible-io/quible-api/app-service-ng/services/ablyService"
	libAPI "github.com/quible-io/quible-api/lib/api"
)

type GetLiveTokenInput struct {
}

type GetLiveTokenOutput struct {
	Body *ably.TokenRequest
}

func (impl *VersionedImpl) RegisterGetLiveToken(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID:   "get-live-token",
				Summary:       "Get live token",
				Description:   "Generate and return Ably `TokenRequest` bound to `live:main` channel",
				Method:        http.MethodGet,
				DefaultStatus: http.StatusOK,
				Errors:        []int{},
				Tags:          []string{"live", "public"},
				Path:          "/live/token",
			},
		),
		func(ctx context.Context, input *GetLiveTokenInput) (*GetLiveTokenOutput, error) {
			capabilities, _ := json.Marshal(&map[string][]string{
				"live:main": {"subscribe", "history"},
			})
			token, err := ablyService.CreateTokenRequest(&ably.TokenParams{
				Capability: string(capabilities),
				ClientID:   "nobody",
			})
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					err,
				)
			}
			return &GetLiveTokenOutput{
				Body: token,
			}, nil
		},
	)
}
