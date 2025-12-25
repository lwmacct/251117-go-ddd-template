package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/bootstrap"
	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/lwmacct/251219-go-pkg-logm/pkg/logm"

	"github.com/urfave/cli/v3"

	httpserver "github.com/lwmacct/251117-go-ddd-template/internal/adapters/http"
)

// action 执行 API 服务器启动逻辑
func action(ctx context.Context, cmd *cli.Command) error {
	// 加载配置 (优先级：默认值 -> 配置文件 -> 环境变量)
	logm.MustInit(logm.PresetAuto()...)
	cfg := cfgm.MustLoadCmd(cmd, config.DefaultConfig(), version.AppRawName)
	// 如果用户显式指定了 --static 参数，则覆盖配置
	if cmd.IsSet("static") {
		cfg.Server.WebDist = cmd.String("static")
	}

	// 初始化容器 (依赖注入) - 使用 DDD+CQRS 架构容器
	opts := &bootstrap.ContainerOptions{
		AutoMigrate: cfg.Data.AutoMigrate, // 从配置读取是否自动迁移
	}
	container, err := bootstrap.NewContainer(ctx, cfg, opts)
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
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down API server...")

	// 使用传入的 ctx 派生 shutdown context，符合 contextcheck 规范
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 关闭日志系统（确保文件正确关闭）
	if err := logm.Close(); err != nil {
		return err
	}

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		return err
	}

	slog.Info("API server exited")
	return nil
}
