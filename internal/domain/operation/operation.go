package operation

// Operation 统一操作标识符。
// 格式：{domain}:{module}.{action}
//
// 域划分：
//   - auth: 认证域（公开访问）
//   - sys:  系统管理域（需管理员权限）
//   - user: 用户自服务域（当前用户权限）
//
// 示例：
//   - auth:login         认证域登录操作
//   - sys:users.create   系统域用户模块创建操作
//   - user:profile.read  用户域个人资料读取操作
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
// Operation-Centric RBAC: Operation code 本身即权限标识符，
// 无需额外的 Permission 字段。权限检查直接使用 Operation code。
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

	// 公开访问标记（默认需要权限检查）
	Public bool // true=公开访问，无需权限检查
}
