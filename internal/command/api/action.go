package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/bootstrap"
	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251125-go-pkg-logger/pkg/logger"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/urfave/cli/v3"

	httpserver "github.com/lwmacct/251117-go-ddd-template/internal/adapters/http"

	pkgconfig "github.com/lwmacct/251207-go-pkg-config/pkg/config"
)

// action 执行 API 服务器启动逻辑
func action(ctx context.Context, cmd *cli.Command) error {
	// 加载配置 (优先级：默认值 -> 配置文件 -> 环境变量)
	if err := logger.InitEnv(); err != nil {
		slog.Warn("初始化日志系统失败，使用默认配置", "error", err)
	}

	cfg, err := config.Load(cmd, pkgconfig.DefaultPaths(version.GetAppRawName()))
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
		cfg.Server.DistWeb = cmd.String("static")
	}

	// 初始化容器 (依赖注入) - 使用 DDD+CQRS 架构容器
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

	// 启动服务器 (在goroutine中)
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

	// 关闭日志系统（确保文件正确关闭）
	if err := logger.Close(); err != nil {
		return err
	}

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("API server exited")
	return nil
}
