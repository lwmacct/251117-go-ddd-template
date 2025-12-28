package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/bootstrap"
	"github.com/lwmacct/251117-go-ddd-template/internal/config"
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
	manager := database.NewMigrationManager(db, bootstrap.GetAllModels())
	if err := manager.ResetWithIndexes(getIndexMigrations()); err != nil {
		slog.Error("Database reset failed", "error", err)
		return err
	}
	slog.Info("Database schema recreated successfully")

	// 2. 填充种子数据（除非 --empty）
	if !cmd.Bool("empty") {
		slog.Info("Running database seeders...")
		seederManager := database.NewSeederManager(db, seeds.DefaultSeeders())
		if err := seederManager.Run(ctx); err != nil {
			slog.Error("Seeding failed", "error", err)
			return err
		}
		slog.Info("Seed data populated successfully")
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
