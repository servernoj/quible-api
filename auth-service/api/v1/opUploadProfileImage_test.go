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
	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/misc"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/store"
	"github.com/rs/zerolog/log"
)

func (suite *TestCases) TestUploadProfileImage() {
	t := suite.T()
	// 1. Import users from CSV file
	store.InsertFromCSV(t, "users", UsersCSV)
	NewTCRequest := func(t *testing.T, contentType, imageFilename, userId, fieldName string) TCRequest {
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
			fmt.Sprintf("Authorization: Bearer %s", GetToken(t, userId, jwt.TokenActionAccess)),
			bytes.NewReader(body.Bytes()),
		}
		return TCRequest{
			Args: args,
			Params: map[string]any{
				"userId":        userId,
				"contentType":   contentType,
				"imageFilename": imageFilename,
			},
		}
	}
	// 2. Define test scenarios
	testCases := TCScenarios{
		"Success": TCData{
			Description: "Successful upload and confirmation in DB",
			Request:     NewTCRequest(t, "image/svg+xml", "TestData/image.svg", "42d29b4b-935d-4f35-b26c-70080107f6d6", "image"),
			Response: TCResponse{
				Status: http.StatusAccepted,
			},
			ExtraTests: []TCExtraTest{
				func(t TCRequest, res *httptest.ResponseRecorder) bool {
					user, err := models.FindUserG(context.Background(), t.Params["userId"].(string))
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
		},
		"FailureOnInvalidMultipartName": TCData{
			Description: "Failure to upload when multipart content field is not named as `image`",
			Request:     NewTCRequest(t, "image/svg+xml", "TestData/image.svg", "42d29b4b-935d-4f35-b26c-70080107f6d6", "invalid"),
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_ImageDataNotPresent),
			},
		},
		"FailureOnInvalidContentType": TCData{
			Description: "Failure to upload when file content type doesn't have prefix `image/`",
			Request:     NewTCRequest(t, "invalid", "TestData/image.svg", "42d29b4b-935d-4f35-b26c-70080107f6d6", "image"),
			Response: TCResponse{
				Status:    http.StatusBadRequest,
				ErrorCode: misc.Of(v1.Err400_ImageDataNotPresent),
			},
		},
	}
	// 3. Run scenarios in sequence
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(suite.TestAPI, http.MethodPut, "/api/v1/user/image"))
	}
}
