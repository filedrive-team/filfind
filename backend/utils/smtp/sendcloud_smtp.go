package smtp

import (
	"gopkg.in/gomail.v2"
	"strconv"
	"sync/atomic"
)

type SendCloudMailClient struct {
	Name               string
	mailVcodeFailCount uint64

	// Set your key information
	ApiUser  string
	ApiKey   string
	Host     string
	Port     string
	From     string
	FromName string
}

func (c *SendCloudMailClient) GetName() string {
	return c.Name
}

func (c *SendCloudMailClient) SendVCodeMailBySendCloundImplementation(mailTo []string, subject string, body string) error {
	port, _ := strconv.Atoi(c.Port)

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(c.From, c.FromName))
	m.SetHeader("To", mailTo...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(c.Host, port, c.ApiUser, c.ApiKey)

	err := d.DialAndSend(m)
	return err
}

func (c *SendCloudMailClient) Send(subject, content string, toAddress string) error {
	mailTo := []string{toAddress}

	err := c.SendVCodeMailBySendCloundImplementation(mailTo, subject, content)
	if err != nil {
		atomic.AddUint64(&c.mailVcodeFailCount, 1)
		failCount := atomic.LoadUint64(&c.mailVcodeFailCount)
		if failCount%100 == 0 {
			// TODO: send warning
		}
		return err
	}
	return nil
}
