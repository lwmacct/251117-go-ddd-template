# Redis 集成

本项目实现了完整的 Redis 集成，提供缓存管理、分布式锁等功能。

## 功能特性

- ✅ Redis 客户端初始化和连接管理
- ✅ 缓存仓储接口（CacheRepository）
- ✅ 自动 JSON 序列化/反序列化
- ✅ 支持 TTL 过期时间
- ✅ 分布式锁（SetNX）
- ✅ 健康检查端点
- ✅ RESTful API 示例（缓存 CRUD 操作）
- ✅ 优雅关闭

## 快速开始

### 1. 启动 Redis 服务

使用 Docker Compose（推荐）：

```bash
docker-compose up -d redis
```

或手动启动：

```bash
# 使用 Docker
docker run -d --name redis -p 6379:6379 redis:7-alpine

# 或使用本地 Redis
redis-server
```

### 2. 配置 Redis 连接

**配置文件方式** (`config.yaml`):

```yaml
data:
  redis:
    url: "redis://localhost:6379/0"
```

**环境变量方式**:

```bash
export APP_DATA_REDIS_URL="redis://localhost:6379/0"
```

**带密码的连接**:

```yaml
# 配置文件
data:
  redis:
    url: "redis://:your-password@localhost:6379/0"
```

```bash
# 环境变量
export APP_DATA_REDIS_URL="redis://:your-password@localhost:6379/0"
```

### 3. 运行应用

```bash
# 使用 Task
task go:run -- api

# 或直接运行
.local/bin/251117-go-ddd-template api
```

## 架构设计

### 代码结构

```
internal/
├── infrastructure/
│   └── redis/
│       ├── client.go           # Redis 客户端初始化
│       └── cache_repository.go # 缓存仓储接口和实现
├── adapters/
│   └── http/
│       └── handler/
│           ├── health.go       # 健康检查处理器
│           └── cache.go        # 缓存操作处理器
└── bootstrap/
    └── container.go            # 依赖注入容器
```

### CacheRepository 接口

```go
type CacheRepository interface {
    // Get 获取缓存值并反序列化
    Get(ctx context.Context, key string, dest interface{}) error

    // Set 序列化并设置缓存值
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

    // Delete 删除缓存值
    Delete(ctx context.Context, key string) error

    // Exists 检查键是否存在
    Exists(ctx context.Context, key string) (bool, error)

    // SetNX 仅当键不存在时设置值（用于分布式锁）
    SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
}
```

## API 端点

### 健康检查

```bash
curl http://localhost:8080/health
```

**响应示例：**

```json
{
  "status": "ok",
  "database": "connected",
  "redis": "connected"
}
```

### 设置缓存

```bash
curl -X POST http://localhost:8080/api/cache \
  -H "Content-Type: application/json" \
  -d '{
    "key": "user:123",
    "value": {"name": "张三", "age": 30},
    "ttl": 60
  }'
```

### 获取缓存

```bash
curl http://localhost:8080/api/cache/user:123
```

**响应示例：**

```json
{
  "key": "user:123",
  "value": {
    "name": "张三",
    "age": 30
  }
}
```

### 删除缓存

```bash
curl -X DELETE http://localhost:8080/api/cache/user:123
```

## 使用示例

### 基本操作

```go
import (
    "context"
    "time"
    redisinfra "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/redis"
)

// 从容器获取 Redis 客户端
redisClient := container.RedisClient

// 创建缓存仓储
cacheRepo := redisinfra.NewCacheRepository(redisClient)

ctx := context.Background()

// 1. 设置缓存
data := map[string]interface{}{
    "name": "张三",
    "age":  30,
}
err := cacheRepo.Set(ctx, "user:123", data, 5*time.Minute)

// 2. 获取缓存
var result map[string]interface{}
err = cacheRepo.Get(ctx, "user:123", &result)
fmt.Printf("Name: %s, Age: %.0f\n", result["name"], result["age"])

// 3. 检查键是否存在
exists, err := cacheRepo.Exists(ctx, "user:123")
if exists {
    fmt.Println("Key exists")
}

// 4. 删除缓存
err = cacheRepo.Delete(ctx, "user:123")
```

### 缓存结构体

```go
type User struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

// 设置结构体到缓存
user := User{ID: 1, Username: "zhangsan", Email: "zhangsan@example.com"}
err := cacheRepo.Set(ctx, "user:1", user, 10*time.Minute)

// 从缓存获取结构体
var cachedUser User
err = cacheRepo.Get(ctx, "user:1", &cachedUser)
```

### 分布式锁

```go
// 尝试获取锁
lockKey := "lock:resource:123"
locked, err := cacheRepo.SetNX(ctx, lockKey, "locked", 10*time.Second)

if locked {
    defer cacheRepo.Delete(ctx, lockKey) // 释放锁

    // 执行需要加锁的操作
    fmt.Println("Lock acquired, processing...")
    time.Sleep(2 * time.Second)
} else {
    fmt.Println("Failed to acquire lock")
}
```

### 缓存更新模式

#### 1. Cache-Aside（旁路缓存）

```go
func GetUser(ctx context.Context, userID uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", userID)

    // 1. 尝试从缓存获取
    var user User
    err := cacheRepo.Get(ctx, cacheKey, &user)
    if err == nil {
        return &user, nil // 缓存命中
    }

    // 2. 缓存未命中，从数据库查询
    user, err = userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存
    cacheRepo.Set(ctx, cacheKey, user, 10*time.Minute)

    return user, nil
}

func UpdateUser(ctx context.Context, user *User) error {
    // 1. 更新数据库
    err := userRepo.Update(ctx, user)
    if err != nil {
        return err
    }

    // 2. 删除缓存（下次查询时重新加载）
    cacheKey := fmt.Sprintf("user:%d", user.ID)
    cacheRepo.Delete(ctx, cacheKey)

    return nil
}
```

#### 2. Write-Through（写穿）

```go
func UpdateUser(ctx context.Context, user *User) error {
    // 1. 更新数据库
    err := userRepo.Update(ctx, user)
    if err != nil {
        return err
    }

    // 2. 更新缓存
    cacheKey := fmt.Sprintf("user:%d", user.ID)
    err = cacheRepo.Set(ctx, cacheKey, user, 10*time.Minute)

    return err
}
```

## 常用场景

### 1. 会话管理

```go
// 存储会话
sessionID := uuid.New().String()
session := map[string]interface{}{
    "user_id":  1,
    "username": "zhangsan",
    "login_at": time.Now(),
}
err := cacheRepo.Set(ctx, "session:"+sessionID, session, 24*time.Hour)

// 获取会话
var session map[string]interface{}
err = cacheRepo.Get(ctx, "session:"+sessionID, &session)

// 删除会话（登出）
err = cacheRepo.Delete(ctx, "session:"+sessionID)
```

### 2. API 限流

```go
func RateLimit(ctx context.Context, userID uint, limit int, window time.Duration) (bool, error) {
    key := fmt.Sprintf("ratelimit:user:%d", userID)

    // 获取当前计数
    var count int
    err := cacheRepo.Get(ctx, key, &count)
    if err != nil {
        // 首次请求
        count = 0
    }

    if count >= limit {
        return false, nil // 超过限制
    }

    // 增加计数
    count++
    err = cacheRepo.Set(ctx, key, count, window)

    return true, err
}

// 使用
allowed, _ := RateLimit(ctx, 1, 100, 1*time.Minute) // 每分钟100次
if !allowed {
    return errors.New("rate limit exceeded")
}
```

### 3. 防止缓存穿透

```go
func GetUser(ctx context.Context, userID uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", userID)

    // 尝试从缓存获取
    var user User
    err := cacheRepo.Get(ctx, cacheKey, &user)
    if err == nil {
        // 检查是否是空值标记
        if user.ID == 0 {
            return nil, errors.New("user not found")
        }
        return &user, nil
    }

    // 从数据库查询
    user, err = userRepo.FindByID(ctx, userID)
    if err != nil {
        // 缓存空值，防止穿透
        emptyUser := User{ID: 0}
        cacheRepo.Set(ctx, cacheKey, emptyUser, 5*time.Minute)
        return nil, err
    }

    // 缓存正常数据
    cacheRepo.Set(ctx, cacheKey, user, 10*time.Minute)
    return user, nil
}
```

### 4. 防止缓存雪崩

```go
func GetUserWithJitter(ctx context.Context, userID uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", userID)

    var user User
    err := cacheRepo.Get(ctx, cacheKey, &user)
    if err == nil {
        return &user, nil
    }

    user, err = userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, err
    }

    // 添加随机抖动，避免同时过期
    baseTTL := 10 * time.Minute
    jitter := time.Duration(rand.Intn(60)) * time.Second
    ttl := baseTTL + jitter

    cacheRepo.Set(ctx, cacheKey, user, ttl)
    return user, nil
}
```

### 5. 排行榜（Sorted Set）

虽然 CacheRepository 不直接支持 Sorted Set，但可以直接使用 Redis 客户端：

```go
import "github.com/redis/go-redis/v9"

// 添加分数
redisClient.ZAdd(ctx, "leaderboard", redis.Z{
    Score:  100.0,
    Member: "user:1",
})

// 获取排名
rank, _ := redisClient.ZRank(ctx, "leaderboard", "user:1").Result()

// 获取 Top 10
top10, _ := redisClient.ZRevRangeWithScores(ctx, "leaderboard", 0, 9).Result()
```

## 性能优化

### 1. Pipeline（批量操作）

```go
// 使用 Pipeline 批量执行命令
pipe := redisClient.Pipeline()

for i := 0; i < 100; i++ {
    key := fmt.Sprintf("key:%d", i)
    pipe.Set(ctx, key, i, time.Hour)
}

_, err := pipe.Exec(ctx)
```

### 2. 连接池配置

```go
// 在初始化 Redis 客户端时配置
opt := &redis.Options{
    Addr:         "localhost:6379",
    PoolSize:     10,  // 连接池大小
    MinIdleConns: 5,   // 最小空闲连接
    MaxRetries:   3,   // 最大重试次数
}
client := redis.NewClient(opt)
```

## 故障排查

### 连接失败

```bash
# 检查 Redis 是否运行
docker ps | grep redis

# 查看 Redis 日志
docker logs go-ddd-redis

# 测试连接
redis-cli ping
# 应返回 PONG

# 使用 redis-cli 连接
redis-cli -h localhost -p 6379
> PING
> INFO
```

### 性能问题

```bash
# 查看慢查询
redis-cli SLOWLOG GET 10

# 查看内存使用
redis-cli INFO memory

# 查看连接数
redis-cli INFO clients

# 监控实时命令
redis-cli MONITOR
```

### 键过期问题

```bash
# 检查键的 TTL
redis-cli TTL key_name

# 查看所有键（小心，生产环境慎用）
redis-cli KEYS *

# 更好的方式：使用 SCAN
redis-cli SCAN 0 MATCH user:* COUNT 100
```

## 最佳实践

1. **合理设置 TTL**
   - 热点数据：较长 TTL（如 1 小时）
   - 一般数据：中等 TTL（如 10 分钟）
   - 冷数据：较短 TTL（如 1 分钟）

2. **键命名规范**
   - 使用冒号分隔：`user:123`、`session:abc`
   - 包含类型信息：`cache:user:123`
   - 避免过长的键名

3. **避免大 Value**
   - 单个值不超过 10MB
   - 大对象拆分存储
   - 使用压缩（如 JSON → MessagePack）

4. **监控和告警**
   - 监控内存使用率
   - 监控命中率
   - 设置内存淘汰策略

5. **数据一致性**
   - 使用 Cache-Aside 模式
   - 更新时先更新数据库
   - 删除缓存或设置短 TTL

## 扩展功能

### 1. Redis Cluster 支持

```go
import "github.com/redis/go-redis/v9"

clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
    Addrs: []string{
        "localhost:7000",
        "localhost:7001",
        "localhost:7002",
    },
})
```

### 2. Redis 哨兵模式

```go
sentinelClient := redis.NewFailoverClient(&redis.FailoverOptions{
    MasterName:    "mymaster",
    SentinelAddrs: []string{"localhost:26379", "localhost:26380"},
})
```

### 3. 发布/订阅

```go
// 订阅
pubsub := redisClient.Subscribe(ctx, "notifications")
ch := pubsub.Channel()

go func() {
    for msg := range ch {
        fmt.Println("Received:", msg.Payload)
    }
}()

// 发布
redisClient.Publish(ctx, "notifications", "Hello, World!")
```

### 4. Lua 脚本

```go
// 原子性操作
script := redis.NewScript(`
    local current = redis.call("GET", KEYS[1])
    if tonumber(current) >= tonumber(ARGV[1]) then
        return redis.call("DECRBY", KEYS[1], ARGV[1])
    else
        return -1
    end
`)

result, err := script.Run(ctx, redisClient, []string{"balance:user:1"}, 100).Result()
```

## 注意事项

1. **自动 JSON 序列化**
   - CacheRepository 自动将值序列化为 JSON
   - 读取时自动反序列化
   - 确保类型匹配

2. **上下文超时**
   - 所有操作都支持 context
   - 建议设置超时：`ctx, cancel := context.WithTimeout(ctx, 2*time.Second)`

3. **优雅关闭**
   - 应用退出时自动关闭 Redis 连接
   - 在 `container.Close()` 中处理

4. **错误处理**
   - 缓存失败不应影响主流程
   - 降级处理：缓存失败时直接查数据库

5. **内存管理**
   - 设置 maxmemory 限制
   - 配置淘汰策略（如 allkeys-lru）

## 下一步

- 了解 [认证授权](/guide/authentication)
- 学习 [PostgreSQL 集成](/guide/postgresql)
- 查看 [项目架构](/guide/architecture)
- 探索 [配置系统](/guide/configuration)
