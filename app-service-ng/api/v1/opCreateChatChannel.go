package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type CreateChatChannelInput struct {
	AuthorizationHeaderResolver
	GroupId string `path:"chatGroupId" format:"uuid"`
	Body    struct {
		// second part of the channel resource name following `chat group`
		Name    string  `json:"name" pattern:"\\w+" doc:"second part of the channel resource name (everything after ':')"`
		Title   string  `json:"title" doc:"human-readable 'title' of the channel"`
		Summary *string `json:"summary,omitempty" doc:"Optional summary, potentially lengthy text"`
	}
}

type CreateChatChannelOutput struct {
	Body models.Chat
}

func (impl *VersionedImpl) RegisterCreateChatChannel(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID:   "create-chat-channel",
				Summary:       "Create chat channel",
				Description:   "Create a channel within the specified chat group (logged in user must own that chat group)",
				Method:        http.MethodPost,
				DefaultStatus: http.StatusCreated,
				Errors: []int{
					http.StatusUnauthorized,
					http.StatusBadRequest,
				},
				Tags: []string{"chat", "protected"},
				Path: "/chat/groups/{chatGroupId}/channels",
			},
		),
		func(ctx context.Context, input *CreateChatChannelInput) (*CreateChatChannelOutput, error) {
			chatGroupFound, err := models.Chats(
				models.ChatWhere.ID.EQ(input.GroupId),
				models.ChatWhere.ParentID.IsNull(),
			).ExistsG(ctx)
			if err != nil || !chatGroupFound {
				return nil, ErrorMap.GetErrorResponse(
					Err404_ChatGroupNotFound,
					err,
				)
			}
			channelFound, _ := models.Chats(
				models.ChatWhere.Resource.EQ(input.Body.Name),
				models.ChatWhere.ParentID.EQ(null.StringFrom(input.GroupId)),
			).ExistsG(ctx)
			if channelFound {
				return nil, ErrorMap.GetErrorResponse(
					Err400_ChannelExists,
				)
			}
			chatChannel := models.Chat{
				Resource:  input.Body.Name,
				Title:     input.Body.Title,
				Summary:   null.StringFromPtr(input.Body.Summary),
				ParentID:  null.StringFrom(input.GroupId),
				OwnerID:   null.StringFromPtr(nil),
				IsPrivate: null.BoolFromPtr(nil),
			}
			if err := chatChannel.InsertG(ctx, boil.Infer()); err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					err,
				)
			}
			return &CreateChatChannelOutput{
				Body: chatChannel,
			}, nil

		},
	)
}
