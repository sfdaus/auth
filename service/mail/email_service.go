package service

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v2"
)

type EmailService interface {
	SendEmail(ctx context.Context, to string, subject string, htmlBody string) error
}

type resendEmailService struct {
	client *resend.Client
	from   string
}

func NewResendEmailService(apiKey, from string) EmailService {
	client := resend.NewClient(apiKey)
	return &resendEmailService{
		client: client,
		from:   from,
	}
}

func (r *resendEmailService) SendEmail(ctx context.Context, to string, subject string, htmlBody string) error {
	params := &resend.SendEmailRequest{
		From:    r.from,
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
	}

	_, err := r.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
