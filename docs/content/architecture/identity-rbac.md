# RBAC 权限系统工作原理

本文档详细介绍本系统的 RBAC（基于角色的访问控制）权限系统的工作原理、实现方式和使用方法。

<!--TOC-->

## Table of Contents

- [概述](#概述) `:68+29`
  - [什么是 RBAC？](#什么是-rbac) `:70+13`
  - [核心概念](#核心概念) `:83+14`
- [权限模型](#权限模型) `:97+181`
  - [领域实体与持久化模型](#领域实体与持久化模型) `:99+98`
  - [权限代码格式](#权限代码格式) `:197+39`
  - [通配符权限匹配](#通配符权限匹配) `:236+33`
  - [预设角色](#预设角色) `:269+9`
- [JWT Token 与权限流转](#jwt-token-与权限流转) `:278+91`
  - [JWT Claims 结构](#jwt-claims-结构) `:280+29`
  - [登录时的权限加载流程](#登录时的权限加载流程) `:309+60`
- [Personal Access Token 认证](#personal-access-token-认证) `:369+94`
  - [PAT vs JWT 对比](#pat-vs-jwt-对比) `:373+13`
  - [PAT Token 格式](#pat-token-格式) `:386+18`
  - [统一认证中间件](#统一认证中间件) `:404+23`
  - [PAT 认证流程](#pat-认证流程) `:427+36`
- [中间件系统](#中间件系统) `:463+266`
  - [中间件执行顺序](#中间件执行顺序) `:465+9`
  - [1. Auth 中间件（统一认证）](#1-auth-中间件统一认证) `:474+36`
  - [2. RBAC 中间件（授权）](#2-rbac-中间件授权) `:510+180`
  - [3. AuditMiddleware（审计日志）](#3-auditmiddleware审计日志) `:690+39`
- [API 路由保护](#api-路由保护) `:729+91`
  - [路由配置](#路由配置) `:731+47`
  - [管理员接口保护示例](#管理员接口保护示例) `:778+24`
  - [用户接口保护示例](#用户接口保护示例) `:802+18`
- [完整的请求流程](#完整的请求流程) `:820+204`
  - [流程图](#流程图) `:822+26`
  - [成功案例：管理员创建用户](#成功案例管理员创建用户) `:848+82`
  - [失败案例：普通用户访问管理员接口](#失败案例普通用户访问管理员接口) `:930+45`
  - [所有权检查案例：用户修改自己的资料](#所有权检查案例用户修改自己的资料) `:975+49`
- [系统初始化](#系统初始化) `:1024+70`
  - [数据库迁移](#数据库迁移) `:1026+23`
  - [RBAC 种子数据](#rbac-种子数据) `:1049+45`
- [使用指南](#使用指南) `:1094+176`
  - [1. 登录获取 Token](#1-登录获取-token) `:1096+30`
  - [2. 为用户分配角色](#2-为用户分配角色) `:1126+44`
  - [3. 创建自定义角色](#3-创建自定义角色) `:1170+13`
  - [4. 为角色分配权限](#4-为角色分配权限) `:1183+45`
  - [5. 查看审计日志](#5-查看审计日志) `:1228+42`
- [最佳实践](#最佳实践) `:1270+89`
  - [安全建议](#安全建议) `:1272+28`
  - [权限设计原则](#权限设计原则) `:1300+18`
  - [Token 管理建议](#token-管理建议) `:1318+22`
  - [审计日志使用](#审计日志使用) `:1340+19`
- [常见问题](#常见问题) `:1359+135`
  - [Q1: 权限变更后何时生效？](#q1-权限变更后何时生效) `:1361+13`
  - [Q2: 如何实现更细粒度的权限控制？](#q2-如何实现更细粒度的权限控制) `:1374+29`
  - [Q3: 如何让用户可以修改自己的资料但不能修改他人？](#q3-如何让用户可以修改自己的资料但不能修改他人) `:1403+17`
  - [Q4: 如何处理跨域请求？](#q4-如何处理跨域请求) `:1420+18`
  - [Q5: Token 被盗用怎么办？](#q5-token-被盗用怎么办) `:1438+23`
  - [Q6: 如何添加新的权限？](#q6-如何添加新的权限) `:1461+33`
- [技术架构总结](#技术架构总结) `:1494+24`
  - [优势](#优势) `:1496+8`
  - [限制](#限制) `:1504+6`
  - [扩展方向](#扩展方向) `:1510+8`
- [参考资料](#参考资料) `:1518+7`

<!--TOC-->

## 概述

### 什么是 RBAC？

RBAC (Role-Based Access Control) 是一种基于角色的访问控制模型，通过将权限分配给角色，再将角色分配给用户的方式来管理系统权限。

本系统采用完整的 RBAC 模型，支持：

- ✅ 多角色支持（一个用户可以拥有多个角色）
- ✅ 细粒度权限控制（基于三段式 `domain:resource:action` 格式）
- ✅ 通配符权限匹配（支持 `admin:users:*`、`admin:*:create`、`*:*:*` 等）
- ✅ 双重认证机制（JWT + Personal Access Token）
- ✅ 完整的审计日志
- ✅ 物理分离的管理员和用户接口

### 核心概念

```
User（用户）
  ↓ 多对多关系
Role（角色）
  ↓ 多对多关系
Permission（权限）
```

- **User（用户）**：系统的使用者
- **Role（角色）**：权限的集合，如 `admin`（管理员）、`user`（普通用户）
- **Permission（权限）**：具体的操作许可，采用三段式格式，如 `admin:users:create`（管理员域-创建用户）

## 权限模型

### 领域实体与持久化模型

#### Domain 层实体（无 GORM Tag）

```go
// internal/domain/role/entity_role.go
type Role struct {
    ID          uint
    Name        string
    DisplayName string
    Description string
    IsSystem    bool
    Permissions []Permission
}

// internal/domain/role/entity_permission.go
type Permission struct {
    ID          uint
    Domain      string
    Resource    string
    Action      string
    Code        string
    Description string
}

// internal/domain/user/entity_user.go
type User struct {
    ID       uint
    Username string
    Email    string
    Password string
    Status   string
    Roles    []role.Role
}
```

- Domain 实体只承载业务字段与行为（如 `User.AssignRole`），不包含任何 ORM/GORM 相关细节。
- 领域行为（激活用户、分配角色等）都在实体方法中完成，Use Case 通过这些方法保持业务规则一致性。

#### Infrastructure 层 `*_model.go`

```go
// internal/infrastructure/persistence/role_model.go
type RoleModel struct {
    ID          uint           `gorm:"primaryKey"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
    Name        string         `gorm:"size:50;uniqueIndex;not null"`
    DisplayName string         `gorm:"size:100;not null"`
    Description string         `gorm:"size:255"`
    IsSystem    bool           `gorm:"default:false;not null"`
    Permissions []PermissionModel `gorm:"many2many:role_permissions"`
}
```

```go
// internal/infrastructure/persistence/permission_model.go
type PermissionModel struct {
    ID          uint           `gorm:"primaryKey"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
    Domain      string         `gorm:"size:50;not null"`
    Resource    string         `gorm:"size:50;not null"`
    Action      string         `gorm:"size:50;not null"`
    Code        string         `gorm:"size:150;uniqueIndex;not null"`
    Description string         `gorm:"size:255"`
}
```

```go
// internal/infrastructure/persistence/user_model.go
type UserModel struct {
    ID       uint           `gorm:"primaryKey"`
    Username string         `gorm:"uniqueIndex;size:50;not null"`
    Email    string         `gorm:"uniqueIndex;size:100;not null"`
    Password string         `gorm:"size:255;not null"`
    Status   string         `gorm:"size:20;default:'active'"`
    Roles    []RoleModel    `gorm:"many2many:user_roles"`
}
```

- 所有数据库约束（索引、关联、软删除）集中在持久化模型，领域层保持纯粹。
- `mapRoleModelsToEntities` / `mapRoleEntitiesToModels` 等辅助方法负责在仓储中完成 Model ↔ Entity 转换。
- 关联表 `user_roles`、`role_permissions` 由 GORM 根据模型声明自动维护。

#### Repository 分离

| 接口                     | 位置                                         | 说明                                                        |
| ------------------------ | -------------------------------------------- | ----------------------------------------------------------- |
| `role.CommandRepository` | `internal/domain/role/command_repository.go` | 写操作：创建/更新/删除角色与权限，全部使用 GORM Model 写入  |
| `role.QueryRepository`   | `internal/domain/role/query_repository.go`   | 读操作：`GetByIDWithPermissions`、`List`、`ExistsByName` 等 |
| `user.CommandRepository` | `internal/domain/user/command_repository.go` | 用户写操作（激活、分配角色）                                |
| `user.QueryRepository`   | `internal/domain/user/query_repository.go`   | 用户读操作（含角色/权限预加载）                             |

每个仓储的 Infrastructure 实现在 `internal/infrastructure/persistence/*_repository.go`，并在构造函数中显式注入，确保 Command/Query 读写分离。

### 权限代码格式

权限代码采用三段式 `domain:resource:action` 格式：

**格式说明**：

- **Domain（域）**：权限所属的领域，如 `admin`（管理后台）、`user`（用户自服务）、`api`（API 接口）
- **Resource（资源）**：操作的对象，如 `users`（用户）、`roles`（角色）、`profile`（个人资料）
- **Action（动作）**：具体操作，如 `create`（创建）、`read`（读取）、`update`（更新）、`delete`（删除）

**权限示例**：

| 权限代码                 | 域    | 资源        | 动作   | 描述                   |
| ------------------------ | ----- | ----------- | ------ | ---------------------- |
| `admin:users:create`     | admin | users       | create | 管理员创建用户         |
| `admin:users:read`       | admin | users       | read   | 管理员查看所有用户     |
| `admin:users:update`     | admin | users       | update | 管理员更新用户信息     |
| `admin:users:delete`     | admin | users       | delete | 管理员删除用户         |
| `admin:roles:read`       | admin | roles       | read   | 管理员查看角色         |
| `admin:settings:create`  | admin | settings    | create | 管理员创建系统配置     |
| `admin:settings:read`    | admin | settings    | read   | 管理员查看系统配置     |
| `admin:settings:update`  | admin | settings    | update | 管理员更新系统配置     |
| `admin:settings:delete`  | admin | settings    | delete | 管理员删除系统配置     |
| `admin:overview:read`    | admin | overview    | read   | 查看系统概览数据       |
| `admin:permissions:read` | admin | permissions | read   | 管理员查看权限列表     |
| `admin:menus:create`     | admin | menus       | create | 管理员创建菜单         |
| `admin:menus:read`       | admin | menus       | read   | 管理员查看菜单         |
| `admin:menus:update`     | admin | menus       | update | 管理员更新菜单         |
| `admin:menus:delete`     | admin | menus       | delete | 管理员删除菜单         |
| `admin:audit_logs:read`  | admin | audit_logs  | read   | 管理员查看审计日志     |
| `user:profile:read`      | user  | profile     | read   | 用户查看自己的资料     |
| `user:profile:update`    | user  | profile     | update | 用户更新自己的资料     |
| `user:password:update`   | user  | password    | update | 用户修改自己的密码     |
| `user:tokens:create`     | user  | tokens      | create | 用户创建个人访问令牌   |
| `user:tokens:read`       | user  | tokens      | read   | 用户查看自己的令牌列表 |
| `user:tokens:delete`     | user  | tokens      | delete | 用户撤销自己的令牌     |
| `api:cache:read`         | api   | cache       | read   | API 读取缓存           |
| `api:cache:write`        | api   | cache       | write  | API 写入缓存           |

### 通配符权限匹配

系统支持在权限检查时使用通配符 `*`，实现灵活的权限委派：

| 通配符权限       | 匹配范围             | 示例                                |
| ---------------- | -------------------- | ----------------------------------- |
| `admin:users:*`  | 用户资源的所有操作   | 可执行 create、read、update、delete |
| `admin:*:create` | 所有资源的创建操作   | 可创建 users、roles 等              |
| `user:*:*`       | 用户域的所有权限     | 用户自服务的全部功能                |
| `*:users:read`   | 所有域的用户读取权限 | admin 和 user 域都可读              |
| `*:*:*`          | 超级管理员           | 系统所有权限                        |

**匹配规则**：

- 用户权限中的 `*` 可以匹配任意值
- 检查时从左到右逐段比对：domain、resource、action
- 任何一段为 `*` 即匹配该段的所有值

**示例**：

```go
// 用户拥有权限: "admin:users:*"
RequirePermission("admin:users:create")  // ✓ 通过
RequirePermission("admin:users:read")    // ✓ 通过
RequirePermission("admin:users:delete")  // ✓ 通过
RequirePermission("admin:roles:create")  // ✗ 不通过（资源不匹配）

// 用户拥有权限: "admin:*:create"
RequirePermission("admin:users:create")  // ✓ 通过
RequirePermission("admin:roles:create")  // ✓ 通过
RequirePermission("admin:users:update")  // ✗ 不通过（动作不匹配）
```

### 预设角色

系统默认提供两个系统角色：

| 角色名  | 显示名称 | 权限              | 说明                   |
| ------- | -------- | ----------------- | ---------------------- |
| `admin` | 管理员   | 所有 admin 域权限 | 拥有管理后台全部权限   |
| `user`  | 普通用户 | 所有 user 域权限  | 拥有用户自服务全部权限 |

## JWT Token 与权限流转

### JWT Claims 结构

登录成功后，系统会生成包含用户权限信息的 JWT Token：

```go
type Claims struct {
    UserID      uint     `json:"user_id"`      // 用户 ID
    Username    string   `json:"username"`     // 用户名
    Email       string   `json:"email"`        // 邮箱
    Roles       []string `json:"roles"`        // 角色列表 ["admin"]
    Permissions []string `json:"permissions"`  // 权限列表 ["users:create", "users:read", ...]
    jwt.RegisteredClaims
}
```

**Token 示例**（解码后）：

```json
{
  "user_id": 1,
  "username": "admin",
  "email": "admin@example.com",
  "roles": ["admin"],
  "permissions": ["admin:users:create", "admin:users:read", "admin:users:update", "admin:users:delete", "admin:roles:create", "admin:roles:read", "admin:roles:update", "admin:roles:delete", "admin:permissions:read", "admin:audit_logs:read"],
  "exp": 1672531200,
  "iat": 1672444800
}
```

### 登录时的权限加载流程

```mermaid
sequenceDiagram
    participant Client
    participant AuthHandler
    participant AuthService
    participant UserRepo
    participant Database
    participant JWTManager

    Client->>AuthHandler: POST /api/auth/login
    AuthHandler->>AuthService: Login(username, password)
    AuthService->>UserRepo: GetByUsernameWithRoles(username)
    UserRepo->>Database: SELECT * FROM users<br/>JOIN user_roles<br/>JOIN roles<br/>JOIN role_permissions<br/>JOIN permissions
    Database-->>UserRepo: User with Roles & Permissions
    UserRepo-->>AuthService: User object
    AuthService->>AuthService: Verify password (bcrypt)
    AuthService->>AuthService: Extract roles & permissions
    AuthService->>JWTManager: GenerateTokenPair(userID, username, email, roles, permissions)
    JWTManager-->>AuthService: Access Token & Refresh Token
    AuthService-->>AuthHandler: AuthResponse
    AuthHandler-->>Client: 200 OK + Tokens
```

**关键代码**（`internal/infrastructure/auth/service.go:100-150`）：

```go
// 1. 查询用户（预加载角色和权限）
u, err := s.userQueryRepo.GetByUsernameWithRoles(ctx, req.Login)

// 2. 验证密码
bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))

// 3. 提取角色和权限
roles := u.GetRoleNames()           // ["admin"]
permissions := u.GetPermissionCodes() // ["users:create", ...]

// 4. 生成 JWT Token
accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(
    u.ID, u.Username, u.Email, roles, permissions
)
```

**数据库查询优化**（使用嵌套预加载）：

```go
// internal/infrastructure/persistence/user_query_repository.go:60-75
func (r *userQueryRepository) GetByUsernameWithRoles(ctx context.Context, username string) (*user.User, error) {
    var u user.User
    err := r.db.WithContext(ctx).
        Preload("Roles.Permissions"). // 嵌套预加载角色及其权限
        Where("username = ?", username).
        First(&u).Error
    return &u, err
}
```

这样一次查询即可获取用户的所有角色和权限信息，避免 N+1 查询问题。

## Personal Access Token 认证

除了 JWT 认证外，系统还支持 Personal Access Token (PAT) 认证方式，适用于 API 集成、CLI 工具、Webhook 回调等场景。

### PAT vs JWT 对比

| 特性     | JWT Token                | Personal Access Token          |
| -------- | ------------------------ | ------------------------------ |
| 用途     | Web 应用、移动应用       | API 集成、CLI 工具、自动化脚本 |
| 有效期   | 短期（默认 1 小时）      | 可选（7/30/90 天或永久）       |
| 刷新机制 | 使用 Refresh Token       | 无需刷新，到期后重新创建       |
| 权限范围 | 用户全部权限             | 用户权限的子集（创建时选择）   |
| 格式     | `Bearer eyJhbGc...`      | `Bearer pat_xxxxx_yyy...`      |
| 存储方式 | 不存储（无状态）         | 数据库存储（SHA-256 哈希）     |
| IP 限制  | 不支持                   | 支持 IP 白名单                 |
| 撤销     | 不支持（除非加入黑名单） | 支持随时撤销                   |

### PAT Token 格式

**格式**: `pat_<5位前缀>_<32位随机字符>`

**示例**: `pat_2Kj9X_aB3cD4eF5gH6iJ7kL8mN9oP0qR1sT2uV3wX4yZ`

**组成部分**：

1. `pat_` - 固定前缀，用于识别 PAT 类型
2. `2Kj9X` - 5 位随机前缀，用于用户快速识别 token
3. `aB3cD...` - 32 位随机字符，实际的 token 内容

**安全设计**：

- 完整 token 仅在创建时显示一次
- 数据库存储 SHA-256 哈希值，不存储明文
- 前缀用于用户识别，不影响安全性

### 统一认证中间件

系统使用统一的 `Auth` 中间件自动识别 JWT 和 PAT：

```go
// 中间件会自动判断 token 类型
func Auth(jwtManager *JWTManager, patService *PATService, userQueryRepo user.QueryRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")  // "Bearer <token>"
        tokenString := extractToken(authHeader)

        // 判断 token 类型
        if strings.HasPrefix(tokenString, "pat_") {
            // PAT 认证
            authenticateWithPAT(c, patService, userQueryRepo, tokenString)
        } else {
            // JWT 认证
            authenticateWithJWT(c, jwtManager, tokenString)
        }
    }
}
```

### PAT 认证流程

**请求示例**:

```bash
curl -X GET http://localhost:8080/api/user/me \
  -H "Authorization: Bearer pat_2Kj9X_aB3cD4eF5gH6iJ7kL8mN9oP0qR1sT2uV3wX4yZ"
```

**认证步骤**:

1. 中间件检测到 `pat_` 前缀，识别为 PAT
2. 使用 SHA-256 哈希 token
3. 在数据库中查找匹配的 token 记录
4. 检查 token 状态：
   - 是否已撤销 (status != "active")
   - 是否已过期 (expires_at < now)
   - IP 是否在白名单内（如果配置了）
5. 加载用户信息和角色
6. 将 token 的权限列表（用户权限的子集）注入 Context
7. 异步更新 `last_used_at` 时间戳

**注入到 Context 的信息**:

```go
c.Set("user_id", u.ID)
c.Set("username", u.Username)
c.Set("email", u.Email)
c.Set("roles", u.GetRoleNames())
c.Set("permissions", pat.Permissions)  // 注意：是 PAT 的权限子集
c.Set("auth_type", "pat")
c.Set("pat_id", pat.ID)
```

详细的 PAT 使用说明请参考 [Personal Access Token 文档](./identity-pat.md)。

## 中间件系统

### 中间件执行顺序

```
请求 → Auth (JWT/PAT) → AuditMiddleware → RBAC Middleware → Handler
        ↓                  ↓                  ↓
     验证 Token           记录日志          检查权限 (三段式)
     注入 Context
```

### 1. Auth 中间件（统一认证）

**文件**: `internal/adapters/http/middleware/jwt.go`

**功能**: 统一验证 JWT Token 和 PAT，并将用户信息注入到 Gin Context

```go
func Auth(jwtManager *JWTManager, patService *PATService, userQueryRepo user.QueryRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 提取 Authorization Header
        authHeader := c.GetHeader("Authorization")  // "Bearer <token>"

        tokenString := extractToken(authHeader)

        // 2. 判断 token 类型并验证
        if strings.HasPrefix(tokenString, "pat_") {
            authenticateWithPAT(c, patService, userQueryRepo, tokenString)
        } else {
            authenticateWithJWT(c, jwtManager, tokenString)
        }

        // 3. 注入用户信息到 Context（关键！）
        // JWT: 包含用户全部权限
        // PAT: 包含 token 创建时选择的权限子集
        c.Set("user_id", userID)
        c.Set("username", username)
        c.Set("email", email)
        c.Set("roles", roles)           // 角色列表
        c.Set("permissions", permissions) // 权限列表（三段式格式）
        c.Set("auth_type", authType)    // "jwt" 或 "pat"

        c.Next()  // 传递给下一个中间件
    }
}
```

### 2. RBAC 中间件（授权）

**文件**: `internal/adapters/http/middleware/rbac.go`

#### RequireRole - 角色检查

检查用户是否拥有指定角色：

```go
func RequireRole(role string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从 Context 读取角色列表
        roles, exists := c.Get("roles")
        rolesList := roles.([]string)

        // 检查是否包含所需角色
        hasRole := false
        for _, r := range rolesList {
            if r == role {
                hasRole = true
                break
            }
        }

        if !hasRole {
            c.JSON(403, gin.H{"error": "forbidden: insufficient permissions"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

**使用示例**:

```go
admin.Use(middleware.RequireRole("admin"))  // 要求 admin 角色
```

#### RequirePermission - 权限检查（支持通配符）

检查用户是否拥有指定权限，支持三段式格式和通配符匹配：

```go
func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        permissions, _ := c.Get("permissions")
        permissionsList := permissions.([]string)

        hasPermission := false
        for _, p := range permissionsList {
            if matchPermission(p, permission) {  // 支持通配符匹配
                hasPermission = true
                break
            }
        }

        if !hasPermission {
            c.JSON(403, gin.H{"error": "forbidden: insufficient permissions"})
            c.Abort()
            return
        }

        c.Next()
    }
}

// matchPermission 支持三段式通配符匹配
func matchPermission(userPerm, requiredPerm string) bool {
    // 精确匹配
    if userPerm == requiredPerm {
        return true
    }

    // 分割为三部分：domain:resource:action
    userParts := strings.Split(userPerm, ":")
    requiredParts := strings.Split(requiredPerm, ":")

    // 必须是三段式格式
    if len(userParts) != 3 || len(requiredParts) != 3 {
        return false
    }

    // 逐段匹配，用户权限中的 * 可匹配任意值
    for i := 0; i < 3; i++ {
        if userParts[i] == "*" {
            continue  // 通配符匹配
        }
        if userParts[i] != requiredParts[i] {
            return false
        }
    }

    return true
}
```

**使用示例**:

```go
// 三段式权限
router.POST("/users",
    middleware.RequirePermission("admin:users:create"),
    handler.CreateUser,
)

// 通配符示例
// 用户拥有 "admin:users:*"，可以通过以下检查：
// - RequirePermission("admin:users:create") ✓
// - RequirePermission("admin:users:read")   ✓
// - RequirePermission("admin:users:delete") ✓
```

#### RequireOwnership - 所有权检查

检查用户是否访问自己的资源：

```go
func RequireOwnership(paramName ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取当前用户 ID
        userID, _ := c.Get("user_id")
        uid := userID.(uint)

        // 获取 URL 参数中的资源 ID
        resourceID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

        // 比较是否为同一用户
        if uint(resourceID) != uid {
            c.JSON(403, gin.H{"error": "forbidden: can only access own resources"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

**使用示例**:

```go
router.PUT("/users/:id",
    middleware.RequireOwnership(),  // 只能修改自己的资料
    handler.UpdateUser,
)
```

#### RequireAdminOrOwnership - 管理员或所有者

允许管理员或资源所有者访问：

```go
func RequireAdminOrOwnership(paramName ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 先检查是否为管理员
        roles, _ := c.Get("roles")
        for _, role := range roles.([]string) {
            if role == "admin" {
                c.Next()  // 管理员直接通过
                return
            }
        }

        // 不是管理员，检查是否为资源所有者
        userID, _ := c.Get("user_id")
        resourceID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

        if uint(resourceID) == userID.(uint) {
            c.Next()  // 是所有者，通过
        } else {
            c.JSON(403, gin.H{"error": "forbidden: insufficient permissions"})
            c.Abort()
        }
    }
}
```

### 3. AuditMiddleware（审计日志）

**文件**: `internal/adapters/http/middleware/audit.go`

**功能**: 自动记录所有写操作（POST、PUT、DELETE）

```go
func AuditMiddleware(repo auditlog.Repository) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 跳过读操作
        if c.Request.Method == "GET" {
            c.Next()
            return
        }

        // 提取用户信息
        userID, _ := c.Get("user_id")
        username, _ := c.Get("username")

        // 执行请求
        c.Next()

        // 异步记录审计日志
        go func() {
            log := &auditlog.AuditLog{
                UserID:     userID.(uint),
                Username:   username.(string),
                Action:     fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
                Resource:   extractResource(c.Request.URL.Path),
                IPAddress:  c.ClientIP(),
                UserAgent:  c.Request.UserAgent(),
                Status:     determineStatus(c.Writer.Status()),
            }
            repo.Create(context.Background(), log)
        }()
    }
}
```

## API 路由保护

### 路由配置

**文件**: `internal/adapters/http/router.go`

```go
// 1. 公开路由（无需认证）
auth := api.Group("/auth")
{
    auth.POST("/register", authHandler.Register)
    auth.POST("/login", authHandler.Login)
    auth.POST("/refresh", authHandler.RefreshToken)
}

// 2. 管理员路由（需要 admin 角色）
admin := api.Group("/admin")
admin.Use(middleware.JWTAuth(jwtManager))          // JWT 认证
admin.Use(middleware.AuditMiddleware(auditLogRepo))// 审计日志
admin.Use(middleware.RequireRole("admin"))         // 角色检查
{
    // 用户管理
    admin.POST("/users", adminUserHandler.CreateUser)
    admin.GET("/users", adminUserHandler.ListUsers)
    admin.GET("/users/:id", adminUserHandler.GetUser)
    admin.PUT("/users/:id", adminUserHandler.UpdateUser)
    admin.DELETE("/users/:id", adminUserHandler.DeleteUser)
    admin.PUT("/users/:id/roles", adminUserHandler.AssignRoles)

    // 角色管理
    admin.POST("/roles", roleHandler.CreateRole)
    admin.GET("/roles", roleHandler.ListRoles)
    // ...

    // 审计日志
    admin.GET("/audit-logs", auditLogHandler.ListLogs)
}

// 3. 用户路由（仅需认证）
userGroup := api.Group("/user")
userGroup.Use(middleware.JWTAuth(jwtManager))
{
    userGroup.GET("/me", userProfileHandler.GetProfile)
    userGroup.PUT("/me", userProfileHandler.UpdateProfile)
    userGroup.PUT("/me/password", userProfileHandler.ChangePassword)
    userGroup.DELETE("/me", userProfileHandler.DeleteAccount)
}
```

### 管理员接口保护示例

所有 `/api/admin/*` 接口都受到多层保护：

1. **JWTAuth**: 必须提供有效的 JWT Token
2. **AuditMiddleware**: 自动记录所有操作
3. **RequireRole("admin")**: 必须拥有 `admin` 角色

**请求示例**:

```bash
curl -X GET http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**权限检查流程**:

```
请求 → JWTAuth → 验证 Token → 提取 roles: ["admin"]
     → AuditMiddleware → 记录日志
     → RequireRole("admin") → 检查 "admin" 是否在 roles 中 → ✓ 通过
     → Handler 执行
```

### 用户接口保护示例

所有 `/api/user/*` 接口仅需要 JWT 认证：

**请求示例**:

```bash
curl -X GET http://localhost:8080/api/user/me \
  -H "Authorization: Bearer <user_token>"
```

**权限检查流程**:

```
请求 → JWTAuth → 验证 Token → 提取 user_id
     → Handler 执行（使用 user_id 查询该用户的数据）
```

## 完整的请求流程

### 流程图

```mermaid
graph TD
    A[用户发起请求] --> B{请求路径}
    B -->|/api/auth/*| C[公开接口<br/>无认证]
    B -->|/api/admin/*| D[JWTAuth 中间件]
    B -->|/api/user/*| D

    D --> E{Token 有效?}
    E -->|否| F[401 Unauthorized]
    E -->|是| G[解析 Claims<br/>注入 Context]

    G --> H{路径是 /admin/*?}
    H -->|是| I[AuditMiddleware<br/>记录日志]
    H -->|否| L[执行 Handler]

    I --> J[RequireRole admin]
    J --> K{有 admin 角色?}
    K -->|否| M[403 Forbidden]
    K -->|是| L

    C --> L
    L --> N[返回响应]
```

### 成功案例：管理员创建用户

**请求**:

```http
POST /api/admin/users HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "password123",
  "full_name": "New User"
}
```

**流程**:

1. **JWTAuth 中间件**:

   ```go
   claims := ValidateToken(token)
   // claims = {
   //   UserID: 1,
   //   Username: "admin",
   //   Roles: ["admin"],
   //   Permissions: ["users:create", "users:read", ...]
   // }

   c.Set("user_id", 1)
   c.Set("username", "admin")
   c.Set("roles", ["admin"])
   c.Set("permissions", ["users:create", ...])
   ```

2. **AuditMiddleware**:

   ```go
   // 异步记录审计日志
   log := AuditLog{
       UserID: 1,
       Username: "admin",
       Action: "POST /api/admin/users",
       Resource: "users",
       IPAddress: "192.168.1.100",
   }
   repo.Create(log)
   ```

3. **RequireRole("admin")**:

   ```go
   roles := c.Get("roles")  // ["admin"]
   // 检查 "admin" 是否在列表中
   hasRole := true  // ✓ 通过
   c.Next()
   ```

4. **Handler 执行**:
   ```go
   adminUserHandler.CreateUser(c)
   // 创建用户逻辑...
   // 密码加密、保存到数据库等
   ```

**响应**:

```json
{
  "message": "user created successfully",
  "user": {
    "id": 10,
    "username": "newuser",
    "email": "newuser@example.com",
    "full_name": "New User",
    "status": "active"
  }
}
```

### 失败案例：普通用户访问管理员接口

**请求**:

```http
GET /api/admin/users HTTP/1.1
Authorization: Bearer <普通用户的 token>
```

**流程**:

1. **JWTAuth 中间件**:

   ```go
   claims := {
       UserID: 5,
       Username: "testuser",
       Roles: ["user"],  // 注意：不是 admin
       Permissions: ["users:read"]
   }

   c.Set("roles", ["user"])
   ```

2. **AuditMiddleware**: 跳过（GET 请求）

3. **RequireRole("admin")**:

   ```go
   roles := ["user"]
   // 检查 "admin" 是否在列表中
   hasRole := false  // ✗ 不通过

   c.JSON(403, gin.H{"error": "forbidden: insufficient permissions"})
   c.Abort()  // 中断请求链，不执行 Handler
   ```

**响应**:

```json
{
  "error": "forbidden: insufficient permissions"
}
```

### 所有权检查案例：用户修改自己的资料

如果使用 `RequireOwnership` 中间件：

**请求**:

```http
PUT /api/users/5 HTTP/1.1
Authorization: Bearer <用户 ID 为 5 的 token>
Content-Type: application/json

{
  "full_name": "Updated Name"
}
```

**流程**:

1. **JWTAuth**: 注入 `user_id = 5`

2. **RequireOwnership**:

   ```go
   userID := c.Get("user_id")      // 5
   resourceID := c.Param("id")     // "5" (URL 参数)

   // 比较 userID 和 resourceID
   if uint(resourceID) == userID {
       c.Next()  // ✓ 是同一用户，通过
   }
   ```

3. **Handler 执行**: 更新用户 ID 为 5 的资料

**如果用户 ID 不匹配**:

```http
PUT /api/users/10 HTTP/1.1
Authorization: Bearer <用户 ID 为 5 的 token>
```

```go
userID = 5
resourceID = 10
// 5 != 10
c.JSON(403, gin.H{"error": "forbidden: can only access own resources"})
c.Abort()
```

## 系统初始化

### 数据库迁移

启动服务时自动迁移：

```bash
go run . serve --auto-migrate
```

或使用 CLI 命令：

```bash
go run . migrate up
```

迁移会创建以下表：

- `users`
- `roles`
- `permissions`
- `user_roles` (关联表)
- `role_permissions` (关联表)
- `audit_logs`

### RBAC 种子数据

运行种子数据初始化：

```bash
go run . seed
```

或在代码中使用：

```go
seeder := database.NewSeeder(db)
seeder.SeedAll(context.Background())
```

**初始化内容**:

1. **8 个权限**:
   - `users:read`, `users:write`, `users:delete`
   - `roles:read`, `roles:write`, `roles:delete`
   - `permissions:read`
   - `audit_logs:read`

2. **2 个角色**:
   - `admin`: 拥有所有权限
   - `user`: 仅 `users:read` 权限

3. **1 个管理员用户**:
   - 用户名: `admin`
   - 密码: `admin123`
   - 邮箱: `admin@example.com`
   - 角色: `admin`

**安全警告**: 首次登录后请立即修改默认密码！

```bash
curl -X PUT http://localhost:8080/api/user/me/password \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "admin123",
    "new_password": "your_secure_password"
  }'
```

## 使用指南

### 1. 登录获取 Token

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "admin",
    "password": "admin123"
  }'
```

**响应**:

```json
{
  "message": "login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com"
    }
  }
}
```

### 2. 为用户分配角色

**查看所有角色**:

```bash
curl -X GET http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer <admin_token>"
```

**响应示例**:

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "admin",
      "display_name": "管理员",
      "description": "系统管理员角色"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_pages": 1,
    "has_more": false
  }
}
```

**为用户分配角色**:

```bash
curl -X PUT http://localhost:8080/api/admin/users/5/roles \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "role_ids": [1, 2]
  }'
```

### 3. 创建自定义角色

```bash
curl -X POST http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "editor",
    "display_name": "编辑",
    "description": "内容编辑权限"
  }'
```

### 4. 为角色分配权限

**查看所有权限**:

```bash
curl -X GET http://localhost:8080/api/admin/permissions \
  -H "Authorization: Bearer <admin_token>"
```

**响应示例**:

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "admin.users.read",
      "display_name": "读取用户",
      "resource": "users",
      "action": "read"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 50,
    "total": 20,
    "total_pages": 1,
    "has_more": false
  }
}
```

**为角色设置权限**:

```bash
curl -X PUT http://localhost:8080/api/admin/roles/3/permissions \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "permission_ids": [1, 2, 3]
  }'
```

### 5. 查看审计日志

```bash
# 查看所有审计日志
curl -X GET http://localhost:8080/api/admin/audit-logs \
  -H "Authorization: Bearer <admin_token>"

# 过滤：查看特定用户的操作
curl -X GET "http://localhost:8080/api/admin/audit-logs?user_id=1&page=1&limit=20" \
  -H "Authorization: Bearer <admin_token>"

# 过滤：查看特定资源的操作
curl -X GET "http://localhost:8080/api/admin/audit-logs?resource=users" \
  -H "Authorization: Bearer <admin_token>"
```

**响应示例**:

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 12,
      "user_id": 1,
      "action": "update",
      "resource": "users",
      "status": "success",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 42,
    "total_pages": 3,
    "has_more": true
  }
}
```

## 最佳实践

### 安全建议

1. **立即修改默认密码**
   - 默认管理员密码 `admin123` 仅用于初始化
   - 首次登录后立即修改

2. **定期刷新 Token**
   - Access Token 默认有效期 1 小时
   - 使用 Refresh Token 定期更新

   ```bash
   curl -X POST http://localhost:8080/api/auth/refresh \
     -H "Content-Type: application/json" \
     -d '{"refresh_token": "<refresh_token>"}'
   ```

3. **使用 HTTPS**
   - 生产环境必须使用 HTTPS
   - 防止 Token 在传输过程中被窃取

4. **最小权限原则**
   - 只分配必要的权限
   - 避免给普通用户分配 `admin` 角色

5. **审计日志监控**
   - 定期检查审计日志
   - 关注异常操作（如大量删除操作）

### 权限设计原则

1. **使用角色而非直接分配权限**

   ```
   ✓ 推荐: User → Role → Permissions
   ✗ 避免: User → Permissions (难以管理)
   ```

2. **粗粒度与细粒度结合**
   - 管理员接口使用 `RequireRole("admin")`（粗粒度，简单）
   - 特殊接口使用 `RequirePermission("resource:action")`（细粒度，精确）

3. **合理命名权限代码**
   - 使用 `resource:action` 格式
   - 资源名用复数（`users`, `roles`）
   - 操作名用动词（`create`, `read`, `update`, `delete`）

### Token 管理建议

1. **Token 失效时间设置**

   ```go
   // configs/config.example.yaml
   jwt:
     access_token_expiry: 1h      # Access Token 有效期
     refresh_token_expiry: 168h   # Refresh Token 有效期 (7天)
   ```

2. **权限变更的处理**
   - **当前行为**: 权限信息存储在 Token 中，变更后需重新登录
   - **改进方案**:
     - 结合 Redis 实现实时权限检查
     - 或强制 Token 失效机制

3. **Token 体积优化**
   - 如果用户权限过多，Token 会很大
   - 考虑仅存储 Roles，权限在后端实时查询
   - 或使用 Redis 缓存权限信息

### 审计日志使用

1. **定期清理**

   ```go
   // 删除 90 天前的日志
   auditLogRepo.DeleteOlderThan(ctx, 90)
   ```

2. **重点监控**
   - 用户创建/删除操作
   - 角色权限变更
   - 大量失败的登录尝试

3. **日志查询优化**
   - 按时间范围查询
   - 按用户 ID 查询
   - 按资源类型查询

## 常见问题

### Q1: 权限变更后何时生效？

**A**: 当前实现中，权限信息存储在 JWT Token 中，因此：

- **角色/权限变更后**: 用户需要重新登录或刷新 Token 才能获取新权限
- **Token 有效期内**: 使用的是 Token 生成时的权限信息

**改进方案**:

1. 缩短 Access Token 有效期（如 15 分钟）
2. 使用 Redis 存储用户权限，每次请求时实时检查
3. 实现 Token 黑名单机制，权限变更时强制 Token 失效

### Q2: 如何实现更细粒度的权限控制？

**A**: 使用 `RequirePermission` 中间件：

```go
// 不同的操作使用不同的权限
router.POST("/articles",
    middleware.RequirePermission("articles:create"),
    handler.CreateArticle,
)
router.PUT("/articles/:id",
    middleware.RequirePermission("articles:update"),
    handler.UpdateArticle,
)
router.DELETE("/articles/:id",
    middleware.RequirePermission("articles:delete"),
    handler.DeleteArticle,
)
```

也可以组合使用：

```go
router.PUT("/articles/:id",
    middleware.RequireAnyPermission("articles:update", "articles:manage"),
    handler.UpdateArticle,
)
```

### Q3: 如何让用户可以修改自己的资料但不能修改他人？

**A**: 使用 `RequireAdminOrOwnership` 中间件：

```go
router.PUT("/users/:id",
    middleware.JWTAuth(jwtManager),
    middleware.RequireAdminOrOwnership(),  // 管理员或本人
    handler.UpdateUser,
)
```

这样：

- 管理员可以修改任何用户的资料
- 普通用户只能修改自己的资料（`user_id == :id`）

### Q4: 如何处理跨域请求？

**A**: 使用 CORS 中间件（已内置）：

```go
// internal/adapters/http/middleware/cors.go
func CORS() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}
```

### Q5: Token 被盗用怎么办？

**A**: 安全措施：

1. **使用 HTTPS**: 防止中间人攻击
2. **短有效期**: Access Token 仅 1 小时有效
3. **IP 绑定**: 可在 Claims 中加入 IP 地址，验证时检查
4. **设备绑定**: 可加入设备指纹
5. **实时撤销**: 实现 Token 黑名单机制

```go
// 增强版 Claims（可选）
type Claims struct {
    UserID      uint
    Username    string
    Roles       []string
    Permissions []string
    IPAddress   string  // 绑定 IP
    DeviceID    string  // 设备标识
    jwt.RegisteredClaims
}
```

### Q6: 如何添加新的权限？

**A**: 两种方式：

1. **通过 API** (推荐):

```bash
curl -X POST http://localhost:8080/api/admin/permissions \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "resource": "articles",
    "action": "publish",
    "code": "articles:publish",
    "description": "发布文章"
  }'
```

2. **修改种子数据** (`internal/infrastructure/database/seeds/rbac_seeder.go`):

```go
permissions := []role.Permission{
    // ... 现有权限
    {Resource: "articles", Action: "publish", Code: "articles:publish", Description: "发布文章"},
}
```

然后重新运行 seed：

```bash
go run . seed
```

## 技术架构总结

### 优势

- ✅ **基于 JWT 的无状态认证**: 不依赖服务器会话，易于水平扩展
- ✅ **完整的 RBAC 模型**: 支持多角色、多权限
- ✅ **灵活的中间件设计**: 可组合使用，易于扩展
- ✅ **自动审计日志**: 无侵入式记录所有管理操作
- ✅ **代码质量高**: 通过 golangci-lint 检查，0 issues

### 限制

- ⚠️ **权限变更延迟**: Token 有效期内使用旧权限（可通过缩短有效期或使用 Redis 改进）
- ⚠️ **Token 体积**: 权限多时 Token 会较大（可优化为仅存储角色）
- ⚠️ **无实时撤销**: 暂未实现 Token 黑名单机制

### 扩展方向

1. **Redis 缓存**: 缓存用户权限，支持实时变更
2. **Token 黑名单**: 支持强制 Token 失效
3. **权限通配符**: 支持 `users:*` 表示所有用户权限
4. **数据权限**: 不仅控制操作，还控制数据范围（如只能查看自己部门的数据）
5. **前端权限**: 根据权限动态显示菜单和按钮

## 参考资料

- 项目源码: `internal/adapters/http/middleware/rbac.go`
- JWT 实现: `internal/infrastructure/auth/jwt.go`
- 路由配置: `internal/adapters/http/router.go`
- 数据模型: `internal/domain/role/model.go`, `internal/domain/user/model.go`
- 种子数据: `internal/infrastructure/database/seeds/rbac_seeder.go`
