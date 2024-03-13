package v1_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/h2non/gock"
	v1 "github.com/quible-io/quible-api/app-service/api/v1"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/suite"
)

func (tc *TestCases) TestListChatGroups(t *testing.T) {
	// 1. Import users from CSV file
	db := tc.DBStore.RetrieveDB(t.Name())
	deps := tc.ServiceAPI.SetContext("opListChatGroups")
	deps.Set("db", db)
	if err := suite.InsertFromCSV(db, "users", UsersCSV); err != nil {
		t.Fatalf("unable to import users data from CSV: %s", err)
	}
	if err := suite.InsertFromCSV(db, "chats", ChatsCSV); err != nil {
		t.Fatalf("unable to import chat data from CSV: %s", err)
	}
	authServiceHost := "http://localhost"
	// 2. Define test scenarios
	testCases := libAPI.TCScenarios{
		"FailureUnreachableAuthService": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure on unavailable auth-service",
				Request:     libAPI.TCRequest{},
				Response: libAPI.TCResponse{
					Status:    http.StatusUnauthorized,
					ErrorCode: v1.Err401_InvalidAccessToken.Ptr(),
				},
			}
		},
		"FailureOnInvalidAuthorization": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure on unauthorized request due to invalid Bearer token",
				Request: libAPI.TCRequest{
					Args: []any{
						"Authorization: invalid",
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusUnauthorized,
					ErrorCode: v1.Err401_AuthServiceError.Ptr(),
				},
				PreHook: func(t *testing.T) any {
					gock.New("").
						Get("/api/v1/user").
						MatchHeader("Authorization", "invalid").
						Reply(http.StatusUnauthorized).
						JSON(nil)
					return nil
				},
				PostHook: func(t *testing.T, a any) {
					gock.Off()
				},
			}
		},
		"SuccessUserWithChatGroups": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Success with non-empty list of chat groups in response",
				Request: libAPI.TCRequest{
					Args: []any{
						"Authorization: valid",
					},
				},
				Envs: libAPI.TCEnv{
					"ENV_URL_AUTH_SERVICE": authServiceHost,
				},
				Response: libAPI.TCResponse{
					Status: http.StatusOK,
				},
				PreHook: func(t *testing.T) any {
					gock.New(authServiceHost).
						Get("/api/v1/user").
						MatchHeader("Authorization", "valid").
						Reply(http.StatusOK).
						JSON(map[string]string{
							// User A
							"id": "9bef41ed-fb10-4791-b02e-96b372c09466",
						})
					return nil
				},
				PostHook: func(t *testing.T, a any) {
					gock.Off()
				},
				ExtraTests: []libAPI.TCExtraTest{
					func(_ libAPI.TCRequest, res *httptest.ResponseRecorder) bool {
						var chats models.ChatSlice
						if err := json.NewDecoder(res.Result().Body).Decode(&chats); err != nil {
							return false
						}
						if len(chats) != 2 {
							return false
						}
						for _, chat := range chats {
							if found, err := models.ChatExists(context.Background(), db, chat.ID); err != nil || !found {
								return false
							}
						}
						return true
					},
				},
			}
		},
		"SuccessUserWithoutChatGroups": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Success with empty list of chat groups in response",
				Request: libAPI.TCRequest{
					Args: []any{
						"Authorization: valid",
					},
				},
				Envs: libAPI.TCEnv{
					"ENV_URL_AUTH_SERVICE": authServiceHost,
				},
				Response: libAPI.TCResponse{
					Status: http.StatusOK,
				},
				PreHook: func(t *testing.T) any {
					gock.New(authServiceHost).
						Get("/api/v1/user").
						MatchHeader("Authorization", "valid").
						Reply(http.StatusOK).
						JSON(map[string]string{
							// User D
							"id": "00e52081-0452-49ba-adbc-34612d3f1259",
						})
					return nil
				},
				PostHook: func(t *testing.T, a any) {
					gock.Off()
				},
				ExtraTests: []libAPI.TCExtraTest{
					func(_ libAPI.TCRequest, res *httptest.ResponseRecorder) bool {
						var chats models.ChatSlice
						if err := json.NewDecoder(res.Result().Body).Decode(&chats); err != nil {
							return false
						}
						return len(chats) == 0
					},
				},
			}
		},
		"FailureOnUnknownUser": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure on a request on behalf of unknown user",
				Request: libAPI.TCRequest{
					Args: []any{
						"Authorization: valid",
					},
				},
				Envs: libAPI.TCEnv{
					"ENV_URL_AUTH_SERVICE": authServiceHost,
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusInternalServerError,
					ErrorCode: v1.Err500_UnknownError.Ptr(),
				},
				PreHook: func(t *testing.T) any {
					gock.New(authServiceHost).
						Get("/api/v1/user").
						MatchHeader("Authorization", "valid").
						Reply(http.StatusOK).
						JSON(map[string]string{
							"id": "unknown-user-id",
						})
					return nil
				},
				PostHook: func(t *testing.T, a any) {
					gock.Off()
				},
			}
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(tc.TestAPI, http.MethodGet, "/chat/groups"))
	}
}
