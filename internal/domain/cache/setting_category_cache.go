package cache

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// SettingCategoryCacheService SettingCategory 缓存服务接口。
//
// 提供按 ID 和 Key 粒度的缓存操作，采用 Cache-Aside 模式：
//   - 读取时先查缓存，未命中再查数据库
//   - 写入后自动失效相关缓存
//
// Key 命名规范：
//   - 按 ID：{prefix}setting_category:id:{id}
//   - 按 Key：{prefix}setting_category:key:{key}
//   - 全量：{prefix}setting_category:all
//
// 默认 TTL：7 天（数据变更不频繁，且有主动失效机制）
//
// 实现位于 [infrastructure/cache.settingCategoryCacheService]。
type SettingCategoryCacheService interface {
	// =========================================================================
	// 单条操作
	// =========================================================================

	// GetByID 根据 ID 获取缓存的 SettingCategory。
	// 缓存未命中返回 nil, nil（不返回错误）。
	// 缓存数据损坏时自动清除并返回 nil, nil。
	GetByID(ctx context.Context, id uint) (*setting.SettingCategory, error)

	// GetByKey 根据 Key 获取缓存的 SettingCategory。
	// 缓存未命中返回 nil, nil（不返回错误）。
	// 缓存数据损坏时自动清除并返回 nil, nil。
	GetByKey(ctx context.Context, key string) (*setting.SettingCategory, error)

	// Set 设置 SettingCategory 缓存。
	// 同时写入 ID 索引和 Key 索引。
	// 使用默认 TTL（7 天）。
	Set(ctx context.Context, category *setting.SettingCategory) error

	// =========================================================================
	// 批量操作
	// =========================================================================

	// GetAll 获取所有缓存的 SettingCategory。
	//
	// 从全量缓存 key 读取，返回按 Order 升序排列的列表。
	// 缓存为空时返回 nil，不返回错误。
	GetAll(ctx context.Context) ([]*setting.SettingCategory, error)

	// SetAll 批量设置 SettingCategory 缓存（用于预热）。
	//
	// 使用 Pipeline 批量写入：
	//   - 每个 Category 写入 ID 索引和 Key 索引
	//   - 写入全量缓存 key
	//
	// categories 为空时直接返回 nil。
	SetAll(ctx context.Context, categories []*setting.SettingCategory) error

	// GetByIDs 批量获取指定 ID 的缓存。
	//
	// 使用 JSON.MGET 一次网络往返获取多个 ID。
	// 返回找到的 Category 列表，未命中的 ID 不在结果中。
	GetByIDs(ctx context.Context, ids []uint) ([]*setting.SettingCategory, error)

	// =========================================================================
	// 删除操作
	// =========================================================================

	// Delete 删除指定的缓存。
	// 同时删除 ID 索引、Key 索引和全量缓存。
	Delete(ctx context.Context, id uint, key string) error

	// DeleteAll 删除所有 SettingCategory 缓存。
	// 使用 SCAN 命令遍历，适用于缓存重建场景。
	DeleteAll(ctx context.Context) error
}
