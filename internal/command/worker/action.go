package worker

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/queue"
	redisinfra "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/redis"
	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/urfave/cli/v3"
)

// action 启动 worker
func action(ctx context.Context, cmd *cli.Command) error {
	cfg := cfgm.MustLoadCmd(cmd, config.DefaultConfig(), version.AppRawName)

	queueName := cmd.String("queue")
	concurrency := cmd.Int("concurrency")

	slog.Info("Starting worker",
		"queue", queueName,
		"concurrency", concurrency,
	)

	// 初始化 Redis 客户端（与 Telemetry 配置联动）
	redisClient, err := redisinfra.NewClient(ctx, cfg.Data.RedisURL, cfg.Telemetry.Enabled)
	if err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		return err
	}
	defer func() {
		if err := redisinfra.Close(redisClient); err != nil {
			slog.Error("Failed to close Redis connection", "error", err)
		}
	}()

	// 创建队列
	q := queue.NewRedisQueue(redisClient, queueName)

	// 创建处理器 (使用默认 handler，实际使用时应该实现自定义的 JobHandler)
	handler := &queue.DefaultJobHandler{}
	processor := queue.NewProcessor(q, handler, concurrency)

	// 启动处理器 (在 goroutine 中)
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
