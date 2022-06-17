package vcodelimit

import (
	"github.com/filedrive-team/filfind/backend/kvcache"
	"sync"
	"time"
)

const (
	vcodeCountEmailPrefix = "email_count_"
	vcodeCountIPPrefix    = "ip_count_"

	VcodeCountLimitPerEmailPer24H = 10
	VcodeCountLimitPerIPPer24H    = 100
)

var lock sync.Mutex

func IsOutOfEmailLimit(email string) bool {
	count, exist := kvcache.GetInt32(vcodeCountEmailPrefix + email)
	if exist {
		return count >= VcodeCountLimitPerEmailPer24H
	}
	return false
}

func IsOutOfIPLimit(ip string) bool {
	count, exist := kvcache.GetInt32(vcodeCountIPPrefix + ip)
	if exist {
		return count >= VcodeCountLimitPerIPPer24H
	}
	return false
}

func IncrCount(email string, ip string) (ok bool) {
	lock.Lock()
	defer lock.Unlock()

	emailKey := vcodeCountEmailPrefix + email
	ipKey := vcodeCountIPPrefix + ip
	count, exist := kvcache.GetInt32(emailKey)
	if exist {
		if count >= VcodeCountLimitPerEmailPer24H {
			return false
		}
		_, err := kvcache.IncrInt32(emailKey, 1)
		if err != nil {
			return false
		}
	} else {
		kvcache.SetInt32(emailKey, 1, 24*time.Hour)
	}

	count, exist = kvcache.GetInt32(ipKey)
	if exist {
		if count >= VcodeCountLimitPerIPPer24H {
			kvcache.DecrInt32(emailKey, 1)
			return false
		}
		_, err := kvcache.IncrInt32(ipKey, 1)
		if err != nil {
			kvcache.DecrInt32(emailKey, 1)
			return false
		}
	} else {
		kvcache.SetInt32(ipKey, 1, 24*time.Hour)
	}
	return true
}
