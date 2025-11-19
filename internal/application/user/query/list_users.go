// Package query 定义用户模块的查询
package query

// ListUsersQuery 获取用户列表查询
type ListUsersQuery struct {
	Offset int
	Limit  int
}
