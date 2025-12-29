package operation

// Operation 统一操作标识符。
// 格式：{domain}.{resource}.{action}，如 admin.users.create
type Operation string

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
// 元数据结构
// ============================================================================

// operationMeta 操作元数据
type operationMeta struct {
	// HTTP 路由
	Method HTTPMethod // HTTP 方法
	Path   string     // 路由路径（Gin 格式），如 /api/admin/users/:id

	// 权限控制
	Permission string // 权限代码，如 admin:users:create（空=公开）
	Role       string // 要求角色，如 admin（空=不限角色）

	// 审计日志（GitHub Audit Log 风格）
	AuditAction    string         // 审计操作标识，如 user.create（空=不审计）
	AuditCategory  string         // 审计分类，如 user
	AuditOperation AuditOperation // 审计操作类型，如 create

	// 显示信息
	Label       string // 中文标签，如 创建用户
	Description string // 英文描述，如 Create new user

	// Swagger 分组
	Group string // Swagger Tags，如 管理员 - 用户管理 (Admin - User)
}
