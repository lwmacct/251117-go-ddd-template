---
paths:
  - "internal/adapters/**/*.go"
---

# Adapters 层规范

<!--TOC-->

## Table of Contents

- [核心原则](#核心原则) `:23+4`
- [文件命名规范](#文件命名规范) `:27+9`
- [禁止事项](#禁止事项) `:36+6`
- [HTTP Handler 规范](#http-handler-规范) `:42+23`
- [Handler 方法规范](#handler-方法规范) `:65+41`
- [统一响应格式](#统一响应格式) `:106+14`
- [Query 参数结构体](#query-参数结构体) `:120+28`
- [目录结构示例](#目录结构示例) `:148+16`

<!--TOC-->

## 核心原则

Adapters 层是接口适配层，**仅做请求绑定和响应转换**，不包含业务逻辑。

## 文件命名规范

| 目录               | 文件类型     | 命名规范            | 示例      |
| ------------------ | ------------ | ------------------- | --------- |
| `http/handler/`    | HTTP Handler | `{模块}.go`（单数） | `user.go` |
| `http/middleware/` | 中间件       | `{功能}.go`         | `auth.go` |
| `http/`            | 路由定义     | `router.go`         | 固定命名  |
| `http/response/`   | 响应工具     | `response.go`       | 固定命名  |

## 禁止事项

- ❌ **禁止在 Handler 中编排业务逻辑**
- ❌ **禁止直接调用 Repository**
- ❌ 禁止直接依赖 Infrastructure 实现

## HTTP Handler 规范

```go
// handler/xxx.go
type XxxHandler struct {
    createXxxHandler *xxx.CreateXxxHandler  // 依赖 Application Handler
    getXxxHandler    *xxx.GetXxxHandler
    listXxxHandler   *xxx.ListXxxHandler
}

func NewXxxHandler(
    createHandler *xxx.CreateXxxHandler,
    getHandler *xxx.GetXxxHandler,
    listHandler *xxx.ListXxxHandler,
) *XxxHandler {
    return &XxxHandler{
        createXxxHandler: createHandler,
        getXxxHandler:    getHandler,
        listXxxHandler:   listHandler,
    }
}
```

## Handler 方法规范

```go
// Create 处理创建请求
func (h *XxxHandler) Create(c *gin.Context) {
    // 1. 请求绑定
    var req xxx.CreateXxxDTO
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request", err)
        return
    }

    // 2. 调用 Application Handler（业务委托）
    result, err := h.createXxxHandler.Handle(c.Request.Context(), xxx.CreateXxxCommand{
        Name: req.Name,
    })
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to create", err)
        return
    }

    // 3. 响应转换
    response.Success(c, http.StatusCreated, "Created successfully", result)
}

// Get 处理获取请求
func (h *XxxHandler) Get(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

    result, err := h.getXxxHandler.Handle(c.Request.Context(), xxx.GetXxxQuery{
        ID: uint(id),
    })
    if err != nil {
        response.Error(c, http.StatusNotFound, "Not found", err)
        return
    }

    response.Success(c, http.StatusOK, "Success", xxx.ToXxxResponseDTO(result))
}
```

## 统一响应格式

**所有 API 必须使用 `response/` 包，禁止 `c.JSON()` 或 `gin.H{}`**

```go
// 正确
response.OK(c, "success", userDTO)
response.Created(c, "created", result)
response.List(c, "success", items, response.NewPaginationMeta(total, page, limit))

// 禁止
c.JSON(200, gin.H{"user": dto})
```

## Query 参数结构体

Query 参数结构体应**内联定义**在对应的 Handler 文件中，遵循就近原则：

```go
// handler/auditlog.go

// ListAuditLogsQuery 审计日志列表查询参数
type ListAuditLogsQuery struct {
    response.PaginationQueryDTO  // 嵌入通用分页参数

    // UserID 按用户 ID 过滤
    UserID *uint `form:"user_id" json:"user_id" binding:"omitempty,gt=0"`
    // Action 操作类型过滤
    Action string `form:"action" json:"action" binding:"omitempty,oneof=create update delete" enums:"create,update,delete"`
}

// ToQuery 转换为 Application 层 Query 对象
func (q *ListAuditLogsQuery) ToQuery() auditlog.ListLogsQuery {
    return auditlog.ListLogsQuery{
        Page:   q.GetPage(),
        Limit:  q.GetLimit(),
        UserID: q.UserID,
        Action: q.Action,
    }
}
```

## 目录结构示例

```
internal/adapters/http/
├── handler/
│   ├── user.go           # User Handler
│   ├── role.go           # Role Handler
│   └── menu.go           # Menu Handler
├── middleware/
│   ├── auth.go           # 认证中间件
│   └── cors.go           # CORS 中间件
├── response/
│   └── response.go       # 统一响应工具
├── router.go             # 路由定义
└── docs/                 # Swagger 文档（自动生成）
```
