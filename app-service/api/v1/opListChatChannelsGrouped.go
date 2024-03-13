package v1

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
)

type ListChatChannelsGroupedInput struct {
	AuthorizationHeaderResolver
}

type ListChatChannelsGroupedOutput struct {
	Body []ChatChannelsGroup
}

type ChatChannelsGroup struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Summary      *string       `json:"summary"`
	ChatChannels []ChatChannel `json:"chatChannels"`
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
			// 0. Dependences
			deps := impl.Deps.GetContext("opListChatChannelsGrouped")
			db := deps.Get("db").(*sql.DB)
			// 1. Get all user's chat channels
			chatChannels, err := chatChannelsForUser(ctx, db, input.UserId)
			if err != nil {
				return nil, err
			}
			// 2. Group chat channels based on `Parent` field in each record
			chatChannelsGroupMap := map[string]*ChatChannelsGroup{}
			for _, chatChannel := range chatChannels {
				chatGroup := chatChannel.Parent
				if _, ok := chatChannelsGroupMap[chatGroup.ID]; !ok {
					chatChannelsGroupMap[chatGroup.ID] = &ChatChannelsGroup{
						ID:           chatGroup.ID,
						Title:        chatGroup.Title,
						Summary:      chatGroup.Summary.Ptr(),
						ChatChannels: []ChatChannel{},
					}
				}
				chatChannelsGroup := chatChannelsGroupMap[chatGroup.ID]
				chatChannelsGroup.ChatChannels = append(
					chatChannelsGroup.ChatChannels,
					chatChannel,
				)
			}
			// 3. Prepare and return the response
			chatChannelsGroupSlice := []ChatChannelsGroup{}
			for _, chatChannelsGroup := range chatChannelsGroupMap {
				chatChannelsGroupSlice = append(chatChannelsGroupSlice, *chatChannelsGroup)
			}
			return &ListChatChannelsGroupedOutput{
				Body: chatChannelsGroupSlice,
			}, nil
		},
	)
}
