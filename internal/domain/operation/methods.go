package operation

// ============================================================================
// Operation 方法
// ============================================================================

// String 返回操作标识符字符串
func (o Operation) String() string { return string(o) }

// Method 返回 HTTP 方法
func (o Operation) Method() HTTPMethod {
	if m, ok := operationRegistry[o]; ok {
		return m.Method
	}
	return ""
}

// Path 返回路由路径
func (o Operation) Path() string {
	if m, ok := operationRegistry[o]; ok {
		return m.Path
	}
	return ""
}

// AuditAction 返回审计操作标识
func (o Operation) AuditAction() string {
	if m, ok := operationRegistry[o]; ok {
		return m.AuditAction
	}
	return ""
}

// AuditCat 返回审计分类
func (o Operation) AuditCat() AuditCategory {
	if m, ok := operationRegistry[o]; ok {
		return m.AuditCategory
	}
	return ""
}

// AuditOp 返回审计操作类型
func (o Operation) AuditOp() AuditOperation {
	if m, ok := operationRegistry[o]; ok {
		return m.AuditOperation
	}
	return ""
}

// Label 返回中文标签
func (o Operation) Label() string {
	if m, ok := operationRegistry[o]; ok {
		return m.Label
	}
	return ""
}

// Description 返回英文描述
func (o Operation) Description() string {
	if m, ok := operationRegistry[o]; ok {
		return m.Description
	}
	return ""
}

// Group 返回 Swagger 分组
func (o Operation) Group() string {
	if m, ok := operationRegistry[o]; ok {
		return m.Group
	}
	return ""
}

// IsPublic 报告操作是否公开（无需权限检查）
//
// Operation-Centric RBAC: 默认所有操作需要权限检查，
// 只有显式标记 Public: true 的操作才是公开的。
func (o Operation) IsPublic() bool {
	if m, ok := operationRegistry[o]; ok {
		return m.Public
	}
	return false
}

// NeedsAudit 报告操作是否需要审计
func (o Operation) NeedsAudit() bool {
	if m, ok := operationRegistry[o]; ok {
		return m.AuditAction != ""
	}
	return false
}

// Valid 报告操作是否在注册表中
func (o Operation) Valid() bool {
	_, ok := operationRegistry[o]
	return ok
}

// ============================================================================
// AuditOperation 方法
// ============================================================================

//nolint:gochecknoglobals // 标签映射是只读配置
var auditOperationLabels = map[AuditOperation]string{
	AuditOpCreate:       "创建",
	AuditOpUpdate:       "更新",
	AuditOpDelete:       "删除",
	AuditOpAccess:       "访问",
	AuditOpAuthenticate: "认证",
}

// Label 返回审计操作的中文标签
func (o AuditOperation) Label() string {
	if label, ok := auditOperationLabels[o]; ok {
		return label
	}
	return string(o)
}

// ============================================================================
// AuditCategory 方法
// ============================================================================

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

// Label 返回审计分类的中文标签
func (c AuditCategory) Label() string {
	if label, ok := auditCategoryLabels[c]; ok {
		return label
	}
	return string(c)
}
