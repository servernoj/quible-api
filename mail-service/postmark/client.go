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

type Client struct {
	http.Client
	Context      context.Context
	ServerToken  string
	AccountToken string
	BaseURL      string
}

func NewClient(ctx context.Context) *Client {
	return &Client{
		Client:       http.Client{Timeout: 10 * time.Second},
		ServerToken:  os.Getenv("ENV_POSTMARK_SERVER_TOKEN"),
		AccountToken: os.Getenv("ENV_POSTMARK_ACCOUNT_TOKEN"),
		BaseURL:      "https://api.postmarkapp.com",
		Context:      ctx,
	}
}

func (client *Client) SendEmail(email EmailDTO) (*PostmarkResponse, error) {
	return doRequest(
		client,
		RequestParams[EmailDTO]{
			Method:  http.MethodPost,
			Path:    "email",
			Payload: &email,
		},
	)
}

func doRequest[T PostmarkPayload](client *Client, params RequestParams[T]) (*PostmarkResponse, error) {

	var requestBody io.Reader
	if params.Payload != nil {
		marshalled, err := json.Marshal(params.Payload)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal request payload: %w", err)
		}
		requestBody = bytes.NewReader(marshalled)
	}

	req, err := http.NewRequestWithContext(
		client.Context,
		params.Method,
		fmt.Sprintf("%s/%s", client.BaseURL, params.Path),
		requestBody,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Postmark-Server-Token", client.ServerToken)
	// req.Header.Add("X-Postmark-Account-Token", client.AccountToken)

	res, err := client.Do(req)
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
