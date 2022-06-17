package kvcache

import (
	"github.com/filedrive-team/filfind/backend/settings"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	globalCache *cache.Cache
)

func InitGlobalCache(conf *settings.AppConfig) {
	globalCache = cache.New(5*time.Minute, 10*time.Minute)
}
