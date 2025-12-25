# Handler 开发规范

<!--TOC-->

## Table of Contents

- [响应封装原则](#响应封装原则) `:17+14`
- [Query 参数结构体](#query-参数结构体) `:31+54`
  - [定义位置](#定义位置) `:33+28`
  - [Swagger 注解](#swagger-注解) `:61+24`
- [Swagger 注解规范](#swagger-注解规范) `:85+35`
  - [注解顺序](#注解顺序) `:91+19`
  - [关键规则](#关键规则) `:110+10`

<!--TOC-->

## 响应封装原则

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

### 定义位置

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

### Swagger 注解

使用结构体引用方式声明 Query 参数：

```go
// @Param        params query handler.ListAuditLogsQuery false "查询参数"
```

**重要约束**：

| 约束                  | 说明                                             |
| --------------------- | ------------------------------------------------ |
| 参数名必须用 `params` | 避免 `query query` 重复词被 golangci-lint 误修改 |
| 禁止 `example` 标签   | OpenAPI 2.0 不支持 query 参数的 example 属性     |
| 参数按字母序生成      | 前端 API 参数顺序按字段名字母排序，非定义顺序    |

```go
// 正确 - 使用 params 作为参数名
// @Param        params query handler.ListUsersQuery false "查询参数"

// 错误 - query query 会被 linter 修改为 query
// @Param        query query handler.ListUsersQuery false "查询参数"
```

## Swagger 注解规范

> **注解直接影响前端代码生成！**
>
> 前端通过 `pnpm api:generate` 从 swagger.json 生成 TypeScript 客户端到 `src/generated/`

### 注解顺序

```go
// @Summary      简短描述
// @Description  详细说明
// @Tags         管理员 - 用户管理 (Admin - User Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        params query handler.ListUsersQuery false "查询参数"
// @Success      200 {object} response.PagedResponse[user.UserDTO] "成功"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Router       /api/admin/users [get]

// 详情接口示例（带 path 参数）
// @Param        id path int true "用户ID"
// @Router       /api/admin/users/{id} [get]
```

### 关键规则

| 规则       | 说明                                                             |
| ---------- | ---------------------------------------------------------------- |
| Tags 格式  | `中文名 (English Name)`，影响生成的 API 类名                     |
| 响应类型   | 必须用 `response.DataResponse[T]` 或 `response.PagedResponse[T]` |
| 路径参数   | 用 `{id}` 而非 `:id`                                             |
| 认证端点   | 必须加 `@Security BearerAuth`                                    |
| Query 参数 | 使用 `@Param params query handler.StructName` 格式               |
| 枚举值     | 使用 `enums:"val1,val2"` 标签，不要用 `example`                  |
