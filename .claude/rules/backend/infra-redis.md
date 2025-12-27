---
paths:
  - "internal/infrastructure/redis/**/*.go"
---

# Redis Infrastructure 规范

## 核心职责

实现 Domain 层 `cache.CommandRepository` 和 `cache.QueryRepository` 接口，提供缓存服务。

## 文件结构

| 文件                  | 职责                         |
| --------------------- | ---------------------------- |
| `client.go`           | Redis 客户端初始化、连接管理 |
| `cache_repository.go` | 实现 Domain 缓存仓储接口     |
| `doc.go`              | 包文档                       |

## 客户端规范

```go
// client.go
// NewClient 创建并初始化 Redis 客户端
// redisURL 格式: redis://[:password@]host:port[/db]
func NewClient(ctx context.Context, redisURL string, enableTracing bool) (*redis.Client, error)

// Close 关闭 Redis 客户端连接
func Close(client *redis.Client) error

// HealthCheck 检查 Redis 连接健康状态
func HealthCheck(ctx context.Context, client *redis.Client) error
```

## 缓存仓储规范

```go
// cache_repository.go - 同时实现 Command 和 Query 接口
type cacheRepository struct {
    client    *redis.Client
    keyPrefix string // key 前缀，所有操作自动拼接
}

// 创建仓储时指定 key 前缀
func NewCacheCommandRepository(client *redis.Client, keyPrefix string) cache.CommandRepository
func NewCacheQueryRepository(client *redis.Client, keyPrefix string) cache.QueryRepository
```

## Key 命名规范

- 使用 `keyPrefix` 区分不同模块的缓存
- 示例：`user:profile:123`、`perm:user:456`

## 依赖方向

```
domain/cache.CommandRepository (接口)
domain/cache.QueryRepository   (接口)
              ↑
infrastructure/redis (实现)
```
