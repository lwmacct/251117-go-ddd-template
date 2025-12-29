---
paths:
  - "internal/application/**/*.go"
---

# Application 层规范

## 依赖方向（核心原则）

**禁止直接引用 Infrastructure 包**，只能依赖 Domain 层定义的接口。

```
Adapters → Application → Domain ← Infrastructure
```

## 文件命名

| 文件类型        | 命名规范          |
| --------------- | ----------------- |
| Command         | `commands.go`     |
| Query           | `queries.go`      |
| Command Handler | `cmd_{action}.go` |
| Query Handler   | `qry_{action}.go` |
| DTO             | `dto.go`          |
| Mapper          | `mapper.go`       |
| Cache           | `cache.go`        |
| Helper          | `helper.go`       |

**命名约定**：

- `cmd_` / `qry_` 前缀即表示 Handler，无需 `_handler` 后缀
- 多实体时加前缀：`cmd_{entity}_{action}.go`
- 多缓存时按实体拆分：`cache_{entity}.go`
- 多 Helper 时按名称拆分：`helper_{name}.go`

## 目录结构

```
internal/application/{module}/
├── commands.go
├── queries.go
├── cmd_{action}.go      # Command Handler
├── qry_{action}.go      # Query Handler
├── dto.go
├── mapper.go
├── cache.go             # 缓存接口（可选）
├── helper.go            # 辅助函数/Builder（可选）
└── doc.go
```
