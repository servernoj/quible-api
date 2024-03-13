package v1

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/quible-io/quible-api/auth-service/services/emailService"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type CreateUserInput struct {
	AutoActivate bool `query:"auto-activation" default:"false"`
	Body         struct {
		PasswordResolver
		Username string `json:"username"`
		Email    string `json:"email" format:"email"`
		FullName string `json:"full_name" minLength:"1"`
		Phone    string `json:"phone" pattern:"^[0-9() +-]{10,}$"`
	}
}

type CreateUserOutput struct {
	Body UserSimplified
}

func (impl *VersionedImpl) RegisterCreateUser(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "post-create-user",
				Summary:     "Create new user",
				Description: "Register new non-activated user record based on provided information",
				Method:      http.MethodPost,
				Errors: []int{
					http.StatusBadRequest,
				},
				DefaultStatus: http.StatusCreated,
				Tags:          []string{"user", "public"},
				Path:          "/user",
			},
		),
		func(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
			// 0. Dependences
			deps := impl.Deps.GetContext("opCreateUser")
			db := deps.Get("db").(*sql.DB)
			// 1. Check if an activated user with the same username or email exists
			foundUser, _ := models.Users(
				qm.Or2(models.UserWhere.Email.EQ(input.Body.Email)),
				qm.Or2(models.UserWhere.Username.EQ(input.Body.Username)),
			).One(ctx, db)
			if foundUser != nil && foundUser.ActivatedAt.Ptr() != nil {
				return nil, ErrorMap.GetErrorResponse(Err400_UserWithEmailOrUsernameExists)
			}
			// 2. Depending of whether [non-activated] user exists
			user := &models.User{}
			if foundUser != nil {
				user = foundUser
			}
			user.Email = input.Body.Email
			user.Phone = input.Body.Phone
			user.Username = input.Body.Username
			user.FullName = input.Body.FullName
			user.HashedPassword = input.Body.hashedPassword
			if len(user.ID) > 0 {
				// 2a. We update existing DB record
				if _, err := user.Update(ctx, db, boil.Infer()); err != nil {
					return nil, ErrorMap.GetErrorResponse(Err500_UnableToUpdateUser, err)
				}
			} else {
				// 2b. We create new DB record
				if err := user.Insert(ctx, db, boil.Infer()); err != nil {
					return nil, ErrorMap.GetErrorResponse(Err500_UnableToRegister, err)
				}
			}
			// 3. Branching...
			if input.AutoActivate && os.Getenv("IS_DEV") == "1" {
				// 3a. Auto-generate if requested on `dev` (and `local`) environments
				user.ActivatedAt = null.TimeFrom(time.Now())
				if _, err := user.Update(ctx, db, boil.Infer()); err != nil {
					return nil, ErrorMap.GetErrorResponse(Err500_UnableToUpdateUser, err)
				}
			} else {
				// 3b. Send activation email otherwise
				token, _ := jwt.GenerateToken(user, jwt.TokenActionActivate, nil)
				var html bytes.Buffer
				emailService.UserActivation(
					user.FullName,
					fmt.Sprintf(
						"%s/forms/activation?token=%s",
						os.Getenv("WEB_CLIENT_URL"),
						token.String(),
					),
					&html,
				)
				if emailSender, ok := deps.Get("mailer").(email.EmailSender); ok {
					if err := emailSender.SendEmail(ctx, email.EmailPayload{
						From:     "no-reply@quible.io",
						To:       user.Email,
						Subject:  "Activate your Quible account",
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
			}
			// 4. Prepare and return the response
			response := &CreateUserOutput{
				Body: UserSimplified{
					ID:       user.ID,
					Username: user.Username,
					Email:    user.Email,
					Phone:    user.Phone,
					FullName: user.FullName,
				},
			}
			return response, nil
		},
	)
}
