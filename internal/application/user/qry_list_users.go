package user

// ListUsersQuery 获取用户列表查询
type ListUsersQuery struct {
	Page   int
	Limit  int
	Search string // 搜索关键词（用户名或邮箱）
}

// GetOffset 计算数据库查询偏移量
func (q ListUsersQuery) GetOffset() int {
	page := max(q.Page, 1)
	return (page - 1) * q.Limit
}
