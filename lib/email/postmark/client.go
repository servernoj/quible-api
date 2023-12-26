package postmark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Client struct {
	http.Client
	apiKey  string
	BaseURL string
	Context context.Context
}

func NewClient(ctx context.Context) *Client {
	return &Client{
		Client:  http.Client{Timeout: 10 * time.Second},
		apiKey:  os.Getenv("ENV_POSTMARK_API_KEY"),
		BaseURL: "https://api.postmarkapp.com",
		Context: ctx,
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
	req.Header.Add("X-Postmark-Server-Token", client.apiKey)

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
		return nil, ParsedResponse
	}

	return &ParsedResponse, nil
}
