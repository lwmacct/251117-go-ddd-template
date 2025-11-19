// Package command 定义角色模块的命令
package command

// CreateRoleCommand 创建角色命令
type CreateRoleCommand struct {
	Name        string // 角色名称（唯一）
	DisplayName string // 显示名称
	Description string // 描述
}

// CreateRoleResult 创建角色结果
type CreateRoleResult struct {
	RoleID      uint
	Name        string
	DisplayName string
}
