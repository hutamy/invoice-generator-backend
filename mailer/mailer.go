package mailer

import (
	"fmt"
	"net/smtp"
)

type SMTPClient struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

func NewSMTPClient(host string, port int, username, password, sender string) *SMTPClient {
	return &SMTPClient{host, port, username, password, sender}
}

func (c *SMTPClient) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", c.Username, c.Password, c.Host)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\n%s",
		c.Sender, to, subject, body,
	))

	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	return smtp.SendMail(addr, auth, c.Sender, []string{to}, msg)
}
