package ports

import "context"

type EmailSender interface {
	SendPasswordReset(ctx context.Context, toEmail, toName, rawToken string) error
}
