---
paths:
  - "internal/infrastructure/cache/**/*.go"
---

# Cache Infrastructure 规范

使用 RedisJSON 存储，实现 Application 层定义的缓存接口。

> [!TIP]
> 环境已预装 Redis，可直接使用 `redis-cli -u $REDIS_URL` 连接。

## 架构决策

```
Application Layer（接口定义 + DTO）
        │ 实现
Infrastructure Layer（Redis 实现）
```

接口定义在 Application 层（而非 Domain 层）是有意的设计：

1. **缓存是技术优化**：没有缓存系统功能完全正常，不属于领域概念
2. **接口与消费者同层**：Application Handler 使用缓存，接口定义在同层更自然
3. **避免 Domain 污染**：Domain 层应保持纯粹，不关心缓存细节
4. **无循环依赖**：Infrastructure → Application 是单向依赖，通过 DI 注入

## 目录结构

```
internal/infrastructure/cache/
├── doc.go                        # 包文档
└── {module}_cache_service.go     # 缓存服务实现
```

**命名约定**：`{module}` 为领域模块名（如 `user`、`setting`）

## RedisJSON 要点

| 操作 | 命令                              | 注意                              |
| ---- | --------------------------------- | --------------------------------- |
| 写入 | `JSON.SET key $ value` + `EXPIRE` | Pipeline 执行                     |
| 读取 | `JSON.GET key $`                  | **返回数组包装**，需 `wrapper[0]` |

## Key 命名

格式：`{prefix}{模块}:{scope}:{id}`

**prefix 配置约定**：

- 有前缀：配置 `"app:"` → key = `app:user:perms:123`
- 无前缀：配置 `""` → key = `user:perms:123`

> prefix 应包含分隔符（如 `app:`），代码直接拼接，无需额外处理空值场景。

```bash
# 调试 (开发模式下前缀为 dev:)
redis-cli -u $REDIS_URL KEYS 'dev:*'
```

## 缓存回写

使用 **Cache-Aside 同步回写**（延迟可忽略，避免竞态）。

## 数据类型约束

**禁止在 Infrastructure 层自定义缓存类型**，必须直接序列化 Application DTO 或 Domain 实体。

```go
// ❌ 禁止：自定义 CacheDTO
type userCacheDTO struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
}

// ✅ 正确：直接序列化 Application DTO 或 Domain 实体
func (s *service) Set(ctx context.Context, user *user.User) error {
    return s.client.JSONSet(ctx, key, "$", user).Err()
}
```

> Domain 实体需有 `json` tags（见 ddd-domain.md）。
