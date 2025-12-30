// Package operation 定义统一操作注册表。
//
// 本包是系统的核心配置中心，集中定义所有 API 操作的元数据：
//   - HTTP 路由（Method, Path）
//   - 权限控制（基于 Operation URN）
//   - 审计日志（AuditAction, AuditCategory, AuditOperation）
//   - 显示信息（Label, Description, Group）
//
// # 单一数据源
//
// 所有操作定义集中在 [operationRegistry]，其他模块通过派生查询获取：
//   - [AllPermissions]: 权限列表
//   - [AllAuditActions]: 审计操作列表
//   - [Operation] 方法: Scope(), Type(), Identifier() 等
//
// # 操作标识符格式（URN 风格）
//
// [Operation] 采用统一资源名称（URN）格式：{scope}:{type}:{action}
//
//	public:auth:login       // 公开登录操作
//	sys:users:create        // 系统管理创建用户
//	self:profile:update     // 用户更新自己资料
//
// Scope 划分：
//   - public: 公开域（无需认证）
//   - sys:    系统管理域（需管理员权限）
//   - self:   用户自服务域（当前用户权限）
//
// # 资源标识符格式（URN 风格）
//
// [Resource] 同样采用 URN 格式：{scope}:{type}:{id}
//
//	sys:user:123            // 系统用户 123
//	self:user:@me           // 当前用户自身
//	*:*:*                   // 所有资源
//
// # 审计日志格式（GitHub 风格）
//
// AuditAction 采用 {category}.{action} 格式，与 GitHub Audit Log 一致：
//
//	user.create             // 创建用户
//	role.set_permissions    // 设置角色权限
//	auth.login              // 用户登录
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
//	op := operation.SysUsersCreate
//	fmt.Println(op.Scope())      // "sys"
//	fmt.Println(op.Type())       // "users"
//	fmt.Println(op.Identifier()) // "create"
//	fmt.Println(op.IsPublic())   // false
//
// # 新增操作流程
//
//  1. 在 constants.go 添加常量
//  2. 在 registry.go 添加元数据
//  3. 在 routes.go 添加路由绑定（一行）
//  4. 实现 Handler + Swagger 注解
//
// 权限、审计、路由中间件将自动生效。
//
// # 设计参考
//
//   - GitHub Audit Log: https://docs.github.com/en/organizations/keeping-your-organization-secure/managing-security-settings-for-your-organization/reviewing-the-audit-log-for-your-organization
//   - CADF (Cloud Auditing Data Federation): https://www.dmtf.org/standards/cadf
package operation
