package registry_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// routerRoute 从 router.go 解析的路由信息
type routerRoute struct {
	Method     string
	Path       string
	Permission string // 从 RequirePermission 提取
}

// extractPermission 从路由行提取权限常量值
func extractPermission(line string, permRe *regexp.Regexp, permConstants map[string]string) string {
	permMatches := permRe.FindStringSubmatch(line)
	if len(permMatches) != 2 {
		return ""
	}
	constName := permMatches[1]
	return permConstants[constName]
}

// parseRouterRoutes 解析 router.go 中的路由定义
func parseRouterRoutes(t *testing.T) []routerRoute {
	t.Helper()

	routerFile := "../router.go"
	content, err := os.ReadFile(routerFile)
	require.NoError(t, err, "failed to read router.go")

	lines := strings.Split(string(content), "\n")

	// 构建 group 前缀映射
	prefixes := map[string]string{
		"r": "",
	}

	// 正则匹配
	groupRe := regexp.MustCompile(`(\w+)\s*:=\s*(\w+)\.Group\("([^"]+)"\)`)
	routeRe := regexp.MustCompile(`(\w+)\.(GET|POST|PUT|DELETE|PATCH)\("([^"]+)"`)
	permRe := regexp.MustCompile(`RequirePermission\(permission\.(\w+)\)`)

	// 加载权限常量映射
	permConstants := loadPermissionConstants(t)

	routes := make([]routerRoute, 0, len(lines)/4) // 预估每 4 行一个路由

	for _, line := range lines {
		// 解析 Group 定义: api := r.Group("/api")
		if matches := groupRe.FindStringSubmatch(line); len(matches) == 4 {
			newGroup := matches[1]
			parentGroup := matches[2]
			groupPath := matches[3]
			if parentPrefix, ok := prefixes[parentGroup]; ok {
				prefixes[newGroup] = parentPrefix + groupPath
			}
		}

		// 解析路由定义: admin.POST("/users", ...)
		matches := routeRe.FindStringSubmatch(line)
		if len(matches) != 4 {
			continue
		}

		groupName := matches[1]
		method := matches[2]
		routePath := matches[3]

		prefix := prefixes[groupName]
		fullPath := prefix + routePath

		// 跳过非 API 路由
		if !strings.HasPrefix(fullPath, "/api") {
			continue
		}

		// 提取权限
		perm := extractPermission(line, permRe, permConstants)

		routes = append(routes, routerRoute{
			Method:     method,
			Path:       fullPath,
			Permission: perm,
		})
	}

	return routes
}

// loadPermissionConstants 加载权限常量映射
func loadPermissionConstants(t *testing.T) map[string]string {
	t.Helper()

	constFile := "../../../domain/permission/constants.go"
	content, err := os.ReadFile(constFile)
	require.NoError(t, err, "failed to read permission constants")

	constants := make(map[string]string)
	constRe := regexp.MustCompile(`(\w+)\s*=\s*"([^"]+)"`)

	for _, match := range constRe.FindAllStringSubmatch(string(content), -1) {
		if len(match) == 3 {
			constants[match[1]] = match[2]
		}
	}

	return constants
}

// TestRouterRoutes_MatchRegistry 验证 router.go 中的路由与 registry 一致
func TestRouterRoutes_MatchRegistry(t *testing.T) {
	routerRoutes := parseRouterRoutes(t)
	require.NotEmpty(t, routerRoutes, "no routes found in router.go")

	// 构建 registry 索引
	registryIndex := make(map[string]registry.Endpoint)
	for _, ep := range registry.All() {
		key := ep.Method + "|" + ep.Path
		registryIndex[key] = ep
	}

	for _, route := range routerRoutes {
		key := route.Method + "|" + route.Path

		t.Run(route.Method+" "+route.Path, func(t *testing.T) {
			ep, exists := registryIndex[key]
			if !assert.True(t, exists, "router route not in registry: %s %s", route.Method, route.Path) {
				return
			}

			// 验证权限一致
			if route.Permission != "" || ep.Permission != "" {
				assert.Equal(t, ep.Permission, route.Permission,
					"permission mismatch:\n  registry: %q\n  router:   %q",
					ep.Permission, route.Permission)
			}
		})
	}
}

// TestRegistry_AllRoutesInRouter 验证 registry 中的每个路由都在 router.go 中注册
func TestRegistry_AllRoutesInRouter(t *testing.T) {
	routerRoutes := parseRouterRoutes(t)

	// 构建 router 路由索引
	routerIndex := make(map[string]routerRoute)
	for _, route := range routerRoutes {
		key := route.Method + "|" + route.Path
		routerIndex[key] = route
	}

	for _, ep := range registry.All() {
		key := ep.Method + "|" + ep.Path

		t.Run(ep.OperationID, func(t *testing.T) {
			route, exists := routerIndex[key]
			if !assert.True(t, exists, "registry endpoint not in router.go: %s %s", ep.Method, ep.Path) {
				return
			}

			// 验证权限一致
			if ep.Permission != "" || route.Permission != "" {
				assert.Equal(t, ep.Permission, route.Permission,
					"permission mismatch:\n  registry: %q\n  router:   %q",
					ep.Permission, route.Permission)
			}
		})
	}
}
