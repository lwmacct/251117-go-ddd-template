// Package query 定义 PAT 查询处理器
package query

// GetTokenQuery 获取 Token 查询
type GetTokenQuery struct {
	UserID  uint
	TokenID uint
}
