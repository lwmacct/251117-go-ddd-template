# 项目概览

基于 Go 的 DDD 模板应用，采用四层架构 + CQRS 模式。Monorepo 结构包含后端(Go)、前端(Vue 3)、文档(VitePress)。

<!--TOC-->

## Table of Contents

- [架构概览](#架构概览) `:15+16`
- [关键文件](#关键文件) `:31+9`
- [项目文档](#项目文档) `:40+6`

<!--TOC-->

## 架构概览

```
internal/
├── adapters/        # 适配器层 - HTTP Handler、中间件、路由
├── application/     # 应用层 - Use Cases (Command/Query Handler)
├── domain/          # 领域层 - 业务模型、Repository 接口
├── infrastructure/  # 基础设施层 - Repository 实现、数据库
├── bootstrap/       # 依赖注入容器
└── commands/        # CLI 命令
```

**依赖方向**: `Adapters → Application → Domain ← Infrastructure`

> 详细规范见 `.claude/rules/ddd-*.md`，编辑对应目录时自动加载。

## 关键文件

| 用途       | 文件                                             |
| ---------- | ------------------------------------------------ |
| 依赖注入   | `internal/bootstrap/container.go`                |
| 路由定义   | `internal/adapters/http/router.go`               |
| 配置管理   | `internal/infrastructure/config/config.go`       |
| 数据库迁移 | `internal/infrastructure/database/migrations.go` |

## 项目文档

VitePress 文档位于 `docs/` 目录：

- 文档索引: `docs/.vitepress/config.ts`
- 架构文档: `docs/architecture/ddd-cqrs.md`
