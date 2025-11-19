package menu

import "context"

// CommandRepository 定义菜单写操作接口
type CommandRepository interface {
	// Create 创建菜单
	Create(ctx context.Context, menu *Menu) error

	// Update 更新菜单
	Update(ctx context.Context, menu *Menu) error

	// Delete 删除菜单
	Delete(ctx context.Context, id uint) error

	// UpdateOrder 批量更新菜单排序
	UpdateOrder(ctx context.Context, menus []struct {
		ID       uint
		Order    int
		ParentID *uint
	}) error
}
