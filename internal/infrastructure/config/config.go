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
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Addr      string `koanf:"addr" comment:"监听地址，格式: host:port，例如 '0.0.0.0:8080' 或 ':8080'"`
	Env       string `koanf:"env" comment:"运行环境: development | production"`
	StaticDir string `koanf:"static_dir" comment:"静态资源目录路径，用于提供前端文件服务 (如 SPA 应用)"`
	DocsDir   string `koanf:"docs_dir" comment:"文档目录路径，用于提供 VitePress 构建的文档服务，通过 /docs 路由访问"`
}

// DataConfig 数据源配置
type DataConfig struct {
	PgsqlURL       string `koanf:"pgsql_url" comment:"PostgreSQL 连接 URL，格式: postgresql://user:password@host:port/dbname?sslmode=disable"`
	RedisURL       string `koanf:"redis_url" comment:"Redis 连接 URL，格式: redis://:password@host:port/db"`
	RedisKeyPrefix string `koanf:"redis_key_prefix" comment:"Redis key 前缀，所有 key 读写都会自动拼接此前缀，例如 'app:'"`
	AutoMigrate    bool   `koanf:"auto_migrate" comment:"是否在应用启动时自动执行数据库迁移 (仅推荐在开发环境使用，生产环境应使用 migrate 命令)"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret             string        `koanf:"secret" comment:"JWT 签名密钥 - ⚠️ 生产环境务必修改! 建议通过环境变量 APP_JWT_SECRET 设置"`
	AccessTokenExpiry  time.Duration `koanf:"access_token_expiry" comment:"访问令牌过期时间 (格式: 15m, 1h, 24h 等)"`
	RefreshTokenExpiry time.Duration `koanf:"refresh_token_expiry" comment:"刷新令牌过期时间 (168h = 7天)"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	DevSecret       string `koanf:"dev_secret" comment:"开发模式密钥 (用于验证码开发模式) - ⚠️ 生产环境务必修改! 建议通过环境变量 APP_AUTH_DEV_SECRET 设置"`
	TwoFAIssuer     string `koanf:"twofa_issuer" comment:"2FA TOTP 发行者名称，显示在用户的验证器应用中"`
	CaptchaRequired bool   `koanf:"captcha_required" comment:"是否需要验证码 (可在生产环境强制开启以提升安全性)"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `koanf:"level" comment:"日志级别: DEBUG | INFO | WARN | ERROR (大小写不敏感)"`
	Format     string `koanf:"format" comment:"输出格式: json | text | color (开发环境推荐 color，生产环境推荐 json)"`
	Output     string `koanf:"output" comment:"输出目标: stdout | stderr | 文件路径 (如 '/var/log/app.log')"`
	AddSource  bool   `koanf:"add_source" comment:"是否添加源代码位置信息 (调试时可开启，会显示文件名和行号)"`
	TimeFormat string `koanf:"time_format" comment:"时间格式: datetime | rfc3339 | rfc3339ms | unix | unixms (datetime 最易读)"`
	Timezone   string `koanf:"timezone" comment:"时区设置，例如 'Asia/Shanghai' 或 '+08:00'"`
}

// Config 应用配置
type Config struct {
	Server ServerConfig `koanf:"server" comment:"服务器配置"`
	Data   DataConfig   `koanf:"data" comment:"数据源配置"`
	JWT    JWTConfig    `koanf:"jwt" comment:"JWT 认证配置"`
	Auth   AuthConfig   `koanf:"auth" comment:"认证配置"`
	Log    LogConfig    `koanf:"log" comment:"日志配置"`
}

// defaultConfig 返回默认配置
func defaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Addr:      "0.0.0.0:8080",
			Env:       "development",
			StaticDir: "web/dist",
			DocsDir:   "docs/.vitepress/dist",
		},
		Data: DataConfig{
			PgsqlURL:       "postgresql://postgres@localhost:5432/app?sslmode=disable",
			RedisURL:       "redis://localhost:6379/0",
			RedisKeyPrefix: "app:",
			AutoMigrate:    false, // 默认关闭自动迁移，生产环境使用 migrate 命令
		},
		JWT: JWTConfig{
			Secret:             "change-me-in-production",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		},
		Auth: AuthConfig{
			DevSecret:       "dev-secret-change-me",
			TwoFAIssuer:     "Go-DDD-Template",
			CaptchaRequired: true, // 默认开启验证码
		},
		Log: LogConfig{
			Level:      "debug",         // 开发环境默认调试级别
			Format:     "color",         // 开发环境默认使用彩色输出
			Output:     "stdout",        // 输出到标准输出
			AddSource:  false,           // 默认不添加源代码位置信息
			TimeFormat: "datetime",      // 默认使用易读的时间格式
			Timezone:   "Asia/Shanghai", // 默认使用上海时区
		},
	}
}

// Load 加载配置，按优先级合并：
// 1. 默认值 (最低优先级)
// 2. 配置文件 (config.yaml)
// 3. 环境变量 (APP_*，最高优先级)
func Load() (*Config, error) {
	k := koanf.New(".")

	// 1️⃣ 加载默认配置 (最低优先级)
	if err := k.Load(structs.Provider(defaultConfig(), "koanf"), nil); err != nil {
		return nil, fmt.Errorf("failed to load default config: %w", err)
	}

	// 2️⃣ 加载 YAML 配置文件 (可选，文件不存在不报错)
	// 按优先级搜索：当前目录 -> config 目录
	configPaths := []string{"config.yaml", "config/config.yaml"}
	configLoaded := false
	for _, path := range configPaths {
		if err := k.Load(file.Provider(path), yaml.Parser()); err == nil {
			slog.Info("Loaded config from file", "path", path)
			configLoaded = true
			break
		}
	}
	if !configLoaded {
		slog.Info("No config file found, using defaults and env vars")
	}

	// 3️⃣ 加载环境变量 (APP_SERVER_ADDR → server.addr)
	if err := k.Load(env.Provider(".", env.Opt{
		Prefix: "APP_",
		TransformFunc: func(key, value string) (string, any) {
			// APP_SERVER_ADDR → server.addr
			// APP_DATA_PGSQL_URL → data.pgsql.url
			key = strings.ReplaceAll(
				strings.ToLower(strings.TrimPrefix(key, "APP_")),
				"_", ".",
			)
			return key, value
		},
	}), nil); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	// 解析到结构体
	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// LoadWithDefaults 加载配置的别名函数 (保持向后兼容)
func LoadWithDefaults() (*Config, error) {
	return Load()
}
