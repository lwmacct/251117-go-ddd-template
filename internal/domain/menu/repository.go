package menu

import (
	"context"
)

// Repository 菜单仓储接口
type Repository interface {
	// Create 创建菜单
	Create(ctx context.Context, menu *Menu) error

	// Update 更新菜单
	Update(ctx context.Context, menu *Menu) error

	// Delete 删除菜单
	Delete(ctx context.Context, id uint) error

	// FindByID 根据ID查找菜单
	FindByID(ctx context.Context, id uint) (*Menu, error)

	// FindAll 查找所有菜单（树形结构）
	FindAll(ctx context.Context) ([]*Menu, error)

	// FindByParentID 根据父ID查找子菜单
	FindByParentID(ctx context.Context, parentID *uint) ([]*Menu, error)

	// UpdateOrder 批量更新菜单排序
	UpdateOrder(ctx context.Context, menus []struct {
		ID       uint
		Order    int
		ParentID *uint
	}) error
}
