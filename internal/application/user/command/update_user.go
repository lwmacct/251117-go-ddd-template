// Package command 定义用户模块的命令
package command

// UpdateUserCommand 更新用户命令
type UpdateUserCommand struct {
	UserID   uint
	FullName *string
	Avatar   *string
	Bio      *string
	Status   *string
}
