---
paths:
  - "internal/infrastructure/database/**/*.go"
---

# Database Infrastructure 规范

## 核心职责

管理数据库连接、迁移和种子数据。

## 文件结构

| 文件                   | 职责           |
| ---------------------- | -------------- |
| `connection.go`        | 数据库连接配置 |
| `migrator.go`          | 迁移执行器     |
| `migration_manager.go` | 迁移版本管理   |
| `seeder.go`            | 种子数据执行器 |
| `slog_logger.go`       | GORM 日志适配  |

## Seeds 目录

| 文件                | 职责          |
| ------------------- | ------------- |
| `registry.go`       | 种子注册表    |
| `user_seeder.go`    | 用户种子数据  |
| `rbac_seeder.go`    | RBAC 种子数据 |
| `setting_seeder.go` | 设置种子数据  |

## 迁移规范

迁移文件位于 `migrations/` 子目录，使用时间戳命名：

```
migrations/
├── 20241201_create_users.go
├── 20241202_create_roles.go
└── 20241203_create_permissions.go
```

## 种子数据规范

```go
// seeds/xxx_seeder.go
type XxxSeeder struct {
    db *gorm.DB
}

func NewXxxSeeder(db *gorm.DB) *XxxSeeder {
    return &XxxSeeder{db: db}
}

func (s *XxxSeeder) Seed() error { ... }
func (s *XxxSeeder) Name() string { return "xxx" }
```
