package command

// CreateUserCommand 创建用户命令
type CreateUserCommand struct {
	Username string
	Email    string
	Password string
	FullName string
	RoleIDs  []uint // 可选：创建时分配角色
}
