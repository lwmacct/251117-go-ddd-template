package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/lwmacct/251117-bd-vmalert/internal/commands/api"
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
