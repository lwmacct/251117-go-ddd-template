---
paths:
  - "internal/infrastructure/persistence/**/*.go"
---

# Persistence 层规范

## 核心职责

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

### Query 装饰策略

| 场景     | 策略                         |
| -------- | ---------------------------- |
| 单条查询 | Cache-Aside + 版本化回写     |
| 批量查询 | 缓存命中 → 未命中查库 → 回写 |
| 列表查询 | 预热完成从缓存过滤，否则查库 |

### Command 装饰策略

写操作后**失效缓存**（非更新）：

```go
func (r *cachedRepo) Update(ctx, entity) error {
    if err := r.delegate.Update(ctx, entity); err != nil {
        return err
    }
    _ = r.cache.Delete(ctx, entity.Key)  // 同步失效
    return nil
}
```

### 跨层失效

当存在多层缓存时，写操作需级联失效：

| 触发操作 | 本层缓存 | 派生层缓存   |
| -------- | -------- | ------------ |
| Update   | 同步删除 | 异步批量删除 |
| Delete   | 同步删除 | 异步批量删除 |

### 异步失效

使用独立 context：

```go
go func() { //nolint:contextcheck
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    _ = r.derivedCache.DeleteByKey(ctx, key)
}()
```
