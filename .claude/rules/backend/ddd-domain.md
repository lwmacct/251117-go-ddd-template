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
├── entity_{module}.go       # 领域实体
├── command_repository.go    # 写仓储接口
├── query_repository.go      # 读仓储接口
├── errors.go                # 错误定义
└── doc.go                   # 包文档（必须）
```

**命名约定**:

- 多实体时次要实体加前缀：`{实体名}_command_repository.go`

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
// command_repository.go
type CommandRepository interface {
    Create(ctx context.Context, e *Entity) error
    Update(ctx context.Context, e *Entity) error
    Delete(ctx context.Context, id uint) error
}

// query_repository.go
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
