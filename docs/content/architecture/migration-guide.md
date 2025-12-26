# 架构迁移指南

本文档记录了从遗留分层实现升级到 **DDD 四层架构 + CQRS 模式** 的过程。

> **状态**: 迁移已完成 (2025-11-19)

<!--TOC-->

## Table of Contents

- [重构概览](#重构概览) `:23+18`
  - [解决的问题](#解决的问题) `:25+9`
  - [迁移成果](#迁移成果) `:34+7`
- [迁移阶段](#迁移阶段) `:41+12`
- [完成模块](#完成模块) `:53+12`
- [最佳实践](#最佳实践) `:65+15`
  - [Use Case 命名](#use-case-命名) `:67+6`
  - [CQRS 适用场景](#cqrs-适用场景) `:73+7`
- [验证清单](#验证清单) `:80+7`

<!--TOC-->

## 重构概览

### 解决的问题

| 问题                    | 解决方案                     |
| ----------------------- | ---------------------------- |
| 缺少 Application 层     | 新增 `internal/application/` |
| 读写操作混合            | CQRS Repository 分离         |
| Domain 贫血模型         | 富领域模型 + 业务方法        |
| Infrastructure 职责过重 | 明确 Domain Service 接口     |

### 迁移成果

- ✅ Application 层: 30 Use Cases (18 Command + 12 Query)
- ✅ CQRS Repository: 8 Command + 8 Query
- ✅ 富领域模型: User、Role 等含业务行为
- ✅ 单一依赖注入容器

## 迁移阶段

| 阶段 | 内容                           | 参考文件                                              |
| ---- | ------------------------------ | ----------------------------------------------------- |
| 1    | 创建 Application 层结构        | `internal/application/*/`                             |
| 2    | 重构 Domain 层 (移除 GORM Tag) | `internal/domain/*/entity_*.go`                       |
| 3    | 实现 CQRS Repository           | `internal/infrastructure/persistence/*_repository.go` |
| 4    | 创建 Use Cases                 | `internal/application/*/cmd_*.go`, `qry_*.go`         |
| 5    | 重构 Infrastructure Service    | `internal/infrastructure/auth/service.go`             |
| 6    | 更新 Adapter 层                | `internal/adapters/http/handler/*.go`                 |
| 7    | 更新依赖注入容器               | `internal/bootstrap/container.go`                     |

## 完成模块

| 模块     | Command | Query | 状态 |
| -------- | ------- | ----- | ---- |
| Auth     | 3       | 1     | ✅   |
| User     | 5       | 5     | ✅   |
| Role     | 4       | 3     | ✅   |
| Menu     | 4       | 2     | ✅   |
| Setting  | 4       | 2     | ✅   |
| PAT      | 2       | 2     | ✅   |
| AuditLog | 0       | 2     | ✅   |

## 最佳实践

### Use Case 命名

- Command: `cmd_{动作}.go` → `Create{Entity}Command`
- Query: `qry_{动作}.go` → `Get{Entity}Query`
- Handler: `*_handler.go`

### CQRS 适用场景

- ✅ 读写频率差异大
- ✅ 需要独立优化读性能
- ✅ 审计/日志等只读场景
- ⚠️ 简单 CRUD 可选

## 验证清单

- [ ] 所有 Handler 只做 HTTP 转换
- [ ] 业务逻辑在 Application 层
- [ ] Domain 无 GORM 依赖
- [ ] Repository 读写分离
- [ ] `go build ./...` 无错误
