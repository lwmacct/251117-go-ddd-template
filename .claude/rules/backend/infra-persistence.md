---
paths:
  - "internal/infrastructure/persistence/**/*.go"
---

# Persistence 层规范

实现 Domain 层 Repository 接口，处理数据库持久化（GORM）。

## 目录结构

```
internal/infrastructure/persistence/
├── doc.go                                # 包文档（必需）
├── generic_repository.go                 # 通用仓储基类
├── {module}_model.go                     # 数据模型
├── {module}_command_repository.go        # 写仓储实现
├── {module}_query_repository.go          # 读仓储实现
├── {module}_cached_{type}_repository.go  # 缓存装饰器（可选）
└── {module}_repositories.go              # 仓储聚合（可选）
```

**命名约定**:

- `{module}` 为领域模块名（如 `user`、`role`、`setting`）
- `{type}` 为仓储类型：`query` 或 `command`
- 复合模块用下划线连接（如 `user_setting`、`setting_category`）

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
