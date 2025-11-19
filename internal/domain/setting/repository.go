package setting

import "context"

// Repository 定义 Setting 仓储接口
type Repository interface {
	// Create 创建配置
	Create(ctx context.Context, s *Setting) error

	// Update 更新配置
	Update(ctx context.Context, s *Setting) error

	// Delete 删除配置
	Delete(ctx context.Context, id uint) error

	// FindByID 根据 ID 查找配置
	FindByID(ctx context.Context, id uint) (*Setting, error)

	// FindByKey 根据 Key 查找配置
	FindByKey(ctx context.Context, key string) (*Setting, error)

	// FindByCategory 根据分类查找配置列表
	FindByCategory(ctx context.Context, category string) ([]*Setting, error)

	// FindAll 查找所有配置
	FindAll(ctx context.Context) ([]*Setting, error)

	// BatchUpsert 批量插入或更新配置
	BatchUpsert(ctx context.Context, settings []*Setting) error
}
