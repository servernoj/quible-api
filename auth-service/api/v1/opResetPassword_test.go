package v1_test

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	"github.com/quible-io/quible-api/auth-service/services/userService"
	"github.com/quible-io/quible-api/lib/env"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/store"
	"github.com/rs/zerolog/log"
)

func (suite *TestCases) TestPasswordReset() {
	t := suite.T()
	// 1. Import users from CSV file
	store.InsertFromCSV(t, "users", UsersCSV)
	// 2. Load environment variables
	env.Setup()
	// 3. Define test scenarios
	testCases := TCScenarios{
		"FailureOnInvalidStepValue": TCData{
			Description: "Failure when the `step` field is set to anything but `define` or `validate`",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"token": "passes.regex.validation",
						"step":  "invalid",
					},
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_InvalidRequest),
			},
		},
		"FailureOnInvalidTokenAction": TCData{
			Description: "Failure token has a wrong action (anything but `password reset`)",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"token": GetToken(t, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess),
						"step":  "validate",
					},
				},
			},
			Response: TCResponse{
				Status:    http.StatusUnauthorized,
				ErrorCode: misc.Of(v1.Err401_InvalidPasswordResetToken),
			},
		},
		"FailureOnNonExistingUser": TCData{
			Description: "Failure when a token of the correct kind references a non-existing user",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						// token with expiration at 2030-01-01 of PasswordReset action but for non-existing user
						"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4OTM0NTYwMDAsInVzZXJJZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImFjdGlvbiI6IlBhc3N3b3JkUmVzZXQiLCJleHRyYUNsYWltcyI6bnVsbH0.5n3jVZDQ1wAHNv5s-uKfyuQM7YZvpOz3aa--ACy76oc",
						"step":  "validate",
					},
				},
			},
			Response: TCResponse{
				Status:    http.StatusExpectationFailed,
				ErrorCode: misc.Of(v1.Err417_UnableToAssociateUser),
			},
		},
		"FailureOnPasswordMismatch": TCData{
			Description: "Failure when password doesn't match its confirmation in the `define` step",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"token":           GetToken(t, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionPasswordReset),
						"step":            "define",
						"password":        "abc123",
						"confirmPassword": "123abc",
					},
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_UnsatisfactoryConfirmPassword),
			},
		},
		"FailureOnWeakPassword": TCData{
			Description: "Failure when the password is too weak (less than 6 characters) in the `define` step",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"token":           GetToken(t, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionPasswordReset),
						"step":            "define",
						"password":        "abc",
						"confirmPassword": "abc",
					},
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_UnsatisfactoryPassword),
			},
		},
		"SuccessOnDefineStep": TCData{
			Description: "Success while setting a new password in the `define` step",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"token":           GetToken(t, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionPasswordReset),
						"password":        "abc123",
						"confirmPassword": "abc123",
						"step":            "define",
					},
				},
			},
			Response: TCResponse{
				Status: http.StatusOK,
			},
			ExtraTests: []TCExtraTest{
				func(req TCRequest, res *httptest.ResponseRecorder) bool {
					user, err := models.FindUserG(context.Background(), "9bef41ed-fb10-4791-b02e-96b372c09466")
					if err != nil {
						log.Error().Err(err).Send()
						return false
					}
					us := new(userService.UserService)
					if err := us.ValidatePassword(user.HashedPassword, "abc123"); err != nil {
						log.Error().Err(err).Send()
						return false
					}
					return true
				},
			},
		},
		"SuccessOnValidationStep": TCData{
			Description: "Success with a correct token in the `validate` step",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"token": GetToken(t, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionPasswordReset),
						"step":  "validate",
					},
				},
			},
			Response: TCResponse{
				Status: http.StatusOK,
			},
		},
	}
	// 4. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(suite.TestAPI, http.MethodPost, "/api/v1/user/password-reset"))
	}
}
