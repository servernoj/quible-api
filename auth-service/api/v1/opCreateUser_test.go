package v1_test

import (
	"net/http"
	"strconv"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/store"
	"github.com/stretchr/testify/assert"
)

func (suite *TestCases) TestCreateUser() {
	t := suite.T()
	testCases := TCScenarios{
		"FailureOnActivatedWithExistingEmail": TCData{
			Description: "Failure of registering existing email on activated user",
			Request: TCRequest{
				Body: map[string]any{
					"username":  "new-username",
					"email":     "userA@gmail.com",
					"password":  "password",
					"phone":     "0123456789",
					"full_name": "existing email",
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_UserWithEmailOrUsernameExists),
			},
		},
		"FailureOnActivatedWithExistingUsername": TCData{
			Description: "Failure of registering existing username on activated user",
			Request: TCRequest{
				Body: map[string]any{
					"username":  "userA",
					"email":     "non-existent@gmail.com",
					"password":  "password",
					"phone":     "0123456789",
					"full_name": "existing username",
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_UserWithEmailOrUsernameExists),
			},
		},
		"FailureOnInvalidEmail": TCData{
			Description: "Failure of registering user with invalid email",
			Request: TCRequest{
				Body: map[string]any{
					"username":  "userD",
					"email":     "not-an-email-address",
					"password":  "password",
					"phone":     "0123456789",
					"full_name": "User D",
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_InvalidEmailFormat),
			},
		},
		"FailureOnInvalidPhone": TCData{
			Description: "Failure of registering user with invalid phone",
			Request: TCRequest{
				Body: map[string]any{
					"username":  "userD",
					"email":     "userD@gmail.com",
					"password":  "password",
					"phone":     "invalid",
					"full_name": "User D",
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_InvalidPhoneFormat),
			},
		},
		"Success": TCData{
			Description: "Failure of registering user with invalid phone",
			Request: TCRequest{
				Body: map[string]any{
					"username":  "userD",
					"email":     "userD@gmail.com",
					"password":  "password",
					"phone":     "0123456789",
					"full_name": "User D",
				},
			},
			Response: TCResponse{
				Status: http.StatusCreated,
			},
		},
	}
	// 1. Import users from CSV file
	store.InsertFromCSV(t, "users", users_as_csv)
	// 2. Try different login scenarios
	for name, scenario := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			response := suite.TestAPI.Post("/api/v1/user", scenario.Request.Body)
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
