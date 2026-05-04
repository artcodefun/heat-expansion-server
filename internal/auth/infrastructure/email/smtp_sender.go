package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"mime"
	"net"
	"net/smtp"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

const dialTimeout = 15 * time.Second

type SMTPSender struct {
	host     string
	user     string
	password string
	from     string
}

func NewSMTPSender(host, user, password, from string) *SMTPSender {
	return &SMTPSender{host: host, user: user, password: password, from: from}
}

func (s *SMTPSender) SendPasswordReset(ctx context.Context, toEmail, toName, code string) error {
	ctx, span := otel.Tracer("heat-expansion-auth").Start(ctx, "email.send_password_reset")
	defer span.End()
	subject := mime.QEncoding.Encode("utf-8", "Heat Expansion — password reset code")
	body := fmt.Sprintf(
		"Hi %s,\r\n\r\nYour password reset code is:\r\n\r\n    %s\r\n\r\nEnter it in the game client. The code expires in 1 hour.\r\n\r\nIf you didn't request this, you can safely ignore this email.",
		toName, code,
	)
	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.from, toEmail, subject, body,
	))

	err := s.send(toEmail, msg)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

func (s *SMTPSender) send(toEmail string, msg []byte) error {
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: dialTimeout},
		"tcp",
		fmt.Sprintf("%s:465", s.host),
		&tls.Config{ServerName: s.host},
	)
	if err != nil {
		return fmt.Errorf("smtp dial: %w", err)
	}

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Close()

	if err := client.Auth(smtp.PlainAuth("", s.user, s.password, s.host)); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}
	if err := client.Mail(s.from); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}
	if err := client.Rcpt(toEmail); err != nil {
		return fmt.Errorf("smtp rcpt: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	return w.Close()
}
