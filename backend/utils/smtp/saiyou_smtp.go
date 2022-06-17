package smtp

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
)

type SaiyouMailClient struct {
	Name               string
	mailVcodeFailCount uint64

	// Set your own key pair
	Appid    string
	Appkey   string
	SignType string
}

func (c *SaiyouMailClient) GetName() string {
	return c.Name
}

type SaiyouResp struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
}

func (c *SaiyouMailClient) Send(subject, content string, toAddress string) error {
	config := make(map[string]string)
	config["appid"] = c.Appid
	config["appkey"] = c.Appkey
	config["signType"] = c.SignType

	submail := CreateSaiyouMailSender(config)
	submail.AddTo(toAddress, toAddress)
	submail.SetSubject(subject)
	submail.SetHtml(content)
	submail.AddVar("name", "leo")
	submail.AddLink("url", "https://www.mysubmail.com")
	submail.AddHeaders("X-Mailer", "SUBMAIL Golang SDK")
	send := submail.Send()

	s := &SaiyouResp{}
	json.Unmarshal([]byte(send), s)

	if s.Status == "error" {
		e := errors.New("saiyou mail error, code=" + strconv.Itoa(s.Code) + ",  msg=" + s.Msg)
		failCount := atomic.AddUint64(&c.mailVcodeFailCount, 1)
		if failCount%100 == 0 {
			// TODO: send warning
		}
		return e

	}
	return nil
}

const sendURL = "https://api.mysubmail.com/mail/send"

func CreateSaiyouMailSender(config map[string]string) *SaiyouMailSend {
	mail := new(SaiyouMailSend)
	mail.appid = config["appid"]
	mail.appkey = config["appkey"]
	mail.signType = config["signType"]
	mail.vars = make(map[string]string)
	mail.links = make(map[string]string)
	mail.headers = make(map[string]string)
	return mail
}

type SaiyouMailSend struct {
	appid        string
	appkey       string
	signType     string
	to           []map[string]string
	from         string
	fromName     string
	reply        string
	cc           []string
	bcc          []string
	subject      string
	html         string
	text         string
	vars         map[string]string
	links        map[string]string
	headers      map[string]string
	asynchronous string
	attachments  []string
	tag          string
}

func (this *SaiyouMailSend) AddTo(address string, name string) {
	item := make(map[string]string)
	item["address"] = address
	item["name"] = name
	this.to = append(this.to, item)
}

func (this *SaiyouMailSend) SetSender(address string, name string) {
	this.from = address
	this.fromName = name
}

func (this *SaiyouMailSend) SetReply(address string) {
	this.reply = address
}

func (this *SaiyouMailSend) AddCc(address string) {
	this.cc = append(this.cc, address)
}

func (this *SaiyouMailSend) AddBcc(address string) {
	this.bcc = append(this.bcc, address)
}

func (this *SaiyouMailSend) SetSubject(subject string) {
	this.subject = subject
}

func (this *SaiyouMailSend) SetHtml(html string) {
	this.html = html
}

func (this *SaiyouMailSend) SetText(text string) {
	this.text = text
}

func (this *SaiyouMailSend) AddVar(key string, val string) {
	this.vars[key] = val
}

func (this *SaiyouMailSend) AddLink(key string, val string) {
	this.links[key] = val
}

func (this *SaiyouMailSend) AddHeaders(key string, val string) {
	this.headers[key] = val
}

func (this *SaiyouMailSend) SetAsynchronous(status bool) {
	if status {
		this.asynchronous = "true"
	} else {
		this.asynchronous = "false"
	}
}

func (this *SaiyouMailSend) AddAttachments(file string) {
	this.attachments = append(this.attachments, file)
}

func (this *SaiyouMailSend) SetTag(tag string) {
	this.tag = tag
}

func (this *SaiyouMailSend) Send() string {
	config := make(map[string]string)
	config["appid"] = this.appid
	config["appkey"] = this.appkey
	config["signType"] = this.signType

	request := make(map[string]string)
	request["appid"] = this.appid
	if len(this.to) > 0 {
		to_list := make([]string, 0, 32)
		for _, item := range this.to {
			to_list = append(to_list, fmt.Sprintf("%s<%s>", item["name"], item["address"]))
		}
		request["to"] = strings.Join(to_list, ",")
	}
	if this.from != "" {
		request["from"] = this.from
	}
	if this.fromName != "" {
		request["from_name"] = this.fromName
	}
	if this.reply != "" {
		request["reply"] = this.reply
	}
	if len(this.cc) > 0 {
		request["cc"] = strings.Join(this.cc, ",")
	}
	if len(this.bcc) > 0 {
		request["bcc"] = strings.Join(this.bcc, ",")
	}
	if this.subject != "" {
		request["subject"] = this.subject
	}
	if this.asynchronous != "" {
		request["asynchronous"] = this.asynchronous
	}
	if this.tag != "" {
		request["tag"] = this.tag
	}
	if this.signType != "normal" {
		request["sign_type"] = this.signType
		request["timestamp"] = GetTimestamp()
		request["sign_version"] = "2"
	}
	request["signature"] = CreateSignature(request, config)

	//V2 digital signature, does not participate in digital signature calculation about html/text/vars/links/headers/attachments fields
	if this.html != "" {
		request["html"] = this.html
	}
	if this.text != "" {
		request["text"] = this.text
	}

	if len(this.vars) > 0 {
		data, err := json.Marshal(this.vars)
		if err == nil {
			request["vars"] = string(data)
		}
	}

	if len(this.links) > 0 {
		data, err := json.Marshal(this.links)
		if err == nil {
			request["links"] = string(data)
		}
	}

	if len(this.headers) > 0 {
		data, err := json.Marshal(this.headers)
		if err == nil {
			request["headers"] = string(data)
		}
	}

	if len(this.attachments) > 0 {
		request["attachments"] = strings.Join(this.attachments, ",")
	}
	return MultipartPost(sendURL, request)
}
