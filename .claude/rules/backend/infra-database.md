---
paths:
  - "internal/infrastructure/database/**/*.go"
---

# Database Infrastructure 规范

## 核心职责

管理数据库连接、迁移和种子数据。

## 文件结构

| 文件                   | 职责                   |
| ---------------------- | ---------------------- |
| `doc.go`               | 包文档（**必需**）     |
| `connection.go`        | 数据库连接配置         |
| `migrator.go`          | 迁移执行器             |
| `migration_manager.go` | 迁移版本管理（历史表） |
| `seeder.go`            | 种子数据执行器         |
| `slog_logger.go`       | GORM 日志适配          |

## Seeds 目录

| 文件                | 职责          |
| ------------------- | ------------- |
| `registry.go`       | 种子注册表    |
| `user_seeder.go`    | 用户种子数据  |
| `rbac_seeder.go`    | RBAC 种子数据 |
| `setting_seeder.go` | 设置种子数据  |

## 迁移原则

本项目采用 **GORM AutoMigrate** 自动迁移方案：

- 模型定义在 `infrastructure/persistence/*_model.go`
- `Migrator` 调用 `db.AutoMigrate()` 执行迁移
- `MigrationManager` 在 `migrations` 表中记录迁移历史

**优势**：单一数据源（Model 即 Schema），无需维护独立迁移文件

**限制**：不支持复杂的自定义迁移脚本（如数据迁移）

## 种子数据原则

- 每个 Seeder 实现 `Seed()` 和 `Name()` 方法
- 通过 `registry.go` 统一注册
- 支持幂等执行（重复运行不报错）
