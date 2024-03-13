package v1

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/quible-io/quible-api/auth-service/services/emailService"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/quible-io/quible-api/lib/models"
)

type InviteUserInput struct {
	AuthorizationHeaderResolver
	Body struct {
		Email    string `json:"email" format:"email"`
		FullName string `json:"full_name" minLength:"1"`
	}
}

type InviteUserOutput struct {
}

func (impl *VersionedImpl) RegisterInviteUser(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "post-invite-user",
				Summary:     "Invite new user",
				Description: "Invite new user by sending an invitation email",
				Method:      http.MethodPost,
				Errors: []int{
					http.StatusBadRequest,
					http.StatusFailedDependency,
				},
				DefaultStatus: http.StatusOK,
				Tags:          []string{"user", "private"},
				Path:          "/user/invite",
			},
		),
		func(ctx context.Context, input *InviteUserInput) (*InviteUserOutput, error) {
			// 0. Dependences
			deps := impl.Deps.GetContext("opInviteUser")
			db := deps.Get("db").(*sql.DB)
			// 1. Identify if a user already exist
			user, _ := models.Users(models.UserWhere.Email.EQ(input.Body.Email)).One(ctx, db)
			if user != nil {
				return nil, ErrorMap.GetErrorResponse(Err400_UserWithEmailExists)
			}
			// 2. Generate and send invitation email
			var html bytes.Buffer
			emailService.UserInvitation(
				input.Body.FullName,
				fmt.Sprintf(
					"%s/forms/register",
					os.Getenv("WEB_CLIENT_URL"),
				),
				&html,
			)
			if emailSender, ok := deps.Get("mailer").(email.EmailSender); ok {
				if err := emailSender.SendEmail(ctx, email.EmailPayload{
					From:     "no-reply@quible.io",
					To:       input.Body.Email,
					Subject:  "Invitation to register Quible account",
					HTMLBody: html.String(),
				}); err != nil {
					return nil, ErrorMap.GetErrorResponse(Err424_UnableToSendEmail, err)
				}
			} else {
				return nil, ErrorMap.GetErrorResponse(
					Err424_UnableToSendEmail,
					errors.New("email client unavailable"),
				)
			}
			// 3. Return empty response to indicate success
			return nil, nil
		},
	)
}
