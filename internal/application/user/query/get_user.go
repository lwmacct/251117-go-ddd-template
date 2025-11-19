// Package query 定义用户模块的查询
package query

// GetUserQuery 获取用户查询
type GetUserQuery struct {
	UserID      uint
	WithRoles   bool // 是否包含角色信息
}
