package v1_test

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/suite"
)

func (tc *TestCases) TestActivateUser(t *testing.T) {
	// 1. Import users from CSV file
	db := tc.DBStore.RetrieveDB(t.Name())
	deps := tc.ServiceAPI.SetContext("opActivateUser")
	deps.Set("db", db)
	if err := suite.InsertFromCSV(db, "users", UsersCSV); err != nil {
		t.Fatalf("unable to import test data from CSV: %s", err)
	}
	// 2. Define test scenarios
	testCases := libAPI.TCScenarios{
		"FailureMalformedToken": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure due to a malformed token in the request body",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"token": "purely-invalid-token",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_InvalidOrMalformedToken.Ptr(),
				},
			}
		},
		"FailureTokenImproperAction": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure due to token made for improper action",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"token": suite.GetToken(t, db, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess),
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusUnauthorized,
					ErrorCode: v1.Err401_InvalidActivationToken.Ptr(),
				},
			}
		},
		"FailureTokenExpired": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure due to an expired token",
				Envs: libAPI.TCEnv{
					"ENV_JWT_SECRET": "secret",
				},
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							// expired on 1/1/2000 at 00:00 UTC
							"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk0NjY4NDgwMCwidXNlcklkIjoiOWJlZjQxZWQtZmIxMC00NzkxLWIwMmUtOTZiMzcyYzA5NDY2IiwiYWN0aW9uIjoiQWN0aXZhdGUiLCJleHRyYUNsYWltcyI6bnVsbH0.R6KGfdxRtYcBNrP3lc5nDgzAoVlYU58w-XUKEmDZ5P8",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusUnauthorized,
					ErrorCode: v1.Err401_InvalidActivationToken.Ptr(),
				},
			}
		},
		"FailureTokenBadSignature": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure due to invalid token signature",
				Envs: libAPI.TCEnv{
					"ENV_JWT_SECRET": "secret",
				},
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							// signed with "wrong-secret", exp set to 1/1/2030 at 00:00 UTC
							"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4OTM0NTYwMDAsInVzZXJJZCI6IjliZWY0MWVkLWZiMTAtNDc5MS1iMDJlLTk2YjM3MmMwOTQ2NiIsImFjdGlvbiI6IkFjdGl2YXRlIiwiZXh0cmFDbGFpbXMiOm51bGx9.ZlDe3hRl5MfpbkHKyFfUDJ_Zsv140KBjWZ7wBBd6-f0",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusUnauthorized,
					ErrorCode: v1.Err401_InvalidActivationToken.Ptr(),
				},
			}
		},
		"FailureTokenBadUser": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure due to non-existing user referenced in the token",
				Envs: libAPI.TCEnv{
					"ENV_JWT_SECRET": "secret",
				},
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							// token for non-existing user "00000000-0000-0000-0000-000000000000" with secret `secret`
							"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4OTM0NTYwMDAsInVzZXJJZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImFjdGlvbiI6IkFjdGl2YXRlIiwiZXh0cmFDbGFpbXMiOm51bGx9.t_pSHH5N4ptUCFKbkXZrnHzsTGkJOCH6bp3aAK-W_GU",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusExpectationFailed,
					ErrorCode: v1.Err417_UnableToAssociateUser.Ptr(),
				},
			}
		},
		"Success": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Success with confirmation of the activation in DB",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							// non-activated user
							"token": suite.GetToken(t, db, "c6174e8a-e12f-4d64-a4fe-a3b0c081bd31", jwt.TokenActionActivate),
						},
					},
				},
				Response: libAPI.TCResponse{
					Status: http.StatusOK,
				},
				ExtraTests: []libAPI.TCExtraTest{
					func(req libAPI.TCRequest, res *httptest.ResponseRecorder) bool {
						userId := "c6174e8a-e12f-4d64-a4fe-a3b0c081bd31"
						user, err := models.FindUser(context.Background(), db, userId)
						if err != nil {
							return false
						}
						return user.ActivatedAt.Ptr() != nil
					},
				},
			}
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(tc.TestAPI, http.MethodPost, "/user/activate"))
	}
}
