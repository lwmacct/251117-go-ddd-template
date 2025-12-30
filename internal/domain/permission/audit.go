package permission

import "strings"

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

// String 返回审计操作的字符串表示。
func (o AuditOperation) String() string {
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

// String 返回审计分类的字符串表示。
func (c AuditCategory) String() string {
	return string(c)
}

// ============================================================================
// Operation → 审计信息派生
// ============================================================================

// typeToCategory type 到审计分类的映射。
//
//nolint:gochecknoglobals // 只读映射表
var typeToCategory = map[string]AuditCategory{
	"users":               AuditCatUser,
	"roles":               AuditCatRole,
	"settings":            AuditCatSetting,
	"settings-categories": AuditCatSetting,
	"cache":               AuditCatCache,
	"auth":                AuditCatAuth,
	"2fa":                 AuditCatAuth, // 特殊映射：2FA 归类为认证
	"profile":             AuditCatProfile,
	"password":            AuditCatProfile, // 特殊映射：密码归类为个人资料
	"account":             AuditCatProfile, // 特殊映射：账户归类为个人资料
	"tokens":              AuditCatToken,
}

// identifierToOperation identifier 到审计操作类型的映射。
//
//nolint:gochecknoglobals // 只读映射表
var identifierToOperation = map[string]AuditOperation{
	"create":       AuditOpCreate,
	"batch-create": AuditOpCreate,
	"update":       AuditOpUpdate,
	"delete":       AuditOpDelete,
	"list":         AuditOpAccess,
	"get":          AuditOpAccess,
	"login":        AuditOpAuthenticate,
	"login2fa":     AuditOpAuthenticate,
	"refresh":      AuditOpAuthenticate,
	"register":     AuditOpCreate,
}

// DeriveAuditCategory 从 Operation 派生审计分类。
//
// 优先使用映射表，未找到时使用 type 的单数形式。
func (o Operation) DeriveAuditCategory() AuditCategory {
	typ := o.Type()
	if cat, ok := typeToCategory[typ]; ok {
		return cat
	}
	return AuditCategory(singularize(typ))
}

// DeriveAuditOperation 从 Operation 派生审计操作类型。
//
// 优先使用映射表，未找到时默认为 update。
func (o Operation) DeriveAuditOperation() AuditOperation {
	id := o.Identifier()
	if op, ok := identifierToOperation[id]; ok {
		return op
	}
	return AuditOpUpdate // 默认
}

// DeriveAuditAction 从 Operation 派生审计操作标识。
//
// 格式：{category}.{identifier}，连字符转下划线。
// 示例：sys:users:create → user.create
func (o Operation) DeriveAuditAction() string {
	cat := o.DeriveAuditCategory()
	id := strings.ReplaceAll(o.Identifier(), "-", "_")
	return string(cat) + "." + id
}

// singularize 将复数形式转换为单数。
// 简单实现：移除尾部 s/es。
func singularize(s string) string {
	if len(s) == 0 {
		return s
	}
	// 特殊情况：以 -categories 结尾
	if before, ok := strings.CutSuffix(s, "-categories"); ok {
		return before + "_category"
	}
	// 简单规则：移除尾部 s
	if strings.HasSuffix(s, "s") {
		return s[:len(s)-1]
	}
	return s
}
