package setting

import "context"

// CommandRepository 定义配置写操作接口
type CommandRepository interface {
	// Create 创建配置
	Create(ctx context.Context, s *Setting) error

	// Update 更新配置
	Update(ctx context.Context, s *Setting) error

	// Delete 删除配置
	Delete(ctx context.Context, id uint) error

	// BatchUpsert 批量插入或更新配置
	BatchUpsert(ctx context.Context, settings []*Setting) error
}
