package precommit_test

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	SuccessDTO  string // from @Success, e.g., "user.UserWithRolesDTO"
}

// parseHandlerAnnotations 解析 handler 目录下所有 Go 文件的 Swagger 注解
func parseHandlerAnnotations(t *testing.T) []handlerAnnotation {
	t.Helper()

	handlerDir := "../adapters/http/handler"
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
	// @Success 200 {object} response.DataResponse[user.UserDTO] "描述"
	// @Success 200 {object} response.DataResponse[[]menu.MenuDTO] "描述" (数组类型)
	// 提取泛型参数中的 DTO 类型，如 user.UserDTO 或 []menu.MenuDTO
	successRe := regexp.MustCompile(`@Success\s+\d+\s+\{object\}\s+response\.\w+\[(\[\])?([^\]]+)\]`)

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
			if matches := successRe.FindStringSubmatch(line); len(matches) == 3 {
				current.SuccessDTO = matches[2] // 第二组是实际 DTO 类型
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

// loadDTOTypes 加载 application 层所有 DTO 类型
func loadDTOTypes(t *testing.T) map[string]bool {
	t.Helper()

	appDir := "../application"
	dtoTypes := make(map[string]bool)
	structRe := regexp.MustCompile(`^type\s+(\w+DTO)\s+struct`)

	err := filepath.Walk(appDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// 获取包名（目录名）
		pkgName := filepath.Base(filepath.Dir(path))

		file, err := os.Open(path) //nolint:gosec // 测试代码
		if err != nil {
			return nil
		}
		defer func() { _ = file.Close() }()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if matches := structRe.FindStringSubmatch(scanner.Text()); len(matches) == 2 {
				// 存储为 "pkg.TypeDTO" 格式
				fullName := pkgName + "." + matches[1]
				dtoTypes[fullName] = true
			}
		}
		return nil
	})

	require.NoError(t, err, "failed to load DTO types")
	return dtoTypes
}

// TestHandlerAnnotations_SuccessDTOExists 验证 @Success 中的 DTO 类型存在于 application 层
func TestHandlerAnnotations_SuccessDTOExists(t *testing.T) {
	annotations := parseHandlerAnnotations(t)
	dtoTypes := loadDTOTypes(t)

	for _, ann := range annotations {
		if !strings.HasPrefix(ann.Path, "/api") || ann.SuccessDTO == "" {
			continue
		}

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			assert.True(t, dtoTypes[ann.SuccessDTO],
				"@Success DTO type not found: %q\n  available types in package: check internal/application/{pkg}/dto.go",
				ann.SuccessDTO)
		})
	}
}
