---
paths: internal/adapters/http/handler/**
---

# Swagger 注解规范

> **注解直接影响前端代码生成！**
>
> 前端通过 `pnpm api:generate` 从 swagger.json 生成 TypeScript 客户端到 `src/generated/`

<!--TOC-->

## Table of Contents

- [注解顺序](#注解顺序) `:21+19`
- [关键规则](#关键规则) `:40+11`
- [Query 参数 Swagger 注解](#query-参数-swagger-注解) `:51+23`

<!--TOC-->

## 注解顺序

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

## 关键规则

| 规则       | 说明                                                             |
| ---------- | ---------------------------------------------------------------- |
| Tags 格式  | `中文名 (English Name)`，影响生成的 API 类名                     |
| 响应类型   | 必须用 `response.DataResponse[T]` 或 `response.PagedResponse[T]` |
| 路径参数   | 用 `{id}` 而非 `:id`                                             |
| 认证端点   | 必须加 `@Security BearerAuth`                                    |
| Query 参数 | 使用 `@Param params query handler.StructName` 格式               |
| 枚举值     | 使用 `enums:"val1,val2"` 标签，不要用 `example`                  |

## Query 参数 Swagger 注解

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
