// Package routes 定义 HTTP 路由配置。
//
// 本包是适配器层的路由配置中心，集中管理所有 API 操作的 HTTP 元数据：
//   - HTTP 路由（Method, Path）
//   - Swagger 注解（Tags, Summary, Description）
//   - 审计开关（Audit bool，详情从 Operation 派生）
//
// # 单一数据源
//
// 所有路由配置集中在 [Registry]，其他模块通过函数获取：
//   - [Method]: 获取 HTTP 方法
//   - [Path]: 获取路由路径
//   - [NeedsAudit]: 判断是否需要审计
//   - [AuditAction]: 获取审计操作标识（派生自 Operation）
//   - [AllOperationDefinitions]: 权限列表（供前端）
//
// # 审计信息派生
//
// 审计详情从 [permission.Operation] 自动派生，无需手动配置：
//   - AuditCategory: 从 Operation.Type() 映射
//   - AuditAction: 格式为 {category}.{identifier}
//   - AuditOperation: 从 Operation.Identifier() 映射
//
// # 依赖关系
//
// 本包依赖 [permission] 包的领域类型：
//   - [permission.Operation]: 操作标识符
//   - [permission.AuditCategory]: 审计分类
//   - [permission.AuditOperation]: 审计操作类型
//
// # 设计原则
//
// HTTP 配置属于适配器层，与领域概念分离：
//   - 领域概念（Operation, Resource, 审计派生）→ domain/permission
//   - HTTP 配置（Method, Path, Swagger, Audit 开关）→ adapters/http/routes
package routes
