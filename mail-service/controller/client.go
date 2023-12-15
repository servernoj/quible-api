package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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

// Template represents an email template on the server
type Template struct {
	// TemplateID: ID of template
	TemplateID int64 `json:"TemplateID"`
	// Name: Name of template
	Name string
	// Subject: The content to use for the Subject when this template is used to send email.
	Subject string
	// HTMLBody: The content to use for the HTMLBody when this template is used to send email.
	HTMLBody string `json:"HtmlBody"`
	// TextBody: The content to use for the TextBody when this template is used to send email.
	TextBody string
	// AssociatedServerID: The ID of the Server with which this template is associated.
	AssociatedServerID int64 `json:"AssociatedServerId"`
	// Active: Indicates that this template may be used for sending email.
	Active bool
}

// TemplateInfo is a limited set of template info returned via Index/Editing endpoints
type TemplateInfo struct {
	// TemplateID: ID of template
	TemplateID int64 `json:"TemplateID"`
	// Name: Name of template
	Name string
	// Active: Indicates that this template may be used for sending email.
	Active bool
}

// GetTemplate fetches a specific template via TemplateID
func (client *Client) GetTemplate(ctx context.Context, templateID string) (Template, error) {
	res := Template{}
	err := client.doRequest(ctx, parameters{
		Method:    "GET",
		Path:      fmt.Sprintf("templates/%s", templateID),
		TokenType: serverToken,
	}, &res)
	return res, err
}

type templatesResponse struct {
	TotalCount int64
	Templates  []TemplateInfo
}

// ValidateTemplateBody contains the template/render model combination to be validated
type ValidateTemplateBody struct {
	Subject                    string
	TextBody                   string
	HTMLBody                   string `json:"HTMLBody"`
	TestRenderModel            map[string]interface{}
	InlineCSSForHTMLTestRender bool `json:"InlineCssForHtmlTestRender"`
}

// ValidateTemplateResponse contains information as to how the validation went
type ValidateTemplateResponse struct {
	AllContentIsValid      bool
	HTMLBody               Validation `json:"HTMLBody"`
	TextBody               Validation
	Subject                Validation
	SuggestedTemplateModel map[string]interface{}
}

// Validation contains the results of a field's validation
type Validation struct {
	ContentIsValid   bool
	ValidationErrors []ValidationError
	RenderedContent  string
}

// ValidationError contains information about the errors which occurred during validation for a given field
type ValidationError struct {
	Message           string
	Line              int
	CharacterPosition int
}

// TemplatedEmail is used to send an email via a template
type TemplatedEmail struct {
	// TemplateID: REQUIRED if TemplateAlias is not specified. - The template id to use when sending this message.
	TemplateID int64 `json:"TemplateId,omitempty"`
	// TemplateAlias: REQUIRED if TemplateID is not specified. - The template alias to use when sending this message.
	TemplateAlias string `json:",omitempty"`
	// TemplateModel: The model to be applied to the specified template to generate HtmlBody, TextBody, and Subject.
	TemplateModel map[string]interface{} `json:",omitempty"`
	// InlineCSS: By default, if the specified template contains an HtmlBody, we will apply the style blocks as inline attributes to the rendered HTML content. You may opt out of this behavior by passing false for this request field.
	InlineCSS bool `json:"InlineCSS,omitempty"`
	// From: The sender email address. Must have a registered and confirmed Sender Signature.
	From string `json:",omitempty"`
	// To: REQUIRED Recipient email address. Multiple addresses are comma separated. Max 50.
	To string `json:",omitempty"`
	// Cc recipient email address. Multiple addresses are comma separated. Max 50.
	Cc string `json:",omitempty"`
	// Bcc recipient email address. Multiple addresses are comma separated. Max 50.
	Bcc string `json:",omitempty"`
	// Tag: Email tag that allows you to categorize outgoing emails and get detailed statistics.
	Tag string `json:",omitempty"`
	// Reply To override email address. Defaults to the Reply To set in the sender signature.
	ReplyTo string `json:",omitempty"`
	// Headers: List of custom headers to include.
	Headers []Header `json:",omitempty"`
	// TrackOpens: Activate open tracking for this email.
	TrackOpens bool `json:",omitempty"`
	// TrackLinks: Activate link tracking. Possible options: "None", "HtmlAndText", "HtmlOnly", "TextOnly".
	TrackLinks string `json:",omitempty"`
	// Attachments: List of attachments
	Attachments []Attachment `json:",omitempty"`
	// MessageStream: MessageStream will default to the outbound message stream ID (Default Transactional Stream) if no message stream ID is provided.
	MessageStream string `json:",omitempty"`
}

// GetTemplates fetches a list of templates on the server
// It returns a TemplateInfo slice, the total template count, and any error that occurred
// Note: TemplateInfo only returns a subset of template attributes, use GetTemplate(id) to
// retrieve all template info.
func (client *Client) GetTemplates(ctx context.Context, count int64, offset int64) ([]TemplateInfo, int64, error) {
	res := templatesResponse{}

	values := &url.Values{}
	values.Add("count", fmt.Sprintf("%d", count))
	values.Add("offset", fmt.Sprintf("%d", offset))

	err := client.doRequest(ctx, parameters{
		Method:    "GET",
		Path:      fmt.Sprintf("templates?%s", values.Encode()),
		TokenType: serverToken,
	}, &res)
	return res.Templates, res.TotalCount, err
}

// CreateTemplate saves a new template to the server
func (client *Client) CreateTemplate(ctx context.Context, template Template) (TemplateInfo, error) {
	res := TemplateInfo{}
	err := client.doRequest(ctx, parameters{
		Method:    "POST",
		Path:      "templates",
		Payload:   template,
		TokenType: serverToken,
	}, &res)
	return res, err
}

// EditTemplate updates details for a specific template with templateID
func (client *Client) EditTemplate(ctx context.Context, templateID string, template Template) (TemplateInfo, error) {
	res := TemplateInfo{}
	err := client.doRequest(ctx, parameters{
		Method:    "PUT",
		Path:      fmt.Sprintf("templates/%s", templateID),
		Payload:   template,
		TokenType: serverToken,
	}, &res)
	return res, err
}

// DeleteTemplate removes a template (with templateID) from the server
func (client *Client) DeleteTemplate(ctx context.Context, templateID string) error {
	res := APIError{}
	err := client.doRequest(ctx, parameters{
		Method:    "DELETE",
		Path:      fmt.Sprintf("templates/%s", templateID),
		TokenType: serverToken,
	}, &res)

	if res.ErrorCode != 0 {
		return res
	}

	return err
}

// ValidateTemplate validates the provided template/render model combination
func (client *Client) ValidateTemplate(ctx context.Context, validateTemplateBody ValidateTemplateBody) (ValidateTemplateResponse, error) {
	res := ValidateTemplateResponse{}
	err := client.doRequest(ctx, parameters{
		Method:    "POST",
		Path:      "templates/validate",
		Payload:   validateTemplateBody,
		TokenType: serverToken,
	}, &res)
	return res, err
}

// SendTemplatedEmail sends an email using a template (TemplateID)
func (client *Client) SendTemplatedEmail(ctx context.Context, email TemplatedEmail) (EmailResponse, error) {
	res := EmailResponse{}
	err := client.doRequest(ctx, parameters{
		Method:    "POST",
		Path:      "email/withTemplate",
		Payload:   email,
		TokenType: serverToken,
	}, &res)
	return res, err
}

// SendTemplatedEmailBatch sends batch email using a template (TemplateID)
func (client *Client) SendTemplatedEmailBatch(ctx context.Context, emails []TemplatedEmail) ([]EmailResponse, error) {
	var res []EmailResponse
	formatEmails := map[string]interface{}{
		"Messages": emails,
	}
	err := client.doRequest(ctx, parameters{
		Method:    "POST",
		Path:      "email/batchWithTemplates",
		Payload:   formatEmails,
		TokenType: serverToken,
	}, &res)
	return res, err
}

// Error returns the error message details
func (res APIError) Error() string {
	return res.Message
}

// initialize a newclient function
func NewClient() *Client {
	serverToken := os.Getenv("ENV_POSTMARK_SERVER_TOKEN")
	accountToken := os.Getenv("ENV_POSTMARK_ACCOUNT_TOKEN")

	return &Client{
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
		ServerToken:  serverToken,
		AccountToken: accountToken,
		BaseURL:      "https://api.postmarkapp.com",
	}
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

// SendEmail logic
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
