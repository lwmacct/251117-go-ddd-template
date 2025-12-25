package bootstrap

import (
	"context"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/database"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/eventbus"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/redis"
)

// newInfrastructureModule 初始化基础设施模块
// 按顺序初始化：Database → Redis → EventBus
func newInfrastructureModule(ctx context.Context, cfg *config.Config, opts *ContainerOptions) (*InfrastructureModule, error) {
	m := &InfrastructureModule{}

	// 1. 初始化数据库连接
	dbConfig := database.DefaultConfig(cfg.Data.PgsqlURL)
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

	// 3. 初始化 Redis 客户端
	redisClient, err := redis.NewClient(ctx, cfg.Data.RedisURL)
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

	return nil
}
