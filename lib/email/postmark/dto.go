package postmark

type AttachmentDTO struct {
	// Name: attachment name
	Name string `json:"name" binding:"required"`
	// Content: Base64 encoded attachment data
	Content string `json:"content" binding:"required"`
	// ContentType: attachment MIME type
	ContentType string `json:"content_type" binding:"required"`
}

type EmailDTO struct {
	// From: REQUIRED The sender email address. Must have a registered and confirmed Sender Signature.
	From string `json:"from" binding:"required,email"`
	// To: REQUIRED Recipient email address. Multiple addresses are comma separated. Max 50.
	To string `json:"to" binding:"required,email"`
	// Cc recipient email address. Multiple addresses are comma separated. Max 50.
	Cc string `json:"cc" binding:"omitempty,email"`
	// Bcc recipient email address. Multiple addresses are comma separated. Max 50.
	Bcc string `json:"bcc" binding:"omitempty,email"`
	// Subject: Email subject
	Subject string `json:"subject" binding:"required"`
	// HTMLBody: HTML email message. REQUIRED, If no TextBody specified
	HTMLBody string `json:"HtmlBody" binding:"required_without=TextBody"`
	// TextBody: Plain text email message. REQUIRED, If no HTMLBody specified
	TextBody string `json:"TextBody" binding:"required_without=HTMLBody"`
	// ReplyTo: Reply To override email address. Defaults to the Reply To set in the sender signature.
	ReplyTo string `json:"reply_to" binding:"omitempty,email"`
	// Attachments: List of attachments
	Attachments []AttachmentDTO `json:"attachments" binding:"dive,required"`
}
