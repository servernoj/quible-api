package v1_test

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	"github.com/quible-io/quible-api/lib/env"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/store"
	"github.com/stretchr/testify/assert"
)

func (suite *TestCases) TestActivateUser() {
	t := suite.T()
	// 1. Import users from CSV file
	store.InsertFromCSV(t, "users", UsersCSV)
	// 2. Load environment variables
	env.Setup()
	// 3. Define test scenarios
	testCases := TCScenarios{
		"FailureMalformedToken": TCData{
			Description: "Failure due to a malformed token in the request body",
			Request: TCRequest{
				Body: map[string]any{
					"token": "purely-invalid-token",
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_InvalidOrMalformedToken),
			},
		},
		"FailureTokenImproperAction": TCData{
			Description: "Failure due to token made for improper action",
			Request: TCRequest{
				Body: map[string]any{
					"token": GetToken(t, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess),
				},
			},
			Response: TCResponse{
				Status:    http.StatusUnauthorized,
				ErrorCode: misc.Of(v1.Err401_InvalidActivationToken),
			},
		},
		"FailureTokenExpired": TCData{
			Description: "Failure due to an expired token",
			Request: TCRequest{
				Body: map[string]any{
					// expired on 1/1/2000 at 00:00 UTC
					"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk0NjY4NDgwMCwidXNlcklkIjoiOWJlZjQxZWQtZmIxMC00NzkxLWIwMmUtOTZiMzcyYzA5NDY2IiwiYWN0aW9uIjoiQWN0aXZhdGUiLCJleHRyYUNsYWltcyI6bnVsbH0.vkCD8TY30l4LJOir8tXuxgpLO_XDzY3pv7k9CazQw9g",
				},
			},
			Response: TCResponse{
				Status:    http.StatusUnauthorized,
				ErrorCode: misc.Of(v1.Err401_InvalidActivationToken),
			},
		},
		"FailureTokenBadSignature": TCData{
			Description: "Failure due to invalid token signature",
			Request: TCRequest{
				Body: map[string]any{
					// signed with "wrong-secret", exp set to 1/1/2030 at 00:00 UTC
					"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4OTM0NTYwMDAsInVzZXJJZCI6IjliZWY0MWVkLWZiMTAtNDc5MS1iMDJlLTk2YjM3MmMwOTQ2NiIsImFjdGlvbiI6IkFjdGl2YXRlIiwiZXh0cmFDbGFpbXMiOm51bGx9.ZlDe3hRl5MfpbkHKyFfUDJ_Zsv140KBjWZ7wBBd6-f0",
				},
			},
			Response: TCResponse{
				Status:    http.StatusUnauthorized,
				ErrorCode: misc.Of(v1.Err401_InvalidActivationToken),
			},
		},
		"Success": TCData{
			Description: "Success with confirmation of the activation in DB",
			Request: TCRequest{
				Body: map[string]any{
					// non-activated user
					"token": GetToken(t, "c6174e8a-e12f-4d64-a4fe-a3b0c081bd31", jwt.TokenActionActivate),
				},
			},
			Response: TCResponse{
				Status: http.StatusOK,
			},
			ExtraTests: []TCExtraTest{
				func(t TCRequest, rr *httptest.ResponseRecorder) bool {
					userId := "c6174e8a-e12f-4d64-a4fe-a3b0c081bd31"
					user, err := models.FindUserG(context.Background(), userId)
					if err != nil {
						return false
					}
					return user.ActivatedAt.Ptr() != nil
				},
			},
		},
	}
	// 4. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			response := suite.TestAPI.Post("/api/v1/user/activate", scenario.Request.Body)
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
