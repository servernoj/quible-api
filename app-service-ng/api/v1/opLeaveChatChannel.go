package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
)

type LeaveChatChannelInput struct {
	AuthorizationHeaderResolver
}

type LeaveChatChannelOutput struct {
}

func (impl *VersionedImpl) RegisterLeaveChatChannel(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID:   "leave-chat-channel",
				Summary:       "Leave chat channel",
				Description:   "Leave previously joined chat channel",
				Method:        http.MethodDelete,
				DefaultStatus: http.StatusNoContent,
				Errors: []int{
					http.StatusUnauthorized,
					http.StatusNotFound,
				},
				Tags: []string{"chat", "protected"},
				Path: "/chat/channels/{channelId}",
			},
		),
		func(ctx context.Context, input *LeaveChatChannelInput) (*LeaveChatChannelOutput, error) {
			return nil, nil
		},
	)
}
