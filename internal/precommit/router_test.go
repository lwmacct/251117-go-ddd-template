package precommit_test

import (
	"testing"

	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
