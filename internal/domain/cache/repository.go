// Package cache 定义缓存领域接口。
//
// 本包定义了缓存操作的领域接口，遵循 CQRS 模式：
//   - [CommandRepository]: 写仓储接口（Set, Delete, SetNX）
//   - [QueryRepository]: 读仓储接口（Get, Exists）
//
// 实现位于 infrastructure/redis 包。
package cache

import (
	"context"
	"time"
)

// CommandRepository 缓存写操作仓储接口
type CommandRepository interface {
	// Set 设置缓存值
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	// Delete 删除缓存值
	Delete(ctx context.Context, key string) error
	// SetNX 仅当键不存在时设置值 (分布式锁)
	SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error)
}

// QueryRepository 缓存读操作仓储接口
type QueryRepository interface {
	// Get 获取缓存值
	Get(ctx context.Context, key string, dest any) error
	// Exists 检查键是否存在
	Exists(ctx context.Context, key string) (bool, error)
}
