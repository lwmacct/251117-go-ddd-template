package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpserver "github.com/lwmacct/251117-go-ddd-template/internal/adapters/http"
	"github.com/lwmacct/251117-go-ddd-template/internal/bootstrap"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/config"
	"github.com/urfave/cli/v3"
)

// Command 定义 API 服务命令
var Command = &cli.Command{
	Name:    "api",
	Aliases: []string{"serve", "server"},
	Usage:   "启动 REST API 服务",
	Description: `
   启动 HTTP API 服务器，提供 RESTful API 接口。
   服务器支持优雅关闭，会等待正在处理的请求完成后再退出。

   配置优先级（从低到高）：
   1. 默认值
   2. 配置文件 (config.yaml)
   3. 环境变量 (APP_SERVER_ADDR)
   4. 命令行参数 (--addr)
	`,
	Action: runAPIServer,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "addr",
			Aliases: []string{"a"},
			Usage:   "服务器监听地址 (格式: host:port)",
		},
		&cli.StringFlag{
			Name:    "static",
			Aliases: []string{"s"},
			Usage:   "静态资源目录路径",
		},
	},
}

// runAPIServer 执行 API 服务器启动逻辑
func runAPIServer(ctx context.Context, cmd *cli.Command) error {
	// 加载配置（优先级：默认值 -> 配置文件 -> 环境变量）
	cfg, err := config.LoadWithDefaults()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// 如果用户显式指定了 --addr 参数，则覆盖配置
	if cmd.IsSet("addr") {
		cfg.Server.Addr = cmd.String("addr")
	}

	// 如果用户显式指定了 --static 参数，则覆盖配置
	if cmd.IsSet("static") {
		cfg.Server.StaticDir = cmd.String("static")
	}

	// 初始化容器（依赖注入）
	opts := &bootstrap.ContainerOptions{
		AutoMigrate: cfg.Data.AutoMigrate, // 从配置读取是否自动迁移
	}
	container, err := bootstrap.NewContainer(cfg, opts)
	if err != nil {
		slog.Error("Failed to initialize container", "error", err)
		os.Exit(1)
	}
	// 确保在退出时关闭所有资源
	defer func() {
		if err := container.Close(); err != nil {
			slog.Error("Failed to close container resources", "error", err)
		}
	}()

	slog.Info("Configuration loaded",
		"server_addr", cfg.Server.Addr,
		"server_env", cfg.Server.Env,
	)

	// 创建并启动 HTTP 服务器
	server := httpserver.NewServer(container.Router, cfg.Server.Addr)

	// 启动服务器（在goroutine中）
	go func() {
		slog.Info("Starting API server", "address", cfg.Server.Addr)
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down API server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("API server exited")
	return nil
}
