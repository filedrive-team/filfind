package errormsg

import (
	"github.com/filedrive-team/filfind/backend/settings"
	"github.com/gin-gonic/gin"
)

type ErrorCode int

const (
	unknownError ErrorCode = iota

	// General error
	SearchFailed
	OperationFailed
	ParamsError
	InternalServerError
	NotFoundError
	// Permissions error
	Expire
	Forbidden
	ReviewForbidden
	// Account error
	Unregister
	HasRegistered
	RegisterError
	WrongPassword
	EmailFormatError
	EmailRegistered
	UserOrPasswordError
	PasswordFormatError
	NewPasswordFormatError
	UserTypeError
	AddressError
	SignatureError
	ActorNotFound
	// VCode error
	SendVCodeFailed
	VCodeExpire
	VCodeError
	VCodeNotExist
	TryAgainVCodeRequest
	TooManyVCodeRequest
	// File upload error
	SaveFileFailed
	UploadFailed
	// Query error
	PageError
	PageSizeError

	OwnProviderError
	ScoreError
	TextLengthError
	CustomTextLengthError
	CustomFieldRequiredError
	FileTypeError
)
const (
	ZH string = "zh"
	EN string = "en"
)

type MultiLang struct {
	ZH string `json:"zh"`
	EN string `json:"en"`
}

// ByCtx - give appropriate message according to custom http header "X-Language"
func ByCtx(c *gin.Context, key ErrorCode) string {
	langcode := c.GetHeader(settings.HTTPLanguageHeader)
	return ByLangcode(langcode, key)
}

func ByLangcode(langcode string, key ErrorCode) string {
	msgData, ok := msgmap[key]
	if !ok {
		return ""
	}
	if langcode == ZH {
		return msgData.ZH
	}

	return msgData.EN
}

var msgmap = map[ErrorCode]MultiLang{
	SearchFailed:        {"查询失败", "Query failed."},
	OperationFailed:     {"操作失败", "Operation failed."},
	ParamsError:         {"参数错误", "Param format error."},
	InternalServerError: {"服务器内部错误", "Internal server error."},
	NotFoundError:       {"找不到记录", "Not found."},

	Expire:          {"登录状态已过期,请重新登录", "Login expired, please login again."},
	Forbidden:       {"拒绝访问", "Forbidden."},
	ReviewForbidden: {"对不起!\n您不能提交此评论，因为该存储供应商没有存储您的交易。", "Sorry!\nYou can't submit this review since this storage provider hasn't stored your deals."},

	Unregister:             {"账户未注册", "Account doesn't exist."},
	HasRegistered:          {"帐户已注册", "Account already signed up."},
	RegisterError:          {"注册失败,请重新尝试或联系客服", "Failed to sign up."},
	WrongPassword:          {"密码错误", "Wrong password."},
	EmailFormatError:       {"邮箱地址不正确", "The Email address is incorrectly formatted."},
	EmailRegistered:        {"邮箱已被注册", "Email has been registered."},
	UserOrPasswordError:    {"用户名或密码错误", "User or password is incorrect."},
	PasswordFormatError:    {"密码格式不正确", "The password format is incorrect."},
	NewPasswordFormatError: {"新密码格式不正确", "The new password format is incorrect."},
	UserTypeError:          {"用户类型错误", "User type error."},
	AddressError:           {"Filecoin地址格式错误", "Filecoin Address format error."},
	SignatureError:         {"签名或消息错误", "Signature error or the message is not canonical."},
	ActorNotFound:          {"地址没有找到", "Actor not found."},

	SendVCodeFailed:      {"发送验证码失败,请检查后重试", "Failed to send the code."},
	VCodeExpire:          {"验证码超时", "Code is invalid."},
	VCodeError:           {"验证码错误", "Code is incorrect."},
	VCodeNotExist:        {"验证码不存在", "Code doesn't exist."},
	TryAgainVCodeRequest: {"请稍后再试", "Try again later."},
	TooManyVCodeRequest:  {"达到验证码次数上限，请明日再试", "Try again later."},

	SaveFileFailed: {"保存文件失败", "Failed to save uploaded file"},
	UploadFailed:   {"上传失败", "Upload failed."},

	PageError:     {"Page不合法", "The Page is illegal."},
	PageSizeError: {"PageSize不合法", "The Page Size is illegal."},

	OwnProviderError:         {"你还没拥有这个存储供应商", "You have not own this storage provider."},
	ScoreError:               {"评分不合法", "The Score is illegal."},
	TextLengthError:          {"文本长度超限", "The length of text is out of limit."},
	CustomTextLengthError:    {"%s长度超限", "The length of %s is out of limit."},
	CustomFieldRequiredError: {"%s字段是必须的", "The %s field is required."},
	FileTypeError:            {"文件类型不合法", "The file type is illegal."},
}
