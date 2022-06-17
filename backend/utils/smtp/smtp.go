package smtp

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"
)

type MailClient interface {
	GetName() string
	Send(subject, content string, toAddress string) error
}

type SmtpMail struct {
	Host      string
	User      string
	Password  string
	UserAlias string

	To          string
	Subject     string
	Content     string
	ContentType string
}

func SendMailBySmtp(smtpMail *SmtpMail) error {
	hp := strings.Split(smtpMail.Host, ":")
	if len(hp) != 2 {
		return errors.New("error host")
	}
	auth := smtp.PlainAuth("", smtpMail.User, smtpMail.Password, hp[0])
	from := smtpMail.User
	if smtpMail.UserAlias != "" {
		from = fmt.Sprintf("%s<%s>", smtpMail.UserAlias, smtpMail.User)
	}
	msg := []byte(fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n%s\r\n\r\n%s", smtpMail.To, from,
		smtpMail.Subject, smtpMail.ContentType, smtpMail.Content))
	err := smtp.SendMail(smtpMail.Host, auth, smtpMail.User, []string{smtpMail.To}, msg)
	return err
}
