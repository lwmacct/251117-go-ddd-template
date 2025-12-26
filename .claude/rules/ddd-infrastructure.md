---
paths:
  - "internal/infrastructure/**/*.go"
---

# Infrastructure 层规范

<!--TOC-->

## Table of Contents

- [核心职责](#核心职责) `:24+4`
- [文件命名规范](#文件命名规范) `:28+19`
  - [persistence 目录](#persistence-目录) `:30+9`
  - [其他目录](#其他目录) `:39+8`
- [持久化 Model 规范](#持久化-model-规范) `:47+21`
- [Repository 实现规范](#repository-实现规范) `:68+23`
- [仓储聚合（可选）](#仓储聚合可选) `:91+17`
- [Domain Service 实现](#domain-service-实现) `:108+16`
- [目录结构示例](#目录结构示例) `:124+16`

<!--TOC-->

## 核心职责

实现 Domain 层定义的接口，处理技术细节（数据库、缓存、外部 API）。

## 文件命名规范

### persistence 目录

| 文件类型     | 命名规范                       | 示例                           |
| ------------ | ------------------------------ | ------------------------------ |
| 持久化 Model | `{模块}_model.go`              | `user_model.go`                |
| 写仓储实现   | `{模块}_command_repository.go` | `user_command_repository.go`   |
| 读仓储实现   | `{模块}_query_repository.go`   | `user_query_repository.go`     |
| 仓储聚合     | `{模块}_repositories.go`       | `user_repositories.go`（可选） |

### 其他目录

| 目录        | 文件类型            | 命名规范        |
| ----------- | ------------------- | --------------- |
| `auth/`     | Domain Service 实现 | `service.go`    |
| `config/`   | 配置管理            | `config.go`     |
| `database/` | 数据库初始化/迁移   | `migrations.go` |

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

## Domain Service 实现

```go
// auth/service.go - 实现 domain/auth.Service 接口
type authService struct {
    jwtManager *JWTManager
}

func NewAuthService(jwtManager *JWTManager) auth.Service {
    return &authService{jwtManager: jwtManager}
}

func (s *authService) HashPassword(password string) (string, error) { ... }
func (s *authService) VerifyPassword(hashedPassword, password string) error { ... }
```

## 目录结构示例

```
internal/infrastructure/
├── persistence/
│   ├── user_model.go                 # GORM Model + 映射函数
│   ├── user_command_repository.go    # 写仓储实现
│   ├── user_query_repository.go      # 读仓储实现
│   └── user_repositories.go          # 仓储聚合（可选）
├── auth/
│   └── service.go                    # 认证服务实现
├── config/
│   └── config.go                     # 配置管理
└── database/
    └── migrations.go                 # 数据库迁移
```
