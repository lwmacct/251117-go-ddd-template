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
| 包文档       | `doc.go`                       | **必需**                       |
| 泛型基类     | `generic_repository.go`        | 可选，减少样板代码             |
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

## 设计原则

### Model 规范

- 必须定义 `TableName()` 方法
- 必须提供 `newXxxModelFromEntity()` 和 `toEntity()` 映射函数
- GORM Tag 只在 Model 中使用，Domain 实体禁止 GORM 依赖

### Repository 规范

- Command Repository 实现写操作（Create/Update/Delete）
- Query Repository 实现读操作（Get/List/Exists）
- `Create()` 方法必须回写生成的 ID 到实体

### 泛型基类（推荐）

使用 `GenericCommandRepository[E, M]` 减少样板代码：

- 统一 CRUD 实现
- 类型安全，编译时检查
- 仅需实现 `toModel` 和 `toEntity` 映射函数

## 目录结构示例

### 单实体模块

```
user_model.go                 # GORM Model + 映射函数
user_command_repository.go    # 写仓储实现
user_query_repository.go      # 读仓储实现
user_repositories.go          # 仓储聚合（可选）
```

### 多实体模块

```
setting_model.go                     # Setting Model
setting_command_repository.go        # Setting 写仓储
user_setting_model.go                # UserSetting Model
user_setting_command_repository.go   # UserSetting 写仓储
```
