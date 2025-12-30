---
paths:
  - "internal/domain/**/*.go"
---

# Domain 层规范

## 核心原则

DDD 架构核心，**不依赖任何外层**。

## 目录结构

```
internal/domain/{module}/
├── entity_{module}.go       # 主实体
├── entity_{xxx}.go          # 次要实体（可选）
├── repository.go            # Repository 接口（Command + Query）
├── cmd_{entity}.go          # 写仓储接口（多实体时）
├── qry_{entity}.go          # 读仓储接口（多实体时）
├── errors.go                # 错误定义
├── value_objects.go         # 值对象（可选）
└── doc.go                   # 包文档（必须）
```

**命名约定**:

- `cmd_` / `qry_` 前缀即表示仓储接口，无需 `_repository` 后缀
- 多实体时按实体拆分：`cmd_{entity}.go`、`qry_{entity}.go`

## 命名规范

| 文件               | 类型      | 规则                             | 示例                               |
| ------------------ | --------- | -------------------------------- | ---------------------------------- |
| `repository.go`    | interface | `Command`/`Query` + `Repository` | `CommandRepository`                |
| `cmd_*.go`         | interface | `{Entity}CommandRepository`      | `SettingCategoryCommandRepository` |
| `qry_*.go`         | interface | `{Entity}QueryRepository`        | `UserSettingQueryRepository`       |
| `entity_*.go`      | struct    | 实体名（首字母大写）             | `User`, `Setting`                  |
| `value_objects.go` | struct    | 值对象名                         | `Scope`, `InputType`               |

> 规范由 `internal/precommit/domain_test.go` 自动检查

## 禁止事项

- ❌ GORM Tag 或 `gorm` 依赖
- ❌ import 外层代码
- ❌ 数据库/Redis/HTTP 等技术实现

## 实体规范

```go
type Entity struct {
    ID   uint
    Name string
    // 无 GORM Tag
}

// 业务行为通过方法体现
func (e *Entity) IsValid() bool { ... }
func (e *Entity) Activate() { ... }
```

## Repository 接口（CQRS）

```go
// repository.go - 合并定义
type CommandRepository interface {
    Create(ctx context.Context, e *Entity) error
    Update(ctx context.Context, e *Entity) error
    Delete(ctx context.Context, id uint) error
}

type QueryRepository interface {
    GetByID(ctx context.Context, id uint) (*Entity, error)
    List(ctx context.Context, offset, limit int) ([]*Entity, error)
}
```

## doc.go 规范

```go
// Package xxx 定义 xxx 领域模型和仓储接口。
//
// 本包定义了：
//   - [Entity]: 实体（富领域模型）
//   - [CommandRepository]: 写仓储接口
//   - [QueryRepository]: 读仓储接口
//
// 依赖倒置：本包仅定义接口，实现位于 infrastructure/persistence。
package xxx
```
