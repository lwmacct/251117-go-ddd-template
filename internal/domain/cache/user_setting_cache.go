package cache

import (
	"context"
)

// UserSettingCacheService 用户设置缓存服务接口。
//
// 存储用户的有效设置值（合并后的结果），而非原始 UserSetting 记录。
// 这样查询时无需再合并 Setting.DefaultValue 和 UserSetting.Value。
//
// Key 格式：{prefix}user:{userID}:setting:{key}
// 默认 TTL：30 分钟（比系统设置缓存更长，因为用户设置变更频率较低）
//
// 多实例安全：
//   - 删除操作直接生效，无需跨实例通知（无本地缓存）
//
// 实现位于 [infrastructure/redis.userSettingCacheService]。
type UserSettingCacheService interface {
	// =========================================================================
	// 单条操作
	// =========================================================================

	// Get 获取用户的有效设置值。
	//
	// userID: 用户 ID
	// key: 设置 key（如 "general.theme"）
	//
	// 缓存未命中返回 nil, nil。
	// 缓存数据损坏时自动清除并返回 nil, nil。
	Get(ctx context.Context, userID uint, key string) (*EffectiveUserSetting, error)

	// Set 缓存用户的有效设置值。
	Set(ctx context.Context, userID uint, value *EffectiveUserSetting) error

	// =========================================================================
	// 批量操作
	// =========================================================================

	// GetByKeys 批量获取用户的有效设置值。
	//
	// 使用 MGET 一次网络往返获取多个 key。
	// 返回 key -> *EffectiveUserSetting 映射，未命中的 key 不在结果中。
	GetByKeys(ctx context.Context, userID uint, keys []string) (map[string]*EffectiveUserSetting, error)

	// SetBatch 批量设置用户的有效设置值。
	//
	// 使用 Pipeline 批量写入。
	// values 为空时直接返回 nil。
	SetBatch(ctx context.Context, userID uint, values []*EffectiveUserSetting) error

	// =========================================================================
	// 删除操作
	// =========================================================================

	// Delete 删除用户的指定设置缓存。
	Delete(ctx context.Context, userID uint, key string) error

	// DeleteByKeys 批量删除用户的指定设置缓存。
	DeleteByKeys(ctx context.Context, userID uint, keys []string) error

	// DeleteByUser 删除用户的所有设置缓存。
	//
	// 使用 SCAN 遍历 {prefix}user:{userID}:setting:* 模式。
	// 用于用户重置所有设置或用户删除场景。
	DeleteByUser(ctx context.Context, userID uint) error

	// DeleteBySettingKey 删除所有用户的某个设置缓存。
	//
	// 当系统默认值变更时调用，使所有用户的该设置缓存失效。
	// 使用 SCAN 遍历 {prefix}user:*:setting:{key} 模式。
	//
	// 注意：这是低频操作（仅管理员修改系统设置时触发），
	// SCAN 遍历的性能影响可接受。
	DeleteBySettingKey(ctx context.Context, key string) error

	// DeleteBySettingKeys 批量删除所有用户的多个设置缓存。
	//
	// 当批量修改系统默认值时调用。
	DeleteBySettingKeys(ctx context.Context, keys []string) error
}

// EffectiveUserSetting 用户有效设置值（缓存数据结构）。
//
// 存储合并后的实际生效值，避免查询时再次合并 Setting + UserSetting。
// 包含 UI 渲染所需的元数据，前端可直接使用。
type EffectiveUserSetting struct {
	// Key 设置键
	Key string `json:"key"`

	// Value 实际生效值（用户值或系统默认值）
	Value any `json:"value"`

	// DefaultValue 系统默认值（用于判断是否自定义、重置功能）
	DefaultValue any `json:"default_value"`

	// IsCustomized 是否被用户自定义
	// true: Value 来自 UserSetting
	// false: Value 等于 DefaultValue
	IsCustomized bool `json:"is_customized"`

	// =========================================================================
	// UI 元数据（透传给前端）
	// =========================================================================

	// ValueType 值类型：string, number, boolean, json
	ValueType string `json:"value_type"`

	// CategoryID 所属分类 ID
	CategoryID uint `json:"category_id"`

	// Group 分类内子分组
	Group string `json:"group"`

	// Label 显示标签
	Label string `json:"label"`

	// UIConfig UI 配置（JSON 字符串）
	UIConfig string `json:"ui_config,omitempty"`

	// Order 排序权重
	Order int `json:"order"`
}
