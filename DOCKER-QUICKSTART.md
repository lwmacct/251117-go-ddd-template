# 🐳 Docker 快速启动指南

**3 分钟部署到生产环境** 🚀

---

## ⚡ 快速开始

```bash
# 1️⃣ 复制配置模板
cp docker-compose.example.yml docker-compose.yml
cp .env.example .env

# 2️⃣ 生成安全密钥并写入 .env
cat >> .env <<EOF
JWT_SECRET=$(openssl rand -base64 32)
DEV_SECRET=$(openssl rand -hex 16)
POSTGRES_PASSWORD=$(openssl rand -base64 24)
GITHUB_USERNAME=your-github-username
PROJECT_NAME=go-ddd-template
EOF

# 3️⃣ 启动服务
docker-compose up -d

# 4️⃣ 执行数据库迁移
docker-compose exec app app migrate up

# 5️⃣ 验证部署
curl http://localhost:8080/health
```

---

## 📂 文件清单

| 文件 | 说明 | 是否提交到 Git |
|-----|------|---------------|
| `docker-compose.example.yml` | Docker Compose 配置模板 | ✅ 是（模板） |
| `docker-compose.yml` | 实际使用的配置（从模板复制） | ❌ 否（.gitignore） |
| `.env.example` | 环境变量模板 | ✅ 是（模板） |
| `.env` | 实际使用的环境变量（包含密钥） | ❌ 否（.gitignore） |
| `DEPLOYMENT.md` | 完整部署文档 | ✅ 是 |

---

## 🔑 必需的环境变量

**这 3 个变量必须设置，否则应用无法启动**：

```bash
# .env 文件
JWT_SECRET=生成的随机密钥（至少32字符）
POSTGRES_PASSWORD=数据库密码
DEV_SECRET=开发模式密钥
```

**生成方法**：

```bash
# JWT 密钥
openssl rand -base64 32

# 数据库密码
openssl rand -base64 24

# 开发密钥
openssl rand -hex 16
```

---

## 🌍 访问地址

| 服务 | URL |
|-----|-----|
| 🎨 前端应用 | http://localhost:8080/ |
| 📘 API 文档 | http://localhost:8080/swagger/index.html |
| 📚 项目文档 | http://localhost:8080/docs/ |
| 💚 健康检查 | http://localhost:8080/health |

---

## 🛠️ 常用命令

```bash
# 查看日志
docker-compose logs -f app

# 查看所有容器状态
docker-compose ps

# 停止服务
docker-compose down

# 重启应用
docker-compose restart app

# 进入容器
docker-compose exec app sh

# 数据库迁移
docker-compose exec app app migrate up
docker-compose exec app app migrate status

# 备份数据库
docker-compose exec postgres pg_dump -U postgres myapp > backup.sql

# 恢复数据库
cat backup.sql | docker-compose exec -T postgres psql -U postgres myapp
```

---

## 🔍 故障排查

### 问题：应用无法启动

```bash
# 查看日志
docker-compose logs app

# 常见原因：缺少环境变量
# 解决：检查 .env 文件是否包含 JWT_SECRET、POSTGRES_PASSWORD、DEV_SECRET
```

### 问题：数据库连接失败

```bash
# 检查 PostgreSQL 状态
docker-compose exec postgres pg_isready -U postgres

# 查看数据库日志
docker-compose logs postgres
```

### 问题：静态文件 404

```bash
# 检查文件是否存在
docker-compose exec app ls -la /apps/data/public/web/
docker-compose exec app ls -la /apps/data/public/docs/

# 检查环境变量
docker-compose exec app env | grep APP_SERVER
```

---

## 📊 配置优先级

应用使用三层配置合并：

```
环境变量 (APP_*)          ← 最高优先级 ⭐
    ↑
配置文件 (config.yaml)    ← 容器内置，包含基础配置
    ↑
代码默认值 (config.go)    ← 最低优先级
```

**容器内置配置**（`/apps/data/configs/config.yaml`）：

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

**⚠️ 注意**：容器配置**不包含敏感信息**，这些必须通过环境变量传递！

---

## 🔐 安全检查清单

部署前确认：

- [ ] ✅ `JWT_SECRET` 已设置为至少 32 字符的强随机字符串
- [ ] ✅ `POSTGRES_PASSWORD` 已修改为强密码
- [ ] ✅ `DEV_SECRET` 已更改
- [ ] ✅ `.env` 文件不在版本控制中
- [ ] ✅ 数据库和 Redis 端口未暴露到公网
- [ ] ✅ 生产环境启用 HTTPS（通过反向代理）

---

## 📚 详细文档

- 完整部署指南：[DEPLOYMENT.md](./DEPLOYMENT.md)
- 开发文档：[README.md](./README.md)
- 架构文档：[docs/architecture/ddd-cqrs.md](./docs/architecture/ddd-cqrs.md)

---

**需要帮助？** 查看 [DEPLOYMENT.md](./DEPLOYMENT.md) 获取详细说明。
