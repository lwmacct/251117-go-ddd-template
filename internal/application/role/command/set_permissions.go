package command

// SetPermissionsCommand 设置角色权限命令
type SetPermissionsCommand struct {
	RoleID        uint
	PermissionIDs []uint
}
