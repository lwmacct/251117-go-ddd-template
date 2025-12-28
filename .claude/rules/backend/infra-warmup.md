---
paths:
  - "internal/infrastructure/warmup/**/*.go"
---

# Warmup Infrastructure 规范

## 核心职责

提供缓存预热服务，在应用启动时批量加载数据到 Redis。

**注意**：本包负责预热逻辑，Redis 具体实现位于 `infrastructure/cache/`。

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
