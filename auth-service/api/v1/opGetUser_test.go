package v1_test

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/store"
)

func (suite *TestCases) TestGetUser() {
	t := suite.T()
	// 1. Import users from CSV file
	store.InsertFromCSV(t, "users", UsersCSV)
	// 2. Define test scenarios
	testCases := TCScenarios{
		"FailureMissingAuthorizationHeader": TCData{
			Description: "Failure due to missing authorization header",
			Request: TCRequest{
				Args: []any{},
			},
			Response: TCResponse{
				Status:    http.StatusUnauthorized,
				ErrorCode: misc.Of(v1.Err401_InvalidAccessToken),
			},
		},
		"FailureInvalidAccessToken": TCData{
			Description: "Failure due to an invalid token in the authorization header",
			Request: TCRequest{
				Args: []any{
					"Authorization: Bearer invalid-token",
				},
			},
			Response: TCResponse{
				Status:    http.StatusUnauthorized,
				ErrorCode: misc.Of(v1.Err401_InvalidAccessToken),
			},
		},
		"FailureTokenBadUser": TCData{
			Description: "Failure due to non-existing user referenced in the token",
			Request: TCRequest{
				Args: []any{
					// access token for non-existing user "00000000-0000-0000-0000-000000000000" with secret `secret`
					fmt.Sprintf(
						"Authorization: Bearer %s",
						"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4OTM0NTYwMDAsImp0aSI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImlzcyI6IlF1aWJsZSIsInVzZXJJZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImFjdGlvbiI6IkFjY2VzcyIsImV4dHJhQ2xhaW1zIjpudWxsfQ.3-JimplIBm4W7YBcgZGp2Efmh5Thv8bQZM5ggrAW0RY",
					),
				},
			},
			Response: TCResponse{
				Status:    http.StatusUnauthorized,
				ErrorCode: misc.Of(v1.Err401_InvalidAccessToken),
			},
			PreHook: func(t *testing.T) any {
				t.Setenv("ENV_JWT_SECRET", "secret")
				return nil
			},
		},
		"Success": TCData{
			Description: "Success + validating returned user against DB",
			Request: TCRequest{
				Args: []any{
					// User A
					fmt.Sprintf("Authorization: Bearer %s", GetToken(t, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess)),
				},
			},
			Response: TCResponse{
				Status: http.StatusOK,
			},
			ExtraTests: []TCExtraTest{
				func(_ TCRequest, res *httptest.ResponseRecorder) bool {
					var got v1.UserSimplified
					if err := json.NewDecoder(res.Result().Body).Decode(&got); err != nil {
						return false
					}
					want := v1.UserSimplified{
						ID:       "9bef41ed-fb10-4791-b02e-96b372c09466",
						Username: "userA",
						Email:    "userA@gmail.com",
						Phone:    "1234567890",
						FullName: "User A",
					}
					return reflect.DeepEqual(got, want)
				},
			},
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(suite.TestAPI, http.MethodGet, "/api/v1/user"))
	}
}
