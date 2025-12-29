// Package container 提供 Fx 模块定义，用于依赖注入。
//
// 每个模块聚合相关的 Provider 和生命周期钩子：
//   - [InfraModule]: 数据库、Redis、EventBus、链路追踪
//   - [CacheModule]: 所有缓存服务
//   - [RepositoryModule]: 带缓存装饰器的 CQRS 仓储
//   - [ServiceModule]: 领域服务和基础设施服务
//   - [UseCaseModule]: 应用层用例处理器
//   - [HTTPModule]: HTTP 处理器和路由
package container

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/database"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/eventbus"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/telemetry"
)

// InfraModule 提供基础设施组件。
//
// 组件：
//   - OpenTelemetry 链路追踪
//   - PostgreSQL 数据库连接
//   - Redis 客户端
//   - 内存事件总线
//
// 生命周期：
//   - OnStart: 初始化链路追踪、连接数据库、连接 Redis
//   - OnStop: 按相反顺序关闭连接
var InfraModule = fx.Module("infra",
	fx.Provide(
		newTelemetry,
		newDatabase,
		newRedisClient,
		newEventBus,
	),
)

// TelemetryResult 包装链路追踪的关闭函数，用于 Fx 生命周期管理。
type TelemetryResult struct {
	fx.Out

	Shutdown telemetry.ShutdownFunc
}

func newTelemetry(lc fx.Lifecycle, cfg *config.Config) (TelemetryResult, error) {
	ctx := context.Background()
	shutdown, err := telemetry.InitTracer(ctx, telemetry.Config{
		ServiceName:    "go-ddd-template",
		ServiceVersion: "1.0.0",
		Environment:    cfg.Server.Env,
		Enabled:        cfg.Telemetry.Enabled,
		ExporterType:   cfg.Telemetry.ExporterType,
		OTLPEndpoint:   cfg.Telemetry.OTLPEndpoint,
		SampleRate:     cfg.Telemetry.SampleRate,
	})
	if err != nil {
		return TelemetryResult{}, err
	}

	if cfg.Telemetry.Enabled {
		slog.Info("OpenTelemetry tracing initialized",
			"exporter", cfg.Telemetry.ExporterType,
			"sample_rate", cfg.Telemetry.SampleRate)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if shutdown != nil {
				if err := shutdown(ctx); err != nil {
					slog.Error("failed to shutdown telemetry", "error", err)
					return err
				}
				slog.Info("OpenTelemetry shutdown completed")
			}
			return nil
		},
	})

	return TelemetryResult{Shutdown: shutdown}, nil
}

func newDatabase(lc fx.Lifecycle, cfg *config.Config, opts *ContainerOptions) (*gorm.DB, error) {
	ctx := context.Background()
	dbConfig := database.DefaultConfig(cfg.Data.PgsqlURL)
	dbConfig.EnableTracing = cfg.Telemetry.Enabled

	db, err := database.NewConnection(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	// 如果启用了自动迁移则执行
	if opts.AutoMigrate {
		if err := runAutoMigrate(db); err != nil {
			return nil, err
		}
	} else {
		slog.Info("Auto-migration disabled, skipping database migration")
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return database.Close(db)
		},
	})

	return db, nil
}

// runAutoMigrate 执行数据库自动迁移和索引创建。
func runAutoMigrate(db *gorm.DB) error {
	slog.Info("Auto-migration enabled, migrating database...")

	if err := db.AutoMigrate(GetAllModels()...); err != nil {
		return err
	}

	// 为 SettingModel 创建索引
	if err := database.CreateIndexes(db, &persistence.SettingModel{}, []string{
		"idx_settings_category_sort",
	}); err != nil {
		return err
	}

	// 为多对多关联表创建索引
	// role_permissions 使用复合主键，PostgreSQL 自动利用前缀索引
	if err := database.CreateJoinTableIndexes(db, []database.JoinTableIndex{
		{Table: "user_roles", Name: "idx_user_roles_user_id", Columns: "user_id"},
		{Table: "user_roles", Name: "idx_user_roles_role_id", Columns: "role_id"},
	}); err != nil {
		return err
	}

	slog.Info("Database migration completed")
	return nil
}

func newRedisClient(lc fx.Lifecycle, cfg *config.Config) (*redis.Client, error) {
	ctx := context.Background()
	client, err := cache.NewClient(ctx, cfg.Data.RedisURL, cfg.Telemetry.Enabled)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return cache.Close(client)
		},
	})

	return client, nil
}

func newEventBus(lc fx.Lifecycle) event.EventBus {
	bus := eventbus.NewInMemoryEventBus()
	slog.Info("Event bus initialized", "type", "InMemoryEventBus")

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return bus.Close()
		},
	})

	return bus
}
