---
paths:
  - "internal/adapters/http/handler/**/*.go"
---

# HTTP Handler 规范

## 核心原则

Handler 是 HTTP 层的适配器，**仅做请求绑定和响应转换**，不包含业务逻辑。

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
    response.Created(c, response.MsgCreated, result)
}
```

## 禁止事项

- ❌ 在 Handler 中编排业务逻辑
- ❌ 直接调用 Repository
- ❌ 直接依赖 Infrastructure 实现
- ❌ 自定义错误常量（必须使用 Domain 层定义）

## 响应规范

必须使用 `response/` 包：

```go
// ✅ 正确
response.OK(c, response.MsgSuccess, data)

// ❌ 禁止
c.JSON(200, gin.H{...})
```

### 响应消息规范

**禁止魔法字符串**，按优先级使用：

1. **通用操作消息**：使用 `response` 包常量

   ```go
   response.OK(c, response.MsgSuccess, result)      // 通用成功
   response.Created(c, response.MsgCreated, result)  // 创建成功
   response.OK(c, response.MsgUpdated, result)       // 更新成功
   response.OK(c, response.MsgDeleted, nil)          // 删除成功
   ```

2. **特定业务消息**：使用 Domain 层 `constants.go` 常量

   ```go
   response.OK(c, authDomain.MsgTwoFARequired, dto)
   ```

3. **错误消息**：使用 `err.Error()` 或 Domain 错误
   ```go
   response.NotFoundMessage(c, err.Error())
   response.BadRequest(c, orgDomain.ErrInvalidOrgID.Error())
   ```

```go
// ✅ 正确
response.OK(c, response.MsgSuccess, result)
response.OK(c, authDomain.MsgTwoFARequired, dto)
response.NotFoundMessage(c, err.Error())

// ❌ 禁止
response.OK(c, "success", data)
response.OK(c, "Two factor authentication required", dto)
response.BadRequest(c, "invalid organization ID")
```

### 响应函数选择

| 场景         | 函数                           | 参数说明                       |
| ------------ | ------------------------------ | ------------------------------ |
| 200 成功     | `OK(c, msg, data)`             | msg 常用 `MsgSuccess`          |
| 200 列表     | `List(c, msg, data, meta)`     | msg 常用 `MsgSuccess`          |
| 201 创建     | `Created(c, msg, data)`        | msg 常用 `MsgCreated`          |
| 204 无内容   | `NoContent(c)`                 | 无参数                         |
| 400 参数错误 | `BadRequest(c, msg, details?)` | msg 来自 `err.Error()`         |
| 400 验证错误 | `ValidationError(c, details)`  | 自动使用 `MsgValidationFailed` |
| 401 未认证   | `Unauthorized(c, msg?)`        | 空=默认常量，或 `err.Error()`  |
| 403 无权限   | `Forbidden(c, msg?)`           | 空=默认常量，或 `err.Error()`  |
| 404 不存在   | `NotFoundMessage(c, msg)`      | msg 来自 `err.Error()`         |
| 409 冲突     | `Conflict(c, msg?)`            | 空=默认常量，或 `err.Error()`  |
| 429 限流     | `TooManyRequests(c)`           | 固定消息                       |
| 500 错误     | `InternalError(c, details?)`   | 固定常量 + 详情                |
| 503 不可用   | `ServiceUnavailable(c, msg?)`  | 空=默认常量                    |

## Query 参数

内联定义在 Handler 文件中：

```go
type ListQuery struct {
    response.PaginationQueryDTO
    Status string `form:"status"`
}

func (q *ListQuery) ToQuery() app.ListQuery { ... }
```
