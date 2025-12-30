---
paths:
  - "internal/infrastructure/cache/warmup/**/*.go"
---

# Warmup Infrastructure 规范

提供缓存预热服务，在应用启动时批量加载数据到 Redis。

**注意**：本包是 `cache/` 的子包，负责预热逻辑；Redis 具体实现位于父包 `cache/`。

## 核心原则

**预热的目的**：确保应用启动后，常用数据已在缓存中，减少首次请求的数据库压力。

**设计简化**：

- 无分布式锁：多实例同时预热是幂等操作
- 无预热标记：每次启动都执行预热，确保数据最新

## 目录结构

```
internal/infrastructure/cache/warmup/
├── doc.go             # 包文档（必需）
└── {module}_warmup.go # 预热服务
```

**命名约定**:

- `{module}` 为领域模块名（如 `setting`、`permission`）

## 预热流程

```go
func (w *Warmer) WarmUp(ctx context.Context) error {
    // 1. 从数据库加载所有数据
    data, err := w.repo.FindAll(ctx)
    if err != nil {
        return fmt.Errorf("failed to load data: %w", err)
    }

    // 2. 批量写入缓存
    if len(data) > 0 {
        if err := w.cache.SetAll(ctx, data); err != nil {
            return fmt.Errorf("failed to warmup cache: %w", err)
        }
    }

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
