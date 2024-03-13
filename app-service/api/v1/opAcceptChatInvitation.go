package v1

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type AcceptChatInvitationInput struct {
	Body struct {
		Token string `json:"token" pattern:"^[^.]+([.][^.]+){2}$"`
	}
}

type AcceptChatInvitationOutput struct {
}

func (impl *VersionedImpl) RegisterAcceptChatInvitation(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID:   "accept-chat-invitation",
				Summary:       "Accept chat invitation",
				Description:   "Accept chat channel invitation initiated by clicking email link",
				Method:        http.MethodPost,
				DefaultStatus: http.StatusOK,
				Errors: []int{
					http.StatusBadRequest,
					http.StatusExpectationFailed,
				},
				Tags: []string{"chat", "public"},
				Path: "/chat/channels/accept",
			},
		),
		func(ctx context.Context, input *AcceptChatInvitationInput) (*AcceptChatInvitationOutput, error) {
			// 0. Dependences
			deps := impl.Deps.GetContext("opAcceptChatInvitation")
			db := deps.Get("db").(*sql.DB)
			// 1. Process invitation token from request body
			tokenClaims, err := jwt.VerifyJWT(input.Body.Token, jwt.TokenActionInvitationToPrivateChat)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err417_InvalidToken,
					err,
				)
			}
			invitorId := tokenClaims["userId"].(string)
			extraClaims := tokenClaims["extraClaims"].(jwt.ExtraClaims)
			inviteeId, ok := extraClaims["inviteeId"].(string)
			if !ok {
				return nil, ErrorMap.GetErrorResponse(
					Err417_InvalidToken,
					errors.New("missing or invalid inviteeId in extraClaims"),
				)
			}
			chatChannelId, ok := extraClaims["chatChannelId"].(string)
			if !ok {
				return nil, ErrorMap.GetErrorResponse(
					Err417_InvalidToken,
					errors.New("missing or invalid chatChannelId in extraClaims"),
				)
			}
			// 2. Locate DB records for chat channel, its holding group and its association with invitee
			chatChannel, err := models.Chats(
				models.ChatWhere.ID.EQ(chatChannelId),
				models.ChatWhere.ParentID.IsNotNull(),
				qm.Load(
					models.ChatRels.Parent,
					models.ChatWhere.OwnerID.EQ(null.StringFrom(invitorId)),
				),
				qm.Load(
					models.ChatRels.ChatUsers,
					models.ChatUserWhere.UserID.EQ(inviteeId),
				),
			).One(ctx, db)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					err,
				)
			}
			if chatChannel.R.Parent == nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					errors.New("qualified parent chat group not found"),
				)
			}
			if len(chatChannel.R.ChatUsers) == 0 {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					errors.New("qualified chat-user association not found"),
				)
			}
			// 3. Update chat-user association (to clear "disabled" flag)
			chatUser := chatChannel.R.ChatUsers[0]
			chatUser.Disabled = false
			if _, err := chatUser.Update(ctx, db, boil.Whitelist("disabled")); err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnableUpdateChatUser,
					err,
				)
			}
			return nil, nil
		},
	)
}
