# 快速开始

本指南将帮助你快速搭建和运行 Go DDD Template 项目。

## 环境要求

- Go 1.25.4 或更高版本
- Docker 和 Docker Compose（用于运行 PostgreSQL 和 Redis）
- Task（可选，用于任务自动化）

## 安装

### 1. 克隆项目

```bash
git clone https://github.com/lwmacct/251117-go-ddd-template.git
cd 251117-go-ddd-template
```

### 2. 启动依赖服务

使用 Docker Compose 启动 PostgreSQL 和 Redis：

```bash
docker-compose up -d
```

这将启动：

- PostgreSQL（端口 5432）
- Redis（端口 6379）

### 3. 配置应用

创建配置文件 `config.yaml`（可选）：

```yaml
server:
  addr: ":8080"
  env: "development"

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

或使用环境变量：

```bash
export APP_SERVER_ADDR=":8080"
export APP_DATA_PGSQL_URL="postgresql://postgres:postgres@localhost:5432/app?sslmode=disable"
export APP_DATA_REDIS_URL="redis://localhost:6379/0"
export APP_JWT_SECRET="your-secret-key"
```

### 4. 运行应用

使用 Task（推荐）：

```bash
# 安装 Task（如果还没安装）
go install github.com/go-task/task/v3/cmd/task@latest

# 构建并运行
task go:run -- api
```

或直接构建：

```bash
# 构建
task go:build

# 运行
.local/bin/251117-go-ddd-template api
```

或使用开发热重载：

```bash
# 安装 air
go install github.com/air-verse/air@latest

# 使用 air 运行（支持热重载）
air
```

## 验证安装

### 健康检查

```bash
curl http://localhost:8080/health
```

应返回：

```json
{
  "status": "ok",
  "database": "connected",
  "redis": "connected"
}
```

### 注册用户

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'
```

### 登录获取 Token

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "testuser",
    "password": "password123"
  }'
```

### 访问受保护的端点

```bash
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 下一步

- 了解[项目架构](/guide/architecture)
- 查看[配置系统](/guide/configuration)
- 学习[认证授权](/guide/authentication)
- 探索 [API 文档](/api/)

## 故障排查

### 数据库连接失败

确保 PostgreSQL 正在运行：

```bash
docker-compose ps
```

检查数据库连接字符串是否正确。

### Redis 连接失败

确保 Redis 正在运行：

```bash
docker-compose ps
```

检查 Redis 连接字符串是否正确。

### 端口被占用

如果 8080 端口被占用，可以通过环境变量或配置文件修改：

```bash
APP_SERVER_ADDR=":9000" task go:run -- api
```
