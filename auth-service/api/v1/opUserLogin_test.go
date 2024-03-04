package v1_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	"github.com/quible-io/quible-api/lib/env"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/store"
)

func (suite *TestCases) TestUserLogin() {
	t := suite.T()
	// 1. Import users from CSV file
	store.InsertFromCSV(t, "users", UsersCSV)
	// 2. Load environment variables
	env.Setup()
	// 3. Define test scenarios
	testCases := TCScenarios{
		"Success": TCData{
			Description: "login with correct credentials and expect success",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"email":    "userA@gmail.com",
						"password": "password",
					},
				},
				Params: map[string]any{
					"email": "userA@gmail.com",
				},
			},
			Response: TCResponse{
				Status: http.StatusOK,
			},
			ExtraTests: []TCExtraTest{
				func(req TCRequest, response *httptest.ResponseRecorder) bool {
					email := req.Params["email"].(string)
					var responseBody v1.UserTokens
					if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
						return false
					}
					user, err := models.Users(
						models.UserWhere.Email.EQ(email),
					).OneG(context.Background())
					if err != nil {
						return false
					}
					claims, err := jwt.VerifyJWT(responseBody.RefreshToken, jwt.TokenActionRefresh)
					if err != nil {
						return false
					}
					userId := claims["userId"].(string)
					refreshTokenId := claims["jti"].(string)
					if user.ID != userId || user.Refresh != refreshTokenId {
						return false
					}
					return true
				},
			},
		},
		"InvalidCredentials": TCData{
			Description: "login with incorrect credentials and expect error",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"email":    "userA@gmail.com",
						"password": "wrong password",
					},
				},
			},
			Response: TCResponse{
				Status:    http.StatusUnauthorized,
				ErrorCode: misc.Of(v1.Err401_InvalidCredentials),
			},
		},
		"UnknownUser": TCData{
			Description: "login non-existing user and expect error",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"email":    "unknown-user@gmail.com",
						"password": "does-not-matter",
					},
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_EmailNotRegistered),
			},
		},
		"InvalidEmailFormat": TCData{
			Description: "login with improperly formatted email",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"email":    "not-an-email-address",
						"password": "password",
					},
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_InvalidEmailFormat),
			},
		},
		"UnactivatedUser": TCData{
			Description: "login with unactivated user",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"email":    "UserC@gmail.com",
						"password": "does-not-matter",
					},
				},
			},
			Response: TCResponse{
				Status:    http.StatusUnauthorized,
				ErrorCode: misc.Of(v1.Err401_UserNotActivated),
			},
		},
	}
	// 4. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(suite.TestAPI, http.MethodPost, "/api/v1/login"))
	}
}
