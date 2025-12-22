// Package command 定义菜单模块的写操作命令。
//
// 本包处理系统菜单的增删改操作：
//   - CreateMenuCommand: 创建菜单项
//   - UpdateMenuCommand: 更新菜单信息
//   - DeleteMenuCommand: 删除菜单（级联删除子菜单）
//   - ReorderMenusCommand: 批量调整菜单排序
//
// 树形结构：通过 ParentID 实现父子关系
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
