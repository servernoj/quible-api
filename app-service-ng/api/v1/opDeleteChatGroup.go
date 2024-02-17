package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/null/v8"
)

type DeleteChatGroupInput struct {
	AuthorizationHeaderResolver
	ChatGroupId string `path:"chatGroupId"`
}

type DeleteChatGroupOutput struct {
}

func (impl *VersionedImpl) RegisterDeleteChatGroup(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID:   "delete-chat-group",
				Summary:       "Delete chat group",
				Description:   "Delete chat group (logged in user must be the owner)",
				Method:        http.MethodDelete,
				DefaultStatus: http.StatusNoContent,
				Errors: []int{
					http.StatusUnauthorized,
					http.StatusNotFound,
				},
				Tags: []string{"chat", "protected"},
				Path: "/chat/groups/{chatGroupId}",
			},
		),
		func(ctx context.Context, input *DeleteChatGroupInput) (*DeleteChatGroupOutput, error) {
			chatGroup, err := models.Chats(
				models.ChatWhere.ID.EQ(input.ChatGroupId),
				models.ChatWhere.ParentID.IsNull(),
				models.ChatWhere.OwnerID.EQ(null.StringFrom(input.UserId)),
			).OneG(ctx)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err404_ChatGroupNotFound,
					err,
				)
			}
			_, err = chatGroup.ParentChats().DeleteAllG(ctx)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					err,
				)
			}
			_, err = chatGroup.DeleteG(ctx)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					err,
				)
			}
			return nil, nil
		},
	)
}
