---
paths:
  - "internal/application/**/*.go"
---

# Application 层规范

## 文件与结构体

| 文件                    | 结构体     | 示例                                      |
| ----------------------- | ---------- | ----------------------------------------- |
| `commands.go`           | `*Command` | `CreateCommand`, `UpdateCommand`          |
| `queries.go`            | `*Query`   | `GetQuery`, `ListQuery`                   |
| `cmd_{操作}_handler.go` | `*Handler` | `cmd_create_handler.go` → `CreateHandler` |
| `qry_{操作}_handler.go` | `*Handler` | `qry_get_handler.go` → `GetHandler`       |
| `dto.go`                | `*DTO`     | `CreateDTO`, `UserDTO`                    |
| `mapper.go`             | -          | Entity → DTO 映射                         |

## 多实体命名

当包含多个实体时，用实体名前缀区分：

| 文件                        | 结构体             | 说明              |
| --------------------------- | ------------------ | ----------------- |
| `cmd_user_set_handler.go`   | `UserSetCommand`   | User 实体的 Set   |
| `qry_user_get_handler.go`   | `UserGetQuery`     | User 实体的 Get   |
| `cmd_user_reset_handler.go` | `UserResetCommand` | User 实体的 Reset |

## 目录结构

```
internal/application/{module}/
├── commands.go                       # 所有 Command 定义
├── queries.go                        # 所有 Query 定义
├── cmd_{action}_handler.go           # 主实体 Handler
├── cmd_{entity}_{action}_handler.go  # 次要实体 Handler
├── qry_{action}_handler.go
├── qry_{entity}_{action}_handler.go
├── dto.go, mapper.go, doc.go
```
