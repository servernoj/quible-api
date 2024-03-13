package v1

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type ActivateUserInput struct {
	Body struct {
		Token string `json:"token" pattern:"^[^.]+([.][^.]+){2}$"`
	}
}

type ActivateUserOutput struct {
}

func (impl *VersionedImpl) RegisterActivateUser(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "post-activate-user",
				Summary:     "Activate user account",
				Description: "Update user record in response to clicking the link in activation email",
				Method:      http.MethodPost,
				Errors: []int{
					http.StatusUnauthorized,
					http.StatusExpectationFailed,
				},
				DefaultStatus: http.StatusOK,
				Tags:          []string{"user", "public"},
				Path:          "/user/activate",
			},
		),
		func(ctx context.Context, input *ActivateUserInput) (*ActivateUserOutput, error) {
			// 0. Dependences
			deps := impl.Deps.GetContext("opActivateUser")
			db := deps.Get("db").(*sql.DB)
			// 1. Identify `userId` from the provided activation token
			tokenClaims, err := jwt.VerifyJWT(input.Body.Token, jwt.TokenActionActivate)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err401_InvalidActivationToken, err)
			}
			userId := tokenClaims["userId"].(string)
			// 2. Locate user in DB
			user, err := models.FindUser(ctx, db, userId)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err417_UnableToAssociateUser, err)
			}
			// 3. Once found, update user record to store activation timestamp
			user.ActivatedAt = null.TimeFrom(time.Now())
			_, err = user.Update(ctx, db, boil.Infer())
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err500_UnableToActivateUser, err)
			}
			return nil, nil
		},
	)
}
