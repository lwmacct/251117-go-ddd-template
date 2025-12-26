# 数据存储

本项目使用 PostgreSQL 作为主数据库，Redis 作为缓存层。

<!--TOC-->

## Table of Contents

- [PostgreSQL 集成](#postgresql-集成) `:28+58`
  - [架构设计](#架构设计) `:32+14`
  - [CQRS 仓储分离](#cqrs-仓储分离) `:46+6`
  - [Domain 与 Model 分离](#domain-与-model-分离) `:52+7`
  - [连接池配置](#连接池配置) `:59+9`
  - [配置](#配置) `:68+10`
  - [最佳实践](#最佳实践) `:78+8`
- [Redis 集成](#redis-集成) `:86+45`
  - [CacheRepository 接口](#cacherepository-接口) `:90+12`
  - [配置](#配置-1) `:102+10`
  - [使用模式](#使用模式) `:112+11`
  - [最佳实践](#最佳实践-1) `:123+8`
- [健康检查](#健康检查) `:131+7`
- [故障排查](#故障排查) `:138+17`
  - [数据库连接失败](#数据库连接失败) `:140+8`
  - [Redis 连接失败](#redis-连接失败) `:148+7`

<!--TOC-->

## PostgreSQL 集成

使用 GORM 实现，遵循 DDD + CQRS 模式。

### 架构设计

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

### CQRS 仓储分离

**CommandRepository（写操作）**: `Create`, `Update`, `Delete`

**QueryRepository（读操作）**: `GetByID`, `GetByUsername`, `List`, `Count`, `Search`

### Domain 与 Model 分离

| 层级                 | 特点                                |
| -------------------- | ----------------------------------- |
| Domain Entity        | 仅业务字段，无 GORM Tag             |
| Infrastructure Model | 完整 GORM Tag（索引、约束、软删除） |

### 连接池配置

**配置文件**: `internal/infrastructure/database/connection.go`

| 环境 | MaxOpenConns | MaxIdleConns |
| ---- | ------------ | ------------ |
| 开发 | 10           | 5            |
| 生产 | 25-100       | 根据负载调整 |

### 配置

```yaml
data:
  pgsql:
    url: "postgresql://postgres:postgres@localhost:5432/app?sslmode=disable"
```

环境变量: `APP_DATA_PGSQL_URL`

### 最佳实践

1. **使用事务**: 多操作业务逻辑
2. **避免 N+1 查询**: 使用 `Preload` 预加载关联
3. **合理索引**: 为常用查询字段添加索引
4. **连接池管理**: 根据负载调整大小
5. **定期维护**: VACUUM、ANALYZE

## Redis 集成

使用 go-redis 实现，提供缓存管理和分布式锁功能。

### CacheRepository 接口

| 方法     | 说明         |
| -------- | ------------ |
| `Get`    | 获取缓存     |
| `Set`    | 设置缓存     |
| `Delete` | 删除缓存     |
| `Exists` | 检查是否存在 |
| `SetNX`  | 分布式锁     |

**实现文件**: `internal/infrastructure/redis/cache_repository.go`

### 配置

```yaml
data:
  redis:
    url: "redis://localhost:6379/0"
```

环境变量: `APP_DATA_REDIS_URL`

### 使用模式

**Cache-Aside（旁路缓存）**:

1. 读取: 先查缓存 → 未命中则查数据库 → 写入缓存
2. 更新: 更新数据库 → 删除缓存

**防止缓存穿透**: 缓存空值（短 TTL）

**防止缓存雪崩**: 为 TTL 添加随机抖动

### 最佳实践

1. **合理设置 TTL**: 热点数据长、冷数据短
2. **键命名规范**: 使用冒号分隔，如 `user:123`
3. **避免大 Value**: 单个值不超过 10MB
4. **缓存失败降级**: 缓存失败不应影响主流程
5. **内存管理**: 设置 maxmemory 和淘汰策略

## 健康检查

```bash
curl http://localhost:8080/health
# 返回: {"status": "ok", "database": "connected", "redis": "connected"}
```

## 故障排查

### 数据库连接失败

```bash
docker ps | grep postgres
docker logs go-ddd-postgres
psql "postgresql://postgres:postgres@localhost:5432/app"
```

### Redis 连接失败

```bash
docker ps | grep redis
redis-cli ping
docker logs go-ddd-redis
```
