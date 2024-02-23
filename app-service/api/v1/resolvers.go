package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/quible-io/quible-api/lib/models"
)

type _resolved struct {
	resolved bool
}

// -- Authorization header containing Bearer access token. Injects `UserId` into `input` struct
type AuthorizationHeaderResolver struct {
	_resolved
	Authorization string `header:"authorization"`
	UserId        string
}

func (input *AuthorizationHeaderResolver) Resolve(ctx huma.Context) (errs []error) {
	if !input.resolved {
		input.resolved = true
	} else {
		return
	}
	request, _ := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s/api/v1/user",
			os.Getenv("ENV_URL_AUTH_SERVICE"),
		),
		http.NoBody,
	)
	request.Header.Add("Authorization", input.Authorization)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		errs = append(errs, &huma.ErrorDetail{
			Message:  "unable to send request to auth-service",
			Location: "header.authorization.request",
			Value:    err,
		})
		return
	}
	body := response.Body
	defer body.Close()
	var data map[string]any
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		errs = append(errs, &huma.ErrorDetail{
			Message:  "unable to parse response from auth-service",
			Location: "auth-service.getUser.body",
			Value:    err,
		})
		return
	}
	if response.StatusCode == http.StatusUnauthorized {
		errs = append(errs, &huma.ErrorDetail{
			Message:  "insufficient privilege",
			Location: "auth-service.getUser.status",
			Value:    response.StatusCode,
		})
		return
	}
	if userId, ok := data["id"].(string); !ok {
		errs = append(errs, &huma.ErrorDetail{
			Message:  "field `id` is not present in the returned object",
			Location: "auth-service.getUser.data",
			Value:    data,
		})
		return
	} else {
		if exists, err := models.UserExistsG(ctx.Context(), userId); err != nil || !exists {
			errs = append(errs, &huma.ErrorDetail{
				Message:  "user not found",
				Location: "db.users",
				Value:    err,
			})
			return
		}
		input.UserId = userId
	}
	return
}