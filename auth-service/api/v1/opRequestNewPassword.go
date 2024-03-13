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
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/rs/zerolog/log"
)

type RequestNewPasswordInput struct {
	Body struct {
		Email string `json:"email" format:"email"`
	}
}

type RequestNewPasswordOutput struct {
}

func (impl *VersionedImpl) RegisterRequestNewPassword(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "post-request-new-password",
				Summary:     "Request new password",
				Description: "Recover forgotten password by submitting associated email address",
				Method:      http.MethodPost,
				Errors: []int{
					http.StatusBadRequest,
					http.StatusFailedDependency,
				},
				DefaultStatus: http.StatusAccepted,
				Tags:          []string{"user", "public"},
				Path:          "/user/request-new-password",
			},
		),
		func(ctx context.Context, input *RequestNewPasswordInput) (*RequestNewPasswordOutput, error) {
			// 0. Dependences
			deps := impl.Deps.GetContext("opRequestNewPassword")
			db := deps.Get("db").(*sql.DB)
			// 1. Locate user record based on provided email
			user, err := models.Users(models.UserWhere.Email.EQ(input.Body.Email)).One(ctx, db)
			if err != nil {
				// We intentionally don't send HTTP error for security reasons
				log.Error().Str("email", input.Body.Email).Msg("Email not registered")
				return nil, nil
			}
			// 2. Generate Password Reset email
			token, _ := jwt.GenerateToken(user, jwt.TokenActionPasswordReset, nil)
			var html bytes.Buffer
			emailService.PasswordReset(
				user.FullName,
				fmt.Sprintf(
					"%s/forms/password-reset?token=%s",
					os.Getenv("WEB_CLIENT_URL"),
					token.String(),
				),
				&html,
			)
			// 3. Send out generated email
			if emailSender, ok := deps.Get("mailer").(email.EmailSender); ok {
				if err := emailSender.SendEmail(ctx, email.EmailPayload{
					From:     "no-reply@quible.io",
					To:       user.Email,
					Subject:  "Password reset",
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
			// 4. Return empty response to indicate success
			return nil, nil
		},
	)
}
