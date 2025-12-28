---
paths:
  - "internal/bootstrap/**/*.go"
---

# Bootstrap 依赖注入规范

## Container 字段顺序

```go
type Container struct {
    // 1. 基础设施（DB, Redis, EventBus）
    // 2. Domain Services
    // 3. Repositories（按模块分组）
    // 4. Use Case Handlers（按模块分组）
    // 5. HTTP Handlers
}
```

## 初始化顺序

```go
func NewContainer(cfg *Config) (*Container, error) {
    // 1️⃣ 基础设施
    // 2️⃣ 缓存服务
    // 3️⃣ Repositories
    // 4️⃣ Domain Services
    // 5️⃣ Use Case Handlers
    // 6️⃣ HTTP Handlers
    // 7️⃣ Router
}
```

## 缓存服务设计规范

### 抽象层次原则

| 层级        | 缓存接口位置    | 存储内容        | 使用者              |
| ----------- | --------------- | --------------- | ------------------- |
| Domain      | `domain/cache/` | Domain 实体     | Repository 装饰器   |
| Application | `domain/cache/` | DTO（合并结果） | Application Handler |

**规则**：

- ❌ Repository 层禁止直接依赖 `*redis.Client`
- ✅ 必须通过 `domain/cache/` 定义的接口访问缓存
- ✅ 缓存服务实现放在 `infrastructure/redis/`

### 工厂函数签名一致性

所有带缓存的仓储工厂函数签名必须统一：

```go
// ✅ 正确 - 仅依赖抽象接口
func newXxxRepositoriesWithCache(db *gorm.DB, cacheServices *CacheServicesModule)

// ❌ 禁止 - 直接依赖基础设施
func newXxxRepositoriesWithCache(db *gorm.DB, redisClient *redis.Client, keyPrefix string, ...)
```

### CacheServicesModule 命名规范

```go
type CacheServicesModule struct {
    // 格式：{Entity} 或 {Entity}{Layer}
    Setting          cache.SettingCacheService           // Domain 实体缓存
    UserSettingQuery cache.UserSettingQueryCacheService  // Domain 实体缓存（Query 层）
    UserSetting      cache.UserSettingCacheService       // Application DTO 缓存
    Permission       cache.PermissionCacheService
}
```

### 新增缓存模块检查清单

1. [ ] 确定缓存内容的抽象层次（Domain 实体 vs Application DTO）
2. [ ] 在 `domain/cache/` 定义接口
3. [ ] 在 `infrastructure/redis/` 实现接口
4. [ ] 在 `CacheServicesModule` 添加字段
5. [ ] 在 `init_cache.go` 初始化服务
6. [ ] 仓储工厂函数仅依赖 `CacheServicesModule`

## 新增模块检查清单

1. Container 结构体添加字段
2. 创建 Command/Query Repository
3. 创建 Use Case Handlers
4. 创建 HTTP Handler
5. 注册路由
