package operation

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

//nolint:gochecknoglobals // 标签映射是只读配置
var auditOperationLabels = map[AuditOperation]string{
	AuditOpCreate:       "创建",
	AuditOpUpdate:       "更新",
	AuditOpDelete:       "删除",
	AuditOpAccess:       "访问",
	AuditOpAuthenticate: "认证",
}

// Label 返回审计操作的中文标签。
func (o AuditOperation) Label() string {
	if label, ok := auditOperationLabels[o]; ok {
		return label
	}
	return string(o)
}

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

//nolint:gochecknoglobals // 标签映射是只读配置
var auditCategoryLabels = map[AuditCategory]string{
	AuditCatAuth:        "认证",
	AuditCatUser:        "用户",
	AuditCatRole:        "角色",
	AuditCatMenu:        "菜单",
	AuditCatSetting:     "配置",
	AuditCatCache:       "缓存",
	AuditCatProfile:     "个人资料",
	AuditCatToken:       "访问令牌",
	AuditCatUserSetting: "用户配置",
}

// Label 返回审计分类的中文标签。
func (c AuditCategory) Label() string {
	if label, ok := auditCategoryLabels[c]; ok {
		return label
	}
	return string(c)
}

// ============================================================================
// 元数据结构
// ============================================================================

// operationMeta 操作元数据。
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
