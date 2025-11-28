// Package query 定义用户模块的读操作查询。
//
// 本包实现 CQRS 模式中的 Query 端，包含：
//   - GetUserQuery: 获取单个用户详情
//   - ListUsersQuery: 分页列表查询
//   - SearchUsersQuery: 关键字搜索
//
// 设计特点：
//   - 查询结果直接返回领域实体或 DTO
//   - 支持按需加载关联数据（如 WithRoles）
//   - Handler 仅调用 QueryRepository，不修改状态
package query

// GetUserQuery 获取用户查询
type GetUserQuery struct {
	UserID    uint
	WithRoles bool // 是否包含角色信息
}
