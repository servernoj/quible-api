package v1

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/models"
)

type GetUserProfileImageInput struct {
	UserId string `path:"userId"`
}

type GetUserProfileImageOutput struct {
	ContentType string `header:"content-type"`
	Body        []byte `doc:"binary content of the user's profile image"`
}

func (impl *VersionedImpl) RegisterGetUserProfileImage(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "get-user-image",
				Summary:     "Get user profile image",
				Description: "Return profile image (binary data) of the requested user",
				Method:      http.MethodGet,
				Errors: []int{
					http.StatusNotFound,
				},
				Tags: []string{"user", "public"},
				Path: "/user/{userId}/image",
			},
		),
		func(ctx context.Context, input *GetUserProfileImageInput) (*GetUserProfileImageOutput, error) {
			// 0. Dependences
			deps := impl.Deps.GetContext("opGetUserProfileImage")
			db := deps.Get("db").(*sql.DB)
			// 1. Retrieve the user based on `userId` path parameter
			user, err := models.FindUser(ctx, db, input.UserId)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err404_UserNotFound, err)
			}
			// 2. Retrieve profile image data
			var imageData ImageData
			if imageDataBytesPtr := user.Image.Ptr(); imageDataBytesPtr != nil {
				if err := json.Unmarshal(*imageDataBytesPtr, &imageData); err != nil {
					return nil, ErrorMap.GetErrorResponse(Err500_UnableToRetrieveProfileImage, err)
				}
			} else {
				return nil, ErrorMap.GetErrorResponse(Err404_UserHasNoImage)
			}
			// 3. Return response with binary image data and appropriate content type header
			response := &GetUserProfileImageOutput{
				ContentType: imageData.ContentType,
				Body:        imageData.BinaryContent,
			}
			return response, nil
		},
	)
}
