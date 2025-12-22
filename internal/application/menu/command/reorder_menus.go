package command

// MenuItem 菜单排序项
type MenuItem struct {
	ID       uint
	Order    int
	ParentID *uint
}

// ReorderMenusCommand 批量更新菜单排序命令
type ReorderMenusCommand struct {
	Menus []MenuItem
}
