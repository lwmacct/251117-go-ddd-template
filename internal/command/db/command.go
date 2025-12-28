// Package db 提供数据库管理命令
package db

import (
	"github.com/urfave/cli/v3"
)

// Command 定义数据库管理命令
var Command = &cli.Command{
	Name:  "db",
	Usage: "数据库管理",
	Description: `
   管理数据库：重置表结构和填充种子数据。

   示例：
   - db reset         重置数据库并填充种子数据
   - db reset --empty 只重置表结构，不填充数据
   - db reset --force 跳过确认提示
	`,
	Commands: []*cli.Command{
		{
			Name:  "reset",
			Usage: "重置数据库（删表 + 重建 + 填充种子数据）",
			Description: `删除所有表并重新创建，默认填充种子数据。

   警告：此操作会删除所有数据，仅适用于开发环境！`,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "empty",
					Usage: "不填充种子数据",
				},
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "跳过确认提示",
				},
			},
			Action: actionReset,
		},
	},
}
