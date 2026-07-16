package notification

import (
	"fmt"
	"net/smtp"

	"github.com/openpayment/gateway/internal/config"
	"github.com/rs/zerolog/log"
)

type EmailSender struct {
	config config.SMTPConfig
	domain string
}

func NewEmailSender(cfg config.SMTPConfig, frontendURL string) *EmailSender {
	return &EmailSender{config: cfg, domain: frontendURL}
}

func (e *EmailSender) SendResetToken(to, token string) error {
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", e.domain, token)

	if e.config.Host == "" {
		log.Warn().
			Str("to", to).
			Str("reset_link", resetLink).
			Msg("SMTP not configured — password reset token logged instead of emailed")
		return nil
	}

	subject := "Password Reset — Open Payment Gateway"
	body := fmt.Sprintf(`Hello,

You requested a password reset for your Open Payment Gateway account.

Click the link below to reset your password:

%s

If you did not request this, please ignore this email.

This link expires in 1 hour.

— Open Payment Gateway`, resetLink)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s",
		e.config.From, to, subject, body)

	addr := fmt.Sprintf("%s:%d", e.config.Host, e.config.Port)

	auth := smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.Host)

	if err := smtp.SendMail(addr, auth, e.config.From, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("sending reset email: %w", err)
	}

	log.Info().Str("to", to).Msg("password reset email sent")
	return nil
}
