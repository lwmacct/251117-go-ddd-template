package setting

import "context"

// QueryRepository 定义配置读操作接口
type QueryRepository interface {
	// FindByID 根据 ID 查找配置
	FindByID(ctx context.Context, id uint) (*Setting, error)

	// FindByKey 根据 Key 查找配置
	FindByKey(ctx context.Context, key string) (*Setting, error)

	// FindByKeys 根据多个 Key 批量查找配置
	FindByKeys(ctx context.Context, keys []string) ([]*Setting, error)

	// FindByCategory 根据分类查找配置列表
	FindByCategory(ctx context.Context, category string) ([]*Setting, error)

	// FindAll 查找所有配置
	FindAll(ctx context.Context) ([]*Setting, error)
}
