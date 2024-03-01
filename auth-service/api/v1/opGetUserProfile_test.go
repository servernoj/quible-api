package v1_test

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/store"
	"github.com/stretchr/testify/assert"
)

func (suite *TestCases) TestGetUserProfile() {
	t := suite.T()
	// 1. Import users from CSV file
	store.InsertFromCSV(t, "users", UsersCSV)
	// 2. Define test scenarios
	testCases := TCScenarios{
		"SuccessWithImage": TCData{
			Description: "Success with image in profile",
			Request: TCRequest{
				Params: map[string]any{
					"userId": "9bef41ed-fb10-4791-b02e-96b372c09466",
				},
			},
			Response: TCResponse{
				Status: http.StatusOK,
			},
			ExtraTests: []TCExtraTest{
				func(_ TCRequest, res *httptest.ResponseRecorder) bool {
					var got v1.UserProfile
					if err := json.NewDecoder(res.Result().Body).Decode(&got); err != nil {
						return false
					}
					dataUrl := "data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz48IS0tIFVwbG9hZGVkIHRvOiBTVkcgUmVwbywgd3d3LnN2Z3JlcG8uY29tLCBHZW5lcmF0b3I6IFNWRyBSZXBvIE1peGVyIFRvb2xzIC0tPg0KPHN2ZyB3aWR0aD0iODAwcHgiIGhlaWdodD0iODAwcHgiIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4NCjxwYXRoIG9wYWNpdHk9IjAuNSIgZD0iTTEyIDIyQzcuMjg1OTUgMjIgNC45Mjg5MyAyMiAzLjQ2NDQ3IDIwLjUzNTVDMiAxOS4wNzExIDIgMTYuNzE0IDIgMTJDMiA3LjI4NTk1IDIgNC45Mjg5MyAzLjQ2NDQ3IDMuNDY0NDdDNC45Mjg5MyAyIDcuMjg1OTUgMiAxMiAyQzE2LjcxNCAyIDE5LjA3MTEgMiAyMC41MzU1IDMuNDY0NDdDMjIgNC45Mjg5MyAyMiA3LjI4NTk1IDIyIDEyQzIyIDE2LjcxNCAyMiAxOS4wNzExIDIwLjUzNTUgMjAuNTM1NUMxOS4wNzExIDIyIDE2LjcxNCAyMiAxMiAyMloiIGZpbGw9IiMxQzI3NEMiLz4NCjwvc3ZnPg=="
					want := v1.UserProfile{
						ID:       "9bef41ed-fb10-4791-b02e-96b372c09466",
						FullName: "User A",
						Image:    &dataUrl,
					}
					return reflect.DeepEqual(got, want)
					// if diff := cmp.Diff(want, got); diff != "" {
					// 	log.Warn().Msg(diff)
					// 	return false
					// }
					// return true
				},
			},
		},
		"SuccessWithoutImage": TCData{
			Description: "Success without image in profile",
			Request: TCRequest{
				Params: map[string]any{
					"userId": "42d29b4b-935d-4f35-b26c-70080107f6d6",
				},
			},
			Response: TCResponse{
				Status: http.StatusOK,
			},
			ExtraTests: []TCExtraTest{
				func(_ TCRequest, res *httptest.ResponseRecorder) bool {
					var got v1.UserProfile
					if err := json.NewDecoder(res.Result().Body).Decode(&got); err != nil {
						return false
					}
					want := v1.UserProfile{
						ID:       "42d29b4b-935d-4f35-b26c-70080107f6d6",
						FullName: "User B",
						Image:    nil,
					}
					return reflect.DeepEqual(got, want)
				},
			},
		},
		"FailureOnNotFoundUser": TCData{
			Description: "Failure due to non-existing user in request",
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
					// user with invalid image data in DB
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
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			path := fmt.Sprintf("/api/v1/user/%s/profile", scenario.Request.Params["userId"])
			response := suite.TestAPI.Get(path, scenario.Request.Headers...)
			// response status
			assert.EqualValues(scenario.Response.Status, response.Code, "response status should match the expectation")
			// error code in case of error
			if scenario.Response.ErrorCode != nil {
				assert.Contains(
					response.Body.String(),
					strconv.Itoa(int(*scenario.Response.ErrorCode)),
					"error code should match expectation",
				)
			}
			// extra tests (if present)
			for _, fn := range scenario.ExtraTests {
				assert.True(
					fn(scenario.Request, response),
				)
			}
		})
	}
}
