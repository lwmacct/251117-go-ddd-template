// Package permission 定义权限系统的核心领域概念。
//
// 本包提供统一资源名称（URN）格式的权限标识符，用于 RBAC 权限控制和审计日志。
//
// # 操作标识符格式
//
// [Operation] 采用 URN 格式：{scope}:{type}:{action}
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
// # 资源标识符格式
//
// [Resource] 同样采用 URN 格式：{scope}:{type}:{id}
//
//	sys:user:123            // 系统用户 123
//	self:user:@me           // 当前用户自身
//	*:*:*                   // 所有资源
//
// # 权限匹配
//
// 使用 [MatchOperation] 和 [MatchResource] 进行模式匹配：
//
//	MatchOperation("sys:*:*", "sys:users:create")     // true
//	MatchResource("self:user:@me", "self:user:123")   // false（需先解析 @me）
//
// # 变量解析
//
// 使用 [Resolver] 替换运行时变量：
//
//	r := NewResolver(map[string]string{"@me": "123"})
//	r.ResolveString("self:user:@me")  // "self:user:123"
//
// # 审计分类
//
// [AuditCategory] 和 [AuditOperation] 用于审计日志分类，遵循 GitHub Audit Log 风格。
package permission
