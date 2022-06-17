package vcodemgr

import (
	"github.com/filedrive-team/filfind/backend/kvcache"
	"time"
)

const (
	emailVcodeRegisterPrefix  = "register_"
	emailVcodeForgetPwdPrefix = "forget_pwd_"
)

func GetEmailVcodeToRegister(email string) (vcode string, b bool) {
	return kvcache.GetObject(emailVcodeRegisterPrefix + email)
}

func SetEmailVcodeToRegister(email string, vcode string) {
	kvcache.SetObject(emailVcodeRegisterPrefix+email, vcode, time.Hour)
}

func DeleteEmailVcodeToRegister(email string) {
	kvcache.DeleteObject(emailVcodeRegisterPrefix + email)
}

func GetEmailVcodeToForgetPwd(email string) (vcode string, b bool) {
	return kvcache.GetObject(emailVcodeForgetPwdPrefix + email)
}

func SetEmailVcodeToForgetPwd(email string, vcode string) {
	kvcache.SetObject(emailVcodeForgetPwdPrefix+email, vcode, time.Hour)
}

func DeleteEmailVcodeToForgetPwd(email string) {
	kvcache.DeleteObject(emailVcodeForgetPwdPrefix + email)
}
