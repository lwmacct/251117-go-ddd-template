# RBAC 权限系统

本文档介绍 RBAC（基于角色的访问控制）权限系统的工作原理和使用方法。

<!--TOC-->

## Table of Contents

- [概述](#概述) `:48+14`
- [权限模型](#权限模型) `:62+41`
  - [三段式权限格式](#三段式权限格式) `:64+20`
  - [通配符匹配规则](#通配符匹配规则) `:84+12`
  - [预设角色](#预设角色) `:96+7`
- [JWT Token 与权限](#jwt-token-与权限) `:103+26`
  - [Claims 结构](#claims-结构) `:105+15`
  - [登录流程](#登录流程) `:120+9`
- [Personal Access Token](#personal-access-token) `:129+23`
  - [PAT vs JWT](#pat-vs-jwt) `:133+10`
  - [Token 格式](#token-格式) `:143+9`
- [中间件系统](#中间件系统) `:152+60`
  - [执行顺序](#执行顺序) `:154+7`
  - [Auth 中间件](#auth-中间件) `:161+12`
  - [RBAC 中间件](#rbac-中间件) `:173+35`
  - [AuditMiddleware](#auditmiddleware) `:208+4`
- [API 路由保护](#api-路由保护) `:212+29`
  - [路由配置示例](#路由配置示例) `:214+27`
- [使用指南](#使用指南) `:241+55`
  - [1. 登录获取 Token](#1-登录获取-token) `:243+8`
  - [2. 查看角色](#2-查看角色) `:251+7`
  - [3. 为用户分配角色](#3-为用户分配角色) `:258+8`
  - [4. 创建自定义角色](#4-创建自定义角色) `:266+8`
  - [5. 查看权限](#5-查看权限) `:274+7`
  - [6. 为角色分配权限](#6-为角色分配权限) `:281+8`
  - [7. 查看审计日志](#7-查看审计日志) `:289+7`
- [最佳实践](#最佳实践) `:296+16`
  - [安全建议](#安全建议) `:298+8`
  - [权限设计原则](#权限设计原则) `:306+6`
- [常见问题](#常见问题) `:312+44`
  - [Q1: 权限变更后何时生效？](#q1-权限变更后何时生效) `:314+5`
  - [Q2: 如何实现细粒度权限？](#q2-如何实现细粒度权限) `:319+11`
  - [Q3: 如何让用户只能修改自己的资料？](#q3-如何让用户只能修改自己的资料) `:330+11`
  - [Q4: Token 被盗用怎么办？](#q4-token-被盗用怎么办) `:341+6`
  - [Q5: 如何添加新权限？](#q5-如何添加新权限) `:347+9`
- [技术架构](#技术架构) `:356+19`

<!--TOC-->

## 概述

RBAC (Role-Based Access Control) 通过将权限分配给角色，再将角色分配给用户来管理系统权限。

**核心关系**：User ↔ Role ↔ Permission（多对多）

**系统特性**：

- ✅ 多角色支持（一用户多角色）
- ✅ 细粒度权限（三段式 `domain:resource:action`）
- ✅ 通配符匹配（`admin:users:*`、`*:*:*`）
- ✅ 双重认证（JWT + PAT）
- ✅ 完整审计日志

## 权限模型

### 三段式权限格式

格式：`domain:resource:action`

| 段       | 说明     | 示例                                      |
| -------- | -------- | ----------------------------------------- |
| Domain   | 权限域   | `admin`（管理后台）、`user`（用户自服务） |
| Resource | 操作对象 | `users`、`roles`、`profile`               |
| Action   | 具体操作 | `create`、`read`、`update`、`delete`      |

**权限示例**：

| 权限代码              | 描述               |
| --------------------- | ------------------ |
| `admin:users:create`  | 管理员创建用户     |
| `admin:users:*`       | 用户资源的所有操作 |
| `admin:*:create`      | 所有资源的创建操作 |
| `user:profile:update` | 用户更新自己资料   |
| `*:*:*`               | 超级管理员权限     |

### 通配符匹配规则

- 用户权限中的 `*` 匹配任意值
- 从左到右逐段比对：domain → resource → action

```
用户权限: "admin:users:*"
RequirePermission("admin:users:create")  ✓
RequirePermission("admin:users:delete")  ✓
RequirePermission("admin:roles:create")  ✗ (资源不匹配)
```

### 预设角色

| 角色    | 权限范围              |
| ------- | --------------------- |
| `admin` | 所有 `admin:*:*` 权限 |
| `user`  | 所有 `user:*:*` 权限  |

## JWT Token 与权限

### Claims 结构

登录成功后生成的 JWT Token 包含：

```json
{
  "user_id": 1,
  "username": "admin",
  "email": "admin@example.com",
  "roles": ["admin"],
  "permissions": ["admin:users:create", "admin:users:read", ...],
  "exp": 1672531200
}
```

### 登录流程

1. 查询用户（预加载角色和权限）
2. 验证密码（bcrypt）
3. 提取角色和权限列表
4. 生成 JWT Token Pair（access + refresh）

数据库使用嵌套预加载 `Preload("Roles.Permissions")` 避免 N+1 查询。

## Personal Access Token

PAT 适用于 API 集成、CLI 工具、自动化脚本等场景。

### PAT vs JWT

| 特性     | JWT                | PAT                       |
| -------- | ------------------ | ------------------------- |
| 用途     | Web/移动应用       | API/CLI/自动化            |
| 有效期   | 短期（1小时）      | 可选（7/30/90天或永久）   |
| 权限范围 | 用户全部权限       | 创建时选择的子集          |
| 格式     | `Bearer eyJhbG...` | `Bearer pat_xxxxx_yyy...` |
| IP 限制  | 不支持             | 支持白名单                |

### Token 格式

`pat_<5位前缀>_<32位随机字符>`

- 完整 token 仅创建时显示一次
- 数据库存储 SHA-256 哈希

详见 [PAT 使用指南](./identity-pat.md)

## 中间件系统

### 执行顺序

```
请求 → Auth (JWT/PAT) → AuditMiddleware → RBAC → Handler
        验证Token       记录日志         检查权限
```

### Auth 中间件

自动识别 JWT 和 PAT，将用户信息注入 Context：

```go
c.Set("user_id", userID)
c.Set("username", username)
c.Set("roles", roles)           // 角色列表
c.Set("permissions", permissions) // 权限列表
c.Set("auth_type", "jwt"|"pat")
```

### RBAC 中间件

**RequireRole** - 角色检查：

```go
admin.Use(middleware.RequireRole("admin"))
```

**RequirePermission** - 权限检查（支持通配符）：

```go
router.POST("/users",
    middleware.RequirePermission("admin:users:create"),
    handler.CreateUser,
)
```

**RequireOwnership** - 所有权检查：

```go
router.PUT("/users/:id",
    middleware.RequireOwnership(),  // 只能修改自己
    handler.UpdateUser,
)
```

**RequireAdminOrOwnership** - 管理员或所有者：

```go
router.PUT("/users/:id",
    middleware.RequireAdminOrOwnership(),
    handler.UpdateUser,
)
```

### AuditMiddleware

自动记录所有写操作（POST/PUT/DELETE）到审计日志。

## API 路由保护

### 路由配置示例

```go
// 公开路由（无需认证）
auth := api.Group("/auth")
auth.POST("/login", authHandler.Login)

// 管理员路由
admin := api.Group("/admin")
admin.Use(middleware.JWTAuth(jwtManager))
admin.Use(middleware.AuditMiddleware(auditLogRepo))
admin.Use(middleware.RequireRole("admin"))
{
    admin.POST("/users", adminUserHandler.CreateUser)
    admin.GET("/users", adminUserHandler.ListUsers)
    admin.PUT("/users/:id/roles", adminUserHandler.AssignRoles)
}

// 用户路由（仅需认证）
userGroup := api.Group("/user")
userGroup.Use(middleware.JWTAuth(jwtManager))
{
    userGroup.GET("/me", userProfileHandler.GetProfile)
    userGroup.PUT("/me", userProfileHandler.UpdateProfile)
}
```

## 使用指南

### 1. 登录获取 Token

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login": "admin", "password": "admin123"}'
```

### 2. 查看角色

```bash
curl -X GET http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer <token>"
```

### 3. 为用户分配角色

```bash
curl -X PUT http://localhost:8080/api/admin/users/5/roles \
  -H "Authorization: Bearer <token>" \
  -d '{"role_ids": [1, 2]}'
```

### 4. 创建自定义角色

```bash
curl -X POST http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer <token>" \
  -d '{"name": "editor", "display_name": "编辑"}'
```

### 5. 查看权限

```bash
curl -X GET http://localhost:8080/api/admin/permissions \
  -H "Authorization: Bearer <token>"
```

### 6. 为角色分配权限

```bash
curl -X PUT http://localhost:8080/api/admin/roles/3/permissions \
  -H "Authorization: Bearer <token>" \
  -d '{"permission_ids": [1, 2, 3]}'
```

### 7. 查看审计日志

```bash
curl -X GET "http://localhost:8080/api/admin/audit-logs?user_id=1" \
  -H "Authorization: Bearer <token>"
```

## 最佳实践

### 安全建议

1. **立即修改默认密码** - 默认 `admin123` 仅用于初始化
2. **定期刷新 Token** - Access Token 默认 1 小时有效
3. **生产环境使用 HTTPS** - 防止 Token 被窃取
4. **最小权限原则** - 只分配必要权限
5. **定期审查审计日志** - 关注异常操作

### 权限设计原则

1. **通过角色分配权限**：User → Role → Permissions
2. **命名规范**：资源用复数（`users`），操作用动词（`create`）
3. **粗细结合**：管理员用角色检查，特殊接口用权限检查

## 常见问题

### Q1: 权限变更后何时生效？

权限存储在 JWT 中，变更后需重新登录或刷新 Token。
改进方案：缩短 Token 有效期、使用 Redis 实时检查。

### Q2: 如何实现细粒度权限？

使用 `RequirePermission` 中间件：

```go
router.DELETE("/articles/:id",
    middleware.RequirePermission("articles:delete"),
    handler.DeleteArticle,
)
```

### Q3: 如何让用户只能修改自己的资料？

使用 `RequireAdminOrOwnership` 中间件：

```go
router.PUT("/users/:id",
    middleware.RequireAdminOrOwnership(),
    handler.UpdateUser,
)
```

### Q4: Token 被盗用怎么办？

- 使用 HTTPS
- 短有效期（1小时）
- 可选：IP 绑定、设备绑定、Token 黑名单

### Q5: 如何添加新权限？

通过 API：

```bash
curl -X POST http://localhost:8080/api/admin/permissions \
  -d '{"resource": "articles", "action": "publish", "code": "articles:publish"}'
```

## 技术架构

**优势**：

- ✅ JWT 无状态认证，易扩展
- ✅ 完整 RBAC 模型
- ✅ 灵活中间件设计
- ✅ 自动审计日志

**限制**：

- ⚠️ 权限变更需重新登录（可用 Redis 改进）
- ⚠️ 权限多时 Token 较大

**相关代码**：

- 中间件：`internal/adapters/http/middleware/rbac.go`
- JWT：`internal/infrastructure/auth/jwt.go`
- 路由：`internal/adapters/http/router.go`
