---
paths:
  - "internal/infrastructure/redis/**/*.go"
---

# Redis Infrastructure 规范

## 核心职责

实现缓存服务接口，提供 Redis 缓存能力。

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
| 客户端       | `client.go`               |
| 缓存服务实现 | `{模块}_cache_service.go` |

## 设计原则

### 1. 接口实现

```go
// Domain 实体缓存
func NewXxxCacheService(client *redis.Client, keyPrefix string) cache.XxxCacheService

// Application DTO 缓存
func NewXxxCacheService(client *redis.Client, keyPrefix string) appsetting.XxxCacheService
```

### 2. 版本化写入（防 Stale Cache）

异步回写场景使用 Lua 脚本原子比较版本号：

```lua
local current = redis.call('GET', KEYS[1])
if current then
    local data = cjson.decode(current)
    if data.v >= tonumber(ARGV[2]) then return 0 end
end
redis.call('SETEX', KEYS[1], ARGV[3], ARGV[1])
return 1
```

### 3. 分布式锁

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
    ID      uint   `json:"id"`
    Version int64  `json:"v"` // 版本号
}

func toXxxCacheDTO(e *Entity) xxxCacheDTO
func (d xxxCacheDTO) toEntity() *Entity
```

### Application DTO 缓存

**两种方案**：

1. **直接序列化**（推荐，如 Application DTO 已有合适 JSON tags）：

   ```go
   func (s *service) set(ctx context.Context, key string, dto []app.SomeDTO) error {
       data, _ := json.Marshal(dto)
       return s.client.Set(ctx, key, data, ttl).Err()
   }
   ```

2. **独立缓存 DTO**（需要隔离变化或版本控制时）：

   ```go
   // 内部未导出 DTO（隔离 Application DTO 变更）
   type someCacheDTO struct {
       Field string `json:"field"`
   }

   func toSomeCacheDTO(dto app.SomeDTO) someCacheDTO
   func fromSomeCacheDTO(c someCacheDTO) app.SomeDTO
   ```

**选择依据**：

- Application DTO 稳定 → 直接序列化
- 需要缓存版本控制 → 独立 DTO + 版本号
- 需要隔离变化（缓存兼容性） → 独立 DTO

## 方法排序

1. 构造函数
2. 导出方法（按接口顺序）
3. 未导出辅助方法
4. 接口实现检查 `var _ = (*impl)(nil)`

## 异步操作

使用独立 context 避免请求取消中断操作：

```go
go func() { //nolint:contextcheck
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    // ...
}()
```
