package v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/quible-io/quible-api/app-service-ng/services/emailService"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type InviteUserInput struct {
	AuthorizationHeaderResolver
	ChatChannelId string `path:"chatChannelId"`
	Body          struct {
		Email string `json:"email" format:"email"`
	}
}

type InviteUserOutput struct {
}

func (impl *VersionedImpl) RegisterInviteUser(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID:   "invite-user",
				Summary:       "Invite user",
				Description:   "Invite a user to join private channel",
				Method:        http.MethodPost,
				DefaultStatus: http.StatusOK,
				Errors: []int{
					http.StatusBadRequest,
					http.StatusUnauthorized,
					http.StatusNotFound,
				},
				Tags: []string{"chat", "protected"},
				Path: "/chat/channels/{chatChannelId}/invite",
			},
		),
		func(ctx context.Context, input *InviteUserInput) (*InviteUserOutput, error) {
			// 1. test if request chat channel exists
			chatChannel, err := models.Chats(
				models.ChatWhere.ID.EQ(input.ChatChannelId),
				models.ChatWhere.ParentID.IsNotNull(),
				qm.Load(
					models.ChatRels.Parent,
				),
			).OneG(ctx)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err404_ChatChannelNotFound,
					err,
				)
			}
			// 1a. ...and it belongs to a private chat group which is owned by the invitor
			chatGroup := chatChannel.R.Parent
			if !chatGroup.IsPrivate.Bool || chatGroup.OwnerID != null.StringFrom(input.UserId) {
				return nil, ErrorMap.GetErrorResponse(
					Err404_ChatChannelNotFound,
					errors.New("holding chat group is not private or is not owned by user"),
				)
			}
			// 2. Find invitee user by provided email
			invitee, err := models.Users(
				models.UserWhere.Email.EQ(input.Body.Email),
			).OneG(ctx)
			if err != nil || invitee == nil {
				return nil, ErrorMap.GetErrorResponse(
					Err400_ChatChannelInviteeNotUser,
					err,
				)
			}
			if invitee.ID == input.UserId {
				return nil, ErrorMap.GetErrorResponse(
					Err400_ChatChannelInviteeOwnsChatGroup,
					err,
				)
			}
			// 4. Test if association between user and channel already exists
			foundChatUser, err := models.FindChatUserG(ctx, input.ChatChannelId, invitee.ID)
			if err != nil && foundChatUser != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err500_UnknownError,
					err,
				)
			}
			if foundChatUser != nil {
				if !foundChatUser.Disabled {
					return nil, ErrorMap.GetErrorResponse(
						Err400_ChatChannelAlreadyJoined,
						errors.New("activate association between invitee and requested private chat channel found"),
					)
				}
			} else {
				// create new one
				chatUser := models.ChatUser{
					ChatID:   input.ChatChannelId,
					UserID:   invitee.ID,
					Disabled: true,
				}
				if err := chatUser.InsertG(ctx, boil.Infer()); err != nil {
					return nil, ErrorMap.GetErrorResponse(
						Err500_UnableCreateChatUser,
						err,
					)
				}
			}
			// 5. Send invitation email
			invitor, _ := models.FindUserG(ctx, input.UserId)
			token, _ := jwt.GenerateToken(
				invitor,
				jwt.TokenActionInvitationToPrivateChat,
				jwt.ExtraClaims{
					"inviteeId":     invitee.ID,
					"chatChannelId": input.ChatChannelId,
				},
			)
			var html bytes.Buffer
			emailService.InviteToPrivateChatGroup(
				invitee.FullName,
				chatChannel.Title,
				chatGroup.Title,
				fmt.Sprintf(
					"%s/forms/accept-private-chat-invitation?token=%s",
					os.Getenv("WEB_CLIENT_URL"),
					token.Token,
				),
				&html,
			)
			if err := email.Send(ctx, email.EmailDTO{
				From:     "no-reply@quible.io",
				To:       invitee.Email,
				Subject:  "Invitation to join private chat channel",
				HTMLBody: html.String(),
			}); err != nil {
				return nil, ErrorMap.GetErrorResponse(
					Err424_UnableToSendEmail,
					err,
				)
			}
			return nil, nil
		},
	)
}
