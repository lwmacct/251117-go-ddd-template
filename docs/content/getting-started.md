# 快速入门

Go DDD Template 是一个基于领域驱动设计（DDD）和 CQRS 模式的企业级应用模板。

<!--TOC-->

## Table of Contents

- [核心价值](#核心价值) `:25+7`
- [环境要求](#环境要求) `:32+6`
- [安装步骤](#安装步骤) `:38+21`
- [验证安装](#验证安装) `:59+14`
- [项目结构](#项目结构) `:73+16`
- [配置管理](#配置管理) `:89+18`
  - [核心配置项](#核心配置项) `:98+9`
- [CLI 命令](#cli-命令) `:107+18`
- [开发工具配置](#开发工具配置) `:125+8`
- [生产部署](#生产部署) `:133+19`
  - [Docker](#docker) `:135+10`
  - [安全建议](#安全建议) `:145+7`
- [下一步](#下一步) `:152+7`

<!--TOC-->

## 核心价值

- **架构清晰**: DDD 四层架构 + CQRS 模式
- **开发高效**: Task 任务自动化、热重载
- **质量保证**: 完整的测试策略和代码规范
- **生产就绪**: Docker 容器化、健康检查

## 环境要求

- Go 1.25.4+
- Docker 和 Docker Compose
- Task（可选，用于任务自动化）

## 安装步骤

```bash
# 1. 克隆项目
git clone https://github.com/lwmacct/251117-go-ddd-template.git
cd 251117-go-ddd-template

# 2. 启动依赖服务
docker-compose up -d

# 3. 数据库迁移
go run main.go migrate up

# 4. 填充种子数据（可选）
go run main.go seed

# 5. 运行应用
task go:run -- api
# 或使用热重载: air
```

## 验证安装

```bash
# 健康检查
curl http://localhost:8080/health

# 登录获取 Token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login": "admin", "password": "password123"}'
```

预置账号: `admin / password123`

## 项目结构

```
251117-go-ddd-template/
├── internal/               # 核心业务代码
│   ├── adapters/          # 适配器层
│   ├── application/       # 应用层
│   ├── domain/            # 领域层
│   ├── infrastructure/    # 基础设施层
│   └── bootstrap/         # 依赖注入
├── src/                   # 前端源代码（Vue 3）
├── docs/                  # 项目文档（VitePress）
├── configs/               # 配置文件
└── main.go                # 应用入口
```

## 配置管理

配置按以下优先级加载（高到低）：

1. 命令行参数
2. 环境变量（`APP_` 前缀）
3. 配置文件（`config.yaml`）
4. 默认值

### 核心配置项

| 配置           | 环境变量             | 说明            |
| -------------- | -------------------- | --------------- |
| server.addr    | `APP_SERVER_ADDR`    | 服务监听地址    |
| data.pgsql_url | `APP_DATA_PGSQL_URL` | PostgreSQL 连接 |
| data.redis_url | `APP_DATA_REDIS_URL` | Redis 连接      |
| jwt.secret     | `APP_JWT_SECRET`     | JWT 签名密钥    |

## CLI 命令

```bash
./251117-go-ddd-template api           # 启动 HTTP 服务
./251117-go-ddd-template migrate up    # 执行迁移
./251117-go-ddd-template migrate down  # 回滚迁移
./251117-go-ddd-template seed          # 填充种子数据
./251117-go-ddd-template worker        # 启动后台任务
```

使用 Task:

```bash
task go:run -- api
task go:run -- migrate up
task go:build
```

## 开发工具配置

| 工具       | 配置文件                  | 说明       |
| ---------- | ------------------------- | ---------- |
| Air        | `.air.toml`               | 热重载     |
| Task       | `Taskfile.yml`            | 任务自动化 |
| Pre-commit | `.pre-commit-config.yaml` | 提交检查   |

## 生产部署

### Docker

```bash
docker build -t go-ddd-template .
docker run -e APP_SERVER_ENV=production \
  -e APP_DATA_PGSQL_URL=... \
  -e APP_JWT_SECRET=... \
  go-ddd-template
```

### 安全建议

- 使用环境变量或密钥管理系统存储敏感信息
- 数据库使用 SSL/TLS 连接
- JWT secret 至少 32 字节随机字符串
- 定期轮换密钥

## 下一步

- [架构文档](./architecture/ddd-cqrs.md) - DDD + CQRS 详解
- [开发指南](./development.md) - 环境设置和测试
- [前端架构](./architecture/frontend.md) - Vue 3 前端
- [身份认证](./architecture/identity.md) - 认证和权限
- Swagger UI (`/swagger/index.html`) - API 文档
