package v1_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/h2non/gock"
	v1 "github.com/quible-io/quible-api/app-service/api/v1"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/store"
)

func (suite *TestCases) TestListChatGroups() {
	t := suite.T()
	// 1. Import data from CSV files
	store.InsertFromCSV(t, "users", UsersCSV)
	store.InsertFromCSV(t, "chats", ChatsCSV)
	// 2. Define test scenarios
	testCases := TCScenarios{
		"FailureUnreachableAuthService": TCData{
			Description: "Failure on unavailable auth-service",
			Request:     TCRequest{},
			Response: TCResponse{
				Status:    http.StatusUnauthorized,
				ErrorCode: misc.Of(v1.Err401_InvalidAccessToken),
			},
		},
		"FailureOnInvalidAuthorization": TCData{
			Description: "Failure on unauthorized request due to invalid Bearer token",
			Request: TCRequest{
				Args: []any{
					"Authorization: invalid",
				},
			},
			Response: TCResponse{
				Status:    http.StatusUnauthorized,
				ErrorCode: misc.Of(v1.Err401_AuthServiceError),
			},
			PreHook: func(t *testing.T) any {
				gock.New(os.Getenv("ENV_URL_AUTH_SERVICE")).
					Get("/api/v1/user").
					MatchHeader("Authorization", "invalid").
					Reply(http.StatusUnauthorized).
					JSON(nil)
				return nil
			},
			PostHook: func(t *testing.T, a any) {
				gock.Off()
			},
		},
		"SuccessUserWithChatGroups": TCData{
			Description: "Success with non-empty list of chat groups in response",
			Request: TCRequest{
				Args: []any{
					"Authorization: valid",
				},
			},
			Response: TCResponse{
				Status: http.StatusOK,
			},
			PreHook: func(t *testing.T) any {
				gock.New(os.Getenv("ENV_URL_AUTH_SERVICE")).
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
			ExtraTests: []TCExtraTest{
				func(_ TCRequest, res *httptest.ResponseRecorder) bool {
					var chats models.ChatSlice
					if err := json.NewDecoder(res.Result().Body).Decode(&chats); err != nil {
						return false
					}
					if len(chats) != 2 {
						return false
					}
					for _, chat := range chats {
						if found, err := models.ChatExistsG(context.Background(), chat.ID); err != nil || !found {
							return false
						}
					}
					return true
				},
			},
		},
		"SuccessUserWithoutChatGroups": TCData{
			Description: "Success with empty list of chat groups in response",
			Request: TCRequest{
				Args: []any{
					"Authorization: valid",
				},
			},
			Response: TCResponse{
				Status: http.StatusOK,
			},
			PreHook: func(t *testing.T) any {
				gock.New(os.Getenv("ENV_URL_AUTH_SERVICE")).
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
			ExtraTests: []TCExtraTest{
				func(_ TCRequest, res *httptest.ResponseRecorder) bool {
					var chats models.ChatSlice
					if err := json.NewDecoder(res.Result().Body).Decode(&chats); err != nil {
						return false
					}
					return len(chats) == 0
				},
			},
		},
		"FailureOnUnknownUser": TCData{
			Description: "Failure on a request on behalf of unknown user",
			Request: TCRequest{
				Args: []any{
					"Authorization: valid",
				},
			},
			Response: TCResponse{
				Status:    http.StatusUnauthorized,
				ErrorCode: misc.Of(v1.Err401_UserNotFound),
			},
			PreHook: func(t *testing.T) any {
				gock.New(os.Getenv("ENV_URL_AUTH_SERVICE")).
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
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(suite.TestAPI, http.MethodGet, "/api/v1/chat/groups"))
	}
}
