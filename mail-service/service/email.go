package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	accountToken = "account"
	serverToken  = "server"
)

type Email struct {
	// From: REQUIRED The sender email address. Must have a registered and confirmed Sender Signature.
	From string `json:",omitempty"`
	// To: REQUIRED Recipient email address. Multiple addresses are comma separated. Max 50.
	To string `json:",omitempty"`
	// Cc recipient email address. Multiple addresses are comma separated. Max 50.
	Cc string `json:",omitempty"`
	// Bcc recipient email address. Multiple addresses are comma separated. Max 50.
	Bcc string `json:",omitempty"`
	// Subject: Email subject
	Subject string `json:",omitempty"`
	// Tag: Email tag that allows you to categorize outgoing emails and get detailed statistics.
	Tag string `json:",omitempty"`
	// HTMLBody: HTML email message. REQUIRED, If no TextBody specified
	HTMLBody string `json:"HtmlBody,omitempty"`
	// TextBody: Plain text email message. REQUIRED, If no HTMLBody specified
	TextBody string `json:",omitempty"`
	// ReplyTo: Reply To override email address. Defaults to the Reply To set in the sender signature.
	ReplyTo string `json:",omitempty"`
	// Headers: List of custom headers to include.
	Headers []Header `json:",omitempty"`
	// TrackOpens: Activate open tracking for this email.
	TrackOpens bool `json:",omitempty"`
	// TrackLinks:Activate link tracking for links in the HTML or Text bodies of this email. Possible options: None HtmlAndText HtmlOnly TextOnly
	TrackLinks string `json:",omitempty"`
	// Attachments: List of attachments
	Attachments []Attachment `json:",omitempty"`
	// Metadata: metadata
	Metadata map[string]string `json:",omitempty"`
	// MessageStream: MessageStream will default to the outbound message stream ID (Default Transactional Stream) if no message stream ID is provided.
	MessageStream string `json:",omitempty"`
}

// Client provides a connection to the Postmark API
type Client struct {
	// HTTPClient is &http.Client{} by default
	HTTPClient *http.Client
	// Server Token: Used for requests that require server level privileges. This token can be found on the Credentials tab under your Postmark server.
	ServerToken string
	// AccountToken: Used for requests that require account level privileges. This token is only accessible by the account owner, and can be found on the Account tab of your Postmark account.
	AccountToken string
	// BaseURL is the root API endpoint
	BaseURL string
}

// Options is an object to hold variable parameters to perform request.
type parameters struct {
	// Method is HTTP method type.
	Method string
	// Path is postfix for URI.
	Path string
	// Payload for the request.
	Payload interface{}
	// TokenType defines which token to use
	TokenType string
}

// Header - an email header
type Header struct {
	// Name: header name
	Name string
	// Value: header value
	Value string
}

// Attachment is an optional encoded file to send along with an email
type Attachment struct {
	// Name: attachment name
	Name string
	// Content: Base64 encoded attachment data
	Content string
	// ContentType: attachment MIME type
	ContentType string
	// ContentId: populate for inlining images with the images cid
	ContentID string `json:",omitempty"`
}

// EmailResponse holds info in response to a send/send-batch request
// Even if API request comes back successful, check the ErrorCode to see if there might be a delivery problem
type EmailResponse struct {
	// To: Recipient email address
	To string
	// SubmittedAt: Timestamp
	SubmittedAt time.Time
	// MessageID: ID of message
	MessageID string
	// ErrorCode: see error codes here (https://postmarkapp.com/developer/api/overview#error-codes)
	ErrorCode int64
	// Message: Response message
	Message string
}

// APIError represents errors returned by Postmark
type APIError struct {
	// ErrorCode: see error codes here (https://postmarkapp.com/developer/api/overview#error-codes)
	ErrorCode int64 `json:"ErrorCode"`
	// Message contains error details
	Message string `json:"Message"`
}

// Error returns the error message details
func (res APIError) Error() string {
	return res.Message
}

// doRequest performs the request to the Postmark API
func (client *Client) doRequest(ctx context.Context, opts parameters, dst interface{}) (err error) {
	url := fmt.Sprintf("%s/%s", client.BaseURL, opts.Path)

	var req *http.Request
	if req, err = http.NewRequestWithContext(
		ctx, opts.Method, url, nil,
	); err != nil {
		return
	}

	if opts.Payload != nil {
		var payloadData []byte
		if payloadData, err = json.Marshal(opts.Payload); err != nil {
			return
		}
		req.Body = io.NopCloser(bytes.NewBuffer(payloadData))
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	switch opts.TokenType {
	case accountToken:
		req.Header.Add("X-Postmark-Account-Token", client.AccountToken)
	default:
		req.Header.Add("X-Postmark-Server-Token", client.ServerToken)
	}

	var res *http.Response
	if res, err = client.HTTPClient.Do(req); err != nil {
		return
	}

	defer func() {
		_ = res.Body.Close()
	}()
	var body []byte
	if body, err = io.ReadAll(res.Body); err != nil {
		return
	}

	if res.StatusCode >= http.StatusBadRequest {
		// If the status code is not a success, attempt to unmarshall the body into the APIError struct.
		var apiErr APIError
		if err = json.Unmarshal(body, &apiErr); err != nil {
			return
		}
		return apiErr
	}

	return json.Unmarshal(body, dst)
}

// SendEmail
func (client *Client) SendEmail(ctx context.Context, email Email) (EmailResponse, error) {
	res := EmailResponse{}
	err := client.doRequest(ctx, parameters{
		Method:    "POST",
		Path:      "email",
		Payload:   email,
		TokenType: serverToken,
	}, &res)

	if res.ErrorCode != 0 {
		return res, fmt.Errorf(`%v %s`, res.ErrorCode, res.Message)
	}

	return res, err
}

// SendEmailBatch sends multiple emails together
// Note, individual emails in the batch can error, so it would be wise to
// range over the responses and sniff for errors
func (client *Client) SendEmailBatch(ctx context.Context, emails []Email) ([]EmailResponse, error) {
	var res []EmailResponse
	err := client.doRequest(ctx, parameters{
		Method:    "POST",
		Path:      "email/batch",
		Payload:   emails,
		TokenType: serverToken,
	}, &res)
	return res, err
}
