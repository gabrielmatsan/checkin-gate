package mail

import "context"




type Attachment struct {
	Filename    string
	Content     []byte
	ContentType string
}

type SendEmailParams struct {
	To          string
	Subject     string
	Body        string
	Attachments []Attachment
}

type EmailService interface {
	Send(ctx context.Context, params SendEmailParams) error
}
