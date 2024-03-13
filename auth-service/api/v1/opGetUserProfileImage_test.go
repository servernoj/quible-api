package v1_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/suite"
	"github.com/rs/zerolog/log"
)

func (tc *TestCases) TestGetUserProfileImage(t *testing.T) {
	// 1. Import users from CSV file
	db := tc.DBStore.RetrieveDB(t.Name())
	deps := tc.ServiceAPI.SetContext("opGetUserProfileImage")
	deps.Set("db", db)
	if err := suite.InsertFromCSV(db, "users", UsersCSV); err != nil {
		t.Fatalf("unable to import test data from CSV: %s", err)
	}
	// 2. Define test scenarios
	testCases := libAPI.TCScenarios{
		"SuccessWithImage": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Success with a valid image in the user’s profile",
				Request: libAPI.TCRequest{
					Params: map[string]any{
						"userId": "9bef41ed-fb10-4791-b02e-96b372c09466",
					},
				},
				Response: libAPI.TCResponse{
					Status: http.StatusOK,
				},
				ExtraTests: []libAPI.TCExtraTest{
					func(req libAPI.TCRequest, res *httptest.ResponseRecorder) bool {
						user, err := models.FindUser(context.Background(), db, req.Params["userId"].(string))
						if err != nil {
							log.Error().Err(err).Send()
							return false
						}
						var imageData v1.ImageData
						if err := json.Unmarshal(user.Image.Bytes, &imageData); err != nil {
							log.Error().Err(err).Send()
							return false
						}
						if imageData.ContentType != res.Result().Header.Get("content-type") {
							log.Error().Msg("unexpected content type")
							return false
						}
						return reflect.DeepEqual(
							imageData.BinaryContent,
							res.Body.Bytes(),
						)
					},
				},
			}
		},
		"FailureOnNoImage": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure when a user doesn't have an image",
				Request: libAPI.TCRequest{
					Params: map[string]any{
						// User B
						"userId": "42d29b4b-935d-4f35-b26c-70080107f6d6",
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusNotFound,
					ErrorCode: v1.Err404_UserHasNoImage.Ptr(),
				},
			}
		},
		"FailureOnNotFoundUser": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure due to a non-existing user in the request",
				Request: libAPI.TCRequest{
					Params: map[string]any{
						// non-existing userId of the correct UUID format
						"userId": "00000000-0000-0000-0000-000000000000",
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusNotFound,
					ErrorCode: v1.Err404_UserNotFound.Ptr(),
				},
			}
		},
		"FailureOnInvalidImageData": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure due to invalid image data stored in DB",
				Request: libAPI.TCRequest{
					Params: map[string]any{
						// user C with invalid image data in DB
						"userId": "c6174e8a-e12f-4d64-a4fe-a3b0c081bd31",
					},
				},
				Response: libAPI.TCResponse{
					Status:    http.StatusInternalServerError,
					ErrorCode: v1.Err500_UnableToRetrieveProfileImage.Ptr(),
				},
			}
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(
			name,
			scenario.GetRunner(
				tc.TestAPI,
				http.MethodGet,
				"/user/%s/image",
				"userId",
			),
		)
	}
}
