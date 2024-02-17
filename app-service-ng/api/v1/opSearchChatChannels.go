package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
)

type SearchChatChannelsInput struct {
}

type SearchChatChannelsOutput struct {
}

func (impl *VersionedImpl) RegisterSearchChatChannels(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID:   "search-chat-channels",
				Summary:       "Search chat channels",
				Description:   "Search public chat channels and report them along with the parent group",
				Method:        http.MethodGet,
				DefaultStatus: http.StatusOK,
				Errors: []int{
					http.StatusUnauthorized,
				},
				Tags: []string{"chat", "public"},
				Path: "/chat/channels/search",
			},
		),
		func(ctx context.Context, input *SearchChatChannelsInput) (*SearchChatChannelsOutput, error) {
			return nil, nil
		},
	)
}
