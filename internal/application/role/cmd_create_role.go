package role

// CreateRoleCommand 创建角色命令
type CreateRoleCommand struct {
	Name        string // 角色名称（唯一）
	DisplayName string // 显示名称
	Description string // 描述
}
