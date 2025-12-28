package cache

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
)

// UserSettingQueryCacheService 用户设置查询缓存接口（Domain 层）。
//
// 存储原始 UserSetting 记录，用于 Repository 层减少数据库查询。
// 采用用户维度全量缓存策略：一次查询缓存用户所有自定义配置。
//
// 与 [UserSettingCacheService] 的区别：
//   - 本接口：存储原始 UserSetting（Domain 实体），Repository 层使用
//   - UserSettingCacheService：存储合并后的 EffectiveUserSetting（DTO），Application 层使用
//
// Key 格式：{prefix}usersetting:user:{userID}
// 默认 TTL：30 分钟
type UserSettingQueryCacheService interface {
	// GetByUser 获取用户的所有自定义配置缓存。
	// 返回 map[settingKey]*UserSetting，未命中返回 nil, nil。
	// 缓存数据损坏时自动清除并返回 nil, nil。
	GetByUser(ctx context.Context, userID uint) (map[string]*setting.UserSetting, error)

	// SetByUser 设置用户的所有自定义配置缓存。
	// settings 为空时缓存空结果（防止缓存穿透）。
	SetByUser(ctx context.Context, userID uint, settings []*setting.UserSetting) error

	// DeleteByUser 删除用户的所有配置缓存。
	DeleteByUser(ctx context.Context, userID uint) error
}
