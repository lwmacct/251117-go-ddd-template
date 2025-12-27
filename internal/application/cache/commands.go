package cache

import "time"

// SetCommand 设置缓存命令
type SetCommand struct {
	Key   string
	Value any
	TTL   time.Duration
}

// DeleteCommand 删除缓存命令
type DeleteCommand struct {
	Key string
}
