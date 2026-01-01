package product

import "context"

// CommandRepository 产品写仓储接口。
type CommandRepository interface {
	// Create 创建产品，成功后回写 ID 到实体。
	Create(ctx context.Context, product *Product) error
	// Update 更新产品。
	Update(ctx context.Context, product *Product) error
	// Delete 删除产品。
	Delete(ctx context.Context, id uint) error
}

// QueryRepository 产品读仓储接口。
type QueryRepository interface {
	// GetByID 根据 ID 获取产品。
	GetByID(ctx context.Context, id uint) (*Product, error)
	// GetByCode 根据产品代码获取产品。
	GetByCode(ctx context.Context, code string) (*Product, error)
	// GetByName 根据名称获取产品。
	GetByName(ctx context.Context, name string) (*Product, error)
	// List 分页获取产品列表。
	List(ctx context.Context, offset, limit int) ([]*Product, error)
	// Count 获取产品总数。
	Count(ctx context.Context) (int64, error)
	// ListByStatus 按状态分页获取产品列表。
	ListByStatus(ctx context.Context, status string, offset, limit int) ([]*Product, error)
	// CountByStatus 按状态统计产品数量。
	CountByStatus(ctx context.Context, status string) (int64, error)
	// ExistsByName 检查产品名称是否存在。
	ExistsByName(ctx context.Context, name string) (bool, error)
	// ExistsByCode 检查产品代码是否存在。
	ExistsByCode(ctx context.Context, code string) (bool, error)
	// ListActive 获取所有激活的产品（用于产品目录）。
	ListActive(ctx context.Context) ([]*Product, error)
}
