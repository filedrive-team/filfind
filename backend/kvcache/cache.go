package kvcache

import (
	"time"
)

func GetObject(key string) (o string, b bool) {
	res, b := globalCache.Get(key)
	if b {
		o = res.(string)
	}
	return
}

func SetObject(key string, object string, duration time.Duration) {
	globalCache.Set(key, object, duration)
}

func DeleteObject(key string) {
	globalCache.Delete(key)
}

func SetInt32(key string, value int32, duration time.Duration) {
	globalCache.Set(key, value, duration)
}

func GetInt32(key string) (v int32, b bool) {
	res, b := globalCache.Get(key)
	if b {
		v = res.(int32)
	}
	return
}

func IncrInt32(key string, n int32) (int32, error) {
	return globalCache.IncrementInt32(key, n)
}

func DecrInt32(key string, n int32) (int32, error) {
	return globalCache.DecrementInt32(key, n)
}
