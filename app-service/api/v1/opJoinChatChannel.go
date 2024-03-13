package v1

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type JoinChatChannelInput struct {
	AuthorizationHeaderResolver
	ChatChannelId string `path:"chatChannelId" format:"uuid"`
}

type JoinChatChannelOutput struct {
}

func (impl *VersionedImpl) RegisterJoinChatChannel(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID:   "join-chat-channel",
				Summary:       "Join chat channel",
				Description:   "Join chat channel associated with public chat group",
				Method:        http.MethodPost,
				DefaultStatus: http.StatusOK,
				Errors: []int{
					http.StatusUnauthorized,
					http.StatusBadRequest,
					http.StatusNotFound,
				},
				Tags: []string{"chat", "protected"},
				Path: "/chat/channels/{chatChannelId}",
			},
		),
		func(ctx context.Context, input *JoinChatChannelInput) (*JoinChatChannelOutput, error) {
			// 0. Dependences
			deps := impl.Deps.GetContext("opJoinChatChannel")
			db := deps.Get("db").(*sql.DB)
			chatChannel, err := models.FindChat(ctx, db, input.ChatChannelId)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err404_ChatChannelNotFound,
					err,
				)
			}
			chatGroup, err := chatChannel.Parent().One(ctx, db)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err404_ChatGroupNotFound,
					err,
				)
			}
			if chatGroup.IsPrivate.Bool {
				return nil, ErrorMap.GetErrorResponse(
					Err400_ChatGroupIsPrivate,
				)
			}
			if chatGroup.OwnerID.String == input.UserId {
				return nil, ErrorMap.GetErrorResponse(
					Err400_ChatGroupIsSelfOwned,
				)
			}
			chatUserFound, err := models.ChatUserExists(ctx, db, input.ChatChannelId, input.UserId)
			if chatUserFound || err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err400_ChatChannelAlreadyJoined,
				)
			}
			chatUser := models.ChatUser{
				ChatID: input.ChatChannelId,
				UserID: input.UserId,
			}
			return nil, chatUser.Insert(ctx, db, boil.Infer())
		},
	)
}
