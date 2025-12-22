package migrate

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/bootstrap"
	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/database"
	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/urfave/cli/v3"
)

// actionUp 执行向上迁移
func actionUp(ctx context.Context, cmd *cli.Command) error {
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

	// 创建迁移管理器
	manager := database.NewMigrationManager(db, bootstrap.GetAllModels())

	slog.Info("Running database migration...")
	if err := manager.Up(); err != nil {
		slog.Error("Migration failed", "error", err)
		return err
	}

	slog.Info("Migration completed successfully")
	return nil
}

// actionStatus 查看迁移状态
func actionStatus(ctx context.Context, cmd *cli.Command) error {
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

	// 创建迁移管理器
	manager := database.NewMigrationManager(db, bootstrap.GetAllModels())

	migrations, err := manager.Status()
	if err != nil {
		slog.Error("Failed to get migration status", "error", err)
		return err
	}

	if len(migrations) == 0 {
		slog.Info("No migrations found")
		return nil
	}

	slog.Info("Migration history:")
	//nolint:forbidigo // CLI 格式化输出，使用 fmt 是合理的
	fmt.Println("\n  ID | Version        | Name          | Applied At")
	//nolint:forbidigo // CLI 格式化输出
	fmt.Println("  ---|----------------|---------------|----------------------------")
	for _, m := range migrations {
		//nolint:forbidigo // CLI 格式化输出
		fmt.Printf("  %-3d| %-14s | %-13s | %s\n",
			m.ID, m.Version, m.Name, m.AppliedAt.Format("2006-01-02 15:04:05"))
	}
	//nolint:forbidigo // CLI 格式化输出
	fmt.Println()

	return nil
}

// actionFresh 删除所有表并重新迁移
func actionFresh(ctx context.Context, cmd *cli.Command) error {
	cfg := cfgm.MustLoadCmd(cmd, config.DefaultConfig(), version.AppRawName)

	// 检查是否为生产环境
	if cfg.Server.Env == "production" && !cmd.Bool("force") {
		slog.Error("Cannot run fresh migration in production environment without --force flag")
		return errors.New("fresh migration is dangerous in production")
	}

	// 如果没有 --force 标志，需要用户确认
	if !cmd.Bool("force") {
		//nolint:forbidigo // CLI 用户交互输出
		fmt.Println("\n⚠️  WARNING: This will delete ALL data in the database!")
		//nolint:forbidigo // CLI 用户交互输出
		fmt.Print("Are you sure you want to continue? (yes/no): ")
		var confirm string
		if _, err := fmt.Scanln(&confirm); err != nil {
			slog.Error("Failed to read input", "error", err)
			return err
		}
		if confirm != "yes" {
			slog.Info("Migration canceled")
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

	// 创建迁移管理器
	manager := database.NewMigrationManager(db, bootstrap.GetAllModels())

	slog.Info("Dropping all tables...")
	if err := manager.Fresh(); err != nil {
		slog.Error("Fresh migration failed", "error", err)
		return err
	}

	slog.Info("Fresh migration completed successfully")
	return nil
}
