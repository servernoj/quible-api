package v1_test

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/store"
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

func (suite *TestCases) TestRequestNewPassword() {
	t := suite.T()
	// 1. Import users from CSV file
	store.InsertFromCSV(t, "users", UsersCSV)
	// 2. Define test scenarios
	testCases := TCScenarios{
		"FailureOnInvalidEmail": TCData{
			Description: "Failure to reset the password for a user with an invalid email",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"email": "not-an-email-address",
					},
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_InvalidEmailFormat),
			},
		},
		"NoEmailForNonExistingUser": TCData{
			Description: "Failure to reset the password for a non-existing user",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"email": "userD@gmail.com",
					},
				},
			},
			Response: TCResponse{
				Status: http.StatusAccepted,
			},
			PreHook: func(t *testing.T) any {
				mockedEmailSender := new(RequestNewPasswordEmailSender)
				suite.ServiceAPI.SetEmailSender(
					mockedEmailSender,
				)
				return mockedEmailSender
			},
			PostHook: func(t *testing.T, state any) {
				mockedEmailSender := state.(*RequestNewPasswordEmailSender)
				mockedEmailSender.AssertNumberOfCalls(t, "SendEmail", 0)
			},
		},
		"FailureOnSendEmail": TCData{
			Description: "Failure due to an error caused by email sender",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"email": "userA@gmail.com",
					},
				},
			},
			Response: TCResponse{
				Status:    http.StatusFailedDependency,
				ErrorCode: misc.Of(v1.Err424_UnableToSendEmail),
			},
			PreHook: func(t *testing.T) any {
				mockedEmailSender := new(RequestNewPasswordEmailSender)
				mockedEmailSender.On("SendEmail", mock.Anything, mock.Anything).Return(errors.New("email delivery failed"))
				suite.ServiceAPI.SetEmailSender(
					mockedEmailSender,
				)
				return mockedEmailSender
			},
			PostHook: func(t *testing.T, state any) {
				mockedEmailSender := state.(*RequestNewPasswordEmailSender)
				mockedEmailSender.AssertNumberOfCalls(t, "SendEmail", 1)
			},
		},
		"Success": TCData{
			Description: "Happy path with mocked email sender",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"email": "userA@gmail.com",
					},
				},
			},
			Response: TCResponse{
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
				suite.ServiceAPI.SetEmailSender(
					mockedEmailSender,
				)
				return mockedEmailSender
			},
			PostHook: func(t *testing.T, state any) {
				mockedEmailSender := state.(*RequestNewPasswordEmailSender)
				mockedEmailSender.AssertNumberOfCalls(t, "SendEmail", 1)
			},
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(suite.TestAPI, http.MethodPost, "/api/v1/user/request-new-password"))
	}
}
