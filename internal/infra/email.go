package infra

import (
	"fmt"
	"net/smtp"
	"rest-fiber/config"
	"strings"
)

type EmailReciever struct {
	Email string
}

type EmailParams struct {
	Subject  string
	Message  string
	Reciever EmailReciever
}

type EmailService interface {
	SendEmail(params EmailParams) error
}

type emailServiceImpl struct {
	env    config.Env
	logger *AppLogger 
}

func NewEmailService(env config.Env, logger *AppLogger ) EmailService {
	return &emailServiceImpl{env, logger}
}
func (s *emailServiceImpl) SendEmail(params EmailParams) error {
	addr := fmt.Sprintf("%s:%s", s.env.SMTPHost, s.env.SMTPPort)
	auth := smtp.PlainAuth(
		"",
		s.env.SMTPEmail,
		s.env.SMTPPassword,
		s.env.SMTPHost,
	)
	headers := []string{
		"From: " + s.env.SMTPEmail,
		"To: " + params.Reciever.Email,
		"Subject: " + params.Subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"utf-8\"",
		"",
		params.Message,
	}

	msg := []byte(strings.Join(headers, "\r\n"))

	return smtp.SendMail(addr, auth, s.env.SMTPEmail, []string{params.Reciever.Email}, msg)
}
