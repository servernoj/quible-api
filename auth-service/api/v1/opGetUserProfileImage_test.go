package v1_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/store"
	"github.com/rs/zerolog/log"
)

func (suite *TestCases) TestGetUserProfileImage() {
	t := suite.T()
	// 1. Import users from CSV file
	store.InsertFromCSV(t, "users", UsersCSV)
	// 2. Define test scenarios
	testCases := TCScenarios{
		"SuccessWithImage": TCData{
			Description: "Success with a valid image in the user’s profile",
			Request: TCRequest{
				Params: map[string]any{
					"userId": "9bef41ed-fb10-4791-b02e-96b372c09466",
				},
			},
			Response: TCResponse{
				Status: http.StatusOK,
			},
			ExtraTests: []TCExtraTest{
				func(req TCRequest, res *httptest.ResponseRecorder) bool {
					user, err := models.FindUserG(context.Background(), req.Params["userId"].(string))
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
		},
		"FailureOnNoImage": TCData{
			Description: "Failure when a user doesn't have an image",
			Request: TCRequest{
				Params: map[string]any{
					// User B
					"userId": "42d29b4b-935d-4f35-b26c-70080107f6d6",
				},
			},
			Response: TCResponse{
				Status:    http.StatusNotFound,
				ErrorCode: misc.Of(v1.Err404_UserHasNoImage),
			},
		},
		"FailureOnNotFoundUser": TCData{
			Description: "Failure due to a non-existing user in the request",
			Request: TCRequest{
				Params: map[string]any{
					// non-existing userId of the correct UUID format
					"userId": "00000000-0000-0000-0000-000000000000",
				},
			},
			Response: TCResponse{
				Status:    http.StatusNotFound,
				ErrorCode: misc.Of(v1.Err404_UserNotFound),
			},
		},
		"FailureOnInvalidImageData": TCData{
			Description: "Failure due to invalid image data stored in DB",
			Request: TCRequest{
				Params: map[string]any{
					// user C with invalid image data in DB
					"userId": "c6174e8a-e12f-4d64-a4fe-a3b0c081bd31",
				},
			},
			Response: TCResponse{
				Status:    http.StatusInternalServerError,
				ErrorCode: misc.Of(v1.Err500_UnableToRetrieveProfileImage),
			},
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(
			name,
			scenario.GetRunner(
				suite.TestAPI,
				http.MethodGet,
				"/api/v1/user/%s/image",
				scenario.Request.Params["userId"],
			),
		)
	}
}
