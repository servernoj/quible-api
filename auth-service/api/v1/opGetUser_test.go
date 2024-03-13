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
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/suite"
)

func (tc *TestCases) TestGetUser(t *testing.T) {
	// 1. Import users from CSV file
	db := tc.DBStore.RetrieveDB(t.Name())
	deps := tc.ServiceAPI.SetContext("opGetUser")
	deps.Set("db", db)
	if err := suite.InsertFromCSV(db, "users", UsersCSV); err != nil {
		t.Fatalf("unable to import test data from CSV: %s", err)
	}
	// 2. Define test scenarios
	testCases := libAPI.TCScenarios{
		"FailureMissingAuthorizationHeader": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure due to missing authorization header",
				Request: libAPI.TCRequest{
					Args: []any{},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusUnauthorized,
					ErrorCode: v1.Err401_InvalidAccessToken.Ptr(),
				},
			}
		},
		"FailureInvalidAccessToken": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure due to an invalid token in the authorization header",
				Request: libAPI.TCRequest{
					Args: []any{
						"Authorization: Bearer invalid-token",
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusUnauthorized,
					ErrorCode: v1.Err401_InvalidAccessToken.Ptr(),
				},
			}
		},
		"FailureTokenBadUser": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure due to non-existing user referenced in the token",
				Request: libAPI.TCRequest{
					Args: []any{
						// access token for non-existing user "00000000-0000-0000-0000-000000000000" with secret `secret`
						fmt.Sprintf(
							"Authorization: Bearer %s",
							"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4OTM0NTYwMDAsImp0aSI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImlzcyI6IlF1aWJsZSIsInVzZXJJZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImFjdGlvbiI6IkFjY2VzcyIsImV4dHJhQ2xhaW1zIjpudWxsfQ.3-JimplIBm4W7YBcgZGp2Efmh5Thv8bQZM5ggrAW0RY",
						),
					},
				},
				Envs: libAPI.TCEnv{
					"ENV_JWT_SECRET": "secret",
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusUnauthorized,
					ErrorCode: v1.Err401_InvalidAccessToken.Ptr(),
				},
			}
		},
		"Success": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Success + validating returned user against DB",
				Request: libAPI.TCRequest{
					Args: []any{
						// User A
						fmt.Sprintf(
							"Authorization: Bearer %s",
							suite.GetToken(t, db, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess),
						),
					},
				},
				Response: libAPI.TCResponse{
					Status: http.StatusOK,
				},
				ExtraTests: []libAPI.TCExtraTest{
					func(_ libAPI.TCRequest, res *httptest.ResponseRecorder) bool {
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
			}
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(tc.TestAPI, http.MethodGet, "/user"))
	}
}
