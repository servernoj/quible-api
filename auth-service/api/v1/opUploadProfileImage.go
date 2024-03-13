package v1

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UploadProfileImageInput struct {
	AuthorizationHeaderResolver
	ContentType string `header:"Content-Type"`
	RawBody     []byte
	*ImageData
}

func (input *UploadProfileImageInput) Resolve(ctx huma.Context) (errs []error) {
	if errs = input.AuthorizationHeaderResolver.Resolve(ctx); len(errs) > 0 {
		return
	}
	// 1. Analyze content-type header
	mediaType, params, err := mime.ParseMediaType(input.ContentType)
	if err != nil {
		log.Error().Err(err).Send()
		errs = append(errs, &huma.ErrorDetail{
			Message:  err.Error(),
			Location: "header.content-type",
			Value:    input.ContentType,
		})
		return
	}
	// 2. For multipart/* content...
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(bytes.NewReader(input.RawBody), params["boundary"])
		// 3. Iterate over all identifiable body parts
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				errs = append(errs, &huma.ErrorDetail{
					Message:  err.Error(),
					Location: "body.part",
					Value:    p,
				})
				return
			}
			// 3a. Read the associated data
			slurp, err := io.ReadAll(p)
			if err != nil {
				errs = append(errs, &huma.ErrorDetail{
					Message:  err.Error(),
					Location: fmt.Sprintf("body.part.%s", p.FormName()),
					Value:    p,
				})
				return
			}
			// 3b. Format resulting ImageData struct to include the parsed information
			if p.FormName() == "image" && strings.HasPrefix(p.Header.Get("Content-Type"), "image") {
				input.ImageData = &ImageData{
					ContentType:   p.Header.Get("Content-Type"),
					BinaryContent: slurp,
				}
				break
			}
		}
	}
	return nil
}

type UploadProfileImageOutput struct {
}

func (impl *VersionedImpl) RegisterUploadProfileImage(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "put-upload-profile-image",
				Summary:     "Upload profile image",
				Description: "Upload profile image for the logged in user",
				Method:      http.MethodPut,
				Errors: []int{
					http.StatusBadRequest,
					http.StatusUnauthorized,
				},
				DefaultStatus: http.StatusAccepted,
				Tags:          []string{"user", "protected"},
				Path:          "/user/image",
			},
		),
		func(ctx context.Context, input *UploadProfileImageInput) (*UploadProfileImageOutput, error) {
			// 0. Dependences
			deps := impl.Deps.GetContext("opUploadProfileImage")
			db := deps.Get("db").(*sql.DB)
			// 1. Analyze result of resolver execution
			if input.ImageData == nil {
				return nil, ErrorMap.GetErrorResponse(Err400_ImageDataNotPresent)
			}
			if len(input.ImageData.BinaryContent) > 1*1024*1024 {
				return nil, ErrorMap.GetErrorResponse(Err400_FileTooLarge)
			}
			// 2. Locate the user record to be updated with image data
			user, err := models.FindUser(ctx, db, input.UserId)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err401_InvalidAccessToken, err)
			}
			b, err := json.Marshal(input.ImageData)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err500_UnableToStoreImage, err)
			}
			user.Image = null.BytesFrom(b)
			// 3. Update user's record
			if _, err := user.Update(ctx, db, boil.Infer()); err != nil {
				return nil, ErrorMap.GetErrorResponse(Err500_UnableToStoreImage, err)
			}
			return nil, nil
		},
	)
}
