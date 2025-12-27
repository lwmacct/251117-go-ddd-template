---
paths:
  - "internal/infrastructure/persistence/**/*.go"
---

# Persistence 层规范

## 核心职责

实现 Domain 层定义的 Repository 接口，处理**数据库持久化**（GORM）。

**注意**：非数据库存储（如内存缓存、聚合查询）应放在独立的 infrastructure 模块中：

- 内存存储 → `infrastructure/captcha/`
- 聚合查询 → `infrastructure/stats/`

## 文件命名规范

| 文件类型     | 命名规范                       | 示例                           |
| ------------ | ------------------------------ | ------------------------------ |
| 持久化 Model | `{模块}_model.go`              | `user_model.go`                |
| 写仓储实现   | `{模块}_command_repository.go` | `user_command_repository.go`   |
| 读仓储实现   | `{模块}_query_repository.go`   | `user_query_repository.go`     |
| 仓储聚合     | `{模块}_repositories.go`       | `user_repositories.go`（可选） |

## 多实体模块

当 Domain 模块包含多个实体时，Infrastructure 层对应扩展：

| 文件类型   | 主实体                         | 次要实体                         |
| ---------- | ------------------------------ | -------------------------------- |
| Model      | `{模块}_model.go`              | `{实体名}_model.go`              |
| 写仓储实现 | `{模块}_command_repository.go` | `{实体名}_command_repository.go` |
| 读仓储实现 | `{模块}_query_repository.go`   | `{实体名}_query_repository.go`   |
| 仓储聚合   | `{模块}_repositories.go`       | `{实体名}_repositories.go`       |

**命名一致性原则**：Infrastructure Model/Repository 文件名应与 Domain 实体名保持一致。

## 持久化 Model 规范

```go
// {模块}_model.go
type XxxModel struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"size:100;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// 表名
func (XxxModel) TableName() string { return "xxx" }

// Entity → Model 映射
func newXxxModelFromEntity(entity *xxx.Xxx) *XxxModel { ... }

// Model → Entity 映射
func (m *XxxModel) toEntity() *xxx.Xxx { ... }
```

## Repository 实现规范

```go
// {模块}_command_repository.go
type xxxCommandRepository struct { db *gorm.DB }

func NewXxxCommandRepository(db *gorm.DB) xxx.CommandRepository {
    return &xxxCommandRepository{db: db}
}

func (r *xxxCommandRepository) Create(ctx context.Context, entity *xxx.Xxx) error {
    model := newXxxModelFromEntity(entity)
    if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
        return err
    }
    // 回写生成的 ID
    if saved := model.toEntity(); saved != nil {
        *entity = *saved
    }
    return nil
}
```

## 仓储聚合（可选）

```go
// {模块}_repositories.go - 便于依赖注入
type XxxRepositories struct {
    Command xxx.CommandRepository
    Query   xxx.QueryRepository
}

func NewXxxRepositories(db *gorm.DB) XxxRepositories {
    return XxxRepositories{
        Command: NewXxxCommandRepository(db),
        Query:   NewXxxQueryRepository(db),
    }
}
```

## 目录结构示例

### 单实体模块

```
internal/infrastructure/persistence/
├── user_model.go                 # GORM Model + 映射函数
├── user_command_repository.go    # 写仓储实现
├── user_query_repository.go      # 读仓储实现
└── user_repositories.go          # 仓储聚合（可选）
```

### 多实体模块

```
internal/infrastructure/persistence/
├── setting_model.go                     # Setting Model
├── setting_command_repository.go        # Setting 写仓储实现
├── setting_query_repository.go          # Setting 读仓储实现
├── setting_repositories.go              # Setting 仓储聚合
├── user_setting_model.go                # UserSetting Model
├── user_setting_command_repository.go   # UserSetting 写仓储实现
├── user_setting_query_repository.go     # UserSetting 读仓储实现
└── user_setting_repositories.go         # UserSetting 仓储聚合
```
