package setting

import (
	"context"
)

// SchemaCacheService Schema 响应缓存服务接口。
//
// 缓存 Setting Schema API 的最终响应 DTO，避免重复的数据库查询和构建逻辑。
// 采用 user+category 维度缓存，支持懒加载场景。
//
// # 架构说明
//
// 本接口定义在 Application 层而非 domain/cache，因为：
//   - 缓存的是 Application 层 DTO（[SchemaCategoryDTO]）
//   - 遵循 DDD 依赖方向：Application 定义接口，Infrastructure 实现
//   - Domain 层不应知道 Application 层的 DTO 结构
//
// # Key 命名规范
//
//   - 用户 Schema：{prefix}schema:user:{userID}:{categoryKey}
//   - 管理员 Schema：{prefix}schema:admin:{categoryKey}
//   - categoryKey 为空时使用 "_all" 表示全量
//
// 默认 TTL：30 分钟
//
// # 缓存失效策略
//
//   - 用户修改自己的设置 → [DeleteUserSchemaAll]
//   - 管理员修改系统设置 → [DeleteByCategoryKey]
//   - Category 结构变更 → [DeleteAll]
//
// 实现位于 [infrastructure/redis.schemaCacheService]。
type SchemaCacheService interface {
	// =========================================================================
	// 用户 Schema 操作
	// =========================================================================

	// GetUserSchema 获取用户 Schema 缓存。
	//
	// userID: 用户 ID
	// categoryKey: 分类键，为空表示全量 Schema
	//
	// 缓存未命中返回 nil, nil。
	// 缓存数据损坏时自动清除并返回 nil, nil。
	GetUserSchema(ctx context.Context, userID uint, categoryKey string) ([]SchemaCategoryDTO, error)

	// SetUserSchema 设置用户 Schema 缓存。
	//
	// schema 为空时也会缓存（防止缓存穿透）。
	SetUserSchema(ctx context.Context, userID uint, categoryKey string, schema []SchemaCategoryDTO) error

	// DeleteUserSchema 删除用户的指定 category Schema 缓存。
	DeleteUserSchema(ctx context.Context, userID uint, categoryKey string) error

	// DeleteUserSchemaAll 删除用户的所有 Schema 缓存。
	//
	// 使用 SCAN 遍历 {prefix}schema:user:{userID}:* 模式。
	// 用于用户修改设置后失效所有相关缓存。
	DeleteUserSchemaAll(ctx context.Context, userID uint) error

	// =========================================================================
	// 管理员 Schema 操作
	// =========================================================================

	// GetAdminSchema 获取管理员 Schema 缓存。
	//
	// categoryKey: 分类键，为空表示全量 Schema
	//
	// 缓存未命中返回 nil, nil。
	GetAdminSchema(ctx context.Context, categoryKey string) ([]SchemaCategoryDTO, error)

	// SetAdminSchema 设置管理员 Schema 缓存。
	SetAdminSchema(ctx context.Context, categoryKey string, schema []SchemaCategoryDTO) error

	// DeleteAdminSchema 删除管理员的指定 category Schema 缓存。
	DeleteAdminSchema(ctx context.Context, categoryKey string) error

	// DeleteAdminSchemaAll 删除管理员的所有 Schema 缓存。
	//
	// 使用 SCAN 遍历 {prefix}schema:admin:* 模式。
	DeleteAdminSchemaAll(ctx context.Context) error

	// =========================================================================
	// 批量失效操作
	// =========================================================================

	// DeleteByCategoryKey 删除所有用户和管理员的指定 category Schema 缓存。
	//
	// 当系统设置定义变更时调用，使所有相关 Schema 缓存失效。
	// 使用 SCAN 遍历 {prefix}schema:*:{categoryKey} 模式。
	//
	// 注意：这是低频操作（仅管理员修改系统设置时触发），
	// SCAN 遍历的性能影响可接受。
	DeleteByCategoryKey(ctx context.Context, categoryKey string) error

	// DeleteAll 删除所有 Schema 缓存。
	//
	// 当 Category 结构变更或执行数据迁移时调用。
	// 使用 SCAN 遍历 {prefix}schema:* 模式。
	DeleteAll(ctx context.Context) error
}
