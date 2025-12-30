---
paths:
  - "internal/infrastructure/cache/**/*.go"
  - "internal/application/**/*cache*.go"
---

# Cache Infrastructure 规范

使用 RedisJSON 存储，缓存只在 Application ↔ Infrastructure 层流转。

> [!TIP]
> 环境已预装 Redis，可直接使用 `redis-cli -u $REDIS_URL` 连接。

## 缓存层次原则

```
Application Layer（接口定义 + DTO/实体缓存）
        │ 实现
Infrastructure Layer（Redis 实现）
```

- ❌ Domain 层不定义缓存接口
- ✅ Application 层统一定义缓存接口（含 DTO 和领域实体）
- ✅ Repository 装饰器使用 Application 层缓存服务

## RedisJSON 要点

| 操作 | 命令                              | 注意                              |
| ---- | --------------------------------- | --------------------------------- |
| 写入 | `JSON.SET key $ value` + `EXPIRE` | Pipeline 执行                     |
| 读取 | `JSON.GET key $`                  | **返回数组包装**，需 `wrapper[0]` |

## 目录结构

**Application 层**（接口定义）：

```
internal/application/{module}/
└── cache.go                    # 缓存接口（或 cache_{entity}.go）
```

**Infrastructure 层**（实现）：

```
internal/infrastructure/cache/
└── {module}_cache_service.go   # Redis 缓存实现
```

**命名约定**:

- `{module}` 为领域模块名（如 `user`、`setting`）
- `{entity}` 用于多缓存场景（如 `cache_user.go`、`cache_role.go`）

## Key 命名

格式：`{prefix}{模块}:{scope}:{id}:{key}`

```bash
# 调试
redis-cli -u $REDIS_URL KEYS 'dev:*'
```

## DTO 规范

- **Domain 实体**：定义独立 `xxxCacheDTO`（避免实体加 JSON tags）
- **Application DTO**：直接序列化（已有 JSON tags）

## 缓存回写

使用 **Cache-Aside 同步回写**（延迟可忽略，避免竞态）。
