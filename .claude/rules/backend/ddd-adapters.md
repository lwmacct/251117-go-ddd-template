---
paths:
  - "internal/adapters/**/*.go"
---

# Adapters 层规范

## 核心原则

接口适配层，**仅做请求绑定和响应转换**，不包含业务逻辑。

## 文件命名

| 目录               | 命名规范      |
| ------------------ | ------------- |
| `http/handler/`    | `{模块}.go`   |
| `http/middleware/` | `{功能}.go`   |
| `http/response/`   | `response.go` |
| `http/`            | `router.go`   |

## 禁止事项

- ❌ 在 Handler 中编排业务逻辑
- ❌ 直接调用 Repository
- ❌ 直接依赖 Infrastructure 实现

## Handler 结构

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
