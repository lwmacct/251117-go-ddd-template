package routes

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/audit"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/permission"
)

// ============================================================================
// 元数据访问函数
// ============================================================================

// Method 返回操作的 HTTP 方法。
func Method(o permission.Operation) HTTPMethod {
	if m, ok := Registry[o]; ok {
		return m.Method
	}
	return ""
}

// Path 返回操作的路由路径。
func Path(o permission.Operation) string {
	if m, ok := Registry[o]; ok {
		return m.Path
	}
	return ""
}

// Tags 返回操作的 Swagger Tags。
func Tags(o permission.Operation) string {
	if m, ok := Registry[o]; ok {
		return m.Tags
	}
	return ""
}

// Summary 返回操作的 Swagger Summary。
func Summary(o permission.Operation) string {
	if m, ok := Registry[o]; ok {
		return m.Summary
	}
	return ""
}

// Description 返回操作的描述。
func Description(o permission.Operation) string {
	if m, ok := Registry[o]; ok {
		return m.Description
	}
	return ""
}

// AuditAction 返回操作的审计操作标识（从 Operation 派生）。
func AuditAction(o permission.Operation) string {
	if m, ok := Registry[o]; ok && m.Audit {
		return audit.DeriveAction(o.Type(), o.Identifier())
	}
	return ""
}

// AuditCategory 返回操作的审计分类（从 Operation 派生）。
func AuditCategory(o permission.Operation) audit.Category {
	if m, ok := Registry[o]; ok && m.Audit {
		return audit.DeriveCategory(o.Type())
	}
	return ""
}

// AuditOperation 返回操作的审计操作类型（从 Operation 派生）。
func AuditOperation(o permission.Operation) audit.Operation {
	if m, ok := Registry[o]; ok && m.Audit {
		return audit.DeriveOperation(o.Identifier())
	}
	return ""
}

// NeedsAudit 报告操作是否需要审计。
func NeedsAudit(o permission.Operation) bool {
	if m, ok := Registry[o]; ok {
		return m.Audit
	}
	return false
}

// Valid 报告操作是否在注册表中。
func Valid(o permission.Operation) bool {
	_, ok := Registry[o]
	return ok
}
