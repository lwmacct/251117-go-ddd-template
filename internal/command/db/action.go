package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/container"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/cache"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/database"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/database/seeds"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/urfave/cli/v3"
)

// getIndexMigrations 返回需要创建的索引配置
func getIndexMigrations() []database.IndexMigration {
	return []database.IndexMigration{
		{
			Model:   &persistence.SettingModel{},
			Indexes: []string{"idx_settings_category_sort"},
		},
	}
}

// getJoinTableIndexes 返回 many2many 关联表需要创建的索引
// GORM AutoMigrate 只创建复合主键，不会为外键列创建单独索引
func getJoinTableIndexes() []database.JoinTableIndex {
	return []database.JoinTableIndex{
		{Table: "user_roles", Name: "idx_user_roles_user_id", Columns: "user_id"},
		{Table: "user_roles", Name: "idx_user_roles_role_id", Columns: "role_id"},
	}
}

// flushRedisCache 清空 Redis 缓存
// 数据库操作后需要清空缓存以避免数据不一致
func flushRedisCache(ctx context.Context, cfg *config.Config) error {
	slog.Info("Flushing Redis cache...")

	redisClient, err := cache.NewClient(ctx, cfg.Data.RedisURL, false)
	if err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		return err
	}
	defer func() {
		if closeErr := cache.Close(redisClient); closeErr != nil {
			slog.Error("Failed to close Redis connection", "error", closeErr)
		}
	}()

	// 使用 FLUSHDB 清空当前数据库
	if err := redisClient.FlushDB(ctx).Err(); err != nil {
		slog.Error("Failed to flush Redis cache", "error", err)
		return err
	}

	slog.Info("Redis cache flushed successfully")
	return nil
}

// actionMigrate 执行数据库迁移（只添加，不删除）
func actionMigrate(ctx context.Context, cmd *cli.Command) error {
	cfg := cfgm.MustLoadCmd(cmd, config.DefaultConfig(), version.AppRawName)

	// 初始化数据库连接
	dbConfig := database.DefaultConfig(cfg.Data.PgsqlURL)
	db, err := database.NewConnection(ctx, dbConfig)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return err
	}
	defer func() {
		if closeErr := database.Close(db); closeErr != nil {
			slog.Error("Failed to close database connection", "error", closeErr)
		}
	}()

	// 1. 执行 AutoMigrate
	slog.Info("Running AutoMigrate...")
	if err := db.AutoMigrate(container.GetAllModels()...); err != nil {
		slog.Error("AutoMigrate failed", "error", err)
		return err
	}
	slog.Info("AutoMigrate completed")

	// 2. 创建 Model 索引
	slog.Info("Creating model indexes...")
	for _, im := range getIndexMigrations() {
		if err := database.CreateIndexes(db, im.Model, im.Indexes); err != nil {
			slog.Error("Failed to create model indexes", "error", err)
			return err
		}
	}

	// 3. 创建关联表索引
	slog.Info("Creating join table indexes...")
	if err := database.CreateJoinTableIndexes(db, getJoinTableIndexes()); err != nil {
		slog.Error("Failed to create join table indexes", "error", err)
		return err
	}

	// 4. 清空 Redis 缓存
	if err := flushRedisCache(ctx, cfg); err != nil {
		return err
	}

	slog.Info("Database migration completed successfully")
	return nil
}

// actionReset 重置数据库（删表 + 重建 + 可选填充种子数据）
func actionReset(ctx context.Context, cmd *cli.Command) error {
	cfg := cfgm.MustLoadCmd(cmd, config.DefaultConfig(), version.AppRawName)

	// 检查是否为生产环境
	if cfg.Server.Env == "production" && !cmd.Bool("force") {
		slog.Error("Cannot reset database in production environment without --force flag")
		return errors.New("database reset is dangerous in production")
	}

	// 如果没有 --force 标志，需要用户确认
	if !cmd.Bool("force") {
		if canceled, err := confirmReset(cmd.Bool("empty")); err != nil {
			return err
		} else if canceled {
			return nil
		}
	}

	// 初始化数据库连接
	dbConfig := database.DefaultConfig(cfg.Data.PgsqlURL)
	db, err := database.NewConnection(ctx, dbConfig)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return err
	}
	defer func() {
		if closeErr := database.Close(db); closeErr != nil {
			slog.Error("Failed to close database connection", "error", closeErr)
		}
	}()

	// 1. 删表 + 重建
	slog.Info("Dropping all tables...")
	manager := database.NewMigrationManager(db, container.GetAllModels())
	if err := manager.ResetWithIndexes(getIndexMigrations()); err != nil {
		slog.Error("Database reset failed", "error", err)
		return err
	}
	slog.Info("Database schema recreated successfully")

	// 2. 创建关联表索引
	slog.Info("Creating join table indexes...")
	if err := database.CreateJoinTableIndexes(db, getJoinTableIndexes()); err != nil {
		slog.Error("Failed to create join table indexes", "error", err)
		return err
	}

	// 3. 填充种子数据（除非 --empty）
	if !cmd.Bool("empty") {
		slog.Info("Running database seeders...")
		seederManager := database.NewSeederManager(db, seeds.DefaultSeeders())
		if err := seederManager.Run(ctx); err != nil {
			slog.Error("Seeding failed", "error", err)
			return err
		}
		slog.Info("Seed data populated successfully")
	}

	// 4. 清空 Redis 缓存
	if err := flushRedisCache(ctx, cfg); err != nil {
		return err
	}

	slog.Info("Database reset completed")
	return nil
}

// confirmReset 显示确认提示并获取用户输入
// 返回 (canceled, error)
func confirmReset(empty bool) (bool, error) {
	//nolint:forbidigo // CLI 用户交互输出
	fmt.Println("\n⚠️  WARNING: This will delete ALL data in the database!")
	if !empty {
		//nolint:forbidigo // CLI 用户交互输出
		fmt.Println("   After reset, seed data will be populated.")
	}
	//nolint:forbidigo // CLI 用户交互输出
	fmt.Print("Are you sure you want to continue? (yes/no): ")

	var confirm string
	if _, err := fmt.Scanln(&confirm); err != nil {
		slog.Error("Failed to read input", "error", err)
		return false, err
	}

	if confirm != "yes" {
		slog.Info("Database reset canceled")
		return true, nil
	}

	return false, nil
}
