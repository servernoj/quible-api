package email

import "context"

type Attachment struct {
	Name        string `json:"name"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
}

type EmailPayload struct {
	From        string        `json:"from"`
	To          string        `json:"to"`
	Subject     string        `json:"subject"`
	HTMLBody    string        `json:"HtmlBody"`
	TextBody    string        `json:"TextBody"`
	Attachments *[]Attachment `json:"attachments,omitempty"`
}

type EmailSender interface {
	SendEmail(context.Context, EmailPayload) error
}
