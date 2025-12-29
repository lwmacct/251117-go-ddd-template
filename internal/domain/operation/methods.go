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

// Permission 返回权限代码
func (o Operation) Permission() string {
	if m, ok := operationRegistry[o]; ok {
		return m.Permission
	}
	return ""
}

// Role 返回要求角色
func (o Operation) Role() string {
	if m, ok := operationRegistry[o]; ok {
		return m.Role
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

// AuditCategory 返回审计分类
func (o Operation) AuditCategory() string {
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

// IsPublic 报告操作是否公开（无需权限）
func (o Operation) IsPublic() bool {
	if m, ok := operationRegistry[o]; ok {
		return m.Permission == ""
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

// NeedsRole 报告操作是否需要角色验证
func (o Operation) NeedsRole() bool {
	if m, ok := operationRegistry[o]; ok {
		return m.Role != ""
	}
	return false
}

// Valid 报告操作是否在注册表中
func (o Operation) Valid() bool {
	_, ok := operationRegistry[o]
	return ok
}
