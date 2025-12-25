package menu

// CreateMenuCommand 创建菜单命令
type CreateMenuCommand struct {
	Title    string
	Path     string
	Icon     string
	ParentID *uint
	Order    int
	Visible  bool
}
