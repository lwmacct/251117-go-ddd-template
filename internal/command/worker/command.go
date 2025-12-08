// Package worker 提供后台任务处理命令
package worker

import (
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/urfave/cli/v3"
)

// Command 定义 worker 命令
var Command = &cli.Command{
	Name:  "worker",
	Usage: "启动后台任务处理器",
	Description: `
   启动 Worker 进程处理队列中的任务。
   支持并发处理和优雅关闭，会等待当前任务完成后再退出。

   示例：
   - 默认配置启动：worker
   - 指定队列和并发数：worker --queue jobs --concurrency 10
	`,
	Action:   action,
	Commands: []*cli.Command{version.Command},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "queue",
			Aliases: []string{"q"},
			Value:   "default",
			Usage:   "队列名称",
		},
		&cli.IntFlag{
			Name:    "concurrency",
			Aliases: []string{"c"},
			Value:   5,
			Usage:   "并发处理数",
		},
	},
}
