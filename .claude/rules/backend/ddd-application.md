---
paths:
  - "internal/application/**/*.go"
---

# Application 层规范

## 文件命名

| 文件类型        | 命名规范          |
| --------------- | ----------------- |
| Command         | `commands.go`     |
| Query           | `queries.go`      |
| Command Handler | `cmd_{action}.go` |
| Query Handler   | `qry_{action}.go` |
| DTO             | `dto.go`          |
| Mapper          | `mapper.go`       |

**命名约定**：`cmd_` / `qry_` 前缀即表示 Handler，无需 `_handler` 后缀。

多实体时加前缀：`cmd_{entity}_{action}.go`

## 目录结构

```
internal/application/{module}/
├── commands.go
├── queries.go
├── cmd_{action}.go      # Command Handler
├── qry_{action}.go      # Query Handler
├── dto.go
├── mapper.go
└── doc.go
```
