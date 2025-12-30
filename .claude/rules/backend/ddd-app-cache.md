---
paths:
  - "internal/application/**/cache*.go"
---

# Application 层缓存接口规范

## 数据类型约束

**缓存接口仅使用 Application DTO**，禁止直接缓存 Domain 实体。

```go
// ✅ 正确
type UserCacheService interface {
    Get(ctx context.Context, userID uint) (*UserWithRolesDTO, error)
}

// ❌ 禁止
type UserCacheService interface {
    Get(ctx context.Context, userID uint) (*user.User, error)
}
```

## 接口命名

| 模式                        | 说明            | 示例                           |
| --------------------------- | --------------- | ------------------------------ |
| `{Entity}CacheService`      | 实体缓存        | `UserWithRolesCacheService`    |
| `{Entity}QueryCacheService` | Repository 缓存 | `UserSettingQueryCacheService` |

## 方法模式

| 操作类型 | 方法前缀                      | 示例                        |
| -------- | ----------------------------- | --------------------------- |
| 读取     | `Get` / `GetBy` / `GetByKeys` | `GetUserWithRoles`          |
| 写入     | `Set` / `SetBatch`            | `SetUserSettings`           |
| 删除     | `Delete` / `DeleteBy`         | `DeleteByUser`, `DeleteAll` |
| 失效     | `Invalidate`                  | `InvalidateUser`            |

## 返回约定

**默认模式**（简洁优先）：

- **缓存未命中**：返回 `nil, nil`（不返回错误）
- **缓存损坏**：自动清除后返回 `nil, nil`

```go
user, err := cache.GetUser(ctx, id)
if user == nil {
    // 回源查询数据库
}
```

**哨兵错误模式**（语义明确，可选）：

当需要区分「未命中」和「Redis 返回空值」时，定义哨兵错误：

```go
// cache.go
var ErrCacheMiss = errors.New("cache: miss")

// 调用方
user, err := cache.GetUser(ctx, id)
if errors.Is(err, ErrCacheMiss) {
    // 回源查询
} else if err != nil {
    // Redis 连接等真正的错误
}
```

> Go 标准库惯例：`sql.ErrNoRows`、`redis.Nil`

**选择建议**：项目内保持一致即可，两种模式各有优劣。
