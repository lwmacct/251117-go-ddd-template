# API 文档概述

Go DDD Template 提供了完整的 RESTful API，支持 JWT 和 PAT 两种认证方式。

<!--TOC-->

## Table of Contents

- [API 基础信息](#api-基础信息) `:38+7`
- [Swagger 文档](#swagger-文档) `:45+7`
- [认证方式](#认证方式) `:52+35`
  - [JWT Token](#jwt-token) `:54+16`
  - [Personal Access Token (PAT)](#personal-access-token-pat) `:70+17`
- [响应格式](#响应格式) `:87+45`
  - [成功响应](#成功响应) `:89+13`
  - [错误响应](#错误响应) `:102+14`
  - [分页响应](#分页响应) `:116+16`
- [HTTP 状态码](#http-状态码) `:132+16`
- [API 模块](#api-模块) `:148+39`
  - [认证管理 (/api/auth)](#认证管理-apiauth) `:150+7`
  - [用户管理 (/api/users)](#用户管理-apiusers) `:157+6`
  - [角色管理 (/api/roles)](#角色管理-apiroles) `:163+6`
  - [菜单管理 (/api/menus)](#菜单管理-apimenus) `:169+6`
  - [PAT 管理 (/api/pat)](#pat-管理-apipat) `:175+6`
  - [审计日志 (/api/audit-logs)](#审计日志-apiaudit-logs) `:181+6`
- [请求示例](#请求示例) `:187+55`
  - [基础认证流程](#基础认证流程) `:189+27`
  - [创建资源](#创建资源) `:216+14`
  - [查询资源](#查询资源) `:230+12`
- [公共参数](#公共参数) `:242+19`
  - [分页参数](#分页参数) `:244+8`
  - [过滤参数](#过滤参数) `:252+9`
- [限流策略](#限流策略) `:261+12`
- [最佳实践](#最佳实践) `:273+23`

<!--TOC-->

## API 基础信息

- **基础 URL**: `http://localhost:8080/api`
- **API 版本**: v1
- **响应格式**: JSON
- **字符编码**: UTF-8

## Swagger 文档

运行应用后，可通过以下地址访问交互式 API 文档：

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **OpenAPI JSON**: http://localhost:8080/swagger/doc.json

## 认证方式

### JWT Token

适用于 Web 应用和移动应用的短期认证：

```bash
# 登录获取 Token
POST /api/auth/login
{
  "login": "username",
  "password": "password"
}

# 使用 Token
Authorization: Bearer <access_token>
```

### Personal Access Token (PAT)

适用于 API 集成和 CLI 工具的长期认证：

```bash
# 创建 PAT
POST /api/pat
{
  "name": "CI/CD Token",
  "permissions": ["api:users:read"],
  "expires_at": "2025-12-31"
}

# 使用 PAT
Authorization: Bearer pat_xxxxx
```

## 响应格式

### 成功响应

```json
{
  "code": 200,
  "status": "success",
  "message": "操作成功",
  "data": {
    // 响应数据
  }
}
```

### 错误响应

```json
{
  "code": 400,
  "status": "error",
  "message": "请求参数错误",
  "details": {
    "field": "email",
    "reason": "invalid format"
  }
}
```

### 分页响应

```json
{
  "code": 200,
  "status": "success",
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5
  }
}
```

## HTTP 状态码

| 状态码 | 说明                  | 使用场景               |
| ------ | --------------------- | ---------------------- |
| 200    | OK                    | 成功的 GET、PUT 请求   |
| 201    | Created               | 成功的 POST 创建请求   |
| 204    | No Content            | 成功的 DELETE 请求     |
| 400    | Bad Request           | 请求参数错误           |
| 401    | Unauthorized          | 未认证或认证失败       |
| 403    | Forbidden             | 无权限访问             |
| 404    | Not Found             | 资源不存在             |
| 409    | Conflict              | 资源冲突（如重复创建） |
| 422    | Unprocessable Entity  | 请求格式正确但语义错误 |
| 429    | Too Many Requests     | 请求频率超限           |
| 500    | Internal Server Error | 服务器内部错误         |

## API 模块

### 认证管理 (/api/auth)

- 用户注册、登录、登出
- Token 刷新
- 密码重置
- 个人信息管理

### 用户管理 (/api/users)

- 用户 CRUD 操作
- 用户状态管理
- 批量操作

### 角色管理 (/api/roles)

- 角色 CRUD 操作
- 权限分配
- 用户角色关联

### 菜单管理 (/api/menus)

- 菜单树结构管理
- 菜单权限配置
- 动态菜单生成

### PAT 管理 (/api/pat)

- PAT 创建和撤销
- 权限范围设置
- 使用记录查询

### 审计日志 (/api/audit-logs)

- 操作日志查询
- 日志导出
- 统计分析

## 请求示例

### 基础认证流程

```bash
# 1. 用户登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "admin",
    "password": "password123"
  }'

# 响应
{
  "code": 200,
  "status": "success",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "expires_in": 3600
  }
}

# 2. 使用 Token 访问受保护资源
curl http://localhost:8080/api/users \
  -H "Authorization: Bearer eyJhbGc..."
```

### 创建资源

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "new@example.com",
    "password": "secure123",
    "full_name": "New User"
  }'
```

### 查询资源

```bash
# 分页查询
curl "http://localhost:8080/api/users?page=1&page_size=10&sort=created_at" \
  -H "Authorization: Bearer <token>"

# 条件过滤
curl "http://localhost:8080/api/users?status=active&role=admin" \
  -H "Authorization: Bearer <token>"
```

## 公共参数

### 分页参数

| 参数      | 类型   | 默认值      | 说明                  |
| --------- | ------ | ----------- | --------------------- |
| page      | int    | 1           | 页码                  |
| page_size | int    | 20          | 每页数量              |
| sort      | string | -created_at | 排序字段（-表示降序） |

### 过滤参数

| 参数           | 类型   | 说明         |
| -------------- | ------ | ------------ |
| q              | string | 全文搜索     |
| status         | string | 状态过滤     |
| created_after  | string | 创建时间起始 |
| created_before | string | 创建时间结束 |

## 限流策略

- **全局限流**: 1000 请求/分钟
- **登录限流**: 5 次/分钟
- **API 限流**: 100 请求/分钟（per token）

超过限流后返回 429 状态码，响应头包含：

- `X-RateLimit-Limit`: 限流阈值
- `X-RateLimit-Remaining`: 剩余配额
- `X-RateLimit-Reset`: 重置时间

## 最佳实践

1. **使用正确的 HTTP 方法**
   - GET: 查询资源
   - POST: 创建资源
   - PUT: 完整更新资源
   - PATCH: 部分更新资源
   - DELETE: 删除资源

2. **Token 管理**
   - 定期刷新 Access Token
   - 安全存储 Refresh Token
   - PAT 仅用于自动化场景

3. **错误处理**
   - 检查响应状态码
   - 解析错误详情
   - 实现重试机制

4. **性能优化**
   - 使用分页避免大量数据传输
   - 合理使用查询参数
   - 启用响应压缩
