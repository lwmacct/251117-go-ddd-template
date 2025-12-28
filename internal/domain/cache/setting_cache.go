package cache

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// SettingCacheService Setting 缓存服务接口。
//
// 提供按 key 粒度的缓存操作，采用 Cache-Aside 模式：
//   - 读取时先查缓存，未命中再查数据库
//   - 写入后自动失效相关缓存
//
// Key 命名规范：{prefix}setting:{key}
// 默认 TTL：10 分钟
//
// 实现位于 [infrastructure/redis.settingCacheService]。
type SettingCacheService interface {
	// Get 获取缓存的 Setting。
	// 缓存未命中返回 nil, nil（不返回错误）。
	// 缓存数据损坏时自动清除并返回 nil, nil。
	Get(ctx context.Context, key string) (*setting.Setting, error)

	// Set 设置 Setting 缓存。
	// 使用默认 TTL（10 分钟）。
	Set(ctx context.Context, s *setting.Setting) error

	// Delete 删除指定 key 的缓存。
	Delete(ctx context.Context, key string) error

	// DeleteByKeys 批量删除缓存。
	// keys 为空时直接返回 nil。
	DeleteByKeys(ctx context.Context, keys []string) error

	// DeleteAll 删除所有 Setting 缓存。
	// 使用 SCAN 命令遍历，适用于缓存重建场景。
	DeleteAll(ctx context.Context) error
}
