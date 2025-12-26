# Go DDD Template

基于领域驱动设计（DDD）和 CQRS 模式的企业级应用模板。

<!--TOC-->

## Table of Contents

- [技术栈](#技术栈) `:19+9`
- [快速开始](#快速开始) `:28+18`
- [项目结构](#项目结构) `:46+13`
- [核心功能](#核心功能) `:59+7`
- [文档](#文档) `:66+11`
- [开发命令](#开发命令) `:77+9`
- [License](#license) `:86+3`

<!--TOC-->

## 技术栈

| 后端                     | 前端                         | 文档      |
| ------------------------ | ---------------------------- | --------- |
| Go 1.25 + Gin + GORM     | Vue 3 + Vite + Vuetify       | VitePress |
| PostgreSQL + Redis       | Pinia + Vue Router           |           |
| JWT (golang-jwt) + Koanf | TypeScript + ESLint          |           |
| Swagger (swag)           | openapi-generator (类型同步) |           |

## 快速开始

```bash
# 1. 启动依赖服务
docker-compose up -d

# 2. 数据库迁移
go run main.go migrate up

# 3. 填充种子数据
go run main.go seed

# 4. 启动服务
air  # 或 go run main.go api
```

**预置账号**: `admin / password123`

## 项目结构

```
internal/
├── adapters/        # HTTP Handler、中间件
├── application/     # Use Cases (Command/Query)
├── domain/          # 业务模型、Repository 接口
├── infrastructure/  # 数据访问、外部服务
└── bootstrap/       # 依赖注入
```

**依赖方向**: `Adapters → Application → Domain ← Infrastructure`

## 核心功能

- JWT + PAT 双重认证
- RBAC 三段式权限 (`resource:action:scope`)
- 审计日志
- PostgreSQL + Redis

## 文档

| 文档        | 位置                                    |
| ----------- | --------------------------------------- |
| 快速入门    | `docs/content/getting-started.md`       |
| DDD 架构    | `docs/content/architecture/ddd-cqrs.md` |
| AI 开发规范 | `.claude/rules/`                        |
| API 文档    | 运行后访问 `/swagger/index.html`        |

启动文档服务: `cd docs && pnpm dev`

## 开发命令

```bash
task go:run -- api       # 启动 API
task go:run -- migrate up # 数据库迁移
task go:build            # 构建二进制
golangci-lint run        # 代码检查
```

## License

MIT
