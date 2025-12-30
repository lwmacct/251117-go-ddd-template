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
├── entity.go                # 主实体
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

## 允许事项

- ✅ `json` tags（标准库，不引入外部依赖）

> 若实体需被缓存，必须添加 `json` tags 以确保序列化行为可控。

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

## 细粒度接口模式（可选）

当 QueryRepository 方法较多且使用场景明确分化时，可拆分为细粒度接口：

```go
// 细粒度接口（遵循 ISP）
type BaseQueryRepository interface { ... }
type AuthQueryRepository interface { ... }
type ValidationQueryRepository interface { ... }

// 聚合接口
type QueryRepository interface {
    BaseQueryRepository
    AuthQueryRepository
    ValidationQueryRepository
}
```

**适用场景**：

- 接口方法超过 10 个
- 调用方有明确的功能子集需求（如认证、验证、列表）
- 需要精细化的 Mock 测试

**设计要点**：

- 实现类只需实现聚合接口
- 调用方可按需注入细粒度接口或聚合接口

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

## 测试规范

实体方法和值对象应有单元测试，验证业务规则。
