package v1

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/models"
)

type GetUserProfileInput struct {
	UserId string `path:"userId"`
}

type GetUserProfileOutput struct {
	Body UserProfile
}

func (impl *VersionedImpl) RegisterGetUserProfile(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "get-user-profile-by-id",
				Summary:     "Get user profile",
				Description: "Return public user profile based on provided `userId`",
				Method:      http.MethodGet,
				Errors: []int{
					http.StatusUnauthorized,
					http.StatusNotFound,
				},
				Tags: []string{"user", "public"},
				Path: "/user/{userId}/profile",
			},
		),
		func(ctx context.Context, input *GetUserProfileInput) (*GetUserProfileOutput, error) {
			// 0. Dependences
			deps := impl.Deps.GetContext("opGetUserProfile")
			db := deps.Get("db").(*sql.DB)
			// 1. Retrieve the user based on `userId` path parameter
			user, err := models.FindUser(ctx, db, input.UserId)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err404_UserNotFound, err)
			}
			// 2. Evaluate `image` field on the found user record and craft data URL containing profile image
			var image *string
			if imageDataBytesPtr := user.Image.Ptr(); imageDataBytesPtr != nil {
				var imageData ImageData
				if err := json.Unmarshal(*imageDataBytesPtr, &imageData); err != nil {
					return nil, ErrorMap.GetErrorResponse(Err500_UnableToRetrieveProfileImage, err)
				}
				imageDataURL := fmt.Sprintf(
					"data:%s;base64,%s",
					imageData.ContentType,
					base64.StdEncoding.EncodeToString(imageData.BinaryContent),
				)
				image = &imageDataURL
			}
			// 3. Return "public" user profile comprised from id, name and profile image (if found)
			response := &GetUserProfileOutput{
				Body: UserProfile{
					ID:       user.ID,
					FullName: user.FullName,
					Image:    image,
				},
			}
			return response, nil
		},
	)
}
