package v1_test

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	"github.com/quible-io/quible-api/auth-service/services/userService"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/suite"
	"github.com/rs/zerolog/log"
)

func (tc *TestCases) TestPasswordReset(t *testing.T) {
	// 1. Import users from CSV file
	db := tc.DBStore.RetrieveDB(t.Name())
	deps := tc.ServiceAPI.SetContext("opResetPassword")
	deps.Set("db", db)
	if err := suite.InsertFromCSV(db, "users", UsersCSV); err != nil {
		t.Fatalf("unable to import test data from CSV: %s", err)
	}
	// 2. Define test scenarios
	testCases := libAPI.TCScenarios{
		"FailureOnInvalidStepValue": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure when the `step` field is set to anything but `define` or `validate`",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"token": "passes.regex.validation",
							"step":  "invalid",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_InvalidRequest.Ptr(),
				},
			}
		},
		"FailureOnInvalidTokenAction": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure token has a wrong action (anything but `password reset`)",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"token": suite.GetToken(t, db, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess),
							"step":  "validate",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusUnauthorized,
					ErrorCode: v1.Err401_InvalidPasswordResetToken.Ptr(),
				},
			}
		},
		"FailureOnNonExistingUser": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure when a token of the correct kind references a non-existing user",
				Envs: libAPI.TCEnv{
					"ENV_JWT_SECRET": "secret",
				},
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							// token for non-existing user "00000000-0000-0000-0000-000000000000" with secret `secret`
							"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4OTM0NTYwMDAsInVzZXJJZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImFjdGlvbiI6IlBhc3N3b3JkUmVzZXQiLCJleHRyYUNsYWltcyI6bnVsbH0.nUt1bT0IJ-D3gYLgSOEE-6f9uiphf61TEVpjeWHikY4",
							"step":  "validate",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusExpectationFailed,
					ErrorCode: v1.Err417_UnableToAssociateUser.Ptr(),
				},
			}
		},
		"FailureOnPasswordMismatch": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure when password doesn't match its confirmation in the `define` step",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"token":           suite.GetToken(t, db, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionPasswordReset),
							"step":            "define",
							"password":        "abc123",
							"confirmPassword": "123abc",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_UnsatisfactoryConfirmPassword.Ptr(),
				},
			}
		},
		"FailureOnMissingPassword": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure when password is not sent in the `define` step",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"token":           suite.GetToken(t, db, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionPasswordReset),
							"step":            "define",
							"confirmPassword": "123abc",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_UnsatisfactoryPassword.Ptr(),
				},
			}
		},
		"FailureOnMissingPasswordConfirmation": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure when password confirmation is not sent in the `define` step",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"token":    suite.GetToken(t, db, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionPasswordReset),
							"step":     "define",
							"password": "123abc",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_UnsatisfactoryConfirmPassword.Ptr(),
				},
			}
		},
		"FailureOnWeakPassword": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure when the password is too weak (less than 6 characters) in the `define` step",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"token":           suite.GetToken(t, db, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionPasswordReset),
							"step":            "define",
							"password":        "abc",
							"confirmPassword": "abc",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_UnsatisfactoryPassword.Ptr(),
				},
			}
		},
		"SuccessOnDefineStep": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Success while setting a new password in the `define` step",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"token":           suite.GetToken(t, db, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionPasswordReset),
							"password":        "abc123",
							"confirmPassword": "abc123",
							"step":            "define",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status: http.StatusOK,
				},
				ExtraTests: []libAPI.TCExtraTest{
					func(req libAPI.TCRequest, res *httptest.ResponseRecorder) bool {
						user, err := models.FindUser(context.Background(), db, "9bef41ed-fb10-4791-b02e-96b372c09466")
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
			}
		},
		"SuccessOnValidationStep": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Success with a correct token in the `validate` step",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"token": suite.GetToken(t, db, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionPasswordReset),
							"step":  "validate",
						},
					},
				},
				Response: libAPI.TCResponse{
					Status: http.StatusOK,
				},
			}
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(tc.TestAPI, http.MethodPost, "/user/password-reset"))
	}
}
