package v1_test

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	"github.com/quible-io/quible-api/lib/env"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/store"
	"github.com/rs/zerolog/log"
)

func (suite *TestCases) TestUpdateUser() {
	t := suite.T()
	// 1. Import users from CSV file
	store.InsertFromCSV(t, "users", UsersCSV)
	// 2. Load environment variables
	env.Setup()
	// 3. Define test scenarios
	testCases := TCScenarios{
		"SuccessOnEmptyRequest": TCData{
			Description: "Success on an empty request body, user record stays intact",
			Request: func(t *testing.T) TCRequest {
				userId := "9bef41ed-fb10-4791-b02e-96b372c09466"
				user, err := models.FindUserG(
					context.Background(),
					userId,
					models.UserColumns.ID,
					models.UserColumns.Email,
					models.UserColumns.Username,
					models.UserColumns.FullName,
					models.UserColumns.Phone,
				)
				if err != nil {
					t.Fatal(err)
				}
				return TCRequest{
					Args: []any{
						map[string]any{},
						fmt.Sprintf(
							"Authorization: Bearer %s",
							GetToken(t, userId, jwt.TokenActionAccess),
						),
					},
					Params: map[string]any{
						"user": user,
					},
				}
			}(t),
			Response: TCResponse{
				Status: http.StatusOK,
			},
			ExtraTests: []TCExtraTest{
				func(req TCRequest, res *httptest.ResponseRecorder) bool {
					wanted := req.Params["user"].(*models.User)
					got, err := models.FindUserG(
						context.Background(),
						wanted.ID,
						models.UserColumns.ID,
						models.UserColumns.Email,
						models.UserColumns.Username,
						models.UserColumns.FullName,
						models.UserColumns.Phone,
					)
					if err != nil {
						log.Error().Err(err).Send()
						return false
					}
					if diff := cmp.Diff(wanted, got); diff != "" {
						log.Warn().Msg(diff)
						return false
					}
					return true
				},
			},
		},
		"SuccessOnCompleteRequest": TCData{
			Description: "Success on valid request to change all fields",
			Request: func(t *testing.T) TCRequest {
				userId := "9bef41ed-fb10-4791-b02e-96b372c09466"
				return TCRequest{
					Args: []any{
						map[string]any{
							"email":     "userD@gmail.com",
							"username":  "userD",
							"phone":     "1111111111",
							"full_name": "User D",
						},
						fmt.Sprintf(
							"Authorization: Bearer %s",
							GetToken(t, userId, jwt.TokenActionAccess),
						),
					},
					Params: map[string]any{
						"userId": userId,
					},
				}
			}(t),
			Response: TCResponse{
				Status: http.StatusOK,
			},
			ExtraTests: []TCExtraTest{
				func(req TCRequest, res *httptest.ResponseRecorder) bool {
					requestData := req.Args[0].(map[string]any)
					wanted := v1.UserSimplified{
						ID:       req.Params["userId"].(string),
						Email:    requestData["email"].(string),
						Username: requestData["username"].(string),
						Phone:    requestData["phone"].(string),
						FullName: requestData["full_name"].(string),
					}
					foundUser, err := models.FindUserG(
						context.Background(),
						wanted.ID,
					)
					got := v1.UserSimplified{
						ID:       foundUser.ID,
						Email:    foundUser.Email,
						Username: foundUser.Username,
						Phone:    foundUser.Phone,
						FullName: foundUser.FullName,
					}
					if err != nil {
						log.Error().Err(err).Send()
						return false
					}
					if diff := cmp.Diff(wanted, got); diff != "" {
						log.Warn().Msg(diff)
						return false
					}
					return true
				},
			},
		},
		"FailureOnEmailFormat": TCData{
			Description: "Failure on invalid email format",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"email": "invalid-email",
					},
					fmt.Sprintf(
						"Authorization: Bearer %s",
						GetToken(t, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess),
					),
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_InvalidEmailFormat),
			},
		},
		"FailureOnPhoneFormat": TCData{
			Description: "Failure on invalid phone format",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"phone": "!_invalid_phone_number_!",
					},
					fmt.Sprintf(
						"Authorization: Bearer %s",
						GetToken(t, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess),
					),
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_InvalidPhoneFormat),
			},
		},
		"FailureOnTooShortFullName": TCData{
			Description: "Failure on too short full name",
			Request: TCRequest{
				Args: []any{
					map[string]any{
						"full_name": "",
					},
					fmt.Sprintf(
						"Authorization: Bearer %s",
						GetToken(t, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess),
					),
				},
			},
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_InvalidRequest),
			},
		},
	}
	// 4. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(suite.TestAPI, http.MethodPatch, "/api/v1/user"))
	}
}
