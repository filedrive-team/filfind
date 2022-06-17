package repo

import (
	"encoding/json"
	"github.com/filedrive-team/filfind/backend/kvcache/vcodemgr"
	"github.com/filedrive-team/filfind/backend/utils/utils"
	logger "github.com/sirupsen/logrus"
	"time"
)

// the vcode used for the development
const DevVcode = "123456"

// Verification code validity period
const VcodeDuration = 30 * time.Minute

type Vcode struct {
	Code       string `json:"code"`        // verification code
	ExpireTime int64  `json:"expire_time"` // expire time
}

func IsExpireVcode(vcode *Vcode) bool {
	return vcode.ExpireTime < time.Now().Unix()
}

func GenerateEmailVcodeToForgetPwd(email string, debug bool) (vcode *Vcode) {
	vcode = &Vcode{
		Code:       utils.GenerateRandNumStr(6),
		ExpireTime: time.Now().Add(VcodeDuration).Unix(),
	}
	if debug {
		vcode.Code = DevVcode
	}
	obj, err := json.Marshal(vcode)
	if err != nil {
		logger.WithError(err).Fatal("json marshal failed")
	}
	vcodemgr.SetEmailVcodeToForgetPwd(email, string(obj))
	return
}

func GetEmailVcodeToForgetPwd(email string) (vcode *Vcode, exist bool) {
	vcode = &Vcode{}
	obj, exist := vcodemgr.GetEmailVcodeToForgetPwd(email)
	if exist {
		err := json.Unmarshal([]byte(obj), vcode)
		if err != nil {
			logger.WithError(err).Error("json unmarshal failed, obj=", obj)
			exist = false
			return
		}
	}
	return
}

// Check whether a verification code has been generated within one minute
func IsFrequentEmailVcodeToForgetPwd(email string) bool {
	vcode := &Vcode{}
	obj, exist := vcodemgr.GetEmailVcodeToForgetPwd(email)
	if exist {
		err := json.Unmarshal([]byte(obj), vcode)
		if err != nil {
			logger.WithError(err).Error("json unmarshal failed, obj=", obj)
			return false
		}
		if time.Now().Unix()+int64(VcodeDuration/time.Second) < vcode.ExpireTime+int64(time.Minute/time.Second) {
			return true
		}
	}
	return false
}

func DeleteEmailVcodeToForgetPwd(email string) {
	vcodemgr.DeleteEmailVcodeToForgetPwd(email)
}
