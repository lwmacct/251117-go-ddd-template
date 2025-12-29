// Package operation 定义统一操作注册表。
//
// 本包是系统的核心配置中心，集中定义所有 API 操作的元数据：
//   - HTTP 路由（Method, Path）
//   - 权限控制（Permission, Role）
//   - 审计日志（AuditAction, AuditCategory, AuditOperation）
//   - 显示信息（Label, Description, Group）
//
// # 单一数据源
//
// 所有操作定义集中在 [operationRegistry]，其他模块通过派生查询获取：
//   - [AllPermissions]: 权限列表（替代 permission/constants.go）
//   - [AllAuditActions]: 审计操作列表
//   - [Operation] 方法: Permission(), AuditAction(), Label() 等
//
// # 操作标识符格式
//
// Operation 采用点分隔的三段式格式：{domain}.{resource}.{action}
//
//	admin.users.create    // 管理员创建用户
//	user.profile.update   // 用户更新资料
//	auth.login            // 用户登录（无 resource）
//
// # 审计日志格式（GitHub 风格）
//
// AuditAction 采用 {category}.{action} 格式，与 GitHub Audit Log 一致：
//
//	user.create           // 创建用户
//	role.set_permissions  // 设置角色权限
//	auth.login            // 用户登录
//
// AuditOperation 为粗粒度操作类型：
//   - create: 创建
//   - update: 更新
//   - delete: 删除
//   - access: 访问
//   - authenticate: 认证
//
// # 使用示例
//
//	op := operation.AdminUsersCreate
//	fmt.Println(op.Permission())  // "admin:users:create"
//	fmt.Println(op.AuditAction()) // "user.create"
//	fmt.Println(op.Label())       // "创建用户"
//	fmt.Println(op.IsPublic())    // false
//	fmt.Println(op.NeedsAudit())  // true
//
// # 新增操作流程
//
//  1. 在 registry.go 添加常量和元数据
//  2. 在 routes.go 添加路由绑定（一行）
//  3. 实现 Handler + Swagger 注解
//
// 权限、审计、路由中间件将自动生效。
//
// # 设计参考
//
//   - GitHub Audit Log: https://docs.github.com/en/organizations/keeping-your-organization-secure/managing-security-settings-for-your-organization/reviewing-the-audit-log-for-your-organization
//   - CADF (Cloud Auditing Data Federation): https://www.dmtf.org/standards/cadf
package operation
