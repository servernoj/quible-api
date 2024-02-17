package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
)

type ListChatChannelsGroupedInput struct {
	AuthorizationHeaderResolver
}

type ListChatChannelsGroupedOutput struct {
}

func (impl *VersionedImpl) RegisterListChatChannelsGrouped(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID:   "list-chat-channels-grouped",
				Summary:       "List chat channels grouped",
				Description:   "List user's chat channels grouped based on their parents",
				Method:        http.MethodGet,
				DefaultStatus: http.StatusOK,
				Errors: []int{
					http.StatusUnauthorized,
				},
				Tags: []string{"chat", "protected"},
				Path: "/chat/channels/grouped",
			},
		),
		func(ctx context.Context, input *ListChatChannelsGroupedInput) (*ListChatChannelsGroupedOutput, error) {
			return nil, nil
		},
	)
}
