// Package permission 定义系统权限常量。
//
// 权限采用三段式格式：domain:resource:action
//
// 域（Domain）分类：
//   - admin: 管理后台权限
//   - user: 用户中心权限
//   - api: API 访问权限
//
// 使用示例：
//
//	middleware.RequirePermission(permission.AdminUsersCreate)
//
// 权限常量命名规范：{Domain}{Resource}{Action}
// 例如：AdminUsersCreate 对应 "admin:users:create"
package permission
