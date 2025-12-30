package operation

import "github.com/lwmacct/251117-go-ddd-template/pkg/urn"

// Operation 统一操作标识符。
// 格式：{scope}:{type}:{action}
//
// Scope 划分：
//   - public: 公开域（无需认证）
//   - sys:    系统管理域（需管理员权限）
//   - self:   用户自服务域（当前用户权限）
//
// 示例：
//   - public:auth:login    公开登录操作
//   - sys:users:create     系统域用户创建操作
//   - self:profile:update  用户更新自己资料
//
// 本类型基于 [urn.URN] 格式，但作为独立类型以支持元数据方法。
type Operation string

// URN 返回底层 URN 类型（用于匹配等操作）。
func (o Operation) URN() urn.URN {
	return urn.URN(o)
}

// Scope 返回操作的 scope。
// "sys:users:create" → "sys"
func (o Operation) Scope() string {
	return o.URN().Scope()
}

// Type 返回操作的 type。
// "sys:users:create" → "users"
func (o Operation) Type() string {
	return o.URN().Type()
}

// Identifier 返回操作的 identifier（action）。
// "sys:users:create" → "create"
func (o Operation) Identifier() string {
	return o.URN().Identifier()
}

// Match 检查操作是否匹配指定模式。
func (o Operation) Match(pattern string) bool {
	return o.URN().Match(urn.URN(pattern))
}

// ============================================================================
// HTTP 方法类型
// ============================================================================

// HTTPMethod HTTP 请求方法。
type HTTPMethod string

// HTTP 方法常量。
const (
	HttpGET    HTTPMethod = "GET"
	HttpPOST   HTTPMethod = "POST"
	HttpPUT    HTTPMethod = "PUT"
	HttpDELETE HTTPMethod = "DELETE"
	HttpPATCH  HTTPMethod = "PATCH"
)

// ============================================================================
// 审计操作类型（粗粒度分类）
// ============================================================================

// AuditOperation 审计操作类型，遵循 GitHub Audit Log 风格。
type AuditOperation string

const (
	AuditOpCreate       AuditOperation = "create"
	AuditOpUpdate       AuditOperation = "update"
	AuditOpDelete       AuditOperation = "delete"
	AuditOpAccess       AuditOperation = "access"
	AuditOpAuthenticate AuditOperation = "authenticate"
)

// ============================================================================
// 审计分类
// ============================================================================

// AuditCategory 审计分类。
type AuditCategory string

const (
	AuditCatAuth        AuditCategory = "auth"
	AuditCatUser        AuditCategory = "user"
	AuditCatRole        AuditCategory = "role"
	AuditCatMenu        AuditCategory = "menu"
	AuditCatSetting     AuditCategory = "setting"
	AuditCatCache       AuditCategory = "cache"
	AuditCatProfile     AuditCategory = "profile"
	AuditCatToken       AuditCategory = "token"
	AuditCatUserSetting AuditCategory = "user_setting"
)

// ============================================================================
// 元数据结构
// ============================================================================

// operationMeta 操作元数据
//
// URN-Centric RBAC: Operation code 本身即权限标识符，
// 公开性通过 scope 判断（public: 前缀）。
type operationMeta struct {
	// HTTP 路由
	Method HTTPMethod // HTTP 方法
	Path   string     // 路由路径（Gin 格式），如 /api/admin/users/:id

	// 审计日志（GitHub Audit Log 风格）
	AuditAction    string         // 审计操作标识，如 user.create（空=不审计）
	AuditCategory  AuditCategory  // 审计分类
	AuditOperation AuditOperation // 审计操作类型，如 create

	// 显示信息
	Label       string // 中文标签，如 创建用户
	Description string // 英文描述，如 Create new user

	// Swagger 分组
	Group string // Swagger Tags，如 管理员 - 用户管理 (Admin - User)
}
