package ws

import (
	"encoding/json"
	"github.com/filedrive-team/filfind/backend/types"
)

type Channel string

const (
	MessageError Channel = "message_error"
	// unsubscribe message
	Login     Channel = "login"
	Logout    Channel = "logout"
	ChatSend  Channel = "chat_send"
	ChatCheck Channel = "chat_check"
	// subscribe message
	Echo               Channel = "echo"
	ChatReceive        Channel = "chat_receive"
	ChatPartners       Channel = "chat_partners"
	ChatHistory        Channel = "chat_history"
	ChatPartnersStatus Channel = "chat_partners_status"
	ChatCheckedStatus  Channel = "chat_checked_status"
	ChatUncheckedList  Channel = "chat_unchecked_list"
	ChatUncheckedTotal Channel = "chat_unchecked_total"

	MaxWsConnsOneUser = 20 // max websocket connections one user
)

type PackageReq struct {
	Channel   Channel         `json:"channel"`
	Subscribe *bool           `json:"subscribe,omitempty"`
	Body      json.RawMessage `json:"body,omitempty"`
}

type PackageResp struct {
	Channel Channel     `json:"channel"`
	Body    interface{} `json:"body,omitempty"`
}

func encodePackageResp(channel Channel, body interface{}) ([]byte, error) {
	pkg := &PackageResp{
		Channel: channel,
		Body:    body,
	}
	return json.Marshal(pkg)
}

const (
	ErrExpired        = "Login expired, please login again."
	ErrForbidden      = "Forbidden."
	ErrInternalServer = "Internal server error."
	ErrUserNotFound   = "User doesn't exist."
	ErrMessageEmpty   = "Empty message."
)

type ErrorMessage struct {
	Error  string `json:"error"`
	Source string `json:"source"`
}

func newNotSupportMessageError(src []byte) *ErrorMessage {
	return newMessageError("not support this message", src)
}

func newMessageError(error string, src []byte) *ErrorMessage {
	return &ErrorMessage{
		Error:  error,
		Source: string(src),
	}
}

type WrapPackageReq struct {
	PackageReq
	Sender *Client
	Raw    []byte
}

// LoginReq is user login request
type LoginReq struct {
	Token string `json:"token"`
}

// LoginResp is user login response
type LoginResp struct {
	Msg string `json:"msg"`
}

// Message is chat message
type Message struct {
	Sender    string          `json:"sender,omitempty"`
	Recipient string          `json:"recipient"`
	Type      int8            `json:"type"` // 0:text message
	Content   string          `json:"content"`
	Timestamp *types.UnixTime `json:"timestamp,omitempty"`
	Checked   bool            `json:"checked"`
}

// HistoryMessageReq is history chat message request
type HistoryMessageReq struct {
	Partner string          `json:"partner" form:"partner" binding:"required"`
	Before  *types.UnixTime `json:"before,omitempty" form:"before"`
	Limit   int             `json:"limit" form:"limit" binding:"required,min=1,max=100" example:"20"`
}

// CheckMessageReq is to check chat message request
type CheckMessageReq struct {
	Partner string `json:"partner"`
}

// CheckedMessage is notify to client update checked status of partner's chat message
type CheckedMessage struct {
	Partner string `json:"partner"`
}

// TotalUncheckedMessage is total unchecked chat messages
type TotalUncheckedMessage struct {
	Total int64 `json:"total"`
}
