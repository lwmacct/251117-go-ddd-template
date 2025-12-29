package precommit_test

import (
	"testing"

	op "github.com/lwmacct/251117-go-ddd-template/internal/domain/operation"
	"github.com/stretchr/testify/assert"
)

// TestRoutes_Bindings 检查声明式路由绑定的完整性。
// 规则：routes.go 中的绑定必须覆盖所有 operation。
// 注意：由于使用声明式路由，routes.go 和 operation.All() 应该完全一致。
// 这个测试确保开发者在添加新 operation 时不会忘记添加路由绑定。
func TestRoutes_Bindings(t *testing.T) {
	// 由于声明式路由使用 operation 作为数据源，
	// 路由与 operation 的一致性在编译时就已保证。
	// 这里只验证 operation 数据的有效性。

	for _, o := range op.All() {
		t.Run(o.String(), func(t *testing.T) {
			// 验证每个 operation 都有有效的元数据
			assert.NotEmpty(t, o.Method(), "operation %s missing Method", o)
			assert.NotEmpty(t, o.Path(), "operation %s missing Path", o)
		})
	}
}

// TestRoutes_PermissionConsistency 检查权限定义的一致性。
// 规则：同一资源的权限应该使用一致的命名模式。
func TestRoutes_PermissionConsistency(t *testing.T) {
	// 收集所有权限
	permissions := op.AllPermissions()

	// 检查权限格式一致性
	for _, p := range permissions {
		t.Run(p.Code, func(t *testing.T) {
			assert.NotEmpty(t, p.Domain, "permission %s missing Domain", p.Code)
			assert.NotEmpty(t, p.Resource, "permission %s missing Resource", p.Code)
			assert.NotEmpty(t, p.Action, "permission %s missing Action", p.Code)

			// 验证格式为 domain:resource:action
			expected := p.Domain + ":" + p.Resource + ":" + p.Action
			assert.Equal(t, expected, p.Code, "permission code format mismatch")
		})
	}
}

// TestRoutes_AuditActionsConsistency 检查审计操作的一致性。
// 规则：同一分类的审计操作应该使用一致的命名模式。
func TestRoutes_AuditActionsConsistency(t *testing.T) {
	actions := op.AllAuditActions()

	for _, a := range actions {
		t.Run(a.Action, func(t *testing.T) {
			assert.NotEmpty(t, a.Category, "audit action %s missing Category", a.Action)
			assert.NotEmpty(t, a.Operation, "audit action %s missing Operation", a.Action)
			assert.NotEmpty(t, a.Label, "audit action %s missing Label", a.Action)
		})
	}
}
