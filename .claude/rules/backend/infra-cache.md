---
paths:
  - "internal/infrastructure/cache/**/*.go"
---

# Cache Infrastructure 规范

## 核心职责

实现缓存服务接口，提供 Redis 缓存能力。**使用 RedisJSON 原生 JSON 类型存储**。

## 存储类型

**必须使用 RedisJSON 模块**（Redis Stack 默认包含）：

| 操作     | 命令                              | 说明                         |
| -------- | --------------------------------- | ---------------------------- |
| 写入     | `JSON.SET key $ value` + `EXPIRE` | Pipeline 执行保证原子性      |
| 读取     | `JSON.GET key $`                  | 返回数组包装 `[actual_data]` |
| 批量读取 | `JSON.MGET keys... $`             | 每个元素是 JSON 字符串或 nil |
| 删除     | `DEL key`                         | 不变                         |

**关键注意点**：`JSON.GET $` 返回数组包装，需要解包：

```go
var wrapper []YourDTO
json.Unmarshal([]byte(data), &wrapper)
if len(wrapper) > 0 {
    return wrapper[0], nil
}
```

## 接口位置原则

**根据缓存内容决定接口定义位置**：

| 缓存内容           | 接口定义位置          | 示例                         |
| ------------------ | --------------------- | ---------------------------- |
| Domain 实体        | `domain/cache/`       | `SettingCacheService`        |
| Application 层 DTO | `application/{模块}/` | `setting.SchemaCacheService` |

**判断依据**：

- Domain 层不应知道 Application 层 DTO 结构
- 缓存 Application DTO 时，接口必须定义在 Application 层

## 文件命名

| 文件类型     | 命名规范                  |
| ------------ | ------------------------- |
| 包文档       | `doc.go`（**必需**）      |
| 客户端       | `redis_client.go`         |
| 缓存服务实现 | `{模块}_cache_service.go` |

## 设计原则

### 1. 接口实现

```go
// Domain 实体缓存
func NewXxxCacheService(client *redis.Client, keyPrefix string) cache.XxxCacheService

// Application DTO 缓存
func NewXxxCacheService(client *redis.Client, keyPrefix string) appsetting.XxxCacheService
```

### 2. 写入模式（RedisJSON）

```go
// Pipeline 写入（推荐）
func (s *service) set(ctx context.Context, key string, data any) error {
    pipe := s.client.Pipeline()
    pipe.JSONSet(ctx, key, "$", data)
    pipe.Expire(ctx, key, ttl)
    _, err := pipe.Exec(ctx)
    return err
}
```

### 3. 读取模式（RedisJSON）

```go
func (s *service) get(ctx context.Context, key string) (*SomeDTO, error) {
    data, err := s.client.JSONGet(ctx, key, "$").Result()
    if errors.Is(err, redis.Nil) {
        return nil, nil // cache miss
    }

    // JSON.GET $ 返回数组包装
    var wrapper []SomeDTO
    if json.Unmarshal([]byte(data), &wrapper) != nil || len(wrapper) == 0 {
        _ = s.client.Del(ctx, key) // 损坏数据自动清除
        return nil, nil
    }
    return &wrapper[0], nil
}
```

### 4. 分布式锁

预热等竞争场景使用 SETNX + TTL：

```go
ok, _ := client.SetNX(ctx, lockKey, "1", 30*time.Second).Result()
```

## Key 命名规范

- 格式：`{prefix}{模块}:{标识}` 或 `{prefix}{模块}:{scope}:{id}:{key}`
- 预热标记：`{prefix}{模块}:_warmed_up`
- 预热锁：`{prefix}{模块}:_warmup_lock`

## DTO 规范

### Domain 实体缓存

缓存 DTO 独立定义，避免 Domain 实体添加 JSON tags：

```go
type xxxCacheDTO struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}

func toXxxCacheDTO(e *Entity) xxxCacheDTO
func (d xxxCacheDTO) toEntity() *Entity
```

### Application DTO 缓存

**直接序列化**（推荐，Application DTO 已有 JSON tags）：

```go
func (s *service) set(ctx context.Context, key string, dto []app.SomeDTO) error {
    pipe := s.client.Pipeline()
    pipe.JSONSet(ctx, key, "$", dto) // 直接传结构体
    pipe.Expire(ctx, key, ttl)
    _, err := pipe.Exec(ctx)
    return err
}
```

**独立缓存 DTO**（仅当需要隔离 Application DTO 变化时）：

```go
type someCacheDTO struct {
    Field string `json:"field"`
}
```

**选择依据**：

- Application DTO 稳定 → **直接序列化**（首选）
- 需要隔离变化（缓存兼容性） → 独立 DTO

## 方法排序

1. 构造函数
2. 导出方法（按接口顺序）
3. 未导出辅助方法
4. 接口实现检查 `var _ = (*impl)(nil)`

## 缓存回写策略

**使用同步回写**（推荐）：

```go
// Cache-Aside 模式：cache miss 后同步回写
func (r *cachedRepo) FindByKey(ctx context.Context, key string) (*Entity, error) {
    // 1. 查缓存
    cached, _ := r.cache.Get(ctx, key)
    if cached != nil {
        return cached, nil
    }

    // 2. 查数据库
    result, err := r.delegate.FindByKey(ctx, key)
    if err != nil || result == nil {
        return result, err
    }

    // 3. 同步回写缓存（无需版本控制）
    if err := r.cache.Set(ctx, result); err != nil {
        slog.Warn("cache set failed", "key", key, "err", err)
    }

    return result, nil
}
```

**为何不使用异步回写**：

- 同步回写延迟可忽略（Redis 写入 < 1ms）
- 避免竞态条件，无需版本控制
- 代码更简单，调试更容易
