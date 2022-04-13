package utils

import "github.com/VictoriaMetrics/fastcache"

var LocalCache *fastcache.Cache

func init() {
	LocalCacheInit()
}

func LocalCacheInit() {
	LocalCache = fastcache.New(1024)
}
