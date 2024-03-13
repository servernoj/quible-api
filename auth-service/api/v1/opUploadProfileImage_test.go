package v1_test

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"reflect"
	"testing"

	v1 "github.com/quible-io/quible-api/auth-service/api/v1"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/suite"
	"github.com/rs/zerolog/log"
)

func (tc *TestCases) TestUploadProfileImage(t *testing.T) {
	// 1. Import users from CSV file
	db := tc.DBStore.RetrieveDB(t.Name())
	deps := tc.ServiceAPI.SetContext("opUploadProfileImage")
	deps.Set("db", db)
	if err := suite.InsertFromCSV(db, "users", UsersCSV); err != nil {
		t.Fatalf("unable to import test data from CSV: %s", err)
	}
	// 2. Define test scenarios
	NewRequest := func(t *testing.T, contentType, imageFilename, userId, fieldName string) libAPI.TCRequest {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		h := make(textproto.MIMEHeader)
		h.Set(
			"Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="image.svg"`, fieldName),
		)
		h.Set(
			"Content-Type",
			contentType,
		)
		part, err := writer.CreatePart(h)
		if err != nil {
			t.Fatal(err)
		}
		file, err := os.Open(imageFilename)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		if _, err := io.Copy(part, file); err != nil {
			t.Fatal(err)
		}
		writer.Close()
		args := []any{
			fmt.Sprintf("Content-Length: %d", body.Len()),
			fmt.Sprintf("Content-Type: multipart/form-data; boundary=%s", writer.Boundary()),
			fmt.Sprintf("Authorization: Bearer %s", suite.GetToken(t, db, userId, jwt.TokenActionAccess)),
			bytes.NewReader(body.Bytes()),
		}
		return libAPI.TCRequest{
			Args: args,
			Params: map[string]any{
				"userId":        userId,
				"contentType":   contentType,
				"imageFilename": imageFilename,
			},
		}
	}
	testCases := libAPI.TCScenarios{
		"Success": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Successful upload and confirmation in DB",
				Request:     NewRequest(t, "image/svg+xml", "TestData/image.svg", "42d29b4b-935d-4f35-b26c-70080107f6d6", "image"),
				Response: libAPI.TCResponse{
					Status: http.StatusAccepted,
				},
				ExtraTests: []libAPI.TCExtraTest{
					func(t libAPI.TCRequest, res *httptest.ResponseRecorder) bool {
						user, err := models.FindUser(context.Background(), db, t.Params["userId"].(string))
						if err != nil {
							log.Error().Err(err).Send()
							return false
						}
						imageDataBytesPtr := user.Image.Ptr()
						if imageDataBytesPtr == nil {
							return false
						}
						var imageData v1.ImageData
						if err := json.Unmarshal(*imageDataBytesPtr, &imageData); err != nil {
							log.Error().Err(err).Send()
							return false
						}
						if imageData.ContentType != t.Params["contentType"].(string) {
							return false
						}
						f, err := os.Open(t.Params["imageFilename"].(string))
						if err != nil {
							log.Error().Err(err).Send()
							return false
						}
						defer f.Close()
						b, err := io.ReadAll(f)
						if err != nil {
							log.Error().Err(err).Send()
							return false
						}
						return reflect.DeepEqual(b, imageData.BinaryContent)
					},
				},
			}
		},
		"FailureOnInvalidMultipartName": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure to upload when multipart content field is not named as `image`",
				Request:     NewRequest(t, "image/svg+xml", "TestData/image.svg", "42d29b4b-935d-4f35-b26c-70080107f6d6", "invalid"),
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_ImageDataNotPresent.Ptr(),
				},
			}
		},
		"FailureOnNonImageContentType": func(t *testing.T) libAPI.TCData {
			return libAPI.TCData{
				Description: "Failure to upload when file content type doesn't have prefix `image/`",
				Request:     NewRequest(t, "invalid", "TestData/image.svg", "42d29b4b-935d-4f35-b26c-70080107f6d6", "image"),
				Response: libAPI.TCResponse{
					Status:    http.StatusBadRequest,
					ErrorCode: v1.Err400_ImageDataNotPresent.Ptr(),
				},
			}
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(tc.TestAPI, http.MethodPut, "/user/image"))
	}
}
