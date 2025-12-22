// Package query 定义菜单模块的读操作查询。
//
// 本包提供菜单数据的查询功能：
//   - GetMenuQuery: 获取单个菜单详情
//   - ListMenusQuery: 获取菜单树（自动构建层级结构）
package query

// GetMenuQuery 获取菜单查询
type GetMenuQuery struct {
	MenuID uint
}
