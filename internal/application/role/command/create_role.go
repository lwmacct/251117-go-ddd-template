// Package command 定义角色模块的写操作命令。
//
// 本包处理 RBAC 角色管理的命令：
//   - CreateRoleCommand: 创建新角色
//   - UpdateRoleCommand: 更新角色信息
//   - DeleteRoleCommand: 删除角色（系统角色禁止删除）
//   - SetPermissionsCommand: 设置角色权限（全量替换）
//
// 业务规则：
//   - 角色名称 (Name) 全局唯一
//   - 系统角色 (IsSystem=true) 不可删除
//   - 权限变更会触发相关用户的缓存失效
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
