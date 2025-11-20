// Package query 定义用户模块的查询
package query

// ListUsersQuery 获取用户列表查询
type ListUsersQuery struct {
	Offset int
	Limit  int
	Search string // 搜索关键词（用户名或邮箱）
}
