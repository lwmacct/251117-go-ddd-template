package operation

// ============================================================================
// Operation 元数据访问函数
// ============================================================================
//
// 由于 Operation 是 pkg/operation 的类型别名，Go 不允许为其添加方法。
// 因此使用包级函数访问注册表中的元数据。

// Method 返回操作的 HTTP 方法。
func Method(o Operation) HTTPMethod {
	if m, ok := operationRegistry.Get(o); ok {
		return m.Method
	}
	return ""
}

// Path 返回操作的路由路径。
func Path(o Operation) string {
	if m, ok := operationRegistry.Get(o); ok {
		return m.Path
	}
	return ""
}

// AuditAction 返回操作的审计操作标识。
func AuditAction(o Operation) string {
	if m, ok := operationRegistry.Get(o); ok {
		return m.AuditAction
	}
	return ""
}

// AuditCat 返回操作的审计分类。
func AuditCat(o Operation) AuditCategory {
	if m, ok := operationRegistry.Get(o); ok {
		return m.AuditCategory
	}
	return ""
}

// AuditOp 返回操作的审计操作类型。
func AuditOp(o Operation) AuditOperation {
	if m, ok := operationRegistry.Get(o); ok {
		return m.AuditOperation
	}
	return ""
}

// Summary 返回操作的中文摘要（@Summary）。
func Summary(o Operation) string {
	if m, ok := operationRegistry.Get(o); ok {
		return m.Summary
	}
	return ""
}

// Description 返回操作的英文描述。
func Description(o Operation) string {
	if m, ok := operationRegistry.Get(o); ok {
		return m.Description
	}
	return ""
}

// Tags 返回操作的 Swagger 分组（@Tags）。
func Tags(o Operation) string {
	if m, ok := operationRegistry.Get(o); ok {
		return m.Tags
	}
	return ""
}

// IsPublic 报告操作是否公开（无需权限检查）。
//
// URN-Centric RBAC: 通过 scope 判断是否公开。
// scope 为 "public" 的操作无需权限检查。
func IsPublic(o Operation) bool {
	return o.Scope() == "public"
}

// NeedsAudit 报告操作是否需要审计。
func NeedsAudit(o Operation) bool {
	if m, ok := operationRegistry.Get(o); ok {
		return m.AuditAction != ""
	}
	return false
}

// Valid 报告操作是否在注册表中。
func Valid(o Operation) bool {
	_, ok := operationRegistry.Get(o)
	return ok
}
