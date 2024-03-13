package v1_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/suite"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
)

type CreateUserEmailSender struct {
	mock.Mock
}

func (m *CreateUserEmailSender) SendEmail(ctx context.Context, emailPayload email.EmailPayload) error {
	args := m.Called(ctx, emailPayload)
	log.Info().Msg("Email sender mocked")
	return args.Error(0)
}

func (tc *TestCases) TestCreateUser(t *testing.T) {
	// 1. Import users from CSV file
	db := tc.DBStore.RetrieveDB(t.Name())
	deps := tc.ServiceAPI.SetContext("opCreateUser")
	deps.Set("db", db)
	if err := suite.InsertFromCSV(db, "users", UsersCSV); err != nil {
		t.Fatalf("unable to import test data from CSV: %s", err)
	}
	// 2. Define test scenarios
	testCases := libAPI.TCScenarios{
		"FailureOnActivatedWithExistingEmail": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure of registering existing email on activated user",
				Request: libAPI.TCRequest{
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
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_UserWithEmailOrUsernameExists.Ptr(),
				},
			}
		},
		"FailureOnActivatedWithExistingUsername": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure of registering existing username on activated user",
				Request: libAPI.TCRequest{
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
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_UserWithEmailOrUsernameExists.Ptr(),
				},
			}
		},
		"FailureOnInvalidEmail": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure of registering user with invalid email",
				Request: libAPI.TCRequest{
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
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_InvalidEmailFormat.Ptr(),
				},
			}
		},
		"FailureOnInvalidPhone": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure of registering user with invalid phone",
				Request: libAPI.TCRequest{
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
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_InvalidPhoneFormat.Ptr(),
				},
			}
		},
		"FailureSendEmail": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure due to error in email sender",
				Request: libAPI.TCRequest{
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
				Response: libAPI.TCResponse{
					Status:    http.StatusFailedDependency,
					ErrorCode: v1.Err424_UnableToSendEmail.Ptr(),
				},
				PreHook: func(t *testing.T) any {
					mockedEmailSender := new(CreateUserEmailSender)
					mockedEmailSender.On("SendEmail", mock.Anything, mock.Anything).Return(errors.New("email delivery failed"))
					deps.Set("mailer", mockedEmailSender)
					return mockedEmailSender
				},
				PostHook: func(t *testing.T, state any) {
					mockedEmailSender := state.(*CreateUserEmailSender)
					mockedEmailSender.AssertNumberOfCalls(t, "SendEmail", 1)
				},
			}
		},
		"SuccessDev": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Happy path with mocked email sender and IS_DEV enabled (auto-activation)",
				Envs: libAPI.TCEnv{
					"IS_DEV": "1",
				},
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"username":  "userE",
							"email":     "userE@gmail.com",
							"password":  "password",
							"phone":     "0123456789",
							"full_name": "User E",
						},
					},
					Params: map[string]any{
						"username":       "userE",
						"email":          "userE@gmail.com",
						"autoActivation": "true",
					},
				},
				Response: libAPI.TCResponse{
					Status: http.StatusCreated,
				},
				PreHook: func(t *testing.T) any {
					mockedEmailSender := new(CreateUserEmailSender)
					deps.Set("mailer", mockedEmailSender)
					return mockedEmailSender
				},
				PostHook: func(t *testing.T, state any) {
					mockedEmailSender := state.(*CreateUserEmailSender)
					mockedEmailSender.AssertNumberOfCalls(t, "SendEmail", 0)
				},
				ExtraTests: []libAPI.TCExtraTest{
					func(req libAPI.TCRequest, response *httptest.ResponseRecorder) bool {
						var responseBody v1.UserSimplified
						if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
							return false
						}
						userInDB, err := models.FindUser(context.Background(), db, responseBody.ID)
						if err != nil {
							return false
						}
						email := req.Params["email"].(string)
						username := req.Params["username"].(string)
						if userInDB.Email != email || userInDB.Username != username || userInDB.ActivatedAt.Ptr() == nil {
							return false
						}
						return true
					},
				},
			}
		},
		"SuccessProd": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Happy path with mocked email sender",
				Request: libAPI.TCRequest{
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
				Response: libAPI.TCResponse{
					Status: http.StatusCreated,
				},
				PreHook: func(t *testing.T) any {
					mockedEmailSender := new(CreateUserEmailSender)
					mockedEmailSender.On("SendEmail", mock.Anything, mock.Anything).Return(nil)
					deps.Set("mailer", mockedEmailSender)
					return mockedEmailSender
				},
				PostHook: func(t *testing.T, state any) {
					mockedEmailSender := state.(*CreateUserEmailSender)
					mockedEmailSender.AssertNumberOfCalls(t, "SendEmail", 1)
				},
				ExtraTests: []libAPI.TCExtraTest{
					func(req libAPI.TCRequest, response *httptest.ResponseRecorder) bool {
						var responseBody v1.UserSimplified
						if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
							return false
						}
						userInDB, err := models.FindUser(context.Background(), db, responseBody.ID)
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
			}
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, func(t *testing.T) {
			tcData := scenario(t)
			url := "/user"
			if tcData.Request.Params["autoActivation"] != nil {
				url = fmt.Sprintf(
					"/user?auto-activation=%s",
					tcData.Request.Params["autoActivation"].(string),
				)
			}
			runner := scenario.GetRunner(tc.TestAPI, http.MethodPost, url)
			runner(t)
		})
	}
}
