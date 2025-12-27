package menu

// CreateCommand 创建菜单命令
type CreateCommand struct {
	Title    string
	Path     string
	Icon     string
	ParentID *uint
	Order    int
	Visible  bool
}

// UpdateCommand 更新菜单命令
type UpdateCommand struct {
	MenuID   uint
	Title    *string
	Path     *string
	Icon     *string
	ParentID *uint
	Order    *int
	Visible  *bool
}

// DeleteCommand 删除菜单命令
type DeleteCommand struct {
	MenuID uint
}

// MenuItemCommand 菜单排序项
type MenuItemCommand struct {
	ID       uint
	Order    int
	ParentID *uint
}

// ReorderCommand 批量更新菜单排序命令
type ReorderCommand struct {
	Menus []MenuItemCommand
}
