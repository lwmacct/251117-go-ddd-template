---
paths:
  - "internal/adapters/**/*.go"
---

# Adapters 层规范

## 核心原则

- 接口适配层，**仅做请求绑定和响应转换**，不包含业务逻辑。

## 目录结构

```
internal/adapters/http/
├── handler/
│   └── {module}.go     # HTTP Handler
├── middleware/
│   └── {feature}.go    # 中间件
├── response/
│   └── response.go     # 响应工具
└── router.go           # 路由定义
```

**命名约定**:

- `{module}` 为领域模块名（如 `user`、`auth`）
- `{feature}` 为功能名（如 `cors`、`logging`）

## 禁止事项

- ❌ 在 Handler 中编排业务逻辑
- ❌ 直接调用 Repository
- ❌ 直接依赖 Infrastructure 实现
- ❌ 自定义错误常量（必须使用 Domain 层定义）

```go
type XxxHandler struct {
    createHandler *app.CreateHandler
    getHandler    *app.GetHandler
}

func (h *XxxHandler) Create(c *gin.Context) {
    // 1. 请求绑定
    var req CreateDTO
    if err := c.ShouldBindJSON(&req); err != nil { ... }

    // 2. 调用 Application Handler
    result, err := h.createHandler.Handle(ctx, Command{...})

    // 3. 响应转换
    response.Created(c, "success", result)
}
```

## 响应规范

必须使用 `response/` 包：

```go
// ✅ 正确
response.OK(c, "success", data)
response.List(c, "success", items, meta)

// ❌ 禁止
c.JSON(200, gin.H{...})
```

## Query 参数

内联定义在 Handler 文件中：

```go
type ListQuery struct {
    response.PaginationDTO
    Status string `form:"status"`
}

func (q *ListQuery) ToQuery() app.ListQuery { ... }
```

## Swagger 注解规范

| 注解              | 规则                        | 示例                                            |
| ----------------- | --------------------------- | ----------------------------------------------- |
| `@Summary`        | 必填                        | `@Summary 创建用户`                             |
| `@Tags`           | 必填，英文格式              | `@Tags Admin - User Management`                 |
| `@Accept`         | 必须为 `json`               | `@Accept json`                                  |
| `@Produce`        | 必须为 `json`               | `@Produce json`                                 |
| `@Security`       | 非公开端点必须有            | `@Security BearerAuth`                          |
| `@Router`         | 必须与 operation 匹配       | `@Router /api/admin/users [post]`               |
| `@Success` DTO    | 必须在 application 层存在   | `@Success 200 {object} user.UserDTO`            |
| `@Param` body DTO | 必须在 application 层存在   | `@Param request body auth.LoginDTO true "登录"` |
| `@Param` query    | 类型必须在 handler 包中定义 | `@Param request query ListQuery true "查询"`    |

> 规范由 `internal/precommit/annotation_test.go` 自动检查

## 路由规范

- 所有 API 路径必须以 `/api/` 开头
- 每个 operation 必须有 `Method` 和 `Path`
- 审计操作必须有 `Category`、`Operation`、`Label`

> 规范由 `internal/precommit/router_test.go` 自动检查
