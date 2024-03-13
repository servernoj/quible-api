package v1

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type CreateChatGroupInput struct {
	AuthorizationHeaderResolver
	Body struct {
		Name      string  `json:"name" pattern:"\\w+" doc:"unique across all chat groups owned by the same user"`
		Title     string  `json:"title" doc:"human-readable 'title' of the chat group"`
		Summary   *string `json:"summary,omitempty"`
		IsPrivate bool    `json:"isPrivate"`
	}
}

type CreateChatGroupOutput struct {
	Body models.Chat
}

func (impl *VersionedImpl) RegisterCreateChatGroup(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID:   "create-chat-group",
				Summary:       "Create chat group",
				Description:   "Create a chat group owned by the logged in user",
				Method:        http.MethodPost,
				DefaultStatus: http.StatusCreated,
				Errors: []int{
					http.StatusUnauthorized,
					http.StatusBadRequest,
				},
				Tags: []string{"chat", "protected"},
				Path: "/chat/groups",
			},
		),
		func(ctx context.Context, input *CreateChatGroupInput) (*CreateChatGroupOutput, error) {
			// 0. Dependences
			deps := impl.Deps.GetContext("opCreateChatGroup")
			db := deps.Get("db").(*sql.DB)
			resource := GROUP_PREFIX + input.Body.Name
			chatGroupFound, err := models.Chats(
				models.ChatWhere.OwnerID.EQ(null.StringFrom(input.UserId)),
				models.ChatWhere.ParentID.IsNull(),
				qm.Expr(
					qm.Or2(models.ChatWhere.Resource.EQ(resource)),
					qm.Or2(models.ChatWhere.Title.EQ(input.Body.Title)),
				),
			).Exists(ctx, db)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					err,
				)
			}
			if chatGroupFound {
				return nil, ErrorMap.GetErrorResponse(
					Err400_ChatGroupExists,
				)
			}
			chatGroup := models.Chat{
				Resource:  resource,
				ParentID:  null.StringFromPtr(nil),
				IsPrivate: null.BoolFrom(input.Body.IsPrivate),
				OwnerID:   null.StringFrom(input.UserId),
				Summary:   null.StringFromPtr(input.Body.Summary),
				Title:     input.Body.Title,
			}
			if err := chatGroup.Insert(ctx, db, boil.Infer()); err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					err,
				)
			}
			return &CreateChatGroupOutput{
				Body: chatGroup,
			}, nil
		},
	)
}
