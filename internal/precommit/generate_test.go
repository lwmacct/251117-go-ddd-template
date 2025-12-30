package precommit_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	op "github.com/lwmacct/251117-go-ddd-template/internal/domain/operation"
	"github.com/stretchr/testify/require"
)

// TestGenerate_SwaggerAnnotations 生成 Swagger 注解。
//
// 使用方法：GENERATE=1 go test -v -run TestGenerate_SwaggerAnnotations ./internal/precommit/
//
// 生成规则：
//   - @Summary ← operationMeta.Summary
//   - @Tags ← operationMeta.Tags
//   - @Accept json
//   - @Produce json
//   - @Security BearerAuth（非 public 操作）
//   - @Router ← operationMeta.Path + Method
//
// 保留手动注解：@Param、@Success、@Failure、@Description
func TestGenerate_SwaggerAnnotations(t *testing.T) {
	if os.Getenv("GENERATE") == "" {
		t.Skip("skipping generation test; set GENERATE=1 to run")
	}

	// 解析 routes.go 获取操作到处理器函数的映射
	routeBindings := parseRouteBindings(t)
	require.NotEmpty(t, routeBindings, "no route bindings found")

	// 按文件分组
	fileHandlers := make(map[string][]routeBinding)
	for _, rb := range routeBindings {
		fileHandlers[rb.handlerFile] = append(fileHandlers[rb.handlerFile], rb)
	}

	// 处理每个 handler 文件
	handlerDir := "../adapters/http/handler"
	for file, bindings := range fileHandlers {
		filePath := filepath.Join(handlerDir, file)
		if err := updateHandlerFile(t, filePath, bindings); err != nil {
			t.Errorf("failed to update %s: %v", file, err)
		} else {
			t.Logf("updated %s (%d handlers)", file, len(bindings))
		}
	}
}

// routeBinding 表示一个路由绑定
type routeBinding struct {
	operation   op.Operation
	handlerFile string // 如 "admin_user.go"
	handlerFunc string // 如 "CreateUser"
	handlerType string // 如 "AdminUserHandler"
}

// parseRouteBindings 解析 routes.go 获取路由绑定
func parseRouteBindings(t *testing.T) []routeBinding {
	t.Helper()

	content, err := os.ReadFile("../adapters/http/routes.go")
	require.NoError(t, err, "failed to read routes.go")

	var bindings []routeBinding

	// 匹配形如: {op.SysUsersCreate, deps.AdminUserHandler.CreateUser},
	bindingRe := regexp.MustCompile(`\{op\.(\w+),\s*deps\.(\w+)\.(\w+)\}`)

	for _, match := range bindingRe.FindAllStringSubmatch(string(content), -1) {
		if len(match) == 4 {
			opName := match[1]
			handlerType := match[2]
			funcName := match[3]

			// 直接通过常量名查找操作
			operation := op.ByConstName(opName)
			if operation == "" {
				continue
			}

			// 推断文件名：AdminUserHandler → admin_user.go
			fileName := handlerTypeToFile(handlerType)

			bindings = append(bindings, routeBinding{
				operation:   operation,
				handlerFile: fileName,
				handlerFunc: funcName,
				handlerType: handlerType,
			})
		}
	}

	return bindings
}

// handlerTypeToFile 将 handler 类型转换为文件名
// AdminUserHandler → admin_user.go
// PATHandler → pat.go (连续大写视为一个单词)
func handlerTypeToFile(handlerType string) string {
	// 特殊映射（不遵循 snake_case 规范的文件）
	exceptions := map[string]string{
		"AuditLogHandler": "auditlog.go",
		"TwoFAHandler":    "twofa.go",
	}
	if file, ok := exceptions[handlerType]; ok {
		return file
	}

	// 移除 Handler 后缀
	name := strings.TrimSuffix(handlerType, "Handler")

	// 驼峰转下划线（处理连续大写）
	var result strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			// 检查是否是连续大写字母的一部分
			prevUpper := name[i-1] >= 'A' && name[i-1] <= 'Z'
			nextLower := i+1 < len(name) && name[i+1] >= 'a' && name[i+1] <= 'z'
			// 只有当前一个字母是小写，或者当前是连续大写的最后一个且下一个是小写时，才加下划线
			if !prevUpper || nextLower {
				result.WriteByte('_')
			}
		}
		result.WriteRune(r)
	}

	return strings.ToLower(result.String()) + ".go"
}

// updateHandlerFile 更新 handler 文件中的注解
func updateHandlerFile(t *testing.T, filePath string, bindings []routeBinding) error {
	t.Helper()

	content, err := os.ReadFile(filePath) //nolint:gosec // G304: test tool with controlled paths
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

	// 构建函数名到绑定的映射
	funcToBinding := make(map[string]routeBinding)
	for _, b := range bindings {
		funcToBinding[b.handlerFunc] = b
	}

	// 第一遍：收集需要替换的区间
	type replacement struct {
		start    int      // 注解块起始行（包含）
		end      int      // func 行（不包含，会在替换时保留）
		newAnnot []string // 新注解内容
	}
	var replacements []replacement //nolint:prealloc // size unknown until iteration

	funcRe := regexp.MustCompile(`^func\s+\([^)]+\)\s+(\w+)\s*\(`)

	for i := range lines {
		matches := funcRe.FindStringSubmatch(lines[i])
		if len(matches) != 2 {
			continue
		}

		funcName := matches[1]
		binding, ok := funcToBinding[funcName]
		if !ok {
			continue
		}

		annotStart := findAnnotationStart(lines, i)
		if annotStart < 0 {
			continue
		}

		preserved := extractPreserved(lines[annotStart:i])
		newAnnot := generateAnnotation(binding.operation, preserved)
		replacements = append(replacements, replacement{
			start:    annotStart,
			end:      i,
			newAnnot: newAnnot,
		})
	}

	// 第二遍：从后向前替换（避免索引偏移）
	for i := len(replacements) - 1; i >= 0; i-- {
		r := replacements[i]
		// 删除旧注解，插入新注解
		newLines := make([]string, 0, len(lines)-r.end+r.start+len(r.newAnnot))
		newLines = append(newLines, lines[:r.start]...)
		newLines = append(newLines, r.newAnnot...)
		newLines = append(newLines, lines[r.end:]...)
		lines = newLines
	}

	//nolint:gosec // G306: intentionally 0644 for readable Go source files
	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
}

// findAnnotationStart 查找 Swagger 注解块的起始行
// 只识别 @ 开头的注解，保留函数文档注释
func findAnnotationStart(lines []string, funcLine int) int {
	// 从 func 行向上搜索
	// 找到第一个 @ 注解行作为注解块结束
	// 继续向上找到非 @ 注解的行作为注解块开始

	annotEnd := -1 // 最后一个 @ 注解行

	for i := funcLine - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])

		// 非注释行，停止搜索
		if !strings.HasPrefix(line, "//") {
			break
		}

		// 检查是否是 Swagger 注解（// 后面跟 @）
		content := strings.TrimPrefix(line, "//")
		content = strings.TrimSpace(content)
		if strings.HasPrefix(content, "@") {
			if annotEnd == -1 {
				annotEnd = i
			}
			continue
		}

		// 非 @ 注解的注释行
		// 如果已找到 @ 注解，这里是注解块的开始边界
		if annotEnd != -1 {
			return i + 1 // 注解块从下一行开始
		}
	}

	// 如果找到了 @ 注解但没有前导注释
	if annotEnd != -1 {
		// 从 annotEnd 向上找到所有连续的 @ 注解
		for i := annotEnd - 1; i >= 0; i-- {
			line := strings.TrimSpace(lines[i])
			if !strings.HasPrefix(line, "//") {
				return i + 1
			}
			content := strings.TrimPrefix(line, "//")
			content = strings.TrimSpace(content)
			if !strings.HasPrefix(content, "@") {
				return i + 1
			}
		}
		return 0
	}

	return -1 // 没有找到 @ 注解
}

// preservedAnnot 保留的注解
type preservedAnnot struct {
	description string
	params      []string
	success     []string
	failure     []string
}

// extractPreserved 提取需要保留的注解
func extractPreserved(lines []string) preservedAnnot {
	var result preservedAnnot
	descRe := regexp.MustCompile(`@Description\s+(.+)$`)

	for _, line := range lines {
		if matches := descRe.FindStringSubmatch(line); len(matches) == 2 {
			result.description = strings.TrimSpace(matches[1])
		}
		if strings.Contains(line, "@Param") {
			result.params = append(result.params, line)
		}
		if strings.Contains(line, "@Success") {
			result.success = append(result.success, line)
		}
		if strings.Contains(line, "@Failure") {
			result.failure = append(result.failure, line)
		}
	}

	return result
}

// generateAnnotation 生成注解块
func generateAnnotation(o op.Operation, preserved preservedAnnot) []string {
	var result []string

	// @Summary
	result = append(result, "// @Summary      "+op.Summary(o))

	// @Description（保留原有或使用 operationMeta）
	if preserved.description != "" {
		result = append(result, "// @Description  "+preserved.description)
	} else if desc := op.Description(o); desc != "" {
		result = append(result, "// @Description  "+desc)
	}

	// @Tags
	result = append(result, "// @Tags         "+op.Tags(o))

	// @Accept/@Produce
	result = append(result, "// @Accept       json")
	result = append(result, "// @Produce      json")

	// @Security（非 public 操作）
	if !op.IsPublic(o) {
		result = append(result, "// @Security     BearerAuth")
	}

	// @Param（保留原有）
	result = append(result, preserved.params...)

	// @Success/@Failure（保留原有）
	result = append(result, preserved.success...)
	result = append(result, preserved.failure...)

	// @Router
	swaggerPath := ginToSwaggerPath(op.Path(o))
	method := strings.ToLower(string(op.Method(o)))
	result = append(result, fmt.Sprintf("// @Router       %s [%s]", swaggerPath, method))

	return result
}

// ginToSwaggerPath 将 Gin 路径 :id 转换为 Swagger 路径 {id}
func ginToSwaggerPath(path string) string {
	re := regexp.MustCompile(`:(\w+)`)
	return re.ReplaceAllString(path, "{$1}")
}
