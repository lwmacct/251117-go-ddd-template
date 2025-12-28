package precommit_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
