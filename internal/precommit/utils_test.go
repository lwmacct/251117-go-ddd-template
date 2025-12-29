package precommit_test

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// ============================================================
// Types
// ============================================================

// structInfo 从文件中提取的结构体信息
type structInfo struct {
	File string
	Name string
}

// funcInfo 从文件中提取的函数信息
type funcInfo struct {
	File string
	Name string
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
	SuccessDTO  string // from @Success, e.g., "user.UserWithRolesDTO"
	ParamDTO    string // from @Param body, e.g., "auth.LoginDTO"
	QueryType   string // from @Param query, e.g., "handler.ListUsersQuery"
}

// ============================================================
// Application Layer Helpers
// ============================================================

// parseStructs 解析 Go 文件中的结构体定义
func parseStructs(t *testing.T, filePath string) []structInfo {
	t.Helper()

	file, err := os.Open(filePath) //nolint:gosec // 测试代码
	if err != nil {
		return nil
	}
	defer func() { _ = file.Close() }()

	structRe := regexp.MustCompile(`^\s*type\s+([A-Za-z_][A-Za-z0-9_]*)\s+struct`)
	var structs []structInfo
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if matches := structRe.FindStringSubmatch(scanner.Text()); len(matches) == 2 {
			structs = append(structs, structInfo{
				File: filepath.Base(filePath),
				Name: matches[1],
			})
		}
	}

	return structs
}

// parseFuncs 解析 Go 文件中的函数定义（仅顶级导出函数，不含方法）
func parseFuncs(t *testing.T, filePath string) []funcInfo {
	t.Helper()

	file, err := os.Open(filePath) //nolint:gosec // 测试代码
	if err != nil {
		return nil
	}
	defer func() { _ = file.Close() }()

	// 匹配顶级导出函数：func FuncName(
	funcRe := regexp.MustCompile(`^func\s+([A-Z][A-Za-z0-9_]*)\s*\(`)
	var funcs []funcInfo
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if matches := funcRe.FindStringSubmatch(scanner.Text()); len(matches) == 2 {
			funcs = append(funcs, funcInfo{
				File: filepath.Base(filePath),
				Name: matches[1],
			})
		}
	}

	return funcs
}

// getApplicationFiles 获取 application 目录下的所有 Go 文件
func getApplicationFiles(t *testing.T) []string {
	t.Helper()

	appDir := "../application"
	var files []string

	err := filepath.Walk(appDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		// 跳过测试文件和 handler 文件
		if strings.HasSuffix(path, "_test.go") || strings.HasSuffix(path, "_handler.go") {
			return nil
		}
		files = append(files, path)
		return nil
	})

	if err != nil {
		t.Logf("warning: failed to walk application directory: %v", err)
	}

	return files
}

// ============================================================
// Handler Annotation Helpers
// ============================================================

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
	// @Param request body auth.LoginDTO true "登录凭证"
	// 提取 body 参数中的 DTO 类型，如 auth.LoginDTO
	paramBodyRe := regexp.MustCompile(`@Param\s+\S+\s+body\s+(\S+)\s+`)
	// @Param params query handler.ListUsersQuery false "查询参数"
	// 提取 query 参数中的结构体类型（带 handler. 前缀或本地类型）
	paramQueryRe := regexp.MustCompile(`@Param\s+\S+\s+query\s+(handler\.\w+|\w+Query)\s+`)

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
			if matches := paramBodyRe.FindStringSubmatch(line); len(matches) == 2 {
				current.ParamDTO = matches[1]
			}
			if matches := paramQueryRe.FindStringSubmatch(line); len(matches) == 2 {
				current.QueryType = matches[1]
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

// loadDTOTypes 使用 go doc 加载 application 层所有 DTO 类型
func loadDTOTypes(t *testing.T) map[string]bool {
	t.Helper()

	dtoTypes := make(map[string]bool)
	typeRe := regexp.MustCompile(`^type\s+(\w+DTO)\s+struct`)

	// 获取 application 下的所有子包
	appDir := "../application"
	entries, err := os.ReadDir(appDir)
	require.NoError(t, err, "failed to read application directory")

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pkgName := entry.Name()
		pkgPath := "./internal/application/" + pkgName + "/"

		// 使用 go doc 获取包的类型列表
		cmd := exec.Command("go", "doc", pkgPath) //nolint:gosec,noctx // 测试代码，无需 context
		cmd.Dir = "../.."                         // 从项目根目录执行（internal/precommit -> 项目根）
		output, err := cmd.Output()
		if err != nil {
			continue // 跳过无法解析的包
		}

		// 解析 go doc 输出中的 DTO 类型
		for line := range strings.SplitSeq(string(output), "\n") {
			if matches := typeRe.FindStringSubmatch(line); len(matches) == 2 {
				fullName := pkgName + "." + matches[1]
				dtoTypes[fullName] = true
			}
		}
	}

	require.NotEmpty(t, dtoTypes, "no DTO types found")
	return dtoTypes
}

// loadHandlerQueryTypes 加载 handler 目录中定义的 Query 结构体类型
func loadHandlerQueryTypes(t *testing.T) map[string]bool {
	t.Helper()

	handlerDir := "../adapters/http/handler"
	queryTypes := make(map[string]bool)
	// 匹配 type XXXQuery struct（不用行首锚点，因为是全文匹配）
	typeRe := regexp.MustCompile(`type\s+(\w+Query)\s+struct`)

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

		content, err := os.ReadFile(path) //nolint:gosec // 测试代码
		if err != nil {
			return nil //nolint:nilerr // 跳过无法读取的文件
		}

		for _, match := range typeRe.FindAllStringSubmatch(string(content), -1) {
			if len(match) == 2 {
				// 存储两种格式：带 handler. 前缀和不带前缀
				queryTypes[match[1]] = true
				queryTypes["handler."+match[1]] = true
			}
		}
		return nil
	})

	require.NoError(t, err, "failed to walk handler directory")
	return queryTypes
}
