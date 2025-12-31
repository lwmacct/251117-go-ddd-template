---
paths:
  - "internal/domain/org/**/*.go"
  - "internal/adapters/http/middleware/org_context.go"
---

# 多租户上下文系统

组织（Org）和团队（Team）的信息携带与权限隔离机制。

## 核心设计决策

**org/team 信息不存储在 JWT Token 中**，而是通过路由参数 + 中间件链动态注入。

| 特性         | 实现方式                          | 优势                   |
| ------------ | --------------------------------- | ---------------------- |
| 多租户支持   | org_id 从路由参数获取             | 用户可同时属于多个组织 |
| 权限即时生效 | 权限从缓存实时查询                | 权限变更无需重新登录   |
| 层级隔离     | OrgContext → TeamContext 链式验证 | team 必须属于当前 org  |

## 信息流转路径

```
请求 /api/org/:org_id/teams/:team_id/members
                │              │
Auth 中间件     │              │
├─ JWT: {user_id, username}    │
└─ Context: {user_id, permissions}
                │              │
OrgContext      ▼              │
├─ 提取路由参数: org_id        │
├─ 验证: user 是 org 成员      │
└─ Context: {org_id, org_role} │
                               │
TeamContext                    ▼
├─ 提取路由参数: team_id
├─ 验证: team 属于 org && user 是 team 成员
└─ Context: {team_id, team_role}
                │
RBAC 中间件     ▼
├─ 构建变量: {@me, @org, @team}
├─ 解析: org.@org → org.123
└─ 权限匹配检查
```

## JWT Token 内容

```go
type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    // ❌ 无 org_id - 用户可属多组织
}
```

## Context 变量注入

| 中间件      | 注入 Key      | 类型           | 来源                |
| ----------- | ------------- | -------------- | ------------------- |
| Auth        | `user_id`     | `uint`         | JWT Claims          |
| Auth        | `permissions` | `[]Permission` | Redis 缓存          |
| OrgContext  | `org_id`      | `uint`         | 路由参数 `:org_id`  |
| OrgContext  | `org_role`    | `string`       | 数据库查询          |
| TeamContext | `team_id`     | `uint`         | 路由参数 `:team_id` |
| TeamContext | `team_role`   | `string`       | 数据库查询          |

## 运行时变量

RBAC 中间件从 Context 构建运行时变量：

```go
func buildContextVars(c *gin.Context) map[string]string {
    vars := map[string]string{}

    // @me - 始终可用
    if userID, ok := c.Get("user_id"); ok {
        vars["@me"] = fmt.Sprint(userID)
    }

    // @org - OrgContext 后可用
    if orgID, ok := c.Get("org_id"); ok {
        vars["@org"] = fmt.Sprint(orgID)
    }

    // @team - TeamContext 后可用
    if teamID, ok := c.Get("team_id"); ok {
        vars["@team"] = fmt.Sprint(teamID)
    }

    return vars
}
```

## 中间件自动应用

路由器根据路径自动判断是否需要中间件：

```go
// 路径含 :org_id → 自动添加 OrgContext
// 路径含 :team_id → 自动添加 TeamContext（必须在 OrgContext 后）
```

## 权限资源解析

检查 `org:members:list` 时，生成候选资源列表：

```go
resources := []string{
    "*:*:*",              // 超级通配符
    "org.@org:*:*",       // 原始模式
    "org.123:*:*",        // 解析后
}

// 有 team 上下文时追加
resources = append(resources,
    "org.@org.team.@team:*:*",
    "org.123.team.456:*:*",
)
```

## 角色权限配置示例

**组织管理员**：

```
operation_pattern: org:*:*
resource_pattern:  org.@org:*:*
```

**团队负责人**：

```
operation_pattern: org:team-members:*
resource_pattern:  org.@org.team.@team:*:*
```

## 代码实现

| 文件                                               | 说明                  |
| -------------------------------------------------- | --------------------- |
| `internal/adapters/http/middleware/org_context.go` | Org/Team 上下文中间件 |
| `internal/adapters/http/middleware/rbac.go`        | 运行时变量构建        |
| `internal/domain/org/entity_member.go`             | 组织成员实体          |
| `internal/domain/org/entity_team.go`               | 团队实体              |
| `internal/infrastructure/auth/jwt.go`              | JWT Claims 定义       |
