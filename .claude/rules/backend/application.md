---
paths:
  - "internal/application/**/*.go"
---

# Application 层规范

## 命名原则

包名提供上下文，类型名不重复：`user.CreateCommand` ✅ / `user.CreateUserCommand` ❌

## 文件与结构体

| 文件                    | 结构体     | 示例                                      |
| ----------------------- | ---------- | ----------------------------------------- |
| `commands.go`           | `*Command` | `CreateCommand`, `UpdateCommand`          |
| `queries.go`            | `*Query`   | `GetQuery`, `ListQuery`                   |
| `cmd_{操作}_handler.go` | `*Handler` | `cmd_create_handler.go` → `CreateHandler` |
| `qry_{操作}_handler.go` | `*Handler` | `qry_get_handler.go` → `GetHandler`       |
| `dto.go`                | `*DTO`     | `CreateDTO`, `UserDTO`                    |
| `mapper.go`             | -          | Entity → DTO 映射                         |

## 目录结构

```
internal/application/{module}/
├── commands.go              # 所有 Command 定义
├── queries.go               # 所有 Query 定义
├── cmd_{action}_handler.go  # Handler 保持分离
├── qry_{action}_handler.go
├── dto.go, mapper.go, doc.go
```

## Command/Query 定义规范

在 `commands.go` 和 `queries.go` 中，按操作类型排列：

```go
package xxx
type CreateCommand struct { ... }
type UpdateCommand struct { ... }
type DeleteCommand struct { ... }
```

```go
package xxx
type GetQuery struct { ... }
type ListQuery struct { ... }
```
