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

## 目录结构

```
internal/application/{module}/
├── commands.go          # Command 定义
├── queries.go           # Query 定义
├── cmd_{action}.go      # Command Handler
├── qry_{action}.go      # Query Handler
├── dto.go               # 数据传输对象
├── mapper.go            # 实体-DTO 映射
├── cache.go             # 缓存接口（可选）
├── helper.go            # 辅助函数/Builder（可选）
└── doc.go               # 包文档
```

**命名约定**:

- `cmd_` / `qry_` 前缀即表示 Handler，无需 `_handler` 后缀
- 多实体时加前缀：`cmd_{entity}_{action}.go`
- 多缓存时按实体拆分：`cache_{entity}.go`
- 多 Helper 时按名称拆分：`helper_{name}.go`

## 命名规范

| 文件          | 类型   | 规则                          | 示例                      |
| ------------- | ------ | ----------------------------- | ------------------------- |
| `commands.go` | struct | 以 `Command` 结尾             | `CreateUserCommand`       |
| `queries.go`  | struct | 以 `Query` 结尾               | `GetUserByIDQuery`        |
| `cmd_*.go`    | struct | 以 `Handler` 结尾             | `CreateUserHandler`       |
| `qry_*.go`    | struct | 以 `Handler` 结尾             | `GetUserHandler`          |
| `dto.go`      | struct | 以 `DTO` 结尾                 | `UserDTO`, `UserListDTO`  |
| `mapper.go`   | func   | `To` 开头 + `DTO`/`DTOs` 结尾 | `ToUserDTO`, `ToUserDTOs` |

> 规范由 `internal/precommit/application_test.go` 自动检查
