---
paths:
  - "internal/application/**/*.go"
---

# Application 层规范

<!--TOC-->

## Table of Contents

- [文件命名规范](#文件命名规范) `:19+12`
- [结构体命名强制规范](#结构体命名强制规范) `:31+10`
- [DTO 规范](#dto-规范) `:41+28`
- [目录结构示例](#目录结构示例) `:69+18`

<!--TOC-->

## 文件命名规范

| 文件类型        | 命名规范                | 示例                         |
| --------------- | ----------------------- | ---------------------------- |
| Command 定义    | `cmd_{操作}.go`         | `cmd_create_user.go`         |
| Command Handler | `cmd_{操作}_handler.go` | `cmd_create_user_handler.go` |
| Query 定义      | `qry_{操作}.go`         | `qry_get_user.go`            |
| Query Handler   | `qry_{操作}_handler.go` | `qry_get_user_handler.go`    |
| DTO 定义        | `dto.go`                | 每个模块固定命名             |
| Mapper          | `mapper.go`             | Entity → DTO 映射函数        |
| 包文档          | `doc.go`                | 每个模块必须包含             |

## 结构体命名强制规范

pre-commit 检查会验证以下规则：

| 文件模式   | 结构体后缀要求 | 示例                                     |
| ---------- | -------------- | ---------------------------------------- |
| `cmd_*.go` | 仅 `*Command`  | `CreateUserCommand`, `UpdateRoleCommand` |
| `qry_*.go` | 仅 `*Query`    | `GetUserQuery`, `ListUsersQuery`         |
| `dto.go`   | 仅 `*DTO`      | `UserDTO`, `CreateUserResultDTO`         |

## DTO 规范

DTO 文件 (`dto.go`) 中应包含以下类型：

1. **请求 DTO** - HTTP 请求绑定

   ```go
   type CreateXxxDTO struct {
       Name string `json:"name" binding:"required"`
   }
   ```

2. **结果 DTO** - Handler 返回值

   ```go
   type CreateXxxResultDTO struct {
       ID uint `json:"id"`
   }
   ```

3. **响应 DTO** - HTTP 响应体
   ```go
   type XxxResponseDTO struct {
       ID   uint   `json:"id"`
       Name string `json:"name"`
   }
   ```

## 目录结构示例

```
internal/application/xxx/
├── cmd_create_xxx.go           # CreateXxxCommand
├── cmd_create_xxx_handler.go   # CreateXxxHandler
├── cmd_update_xxx.go           # UpdateXxxCommand
├── cmd_update_xxx_handler.go   # UpdateXxxHandler
├── cmd_delete_xxx.go           # DeleteXxxCommand
├── cmd_delete_xxx_handler.go   # DeleteXxxHandler
├── qry_get_xxx.go              # GetXxxQuery
├── qry_get_xxx_handler.go      # GetXxxHandler
├── qry_list_xxx.go             # ListXxxQuery
├── qry_list_xxx_handler.go     # ListXxxHandler
├── dto.go                      # 所有 DTO 定义
├── mapper.go                   # Entity → DTO 映射
└── doc.go                      # 包文档
```
