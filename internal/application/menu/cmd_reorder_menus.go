package menu

// MenuItemCommand 菜单排序项
type MenuItemCommand struct {
	ID       uint
	Order    int
	ParentID *uint
}

// ReorderMenusCommand 批量更新菜单排序命令
type ReorderMenusCommand struct {
	Menus []MenuItemCommand
}
