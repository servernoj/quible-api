package email

import (
	"context"
	"fmt"

	"github.com/quible-io/quible-api/lib/email/postmark"
	"github.com/rs/zerolog/log"
)

type EmailDTO = postmark.EmailDTO

func Send(ctx context.Context, email EmailDTO) error {

	var NewClient = postmark.NewClient

	response, err := NewClient(ctx).SendEmail(email)

	if err != nil || response.ErrorCode != 0 {
		return fmt.Errorf("unable to send email via Postmark: %w", err)
	}

	log.Info().Msgf("Email sent: %s", response.MessageID)
	return nil
}
