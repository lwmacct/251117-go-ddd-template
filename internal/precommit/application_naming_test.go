package precommit_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestApplicationNaming_CommandSuffix 验证 commands.go 中的结构体以 Command 结尾
// 注意：cmd_*.go 文件是 Handler，由 TestApplicationNaming_HandlerSuffix 验证
func TestApplicationNaming_CommandSuffix(t *testing.T) {
	files := getApplicationFiles(t)

	for _, file := range files {
		filename := filepath.Base(file)

		// 只检查 commands.go 文件
		if filename != "commands.go" {
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

// TestApplicationNaming_QuerySuffix 验证 queries.go 中的结构体以 Query 结尾
// 注意：qry_*.go 文件是 Handler，由 TestApplicationNaming_HandlerSuffix 验证
func TestApplicationNaming_QuerySuffix(t *testing.T) {
	files := getApplicationFiles(t)

	for _, file := range files {
		filename := filepath.Base(file)

		// 只检查 queries.go 文件
		if filename != "queries.go" {
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

// TestApplicationNaming_HandlerSuffix 验证 cmd_*.go / qry_*.go 中至少有一个 Handler 结构体
// 允许文件中存在辅助结构体（如 Query 参数、Result 类型等）
func TestApplicationNaming_HandlerSuffix(t *testing.T) {
	files := getApplicationFiles(t)

	for _, file := range files {
		filename := filepath.Base(file)

		// 只检查 cmd_*.go 或 qry_*.go 文件
		if !strings.HasPrefix(filename, "cmd_") && !strings.HasPrefix(filename, "qry_") {
			continue
		}

		structs := parseStructs(t, file)
		hasHandler := false
		for _, s := range structs {
			if strings.HasSuffix(s.Name, "Handler") {
				hasHandler = true
				break
			}
		}

		t.Run(filename, func(t *testing.T) {
			assert.True(t, hasHandler,
				"file %s should contain at least one struct ending with 'Handler'", filename)
		})
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

// TestApplicationNaming_MapperFuncNaming 验证 mapper.go 中的函数以 To 开头、DTO 结尾
func TestApplicationNaming_MapperFuncNaming(t *testing.T) {
	files := getApplicationFiles(t)

	for _, file := range files {
		filename := filepath.Base(file)

		// 只检查 mapper.go 文件
		if filename != "mapper.go" {
			continue
		}

		funcs := parseFuncs(t, file)
		for _, f := range funcs {
			t.Run(f.File+"/"+f.Name, func(t *testing.T) {
				assert.True(t, strings.HasPrefix(f.Name, "To"),
					"func %q in %s should start with 'To'", f.Name, f.File)
				assert.True(t, strings.HasSuffix(f.Name, "DTO"),
					"func %q in %s should end with 'DTO'", f.Name, f.File)
			})
		}
	}
}
