package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/models"
)

type LeaveChatChannelInput struct {
	AuthorizationHeaderResolver
	ChatChannelId string `path:"chatChannelId" format:"uuid"`
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
				Path: "/chat/channels/{chatChannelId}",
			},
		),
		func(ctx context.Context, input *LeaveChatChannelInput) (*LeaveChatChannelOutput, error) {
			chatUser, err := models.ChatUsers(
				models.ChatUserWhere.ChatID.EQ(input.ChatChannelId),
				models.ChatUserWhere.UserID.EQ(input.UserId),
			).OneG(ctx)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err404_ChatChannelNotFound,
					err,
				)
			}
			if _, err := chatUser.DeleteG(ctx); err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					err,
				)
			}
			return nil, nil
		},
	)
}