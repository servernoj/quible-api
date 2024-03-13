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
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (tc *TestCases) TestRefreshToken(t *testing.T) {
	// 1. Import users from CSV file
	db := tc.DBStore.RetrieveDB(t.Name())
	deps := tc.ServiceAPI.SetContext("opRefreshToken")
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
							"refresh_token": "purely-invalid-token",
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
							"refresh_token": suite.GetToken(t, db, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess),
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusUnauthorized,
					ErrorCode: v1.Err401_InvalidRefreshToken.Ptr(),
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
							"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk0NjY4NDgwMCwidXNlcklkIjoiOWJlZjQxZWQtZmIxMC00NzkxLWIwMmUtOTZiMzcyYzA5NDY2IiwiYWN0aW9uIjoiUmVmcmVzaCIsImV4dHJhQ2xhaW1zIjpudWxsfQ.awLM2azz_ooskVVOOTz5ecXBsUNygT8MgFgc8_hspLY",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusUnauthorized,
					ErrorCode: v1.Err401_InvalidRefreshToken.Ptr(),
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
							"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4OTM0NTYwMDAsImp0aSI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImlzcyI6IlF1aWJsZSIsInVzZXJJZCI6IjliZWY0MWVkLWZiMTAtNDc5MS1iMDJlLTk2YjM3MmMwOTQ2NiIsImFjdGlvbiI6IlJlZnJlc2giLCJleHRyYUNsYWltcyI6bnVsbH0.9J6YGWT9VzuozCKsmeBHqIDyEY3I_dlb1dO9lotaLL8",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusUnauthorized,
					ErrorCode: v1.Err401_InvalidRefreshToken.Ptr(),
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
							"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4OTM0NTYwMDAsImp0aSI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImlzcyI6IlF1aWJsZSIsInVzZXJJZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImFjdGlvbiI6IlJlZnJlc2giLCJleHRyYUNsYWltcyI6bnVsbH0.unkFdy1NsJ3jC-XtsF7DUP6HuGfV0TTlqICkHA5Pf_k",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusUnauthorized,
					ErrorCode: v1.Err401_InvalidRefreshToken.Ptr(),
				},
			}
		},
		"Success": func(t *testing.T) libAPI.TCData {
			ctx := context.Background()
			userId := "9bef41ed-fb10-4791-b02e-96b372c09466"
			refreshTokenSent := suite.GetToken(t, db, userId, jwt.TokenActionRefresh)
			claims, _ := jwt.VerifyJWT(refreshTokenSent, jwt.TokenActionRefresh)
			tokenId := claims["jti"].(string)
			user, _ := models.FindUser(ctx, db, userId)
			user.Refresh = tokenId
			_, _ = user.Update(ctx, db, boil.Whitelist("refresh"))
			return libAPI.TCData{
				Description: "Success with confirmation of user record modification",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"refresh_token": refreshTokenSent,
						},
					},
				},
				Response: libAPI.TCResponse{
					Status: http.StatusOK,
				},
				ExtraTests: []libAPI.TCExtraTest{
					func(req libAPI.TCRequest, res *httptest.ResponseRecorder) bool {
						user, err := models.FindUser(context.Background(), db, userId)
						if err != nil {
							return false
						}
						var responseBody v1.UserTokens
						if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
							return false
						}
						if refreshTokenSent == responseBody.RefreshToken {
							return false
						}
						claims, err := jwt.VerifyJWT(responseBody.RefreshToken, jwt.TokenActionRefresh)
						if err != nil {
							return false
						}
						userIdFromToken := claims["userId"].(string)
						refreshTokenId := claims["jti"].(string)
						if user.ID != userIdFromToken || user.Refresh != refreshTokenId {
							return false
						}
						return true
					},
				},
			}
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(tc.TestAPI, http.MethodPost, "/user/refresh"))
	}
}
