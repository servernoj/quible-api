package v1_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/quible-io/quible-api/lib/env"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/store"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
)

type MockedEmailSender struct {
	mock.Mock
}

func (m *MockedEmailSender) SendEmail(ctx context.Context, emailPayload email.EmailPayload) error {
	args := m.Called(ctx, emailPayload)
	log.Info().Msg("Email sender mocked")
	return args.Error(0)
}

func (suite *TestCases) TestCreateUser() {
	t := suite.T()
	// 1. Import users from CSV file
	store.InsertFromCSV(t, "users", UsersCSV)
	// 2. Load environment variables
	env.Setup()
	// 3. Define test scenarios
	testCases := TCScenarios{
		"FailureOnActivatedWithExistingEmail": TCData{
			Description: "Failure of registering existing email on activated user",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"username":  "new-username",
						"email":     "userA@gmail.com",
						"password":  "password",
						"phone":     "0123456789",
						"full_name": "existing email",
					},
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
				Args: []any{
					map[string]any{
						"username":  "userA",
						"email":     "non-existent@gmail.com",
						"password":  "password",
						"phone":     "0123456789",
						"full_name": "existing username",
					},
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
				Args: []any{
					map[string]any{
						"username":  "userD",
						"email":     "not-an-email-address",
						"password":  "password",
						"phone":     "0123456789",
						"full_name": "User D",
					},
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
				Args: []any{
					map[string]any{
						"username":  "userD",
						"email":     "userD@gmail.com",
						"password":  "password",
						"phone":     "invalid",
						"full_name": "User D",
					},
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_InvalidPhoneFormat),
			},
		},
		"FailureSendEmail": TCData{
			Description: "Failure due to error in email sender",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"username":  "userD",
						"email":     "userD@gmail.com",
						"password":  "password",
						"phone":     "0123456789",
						"full_name": "User D",
					},
				},
			},
			Response: TCResponse{
				Status:    http.StatusFailedDependency,
				ErrorCode: misc.Of(v1.Err424_UnableToSendEmail),
			},
			PreHook: func(t *testing.T) any {
				mockedEmailSender := new(MockedEmailSender)
				mockedEmailSender.On("SendEmail", mock.Anything, mock.Anything).Return(errors.New("email delivery failed"))
				suite.ServiceAPI.SetEmailSender(
					mockedEmailSender,
				)
				return mockedEmailSender
			},
			PostHook: func(t *testing.T, state any) {
				mockedEmailSender := state.(*MockedEmailSender)
				mockedEmailSender.AssertNumberOfCalls(t, "SendEmail", 1)
			},
		},
		"Success": TCData{
			Description: "Happy path with mocked email sender",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"username":  "userD",
						"email":     "userD@gmail.com",
						"password":  "password",
						"phone":     "0123456789",
						"full_name": "User D",
					},
				},
				Params: map[string]any{
					"username": "userD",
					"email":    "userD@gmail.com",
				},
			},
			Response: TCResponse{
				Status: http.StatusCreated,
			},
			PreHook: func(t *testing.T) any {
				mockedEmailSender := new(MockedEmailSender)
				mockedEmailSender.On("SendEmail", mock.Anything, mock.Anything).Return(nil)
				suite.ServiceAPI.SetEmailSender(
					mockedEmailSender,
				)
				return mockedEmailSender
			},
			PostHook: func(t *testing.T, state any) {
				mockedEmailSender := state.(*MockedEmailSender)
				mockedEmailSender.AssertNumberOfCalls(t, "SendEmail", 1)
			},
			ExtraTests: []TCExtraTest{
				func(req TCRequest, response *httptest.ResponseRecorder) bool {
					var responseBody v1.UserSimplified
					if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
						return false
					}
					userInDB, err := models.FindUserG(context.Background(), responseBody.ID)
					if err != nil {
						return false
					}
					email := req.Params["email"].(string)
					username := req.Params["username"].(string)
					if userInDB.Email != email || userInDB.Username != username || userInDB.ActivatedAt.Ptr() != nil {
						return false
					}
					return true
				},
			},
		},
	}
	// 4. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(suite.TestAPI, http.MethodPost, "/api/v1/user"))
	}
}
