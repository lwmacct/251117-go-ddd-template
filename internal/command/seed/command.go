// Package seed 提供数据库种子命令
package seed

import (
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
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
	Action:   action,
	Commands: []*cli.Command{version.Command},
}
