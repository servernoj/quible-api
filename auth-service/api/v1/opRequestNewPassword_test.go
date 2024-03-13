package v1_test

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/quible-io/quible-api/lib/suite"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
)

type RequestNewPasswordEmailSender struct {
	mock.Mock
}

func (m *RequestNewPasswordEmailSender) SendEmail(ctx context.Context, emailPayload email.EmailPayload) error {
	args := m.Called(ctx, emailPayload)
	log.Info().Msg("Email sender mocked")
	return args.Error(0)
}

func (tc *TestCases) TestRequestNewPassword(t *testing.T) {
	// 1. Import users from CSV file
	db := tc.DBStore.RetrieveDB(t.Name())
	deps := tc.ServiceAPI.SetContext("opRequestNewPassword")
	deps.Set("db", db)
	if err := suite.InsertFromCSV(db, "users", UsersCSV); err != nil {
		t.Fatalf("unable to import test data from CSV: %s", err)
	}
	// 2. Define test scenarios
	testCases := libAPI.TCScenarios{
		"FailureOnInvalidEmail": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure to reset the password for a user with an invalid email",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"email": "not-an-email-address",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_InvalidEmailFormat.Ptr(),
				},
			}
		},
		"NoEmailForNonExistingUser": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure to reset the password for a non-existing user",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"email": "userD@gmail.com",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status: http.StatusAccepted,
				},
				PreHook: func(t *testing.T) any {
					mockedEmailSender := new(RequestNewPasswordEmailSender)
					deps.Set("mailer", mockedEmailSender)
					return mockedEmailSender
				},
				PostHook: func(t *testing.T, state any) {
					mockedEmailSender := state.(*RequestNewPasswordEmailSender)
					mockedEmailSender.AssertNumberOfCalls(t, "SendEmail", 0)
				},
			}
		},
		"FailureOnSendEmail": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure due to an error caused by email sender",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"email": "userA@gmail.com",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusFailedDependency,
					ErrorCode: v1.Err424_UnableToSendEmail.Ptr(),
				},
				PreHook: func(t *testing.T) any {
					mockedEmailSender := new(RequestNewPasswordEmailSender)
					mockedEmailSender.On("SendEmail", mock.Anything, mock.Anything).Return(errors.New("email delivery failed"))
					deps.Set("mailer", mockedEmailSender)
					return mockedEmailSender
				},
				PostHook: func(t *testing.T, state any) {
					mockedEmailSender := state.(*RequestNewPasswordEmailSender)
					mockedEmailSender.AssertNumberOfCalls(t, "SendEmail", 1)
				},
			}
		},
		"Success": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Happy path with mocked email sender",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"email": "userA@gmail.com",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status: http.StatusAccepted,
				},
				PreHook: func(t *testing.T) any {
					mockedEmailSender := new(RequestNewPasswordEmailSender)
					mockedEmailSender.On(
						"SendEmail",
						mock.Anything,
						mock.MatchedBy(
							func(payload email.EmailPayload) bool {
								got := email.EmailPayload{
									From:    payload.From,
									To:      payload.To,
									Subject: payload.Subject,
								}
								wanted := email.EmailPayload{
									From:    "no-reply@quible.io",
									To:      "userA@gmail.com",
									Subject: "Password reset",
								}
								return reflect.DeepEqual(wanted, got)
							},
						),
					).Return(nil)
					deps.Set("mailer", mockedEmailSender)
					return mockedEmailSender
				},
				PostHook: func(t *testing.T, state any) {
					mockedEmailSender := state.(*RequestNewPasswordEmailSender)
					mockedEmailSender.AssertNumberOfCalls(t, "SendEmail", 1)
				},
			}
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(tc.TestAPI, http.MethodPost, "/user/request-new-password"))
	}
}
