package user

// CreateUserCommand 创建用户命令
type CreateUserCommand struct {
	Username string
	Email    string
	Password string
	FullName string
	Status   *string // 可选：初始状态，默认 "active"
	RoleIDs  []uint  // 可选：创建时分配角色
}
