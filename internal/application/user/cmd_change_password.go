package user

// ChangePasswordCommand 修改密码命令
type ChangePasswordCommand struct {
	UserID      uint
	OldPassword string
	NewPassword string
}
