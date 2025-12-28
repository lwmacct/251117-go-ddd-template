---
paths:
  - "internal/infrastructure/warmup/**/*.go"
---

# Warmup Infrastructure 规范

提供缓存预热服务，在应用启动时批量加载数据到 Redis。

**注意**：本包负责预热逻辑，Redis 具体实现位于 `infrastructure/cache/`。

## 核心原则

**预热的唯一目的**：防止多实例启动时重复写入缓存。

**预热不是**：

- 缓存有效性的判断依据
- 决定是否查询数据库的条件
- 空数据合法性的证明

## 文件命名

| 文件类型 | 命名规范             |
| -------- | -------------------- |
| 包文档   | `doc.go`（**必需**） |
| 预热服务 | `{模块}_warmup.go`   |

## 多实例预热安全

使用分布式锁 + 双重检查：

```go
func (w *Warmer) WarmUp(ctx context.Context) error {
    // 1. 尝试获取锁
    acquired, release := w.cache.TryAcquireWarmupLock(ctx)
    if !acquired {
        return w.waitForWarmUp(ctx)
    }
    defer release()

    // 2. 双重检查
    if w.cache.IsWarmedUp(ctx) {
        return nil
    }

    // 3. 执行预热
    data, _ := w.repo.FindAll(ctx)
    w.cache.SetAll(ctx, data)
    w.cache.SetWarmedUp(ctx)
    return nil
}
```

## TTL 一致性警告

| 缓存类型 | 典型 TTL | 风险     |
| -------- | -------- | -------- |
| 实体缓存 | 10 分钟  | 先过期   |
| 预热标记 | 24 小时  | 仍然存在 |

**后果**：实体缓存过期后，`IsWarmedUp()` 仍返回 true，空缓存被误判为有效数据。

## 禁止事项

❌ **禁止信任 `IsWarmedUp + 空缓存` 组合**：

```go
// ❌ 错误：空缓存 + 已预热 = 有效空数据
if len(cachedMap) == 0 && cache.IsWarmedUp(ctx) {
    return []Entity{}, true // 危险！实体可能已过期
}

// ✅ 正确：空缓存总是回退到数据库
if len(cachedMap) == 0 {
    return nil, false // 触发数据库查询
}
```

## 依赖注入原则

预热服务依赖**原始仓储**（非缓存装饰器），避免循环依赖：

```go
func NewWarmer(
    repo domain.QueryRepository,  // 原始仓储
    cache cache.CacheService,
) *Warmer
```

## 失败降级

预热失败不阻塞启动，降级为惰性加载：

```go
if err := warmer.WarmUp(ctx); err != nil {
    slog.Warn("warmup failed, using lazy loading", "err", err)
}
```

## 方法排序

1. 构造函数
2. 导出方法（`WarmUp`, `WarmUpWithTimeout`）
3. 未导出辅助方法
