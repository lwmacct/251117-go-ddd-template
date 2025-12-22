package menu

import "context"

// QueryRepository 定义菜单读操作接口
type QueryRepository interface {
	// FindByID 根据ID查找菜单
	FindByID(ctx context.Context, id uint) (*Menu, error)

	// FindAll 查找所有菜单（树形结构）
	FindAll(ctx context.Context) ([]*Menu, error)

	// FindByParentID 根据父ID查找子菜单
	FindByParentID(ctx context.Context, parentID *uint) ([]*Menu, error)
}
