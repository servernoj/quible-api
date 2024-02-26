package postmark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/quible-io/quible-api/lib/email"
)

type Option func(client *PostmarkClient)

func WithBaseUrl(BaseURL string) Option {
	return func(postmarkClient *PostmarkClient) {
		postmarkClient.BaseURL = BaseURL
	}
}
func WithHttpClient(httpClient http.Client) Option {
	return func(postmarkClient *PostmarkClient) {
		postmarkClient.Client = httpClient
	}
}

type PostmarkClient struct {
	http.Client
	apiKey  string
	BaseURL string
}

func NewClient(options ...Option) email.EmailSender {
	postmarkClient := PostmarkClient{
		Client:  http.Client{Timeout: 10 * time.Second},
		apiKey:  os.Getenv("ENV_POSTMARK_API_KEY"),
		BaseURL: "https://api.postmarkapp.com",
	}
	for _, option := range options {
		option(&postmarkClient)
	}
	return &postmarkClient
}

func (postmarkClient *PostmarkClient) SendEmail(ctx context.Context, emailPayload email.EmailPayload) error {
	b, err := json.Marshal(emailPayload)
	if err != nil {
		return fmt.Errorf("unable to marshal payload: %w", err)
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/email", postmarkClient.BaseURL),
		bytes.NewBuffer(b),
	)
	if err != nil {
		return fmt.Errorf("unable to prepare request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Postmark-Server-Token", postmarkClient.apiKey)
	res, err := postmarkClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable send request: %w", err)
	}
	defer res.Body.Close()
	var parsedResponse Response
	if err := json.NewDecoder(res.Body).Decode(&parsedResponse); err != nil {
		return fmt.Errorf("unable to parse response: %w", err)
	}
	if res.StatusCode >= 400 {
		return parsedResponse
	}
	return nil
}
