# 🚀 Docker 部署指南

本文档说明如何使用 Docker 部署应用到生产环境。

## 📋 目录

- [快速开始](#快速开始)
- [配置说明](#配置说明)
- [环境变量](#环境变量)
- [数据库迁移](#数据库迁移)
- [监控和日志](#监控和日志)
- [故障排查](#故障排查)

---

## ⚡ 快速开始

### 1️⃣ 准备配置文件

```bash
# 复制配置文件模板
cp docker-compose.example.yml docker-compose.yml
cp .env.example .env
```

### 2️⃣ 修改敏感配置

编辑 `.env` 文件，**必须修改**以下配置：

```bash
# 生成强随机密钥
JWT_SECRET=$(openssl rand -base64 32)
DEV_SECRET=$(openssl rand -hex 16)
POSTGRES_PASSWORD=$(openssl rand -base64 24)

# 写入 .env 文件
cat > .env <<EOF
# GitHub 配置
GITHUB_USERNAME=your-github-username
PROJECT_NAME=go-ddd-template
VERSION=latest

# 应用端口
APP_PORT=8080

# 数据库配置
POSTGRES_USER=postgres
POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
POSTGRES_DB=myapp

# JWT 配置
JWT_SECRET=${JWT_SECRET}

# 认证配置
DEV_SECRET=${DEV_SECRET}
EOF
```

### 3️⃣ 启动服务

```bash
# 拉取最新镜像
docker-compose pull

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f app
```

### 4️⃣ 执行数据库迁移

```bash
# 等待数据库就绪
docker-compose exec postgres pg_isready -U postgres

# 执行迁移（自动创建表和初始数据）
docker-compose exec app app migrate up

# 验证迁移状态
docker-compose exec app app migrate status
```

### 5️⃣ 验证部署

```bash
# 健康检查
curl http://localhost:8080/health

# API 文档
open http://localhost:8080/swagger/index.html

# 前端应用
open http://localhost:8080/

# 项目文档
open http://localhost:8080/docs/
```

---

## 🔧 配置说明

### 配置优先级

应用使用 **三层配置合并**，优先级从低到高：

```
1. 代码默认值（config.go 中的 defaultConfig）
   ↓
2. 配置文件（容器内 /apps/data/configs/config.yaml）
   ↓
3. 环境变量（APP_* 前缀）← 最高优先级
```

### 容器内配置文件

Docker 镜像内置了基础配置文件 `/apps/data/configs/config.yaml`：

```yaml
server:
  addr: "0.0.0.0:8080"
  env: "production"
  static_dir: "/apps/data/public/web"
  docs_dir: "/apps/data/public/docs"

data:
  redis_key_prefix: "myapp:"
  auto_migrate: false

jwt:
  access_token_expiry: "15m"
  refresh_token_expiry: "168h"

auth:
  twofa_issuer: "Go-DDD-Template"
  captcha_required: true
```

**⚠️ 注意**：配置文件中**不包含任何敏感信息**（数据库密码、JWT 密钥等），这些必须通过环境变量传递。

---

## 🌍 环境变量

### 必需配置

以下环境变量**必须设置**，否则应用无法启动：

| 环境变量 | 说明 | 生成方法 | 示例 |
|---------|------|---------|------|
| `JWT_SECRET` | JWT 签名密钥（至少 32 字符） | `openssl rand -base64 32` | `R3Jh...` |
| `POSTGRES_PASSWORD` | PostgreSQL 数据库密码 | `openssl rand -base64 24` | `xK9m...` |
| `DEV_SECRET` | 开发模式密钥（验证码开发模式） | `openssl rand -hex 16` | `a3f2...` |

**完整的数据库连接示例**：

```bash
# 自动构建的完整连接串
APP_DATA_PGSQL_URL="postgresql://postgres:${POSTGRES_PASSWORD}@postgres:5432/myapp?sslmode=disable"
```

### 可选配置

以下环境变量**有默认值**，通常不需要修改：

<details>
<summary>点击展开完整列表</summary>

#### 服务器配置

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `APP_SERVER_ADDR` | 监听地址 | `0.0.0.0:8080` |
| `APP_SERVER_ENV` | 运行环境 | `production` |
| `APP_SERVER_STATIC_DIR` | 前端静态资源路径 | `/apps/data/public/web` |
| `APP_SERVER_DOCS_DIR` | 文档静态资源路径 | `/apps/data/public/docs` |

#### 数据源配置

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `APP_DATA_PGSQL_URL` | PostgreSQL 连接 URL | *必需* |
| `APP_DATA_REDIS_URL` | Redis 连接 URL | `redis://redis:6379/0` |
| `APP_DATA_REDIS_KEY_PREFIX` | Redis 键前缀 | `myapp:` |
| `APP_DATA_AUTO_MIGRATE` | 自动迁移（⚠️ 生产环境不推荐） | `false` |

#### JWT 配置

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `APP_JWT_SECRET` | JWT 签名密钥 | *必需* |
| `APP_JWT_ACCESS_TOKEN_EXPIRY` | Access Token 有效期 | `15m` |
| `APP_JWT_REFRESH_TOKEN_EXPIRY` | Refresh Token 有效期 | `168h`（7 天） |

#### 认证配置

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `APP_AUTH_DEV_SECRET` | 开发模式密钥 | *必需* |
| `APP_AUTH_TWOFA_ISSUER` | 2FA TOTP 发行者名称 | `Go-DDD-Template` |
| `APP_AUTH_CAPTCHA_REQUIRED` | 是否需要验证码 | `true` |

</details>

### 环境变量命名规则

环境变量使用 `APP_` 前缀，下划线分隔，自动映射到配置结构：

```
APP_SERVER_ADDR       → server.addr
APP_DATA_PGSQL_URL    → data.pgsql_url
APP_JWT_SECRET        → jwt.secret
```

---

## 🗄️ 数据库迁移

### 自动迁移（仅开发环境）

设置环境变量启用自动迁移：

```yaml
# docker-compose.yml
environment:
  APP_DATA_AUTO_MIGRATE: "true"
```

⚠️ **不推荐在生产环境使用**，可能导致数据丢失或服务中断。

### 手动迁移（推荐）

生产环境应使用迁移命令手动控制：

```bash
# 查看迁移状态
docker-compose exec app app migrate status

# 执行迁移（升级到最新版本）
docker-compose exec app app migrate up

# 回滚一个版本
docker-compose exec app app migrate down

# 重置数据库（⚠️ 危险操作）
docker-compose exec app app migrate reset
```

### 备份数据库

```bash
# 备份数据库
docker-compose exec postgres pg_dump -U postgres myapp > backup.sql

# 恢复数据库
cat backup.sql | docker-compose exec -T postgres psql -U postgres myapp
```

---

## 📊 监控和日志

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看应用日志
docker-compose logs -f app

# 查看数据库日志
docker-compose logs -f postgres

# 查看最近 100 行日志
docker-compose logs --tail=100 app
```

### 健康检查

```bash
# HTTP 健康检查
curl http://localhost:8080/health

# 查看健康检查状态
docker-compose ps
```

**健康检查响应示例**：

```json
{
  "status": "ok",
  "timestamp": "2025-01-15T10:30:00Z",
  "checks": {
    "database": "ok",
    "redis": "ok"
  }
}
```

### 服务状态

```bash
# 查看容器状态
docker-compose ps

# 查看资源使用情况
docker stats
```

---

## 🔍 故障排查

### 问题：应用无法启动

**检查日志**：

```bash
docker-compose logs app
```

**常见原因**：

1. **缺少必需的环境变量**：
   ```
   Error: JWT_SECRET is required
   ```
   **解决**：检查 `.env` 文件，确保设置了 `JWT_SECRET`、`POSTGRES_PASSWORD`、`DEV_SECRET`

2. **数据库连接失败**：
   ```
   Error: failed to connect to database
   ```
   **解决**：
   ```bash
   # 检查 PostgreSQL 是否就绪
   docker-compose exec postgres pg_isready -U postgres

   # 查看数据库日志
   docker-compose logs postgres
   ```

### 问题：静态文件 404

**检查文件是否存在**：

```bash
# 检查前端文件
docker-compose exec app ls -la /apps/data/public/web/

# 检查文档文件
docker-compose exec app ls -la /apps/data/public/docs/
```

**检查配置**：

```bash
# 验证环境变量
docker-compose exec app env | grep APP_SERVER

# 查看配置文件
docker-compose exec app cat /apps/data/configs/config.yaml
```

### 问题：权限被拒绝

**检查用户权限**：

```bash
# 查看文件所有者
docker-compose exec app ls -la /apps/data/

# 如果需要，修改所有权（容器重启后可能重置）
docker-compose exec app chown -R nobody:nobody /apps/data/
```

### 问题：Redis 连接失败

**检查 Redis 状态**：

```bash
# 测试 Redis 连接
docker-compose exec redis redis-cli ping

# 如果启用了密码认证
docker-compose exec redis redis-cli -a "${REDIS_PASSWORD}" ping

# 查看 Redis 日志
docker-compose logs redis
```

### 问题：环境变量未生效

**验证环境变量**：

```bash
# 列出所有 APP_* 环境变量
docker-compose exec app env | grep APP_

# 检查 .env 文件
cat .env
```

**配置优先级提醒**：环境变量 > 配置文件 > 默认值

---

## 🔐 生产环境安全检查清单

部署前请确认：

- [ ] ✅ `JWT_SECRET` 已设置为至少 32 字符的强随机字符串
- [ ] ✅ `POSTGRES_PASSWORD` 已修改为强密码
- [ ] ✅ `DEV_SECRET` 已更改（或确认不会在生产环境使用开发模式）
- [ ] ✅ `.env` 文件不在版本控制中（已在 `.gitignore`）
- [ ] ✅ PostgreSQL 端口未暴露到公网（仅容器内部访问）
- [ ] ✅ Redis 端口未暴露到公网（仅容器内部访问）
- [ ] ✅ 启用 HTTPS（在反向代理层配置，如 Nginx/Traefik）
- [ ] ✅ 配置防火墙规则（仅允许必要的端口）
- [ ] ✅ 设置日志轮转（避免磁盘空间耗尽）
- [ ] ✅ 配置备份策略（数据库定期备份）

---

## 🌐 访问地址

服务启动后可访问：

| 服务 | 地址 | 说明 |
|-----|------|------|
| 前端应用 | http://localhost:8080/ | Vue 3 SPA 应用 |
| API 文档 | http://localhost:8080/swagger/index.html | Swagger UI |
| 项目文档 | http://localhost:8080/docs/ | VitePress 文档 |
| 健康检查 | http://localhost:8080/health | 系统健康状态 |

---

## 📚 相关文档

- [开发指南](./README.md)
- [API 文档](./internal/adapters/http/docs/)
- [架构文档](./docs/architecture/ddd-cqrs.md)
- [配置说明](./internal/infrastructure/config/config.go)

---

## 🆘 获取帮助

如有问题，请：

1. 查看 [故障排查](#故障排查) 章节
2. 查看应用日志：`docker-compose logs -f app`
3. 提交 Issue 到项目仓库
