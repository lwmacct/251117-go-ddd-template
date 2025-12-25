package role

// UpdateRoleCommand 更新角色命令
type UpdateRoleCommand struct {
	RoleID      uint
	DisplayName *string // 可选：显示名称
	Description *string // 可选：描述
}
