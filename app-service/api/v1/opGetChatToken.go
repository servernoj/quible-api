package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ably/ably-go/ably"
	"github.com/danielgtaylor/huma/v2"
	"github.com/quible-io/quible-api/app-service/services/ablyService"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/null/v8"
)

type GetChatTokenInput struct {
	AuthorizationHeaderResolver
}

type GetChatTokenOutput struct {
	Body *ably.TokenRequest
}

func (impl *VersionedImpl) RegisterGetChatToken(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID:   "get-chat-token",
				Summary:       "Get chat token",
				Description:   "Generate and return Ably `TokenRequest` associated with the logged in user",
				Method:        http.MethodGet,
				DefaultStatus: http.StatusOK,
				Errors: []int{
					http.StatusUnauthorized,
				},
				Tags: []string{"chat", "protected"},
				Path: "/chat/token",
			},
		),
		func(ctx context.Context, input *GetChatTokenInput) (*GetChatTokenOutput, error) {
			// 1. Compute map of capabilities
			capabilities := map[string][]string{}
			// 1a. Process implied capabilities from self-owned chat groups
			chatGroups, err := models.Chats(
				models.ChatWhere.ParentID.IsNull(),
				models.ChatWhere.OwnerID.EQ(null.StringFrom(input.UserId)),
			).AllG(ctx)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					err,
				)
			}
			for _, chatGroup := range chatGroups {
				resource := chatGroup.Resource + ":*"
				capabilities[resource] = AccessReadWrite
			}
			// 1b. Process joined channels
			chatUsers, err := models.ChatUsers(
				models.ChatUserWhere.UserID.EQ(input.UserId),
				models.ChatUserWhere.Disabled.EQ(false),
			).AllG(ctx)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					err,
				)
			}
			for _, item := range chatUsers {
				chatId := item.ChatID
				access := AccessReadWrite
				if item.IsRo {
					access = AccessReadOnly
				}
				chat, err := models.FindChatG(ctx, chatId)
				if err != nil {
					return nil, ErrorMap.GetErrorResponse(
						Err500_UnknownError,
						err,
					)
				}
				parentChatGroup, err := chat.Parent().OneG(ctx)
				if err != nil {
					return nil, ErrorMap.GetErrorResponse(
						Err500_UnknownError,
						err,
					)
				}
				resource := parentChatGroup.Resource + ":" + chat.Resource
				if accessFound, ok := capabilities[resource]; ok && len(accessFound) > len(access) {
					access = accessFound
				}
				capabilities[resource] = access
			}
			// 2. Prepare and return `TokenRequest` in response
			marshalledCapabilities, _ := json.Marshal(&capabilities)
			token, err := ablyService.CreateTokenRequest(&ably.TokenParams{
				Capability: string(marshalledCapabilities),
				ClientID:   input.UserId,
			})
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					err,
				)
			}
			return &GetChatTokenOutput{
				Body: token,
			}, nil
		},
	)
}
