package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// TestGenerateExample 生成 config/config.example.yaml 配置示例文件
//
// 此测试会从 defaultConfig() 提取默认配置值，并生成符合规范的 YAML 配置示例文件。
// 设计为可被 pre-commit hook 调用，在 config.go 变更时自动执行。
//
// 运行方式:
//
//	go test -v -run TestGenerateExample ./internal/infrastructure/config/...
func TestGenerateExample(t *testing.T) {
	// 获取项目根目录
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("无法找到项目根目录: %v", err)
	}

	// 获取默认配置
	cfg := defaultConfig()

	// 生成 YAML 内容
	var buf bytes.Buffer
	writeConfigYAML(&buf, cfg)

	// 写入文件
	outputPath := filepath.Join(projectRoot, "config", "config.example.yaml")
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		t.Fatalf("写入配置文件失败: %v", err)
	}

	t.Logf("✅ 已生成配置示例文件: %s", outputPath)
}

// TestConfigKeysValid 验证 config.yaml 不包含 config.example.yaml 中不存在的配置项
//
// 此测试确保用户的配置文件不会有未知的配置项，防止因拼写错误或过时配置导致的问题。
// 如果 config.yaml 不存在，测试会跳过。
func TestConfigKeysValid(t *testing.T) {
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("无法找到项目根目录: %v", err)
	}

	configPath := filepath.Join(projectRoot, "config", "config.yaml")
	examplePath := filepath.Join(projectRoot, "config", "config.example.yaml")

	// 如果 config.yaml 不存在，跳过测试
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skip("config.yaml 不存在，跳过验证")
	}

	// 加载 config.example.yaml 获取有效的 keys
	exampleKeys, err := loadYAMLKeys(examplePath)
	if err != nil {
		t.Fatalf("无法加载 config.example.yaml: %v", err)
	}

	// 加载 config.yaml 获取用户配置的 keys
	configKeys, err := loadYAMLKeys(configPath)
	if err != nil {
		t.Fatalf("无法加载 config.yaml: %v", err)
	}

	// 检查 config.yaml 中是否有 config.example.yaml 不存在的 keys
	var invalidKeys []string
	for _, key := range configKeys {
		if !containsKey(exampleKeys, key) {
			invalidKeys = append(invalidKeys, key)
		}
	}

	if len(invalidKeys) > 0 {
		t.Errorf("config.yaml 包含以下无效配置项 (在 config.example.yaml 中不存在):\n")
		for _, key := range invalidKeys {
			t.Errorf("  - %s", key)
		}
		t.Errorf("\n请检查拼写或从 config.example.yaml 中确认有效的配置项")
	}
}

// loadYAMLKeys 加载 YAML 文件并返回所有配置键的扁平化列表
func loadYAMLKeys(path string) ([]string, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("加载文件失败: %w", err)
	}

	return k.Keys(), nil
}

// containsKey 检查 keys 列表中是否包含指定的 key
func containsKey(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

// TestStructTags 验证 Config 结构体字段与 koanf 标签一致性
func TestStructTags(t *testing.T) {
	cfg := defaultConfig()
	cfgType := reflect.TypeOf(cfg)

	// 验证所有字段都有 koanf 标签
	checkKoanfTags(t, cfgType, "Config")
}

// findProjectRoot 通过查找 go.mod 文件定位项目根目录
func findProjectRoot() (string, error) {
	// 获取当前测试文件所在目录
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("无法获取当前文件路径")
	}

	dir := filepath.Dir(filename)

	// 向上查找 go.mod 文件
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("未找到 go.mod 文件")
		}
		dir = parent
	}
}

// writeConfigYAML 将配置结构体转换为带注释的 YAML 格式
func writeConfigYAML(buf *bytes.Buffer, cfg Config) {
	// 写入文件头注释
	buf.WriteString(`# 示例配置文件
# 复制此文件为 config.yaml 并根据需要修改
#
# 此文件与 internal/infrastructure/config/config.go 中的 defaultConfig() 保持同步
# 所有配置项都可以通过环境变量覆盖 (环境变量前缀：APP_)
# 例如：APP_SERVER_ADDR=:8080 会覆盖 server.addr 的值

`)

	// 写入各配置段
	writeServerConfig(buf, cfg.Server)
	buf.WriteString("\n")
	writeDataConfig(buf, cfg.Data)
	buf.WriteString("\n")
	writeJWTConfig(buf, cfg.JWT)
	buf.WriteString("\n")
	writeAuthConfig(buf, cfg.Auth)
	buf.WriteString("\n")
	writeLogConfig(buf, cfg.Log)
}

// writeServerConfig 写入服务器配置段
func writeServerConfig(buf *bytes.Buffer, cfg ServerConfig) {
	buf.WriteString("# 服务器配置\n")
	buf.WriteString("server:\n")
	writeField(buf, "addr", cfg.Addr, `监听地址，格式: host:port，例如 "0.0.0.0:8080" 或 ":8080"`)
	writeField(buf, "env", cfg.Env, "运行环境: development | production")
	writeField(buf, "static_dir", cfg.StaticDir, "静态资源目录路径，用于提供前端文件服务 (如 SPA 应用)")
	writeField(buf, "docs_dir", cfg.DocsDir, "文档目录路径，用于提供 VitePress 构建的文档服务，通过 /docs 路由访问")
}

// writeDataConfig 写入数据源配置段
func writeDataConfig(buf *bytes.Buffer, cfg DataConfig) {
	buf.WriteString("# 数据源配置\n")
	buf.WriteString("data:\n")
	writeField(buf, "pgsql_url", cfg.PgsqlURL, "PostgreSQL 连接 URL，格式: postgresql://user:password@host:port/dbname?sslmode=disable")
	writeField(buf, "redis_url", cfg.RedisURL, "Redis 连接 URL，格式: redis://:password@host:port/db")
	writeField(buf, "redis_key_prefix", cfg.RedisKeyPrefix, `Redis key 前缀，所有 key 读写都会自动拼接此前缀，例如 "myapp:"`)
	writeFieldBool(buf, "auto_migrate", cfg.AutoMigrate, "是否在应用启动时自动执行数据库迁移 (仅推荐在开发环境使用，生产环境应使用 migrate 命令)")
}

// writeJWTConfig 写入 JWT 配置段
func writeJWTConfig(buf *bytes.Buffer, cfg JWTConfig) {
	buf.WriteString("# JWT 认证配置\n")
	buf.WriteString("jwt:\n")
	writeFieldSensitive(buf, "secret", cfg.Secret, "JWT 签名密钥 - ⚠️ 生产环境务必修改! 建议通过环境变量 APP_JWT_SECRET 设置")
	writeFieldDuration(buf, "access_token_expiry", cfg.AccessTokenExpiry, "访问令牌过期时间 (格式: 15m, 1h, 24h 等)")
	writeFieldDuration(buf, "refresh_token_expiry", cfg.RefreshTokenExpiry, "刷新令牌过期时间 (168h = 7天)")
}

// writeAuthConfig 写入认证配置段
func writeAuthConfig(buf *bytes.Buffer, cfg AuthConfig) {
	buf.WriteString("# 认证配置\n")
	buf.WriteString("auth:\n")
	writeFieldSensitive(buf, "dev_secret", cfg.DevSecret, "开发模式密钥 (用于验证码开发模式) - ⚠️ 生产环境务必修改! 建议通过环境变量 APP_AUTH_DEV_SECRET 设置")
	writeField(buf, "twofa_issuer", cfg.TwoFAIssuer, "2FA TOTP 发行者名称，显示在用户的验证器应用中")
	writeFieldBool(buf, "captcha_required", cfg.CaptchaRequired, "是否需要验证码 (可在生产环境强制开启以提升安全性)")
}

// writeLogConfig 写入日志配置段
func writeLogConfig(buf *bytes.Buffer, cfg LogConfig) {
	buf.WriteString("# 日志配置\n")
	buf.WriteString("log:\n")
	writeField(buf, "level", cfg.Level, "日志级别: DEBUG | INFO | WARN | ERROR (大小写不敏感)")
	writeField(buf, "format", cfg.Format, "输出格式: json | text | color (开发环境推荐 color，生产环境推荐 json)")
	writeField(buf, "output", cfg.Output, `输出目标: stdout | stderr | 文件路径 (如 "/var/log/app.log")`)
	writeFieldBool(buf, "add_source", cfg.AddSource, "是否添加源代码位置信息 (调试时可开启，会显示文件名和行号)")
	writeField(buf, "time_format", cfg.TimeFormat, "时间格式: datetime | rfc3339 | rfc3339ms | unix | unixms (datetime 最易读)")
	writeField(buf, "timezone", cfg.Timezone, `时区设置，例如 "Asia/Shanghai" 或 "+08:00"`)
}

// writeField 写入字符串字段
func writeField(buf *bytes.Buffer, key, value, comment string) {
	fmt.Fprintf(buf, "  %s: %q # %s\n", key, value, comment)
}

// writeFieldBool 写入布尔字段
func writeFieldBool(buf *bytes.Buffer, key string, value bool, comment string) {
	fmt.Fprintf(buf, "  %s: %t # %s\n", key, value, comment)
}

// writeFieldDuration 写入 Duration 字段（转换为字符串格式）
func writeFieldDuration(buf *bytes.Buffer, key string, value time.Duration, comment string) {
	fmt.Fprintf(buf, "  %s: %q # %s\n", key, formatDuration(value), comment)
}

// writeFieldSensitive 写入敏感字段（添加安全警告）
func writeFieldSensitive(buf *bytes.Buffer, key, value, comment string) {
	fmt.Fprintf(buf, "  %s: %q # %s\n", key, value, comment)
}

// formatDuration 将 time.Duration 转换为人类可读的字符串格式
func formatDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	// 尝试简化为最大单位
	if d%(24*time.Hour) == 0 {
		hours := d / time.Hour
		return fmt.Sprintf("%dh", hours)
	}
	if d%time.Hour == 0 {
		return fmt.Sprintf("%dh", d/time.Hour)
	}
	if d%time.Minute == 0 {
		return fmt.Sprintf("%dm", d/time.Minute)
	}
	if d%time.Second == 0 {
		return fmt.Sprintf("%ds", d/time.Second)
	}

	return d.String()
}

// checkKoanfTags 递归检查结构体字段的 koanf 标签
func checkKoanfTags(t *testing.T, typ reflect.Type, path string) {
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldPath := path + "." + field.Name

		// 检查 koanf 标签
		koanfTag := field.Tag.Get("koanf")
		if koanfTag == "" {
			t.Errorf("字段 %s 缺少 koanf 标签", fieldPath)
		}

		// 递归检查嵌套结构体
		if field.Type.Kind() == reflect.Struct && !strings.HasPrefix(field.Type.String(), "time.") {
			checkKoanfTags(t, field.Type, fieldPath)
		}
	}
}
