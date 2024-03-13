package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog/log"
)

// -- Authorization header containing Bearer access token. Injects `UserId` into `input` struct
type AuthorizationHeaderResolver struct {
	Authorization string `header:"authorization"`
	UserId        string
}

func (input *AuthorizationHeaderResolver) Resolve(ctx huma.Context) (errs []error) {
	// 1. Prepare request
	log.Info().Msg(os.Getenv("ENV_URL_AUTH_SERVICE"))
	request, _ := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s/api/v1/user",
			os.Getenv("ENV_URL_AUTH_SERVICE"),
		),
		http.NoBody,
	)
	request.Header.Add("Authorization", input.Authorization)
	// 2. Initialize HTTP client
	var httpClient = http.DefaultClient
	// 3. Perform the request
	response, err := httpClient.Do(request)
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
	// 4. Parse the response and check if status is not 401
	var data struct {
		ID string `json:"id"`
	}
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
	input.UserId = data.ID
	return
}
