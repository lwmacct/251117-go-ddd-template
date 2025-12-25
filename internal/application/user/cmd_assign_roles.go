package user

// AssignRolesCommand 分配用户角色命令
type AssignRolesCommand struct {
	UserID  uint
	RoleIDs []uint
}
