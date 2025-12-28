// Package container provides Fx module definitions for dependency injection.
//
// Each module groups related providers and lifecycle hooks:
//   - [InfraModule]: Database, Redis, EventBus, Telemetry
//   - [CacheModule]: All cache services
//   - [RepositoryModule]: CQRS repositories with cache decorators
//   - [ServiceModule]: Domain and infrastructure services
//   - [UseCaseModule]: Application use case handlers
//   - [HTTPModule]: HTTP handlers and router
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

// InfraModule provides infrastructure components.
//
// Components:
//   - OpenTelemetry tracing
//   - PostgreSQL database connection
//   - Redis client
//   - In-memory event bus
//
// Lifecycle:
//   - OnStart: Initialize telemetry, connect DB, connect Redis
//   - OnStop: Close connections in reverse order
var InfraModule = fx.Module("infra",
	fx.Provide(
		newTelemetry,
		newDatabase,
		newRedisClient,
		newEventBus,
	),
)

// TelemetryResult wraps telemetry shutdown function for Fx lifecycle.
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

	// Auto-migrate if enabled
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

// runAutoMigrate 执行数据库自动迁移和索引创建
func runAutoMigrate(db *gorm.DB) error {
	slog.Info("Auto-migration enabled, migrating database...")

	if err := db.AutoMigrate(GetAllModels()...); err != nil {
		return err
	}

	// Create indexes for SettingModel
	if err := database.CreateIndexes(db, &persistence.SettingModel{}, []string{
		"idx_settings_category_sort",
	}); err != nil {
		return err
	}

	// Create indexes for many2many join tables
	if err := database.CreateJoinTableIndexes(db, []database.JoinTableIndex{
		{Table: "user_roles", Name: "idx_user_roles_user_id", Columns: "user_id"},
		{Table: "user_roles", Name: "idx_user_roles_role_id", Columns: "role_id"},
		{Table: "role_permissions", Name: "idx_role_permissions_role_model_id", Columns: "role_model_id"},
		{Table: "role_permissions", Name: "idx_role_permissions_permission_model_id", Columns: "permission_model_id"},
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
