# Redis 集成

本项目使用 go-redis 实现 Redis 集成，提供缓存管理和分布式锁功能。

<!--TOC-->

## Table of Contents

- [功能特性](#功能特性) `:36+10`
- [快速开始](#快速开始) `:46+26`
  - [1. 启动 Redis](#1-启动-redis) `:48+6`
  - [2. 配置连接](#2-配置连接) `:54+12`
  - [3. 运行应用](#3-运行应用) `:66+6`
- [架构设计](#架构设计) `:72+26`
  - [目录结构](#目录结构) `:74+8`
  - [CacheRepository 接口](#cacherepository-接口) `:82+16`
- [API 端点](#api-端点) `:98+23`
  - [健康检查](#健康检查) `:100+7`
  - [缓存操作](#缓存操作) `:107+14`
- [使用模式](#使用模式) `:121+25`
  - [Cache-Aside（旁路缓存）](#cache-aside旁路缓存) `:123+5`
  - [分布式锁](#分布式锁) `:128+10`
  - [防止缓存穿透](#防止缓存穿透) `:138+4`
  - [防止缓存雪崩](#防止缓存雪崩) `:142+4`
- [性能优化](#性能优化) `:146+16`
  - [Pipeline 批量操作](#pipeline-批量操作) `:148+10`
  - [连接池配置](#连接池配置) `:158+4`
- [故障排查](#故障排查) `:162+18`
  - [连接失败](#连接失败) `:164+8`
  - [性能问题](#性能问题) `:172+8`
- [最佳实践](#最佳实践) `:180+16`
  - [扩展功能](#扩展功能) `:188+8`

<!--TOC-->

## 功能特性

- ✅ Redis 客户端初始化和连接管理
- ✅ CacheRepository 接口（Get/Set/Delete/Exists/SetNX）
- ✅ 自动 JSON 序列化/反序列化
- ✅ TTL 过期时间支持
- ✅ 分布式锁（SetNX）
- ✅ 健康检查端点
- ✅ 优雅关闭

## 快速开始

### 1. 启动 Redis

```bash
docker-compose up -d redis
```

### 2. 配置连接

```yaml
# config.yaml
data:
  redis:
    url: "redis://localhost:6379/0"
    # 带密码: "redis://:password@localhost:6379/0"
```

或使用环境变量：`APP_DATA_REDIS_URL`

### 3. 运行应用

```bash
task go:run -- api
```

## 架构设计

### 目录结构

```
internal/infrastructure/redis/
├── client.go           # Redis 客户端初始化
└── cache_repository.go # 缓存仓储接口和实现
```

### CacheRepository 接口

```go
type CacheRepository interface {
    Get(ctx context.Context, key string, dest interface{}) error
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
    SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
}
```

- 自动 JSON 序列化/反序列化
- 支持任意结构体类型
- `SetNX` 用于分布式锁

## API 端点

### 健康检查

```bash
curl http://localhost:8080/health
# 返回: {"status": "ok", "redis": "connected"}
```

### 缓存操作

```bash
# 设置缓存
curl -X POST http://localhost:8080/api/cache \
  -d '{"key": "user:123", "value": {"name": "张三"}, "ttl": 60}'

# 获取缓存
curl http://localhost:8080/api/cache/user:123

# 删除缓存
curl -X DELETE http://localhost:8080/api/cache/user:123
```

## 使用模式

### Cache-Aside（旁路缓存）

1. 读取：先查缓存 → 未命中则查数据库 → 写入缓存
2. 更新：更新数据库 → 删除缓存（下次查询重新加载）

### 分布式锁

```go
locked, _ := cacheRepo.SetNX(ctx, "lock:resource", "locked", 10*time.Second)
if locked {
    defer cacheRepo.Delete(ctx, "lock:resource")
    // 执行需要加锁的操作
}
```

### 防止缓存穿透

缓存空值（短 TTL），避免大量请求穿透到数据库。

### 防止缓存雪崩

为 TTL 添加随机抖动，避免同时过期。

## 性能优化

### Pipeline 批量操作

```go
pipe := redisClient.Pipeline()
for i := 0; i < 100; i++ {
    pipe.Set(ctx, fmt.Sprintf("key:%d", i), i, time.Hour)
}
pipe.Exec(ctx)
```

### 连接池配置

PoolSize=10, MinIdleConns=5, MaxRetries=3（可在初始化时配置）

## 故障排查

### 连接失败

```bash
docker ps | grep redis       # 检查是否运行
redis-cli ping               # 应返回 PONG
docker logs go-ddd-redis     # 查看日志
```

### 性能问题

```bash
redis-cli SLOWLOG GET 10     # 查看慢查询
redis-cli INFO memory        # 查看内存
redis-cli INFO clients       # 查看连接数
```

## 最佳实践

1. **合理设置 TTL**：热点数据长、冷数据短
2. **键命名规范**：使用冒号分隔，如 `user:123`、`session:abc`
3. **避免大 Value**：单个值不超过 10MB
4. **缓存失败降级**：缓存失败不应影响主流程
5. **内存管理**：设置 maxmemory 和淘汰策略（allkeys-lru）

### 扩展功能

- Redis Cluster / Sentinel 支持
- 发布/订阅
- Lua 脚本原子操作
- Sorted Set 排行榜

详见 go-redis 官方文档。
