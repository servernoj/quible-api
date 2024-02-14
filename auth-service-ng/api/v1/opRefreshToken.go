package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/quible-io/quible-api/auth-service-ng/api"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type RefreshTokenInput struct {
	Body struct {
		RefreshToken string `json:"refresh_token"`
	}
}

type RefreshTokenOutput struct {
	Body UserTokens
}

func (impl *VersionedImpl) RegisterRefreshToken(api huma.API, vc api.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "post-refresh-token",
				Summary:     "Refresh tokens",
				Description: "Use provided `refresh` token to generate new pair of access/refresh tokens",
				Method:      http.MethodPost,
				Errors: []int{
					http.StatusBadRequest,
					http.StatusUnauthorized,
				},
				DefaultStatus: http.StatusOK,
				Tags:          []string{"user", "public"},
				Path:          "/user/refresh",
			},
		),
		func(ctx context.Context, input *RefreshTokenInput) (*RefreshTokenOutput, error) {
			// 1. Process and validate provided refresh token
			claims, err := jwt.VerifyJWT(input.Body.RefreshToken, jwt.TokenActionRefresh)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err401_InvalidRefreshToken, err)
			}
			userId := claims["userId"].(string)
			// 2. Retrieve user record associated with the refresh token
			user, err := models.FindUserG(ctx, userId)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err401_InvalidRefreshToken, err)
			}
			// 3. Compare provided refresh token with the one registered for the identified user
			if refreshTokenId := claims["jti"].(string); user.Refresh != refreshTokenId {
				return nil, ErrorMap.GetErrorResponse(Err401_InvalidRefreshToken)
			}
			// 4. Generate tokens
			accessToken, err := jwt.GenerateToken(user, jwt.TokenActionAccess, nil)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err500_UnableToGenerateToken, err)
			}
			refreshToken, err := jwt.GenerateToken(user, jwt.TokenActionRefresh, nil)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err500_UnableToGenerateToken, err)
			}
			// 5. Update user's record to reference freshly generated refresh token
			user.Refresh = refreshToken.ID
			if _, err := user.UpdateG(ctx, boil.Infer()); err != nil {
				return nil, ErrorMap.GetErrorResponse(Err500_UnableToUpdateUser, err)
			}
			// 6. Return both [newly generated] access and refresh tokens
			response := &RefreshTokenOutput{
				Body: UserTokens{
					AccessToken:  accessToken.String(),
					RefreshToken: refreshToken.String(),
				},
			}
			return response, nil
		},
	)
}
