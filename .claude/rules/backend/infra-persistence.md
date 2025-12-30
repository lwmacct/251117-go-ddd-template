---
paths:
  - "internal/infrastructure/persistence/**/*.go"
---

# Persistence 层规范

实现 Domain 层 Repository 接口，处理数据库持久化（GORM）。

## 文件命名

| 文件类型   | 命名规范                             |
| ---------- | ------------------------------------ |
| 包文档     | `doc.go`（**必需**）                 |
| Model      | `{模块}_model.go`                    |
| 写仓储     | `{模块}_command_repository.go`       |
| 读仓储     | `{模块}_query_repository.go`         |
| 缓存装饰器 | `{模块}_cached_{type}_repository.go` |
| 仓储聚合   | `{模块}_repositories.go`（可选）     |

## Model 规范

- 必须定义 `TableName()` 方法
- 必须提供 `toModel()` 和 `toEntity()` 映射函数
- GORM Tag 只在 Model 中使用，Domain 实体禁止 GORM 依赖

## Repository 规范

- Command：写操作（Create/Update/Delete）
- Query：读操作（Get/List/Exists）
- `Create()` 必须回写生成的 ID 到实体

## 关联约束

禁止物理外键，使用逻辑关联：

```go
// ❌ 禁止
CategoryID uint `gorm:"constraint:OnDelete:CASCADE"`

// ✅ 允许
CategoryID uint `gorm:"index;not null"`
```

## 缓存装饰器

使用装饰器模式添加缓存层，对调用方透明。

### 结构

```go
type cachedQueryRepository struct {
    delegate domain.QueryRepository
    cache    cache.CacheService
}

func NewCachedQueryRepository(delegate, cache) domain.QueryRepository
```

### 策略

| 操作类型     | 策略        | 说明                           |
| ------------ | ----------- | ------------------------------ |
| 读 (Query)   | Cache-Aside | 先查缓存，未命中查库，同步回写 |
| 写 (Command) | Invalidate  | 先写库，成功后失效缓存         |

### 关联缓存失效

当写操作影响多个缓存时，需级联失效（均使用**同步失效**）：

```go
func (r *cachedRepo) Upsert(ctx, entity) error {
    if err := r.delegate.Upsert(ctx, entity); err != nil {
        return err
    }
    // 同步失效多个关联缓存
    _ = r.effectiveCache.Delete(ctx, entity.Key)
    _ = r.queryCache.DeleteByUser(ctx, entity.UserID)
    _ = r.schemaCache.DeleteAll(ctx, entity.UserID)
    return nil
}
```

**特例**：全局缓存批量失效可用异步（如 `DeleteAll`），避免阻塞请求：

```go
go func() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    _ = r.globalCache.DeleteAll(ctx)
}()
```
