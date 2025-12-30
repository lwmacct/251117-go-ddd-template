---
paths:
  - "internal/domain/permission/**/*.go"
  - "internal/adapters/http/middleware/rbac.go"
---

# URN 风格 RBAC 权限系统

基于统一资源名称 (URN) 格式的细粒度权限控制系统。

## URN 格式

### 三段式结构

```
{scope}:{type}:{identifier}
   │       │        │
   │       │        └─ 动作(create) 或 资源ID(123/@me/*)
   │       └─ 模块/资源类型 (users, profile, tokens)
   └─ 作用域 (sys, self, public, org.acme)
```

**分隔符**：`:` 分隔主要部分，`.` 分隔 scope 层级

### Operation 与 Resource

| 类型      | 用途     | 示例               |
| --------- | -------- | ------------------ |
| Operation | 操作标识 | `sys:users:create` |
| Resource  | 资源标识 | `sys:user:123`     |

## Scope 体系

| Scope    | 说明             | 示例 Operation        |
| -------- | ---------------- | --------------------- |
| `public` | 公开（无需认证） | `public:auth:login`   |
| `sys`    | 系统管理域       | `sys:users:create`    |
| `self`   | 用户自服务域     | `self:profile:update` |

**层级 Scope**（多租户场景）：

```
org.{org_id}                    # 组织级
org.{org_id}.team.{team_id}     # 团队级
sys.admin                       # 系统管理子域
```

## 通配符匹配

### 基本通配符

| 模式          | 说明                 |
| ------------- | -------------------- |
| `*:*:*`       | 匹配所有（超级管理） |
| `sys:*:*`     | sys 域所有操作       |
| `sys:users:*` | 用户模块所有操作     |

### Scope 层级通配符

`.*` 后缀匹配 scope 及其所有子域：

| 模式             | 匹配                                           |
| ---------------- | ---------------------------------------------- |
| `sys.*:*:*`      | `sys`、`sys.admin`、`sys.readonly`             |
| `org.acme.*:*:*` | `org.acme`、`org.acme.team.dev`、`org.acme.qa` |

## 运行时变量

| 变量   | 说明             | 示例                              |
| ------ | ---------------- | --------------------------------- |
| `@me`  | 当前用户 ID      | `self:user:@me` → `self:user:123` |
| `@org` | 当前用户所属组织 | `org.@org:*:*` → `org.acme:*:*`   |

```go
r := NewResolver(map[string]string{"@me": "123"})
r.ResolveString("self:user:@me")  // "self:user:123"
```

## 权限配置

### 角色权限表

数据库 `role_permissions` 表结构：

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
operation_pattern: self:*:*
resource_pattern:  self:user:@me
```

## 代码实现

### 核心文件

| 文件                                        | 说明           |
| ------------------------------------------- | -------------- |
| `internal/domain/permission/operation.go`   | Operation 类型 |
| `internal/domain/permission/resource.go`    | Resource 类型  |
| `internal/domain/permission/matcher.go`     | 匹配算法       |
| `internal/domain/permission/resolver.go`    | 变量替换       |
| `internal/domain/permission/constants.go`   | Operation 常量 |
| `internal/adapters/http/middleware/rbac.go` | RBAC 中间件    |

> 审计派生逻辑已迁移到 `internal/domain/audit/` 包。
