package smtp

import (
	"github.com/filedrive-team/filfind/backend/settings"
	"github.com/filedrive-team/filfind/backend/utils/smtp"
	logger "github.com/sirupsen/logrus"
	"strings"
)

const (
	ResetPwdVCodeSubject = "Reset Password Verification"

	FromAlias = "FilFind"
)

var mailClients []smtp.MailClient
var officialWebsite string
var officialEmail string

func Setup(conf *settings.AppConfig) {
	officialWebsite = conf.App.OfficialWebsite
	officialEmail = conf.App.OfficialEmail
	basic := conf.Smtp.Basic
	if len(basic.Host) > 0 && len(basic.User) > 0 && len(basic.Password) > 0 {
		client := &smtp.BasicMailClient{
			Name:      "basic",
			User:      basic.User,
			Password:  basic.Password,
			Host:      basic.Host,
			UserAlias: FromAlias,
		}
		logger.Infof("%s mail client add to mailClients", client.GetName())

		mailClients = append(mailClients, client)
	}
}

func SendVCodeMail(subject string, vcode string, toAddress string) error {
	for k, v := range mailClients {
		if err := v.Send(subject, genVCodeEmailContent(vcode), toAddress); err != nil {
			logger.Errorf("%s mail send vcode fail: %s", v.GetName(), err.Error())
			if k == len(mailClients)-1 {
				logger.Error("all mail send vcode fail")
				return err
			}
			continue
		} else {
			logger.Infof("%s mail send vcode success", v.GetName())
			return nil
		}
	}

	return nil
}

func genVCodeEmailContent(vcode string) string {
	template := `
<div style="margin:0px auto;max-width:590px;">
  <p>Your verification code:<br />
  <p> <span style="font-size:18px;color:rgb(105 177 235);">$vcode</span><br />
  <p>The verification code will be valid for 30 minutes.
    Please do not share this code with anyone.&nbsp;
    Don’t recognize this activity? Please ignore this.&nbsp;<br />
  <p>FilFind Official Website：<a href=3D"https://filfind.io"><strong>$officialWebsite</strong></a><br />
  <p>Connect Email：$officialEmail<br />
</div>
`
	str := strings.ReplaceAll(template, "$officialWebsite", officialWebsite)
	str = strings.ReplaceAll(str, "$officialEmail", officialEmail)
	return strings.ReplaceAll(str, "$vcode", vcode)
}
