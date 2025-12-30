package precommit_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/routes"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/permission"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAnnotation_MatchOperation 检查 handler @Router 注解与 operation 的一致性。
// 规则：每个 handler 的 @Router 路径必须与 operation 匹配。
func TestAnnotation_MatchOperation(t *testing.T) {
	annotations := parseHandlerAnnotations(t)
	require.NotEmpty(t, annotations, "no handler annotations found")

	// 构建 operation 索引 (method|swaggerPath -> operation)
	operationIndex := make(map[string]permission.Operation)
	for _, o := range routes.All() {
		// 将 Gin 路径 (:id) 转换为 Swagger 路径 ({id}) 以便比较
		swaggerPath := regexp.MustCompile(`:(\w+)`).ReplaceAllString(routes.Path(o), "{$1}")
		key := string(routes.Method(o)) + "|" + swaggerPath
		operationIndex[key] = o
	}

	for _, ann := range annotations {
		// 跳过非 API 路由（如 /health）
		if !strings.HasPrefix(ann.Path, "/api") {
			continue
		}

		key := ann.Method + "|" + ann.Path

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			// 检查路由是否在 operation 中
			_, exists := operationIndex[key]
			assert.True(t, exists, "handler route not in operation registry: %s %s", ann.Method, ann.Path)
		})
	}
}

// TestAnnotation_OperationCoverage 检查 operation 端点是否都有对应的 handler 注解。
// 规则：operation 中的每个端点都必须有带 @Router 注解的 handler。
func TestAnnotation_OperationCoverage(t *testing.T) {
	annotations := parseHandlerAnnotations(t)

	// 构建 handler 注解索引
	handlerIndex := make(map[string]bool)
	for _, ann := range annotations {
		key := ann.Method + "|" + ann.Path
		handlerIndex[key] = true
	}

	for _, o := range routes.All() {
		// 将 Gin 路径 (:id) 转换为 Swagger 路径 ({id})
		swaggerPath := regexp.MustCompile(`:(\w+)`).ReplaceAllString(routes.Path(o), "{$1}")
		key := string(routes.Method(o)) + "|" + swaggerPath

		t.Run(o.String(), func(t *testing.T) {
			assert.True(t, handlerIndex[key],
				"operation missing handler annotation: %s %s (Operation: %s)",
				routes.Method(o), routes.Path(o), o)
		})
	}
}

// TestAnnotation_RequiredFields 检查 Swagger 注解必填字段。
// 规则：每个 API 端点必须有 @Summary、@Tags、@Accept、@Produce。
func TestAnnotation_RequiredFields(t *testing.T) {
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

// TestAnnotation_SecurityRequired 检查非公开端点的 @Security 注解。
// 规则：除公开端点外，所有 API 都必须有 @Security BearerAuth。
func TestAnnotation_SecurityRequired(t *testing.T) {
	annotations := parseHandlerAnnotations(t)

	// 从 operation 获取公开端点列表
	publicPaths := make(map[string]bool)
	for _, o := range routes.All() {
		if o.IsPublic() {
			publicPaths[routes.Path(o)] = true
		}
	}

	for _, ann := range annotations {
		if !strings.HasPrefix(ann.Path, "/api") {
			continue
		}

		// 将 Swagger 路径 ({id}) 转换为 Gin 路径 (:id) 以便查找
		ginPath := regexp.MustCompile(`\{(\w+)\}`).ReplaceAllString(ann.Path, ":$1")

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			if publicPaths[ginPath] {
				// 公开端点不应有 @Security（或可选）
				return
			}

			// 需要认证的端点必须有 @Security BearerAuth
			assert.Equal(t, "BearerAuth", ann.Security,
				"non-public endpoint should have @Security BearerAuth: %s %s", ann.Method, ann.Path)
		})
	}
}

// TestAnnotation_TagsFormat 检查 @Tags 格式规范。
// 规则：格式为英文分组名（如「Admin - User Management」「Authentication」）。
func TestAnnotation_TagsFormat(t *testing.T) {
	annotations := parseHandlerAnnotations(t)
	// 匹配英文 Tags 格式：字母开头，可包含空格和连字符
	tagsRe := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9\s\-]*$`)

	for _, ann := range annotations {
		if !strings.HasPrefix(ann.Path, "/api") || ann.Tags == "" {
			continue
		}

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			assert.True(t, tagsRe.MatchString(ann.Tags),
				"@Tags should be English format (e.g. 'Admin - User Management'): got %q", ann.Tags)
		})
	}
}

// TestAnnotation_ContentType 检查 @Accept 和 @Produce 值。
// 规则：必须为 json。
func TestAnnotation_ContentType(t *testing.T) {
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

// TestAnnotation_SuccessDTOExists 检查 @Success 中的 DTO 类型是否存在。
// 规则：DTO 类型必须在 internal/application/{pkg}/dto.go 中定义。
// 例外：routes.* 类型（定义在 adapters/http/routes，是路由配置的一部分）。
func TestAnnotation_SuccessDTOExists(t *testing.T) {
	annotations := parseHandlerAnnotations(t)
	dtoTypes := loadDTOTypes(t)

	for _, ann := range annotations {
		if !strings.HasPrefix(ann.Path, "/api") || ann.SuccessDTO == "" {
			continue
		}

		// 跳过 routes.* 类型（定义在 adapters 层）
		if strings.HasPrefix(ann.SuccessDTO, "routes.") {
			continue
		}

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			assert.True(t, dtoTypes[ann.SuccessDTO],
				"@Success DTO type not found: %q\n  available types in package: check internal/application/{pkg}/dto.go",
				ann.SuccessDTO)
		})
	}
}

// TestAnnotation_ParamDTOExists 检查 @Param body 中引用 application 层 DTO 是否存在。
// 规则：带包前缀的 DTO（如 auth.LoginDTO）必须在 application 层定义。
// 注意：无包前缀的本地类型（如 CreateSettingRequest）不检查。
func TestAnnotation_ParamDTOExists(t *testing.T) {
	annotations := parseHandlerAnnotations(t)
	dtoTypes := loadDTOTypes(t)

	for _, ann := range annotations {
		if !strings.HasPrefix(ann.Path, "/api") || ann.ParamDTO == "" {
			continue
		}

		// 只检查带包前缀的类型（如 auth.LoginDTO）
		// 跳过 handler 本地定义的类型（无 . 分隔符）
		if !strings.Contains(ann.ParamDTO, ".") {
			continue
		}

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			assert.True(t, dtoTypes[ann.ParamDTO],
				"@Param body DTO type not found: %q\n  available types in package: check internal/application/{pkg}/dto.go",
				ann.ParamDTO)
		})
	}
}

// TestAnnotation_QueryTypeExists 检查 @Param query 中的结构体类型是否存在。
// 规则：Query 类型必须在 handler 包中定义。
func TestAnnotation_QueryTypeExists(t *testing.T) {
	annotations := parseHandlerAnnotations(t)
	queryTypes := loadHandlerQueryTypes(t)

	for _, ann := range annotations {
		if !strings.HasPrefix(ann.Path, "/api") || ann.QueryType == "" {
			continue
		}

		t.Run(ann.File+"/"+ann.Method+ann.Path, func(t *testing.T) {
			assert.True(t, queryTypes[ann.QueryType],
				"@Param query type not found: %q\n  check handler file for type definition",
				ann.QueryType)
		})
	}
}
