# 配置系统

本项目使用 [Koanf](https://github.com/knadh/koanf) 作为配置管理库，提供灵活的多层配置支持。

## 配置优先级

配置按以下优先级加载（从低到高）：

1. **默认值** - 在代码中定义的默认配置
2. **配置文件** - `config.yaml` 或 `configs/config.yaml`
3. **环境变量** - 前缀为 `APP_` 的环境变量
4. **命令行参数** - CLI 标志

后面的配置会覆盖前面的配置。

## 配置文件

### 文件位置

配置文件会按以下顺序搜索：

1. 工作目录中的 `config.yaml`
2. `configs/config.yaml`

### 配置示例

创建 `config.yaml` 文件：

```yaml
server:
  addr: ":8080"
  env: "development"
  static_dir: "web/dist"

data:
  pgsql:
    url: "postgresql://postgres:postgres@localhost:5432/app?sslmode=disable"
  redis:
    url: "redis://localhost:6379/0"

jwt:
  secret: "your-secret-key-change-in-production"
  access_token_expiry: "15m"
  refresh_token_expiry: "168h"
```

## 环境变量

环境变量使用 `APP_` 前缀，使用下划线分隔层级。

### 格式规则

- 前缀：`APP_`
- 层级分隔：使用 `_`（下划线）
- 大小写：不敏感，但推荐全大写

### 示例

```bash
# 服务器配置
export APP_SERVER_ADDR=":8080"
export APP_SERVER_ENV="production"
export APP_SERVER_STATIC_DIR="web/dist"

# 数据库配置
export APP_DATA_PGSQL_URL="postgresql://user:pass@localhost:5432/db?sslmode=disable"

# Redis 配置
export APP_DATA_REDIS_URL="redis://localhost:6379/0"

# JWT 配置
export APP_JWT_SECRET="your-secret-key"
export APP_JWT_ACCESS_TOKEN_EXPIRY="15m"
export APP_JWT_REFRESH_TOKEN_EXPIRY="168h"
```

## 命令行参数

使用命令行参数覆盖配置：

```bash
# 指定监听地址
.local/bin/251117-go-ddd-template api --addr :9000

# 使用环境变量和命令行参数组合
APP_JWT_SECRET="secret" .local/bin/251117-go-ddd-template api --addr :9000
```

## 配置结构

完整的配置结构定义在 `internal/infrastructure/config/config.go`：

```go
type Config struct {
    Server ServerConfig `koanf:"server"`
    Data   DataConfig   `koanf:"data"`
    JWT    JWTConfig    `koanf:"jwt"`
}

type ServerConfig struct {
    Addr      string `koanf:"addr"`
    Env       string `koanf:"env"`
    StaticDir string `koanf:"static_dir"`
}

type DataConfig struct {
    PgSQL PgSQLConfig `koanf:"pgsql"`
    Redis RedisConfig `koanf:"redis"`
}

type PgSQLConfig struct {
    URL string `koanf:"url"`
}

type RedisConfig struct {
    URL string `koanf:"url"`
}

type JWTConfig struct {
    Secret              string `koanf:"secret"`
    AccessTokenExpiry   string `koanf:"access_token_expiry"`
    RefreshTokenExpiry  string `koanf:"refresh_token_expiry"`
}
```

## 默认配置

如果没有提供配置文件，系统会使用以下默认值：

```go
server:
  addr: ":8080"
  env: "development"
  static_dir: "web/dist"

data:
  pgsql:
    url: "postgresql://postgres:postgres@localhost:5432/app?sslmode=disable"
  redis:
    url: "redis://localhost:6379/0"

jwt:
  secret: "change-me-in-production"
  access_token_expiry: "15m"
  refresh_token_expiry: "168h"
```

## 配置加载流程

1. **加载默认配置** - 从 `defaultConfig()` 函数
2. **加载配置文件** - 查找并解析 YAML 文件
3. **加载环境变量** - 解析 `APP_` 前缀的环境变量
4. **应用命令行参数** - 覆盖特定配置项

## 使用配置

在代码中访问配置：

```go
// 通过 Container 获取配置
config := container.Config

// 访问配置值
addr := config.Server.Addr
dbURL := config.Data.PgSQL.URL
jwtSecret := config.JWT.Secret
```

## 生产环境配置

### Docker 环境

使用环境变量文件 `.env`：

```bash
APP_SERVER_ENV=production
APP_SERVER_ADDR=:8080
APP_DATA_PGSQL_URL=postgresql://user:pass@postgres:5432/db
APP_DATA_REDIS_URL=redis://redis:6379/0
APP_JWT_SECRET=your-production-secret-key
```

在 `docker-compose.yml` 中引用：

```yaml
services:
  app:
    env_file:
      - .env
```

### Kubernetes 环境

使用 ConfigMap 和 Secret：

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  APP_SERVER_ADDR: ":8080"
  APP_SERVER_ENV: "production"
  APP_DATA_PGSQL_URL: "postgresql://..."
  APP_DATA_REDIS_URL: "redis://..."
---
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
stringData:
  APP_JWT_SECRET: "your-secret-key"
```

## 配置验证

启动时会自动验证关键配置：

- 数据库连接字符串有效性
- Redis 连接字符串有效性
- JWT Secret 不为空（生产环境）

## 添加新配置

1. **修改配置结构**：在 `internal/infrastructure/config/config.go` 中添加字段
2. **更新默认值**：在 `defaultConfig()` 函数中添加默认值
3. **更新示例配置**：运行 `sync-config-example` 技能更新 `configs/config.example.yaml`
4. **使用新配置**：在代码中通过 `container.Config` 访问

### 示例：添加日志配置

```go
// 1. 修改 Config 结构
type Config struct {
    Server ServerConfig `koanf:"server"`
    Data   DataConfig   `koanf:"data"`
    JWT    JWTConfig    `koanf:"jwt"`
    Log    LogConfig    `koanf:"log"`  // 新增
}

type LogConfig struct {
    Level  string `koanf:"level"`
    Format string `koanf:"format"`
}

// 2. 更新默认配置
func defaultConfig() *Config {
    return &Config{
        // ...
        Log: LogConfig{
            Level:  "info",
            Format: "json",
        },
    }
}
```

然后运行：

```bash
# 运行技能更新示例配置
claude skill sync-config-example
```

## 配置安全

### 敏感信息处理

- ❌ **不要**将敏感信息（如密码、密钥）提交到 Git
- ✅ 使用环境变量或密钥管理服务
- ✅ 使用 `.env` 文件（添加到 `.gitignore`）
- ✅ 在生产环境使用 Kubernetes Secrets 或类似服务

### 配置文件安全

```bash
# 将敏感配置添加到 .gitignore
echo "config.yaml" >> .gitignore
echo ".env" >> .gitignore
```

## 最佳实践

1. **使用环境变量**：生产环境优先使用环境变量
2. **配置文件用于开发**：开发环境可以使用配置文件
3. **提供示例配置**：维护 `configs/config.example.yaml`
4. **验证配置**：启动时验证必需的配置项
5. **文档化配置**：为每个配置项添加注释说明

## 下一步

- 了解[认证授权](/guide/authentication)
- 学习 [PostgreSQL 集成](/guide/postgresql)
- 探索 [Redis 缓存](/guide/redis)
