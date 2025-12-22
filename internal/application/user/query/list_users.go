package query

// ListUsersQuery 获取用户列表查询
type ListUsersQuery struct {
	Offset int
	Limit  int
	Search string // 搜索关键词（用户名或邮箱）
}
