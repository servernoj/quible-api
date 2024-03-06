package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type ResetPasswordInput struct {
	Body struct {
		Token           string  `json:"token" pattern:"^[^.]+([.][^.]+){2}$"`
		Step            string  `json:"step" enum:"validate,define"`
		Password        *string `json:"password,omitempty" doc:"at least 6 characters long"`
		ConfirmPassword *string `json:"confirmPassword,omitempty"`
	}
	hashedPassword string
}

func (input *ResetPasswordInput) Resolve(ctx huma.Context) (errs []error) {
	// 1. Expect password fields set only when Step is "define"
	if input.Body.Step == "define" {
		// 1a. Both password fields must be present...
		if input.Body.Password == nil {
			errs = append(errs, &huma.ErrorDetail{
				Message:  "Missing `password` value",
				Location: "body.password",
				Value:    input.Body.Password,
			})
		}
		if input.Body.ConfirmPassword == nil {
			errs = append(errs, &huma.ErrorDetail{
				Message:  "Missing `confirmPassword` value",
				Location: "body.confirmPassword",
				Value:    input.Body.ConfirmPassword,
			})
		}
		if len(errs) > 0 {
			return
		}
		// 1b. They should hold the same value...
		if *input.Body.Password != *input.Body.ConfirmPassword {
			errs = append(errs, &huma.ErrorDetail{
				Message:  "Password should match its confirmation value",
				Location: "body.confirmPassword",
				Value:    input.Body,
			})
			return
		}
		// 1c. Password should satisfy complexity requirements confirmed by the resolver...
		passwordResolver := &PasswordResolver{
			Password: *input.Body.Password,
		}
		if passwordResolveErrors := passwordResolver.Resolve(ctx); len(passwordResolveErrors) > 0 {
			errs = passwordResolveErrors
			return
		}
		// 1d. Compute and set HashedPassword field to be used in operation
		input.hashedPassword = passwordResolver.hashedPassword
	}
	return nil
}

type ResetPasswordOutput struct {
}

func (impl *VersionedImpl) RegisterResetPassword(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "post-reset-password",
				Summary:     "Reset password",
				Description: "Accept new password and set it for the user identified by the provided token",
				Method:      http.MethodPost,
				Errors: []int{
					http.StatusUnauthorized,
					http.StatusExpectationFailed,
				},
				DefaultStatus: http.StatusOK,
				Tags:          []string{"user", "public"},
				Path:          "/user/password-reset",
			},
		),
		func(ctx context.Context, input *ResetPasswordInput) (*ResetPasswordOutput, error) {
			// 1. Identify `userId` from the provided activation token
			tokenClaims, err := jwt.VerifyJWT(input.Body.Token, jwt.TokenActionPasswordReset)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err401_InvalidPasswordResetToken, err)
			}
			userId := tokenClaims["userId"].(string)
			// 2. Retrieve associated user record
			user, err := models.FindUserG(ctx, userId)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err417_UnableToAssociateUser, err)
			}
			// 3. Branch based on the value of `input.Body.Step`
			if input.Body.Step == "define" {
				user.HashedPassword = input.hashedPassword
				if _, err := user.UpdateG(ctx, boil.Infer()); err != nil {
					return nil, ErrorMap.GetErrorResponse(Err500_UnableToResetPassword, err)
				}
			}
			return nil, nil
		},
	)
}
