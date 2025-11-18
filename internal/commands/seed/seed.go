// Package seed 提供数据库种子命令
package seed

import (
	"context"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/database"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/database/seeds"
	"github.com/urfave/cli/v3"
)

// Command 定义种子命令
var Command = &cli.Command{
	Name:  "seed",
	Usage: "填充数据库种子数据",
	Description: `
   填充数据库种子数据，用于开发和测试环境。
   包含示例用户、演示数据等。
	`,
	Action: runSeed,
}

// getAllSeeders 返回所有可用的种子
func getAllSeeders() []database.Seeder {
	return []database.Seeder{
		&seeds.UserSeeder{},
		// 未来可以添加更多种子...
	}
}

// runSeed 执行种子数据填充
func runSeed(ctx context.Context, cmd *cli.Command) error {
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

	// 创建种子管理器
	manager := database.NewSeederManager(db, getAllSeeders())

	slog.Info("Running database seeders...")
	if err := manager.Run(ctx); err != nil {
		slog.Error("Seeding failed", "error", err)
		return err
	}

	slog.Info("Seeding completed successfully")
	return nil
}
