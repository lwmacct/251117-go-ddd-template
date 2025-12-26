# 身份认证与权限

面向后端开发者的身份体系功能手册，聚焦认证、授权和 Token 管理。

<!--TOC-->

## Table of Contents

- [认证机制](#认证机制) `:29+46`
  - [JWT Token 流程](#jwt-token-流程) `:31+12`
  - [功能特性](#功能特性) `:43+8`
  - [架构设计](#架构设计) `:51+12`
  - [API 端点](#api-端点) `:63+12`
- [RBAC 权限系统](#rbac-权限系统) `:75+45`
  - [三段式格式](#三段式格式) `:79+14`
  - [通配符匹配](#通配符匹配) `:93+6`
  - [中间件](#中间件) `:99+10`
  - [路由保护](#路由保护) `:109+4`
  - [最佳实践](#最佳实践) `:113+7`
- [Personal Access Token (PAT)](#personal-access-token-pat) `:120+38`
  - [PAT vs JWT](#pat-vs-jwt) `:124+10`
  - [Token 格式](#token-格式) `:134+9`
  - [API 端点](#api-端点-1) `:143+8`
  - [最佳实践](#最佳实践-1) `:151+7`
- [安全配置](#安全配置) `:158+18`

<!--TOC-->

## 认证机制

### JWT Token 流程

**登录流程**: 登录 → 查询用户角色 → 聚合权限 → 生成 JWT

Token Claims 包含:

| 字段          | 说明         |
| ------------- | ------------ |
| `sub`         | 用户 ID      |
| `roles`       | 角色名称列表 |
| `permissions` | 权限列表     |

### 功能特性

- 用户注册（用户名/邮箱唯一性验证）
- 用户登录（支持用户名或邮箱）
- Token 刷新（Access 15分钟，Refresh 7天）
- bcrypt 密码加密
- 用户状态检查（仅 active 可登录）

### 架构设计

```
internal/
├── domain/auth/           # 认证领域服务接口
├── application/auth/      # Register/Login/Refresh Use Case
├── infrastructure/auth/   # JWT 管理器、服务实现
├── adapters/http/
│   ├── handler/auth.go    # HTTP Handler
│   └── middleware/jwt.go  # JWT 鉴权中间件
```

### API 端点

**公开端点**:

| 方法 | 路径                 | 说明         |
| ---- | -------------------- | ------------ |
| POST | `/api/auth/register` | 注册新用户   |
| POST | `/api/auth/login`    | 用户登录     |
| POST | `/api/auth/refresh`  | 刷新访问令牌 |

详细 API 请参阅 Swagger UI (`/swagger/index.html`)

## RBAC 权限系统

基于角色的访问控制 (Role-Based Access Control)。

### 三段式格式

```
{resource}:{action}:{scope}
```

| 段       | 说明     | 示例                                 |
| -------- | -------- | ------------------------------------ |
| resource | 资源类型 | `user`, `role`, `menu`               |
| action   | 操作类型 | `create`, `read`, `update`, `delete` |
| scope    | 作用范围 | `*` (全部), `own` (自己)             |

**示例**: `user:read:*` (读取所有用户), `user:update:own` (更新自己)

### 通配符匹配

- `*:*:*` - 超级管理员权限
- `user:*:*` - 用户模块全部权限
- `*:read:*` - 所有模块读取权限

### 中间件

执行顺序: `Auth → RBAC → Handler`

| 中间件 | 位置                                         | 职责       |
| ------ | -------------------------------------------- | ---------- |
| Auth   | `internal/adapters/http/middleware/auth.go`  | 验证 Token |
| RBAC   | `internal/adapters/http/middleware/rbac.go`  | 检查权限   |
| Audit  | `internal/adapters/http/middleware/audit.go` | 记录日志   |

### 路由保护

**路由配置**: `internal/adapters/http/router.go`

### 最佳实践

- 遵循最小权限原则
- 使用角色而非直接分配权限
- 敏感操作启用审计日志
- Token 有效期不超过 24 小时

## Personal Access Token (PAT)

长期凭证，作为 JWT Token 的替代方案，用于 API 集成、CLI 和自动化脚本。

### PAT vs JWT

| 特性   | JWT Token     | Personal Access Token   |
| ------ | ------------- | ----------------------- |
| 用途   | Web/移动应用  | API 集成、CLI、脚本     |
| 有效期 | 短期（1小时） | 长期（7/30/90天或永久） |
| 刷新   | Refresh Token | 无需刷新                |
| 权限   | 用户全部权限  | 用户权限的子集          |
| 删除   | 不支持        | 支持即时删除            |

### Token 格式

**格式**: `pat_<5位前缀>_<32位随机字符>`

**安全特性**:

- 完整 token 仅在创建时显示一次
- 数据库存储 SHA-256 哈希值

### API 端点

| 方法   | 路径                   | 说明            |
| ------ | ---------------------- | --------------- |
| POST   | `/api/user/tokens`     | 创建 Token      |
| GET    | `/api/user/tokens`     | 查看 Token 列表 |
| DELETE | `/api/user/tokens/:id` | 删除 Token      |

### 最佳实践

- 只授予完成任务所需的最小权限
- 测试环境 7 天，生产环境 90 天
- 使用环境变量或密钥管理服务存储
- 每 90 天轮换生产环境 Token

## 安全配置

**JWT 配置**:

| 配置项               | 环境变量                       | 说明                     |
| -------------------- | ------------------------------ | ------------------------ |
| secret               | `APP_JWT_SECRET`               | 签名密钥（至少 32 字节） |
| access_token_expiry  | `APP_JWT_ACCESS_TOKEN_EXPIRY`  | 访问令牌有效期           |
| refresh_token_expiry | `APP_JWT_REFRESH_TOKEN_EXPIRY` | 刷新令牌有效期           |

**安全特性**:

| 特性       | 说明                                   |
| ---------- | -------------------------------------- |
| 密码加密   | bcrypt (cost=10)                       |
| Token 签名 | HMAC-SHA256                            |
| 唯一性约束 | 用户名、邮箱数据库层面强制唯一         |
| 错误处理   | 登录失败返回通用 "invalid credentials" |
