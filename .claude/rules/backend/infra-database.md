---
paths:
  - "internal/infrastructure/database/**/*.go"
---

# Database Infrastructure 规范

## 核心职责

管理数据库连接、迁移和种子数据。

## 文件命名

| 文件类型   | 命名规范             |
| ---------- | -------------------- |
| 包文档     | `doc.go`（**必需**） |
| 连接管理   | `connection.go`      |
| 迁移执行器 | `migrator.go`        |
| 种子执行器 | `seeder.go`          |

## 迁移原则

使用 GORM AutoMigrate：

- Model 即 Schema（单一数据源）
- 记录迁移历史到 `migrations` 表

## 种子数据原则

- 实现 `Seed()` 和 `Name()` 方法
- 通过 registry 统一注册
- 幂等执行
