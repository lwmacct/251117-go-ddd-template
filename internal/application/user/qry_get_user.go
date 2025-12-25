package user

// GetUserQuery 获取用户查询
type GetUserQuery struct {
	UserID    uint
	WithRoles bool // 是否包含角色信息
}
