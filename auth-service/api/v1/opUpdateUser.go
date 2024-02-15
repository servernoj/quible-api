package v1

import (
	"context"
	"net/http"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	libAPI "github.com/quible-io/quible-api/lib/api"
	"github.com/quible-io/quible-api/lib/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UpdateUserInput struct {
	AuthorizationHeaderResolver
	Body struct {
		Username *string `json:"username,omitempty"`
		Email    *string `json:"email,omitempty" format:"email"`
		FullName *string `json:"full_name,omitempty" minLength:"1"`
		Phone    *string `json:"phone,omitempty" pattern:"^[0-9() +-]{10,}$"`
	}
}

type UpdateUserOutput struct {
	Body UserSimplified
}

func (impl *VersionedImpl) RegisterUpdateUser(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "patch-update-user",
				Summary:     "Patch user record",
				Description: "Update user record with provided details",
				Method:      http.MethodPatch,
				Errors: []int{
					http.StatusBadRequest,
					http.StatusUnauthorized,
				},
				DefaultStatus: http.StatusOK,
				Tags:          []string{"user", "protected"},
				Path:          "/user",
			},
		),
		func(ctx context.Context, input *UpdateUserInput) (*UpdateUserOutput, error) {
			// 1. Retrieve the user record for update
			user, err := models.FindUserG(ctx, input.UserId)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err401_UserNotFound, err)
			}
			// 2. Update user record with respect to provided PATCH data
			patchDataType := reflect.TypeOf(input.Body)
			patchDataValue := reflect.ValueOf(input.Body)
			userValue := reflect.ValueOf(user).Elem()
			for i := 0; i < patchDataValue.NumField(); i++ {
				fieldName := patchDataType.Field(i).Name
				fieldValue := patchDataValue.Field(i).Elem()
				if fieldValue.IsValid() {
					target := userValue.FieldByName(fieldName)
					if target.Kind().String() == "struct" && target.Type().String() == "null.String" {
						target = target.FieldByName("String")
					}
					if target.CanSet() {
						target.SetString(fieldValue.String())
					}
				}
			}
			// 3. Store updated user record
			if _, err := user.UpdateG(ctx, boil.Infer()); err != nil {
				return nil, ErrorMap.GetErrorResponse(Err500_UnableToUpdateUser, err)
			}
			// 4. Prepare and return the response
			response := &UpdateUserOutput{
				Body: UserSimplified{
					ID:       user.ID,
					Username: user.Username,
					Email:    user.Email,
					Phone:    user.Phone,
					FullName: user.FullName,
				},
			}
			return response, nil
		},
	)
}
