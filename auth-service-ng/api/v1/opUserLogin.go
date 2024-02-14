package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/quible-io/quible-api/auth-service-ng/api"
	"github.com/quible-io/quible-api/auth-service-ng/services/userService"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type LoginInput struct {
	Body struct {
		Email    string `json:"email" format:"email"`
		Password string `json:"password"`
	}
}

type UserTokens struct {
	AccessToken  string `json:"access_token" doc:"access token to be used to authenticate other API calls"`
	RefreshToken string `json:"refresh_token" doc:"refresh token to be used to renew/refresh access token without re-submitting user credentials"`
}

type LoginOutput struct {
	Body UserTokens
}

func (impl *VersionedImpl) RegisterLogin(api huma.API, vc api.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "post-login",
				Summary:     "Login user",
				Description: "Login user based on provided credentials (email/password)",
				Method:      http.MethodPost,
				Errors: []int{
					http.StatusBadRequest,
					http.StatusUnauthorized,
				},
				Tags: []string{"user", "public"},
				Path: "/login",
			},
		),
		func(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
			// 1. Locate the user in DB
			foundUser, err := models.Users(
				models.UserWhere.Email.EQ(input.Body.Email),
			).OneG(ctx)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err400_EmailNotRegistered, err)
			}
			// 2. Check if the user is activated
			if foundUser.ActivatedAt.Ptr() == nil {
				return nil, ErrorMap.GetErrorResponse(Err401_UserNotActivated)
			}
			// 3. Compare the stored password hash with the hash computed from the provided password
			us := userService.UserService{}
			if err := us.ValidatePassword(foundUser.HashedPassword, input.Body.Password); err != nil {
				return nil, ErrorMap.GetErrorResponse(Err401_InvalidCredentials, err)
			}
			// 4. Generate tokens
			accessToken, err := jwt.GenerateToken(foundUser, jwt.TokenActionAccess, nil)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err500_UnableToGenerateToken, err)
			}
			refreshToken, err := jwt.GenerateToken(foundUser, jwt.TokenActionRefresh, nil)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err500_UnableToGenerateToken, err)
			}
			// 5. Update user's record to reference freshly generated refresh token
			foundUser.Refresh = refreshToken.ID
			if _, err := foundUser.UpdateG(ctx, boil.Infer()); err != nil {
				return nil, ErrorMap.GetErrorResponse(Err500_UnableToUpdateUser, err)
			}
			// 6. Prepare and return the response
			response := &LoginOutput{}
			response.Body.AccessToken = accessToken.String()
			response.Body.RefreshToken = refreshToken.String()
			return response, nil
		},
	)
}
