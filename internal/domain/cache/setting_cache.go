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
//
// 默认 TTL：7 天（数据变更不频繁，且有主动失效机制）
//
// 实现位于 [infrastructure/cache.settingCacheService]。
type SettingCacheService interface {
	// =========================================================================
	// 单条操作
	// =========================================================================

	// Get 获取缓存的 Setting。
	// 缓存未命中返回 nil, nil（不返回错误）。
	// 缓存数据损坏时自动清除并返回 nil, nil。
	Get(ctx context.Context, key string) (*setting.Setting, error)

	// Set 设置 Setting 缓存。
	// 使用默认 TTL（7 天）。
	Set(ctx context.Context, s *setting.Setting) error

	// =========================================================================
	// 批量操作
	// =========================================================================

	// GetAll 获取所有缓存的 Setting。
	//
	// 使用 SCAN 遍历 {prefix}setting:* 模式的 key。
	// 返回 key -> *Setting 的映射，跳过预热标记等非数据 key。
	// 缓存为空时返回空 map，不返回错误。
	GetAll(ctx context.Context) (map[string]*setting.Setting, error)

	// SetAll 批量设置 Setting 缓存（用于预热）。
	//
	// 使用 Pipeline 批量写入，提高性能。
	// settings 为空时直接返回 nil。
	SetAll(ctx context.Context, settings []*setting.Setting) error

	// GetByKeys 批量获取指定 key 的缓存。
	//
	// 使用 MGET 一次网络往返获取多个 key。
	// 返回 key -> *Setting 映射，未命中的 key 不在结果中。
	GetByKeys(ctx context.Context, keys []string) (map[string]*setting.Setting, error)

	// =========================================================================
	// 删除操作
	// =========================================================================

	// Delete 删除指定 key 的缓存。
	Delete(ctx context.Context, key string) error

	// DeleteByKeys 批量删除缓存。
	// keys 为空时直接返回 nil。
	DeleteByKeys(ctx context.Context, keys []string) error

	// DeleteAll 删除所有 Setting 缓存。
	// 使用 SCAN 命令遍历，适用于缓存重建场景。
	DeleteAll(ctx context.Context) error
}
