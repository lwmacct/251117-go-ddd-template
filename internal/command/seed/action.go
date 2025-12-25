package seed

import (
	"context"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/database"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/database/seeds"
	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/urfave/cli/v3"
)

// action 执行种子数据填充
func action(ctx context.Context, cmd *cli.Command) error {
	// cfg := config.LoadCmd(cmd)
	cfg := cfgm.MustLoadCmd(cmd, config.DefaultConfig(), version.AppRawName)

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
	manager := database.NewSeederManager(db, seeds.DefaultSeeders())

	slog.Info("Running database seeders...")
	if err := manager.Run(ctx); err != nil {
		slog.Error("Seeding failed", "error", err)
		return err
	}

	slog.Info("Seeding completed successfully")
	return nil
}
