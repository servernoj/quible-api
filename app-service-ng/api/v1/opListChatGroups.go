package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
)

type ListChatGroupsInput struct {
	AuthorizationHeaderResolver
}

type ListChatGroupsOutput struct {
}

func (impl *VersionedImpl) RegisterListChatGroups(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID:   "list-chat-groups",
				Summary:       "List chat groups",
				Description:   "List user's chat groups",
				Method:        http.MethodGet,
				DefaultStatus: http.StatusOK,
				Errors: []int{
					http.StatusUnauthorized,
				},
				Tags: []string{"chat", "protected"},
				Path: "/chat/groups",
			},
		),
		func(ctx context.Context, input *ListChatGroupsInput) (*ListChatGroupsOutput, error) {
			return nil, nil
		},
	)
}
