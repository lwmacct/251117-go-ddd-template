// Package worker 提供后台任务处理命令
package worker

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/queue"
	redisinfra "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/redis"
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
	Action: runWorker,
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

// runWorker 启动 worker
func runWorker(ctx context.Context, cmd *cli.Command) error {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return err
	}

	queueName := cmd.String("queue")
	concurrency := cmd.Int("concurrency")

	slog.Info("Starting worker",
		"queue", queueName,
		"concurrency", concurrency,
	)

	// 初始化 Redis 客户端
	redisClient, err := redisinfra.NewClient(ctx, cfg.Data.RedisURL)
	if err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		return err
	}
	defer redisinfra.Close(redisClient)

	// 创建队列
	q := queue.NewRedisQueue(redisClient, queueName)

	// 创建处理器（使用默认 handler，实际使用时应该实现自定义的 JobHandler）
	handler := &queue.DefaultJobHandler{}
	processor := queue.NewProcessor(q, handler, concurrency)

	// 启动处理器（在 goroutine 中）
	go processor.Start(ctx)

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down worker...")
	processor.Stop()

	slog.Info("Worker stopped")
	return nil
}
