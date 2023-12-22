package postmark

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Mailer struct {
	http.Client
	Context      context.Context
	ServerToken  string
	AccountToken string
	BaseURL      string
}

func NewMailer(ctx context.Context) *Mailer {
	return &Mailer{
		Client:       http.Client{Timeout: 10 * time.Second},
		ServerToken:  os.Getenv("ENV_POSTMARK_SERVER_TOKEN"),
		AccountToken: os.Getenv("ENV_POSTMARK_ACCOUNT_TOKEN"),
		BaseURL:      "https://api.postmarkapp.com",
		Context:      ctx,
	}
}

func (mailer *Mailer) SendEmail(email EmailDTO) (*PostmarkResponse, error) {
	return doRequest(
		mailer,
		RequestParams[EmailDTO]{
			Method:  http.MethodPost,
			Path:    "email",
			Payload: &email,
		},
	)
}

func doRequest[T PostmarkPayload](mailer *Mailer, params RequestParams[T]) (*PostmarkResponse, error) {

	var requestBody io.Reader
	if params.Payload != nil {
		marshalled, err := json.Marshal(params.Payload)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal request payload: %w", err)
		}
		requestBody = bytes.NewReader(marshalled)
	}

	req, err := http.NewRequestWithContext(
		mailer.Context,
		params.Method,
		fmt.Sprintf("%s/%s", mailer.BaseURL, params.Path),
		requestBody,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Postmark-Server-Token", mailer.ServerToken)
	// req.Header.Add("X-Postmark-Account-Token", client.AccountToken)

	res, err := mailer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to execute request to %q: %w", req.URL, err)
	}

	defer res.Body.Close()
	var ParsedResponse PostmarkResponse
	if err := json.NewDecoder(res.Body).Decode(&ParsedResponse); err != nil {
		return nil, fmt.Errorf("unable to parse response from %q: %w", req.URL, err)
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("api error: %w", errors.New(ParsedResponse.String()))
	}

	return &ParsedResponse, nil
}
