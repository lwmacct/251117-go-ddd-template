# 配置管理

Go DDD Template 使用 Koanf 作为配置管理库，支持多层配置优先级和多种配置源。

<!--TOC-->

## Table of Contents

- [配置优先级](#配置优先级) `:29+9`
- [配置文件](#配置文件) `:38+28`
  - [配置示例](#配置示例) `:40+26`
- [环境变量](#环境变量) `:66+22`
- [配置结构](#配置结构) `:88+35`
  - [Server 配置](#server-配置) `:90+9`
  - [Data 配置](#data-配置) `:99+9`
  - [JWT 配置](#jwt-配置) `:108+8`
  - [Log 配置](#log-配置) `:116+7`
- [生产环境配置](#生产环境配置) `:123+35`
  - [使用 Docker](#使用-docker) `:125+9`
  - [使用 Kubernetes](#使用-kubernetes) `:134+14`
  - [使用 Systemd](#使用-systemd) `:148+10`
- [配置验证](#配置验证) `:158+8`
- [动态配置](#动态配置) `:166+8`
- [安全建议](#安全建议) `:174+17`
- [配置最佳实践](#配置最佳实践) `:191+19`

<!--TOC-->

## 配置优先级

配置按以下优先级加载（高到低）：

1. **命令行参数** - 最高优先级
2. **环境变量** - `APP_` 前缀
3. **配置文件** - `config.yaml`
4. **默认值** - 代码中定义

## 配置文件

### 配置示例

```yaml
# configs/config.yaml
server:
  addr: "0.0.0.0:8080"
  env: "development"
  static_dir: "web/dist"
  docs_dir: "docs/.vitepress/dist"

data:
  pgsql_url: "postgresql://postgres@localhost:5432/myapp?sslmode=disable"
  redis_url: "redis://localhost:6379/0"
  redis_key_prefix: "myapp:"
  auto_migrate: true

jwt:
  secret: "your-secret-key"
  access_token_expiry: "15m"
  refresh_token_expiry: "168h"

log:
  level: "info"
  format: "text" # text or json
```

## 环境变量

所有配置项都可通过环境变量覆盖，使用 `APP_` 前缀和下划线分隔：

```bash
# 服务器配置
export APP_SERVER_ADDR="0.0.0.0:9000"
export APP_SERVER_ENV="production"

# 数据库配置
export APP_DATA_PGSQL_URL="postgresql://user:pass@host:5432/db"
export APP_DATA_REDIS_URL="redis://:password@host:6379/0"

# JWT 配置
export APP_JWT_SECRET="production-secret"
export APP_JWT_ACCESS_TOKEN_EXPIRY="30m"

# 日志配置
export APP_LOG_LEVEL="debug"
export APP_LOG_FORMAT="json"
```

## 配置结构

### Server 配置

| 字段       | 类型   | 默认值                 | 说明         |
| ---------- | ------ | ---------------------- | ------------ |
| addr       | string | ":8080"                | 服务监听地址 |
| env        | string | "development"          | 运行环境     |
| static_dir | string | "web/dist"             | 静态文件目录 |
| docs_dir   | string | "docs/.vitepress/dist" | 文档目录     |

### Data 配置

| 字段             | 类型   | 默认值 | 说明                  |
| ---------------- | ------ | ------ | --------------------- |
| pgsql_url        | string | -      | PostgreSQL 连接字符串 |
| redis_url        | string | -      | Redis 连接字符串      |
| redis_key_prefix | string | "app:" | Redis 键前缀          |
| auto_migrate     | bool   | false  | 自动执行数据库迁移    |

### JWT 配置

| 字段                 | 类型   | 默认值 | 说明           |
| -------------------- | ------ | ------ | -------------- |
| secret               | string | -      | JWT 签名密钥   |
| access_token_expiry  | string | "15m"  | 访问令牌有效期 |
| refresh_token_expiry | string | "168h" | 刷新令牌有效期 |

### Log 配置

| 字段   | 类型   | 默认值 | 说明                 |
| ------ | ------ | ------ | -------------------- |
| level  | string | "info" | 日志级别             |
| format | string | "text" | 日志格式 (text/json) |

## 生产环境配置

### 使用 Docker

```dockerfile
# Dockerfile
ENV APP_SERVER_ENV=production
ENV APP_LOG_FORMAT=json
ENV APP_DATA_AUTO_MIGRATE=false
```

### 使用 Kubernetes

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  APP_SERVER_ENV: "production"
  APP_LOG_LEVEL: "info"
  APP_LOG_FORMAT: "json"
```

### 使用 Systemd

```ini
# /etc/systemd/system/app.service
[Service]
Environment="APP_SERVER_ENV=production"
Environment="APP_LOG_FORMAT=json"
EnvironmentFile=/etc/app/env
```

## 配置验证

应用启动时会自动验证必需的配置项：

- JWT secret 不能为空
- 数据库连接必须有效
- Redis 连接必须有效

## 动态配置

某些配置支持运行时重载：

- 日志级别（通过 API 调整）
- 缓存策略
- 限流参数

## 安全建议

1. **生产环境密钥管理**
   - 使用环境变量或密钥管理系统
   - 不要将密钥提交到代码仓库
   - 定期轮换密钥

2. **数据库连接**
   - 使用 SSL/TLS 连接
   - 限制数据库访问 IP
   - 使用只读用户进行查询

3. **Redis 安全**
   - 设置访问密码
   - 启用 TLS 连接
   - 限制命令集

## 配置最佳实践

1. **环境分离**

   ```
   configs/
   ├── config.dev.yaml      # 开发环境
   ├── config.test.yaml     # 测试环境
   └── config.prod.yaml     # 生产环境
   ```

2. **敏感信息管理**
   - 使用 `.env` 文件（本地开发）
   - 使用 Vault/AWS Secrets Manager（生产环境）

3. **配置文档化**
   - 提供完整的配置示例
   - 注释说明每个配置项
   - 记录默认值和约束
