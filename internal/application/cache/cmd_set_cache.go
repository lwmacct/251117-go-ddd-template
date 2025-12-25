package cache

import "time"

// SetCacheCommand 设置缓存命令
type SetCacheCommand struct {
	Key   string
	Value any
	TTL   time.Duration
}
