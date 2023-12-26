package email

import (
	"context"
	"fmt"
	"log"

	"github.com/quible-io/quible-api/lib/email/postmark"
)

func Send(ctx context.Context, email postmark.EmailDTO) error {

	var NewClient = postmark.NewClient

	response, err := NewClient(ctx).SendEmail(email)

	if err != nil || response.ErrorCode != 0 {
		return fmt.Errorf("unable to send email via Postmark: %w", err)
	}

	log.Printf("Email sent: %s", response.MessageID)
	return nil
}
