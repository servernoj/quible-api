package v1_test

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/suite"
)

func (tc *TestCases) TestGetUserProfile(t *testing.T) {
	// 1. Import users from CSV file
	db := tc.DBStore.RetrieveDB(t.Name())
	deps := tc.ServiceAPI.SetContext("opGetUserProfile")
	deps.Set("db", db)
	if err := suite.InsertFromCSV(db, "users", UsersCSV); err != nil {
		t.Fatalf("unable to import test data from CSV: %s", err)
	}
	// 2. Define test scenarios
	testCases := libAPI.TCScenarios{
		"SuccessWithImage": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Success with image in profile",
				Request: libAPI.TCRequest{
					Params: map[string]any{
						"userId": "9bef41ed-fb10-4791-b02e-96b372c09466",
					},
				},
				Response: libAPI.TCResponse{
					Status: http.StatusOK,
				},
				ExtraTests: []libAPI.TCExtraTest{
					func(_ libAPI.TCRequest, res *httptest.ResponseRecorder) bool {
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
					},
				},
			}
		},
		"SuccessWithoutImage": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Success without image in profile",
				Request: libAPI.TCRequest{
					Params: map[string]any{
						"userId": "42d29b4b-935d-4f35-b26c-70080107f6d6",
					},
				},
				Response: libAPI.TCResponse{
					Status: http.StatusOK,
				},
				ExtraTests: []libAPI.TCExtraTest{
					func(_ libAPI.TCRequest, res *httptest.ResponseRecorder) bool {
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
			}
		},
		"FailureOnNotFoundUser": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure due to non-existing user in request",
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
						// user with invalid image data in DB
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
				"/user/%s/profile",
				"userId",
			),
		)
	}
}
