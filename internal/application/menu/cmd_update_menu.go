package menu

// UpdateMenuCommand 更新菜单命令
type UpdateMenuCommand struct {
	MenuID   uint
	Title    *string
	Path     *string
	Icon     *string
	ParentID *uint
	Order    *int
	Visible  *bool
}
