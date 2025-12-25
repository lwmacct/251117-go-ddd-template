# PostgreSQL 集成

本项目使用 GORM 实现 PostgreSQL 数据库集成，遵循 DDD + CQRS 模式。

<!--TOC-->

## Table of Contents

- [功能特性](#功能特性) `:30+10`
- [快速开始](#快速开始) `:40+27`
  - [1. 启动 PostgreSQL](#1-启动-postgresql) `:42+6`
  - [2. 配置连接](#2-配置连接) `:48+11`
  - [3. 运行应用](#3-运行应用) `:59+8`
- [架构设计](#架构设计) `:67+39`
  - [目录结构](#目录结构) `:69+14`
  - [CQRS 仓储接口](#cqrs-仓储接口) `:83+14`
  - [Domain 与 Model 分离](#domain-与-model-分离) `:97+9`
- [连接池配置](#连接池配置) `:106+14`
- [API 端点](#api-端点) `:120+13`
  - [健康检查](#健康检查) `:122+7`
  - [用户 CRUD](#用户-crud) `:129+4`
- [特性说明](#特性说明) `:133+9`
- [故障排查](#故障排查) `:142+18`
  - [连接失败](#连接失败) `:144+8`
  - [性能问题](#性能问题) `:152+8`
- [最佳实践](#最佳实践) `:160+7`

<!--TOC-->

## 功能特性

- ✅ 数据库连接管理和连接池配置
- ✅ 自动迁移支持
- ✅ CQRS 仓储分离（Command/Query）
- ✅ Domain 与 Model 分离（无 GORM Tag 污染）
- ✅ 软删除支持
- ✅ bcrypt 密码加密
- ✅ 健康检查（含连接池统计）

## 快速开始

### 1. 启动 PostgreSQL

```bash
docker-compose up -d postgres
```

### 2. 配置连接

```yaml
# config.yaml
data:
  pgsql:
    url: "postgresql://postgres:postgres@localhost:5432/app?sslmode=disable"
```

或使用环境变量：`APP_DATA_PGSQL_URL`

### 3. 运行应用

```bash
task go:run -- api
```

应用启动时自动执行数据库迁移。

## 架构设计

### 目录结构

```
internal/
├── domain/user/
│   ├── entity_user.go          # 领域模型（无 GORM Tag）
│   ├── command_repository.go   # 写接口
│   └── query_repository.go     # 读接口
├── infrastructure/persistence/
│   ├── user_model.go           # GORM Model + 映射函数
│   ├── user_command_repository.go
│   └── user_query_repository.go
```

### CQRS 仓储接口

**CommandRepository（写操作）**：

- `Create`, `Update`, `Delete`
- `AssignRoles`, `UpdatePassword`, `UpdateStatus`

**QueryRepository（读操作）**：

- `GetByID`, `GetByUsername`, `GetByEmail`
- `GetByIDWithRoles`（预加载关联）
- `List`, `Count`, `Search`
- `ExistsByUsername`, `ExistsByEmail`

### Domain 与 Model 分离

| 层级                 | 特点                                |
| -------------------- | ----------------------------------- |
| Domain Entity        | 仅业务字段，无 GORM Tag             |
| Infrastructure Model | 完整 GORM Tag（索引、约束、软删除） |

仓储实现负责 Entity ↔ Model 双向映射。

## 连接池配置

```go
// internal/infrastructure/database/connection.go
sqlDB.SetMaxOpenConns(25)              // 最大连接数
sqlDB.SetMaxIdleConns(10)              // 空闲连接数
sqlDB.SetConnMaxLifetime(5 * time.Minute)
```

| 环境 | MaxOpenConns | MaxIdleConns |
| ---- | ------------ | ------------ |
| 开发 | 10           | 5            |
| 生产 | 25-100       | 根据负载调整 |

## API 端点

### 健康检查

```bash
curl http://localhost:8080/health
# 返回: {"status": "ok", "database": "connected", "db_stats": {...}}
```

### 用户 CRUD

详细 API 请在服务运行后通过 Swagger UI (`/swagger/index.html`) 查看。

## 特性说明

| 特性         | 说明                                                   |
| ------------ | ------------------------------------------------------ |
| **自动迁移** | 启动时自动创建/更新表结构，不删除已存在的列            |
| **软删除**   | GORM `DeletedAt` 字段，查询自动排除已删除记录          |
| **密码加密** | bcrypt 加密，`User.HashPassword()` / `CheckPassword()` |
| **数据验证** | Gin binding 标签验证请求数据                           |

## 故障排查

### 连接失败

```bash
docker ps | grep postgres          # 检查是否运行
docker logs go-ddd-postgres        # 查看日志
psql "postgresql://postgres:postgres@localhost:5432/app"  # 测试连接
```

### 性能问题

```bash
curl http://localhost:8080/health | jq '.db_stats'  # 检查连接池状态
```

优化建议：添加索引、调整连接池、使用 `Select` 只查询需要的字段、使用 `Preload` 避免 N+1。

## 最佳实践

1. **使用事务**：多操作业务逻辑
2. **避免 N+1 查询**：使用 `Preload` 预加载关联
3. **合理索引**：为常用查询字段添加索引
4. **连接池管理**：根据负载调整大小
5. **定期维护**：VACUUM、ANALYZE
