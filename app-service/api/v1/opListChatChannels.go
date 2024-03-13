package v1

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
)

type ListChatChannelsInput struct {
	AuthorizationHeaderResolver
}

type ListChatChannelsOutput struct {
	Body any
}

func (impl *VersionedImpl) RegisterListChatChannels(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID:   "list-chat-channels",
				Summary:       "List chat channels",
				Description:   "List user's chat channels",
				Method:        http.MethodGet,
				DefaultStatus: http.StatusOK,
				Errors: []int{
					http.StatusUnauthorized,
				},
				Tags: []string{"chat", "protected"},
				Path: "/chat/channels",
			},
		),
		func(ctx context.Context, input *ListChatChannelsInput) (*ListChatChannelsOutput, error) {
			// 0. Dependences
			deps := impl.Deps.GetContext("opListChatChannels")
			db := deps.Get("db").(*sql.DB)
			chatChannels, err := chatChannelsForUser(ctx, db, input.UserId)
			if err != nil {
				return nil, err
			}
			return &ListChatChannelsOutput{
				Body: chatChannels,
			}, nil
		},
	)
}
