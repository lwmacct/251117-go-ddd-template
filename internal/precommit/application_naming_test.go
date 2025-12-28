package precommit_test

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// structInfo 从文件中提取的结构体信息
type structInfo struct {
	File string
	Name string
}

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

// TestApplicationNaming_CommandSuffix 验证 commands.go / cmd_*.go 中的结构体以 Command 结尾
func TestApplicationNaming_CommandSuffix(t *testing.T) {
	files := getApplicationFiles(t)

	for _, file := range files {
		filename := filepath.Base(file)

		// 只检查 commands.go 或 cmd_*.go 文件
		if filename != "commands.go" && !strings.HasPrefix(filename, "cmd_") {
			continue
		}

		structs := parseStructs(t, file)
		for _, s := range structs {
			t.Run(s.File+"/"+s.Name, func(t *testing.T) {
				assert.True(t, strings.HasSuffix(s.Name, "Command"),
					"struct %q in %s should end with 'Command'", s.Name, s.File)
			})
		}
	}
}

// TestApplicationNaming_QuerySuffix 验证 queries.go / qry_*.go 中的结构体以 Query 结尾
func TestApplicationNaming_QuerySuffix(t *testing.T) {
	files := getApplicationFiles(t)

	for _, file := range files {
		filename := filepath.Base(file)

		// 只检查 queries.go 或 qry_*.go 文件
		if filename != "queries.go" && !strings.HasPrefix(filename, "qry_") {
			continue
		}

		structs := parseStructs(t, file)
		for _, s := range structs {
			t.Run(s.File+"/"+s.Name, func(t *testing.T) {
				assert.True(t, strings.HasSuffix(s.Name, "Query"),
					"struct %q in %s should end with 'Query'", s.Name, s.File)
			})
		}
	}
}

// TestApplicationNaming_DTOSuffix 验证 dto.go 中的结构体以 DTO 结尾
func TestApplicationNaming_DTOSuffix(t *testing.T) {
	files := getApplicationFiles(t)

	for _, file := range files {
		filename := filepath.Base(file)

		// 只检查 dto.go 文件
		if filename != "dto.go" {
			continue
		}

		structs := parseStructs(t, file)
		for _, s := range structs {
			t.Run(s.File+"/"+s.Name, func(t *testing.T) {
				assert.True(t, strings.HasSuffix(s.Name, "DTO"),
					"struct %q in %s should end with 'DTO'", s.Name, s.File)
			})
		}
	}
}
