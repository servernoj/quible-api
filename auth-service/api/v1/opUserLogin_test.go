package v1_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/store"
	"github.com/stretchr/testify/assert"
)

//go:embed TestData/users.csv
var users_as_csv string

func (suite *TestCases) TestUserLogin() {
	t := suite.T()
	t.Parallel()
	testCases := TCScenarios{
		"Success": TCData{
			Description: "login with correct credentials and expect success",
			Request: TCRequest{
				Body: map[string]any{
					"email":    "userA@gmail.com",
					"password": "password",
				},
			},
			Response: TCResponse{
				Status: http.StatusOK,
			},
			ExtraTests: []TCExtraTest{
				func(req TCRequest, response *httptest.ResponseRecorder) bool {
					email := req.Body["email"].(string)
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
				Body: map[string]any{
					"email":    "userA@gmail.com",
					"password": "wrong password",
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
				Body: map[string]any{
					"email":    "unknown-user@gmail.com",
					"password": "does-not-matter",
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
				Body: map[string]any{
					"email":    "not-an-email-address",
					"password": "password",
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_InvalidEmailFormat),
			},
		},
	}
	// 1. Import users from CSV file
	store.InsertFromCSV(t, "users", users_as_csv)
	// 2. Try different login scenarios
	for name, scenario := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			response := suite.TestAPI.Post("/api/v1/login", scenario.Request.Body)
			// response status
			assert.EqualValues(scenario.Response.Status, response.Code, "response status should match the expectation")
			// error code in case of error
			if scenario.Response.ErrorCode != nil {
				assert.Contains(
					response.Body.String(),
					strconv.Itoa(int(*scenario.Response.ErrorCode)),
					"error code should match expectation",
				)
			}
			// extra tests (if present)
			for _, fn := range scenario.ExtraTests {
				assert.True(
					fn(scenario.Request, response),
				)
			}
		})
	}
}
