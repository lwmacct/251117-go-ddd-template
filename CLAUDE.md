# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此仓库中工作时提供指导。

## 📚 重要：如何查看项目文档

本项目拥有完整的 **VitePress 2.0 文档系统** (位于 `docs/` 目录) ，所有详细的架构、API、配置、开发指南等内容都在文档中维护。

### 查看文档结构

- 文档索引文件：`docs/.vitepress/config.ts`
- 此文件定义了完整的导航和侧边栏配置，包含所有可用的文档页面

### 当需要了解项目详细信息时：

1. 查看 `docs/.vitepress/config.ts` 了解有哪些文档
2. 在 `docs/` 目录下直接阅读对应的 Markdown 文件
3. 修改代码时，同步更新相关文档

## 项目概述

基于 Go 的 DDD (领域驱动设计) 模板应用，使用 Gin 提供 HTTP 服务，遵循整洁架构原则。

**技术栈**：

- 框架：Gin (HTTP)、urfave/cli v3 (CLI)
- 数据库：PostgreSQL + GORM
- 缓存：Redis
- 认证：JWT (golang-jwt/jwt/v5)
- 配置：Koanf

## 架构概览

本项目遵循 DDD (领域驱动设计) 和整洁架构原则。

**分层结构**：

- `internal/commands/` - CLI 命令 (入口点)
- `internal/adapters/` - 外部接口 (HTTP、gRPC 等)
- `internal/domain/` - 领域层 (业务逻辑)
- `internal/infrastructure/` - 技术实现 (数据库、Redis、配置等)
- `internal/bootstrap/` - 依赖注入容器

**关键设计**：

- 依赖注入容器 (`bootstrap.Container`)
- 仓储模式 (Repository Pattern)
- 配置系统 (Koanf，多层优先级)
- JWT 认证授权

> 📖 **详细架构说明**：查看文档 `/guide/architecture`

## 配置系统

配置优先级 (从低到高) ：

1. 默认值 → 2. 配置文件 → 3. 环境变量 (前缀 `APP_`) → 4. 命令行参数

环境变量示例：

```bash
APP_SERVER_ADDR=:8080
APP_DATA_PGSQL_URL=postgresql://user:pass@host:5432/db
APP_DATA_REDIS_URL=redis://localhost:6379/0
APP_JWT_SECRET=your-secret-key
```

**重要**：修改 `internal/infrastructure/config/config.go` 中的配置结构后，运行 `sync-config-example` 技能更新示例配置文件。

> 📖 **详细配置说明**：查看文档 `/guide/configuration`

## 扩展应用

添加新功能的快速参考：

1. **新 HTTP 端点**：

   - 创建 handler：`internal/adapters/http/handler/<name>.go`
   - 注册路由：`internal/adapters/http/router.go`

2. **新领域模型**：

   - 创建模型：`internal/domain/<name>/model.go`
   - 定义仓储接口：`internal/domain/<name>/repository.go`
   - 实现仓储：`internal/infrastructure/persistence/<name>_repository.go`
   - 注入依赖：`internal/bootstrap/container.go`

3. **新配置项**：

   - 更新：`internal/infrastructure/config/config.go`
   - 运行：`sync-config-example` 技能

4. **新 CLI 命令**：
   - 创建：`internal/commands/<name>/`
   - 注册：`main.go` 中的 `buildCommands()`

> 📖 **详细扩展指南**：查看文档 `/guide/architecture` 和 `/guide/contributing`

## 项目结构 (Monorepo)

```
.
├── internal/          # 后端核心代码 (Go)
├── web/               # 前端项目 (Vue 3，独立的 package.json)
├── docs/              # VitePress 文档 (独立的 package.json)
├── configs/           # 配置文件
├── docker-compose.yml # PostgreSQL + Redis
├── Taskfile.yaml      # 任务自动化
├── .air.toml          # 热重载配置
└── main.go            # 应用入口
```

## 已实现功能

✅ DDD 分层架构 + 整洁架构
✅ HTTP 服务器 (Gin) + 优雅关闭
✅ JWT 认证授权系统
✅ PostgreSQL (GORM ORM + 自动迁移)
✅ Redis 缓存 + 分布式锁
✅ 配置管理 (Koanf 多层优先级)
✅ 用户管理 (CRUD + 软删除 + 分页)
✅ 依赖注入容器
✅ 仓储模式
✅ 健康检查
✅ VitePress 文档系统

## 待实现功能

- 应用服务层 (Application Layer)
- 权限和角色管理 (RBAC)
- 结构化日志系统 (zap/zerolog)
- 单元测试和集成测试
- API 文档自动生成 (Swagger/OpenAPI)
- 分布式追踪 (OpenTelemetry)
- 监控和指标 (Prometheus + Grafana)

---

**记住：遇到问题或需要详细信息时，优先查看 VitePress 文档 (`docs/` 目录) ！**
