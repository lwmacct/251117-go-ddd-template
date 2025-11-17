# Redis 实现说明

## 概述

本项目已实现完整的 Redis 集成，包括：
- Redis 客户端初始化和连接管理
- 缓存仓储接口（CacheRepository）
- 健康检查端点（包含 Redis 状态）
- RESTful API 示例（缓存 CRUD 操作）

## 快速开始

### 1. 启动 Redis 服务

使用 Docker Compose：

```bash
docker-compose up -d redis
```

或手动启动 Redis：

```bash
# 使用 Docker
docker run -d --name redis -p 6379:6379 redis:7-alpine

# 或使用本地 Redis
redis-server
```

### 2. 配置 Redis 连接

方式一：配置文件（`config.yaml`）
```yaml
data:
  redis_url: "redis://localhost:6379/0"
```

方式二：环境变量
```bash
export APP_DATA_REDIS_URL="redis://localhost:6379/0"
```

方式三：带密码的连接
```bash
# 配置文件
data:
  redis_url: "redis://:your-password@localhost:6379/0"

# 环境变量
export APP_DATA_REDIS_URL="redis://:your-password@localhost:6379/0"
```

### 3. 运行应用

```bash
# 使用 task
task go:run -- api

# 或直接运行
.local/bin/bd-vmalert api
```

## API 端点

### 健康检查（包含 Redis 状态）

```bash
curl http://localhost:8080/health
```

响应示例：
```json
{
  "status": "ok",
  "checks": {
    "redis": {
      "status": "healthy"
    }
  }
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

响应示例：
```json
{
  "key": "user:123",
  "value": {
    "age": 30,
    "name": "张三"
  }
}
```

### 删除缓存

```bash
curl -X DELETE http://localhost:8080/api/cache/user:123
```

## 代码结构

```
internal/
├── infrastructure/
│   └── redis/
│       ├── client.go           # Redis 客户端初始化
│       └── cache_repository.go # 缓存仓储接口和实现
├── adapters/
│   └── http/
│       ├── handler_health.go   # 健康检查处理器
│       └── handler_cache.go    # 缓存操作处理器
└── bootstrap/
    └── container.go            # 依赖注入容器（包含 Redis 客户端）
```

## 使用示例

### 在代码中使用 Redis

```go
import (
    "context"
    "time"
    redisinfra "github.com/lwmacct/251117-bd-vmalert/internal/infrastructure/redis"
)

// 从容器获取 Redis 客户端
redisClient := container.RedisClient

// 创建缓存仓储
cacheRepo := redisinfra.NewCacheRepository(redisClient)

// 设置缓存
ctx := context.Background()
data := map[string]interface{}{"name": "测试"}
err := cacheRepo.Set(ctx, "mykey", data, 5*time.Minute)

// 获取缓存
var result map[string]interface{}
err = cacheRepo.Get(ctx, "mykey", &result)

// 删除缓存
err = cacheRepo.Delete(ctx, "mykey")

// 检查键是否存在
exists, err := cacheRepo.Exists(ctx, "mykey")

// 分布式锁（仅当键不存在时设置）
ok, err := cacheRepo.SetNX(ctx, "lock:resource", "locked", 10*time.Second)
```

## CacheRepository 接口

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

## 注意事项

1. **自动 JSON 序列化**：CacheRepository 会自动将值序列化为 JSON 存储，读取时自动反序列化
2. **上下文超时**：所有操作都支持 context，可以设置超时和取消
3. **优雅关闭**：应用退出时会自动关闭 Redis 连接
4. **错误处理**：所有操作都返回详细的错误信息
5. **连接池**：go-redis 内置连接池管理，无需手动管理连接

## 扩展建议

1. **添加更多仓储**：参考 `cache_repository.go` 创建其他数据仓储
2. **实现分布式锁**：使用 `SetNX` 方法实现分布式锁机制
3. **缓存预热**：在应用启动时预加载热点数据
4. **缓存失效策略**：实现缓存失效和更新策略
5. **监控和指标**：添加 Redis 操作的监控和性能指标
