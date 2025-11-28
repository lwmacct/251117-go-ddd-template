// API 服务入口点。
//
// 本程序启动 HTTP REST API 服务器，提供：
//   - 用户认证（JWT + PAT 双模式）
//   - RBAC 权限管理
//   - 审计日志
//   - 系统配置管理
//
// 使用方式：
//
//	# 启动服务（默认读取 config.yaml）
//	./api
//
//	# 指定配置文件
//	./api --config /path/to/config.yaml
//
//	# 启用自动数据库迁移（仅开发环境）
//	./api --auto-migrate
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/lwmacct/251117-go-ddd-template/internal/commands/api"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:        "api",
		Version:     "1.0.0",
		Usage:       "REST API 服务",
		Description: "启动 HTTP API 服务器",
		Commands: []*cli.Command{
			api.Command,
		},
		// 默认执行 API 命令
		Action: api.Command.Action,
		Flags:  api.Command.Flags,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error("API command failed to run", "error", err)
		os.Exit(1)
	}
}
