---
paths:
  - "internal/domain/**/errors.go"
---

# Domain errors.go 规范

## 核心原则

- Domain 层定义业务错误，不包含 HTTP 状态码
- Handler 层通过 `errors.Is()` 检查并映射到响应

## 错误定义格式

- 使用标准库 `errors.New()` 定义
- 错误常量命名以 `Err` 开头，驼峰命名
- 每个错误必须带中文注释
- 错误消息使用英文（面向前端展示）

```go
package xxx

import "errors"

// Xxx 相关错误
var (
	// ErrXxxNotFound xxx 不存在
	ErrXxxNotFound = errors.New("xxx not found")

	// ErrXxxAlreadyExists xxx 已存在
	ErrXxxAlreadyExists = errors.New("xxx already exists")
)
```

## 分组规范

不同分类的错误用独立的 `var()` 包裹，按功能分组：

```go
// 主实体相关错误
var (
	// ErrEntityNotFound 实体不存在
	ErrEntityNotFound = errors.New("entity not found")

	// ErrEntityAlreadyExists 实体已存在
	ErrEntityAlreadyExists = errors.New("entity already exists")
)

// 业务约束相关错误
var (
	// ErrCannotDeleteEntity 不能删除实体
	ErrCannotDeleteEntity = errors.New("cannot delete entity")

	// ErrInvalidEntityState 无效的实体状态
	ErrInvalidEntityState = errors.New("invalid entity state")
)
```

## 常见错误类型分组建议

| 分组        | 典型错误示例                                           |
| ----------- | ------------------------------------------------------ |
| 主实体 CRUD | `ErrXxxNotFound`, `ErrXxxAlreadyExists`                |
| 关联实体    | `ErrYyyNotFound`, `ErrYyyAlreadyExists`                |
| 业务约束    | `ErrCannotDeleteXxx`, `ErrInvalidXxx`                  |
| 系统保护    | `ErrCannotDeleteSystemXxx`, `ErrCannotModifySystemXxx` |

## HTTP 状态码映射

Handler 层通过 `errors.Is()` 检查 Domain 错误，映射到对应 HTTP 状态码：

| HTTP 状态 | 响应函数              | 典型错误后缀                              |
| --------- | --------------------- | ----------------------------------------- |
| 400       | `response.BadRequest` | `*Invalid*`, `*Mismatch*`, `*Weak*`       |
| 403       | `response.Forbidden`  | `*Cannot*`, `*Suspended*`, `NoPermission` |
| 404       | `response.NotFound`   | `*NotFound`                               |
| 409       | `response.Conflict`   | `*AlreadyExists`, `*NameAlreadyExists`    |

## Handler 使用示例

```go
import (
	d_org "github.com/.../internal/domain/org"
	"github.com/.../internal/adapters/http/response"
)

func (h *Handler) Delete(c *gin.Context) {
	// ...
	err := h.deleteHandler.Handle(ctx, cmd)
	if errors.Is(err, d_org.ErrOrgHasMembers) {
		response.BadRequest(c, err.Error())
		return
	}
	if errors.Is(err, d_org.ErrOrgNotFound) {
		response.NotFound(c, "Organization")
		return
	}
	// ...
}
```
