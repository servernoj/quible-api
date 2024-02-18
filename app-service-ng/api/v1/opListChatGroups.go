package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/null/v8"
)

type ListChatGroupsInput struct {
	AuthorizationHeaderResolver
}

type ListChatGroupsOutput struct {
	Body models.ChatSlice
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
			chatGroups, err := models.Chats(
				models.ChatWhere.ParentID.IsNull(),
				models.ChatWhere.OwnerID.EQ(null.StringFrom(input.UserId)),
			).AllG(ctx)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
				)
			}
			return &ListChatGroupsOutput{
				Body: chatGroups,
			}, nil
		},
	)
}
