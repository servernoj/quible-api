package v1_test

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"testing"

	"github.com/quible-io/quible-api/lib/jwt"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/quible-io/quible-api/lib/store"
	"github.com/rs/zerolog/log"
)

func NewTCRequest(t *testing.T) TCRequest {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	h := make(textproto.MIMEHeader)
	h.Set(
		"Content-Disposition",
		`form-data; name="image"; filename="image.svg"`,
	)
	h.Set(
		"Content-Type",
		"image/svg+xml",
	)
	part, err := writer.CreatePart(h)
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open("TestData/image.svg")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if _, err := io.Copy(part, file); err != nil {
		t.Fatal(err)
	}
	writer.Close()
	userId := "42d29b4b-935d-4f35-b26c-70080107f6d6"
	args := []any{
		fmt.Sprintf("Content-Length: %d", body.Len()),
		fmt.Sprintf("Content-Type: multipart/form-data; boundary=%s", writer.Boundary()),
		fmt.Sprintf("Authorization: Bearer %s", GetToken(t, userId, jwt.TokenActionAccess)),
		bytes.NewReader(body.Bytes()),
	}
	return TCRequest{
		Method: http.MethodPut,
		Path:   "/api/v1/user/image",
		Args:   args,
		Params: map[string]any{
			"userId": userId,
		},
	}
}

func (suite *TestCases) TestUploadProfileImage() {
	t := suite.T()
	// 1. Import users from CSV file
	store.InsertFromCSV(t, "users", UsersCSV)
	// 2. Define test scenarios
	testCases := TCScenarios{
		"Success": TCData{
			Description: "login with correct credentials and expect success",
			Request:     NewTCRequest(t),
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
					return user.Image.Ptr() != nil
				},
			},
		},
	}
	// 3. Try different login scenarios
	for name, scenario := range testCases {
		t.Run(name, scenario.GetRunner(suite.TestAPI))
	}
}
