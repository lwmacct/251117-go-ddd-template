// Package migrate 提供数据库迁移命令
package migrate

import (
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
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
		version.Command,
		{
			Name:        "up",
			Usage:       "执行数据库迁移",
			Description: `执行数据库迁移，创建或更新所有表结构, 该命令会自动创建迁移记录表，并记录每次迁移的版本和时间。`,
			Action:      actionUp,
		},
		{
			Name:        "status",
			Usage:       "查看迁移状态",
			Description: `查看已执行的迁移记录，包括版本号、名称和执行时间。`,
			Action:      actionStatus,
		},
		{
			Name:        "fresh",
			Usage:       "删除所有表并重新迁移 (危险操作！) ",
			Description: `删除所有表并重新执行迁移, 警告：此操作会删除所有数据，仅适用于开发环境！`,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "force",
					Usage: "强制执行 (不询问确认) ",
				},
			},
			Action: actionFresh,
		},
	},
}
