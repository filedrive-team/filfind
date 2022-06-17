package smtp

import (
	"encoding/base64"
	"strings"
	"sync/atomic"
)

type BasicMailClient struct {
	Name          string
	mailFailCount uint64

	// set your account
	User      string
	Password  string
	Host      string
	UserAlias string
}

func (c *BasicMailClient) GetName() string {
	return c.Name
}

func (c *BasicMailClient) Send(subject, content string, toAddress string) error {
	var builder strings.Builder
	builder.WriteString("=?UTF-8?B?")
	builder.WriteString(base64.StdEncoding.EncodeToString([]byte(subject)))
	builder.WriteString("?=")

	smtpMail := &SmtpMail{
		User:      c.User,
		Password:  c.Password,
		Host:      c.Host,
		UserAlias: c.UserAlias,
	}

	smtpMail.To = toAddress
	smtpMail.Subject = builder.String()
	smtpMail.Content = content
	smtpMail.ContentType = "Content-Type: text/html; charset=UTF-8"
	err := SendMailBySmtp(smtpMail)
	if err != nil {
		failCount := atomic.AddUint64(&c.mailFailCount, 1)
		if failCount%100 == 0 {
			// TODO: send warning
		}
		return err
	}
	return nil
}
