package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/lwmacct/251117-go-ddd-template/internal/commands/api"
	"github.com/lwmacct/251117-go-ddd-template/internal/commands/migrate"
	"github.com/lwmacct/251117-go-ddd-template/internal/commands/seed"
	"github.com/lwmacct/251117-go-ddd-template/internal/commands/worker"
	"github.com/urfave/cli/v3"
)

// buildCommands 根据环境变量条件性构建命令列表
func buildCommands() []*cli.Command {
	commands := []*cli.Command{
		api.Command,     // 🟢 API Service - REST API 服务
		migrate.Command, // 🔧 Database Migration - 数据库迁移工具
		seed.Command,    // 🌱 Database Seeder - 数据库种子数据填充
		worker.Command,  // 🔄 Queue Worker - 后台任务处理器
	}

	if os.Getenv("SHOW_CLI_ITEM") == "1" {
		// 可以在这里添加额外的调试或开发命令
		commands = append([]*cli.Command{}, commands...)
	}

	return commands
}

func main() {
	app := &cli.Command{
		Name:        "go-ddd-skeleton",
		Version:     "1.0.3",
		Usage:       "DDD 架构的 Golang 应用示例",
		Description: `这是一个基于 Domain-Driven Design (DDD) 的 Golang 应用程序。包含用户认证、订单管理等核心功能。`,
		Commands:    buildCommands(),
		Authors: []any{
			map[string]string{
				"name":  "Your Name",
				"email": "your.email@example.com",
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		slog.Error("Application failed to run", "error", err)
		os.Exit(1)
	}
}
