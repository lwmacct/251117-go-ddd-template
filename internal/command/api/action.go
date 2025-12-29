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

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/container"
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

	slog.Info("Configuration loaded",
		"server_addr", cfg.Server.Addr,
		"server_env", cfg.Server.Env,
	)

	// 初始化选项
	opts := &container.ContainerOptions{
		AutoMigrate: cfg.Data.AutoMigrate,
	}

	// 用于接收 Router 的变量
	var router *gin.Engine

	// 创建 Fx 应用
	app := fx.New(
		fx.NopLogger,

		// 提供配置
		fx.Supply(cfg),
		fx.Supply(opts),

		// 核心模块
		container.InfraModule,
		container.CacheModule,
		container.RepositoryModule,
		container.ServiceModule,
		container.UseCaseModule,
		container.HTTPModule,

		// 启动钩子
		fx.Invoke(
			container.RegisterEventHandlers,
		),

		// 提取 Router
		fx.Populate(&router),

		fx.StartTimeout(30*time.Second),
		fx.StopTimeout(30*time.Second),
	)

	// 启动 Fx 应用
	startCtx, startCancel := context.WithTimeout(ctx, 30*time.Second)
	defer startCancel()

	if err := app.Start(startCtx); err != nil {
		slog.Error("Failed to start application", "error", err)
		return err
	}

	// 创建并启动 HTTP 服务器
	server := httpserver.NewServer(router, cfg.Server.Addr)

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

	// 关闭 HTTP 服务器
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	// 关闭 Fx 应用 (停止所有组件)
	// 使用独立 context：原 ctx 可能已取消，需要新超时
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer stopCancel()

	if err := app.Stop(stopCtx); err != nil { //nolint:contextcheck // 关闭时需要独立超时
		slog.Error("Failed to stop application", "error", err)
		return err
	}

	// 关闭日志系统
	if err := logm.Close(); err != nil {
		return err
	}

	slog.Info("API server exited")
	return nil
}
