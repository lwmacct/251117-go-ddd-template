# Handler 开发规范

<!--TOC-->

## Table of Contents

- [响应封装原则](#响应封装原则) `:14+14`
- [Swagger 注解规范](#swagger-注解规范) `:28+30`
  - [注解顺序](#注解顺序) `:34+16`
  - [关键规则](#关键规则) `:50+8`

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

## Swagger 注解规范

> **⚠️ 注解直接影响前端代码生成！**
>
> 前端通过 `pnpm api:generate` 从 swagger.json 生成 TypeScript 客户端到 `src/api/generated/`

### 注解顺序

```go
// @Summary      简短描述
// @Description  详细说明
// @Tags         管理员 - 用户管理 (Admin - User Management)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "用户ID"
// @Param        request body user.CreateUserDTO true "请求体"
// @Success      200 {object} response.DataResponse[user.UserDTO] "成功"
// @Failure      400 {object} response.ErrorResponse "参数错误"
// @Router       /api/admin/users/{id} [get]
```

### 关键规则

| 规则      | 说明                                                             |
| --------- | ---------------------------------------------------------------- |
| Tags 格式 | `中文名 (English Name)`，影响生成的 API 类名                     |
| 响应类型  | 必须用 `response.DataResponse[T]` 或 `response.PagedResponse[T]` |
| 路径参数  | 用 `{id}` 而非 `:id`                                             |
| 认证端点  | 必须加 `@Security BearerAuth`                                    |
