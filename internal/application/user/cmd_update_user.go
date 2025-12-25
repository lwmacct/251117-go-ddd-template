package user

// UpdateUserCommand 更新用户命令
type UpdateUserCommand struct {
	UserID   uint
	Username *string
	Email    *string
	FullName *string
	Avatar   *string
	Bio      *string
	Status   *string
}
