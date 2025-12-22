// Package api provides the HTTP API server command.
//
// @title           Go DDD Template API
// @version         1.0.3
// @description     基于 DDD+CQRS 架构的 Go Web 应用 API 文档
// @termsOfService  https://github.com/lwmacct/251117-go-ddd-template
//
// @contact.name   API Support
// @contact.url    https://github.com/lwmacct/251117-go-ddd-template/issues
// @contact.email  lwmacct@icloud.com
//
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
//
// @host      localhost:40012
// @BasePath  /api
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token authentication. Format: "Bearer {token}"
//
// @externalDocs.description  GitHub Repository
// @externalDocs.url          https://github.com/lwmacct/251117-go-ddd-template
package api

import (
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/urfave/cli/v3"
)

// Command 定义 API 服务命令
var Command = &cli.Command{
	Name:    "api",
	Aliases: []string{"serve", "server"},
	Usage:   "启动 REST API 服务",
	Description: `
   启动 HTTP API 服务器，提供 RESTful API 接口。
   服务器支持优雅关闭，会等待正在处理的请求完成后再退出。

   配置优先级 (从低到高) ：
   1. 默认值
   2. 配置文件 (config.yaml)
   3. 环境变量 (APP_SERVER_ADDR)
   4. 命令行参数 (--addr)
	`,
	Action:   action,
	Commands: []*cli.Command{version.Command},

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "addr",
			Aliases: []string{"a"},
			Usage:   "服务器监听地址 (格式: host:port)",
		},
		&cli.StringFlag{
			Name:    "static",
			Aliases: []string{"s"},
			Usage:   "静态资源目录路径",
		},
	},
}
