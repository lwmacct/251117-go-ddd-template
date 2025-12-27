package setting

import "context"

// CommandRepository 配置定义写操作接口。
type CommandRepository interface {
	// Create 创建配置定义
	Create(ctx context.Context, setting *Setting) error

	// Update 更新配置定义
	Update(ctx context.Context, setting *Setting) error

	// Delete 删除配置定义
	Delete(ctx context.Context, key string) error

	// BatchUpsert 批量插入或更新配置定义
	BatchUpsert(ctx context.Context, settings []*Setting) error
}
