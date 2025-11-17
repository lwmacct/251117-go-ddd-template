# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此仓库中工作时提供指导。

## 项目概述

基于 Go 的 DDD（领域驱动设计）模板应用，使用 Gin 提供 HTTP 服务，Koanf 进行配置管理。项目遵循整洁架构原则，职责分离清晰。

## 常用命令

本项目使用 [Taskfile](https://taskfile.dev) 进行任务自动化。使用 `task -a` 查看所有可用任务。

### 开发

- `task go:build` - 构建应用程序二进制文件到 `.local/bin/`
- `task go:run -- api` - 构建并运行 API 服务器
- `air` - 使用热重载进行开发（使用 `.air.toml` 配置）
- `docker-compose up -d` - 启动 PostgreSQL 和 Redis 服务

### 运行应用

```bash
# 使用 task
task go:run -- api

# 构建后直接执行
.local/bin/<app-name> api

# 使用自定义地址
.local/bin/<app-name> api --addr :9000

# 使用环境变量
APP_SERVER_ADDR=:8080 .local/bin/<app-name> api
```

### Git 操作

- `task git:push` - 推送代码到远程仓库
- `task git:tag:next` - 创建下一个版本标签
- `task git:clear` - 清除提交历史（谨慎使用）

### 发布

- `task go:release` - 构建并推送所有架构的 Docker 镜像
- `task go:release:x86_64` - 专门构建 x86_64 架构

### 文档

本项目使用 VitePress 2.0 构建文档，文档源文件位于 `docs/` 目录，拥有独立的 `package.json`。

- `cd docs && npm run dev` - 启动文档开发服务器（http://localhost:5173/docs/）
- `cd docs && npm run build` - 构建文档静态文件到 `docs/.vitepress/dist/`
- `cd docs && npm run preview` - 预览构建后的文档

**重要提示**：
- docs/ 目录有自己的 package.json 和 node_modules/
- 文档依赖独立于后端和前端项目
- GitHub Actions 会自动部署文档到 GitHub Pages

## 架构

### 分层结构

```
internal/
├── commands/          # CLI 命令（入口点）
│   └── api/          # API 服务器命令
├── adapters/         # 外部接口（HTTP、gRPC 等）
│   └── http/
│       ├── handler/          # HTTP 处理器
│       │   ├── auth.go       # 认证处理器（注册、登录、刷新）
│       │   ├── health.go     # 健康检查
│       │   ├── cache.go      # 缓存操作
│       │   └── user.go       # 用户管理
│       ├── middleware/
│       │   └── jwt.go        # JWT 认证中间件
│       ├── router.go
│       └── server.go
├── bootstrap/        # 应用初始化和依赖注入容器
│   └── container.go
├── domain/           # 领域层
│   └── user/         # 用户领域
│       ├── model.go      # 用户模型和 DTO
│       └── repository.go # 用户仓储接口
├── infrastructure/   # 技术实现
│   ├── auth/         # 认证基础设施
│   │   ├── jwt.go        # JWT 管理器
│   │   └── service.go    # 认证服务
│   ├── config/       # 配置管理（Koanf）
│   ├── database/     # 数据库连接和迁移
│   │   ├── connection.go # PostgreSQL 连接管理
│   │   └── migrator.go   # 数据库迁移
│   ├── persistence/  # 仓储实现
│   │   └── user_repository.go # 用户仓储 GORM 实现
│   └── redis/        # Redis 客户端和仓储
│       ├── client.go           # Redis 连接管理
│       └── cache_repository.go # 缓存仓储接口
└── shared/          # 共享工具
    └── errors/      # 自定义错误类型
```

### 关键设计模式

1. **依赖注入容器** (`internal/bootstrap/container.go`)

   - 在一个地方初始化所有依赖
   - 当前包含：Config、DB、RedisClient、UserRepository、JWTManager、AuthService、Router
   - 添加新服务时扩展此处
   - 提供 `Close()` 方法优雅关闭资源
   - 启动时自动执行数据库迁移

2. **配置系统** (基于 Koanf，见 `internal/infrastructure/config/`)

   - 多层优先级：默认值 → 配置文件 → 环境变量 → 命令行参数
   - 环境变量格式：`APP_<SECTION>_<KEY>`（例如：`APP_SERVER_ADDR`）
   - 配置文件搜索路径：`config.yaml`、`configs/config.yaml`
   - **重要**：修改 `config.go` 中的 `defaultConfig()` 时，运行 `sync-config-example` 技能来更新 `configs/config.example.yaml`

3. **CLI 结构** (urfave/cli v3)
   - 主入口：`main.go` → `buildCommands()`
   - 每个命令位于 `internal/commands/<name>/`
   - 当前只有 `api` 命令用于 REST API 服务器

### HTTP 层

- 框架：Gin
- 路由设置：`internal/adapters/http/router.go`
- 服务器封装：`internal/adapters/http/server.go`（处理优雅关闭）
- 处理器目录：`internal/adapters/http/handler/`
  - `auth.go` - 用户认证（注册、登录、刷新、当前用户）
  - `health.go` - 健康检查（数据库和 Redis）
  - `cache.go` - 缓存操作示例
  - `user.go` - 用户管理 CRUD
- 中间件：`internal/adapters/http/middleware/`
  - `jwt.go` - JWT 认证中间件
- 静态文件服务：通过 `Server.StaticDir` 配置（默认 `web/dist`，用于 SPA）

### 认证系统

- JWT 管理器：`internal/infrastructure/auth/jwt.go`
- 认证服务：`internal/infrastructure/auth/service.go`
- 功能：
  - 用户注册（用户名/邮箱唯一性验证）
  - 用户登录（支持用户名或邮箱）
  - Token 刷新（访问令牌 15 分钟，刷新令牌 7 天）
  - JWT 中间件保护路由
  - bcrypt 密码加密
  - 用户状态检查（仅 active 用户可登录）
- 使用文档：见 `docs/authentication.md`

### PostgreSQL 集成

- 连接管理：`internal/infrastructure/database/connection.go`
- 迁移支持：`internal/infrastructure/database/migrator.go`
- 功能：
  - 连接池管理（最大 25，空闲 10，生命周期 5 分钟）
  - 自动迁移（应用启动时）
  - 健康检查和连接池统计
  - 优雅关闭
- 使用文档：见 `docs/postgresql.md`

### Redis 集成

- 客户端：`internal/infrastructure/redis/client.go`
- 仓储接口：`internal/infrastructure/redis/cache_repository.go`
- 功能：
  - 自动 JSON 序列化/反序列化
  - 支持 TTL 过期时间
  - 分布式锁（SetNX）
  - 健康检查
  - 优雅关闭
- 使用文档：见 `docs/redis.md`

### 领域层（DDD）

- 用户领域：`internal/domain/user/`
  - `model.go` - 用户模型、DTO、业务逻辑
  - `repository.go` - 仓储接口定义
- 仓储实现：`internal/infrastructure/persistence/`
  - `user_repository.go` - GORM 实现
- 特性：
  - 软删除支持
  - bcrypt 密码加密
  - 分页查询
  - 唯一性约束（username、email）

### 扩展应用

添加新功能时：

1. **新 HTTP 端点**：
   - 创建 handler 在 `internal/adapters/http/handler/`
   - 在 `internal/adapters/http/router.go` 注册路由
2. **新 CLI 命令**：在 `internal/commands/<name>/` 创建并在 `main.go` 中注册
3. **新配置项**：
   - 添加到 `internal/infrastructure/config/config.go` 的 `Config` 结构体
   - 更新 `defaultConfig()` 函数
   - 运行 `sync-config-example` 技能更新示例配置
4. **新依赖**：添加到 `bootstrap.Container` 并在 `NewContainer()` 中初始化
5. **新领域模型**：
   - 在 `internal/domain/<name>/` 创建模型和仓储接口
   - 在 `internal/infrastructure/persistence/` 实现仓储
   - 在 `bootstrap.Container` 中注入

## API 端点快速参考

### 公开端点（无需认证）

- `POST /api/auth/register` - 注册新用户
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/refresh` - 刷新访问令牌
- `GET /health` - 健康检查
- `POST/GET/DELETE /api/cache/*` - 缓存操作

### 受保护端点（需要 JWT）

- `GET /api/auth/me` - 获取当前用户信息
- `GET /api/users` - 获取用户列表（分页）
- `GET /api/users/:id` - 获取用户详情
- `PUT /api/users/:id` - 更新用户
- `DELETE /api/users/:id` - 删除用户（软删除）

详细文档：

- 认证 API：`docs/authentication.md`
- PostgreSQL 使用：`docs/postgresql.md`
- Redis 使用：`docs/redis.md`

## 配置

配置优先级（从低到高）：

1. 默认值（`internal/infrastructure/config/config.go:defaultConfig()`）
2. 配置文件（`config.yaml` 或 `configs/config.yaml`）
3. 环境变量（前缀：`APP_`）
4. 命令行参数

环境变量示例：

```bash
APP_SERVER_ADDR=:8080
APP_SERVER_ENV=production
APP_DATA_PGSQL_URL=postgresql://user:pass@host:5432/db
APP_DATA_REDIS_URL=redis://localhost:6379/0
APP_JWT_SECRET=your-secret-key-change-in-production
APP_JWT_ACCESS_TOKEN_EXPIRY=15m
APP_JWT_REFRESH_TOKEN_EXPIRY=168h
```

## 开发环境

- 支持 Dev Container（`.devcontainer/`）
- Air 热重载（`.air.toml`）
- Docker Compose（`docker-compose.yml`）：PostgreSQL + Redis
- 模块路径：`github.com/lwmacct/251117-go-ddd-template`

### 项目结构

本项目是一个 Monorepo，包含三个子项目：

1. **根目录（后端）** - Go DDD 应用
   - 依赖管理：`go.mod`、`go.sum`
   - 主入口：`main.go`
   - 核心代码：`internal/`

2. **web/** - 前端项目（Vue 3）
   - 依赖管理：`web/package.json`
   - 独立的 `node_modules/`
   - 构建产物：`web/dist/`（被后端作为静态文件服务）

3. **docs/** - VitePress 文档
   - 依赖管理：`docs/package.json`
   - 独立的 `node_modules/`
   - 构建产物：`docs/.vitepress/dist/`（部署到 GitHub Pages）

### 启动开发环境

#### 后端开发

```bash
# 1. 启动数据库和 Redis
docker-compose up -d

# 2. 运行应用（自动迁移数据库）
task go:run -- api

# 3. 健康检查
curl http://localhost:8080/health

# 4. 测试认证
# 注册用户
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'

# 登录获取 token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "testuser",
    "password": "password123"
  }'

# 使用 token 访问受保护端点
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### 文档开发

```bash
# 进入文档目录
cd docs

# 首次安装依赖
npm install

# 启动开发服务器（热重载）
npm run dev
# 访问 http://localhost:5173/docs/

# 构建文档
npm run build

# 预览构建结果
npm run preview
```

## 模块信息

- Go 版本：1.25.4
- 主要依赖：Gin、Koanf、urfave/cli/v3、GORM、go-redis/v9、golang-jwt/jwt/v5、bcrypt
- 模块路径：`github.com/lwmacct/251117-go-ddd-template`

## 已实现功能

- ✅ HTTP 服务器（Gin）+ 优雅关闭
- ✅ 配置管理（Koanf）- 多层优先级（默认值/文件/环境变量/CLI）
- ✅ PostgreSQL 集成 - GORM ORM + 自动迁移 + 连接池
- ✅ Redis 集成 - 缓存仓储 + JSON 序列化 + 分布式锁
- ✅ JWT 认证授权 - 注册/登录/刷新 + JWT 中间件 + bcrypt 加密
- ✅ 用户管理 - 完整的 CRUD API + 软删除 + 分页查询
- ✅ 健康检查 - 数据库和 Redis 状态监控
- ✅ DDD 分层架构 - Domain/Infrastructure/Adapters/Bootstrap
- ✅ 仓储模式（Repository Pattern）- 接口驱动设计
- ✅ 依赖注入容器 - 集中管理所有依赖
- ✅ Docker Compose - PostgreSQL + Redis 开发环境
- ✅ VitePress 文档 - 独立的文档项目 + 自动部署到 GitHub Pages

## 待实现功能

- 应用服务层（Application Layer）- 业务逻辑层
- 权限和角色管理（RBAC）- 基于角色的访问控制
- 结构化日志系统 - 使用 zap 或 zerolog
- 单元测试和集成测试 - 完整的测试覆盖
- API 自动文档 - Swagger/OpenAPI 规范自动生成
- 分布式追踪 - OpenTelemetry 集成
- 监控和指标 - Prometheus + Grafana
