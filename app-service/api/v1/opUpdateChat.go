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

type UpdateChatInput struct {
	AuthorizationHeaderResolver
	ChatId string `path:"chatId"`
	Body   struct {
		Name      *string `json:"name,omitempty" pattern:"\\w+" doc:"resource name"`
		Title     *string `json:"title,omitempty" doc:"human-readable 'title' of the chat group/channel"`
		Summary   *string `json:"summary,omitempty" doc:"optional summary, to clear send empty string"`
		IsPrivate *bool   `json:"isPrivate,omitempty" doc:"only for chat groups"`
	}
}

type UpdateChatOutput struct {
	Body models.Chat
}

func (impl *VersionedImpl) RegisterUpdateChat(api huma.API, vc libAPI.VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "patch-update-chat",
				Summary:     "Update chat record",
				Description: "Update chat record (group/channel) with provided details",
				Method:      http.MethodPatch,
				Errors: []int{
					http.StatusBadRequest,
					http.StatusUnauthorized,
				},
				DefaultStatus: http.StatusOK,
				Tags:          []string{"chat", "protected"},
				Path:          "/chat/{chatId}",
			},
		),
		func(ctx context.Context, input *UpdateChatInput) (*UpdateChatOutput, error) {
			// 1. Retrieve the chat record for update
			chat, err := models.FindChatG(ctx, input.ChatId)
			if err != nil {
				return nil, ErrorMap.GetErrorResponse(Err404_ChatRecordNotFound, err)
			}
			if chat.ParentID.Ptr() != nil && input.Body.IsPrivate != nil {
				// for chat channels IsPrivate cannot be set
				return nil, ErrorMap.GetErrorResponse(Err400_OnlyForChatGroups, err)
			}
			// 2. Update chat record with respect to provided PATCH data
			patchDataType := reflect.TypeOf(input.Body)
			patchDataValue := reflect.ValueOf(input.Body)
			userValue := reflect.ValueOf(chat).Elem()
			for i := 0; i < patchDataValue.NumField(); i++ {
				fieldName := patchDataType.Field(i).Name
				fieldValue := patchDataValue.Field(i).Elem()
				if fieldValue.IsValid() {
					target := userValue.FieldByName(fieldName)
					if target.Kind().String() == "struct" {
						switch target.Type().String() {
						case "null.String":
							target.FieldByName("String").SetString(fieldValue.String())
							target.FieldByName("Valid").SetBool(fieldValue.String() != "")
						case "null.Bool":
							target.FieldByName("Bool").SetBool(fieldValue.Bool())
							target.FieldByName("Valid").SetBool(true)
						}
					}
				}
			}
			// 3. Store updated chat record
			if _, err := chat.UpdateG(ctx, boil.Infer()); err != nil {
				return nil, ErrorMap.GetErrorResponse(Err500_UnableUpdateChatRecord, err)
			}
			// 4. Prepare and return the response
			response := &UpdateChatOutput{
				Body: *chat,
			}
			return response, nil
		},
	)
}
