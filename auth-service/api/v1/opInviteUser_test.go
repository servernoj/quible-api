package v1_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	"github.com/quible-io/quible-api/lib/email"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/store"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
)

type InviteUserEmailSender struct {
	mock.Mock
}

func (m *InviteUserEmailSender) SendEmail(ctx context.Context, emailPayload email.EmailPayload) error {
	args := m.Called(ctx, emailPayload)
	log.Info().Msg("Email sender mocked")
	return args.Error(0)
}

func (suite *TestCases) TestInviteUser() {
	t := suite.T()
	// 1. Import users from CSV file
	store.InsertFromCSV(t, "users", UsersCSV)
	// 2. Define test scenarios
	testCases := TCScenarios{
		"FailureOnInvalidEmail": TCData{
			Description: "Failure to invite a user with an invalid email",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"email":     "not-an-email-address",
						"full_name": "User D",
					},
					fmt.Sprintf("Authorization: Bearer %s", GetToken(t, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess)),
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_InvalidEmailFormat),
			},
		},
		"FailureOnInvalidName": TCData{
			Description: "Failure to invite a user with an empty (invalid) name",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"email":     "userD@gmail.com",
						"full_name": "",
					},
					fmt.Sprintf("Authorization: Bearer %s", GetToken(t, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess)),
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_InvalidRequest),
			},
		},
		"FailureToInviteExistingUser": TCData{
			Description: "Failure to invite an existing user",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"email":     "userB@gmail.com",
						"full_name": "user B",
					},
					fmt.Sprintf("Authorization: Bearer %s", GetToken(t, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess)),
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_UserWithEmailExists),
			},
		},
		"FailureSendEmail": TCData{
			Description: "Failure due to an error caused by email sender",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"email":     "userD@gmail.com",
						"full_name": "User D",
					},
					fmt.Sprintf("Authorization: Bearer %s", GetToken(t, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess)),
				},
			},
			Response: TCResponse{
				Status:    http.StatusFailedDependency,
				ErrorCode: misc.Of(v1.Err424_UnableToSendEmail),
			},
			PreHook: func(t *testing.T) any {
				mockedEmailSender := new(InviteUserEmailSender)
				mockedEmailSender.On("SendEmail", mock.Anything, mock.Anything).Return(errors.New("email delivery failed"))
				suite.ServiceAPI.SetEmailSender(
					mockedEmailSender,
				)
				return mockedEmailSender
			},
			PostHook: func(t *testing.T, state any) {
				mockedEmailSender := state.(*InviteUserEmailSender)
				mockedEmailSender.AssertNumberOfCalls(t, "SendEmail", 1)
			},
		},
		"Success": TCData{
			Description: "Happy path with mocked email sender",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"email":     "userD@gmail.com",
						"full_name": "User D",
					},
					fmt.Sprintf("Authorization: Bearer %s", GetToken(t, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess)),
				},
			},
			Response: TCResponse{
				Status: http.StatusOK,
			},
			PreHook: func(t *testing.T) any {
				mockedEmailSender := new(InviteUserEmailSender)
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
								To:      "userD@gmail.com",
								Subject: "Invitation to register Quible account",
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
				mockedEmailSender := state.(*InviteUserEmailSender)
				mockedEmailSender.AssertNumberOfCalls(t, "SendEmail", 1)
			},
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(suite.TestAPI, http.MethodPost, "/api/v1/user/invite"))
	}
}
