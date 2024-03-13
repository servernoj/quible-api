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
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/suite"
	"github.com/rs/zerolog/log"
)

func (tc *TestCases) TestUpdateUser(t *testing.T) {
	// 1. Import users from CSV file
	db := tc.DBStore.RetrieveDB(t.Name())
	deps := tc.ServiceAPI.SetContext("opUpdateUser")
	deps.Set("db", db)
	if err := suite.InsertFromCSV(db, "users", UsersCSV); err != nil {
		t.Fatalf("unable to import test data from CSV: %s", err)
	}
	// 2. Define test scenarios
	testCases := libAPI.TCScenarios{
		"SuccessOnCompleteRequest": func(t *testing.T) libAPI.TCData {
			userId := "9bef41ed-fb10-4791-b02e-96b372c09466"
			return libAPI.TCData{
				Description: "Success on valid request to change all fields",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"email":     "userD@gmail.com",
							"username":  "userD",
							"phone":     "1111111111",
							"full_name": "User D",
						},
						fmt.Sprintf(
							"Authorization: Bearer %s",
							suite.GetToken(t, db, userId, jwt.TokenActionAccess),
						),
					},
					Params: map[string]any{
						"userId": userId,
					},
				},
				Response: libAPI.TCResponse{
					Status: http.StatusOK,
				},
				ExtraTests: []libAPI.TCExtraTest{
					func(req libAPI.TCRequest, res *httptest.ResponseRecorder) bool {
						requestData := req.Args[0].(map[string]any)
						wanted := v1.UserSimplified{
							ID:       req.Params["userId"].(string),
							Email:    requestData["email"].(string),
							Username: requestData["username"].(string),
							Phone:    requestData["phone"].(string),
							FullName: requestData["full_name"].(string),
						}
						foundUser, err := models.FindUser(
							context.Background(),
							db,
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
			}
		},
		"SuccessOnEmptyRequest": func(t *testing.T) libAPI.TCData {
			userId := "9bef41ed-fb10-4791-b02e-96b372c09466"
			user, err := models.FindUser(
				context.Background(),
				db,
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
			return libAPI.TCData{
				Description: "Success on an empty request body, user record stays intact",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{},
						fmt.Sprintf(
							"Authorization: Bearer %s",
							suite.GetToken(t, db, userId, jwt.TokenActionAccess),
						),
					},
					Params: map[string]any{
						"user": user,
					},
				},
				Response: libAPI.TCResponse{
					Status: http.StatusOK,
				},
				ExtraTests: []libAPI.TCExtraTest{
					func(req libAPI.TCRequest, res *httptest.ResponseRecorder) bool {
						wanted := req.Params["user"].(*models.User)
						got, err := models.FindUser(
							context.Background(),
							db,
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
			}
		},
		"FailureOnEmailFormat": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure on invalid email format",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"email": "invalid-email",
						},
						fmt.Sprintf(
							"Authorization: Bearer %s",
							suite.GetToken(t, db, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess),
						),
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_InvalidEmailFormat.Ptr(),
				},
			}
		},
		"FailureOnPhoneFormat": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure on invalid phone format",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"phone": "!_invalid_phone_number_!",
						},
						fmt.Sprintf(
							"Authorization: Bearer %s",
							suite.GetToken(t, db, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess),
						),
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_InvalidPhoneFormat.Ptr(),
				},
			}
		},
		"FailureOnTooShortFullName": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure on too short full name",
				Request: libAPI.TCRequest{
					Args: []any{
						map[string]any{
							"full_name": "",
						},
						fmt.Sprintf(
							"Authorization: Bearer %s",
							suite.GetToken(t, db, "9bef41ed-fb10-4791-b02e-96b372c09466", jwt.TokenActionAccess),
						),
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_InvalidRequest.Ptr(),
				},
			}
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(tc.TestAPI, http.MethodPatch, "/user"))
	}
}
