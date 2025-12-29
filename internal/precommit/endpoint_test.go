package precommit_test

import (
	"strings"
	"testing"

	op "github.com/lwmacct/251117-go-ddd-template/internal/domain/operation"
	"github.com/stretchr/testify/assert"
)

// TestOperation_PermissionsConsistent 检查操作的权限定义一致性。
// 规则：每个有权限的操作，其权限必须在 AllPermissions() 中存在。
func TestOperation_PermissionsConsistent(t *testing.T) {
	// 构建已定义权限集合
	defined := make(map[string]bool)
	for _, p := range op.AllPermissions() {
		defined[p.Code] = true
	}

	// 验证每个操作的权限都在派生列表中
	for _, o := range op.All() {
		perm := o.Permission()
		if perm == "" {
			continue // 公开端点无需权限
		}
		assert.True(t, defined[perm],
			"undefined permission %q in operation %s", perm, o)
	}
}

// TestOperation_OperationIDsUnique 检查 Operation 常量唯一性。
// 规则：所有操作的字符串值不能重复。
func TestOperation_OperationIDsUnique(t *testing.T) {
	seen := make(map[string]bool)

	for _, o := range op.All() {
		id := o.String()
		if seen[id] {
			t.Errorf("duplicate Operation: %s", id)
		}
		seen[id] = true
	}
}

// TestOperation_OperationIDFormat 检查 Operation 格式规范。
// 规则：格式为 domain.resource.action，domain 必须是 admin/user/auth。
func TestOperation_OperationIDFormat(t *testing.T) {
	validDomains := map[string]bool{"admin": true, "user": true, "auth": true}

	for _, o := range op.All() {
		parts := strings.Split(o.String(), ".")
		assert.GreaterOrEqual(t, len(parts), 2,
			"Operation should have at least 2 parts (domain.action): %s", o)

		// 验证 domain 是已知值
		assert.True(t, validDomains[parts[0]],
			"Operation domain should be admin/user/auth: %s", o)
	}
}

// TestOperation_PathFormat 检查路径格式规范。
// 规则：以 /api 开头，参数使用 :param 格式（Gin 风格）。
func TestOperation_PathFormat(t *testing.T) {
	for _, o := range op.All() {
		path := o.Path()
		// 验证以 /api 开头
		assert.True(t, strings.HasPrefix(path, "/api"),
			"path should start with /api: %s (Operation: %s)", path, o)

		// 验证参数格式为 :param（Gin 风格），不是 {param}（Swagger 风格）
		assert.NotContains(t, path, "{",
			"path should use :param not {param}: %s (Operation: %s)", path, o)
	}
}

// TestOperation_MethodValid 检查 HTTP 方法有效性。
// 规则：必须是 GET/POST/PUT/DELETE/PATCH 之一。
func TestOperation_MethodValid(t *testing.T) {
	validMethods := map[op.HTTPMethod]bool{
		op.HttpGET: true, op.HttpPOST: true, op.HttpPUT: true,
		op.HttpDELETE: true, op.HttpPATCH: true,
	}

	for _, o := range op.All() {
		method := o.Method()
		assert.True(t, validMethods[method],
			"invalid HTTP method %q for operation %s", method, o)
	}
}

// TestOperation_AuditActionFormat 检查审计操作格式。
// 规则：格式为 category.action（GitHub 风格）。
func TestOperation_AuditActionFormat(t *testing.T) {
	for _, o := range op.All() {
		action := o.AuditAction()
		if action == "" {
			continue // 不需要审计的操作
		}

		assert.Contains(t, action, ".",
			"AuditAction should be category.action format: %s (Operation: %s)", action, o)
	}
}

// TestOperation_MetadataComplete 检查操作元数据完整性。
// 规则：每个操作都必须有 Label 和 Description。
func TestOperation_MetadataComplete(t *testing.T) {
	for _, o := range op.All() {
		assert.NotEmpty(t, o.Label(),
			"missing Label for operation %s", o)
		assert.NotEmpty(t, o.Description(),
			"missing Description for operation %s", o)
	}
}
