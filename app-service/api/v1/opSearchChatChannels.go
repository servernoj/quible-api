package v1

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SearchChatChannelsInput struct {
	Q string `query:"q" doc:"search term to be partially matched agains chat group title"`
}

type SearchChatChannelsOutput struct {
	Body []ChatChannelsGroup
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
			// 0. Dependences
			deps := impl.Deps.GetContext("opSearchChatChannels")
			db := deps.Get("db").(*sql.DB)
			// 1. Find all matching chat groups
			chatGroups, err := models.Chats(
				models.ChatWhere.Title.ILIKE("%"+input.Q+"%"),
				models.ChatWhere.ParentID.IsNull(),
				models.ChatWhere.IsPrivate.EQ(null.BoolFrom(false)),
				qm.Load(
					models.ChatRels.ParentChats,
				),
			).All(ctx, db)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					err,
				)
			}
			// 2. Extract associated chat channels
			chatChannelsGroupSlice := []ChatChannelsGroup{}
			for _, chatGroup := range chatGroups {
				chatChannels := make([]ChatChannel, len(chatGroup.R.ParentChats))
				for idx, chat := range chatGroup.R.ParentChats {
					chatChannels[idx] = ChatChannel{
						ID:       chat.ID,
						Title:    chat.Title,
						Resource: chatGroup.Resource + ":" + chat.Resource,
					}
				}
				chatChannelsGroupSlice = append(
					chatChannelsGroupSlice,
					ChatChannelsGroup{
						ID:           chatGroup.ID,
						Title:        chatGroup.Title,
						Summary:      chatGroup.Summary.Ptr(),
						ChatChannels: chatChannels,
					},
				)
			}
			// 3. Prepare and return the response
			return &SearchChatChannelsOutput{
				Body: chatChannelsGroupSlice,
			}, nil
		},
	)
}
