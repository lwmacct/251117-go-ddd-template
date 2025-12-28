package setting

import "context"

// QueryRepository 配置定义读操作接口。
type QueryRepository interface {
	// FindByKey 根据 Key 查找配置定义
	FindByKey(ctx context.Context, key string) (*Setting, error)

	// FindByKeys 根据多个 Key 批量查找配置定义
	FindByKeys(ctx context.Context, keys []string) ([]*Setting, error)

	// FindByCategoryID 根据分类 ID 查找配置定义列表
	FindByCategoryID(ctx context.Context, categoryID uint) ([]*Setting, error)

	// FindAll 查找所有配置定义
	FindAll(ctx context.Context) ([]*Setting, error)

	// ExistsByKey 检查 Key 是否已存在
	ExistsByKey(ctx context.Context, key string) (bool, error)
}
