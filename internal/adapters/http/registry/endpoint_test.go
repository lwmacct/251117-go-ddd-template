package registry_test

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/registry"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/permission"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// handlerAnnotation 从 handler 文件解析的注解信息
type handlerAnnotation struct {
	File        string
	Method      string // from @Router [method]
	Path        string // from @Router path
	Permission  string // from @x-permission {"scope":"..."}
	Summary     string // from @Summary
	Description string // from @Description
	Tags        string // from @Tags
	Security    string // from @Security
	Accept      string // from @Accept
	Produce     string // from @Produce
}

// parseHandlerAnnotations 解析 handler 目录下所有 Go 文件的 Swagger 注解
func parseHandlerAnnotations(t *testing.T) []handlerAnnotation {
	t.Helper()

	handlerDir := "../handler"
	var annotations []handlerAnnotation

	// 正则匹配
	routerRe := regexp.MustCompile(`@Router\s+(\S+)\s+\[(\w+)\]`)
	permRe := regexp.MustCompile(`@x-permission\s+\{"scope":"([^"]+)"\}`)
	summaryRe := regexp.MustCompile(`@Summary\s+(.+)$`)
	descRe := regexp.MustCompile(`@Description\s+(.+)$`)
	tagsRe := regexp.MustCompile(`@Tags\s+(.+)$`)
	securityRe := regexp.MustCompile(`@Security\s+(\S+)`)
	acceptRe := regexp.MustCompile(`@Accept\s+(\S+)`)
	produceRe := regexp.MustCompile(`@Produce\s+(\S+)`)

	err := filepath.Walk(handlerDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, err := os.Open(path) //nolint:gosec // 测试代码，路径来自 filepath.Walk
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()

		var current handlerAnnotation
		current.File = filepath.Base(path)
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()

			// 解析各类注解
			if matches := routerRe.FindStringSubmatch(line); len(matches) == 3 {
				current.Path = matches[1]
				current.Method = strings.ToUpper(matches[2])
			}
			if matches := permRe.FindStringSubmatch(line); len(matches) == 2 {
				current.Permission = matches[1]
			}
			if matches := summaryRe.FindStringSubmatch(line); len(matches) == 2 {
				current.Summary = strings.TrimSpace(matches[1])
			}
			if matches := descRe.FindStringSubmatch(line); len(matches) == 2 {
				current.Description = strings.TrimSpace(matches[1])
			}
			if matches := tagsRe.FindStringSubmatch(line); len(matches) == 2 {
				current.Tags = strings.TrimSpace(matches[1])
			}
			if matches := securityRe.FindStringSubmatch(line); len(matches) == 2 {
				current.Security = matches[1]
			}
			if matches := acceptRe.FindStringSubmatch(line); len(matches) == 2 {
				current.Accept = matches[1]
			}
			if matches := produceRe.FindStringSubmatch(line); len(matches) == 2 {
				current.Produce = matches[1]
			}

			// 遇到 func 定义，保存当前注解
			if strings.HasPrefix(strings.TrimSpace(line), "func ") && current.Path != "" {
				annotations = append(annotations, current)
				current = handlerAnnotation{File: filepath.Base(path)}
			}
		}

		return scanner.Err()
	})

	require.NoError(t, err, "failed to parse handler files")
	return annotations
}

// TestHandlerAnnotations_MatchRegistry 验证 handler @Router 注解与 registry 一致
func TestHandlerAnnotations_MatchRegistry(t *testing.T) {
	annotations := parseHandlerAnnotations(t)
	require.NotEmpty(t, annotations, "no handler annotations found")

	// 构建 registry 索引 (method|path -> endpoint)
	registryIndex := make(map[string]registry.Endpoint)
	for _, ep := range registry.All() {
		// 将 Gin 路径 (:id) 转换为 Swagger 路径 ({id}) 以便比较
		swaggerPath := regexp.MustCompile(`:(\w+)`).ReplaceAllString(ep.Path, "{$1}")
		key := ep.Method + "|" + swaggerPath
		registryIndex[key] = ep
	}

	for _, ann := range annotations {
		// 跳过非 API 路由（如 /health）
		if !strings.HasPrefix(ann.Path, "/api") {
			continue
		}

		key := ann.Method + "|" + ann.Path

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			// 检查路由是否在 registry 中
			ep, exists := registryIndex[key]
			if !assert.True(t, exists, "handler route not in registry: %s %s", ann.Method, ann.Path) {
				return
			}

			// 检查权限是否一致
			if ann.Permission != "" || ep.Permission != "" {
				assert.Equal(t, ep.Permission, ann.Permission,
					"permission mismatch for %s %s\n  registry: %q\n  handler:  %q",
					ann.Method, ann.Path, ep.Permission, ann.Permission)
			}
		})
	}
}

// TestRegistry_AllRoutesHaveHandlerAnnotation 验证 registry 中的每个路由都有对应的 handler 注解
func TestRegistry_AllRoutesHaveHandlerAnnotation(t *testing.T) {
	annotations := parseHandlerAnnotations(t)

	// 构建 handler 注解索引
	handlerIndex := make(map[string]bool)
	for _, ann := range annotations {
		key := ann.Method + "|" + ann.Path
		handlerIndex[key] = true
	}

	for _, ep := range registry.All() {
		// 将 Gin 路径 (:id) 转换为 Swagger 路径 ({id})
		swaggerPath := regexp.MustCompile(`:(\w+)`).ReplaceAllString(ep.Path, "{$1}")
		key := ep.Method + "|" + swaggerPath

		t.Run(ep.OperationID, func(t *testing.T) {
			assert.True(t, handlerIndex[key],
				"registry endpoint missing handler annotation: %s %s (OperationID: %s)",
				ep.Method, ep.Path, ep.OperationID)
		})
	}
}

// TestHandlerAnnotations_RequiredFields 验证每个端点都有必需的 Swagger 注解
func TestHandlerAnnotations_RequiredFields(t *testing.T) {
	annotations := parseHandlerAnnotations(t)

	for _, ann := range annotations {
		if !strings.HasPrefix(ann.Path, "/api") {
			continue
		}

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			assert.NotEmpty(t, ann.Summary,
				"missing @Summary for %s %s", ann.Method, ann.Path)
			assert.NotEmpty(t, ann.Tags,
				"missing @Tags for %s %s", ann.Method, ann.Path)
			assert.NotEmpty(t, ann.Accept,
				"missing @Accept for %s %s", ann.Method, ann.Path)
			assert.NotEmpty(t, ann.Produce,
				"missing @Produce for %s %s", ann.Method, ann.Path)
		})
	}
}

// TestHandlerAnnotations_SecurityRequired 验证需要认证的端点有 @Security 注解
func TestHandlerAnnotations_SecurityRequired(t *testing.T) {
	annotations := parseHandlerAnnotations(t)

	// 公开端点列表（不需要 @Security）
	publicPaths := map[string]bool{
		"/api/auth/register":   true,
		"/api/auth/login":      true,
		"/api/auth/login/2fa":  true,
		"/api/auth/refresh":    true,
		"/api/auth/captcha":    true,
		"/api/auth/2fa/setup":  true, // 需要 JWT 但无权限检查
		"/api/auth/2fa/verify": true,
	}

	for _, ann := range annotations {
		if !strings.HasPrefix(ann.Path, "/api") {
			continue
		}

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			if publicPaths[ann.Path] {
				// 公开端点不应有 @Security（或可选）
				return
			}

			// 需要认证的端点必须有 @Security BearerAuth
			assert.Equal(t, "BearerAuth", ann.Security,
				"non-public endpoint should have @Security BearerAuth: %s %s", ann.Method, ann.Path)
		})
	}
}

// TestHandlerAnnotations_TagsFormat 验证 @Tags 格式：中文名 (English Name)
func TestHandlerAnnotations_TagsFormat(t *testing.T) {
	annotations := parseHandlerAnnotations(t)
	tagsRe := regexp.MustCompile(`^.+\s+\([A-Za-z].*\)$`)

	for _, ann := range annotations {
		if !strings.HasPrefix(ann.Path, "/api") || ann.Tags == "" {
			continue
		}

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			assert.True(t, tagsRe.MatchString(ann.Tags),
				"@Tags should match format '中文名 (English Name)': got %q", ann.Tags)
		})
	}
}

// TestHandlerAnnotations_ContentType 验证 @Accept 和 @Produce 为 json
func TestHandlerAnnotations_ContentType(t *testing.T) {
	annotations := parseHandlerAnnotations(t)

	for _, ann := range annotations {
		if !strings.HasPrefix(ann.Path, "/api") {
			continue
		}

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			if ann.Accept != "" {
				assert.Equal(t, "json", ann.Accept,
					"@Accept should be 'json': got %q", ann.Accept)
			}
			if ann.Produce != "" {
				assert.Equal(t, "json", ann.Produce,
					"@Produce should be 'json': got %q", ann.Produce)
			}
		})
	}
}

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
