// Package migrate 提供数据库迁移命令
package migrate

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/bootstrap"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/database"
	"github.com/urfave/cli/v3"
)

// Command 定义迁移命令
var Command = &cli.Command{
	Name:  "migrate",
	Usage: "数据库迁移管理",
	Description: `
   管理数据库迁移，包括执行、回滚、查看状态等操作。

   子命令：
   - up     执行数据库迁移
   - status 查看迁移状态
   - fresh  删除所有表并重新迁移 (危险！仅开发环境使用) 
	`,
	Commands: []*cli.Command{
		upCommand,
		statusCommand,
		freshCommand,
	},
}

var upCommand = &cli.Command{
	Name:  "up",
	Usage: "执行数据库迁移",
	Description: `
   执行数据库迁移，创建或更新所有表结构。
   该命令会自动创建迁移记录表，并记录每次迁移的版本和时间。
	`,
	Action: runUp,
}

var statusCommand = &cli.Command{
	Name:  "status",
	Usage: "查看迁移状态",
	Description: `
   查看已执行的迁移记录，包括版本号、名称和执行时间。
	`,
	Action: runStatus,
}

var freshCommand = &cli.Command{
	Name:  "fresh",
	Usage: "删除所有表并重新迁移 (危险操作！) ",
	Description: `
   删除所有表并重新执行迁移。
   警告：此操作会删除所有数据，仅适用于开发环境！
	`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "force",
			Usage: "强制执行 (不询问确认) ",
		},
	},
	Action: runFresh,
}

// runUp 执行向上迁移
func runUp(ctx context.Context, cmd *cli.Command) error {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return err
	}

	// 初始化数据库连接
	dbConfig := database.DefaultConfig(cfg.Data.PgsqlURL)
	db, err := database.NewConnection(ctx, dbConfig)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return err
	}
	defer func() {
		if err := database.Close(db); err != nil {
			slog.Error("Failed to close database connection", "error", err)
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

// runStatus 查看迁移状态
func runStatus(ctx context.Context, cmd *cli.Command) error {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return err
	}

	// 初始化数据库连接
	dbConfig := database.DefaultConfig(cfg.Data.PgsqlURL)
	db, err := database.NewConnection(ctx, dbConfig)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return err
	}
	defer func() {
		if err := database.Close(db); err != nil {
			slog.Error("Failed to close database connection", "error", err)
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
	fmt.Println("\n  ID | Version        | Name          | Applied At")
	fmt.Println("  ---|----------------|---------------|----------------------------")
	for _, m := range migrations {
		fmt.Printf("  %-3d| %-14s | %-13s | %s\n",
			m.ID, m.Version, m.Name, m.AppliedAt.Format("2006-01-02 15:04:05"))
	}
	fmt.Println()

	return nil
}

// runFresh 删除所有表并重新迁移
func runFresh(ctx context.Context, cmd *cli.Command) error {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return err
	}

	// 检查是否为生产环境
	if cfg.Server.Env == "production" && !cmd.Bool("force") {
		slog.Error("Cannot run fresh migration in production environment without --force flag")
		return fmt.Errorf("fresh migration is dangerous in production")
	}

	// 如果没有 --force 标志，需要用户确认
	if !cmd.Bool("force") {
		fmt.Println("\n⚠️  WARNING: This will delete ALL data in the database!")
		fmt.Print("Are you sure you want to continue? (yes/no): ")
		var confirm string
		if _, err := fmt.Scanln(&confirm); err != nil {
			slog.Error("Failed to read input", "error", err)
			return err
		}
		if confirm != "yes" {
			slog.Info("Migration cancelled")
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
		if err := database.Close(db); err != nil {
			slog.Error("Failed to close database connection", "error", err)
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
