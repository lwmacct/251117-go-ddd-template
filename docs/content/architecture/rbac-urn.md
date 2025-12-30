# URN 风格 RBAC 权限系统

基于统一资源名称 (URN) 格式的细粒度权限控制系统，统一 Operation 和 Resource 表达方式。

<!--TOC-->

## Table of Contents

- [设计理念](#设计理念) `:34+12`
- [URN 格式](#urn-格式) `:46+31`
  - [三段式结构](#三段式结构) `:48+22`
  - [分隔符规则](#分隔符规则) `:70+7`
- [Scope 体系](#scope-体系) `:77+25`
  - [预定义 Scope](#预定义-scope) `:79+8`
  - [层级 Scope](#层级-scope) `:87+15`
- [通配符匹配](#通配符匹配) `:102+38`
  - [基本通配符](#基本通配符) `:104+11`
  - [Scope 层级通配符](#scope-层级通配符) `:115+9`
  - [匹配示例](#匹配示例) `:124+16`
- [运行时变量](#运行时变量) `:140+30`
  - [支持的变量](#支持的变量) `:142+7`
  - [替换时机](#替换时机) `:149+21`
- [权限配置](#权限配置) `:170+44`
  - [角色权限表](#角色权限表) `:172+9`
  - [预置角色示例](#预置角色示例) `:181+33`
- [代码实现](#代码实现) `:214+84`
  - [核心文件](#核心文件) `:216+12`
  - [Operation 常量](#operation-常量) `:228+31`
  - [中间件集成](#中间件集成) `:259+39`
- [与旧格式对比](#与旧格式对比) `:298+19`

<!--TOC-->

## 设计理念

**统一 Operation 和 Resource 格式**，采用 URN 风格，核心优势：

| 特性           | 说明                                                    |
| -------------- | ------------------------------------------------------- |
| **统一格式**   | Operation 和 Resource 共用 `scope:type:identifier` 格式 |
| **分隔符简洁** | 只用 `:` 和 `.`，认知负担低                             |
| **多级 Scope** | 支持 `sys.admin`、`org.acme.team.dev` 等层级            |
| **运行时变量** | `@me`、`@org` 支持动态资源匹配                          |
| **通配符灵活** | `scope.*` 支持子域匹配                                  |

## URN 格式

### 三段式结构

```
{scope}:{type}:{identifier}
   │       │        │
   │       │        └─ 动作(create) 或 资源ID(123/@me/*)
   │       └─ 模块/资源类型 (users, profile, tokens)
   └─ 作用域 (sys, self, public, org.acme)
```

**Operation 示例**：

- `sys:users:create` - 系统级用户创建操作
- `self:profile:update` - 用户更新自己的资料
- `public:auth:login` - 公开的登录操作

**Resource 示例**：

- `sys:user:123` - 系统用户 ID 123
- `self:user:@me` - 当前登录用户
- `org.acme:project:*` - acme 组织的所有项目

### 分隔符规则

| 分隔符 | 用途            | 示例                             |
| ------ | --------------- | -------------------------------- |
| `:`    | 分隔主要部分    | `sys:users:create`               |
| `.`    | 分隔 scope 层级 | `org.acme.team.dev:tasks:create` |

## Scope 体系

### 预定义 Scope

| Scope    | 说明               | 典型 Operation        | 典型 Resource   |
| -------- | ------------------ | --------------------- | --------------- |
| `sys`    | 系统级（平台管理） | `sys:users:create`    | `sys:user:*`    |
| `self`   | 用户自身           | `self:profile:update` | `self:user:@me` |
| `public` | 公开（无需权限）   | `public:auth:login`   | -               |

### 层级 Scope

支持多级 scope 用于多租户场景：

```
org.{org_id}                    # 组织级
org.{org_id}.team.{team_id}     # 团队级
sys.admin                       # 系统管理子域
```

**示例**：

- `org.acme:users:list` - acme 组织的用户列表操作
- `org.acme.team.dev:tasks:create` - acme 组织 dev 团队的任务创建

## 通配符匹配

### 基本通配符

`*` 匹配当前段任意值：

| 模式          | 说明                   |
| ------------- | ---------------------- |
| `*:*:*`       | 匹配所有（超级管理员） |
| `sys:*:*`     | sys 域所有操作         |
| `sys:users:*` | sys 域用户模块所有操作 |
| `*:*:read`    | 所有域所有模块的读操作 |

### Scope 层级通配符

`.*` 后缀匹配 scope 及其所有子域：

| 模式             | 匹配                                                |
| ---------------- | --------------------------------------------------- |
| `sys.*:*:*`      | `sys`、`sys.admin`、`sys.readonly`                  |
| `org.acme.*:*:*` | `org.acme`、`org.acme.team.dev`、`org.acme.team.qa` |

### 匹配示例

```go
// Operation 匹配
"sys:users:create".Match("*:*:*")           // true - 超级通配
"sys:users:create".Match("sys:*:*")         // true - 域通配
"sys:users:create".Match("sys:users:*")     // true - 模块通配
"sys:users:create".Match("sys:users:create") // true - 精确匹配
"sys.admin:config:update".Match("sys.*:*:*") // true - 子域通配

// Resource 匹配
"sys:user:123".Match("*:*:*")               // true
"sys:user:123".Match("sys:user:*")          // true
"self:user:456".Match("self:user:@me")      // 需先替换 @me
```

## 运行时变量

### 支持的变量

| 变量   | 说明             | 替换结果示例                      |
| ------ | ---------------- | --------------------------------- |
| `@me`  | 当前用户 ID      | `self:user:@me` → `self:user:123` |
| `@org` | 当前用户所属组织 | `org.@org:*:*` → `org.acme:*:*`   |

### 替换时机

在 RBAC 中间件检查权限时，自动替换资源模式中的变量：

```go
// internal/domain/urn/resolver.go
type Resolver struct {
    UserID uint
    OrgID  string
}

func (r *Resolver) Resolve(u URN) URN {
    s := string(u)
    s = strings.ReplaceAll(s, "@me", strconv.FormatUint(uint64(r.UserID), 10))
    if r.OrgID != "" {
        s = strings.ReplaceAll(s, "@org", r.OrgID)
    }
    return URN(s)
}
```

## 权限配置

### 角色权限表

数据库中通过 `role_permissions` 表存储：

| 字段                | 说明     | 示例          |
| ------------------- | -------- | ------------- |
| `operation_pattern` | 操作模式 | `sys:users:*` |
| `resource_pattern`  | 资源模式 | `*:*:*`       |

### 预置角色示例

**超级管理员 (admin)**：

```
operation_pattern: *:*:*
resource_pattern:  *:*:*
```

**普通用户 (user)**：

```
# self 域操作，仅限自身资源
operation_pattern: self:*:*
resource_pattern:  self:user:@me

# 公开操作，所有资源
operation_pattern: public:*:*
resource_pattern:  *:*:*
```

**系统管理员（无超级权限）**：

```
# sys 域所有操作
operation_pattern: sys:*:*
resource_pattern:  *:*:*

# self 域操作
operation_pattern: self:*:*
resource_pattern:  self:user:@me
```

## 代码实现

### 核心文件

| 文件                                        | 说明               |
| ------------------------------------------- | ------------------ |
| `internal/domain/urn/urn.go`                | URN 类型定义和解析 |
| `internal/domain/urn/matcher.go`            | URN 匹配算法       |
| `internal/domain/urn/resolver.go`           | 运行时变量替换     |
| `internal/domain/operation/operation.go`    | Operation 类型     |
| `internal/domain/operation/resource.go`     | Resource 类型      |
| `internal/domain/operation/constants.go`    | Operation 常量定义 |
| `internal/adapters/http/middleware/rbac.go` | RBAC 中间件        |

### Operation 常量

```go
// internal/domain/operation/constants.go

// 系统级操作 (sys:*)
const (
    SysUsersCreate Operation = "sys:users:create"
    SysUsersRead   Operation = "sys:users:read"
    SysUsersUpdate Operation = "sys:users:update"
    SysUsersDelete Operation = "sys:users:delete"
    SysRolesCreate Operation = "sys:roles:create"
    // ...
)

// 公开操作 (public:*)
const (
    PublicAuthLogin    Operation = "public:auth:login"
    PublicAuthRegister Operation = "public:auth:register"
    PublicAuthRefresh  Operation = "public:auth:refresh"
)

// 用户自身操作 (self:*)
const (
    SelfProfileGet     Operation = "self:profile:get"
    SelfProfileUpdate  Operation = "self:profile:update"
    SelfTokensCreate   Operation = "self:tokens:create"
    SelfSettingsSet    Operation = "self:settings:set"
)
```

### 中间件集成

```go
// internal/adapters/http/middleware/rbac.go

func RequireOperation(op operation.Operation) gin.HandlerFunc {
    return func(c *gin.Context) {
        permList := getPermissions(c)
        operationStr := string(op)

        // 确定要检查的资源
        resourcesToCheck := []string{"*:*:*"}
        if op.Scope() == "self" {
            // self: 域操作添加 self:user:@me 资源
            resolver := &urn.Resolver{UserID: getUserID(c)}
            resourcesToCheck = append(resourcesToCheck,
                "self:user:@me",
                resolver.ResolveString("self:user:@me"),
            )
        }

        // 检查权限
        for _, p := range permList {
            if operation.MatchOperation(p.OperationPattern, operationStr) {
                for _, res := range resourcesToCheck {
                    if operation.MatchResource(p.ResourcePattern, res) {
                        c.Next()
                        return
                    }
                }
            }
        }

        response.Forbidden(c, "Insufficient permissions")
        c.Abort()
    }
}
```

## 与旧格式对比

| 旧格式                | 新格式 (URN)          | 改进                 |
| --------------------- | --------------------- | -------------------- |
| `sys:users.create`    | `sys:users:create`    | 统一分隔符 `:`       |
| `user:profile.update` | `self:profile:update` | `self` 语义更清晰    |
| `user/123`            | `sys:user:123`        | Resource 带 scope    |
| `user/self`           | `self:user:@me`       | 明确的占位符         |
| `Public: true` 字段   | `public:` 前缀        | scope 自动推断公开性 |

**公开操作判断**：

```go
// 旧方式：需要在 registry 中标记 Public: true
// 新方式：通过 scope 自动判断
func (o Operation) IsPublic() bool {
    return o.Scope() == "public"
}
```
