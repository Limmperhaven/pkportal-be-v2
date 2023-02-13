package client

import (
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/config"
	"net/smtp"
)

type MailClient struct {
	auth    smtp.Auth
	address string
	url     string
}

func NewMailClient(cfg *config.SMTP) *MailClient {
	return &MailClient{
		auth:    smtp.PlainAuth("", cfg.Address, cfg.Password, cfg.Host),
		address: cfg.Address,
		url:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	}
}

func (m *MailClient) SendTextEmail(subject, message string, to []string) error {
	return m.sendEmail(subject, message, to, "test/plain")
}

func (m *MailClient) SendHTMLEmail(subject, message string, to []string) error {
	return m.sendEmail(subject, message, to, "test/html")
}

func (m *MailClient) sendEmail(subject, message string, to []string, contentType string) error {
	msg := fmt.Sprintf(
		"From: %s\r\n"+
			"Subject: %s\r\n\r\n"+
			"%s\r\n", m.address, subject, message)
	return smtp.SendMail(m.url, m.auth, m.address, to, []byte(msg))
}
