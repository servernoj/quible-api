package v1_test

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/store"
	"github.com/stretchr/testify/assert"
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
				Headers: []any{},
			},
			Response: TCResponse{
				Status:    http.StatusUnauthorized,
				ErrorCode: misc.Of(v1.Err401_InvalidAccessToken),
			},
		},
		"FailureInvalidAccessToken": TCData{
			Description: "Failure due to an invalid token in the authorization header",
			Request: TCRequest{
				Headers: []any{
					"Authorization: Bearer invalid-token",
				},
			},
			Response: TCResponse{
				Status:    http.StatusUnauthorized,
				ErrorCode: misc.Of(v1.Err401_InvalidAccessToken),
			},
		},
		"Success": TCData{
			Description: "Success + validating returned user against DB",
			Request: TCRequest{
				Headers: []any{
					// User A
					fmt.Sprintf("Authorization: Bearer %s", GetToken(t, "9bef41ed-fb10-4791-b02e-96b372c09466")),
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
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			response := suite.TestAPI.Get("/api/v1/user", scenario.Request.Headers...)
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
