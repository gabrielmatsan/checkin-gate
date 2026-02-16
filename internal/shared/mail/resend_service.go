package mail

import (
	"bytes"
	"context"
	"fmt"

	"github.com/wneessen/go-mail"
)

const (
	resendHost = "smtp.resend.com"
	resendPort = 587
	resendUser = "resend"
)

type ResendService struct {
	apiKey string
	from   string
}

func NewResendService(apiKey, from string) *ResendService {
	return &ResendService{
		apiKey: apiKey,
		from:   from,
	}
}

func (s *ResendService) Send(ctx context.Context, params SendEmailParams) error {
	msg := mail.NewMsg()

	if err := msg.From(s.from); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}

	if err := msg.To(params.To); err != nil {
		return fmt.Errorf("failed to set to address: %w", err)
	}

	msg.Subject(params.Subject)
	msg.SetBodyString(mail.TypeTextHTML, params.Body)

	for _, attachment := range params.Attachments {
		if err := msg.AttachReader(attachment.Filename, bytes.NewReader(attachment.Content)); err != nil {
			return fmt.Errorf("failed to attach attachment: %w", err)
		}
	}

	client, err := mail.NewClient(resendHost,
		mail.WithPort(resendPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(resendUser),
		mail.WithPassword(s.apiKey),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
	)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}

	if err := client.DialAndSendWithContext(ctx, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
