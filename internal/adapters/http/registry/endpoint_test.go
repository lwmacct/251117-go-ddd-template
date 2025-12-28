package registry_test

import (
	"strings"
	"testing"

	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/registry"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/permission"
	"github.com/stretchr/testify/assert"
)

// TestEndpoints_PermissionsExist 验证 registry 中的每个 Permission 都已定义为常量。
func TestEndpoints_PermissionsExist(t *testing.T) {
	// 构建已定义权限集合
	defined := make(map[string]bool)
	for _, p := range permission.AllDefinitions() {
		defined[p.Code] = true
	}

	// 验证每个 endpoint 的 permission 都已定义
	for _, ep := range registry.All() {
		if ep.Permission == "" {
			continue // 公开端点无需权限
		}
		assert.True(t, defined[ep.Permission],
			"undefined permission %q in endpoint %s %s (OperationID: %s)",
			ep.Permission, ep.Method, ep.Path, ep.OperationID)
	}
}

// TestEndpoints_OperationIDsUnique 验证所有 OperationID 唯一。
func TestEndpoints_OperationIDsUnique(t *testing.T) {
	seen := make(map[string]string) // operationID -> first endpoint path

	for _, ep := range registry.All() {
		if first, exists := seen[ep.OperationID]; exists {
			t.Errorf("duplicate OperationID %q:\n  first: %s\n  duplicate: %s %s",
				ep.OperationID, first, ep.Method, ep.Path)
		}
		seen[ep.OperationID] = ep.Method + " " + ep.Path
	}
}

// TestEndpoints_OperationIDFormat 验证 OperationID 格式为 domain.resource.action。
func TestEndpoints_OperationIDFormat(t *testing.T) {
	for _, ep := range registry.All() {
		parts := strings.Split(ep.OperationID, ".")
		assert.GreaterOrEqual(t, len(parts), 2,
			"OperationID should have at least 2 parts (domain.action): %s", ep.OperationID)

		// 验证 domain 是已知值
		validDomains := map[string]bool{"admin": true, "user": true, "auth": true}
		assert.True(t, validDomains[parts[0]],
			"OperationID domain should be admin/user/auth: %s", ep.OperationID)
	}
}

// TestEndpoints_PathFormat 验证路径格式规范。
func TestEndpoints_PathFormat(t *testing.T) {
	for _, ep := range registry.All() {
		// 验证以 /api 开头
		assert.True(t, strings.HasPrefix(ep.Path, "/api"),
			"path should start with /api: %s (OperationID: %s)", ep.Path, ep.OperationID)

		// 验证参数格式为 :param（Gin 风格），不是 {param}（Swagger 风格）
		assert.NotContains(t, ep.Path, "{",
			"path should use :param not {param}: %s (OperationID: %s)", ep.Path, ep.OperationID)
	}
}

// TestEndpoints_MethodValid 验证 HTTP 方法有效。
func TestEndpoints_MethodValid(t *testing.T) {
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "DELETE": true, "PATCH": true,
	}

	for _, ep := range registry.All() {
		assert.True(t, validMethods[ep.Method],
			"invalid HTTP method %q for endpoint %s (OperationID: %s)",
			ep.Method, ep.Path, ep.OperationID)
	}
}
