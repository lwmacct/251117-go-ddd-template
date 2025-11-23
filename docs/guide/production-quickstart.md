# 生产环境快速部署

本指南将帮助你在 **5 分钟内** 使用 Docker Compose 部署 Go DDD Template 应用到生产环境。

> 💡 **适用场景**: 单机或小型生产环境、快速测试部署
> 🎯 **目标用户**: 运维人员、系统管理员

## 前置要求

- Docker 20.10+ 和 Docker Compose 2.0+
- 至少 2GB 可用内存
- 至少 10GB 可用磁盘空间

检查环境：

```bash
docker --version
docker-compose --version
```

## 快速开始

### 1. 获取项目文件

```bash
# 克隆项目
git clone https://github.com/your-org/go-ddd-template.git
cd go-ddd-template

# 或从发布包中解压
# tar -xzf go-ddd-template-v1.0.0.tar.gz
# cd go-ddd-template
```

### 2. 配置环境变量

复制环境变量模板：

```bash
cp .env.example .env
```

**必须修改**以下配置项（打开 `.env` 文件编辑）：

```bash
# ⚠️ 数据库密码 - 必须修改
POSTGRES_PASSWORD=your-strong-password-here

# ⚠️ JWT 密钥 - 必须修改（至少 32 字符）
JWT_SECRET=your-very-secret-jwt-key-change-in-production

# ⚠️ 开发密钥 - 必须修改
DEV_SECRET=your-dev-secret-change-me
```

**生成强密码**：

```bash
# 生成 PostgreSQL 密码
openssl rand -base64 24

# 生成 JWT 密钥
openssl rand -base64 32

# 生成开发密钥
openssl rand -hex 16
```

将生成的值替换到 `.env` 文件中对应的位置。

### 3. 启动服务

```bash
# 启动所有服务（后台运行）
docker-compose up -d

# 查看服务状态
docker-compose ps
```

预期输出：

```
NAME                COMMAND                  SERVICE    STATUS
go-ddd-app          "/app/app api"           app        Up (healthy)
go-ddd-postgres     "docker-entrypoint..."   postgres   Up (healthy)
go-ddd-redis        "docker-entrypoint..."   redis      Up (healthy)
```

### 4. 验证部署

```bash
# 检查健康状态
curl http://localhost:8080/health

# 预期响应: {"status":"ok"}
```

访问应用：

- **API 服务**: `http://localhost:8080/api/`
- **API 文档**: `http://localhost:8080/api/swagger/index.html`
- **前端应用**: `http://localhost:8080/`

🎉 **部署完成！** 应用现在已经在生产环境运行。

## 环境变量配置详解

### 必需配置（生产环境必须修改）

| 变量                 | 说明                     | 生成方法                     | 示例                                      |
| -------------------- | ------------------------ | ---------------------------- | ----------------------------------------- |
| `POSTGRES_PASSWORD`  | PostgreSQL 数据库密码    | `openssl rand -base64 24`    | `abc123XYZ...` (24 字符以上)              |
| `JWT_SECRET`         | JWT 签名密钥             | `openssl rand -base64 32`    | `your-secret...` (32 字符以上)            |
| `DEV_SECRET`         | 开发模式密钥             | `openssl rand -hex 16`       | `a1b2c3d4...` (16 字节 hex)               |

### 可选配置

| 变量                 | 默认值              | 说明                                 |
| -------------------- | ------------------- | ------------------------------------ |
| `APP_PORT`           | `8080`              | 应用端口（宿主机映射端口）           |
| `REDIS_PASSWORD`     | (空)                | Redis 密码，默认不使用密码           |
| `TWOFA_ISSUER`       | `Go-DDD-Template`   | 2FA 验证器中显示的应用名称           |
| `CAPTCHA_REQUIRED`   | `true`              | 是否启用验证码                       |

### 高级配置（通常不需要修改）

如需覆盖默认配置，在 `.env` 文件中取消注释并修改：

```bash
# 服务器配置
# APP_SERVER_ADDR=0.0.0.0:8080
# APP_SERVER_ENV=production

# 数据库配置
# APP_DATA_AUTO_MIGRATE=false  # 生产环境建议关闭自动迁移

# JWT 配置
# APP_JWT_ACCESS_TOKEN_EXPIRY=15m    # 访问令牌有效期
# APP_JWT_REFRESH_TOKEN_EXPIRY=168h  # 刷新令牌有效期（7天）
```

完整的环境变量列表请参考 [配置系统文档](/guide/configuration)。

## 数据库初始化

首次部署时需要初始化数据库：

```bash
# 进入应用容器
docker exec -it go-ddd-app sh

# 运行数据库迁移
./app migrate up

# 创建初始管理员用户（可选）
./app seed

# 退出容器
exit
```

## 服务管理

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看应用日志
docker-compose logs -f app

# 查看数据库日志
docker-compose logs -f postgres
```

### 重启服务

```bash
# 重启应用
docker-compose restart app

# 重启所有服务
docker-compose restart
```

### 停止服务

```bash
# 停止所有服务
docker-compose stop

# 停止并删除容器（保留数据）
docker-compose down

# 停止并删除容器和数据卷（⚠️ 会删除数据库数据）
docker-compose down -v
```

### 更新应用

```bash
# 拉取最新代码
git pull

# 重新构建并启动
docker-compose up -d --build

# 运行数据库迁移（如果有schema变更）
docker exec -it go-ddd-app ./app migrate up
```

## 数据备份

### PostgreSQL 备份

```bash
# 手动备份数据库
docker exec go-ddd-postgres pg_dump -U postgres app > backup_$(date +%Y%m%d_%H%M%S).sql

# 恢复数据库
docker exec -i go-ddd-postgres psql -U postgres app < backup_20250123_120000.sql
```

### 自动备份（使用 cron）

创建备份脚本 `/usr/local/bin/backup-db.sh`：

```bash
#!/bin/bash
BACKUP_DIR="/backup/postgres"
mkdir -p $BACKUP_DIR
docker exec go-ddd-postgres pg_dump -U postgres app > \
  $BACKUP_DIR/backup_$(date +%Y%m%d_%H%M%S).sql

# 保留最近 7 天的备份
find $BACKUP_DIR -name "backup_*.sql" -mtime +7 -delete
```

添加到 crontab（每天凌晨 2 点备份）：

```bash
0 2 * * * /usr/local/bin/backup-db.sh
```

## 监控和维护

### 健康检查

```bash
# 检查应用健康状态
curl http://localhost:8080/health

# 检查数据库连接
docker exec go-ddd-postgres psql -U postgres -c "SELECT 1"

# 检查 Redis 连接
docker exec go-ddd-redis redis-cli ping
```

### 资源监控

```bash
# 查看容器资源使用情况
docker stats

# 查看磁盘使用情况
docker system df
```

### 清理未使用资源

```bash
# 清理未使用的镜像和容器
docker system prune -a

# 清理未使用的数据卷
docker volume prune
```

## 故障排查

### 问题 1: 容器启动失败

**症状**: `docker-compose ps` 显示容器状态为 `Exited`

**排查步骤**:

```bash
# 查看容器日志
docker-compose logs app

# 常见原因：
# 1. 环境变量未设置 - 检查 .env 文件
# 2. 端口被占用 - 修改 APP_PORT
# 3. 内存不足 - 增加系统内存或调整 Docker 限制
```

**解决方案**:

```bash
# 检查环境变量
cat .env

# 检查端口占用
sudo lsof -i :8080

# 修改端口后重新启动
docker-compose up -d
```

### 问题 2: 数据库连接失败

**症状**: 应用日志显示 `connection refused` 或 `database not found`

**排查步骤**:

```bash
# 检查数据库状态
docker-compose ps postgres

# 检查数据库日志
docker-compose logs postgres

# 测试数据库连接
docker exec go-ddd-postgres psql -U postgres -c "SELECT version()"
```

**解决方案**:

```bash
# 等待数据库完全启动（健康检查通过）
docker-compose up -d
sleep 10

# 检查环境变量中的数据库连接字符串
grep POSTGRES .env
```

### 问题 3: 502 Bad Gateway

**症状**: 访问应用返回 502 错误

**排查步骤**:

```bash
# 检查应用是否正常运行
docker-compose ps app

# 检查应用日志
docker-compose logs app

# 检查健康状态
curl http://localhost:8080/health
```

**解决方案**:

```bash
# 重启应用
docker-compose restart app

# 如果问题持续，查看详细日志
docker-compose logs -f app
```

### 问题 4: 磁盘空间不足

**症状**: 容器无法启动，日志显示 `no space left on device`

**排查步骤**:

```bash
# 检查磁盘使用情况
df -h
docker system df
```

**解决方案**:

```bash
# 清理 Docker 资源
docker system prune -a --volumes

# 清理旧的数据库备份
find /backup -name "*.sql" -mtime +30 -delete
```

### 问题 5: JWT 认证失败

**症状**: 登录后立即被退出，或 API 返回 `401 Unauthorized`

**排查步骤**:

```bash
# 检查 JWT_SECRET 是否设置
grep JWT_SECRET .env

# 检查应用日志中的认证错误
docker-compose logs app | grep -i "jwt\|auth"
```

**解决方案**:

```bash
# 确保 JWT_SECRET 已设置且长度足够
# 编辑 .env 文件，设置强密钥
JWT_SECRET=$(openssl rand -base64 32)

# 重启应用使配置生效
docker-compose restart app
```

## 安全加固

### 1. 使用 HTTPS

建议使用 Nginx 或 Traefik 作为反向代理，配置 SSL 证书：

- 使用 Let's Encrypt 获取免费 SSL 证书
- 配置 HTTPS 重定向
- 启用 HTTP/2

详见 [应用部署指南 - Nginx 配置](/guide/application-deployment#nginx-反向代理配置)

### 2. 限制外部访问

```yaml
# docker-compose.yml 中移除端口映射
services:
  postgres:
    # 注释掉或删除这一行
    # ports:
    #   - "5432:5432"
```

### 3. 使用 Docker Secrets（Docker Swarm）

对于更高安全要求的场景，使用 Docker Secrets 管理敏感信息，而不是环境变量文件。

### 4. 定期更新

```bash
# 定期拉取最新镜像
docker-compose pull

# 重新构建应用
docker-compose up -d --build
```

## 性能优化

### 1. 启用 Redis 缓存

Redis 已包含在 docker-compose.yml 中，应用会自动使用。

### 2. 配置数据库连接池

在 `.env` 文件中添加：

```bash
# 数据库连接池配置（通过应用配置）
APP_DATA_PGSQL_URL=postgresql://postgres:password@postgres:5432/app?sslmode=disable&pool_max_conns=25
```

### 3. 增加容器资源限制

编辑 `docker-compose.yml`，为应用容器增加资源：

```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
```

## 下一步

- 📊 [配置监控和日志收集](#) (Prometheus + Grafana)
- 🔐 [配置 HTTPS 和域名](#) (Nginx + Let's Encrypt)
- ☸️ [迁移到 Kubernetes](/guide/application-deployment#kubernetes-部署) (大规模部署)
- 📖 [了解配置系统](/guide/configuration) (环境变量详解)
- 🔧 [CLI 命令参考](/guide/cli-commands) (数据库迁移、用户管理)

## 获取帮助

- 📚 [完整部署文档](/guide/application-deployment)
- 🐛 [问题反馈](https://github.com/your-org/go-ddd-template/issues)
- 💬 [社区讨论](https://github.com/your-org/go-ddd-template/discussions)
