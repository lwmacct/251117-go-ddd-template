// Package command 定义菜单命令处理器
package command

// CreateMenuCommand 创建菜单命令
type CreateMenuCommand struct {
	Title    string
	Path     string
	Icon     string
	ParentID *uint
	Order    int
	Visible  bool
}

// CreateMenuResult 创建菜单结果
type CreateMenuResult struct {
	ID uint
}
