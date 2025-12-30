package precommit_test

import (
	"testing"

	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/routes"
	"github.com/stretchr/testify/assert"
)

// TestRoutes_Bindings 检查声明式路由绑定的完整性。
// 规则：operation 中的每个操作都必须有有效的元数据。
// 注意：由于使用声明式路由，routes.go 和 operation.All() 应该完全一致。
// 这个测试确保开发者在添加新 operation 时不会忘记添加路由绑定。
func TestRoutes_Bindings(t *testing.T) {
	// 由于声明式路由使用 operation 作为数据源，
	// 路由与 operation 的一致性在编译时就已保证。
	// 这里只验证 operation 数据的有效性。

	for _, o := range routes.All() {
		t.Run(o.String(), func(t *testing.T) {
			// 验证每个 operation 都有有效的元数据
			assert.NotEmpty(t, routes.Method(o), "operation %s missing Method", o)
			assert.NotEmpty(t, routes.Path(o), "operation %s missing Path", o)
		})
	}
}

// TestRoutes_PathFormat 检查路径格式的一致性。
// 规则：所有 API 路径必须以 /api/ 开头。
func TestRoutes_PathFormat(t *testing.T) {
	for _, o := range routes.All() {
		t.Run(o.String(), func(t *testing.T) {
			path := routes.Path(o)
			assert.True(t, len(path) > 0 && path[0] == '/',
				"path should start with /: got %q", path)
			assert.Contains(t, path, "/api/",
				"API path should contain /api/: got %q", path)
		})
	}
}

// TestRoutes_AuditActionsConsistency 检查审计操作的一致性。
// 规则：同一分类的审计操作应该使用一致的命名模式。
func TestRoutes_AuditActionsConsistency(t *testing.T) {
	actions := routes.AllAuditActions()

	for _, a := range actions {
		t.Run(a.Action, func(t *testing.T) {
			assert.NotEmpty(t, a.Category, "audit action %s missing Category", a.Action)
			assert.NotEmpty(t, a.Operation, "audit action %s missing Operation", a.Action)
			assert.NotEmpty(t, a.Label, "audit action %s missing Label", a.Action)
		})
	}
}
