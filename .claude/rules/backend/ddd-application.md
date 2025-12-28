---
paths:
  - "internal/application/**/*.go"
---

# Application 层规范

## 文件命名

| 文件类型 | 命名规范                                              |
| -------- | ----------------------------------------------------- |
| Command  | `commands.go`                                         |
| Query    | `queries.go`                                          |
| Handler  | `cmd_{action}_handler.go` / `qry_{action}_handler.go` |
| DTO      | `dto.go`                                              |
| Mapper   | `mapper.go`                                           |

多实体时加前缀：`cmd_{entity}_{action}_handler.go`

## 目录结构

```
internal/application/{module}/
├── commands.go
├── queries.go
├── cmd_{action}_handler.go
├── qry_{action}_handler.go
├── dto.go
├── mapper.go
└── doc.go
```
