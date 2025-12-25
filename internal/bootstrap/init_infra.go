package bootstrap

import (
	"context"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/database"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/eventbus"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/redis"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/telemetry"
)

// newInfrastructureModule 初始化基础设施模块
// 按顺序初始化：Telemetry → Database → Redis → EventBus
func newInfrastructureModule(ctx context.Context, cfg *config.Config, opts *ContainerOptions) (*InfrastructureModule, error) {
	m := &InfrastructureModule{}

	// 0. 初始化 OpenTelemetry（最先初始化，以便后续组件使用）
	telemetryShutdown, err := telemetry.InitTracer(ctx, telemetry.Config{
		ServiceName:    "go-ddd-template",
		ServiceVersion: "1.0.0",
		Environment:    cfg.Server.Env,
		Enabled:        cfg.Telemetry.Enabled,
		ExporterType:   cfg.Telemetry.ExporterType,
		OTLPEndpoint:   cfg.Telemetry.OTLPEndpoint,
		SampleRate:     cfg.Telemetry.SampleRate,
	})
	if err != nil {
		return nil, err
	}
	m.TelemetryShutdown = telemetryShutdown
	if cfg.Telemetry.Enabled {
		slog.Info("OpenTelemetry tracing initialized",
			"exporter", cfg.Telemetry.ExporterType,
			"sample_rate", cfg.Telemetry.SampleRate)
	}

	// 1. 初始化数据库连接
	dbConfig := database.DefaultConfig(cfg.Data.PgsqlURL)
	dbConfig.EnableTracing = cfg.Telemetry.Enabled // 联动 telemetry 配置
	db, err := database.NewConnection(ctx, dbConfig)
	if err != nil {
		return nil, err
	}
	m.DB = db

	// 2. 条件性执行自动迁移
	if opts.AutoMigrate {
		slog.Info("Auto-migration enabled, migrating database...")
		migrator := database.NewMigrator(db)
		if err = migrator.AutoMigrate(GetAllModels()...); err != nil {
			return nil, err
		}
		slog.Info("Database migration completed")
	} else {
		slog.Info("Auto-migration disabled, skipping database migration")
	}

	// 3. 初始化 Redis 客户端（与 Telemetry 配置联动）
	redisClient, err := redis.NewClient(ctx, cfg.Data.RedisURL, cfg.Telemetry.Enabled)
	if err != nil {
		return nil, err
	}
	m.RedisClient = redisClient

	// 4. 初始化事件总线
	m.EventBus = eventbus.NewInMemoryEventBus()
	slog.Info("Event bus initialized", "type", "InMemoryEventBus")

	return m, nil
}

// Close 关闭基础设施模块的所有资源
// 关闭顺序与初始化顺序相反
func (m *InfrastructureModule) Close() error {
	// 关闭事件总线
	if m.EventBus != nil {
		if err := m.EventBus.Close(); err != nil {
			slog.Error("failed to close event bus", "error", err)
		}
	}

	// 关闭 Redis 连接
	if err := redis.Close(m.RedisClient); err != nil {
		return err
	}

	// 关闭数据库连接
	if err := database.Close(m.DB); err != nil {
		return err
	}

	// 关闭 OpenTelemetry（最后关闭，确保其他组件的 span 都已导出）
	if m.TelemetryShutdown != nil {
		if err := m.TelemetryShutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown telemetry", "error", err)
		} else {
			slog.Info("OpenTelemetry shutdown completed")
		}
	}

	return nil
}
