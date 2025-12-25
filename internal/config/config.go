// Package config 提供应用配置管理
//
// 配置加载优先级 (从低到高) ：
// 1. 默认值 - defaultConfig() 函数中定义
// 2. 配置文件 - config.yaml 或 config/config.yaml
// 3. 环境变量 - 以 APP_ 为前缀 (如 APP_SERVER_ADDR)
//
// 重要提示：
// - 如果修改了 defaultConfig() 中的默认值，请同步更新 config/config.example.yaml 示例文件
// - 配置文件路径硬编码在 Load() 函数中：[]string{"config.yaml", "config/config.yaml"}
package config

import (
	"time"
)

// Server 服务器配置
type Server struct {
	Addr     string `koanf:"addr" comment:"监听地址，格式: host:port，例如 '0.0.0.0:8080' 或 ':8080'"`
	Env      string `koanf:"env" comment:"运行环境: development | production"`
	WebDist  string `koanf:"web-dist" comment:"静态资源目录路径，用于提供前端文件服务 (如 SPA 应用)"`
	DocsDist string `koanf:"docs-dist" comment:"文档目录路径，用于提供 VitePress 构建的文档服务，通过 /docs 路由访问"`
}

// Data 数据源配置
type Data struct {
	PgsqlURL       string `koanf:"pgsql_url" comment:"PostgreSQL 连接 URL，格式: postgresql://user:password@host:port/dbname?sslmode=disable"`
	RedisURL       string `koanf:"redis_url" comment:"Redis 连接 URL，格式: redis://:password@host:port/db"`
	RedisKeyPrefix string `koanf:"redis_key_prefix" comment:"Redis key 前缀，所有 key 读写都会自动拼接此前缀，例如 'app:'"`
	AutoMigrate    bool   `koanf:"auto_migrate" comment:"是否在应用启动时自动执行数据库迁移 (仅推荐在开发环境使用，生产环境应使用 migrate 命令)"`
}

// JWT JWT配置
type JWT struct {
	Secret             string        `koanf:"secret" comment:"JWT 签名密钥 - ⚠️ 生产环境务必修改! 建议通过环境变量 APP_JWT_SECRET 设置"`
	AccessTokenExpiry  time.Duration `koanf:"access_token_expiry" comment:"访问令牌过期时间 (格式: 15m, 1h, 24h 等)"`
	RefreshTokenExpiry time.Duration `koanf:"refresh_token_expiry" comment:"刷新令牌过期时间 (168h = 7天)"`
}

// Auth 认证配置
type Auth struct {
	DevSecret       string `koanf:"dev_secret" comment:"开发模式密钥 (用于验证码开发模式) - ⚠️ 生产环境务必修改! 建议通过环境变量 APP_AUTH_DEV_SECRET 设置"`
	TwoFAIssuer     string `koanf:"twofa_issuer" comment:"2FA TOTP 发行者名称，显示在用户的验证器应用中"`
	CaptchaRequired bool   `koanf:"captcha_required" comment:"是否需要验证码 (可在生产环境强制开启以提升安全性)"`
}

// Config 应用配置
type Config struct {
	Server Server `koanf:"server" comment:"服务器配置"`
	Data   Data   `koanf:"data" comment:"数据源配置"`
	JWT    JWT    `koanf:"jwt" comment:"JWT 认证配置"`
	Auth   Auth   `koanf:"auth" comment:"认证配置"`
}

// DefaultConfig 返回默认配置
// 注意：internal/command/command.go 中的 Defaults 变量引用此函数以实现单一配置来源。
func DefaultConfig() Config {
	return Config{
		Server: Server{
			Addr:     "0.0.0.0:8080",
			Env:      "development",
			WebDist:  "dist",
			DocsDist: "docs/.vitepress/dist",
		},
		Data: Data{
			PgsqlURL:       "postgresql://postgres@localhost:5432/app?sslmode=disable",
			RedisURL:       "redis://localhost:6379/0",
			RedisKeyPrefix: "app:",
			AutoMigrate:    false, // 默认关闭自动迁移，生产环境使用 migrate 命令
		},
		JWT: JWT{
			Secret:             "change-me-in-production",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		},
		Auth: Auth{
			DevSecret:       "dev-secret-change-me",
			TwoFAIssuer:     "Go-DDD-Template",
			CaptchaRequired: true, // 默认开启验证码
		},
	}
}
