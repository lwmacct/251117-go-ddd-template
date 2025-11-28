// Package query 定义角色模块的读操作查询。
//
// 本包提供 RBAC 角色和权限的查询功能：
//   - GetRoleQuery: 获取单个角色详情（含权限列表）
//   - ListRolesQuery: 分页列表查询
//   - ListPermissionsQuery: 获取所有可用权限
package query

// GetRoleQuery 获取角色查询
type GetRoleQuery struct {
	RoleID uint
}
