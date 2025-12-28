// Package registry 定义 API 端点注册表。
//
// 提供统一的端点元数据管理，包括：
//   - Operation ID: 唯一操作标识符（用于日志、追踪、文档）
//   - Permission: 所需权限
//   - 其他元数据（描述、审计级别等）
//
// 使用示例：
//
//	ep := registry.ByPath("POST", "/api/admin/users")
//	fmt.Println(ep.OperationID) // "admin.users.create"
package registry
