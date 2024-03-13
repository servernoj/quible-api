package v1_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/suite"
)

func (tc *TestCases) TestUserLogin(t *testing.T) {
	// 1. Import users from CSV file
	db := tc.DBStore.RetrieveDB(t.Name())
	deps := tc.ServiceAPI.SetContext("opUserLogin")
	deps.Set("db", db)
	if err := suite.InsertFromCSV(db, "users", UsersCSV); err != nil {
		t.Fatalf("unable to import test data from CSV: %s", err)
	}
	// 2. Define test scenarios
	testCases := libAPI.TCScenarios{
		"Success": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "login with correct credentials and expect success",
				Request: libAPI.TCRequest{
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
				Response: libAPI.TCResponse{
					Status: http.StatusOK,
				},
				ExtraTests: []libAPI.TCExtraTest{
					func(req libAPI.TCRequest, response *httptest.ResponseRecorder) bool {
						email := req.Params["email"].(string)
						var responseBody v1.UserTokens
						if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
							return false
						}
						user, err := models.Users(
							models.UserWhere.Email.EQ(email),
						).One(context.Background(), db)
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
			}
		},
		"InvalidCredentials": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "login with incorrect credentials and expect error",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"email":    "userA@gmail.com",
							"password": "wrong password",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusUnauthorized,
					ErrorCode: v1.Err401_InvalidCredentials.Ptr(),
				},
			}
		},
		"UnknownUser": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "login non-existing user and expect error",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"email":    "unknown-user@gmail.com",
							"password": "does-not-matter",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_EmailNotRegistered.Ptr(),
				},
			}
		},
		"InvalidEmailFormat": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "login with improperly formatted email",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"email":    "not-an-email-address",
							"password": "password",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_InvalidEmailFormat.Ptr(),
				},
			}
		},
		"UnactivatedUser": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "login with unactivated user",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"email":    "UserC@gmail.com",
							"password": "does-not-matter",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusUnauthorized,
					ErrorCode: v1.Err401_UserNotActivated.Ptr(),
				},
			}
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(tc.TestAPI, http.MethodPost, "/login"))
	}
}
