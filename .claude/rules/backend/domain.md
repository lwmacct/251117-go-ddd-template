---
paths:
  - "internal/domain/**/*.go"
---

# Domain 层规范

## 核心原则

Domain 层是 DDD 架构的核心，**不依赖任何外层**。

## 文件命名规范

| 文件类型   | 命名规范                | 示例               |
| ---------- | ----------------------- | ------------------ |
| 实体模型   | `entity_{模块}.go`      | `entity_user.go`   |
| 写仓储接口 | `command_repository.go` | 每个模块固定命名   |
| 读仓储接口 | `query_repository.go`   | 每个模块固定命名   |
| 错误定义   | `errors.go`             | 领域错误           |
| 值对象     | `value_objects.go`      | 复杂领域需要时使用 |
| 包文档     | `doc.go`                | **必须包含**       |

## 多实体命名

当模块包含多个实体时，用实体名区分：

| 文件类型   | 主实体（固定命名）      | 次要实体（实体名前缀）           |
| ---------- | ----------------------- | -------------------------------- |
| 实体模型   | `entity_{模块}.go`      | `entity_{实体名}.go`             |
| 写仓储接口 | `command_repository.go` | `{实体名}_command_repository.go` |
| 读仓储接口 | `query_repository.go`   | `{实体名}_query_repository.go`   |

**示例**（setting 模块）：

| 实体        | 文件                                 | 接口                           |
| ----------- | ------------------------------------ | ------------------------------ |
| Setting     | `entity_setting.go`                  | `CommandRepository`            |
|             | `command_repository.go`              | `QueryRepository`              |
|             | `query_repository.go`                |                                |
| UserSetting | `entity_user_setting.go`             | `UserSettingCommandRepository` |
|             | `user_setting_command_repository.go` | `UserSettingQueryRepository`   |
|             | `user_setting_query_repository.go`   |                                |

## 禁止事项

- ❌ **禁止任何 GORM Tag 或 `gorm` 依赖**
- ❌ 禁止 import 外层代码（Infrastructure/Adapters）
- ❌ 禁止包含数据库/Redis/HTTP 等技术实现

## 实体规范（富领域模型）

```go
// entity_xxx.go - 仅含业务字段和行为方法
type Xxx struct {
    ID   uint
    Name string
    // 无 GORM Tag！
}

// 业务行为通过方法体现
func (x *Xxx) IsValid() bool { ... }
func (x *Xxx) Activate() { ... }
```

## Repository 接口规范（CQRS 分离）

```go
// command_repository.go - 写操作
type CommandRepository interface {
    Create(ctx context.Context, entity *Xxx) error
    Update(ctx context.Context, entity *Xxx) error
    Delete(ctx context.Context, id uint) error
}

// query_repository.go - 读操作
type QueryRepository interface {
    GetByID(ctx context.Context, id uint) (*Xxx, error)
    List(ctx context.Context, offset, limit int) ([]*Xxx, error)
    ExistsByName(ctx context.Context, name string) (bool, error)
}
```

## doc.go 包文档规范

```go
// Package xxx 定义 xxx 领域模型和仓储接口。
//
// 本包是 xxx 管理的领域层核心，定义了：
//   - [Xxx]: xxx 实体（富领域模型）
//   - [CommandRepository]: 写仓储接口
//   - [QueryRepository]: 读仓储接口
//   - xxx 领域错误（见 errors.go）
//
// 依赖倒置：
// 本包仅定义接口，实现位于 infrastructure/persistence 包。
package xxx
```

## 目录结构示例

### 单实体模块

```
internal/domain/user/
├── entity_user.go           # User 实体/领域行为
├── command_repository.go    # 写仓储接口
├── query_repository.go      # 读仓储接口
├── errors.go                # 领域错误
└── doc.go                   # 包文档（必须）
```

### 多实体模块

```
internal/domain/setting/
├── entity_setting.go                    # 主实体：Setting
├── entity_user_setting.go               # 次要实体：UserSetting
├── command_repository.go                # Setting 写仓储接口
├── query_repository.go                  # Setting 读仓储接口
├── user_setting_command_repository.go   # UserSetting 写仓储接口
├── user_setting_query_repository.go     # UserSetting 读仓储接口
├── errors.go
└── doc.go
```
