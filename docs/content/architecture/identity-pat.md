# Personal Access Token (PAT) 使用指南

本文档详细介绍 Personal Access Token (个人访问令牌) 的使用方法、最佳实践和安全建议。

<!--TOC-->

## Table of Contents

- [什么是 Personal Access Token?](#什么是-personal-access-token) `:48+21`
  - [PAT vs JWT 对比](#pat-vs-jwt-对比) `:58+11`
- [Token 格式](#token-格式) `:69+12`
- [创建 Personal Access Token](#创建-personal-access-token) `:81+93`
  - [前提条件](#前提条件) `:83+5`
  - [创建 Token](#创建-token) `:88+57`
  - [权限选择](#权限选择) `:145+29`
- [使用 Personal Access Token](#使用-personal-access-token) `:174+105`
  - [基本用法](#基本用法) `:176+11`
  - [在不同场景中使用](#在不同场景中使用) `:187+92`
- [管理 Personal Access Tokens](#管理-personal-access-tokens) `:279+74`
  - [查看 Token 列表](#查看-token-列表) `:281+40`
  - [查看单个 Token 详情](#查看单个-token-详情) `:321+9`
  - [删除 Token](#删除-token) `:330+23`
- [最佳实践](#最佳实践) `:353+93`
  - [1. 最小权限原则](#1-最小权限原则) `:355+16`
  - [2. 设置合理的有效期](#2-设置合理的有效期) `:371+19`
  - [3. 使用描述性名称](#3-使用描述性名称) `:390+15`
  - [4. 配置 IP 白名单](#4-配置-ip-白名单) `:405+10`
  - [5. 定期轮换 Token](#5-定期轮换-token) `:415+6`
  - [6. 安全存储](#6-安全存储) `:421+25`
- [安全建议](#安全建议) `:446+31`
  - [Token 泄露处理](#token-泄露处理) `:448+15`
  - [监控建议](#监控建议) `:463+8`
  - [IP 白名单注意事项](#ip-白名单注意事项) `:471+6`
- [Token 状态说明](#token-状态说明) `:477+10`
- [常见问题](#常见问题) `:487+71`
  - [Q1: Token 创建后忘记保存，如何找回？](#q1-token-创建后忘记保存如何找回) `:489+7`
  - [Q2: 如何增加 Token 的权限？](#q2-如何增加-token-的权限) `:496+8`
  - [Q3: Token 过期后会自动删除吗？](#q3-token-过期后会自动删除吗) `:504+4`
  - [Q4: 使用 PAT 时出现 403 错误？](#q4-使用-pat-时出现-403-错误) `:508+17`
  - [Q5: JWT 和 PAT 可以同时使用吗？](#q5-jwt-和-pat-可以同时使用吗) `:525+13`
  - [Q6: 如何批量管理 Token？](#q6-如何批量管理-token) `:538+8`
  - [Q7: PAT 支持通配符权限吗？](#q7-pat-支持通配符权限吗) `:546+12`
- [相关文档](#相关文档) `:558+6`
- [技术实现](#技术实现) `:564+32`

<!--TOC-->

## 什么是 Personal Access Token?

Personal Access Token (PAT) 是一种用于 API 认证的长期凭证，作为 JWT Token 的替代方案，特别适用于：

- **API 集成**：第三方应用调用您的 API
- **CLI 工具**：命令行工具自动化操作
- **自动化脚本**：定时任务、部署脚本
- **Webhook 回调**：接收和处理 Webhook 事件
- **测试环境**：API 测试和调试

### PAT vs JWT 对比

| 特性        | JWT Token     | Personal Access Token   |
| ----------- | ------------- | ----------------------- |
| **用途**    | Web/移动应用  | API 集成、CLI、脚本     |
| **有效期**  | 短期（1小时） | 长期（7/30/90天或永久） |
| **刷新**    | Refresh Token | 无需刷新                |
| **权限**    | 用户全部权限  | 用户权限的子集          |
| **删除**    | 不支持        | 支持即时删除            |
| **IP 限制** | 不支持        | 支持 IP 白名单          |

## Token 格式

**格式**: `pat_<5位前缀>_<32位随机字符>`

**示例**: `pat_2Kj9X_aB3cD4eF5gH6iJ7kL8mN9oP0qR1sT2uV3wX4yZ`

**安全特性**：

- 完整 token 仅在创建时显示一次
- 数据库存储 SHA-256 哈希值
- 前缀用于快速识别，不影响安全性

## 创建 Personal Access Token

### 前提条件

1. 已登录并获取 JWT Token
2. 拥有 `user:tokens:create` 权限

### 创建 Token

**API 端点**: `POST /api/user/tokens`

**请求示例**:

```bash
curl -X POST http://localhost:8080/api/user/tokens \
  -H "Authorization: Bearer <your_jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My API Integration",
    "permissions": [
      "user:profile:read",
      "user:profile:update",
      "api:cache:read"
    ],
    "expires_in": 90,
    "ip_whitelist": ["192.168.1.100", "10.0.0.50"],
    "description": "Token for production API server"
  }'
```

**请求参数**:

| 字段           | 类型     | 必填 | 说明                                 |
| -------------- | -------- | ---- | ------------------------------------ |
| `name`         | string   | ✓    | Token 名称（1-100 字符）             |
| `permissions`  | string[] | ✓    | 权限列表（必须是用户已有权限的子集） |
| `expires_in`   | int      | ✗    | 有效期天数（7/30/90 或 null=永久）   |
| `ip_whitelist` | string[] | ✗    | IP 白名单（可选）                    |
| `description`  | string   | ✗    | 描述信息                             |

**响应示例**:

```json
{
  "message": "token created successfully",
  "data": {
    "token": "pat_2Kj9X_aB3cD4eF5gH6iJ7kL8mN9oP0qR1sT2uV3wX4yZ",
    "id": 1,
    "name": "My API Integration",
    "token_prefix": "pat_2Kj9X",
    "permissions": ["user:profile:read", "user:profile:update", "api:cache:read"],
    "expires_at": "2025-02-19T10:30:00Z",
    "created_at": "2024-11-20T10:30:00Z"
  },
  "warning": "Please save this token now. You won't be able to see it again!"
}
```

**重要提示**:

- ⚠️ **完整 token 仅显示一次**，请立即保存
- Token 创建后无法查看明文，只能删除并重新创建
- 权限必须是您当前拥有权限的子集

### 权限选择

创建 Token 时，您可以选择需要的权限。系统会验证：

1. 所选权限必须是您已有权限的子集
2. 至少选择一个权限

**权限验证示例**:

```bash
# 假设用户拥有以下权限：
# - user:profile:read
# - user:profile:update
# - user:password:update
# - user:tokens:create
# - user:tokens:read
# - user:tokens:delete

# ✓ 有效请求（权限是子集）
{
  "permissions": ["user:profile:read", "user:profile:update"]
}

# ✗ 无效请求（包含未拥有的权限）
{
  "permissions": ["admin:users:create"]  // 用户没有 admin 域权限
}
```

## 使用 Personal Access Token

### 基本用法

使用 PAT 访问 API 时，将 token 放在 `Authorization` 头中：

```bash
curl -X GET http://localhost:8080/api/user/me \
  -H "Authorization: Bearer pat_2Kj9X_aB3cD4eF5gH6iJ7kL8mN9oP0qR1sT2uV3wX4yZ"
```

系统会自动识别 `pat_` 前缀并使用 PAT 认证流程。

### 在不同场景中使用

#### 1. cURL 命令

```bash
# 读取个人资料
curl -X GET http://localhost:8080/api/user/me \
  -H "Authorization: Bearer pat_2Kj9X_..."

# 更新个人资料
curl -X PUT http://localhost:8080/api/user/me \
  -H "Authorization: Bearer pat_2Kj9X_..." \
  -H "Content-Type: application/json" \
  -d '{"full_name": "Updated Name"}'
```

#### 2. JavaScript/Node.js

```javascript
const PAT = "pat_2Kj9X_aB3cD4eF5gH6iJ7kL8mN9oP0qR1sT2uV3wX4yZ";

// 使用 fetch
fetch("http://localhost:8080/api/user/me", {
  headers: {
    Authorization: `Bearer ${PAT}`,
  },
})
  .then((res) => res.json())
  .then((data) => console.log(data));

// 使用 axios
const axios = require("axios");

axios
  .get("http://localhost:8080/api/user/me", {
    headers: {
      Authorization: `Bearer ${PAT}`,
    },
  })
  .then((response) => console.log(response.data));
```

#### 3. Python

```python
import requests

PAT = 'pat_2Kj9X_aB3cD4eF5gH6iJ7kL8mN9oP0qR1sT2uV3wX4yZ'

headers = {
    'Authorization': f'Bearer {PAT}'
}

# GET 请求
response = requests.get('http://localhost:8080/api/user/me', headers=headers)
print(response.json())

# POST 请求
data = {'full_name': 'Updated Name'}
response = requests.put('http://localhost:8080/api/user/me', json=data, headers=headers)
print(response.json())
```

#### 4. Go

```go
package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    pat := "pat_2Kj9X_aB3cD4eF5gH6iJ7kL8mN9oP0qR1sT2uV3wX4yZ"

    req, _ := http.NewRequest("GET", "http://localhost:8080/api/user/me", nil)
    req.Header.Set("Authorization", "Bearer "+pat)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

## 管理 Personal Access Tokens

### 查看 Token 列表

**API 端点**: `GET /api/user/tokens`

```bash
curl -X GET http://localhost:8080/api/user/tokens \
  -H "Authorization: Bearer <your_jwt_token>"
```

**响应示例**:

```json
{
  "message": "tokens retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "My API Integration",
      "token_prefix": "pat_2Kj9X",
      "permissions": ["user:profile:read", "user:profile:update"],
      "expires_at": "2025-02-19T10:30:00Z",
      "last_used_at": "2024-11-20T15:20:00Z",
      "status": "active",
      "created_at": "2024-11-20T10:30:00Z"
    },
    {
      "id": 2,
      "name": "CLI Tool",
      "token_prefix": "pat_9XyZ3",
      "permissions": ["api:cache:read"],
      "expires_at": null,
      "last_used_at": null,
      "status": "active",
      "created_at": "2024-11-19T08:00:00Z"
    }
  ],
  "count": 2
}
```

### 查看单个 Token 详情

**API 端点**: `GET /api/user/tokens/:id`

```bash
curl -X GET http://localhost:8080/api/user/tokens/1 \
  -H "Authorization: Bearer <your_jwt_token>"
```

### 删除 Token

**API 端点**: `DELETE /api/user/tokens/:id`

```bash
curl -X DELETE http://localhost:8080/api/user/tokens/1 \
  -H "Authorization: Bearer <your_jwt_token>"
```

**响应**:

```json
{
  "message": "token disabled successfully"
}
```

**删除后**：

- Token 状态变为 `disabled`
- 立即失效，无法再用于 API 认证
- 记录保留在数据库中（软删除）

## 最佳实践

### 1. 最小权限原则

只授予 Token 完成任务所需的最小权限：

```bash
# ✓ 推荐：只读权限
{
  "permissions": ["user:profile:read", "api:cache:read"]
}

# ✗ 避免：授予所有权限
{
  "permissions": [所有用户权限]  // 过度授权
}
```

### 2. 设置合理的有效期

```bash
# ✓ 推荐：短期 Token 用于测试
{
  "expires_in": 7  // 7 天
}

# ✓ 推荐：中期 Token 用于生产
{
  "expires_in": 90  // 90 天
}

# ⚠️ 慎用：永久 Token
{
  "expires_in": null  // 需要定期审查
}
```

### 3. 使用描述性名称

```bash
# ✓ 推荐
{
  "name": "Production API Server - Cache Access",
  "description": "Used by prod-server-01 for cache operations"
}

# ✗ 避免
{
  "name": "token1"  // 无法识别用途
}
```

### 4. 配置 IP 白名单

如果 Token 只在特定服务器使用，启用 IP 限制：

```bash
{
  "ip_whitelist": ["192.168.1.100", "10.0.0.50"]
}
```

### 5. 定期轮换 Token

- 每 90 天轮换一次生产环境的 Token
- 删除旧 Token 前先部署新 Token
- 保持至少有一个备用 Token

### 6. 安全存储

**推荐方式**：

- ✓ 环境变量（生产环境）
- ✓ 密钥管理服务（AWS Secrets Manager, HashiCorp Vault）
- ✓ CI/CD 密钥存储

**避免**：

- ✗ 硬编码在代码中
- ✗ 提交到版本控制系统
- ✗ 明文存储在配置文件

**示例（环境变量）**：

```bash
# .env 文件（不要提交到 git）
API_TOKEN=pat_2Kj9X_aB3cD4eF5gH6iJ7kL8mN9oP0qR1sT2uV3wX4yZ

# 使用
export API_TOKEN=$(cat .env | grep API_TOKEN | cut -d '=' -f2)
curl -H "Authorization: Bearer $API_TOKEN" http://api.example.com/endpoint
```

## 安全建议

### Token 泄露处理

如果怀疑 Token 泄露：

1. **立即删除** 受影响的 Token

   ```bash
   curl -X DELETE http://localhost:8080/api/user/tokens/<token_id> \
     -H "Authorization: Bearer <jwt_token>"
   ```

2. **创建新 Token** 替换
3. **检查审计日志** 查看是否有异常访问
4. **更新部署** 使用新 Token

### 监控建议

定期检查：

- Token 的最后使用时间（`last_used_at`）
- 长期未使用的 Token（考虑删除）
- 即将过期的 Token（提前轮换）

### IP 白名单注意事项

- 如果服务器 IP 可能变化，不要使用 IP 限制
- 使用负载均衡时，需要添加所有出口 IP
- IP 限制失败会返回 `403 Forbidden`

## Token 状态说明

| 状态       | 说明     | 可用性                 |
| ---------- | -------- | ---------------------- |
| `active`   | 正常激活 | ✓ 可用                 |
| `disabled` | 已禁用   | ✗ 不可用（可重新启用） |
| `expired`  | 已过期   | ✗ 不可用               |

系统会自动标记过期的 Token，您可以禁用或删除不再需要的 Token。

## 常见问题

### Q1: Token 创建后忘记保存，如何找回？

**A**: 无法找回。Token 明文仅在创建时显示一次，数据库只存储 SHA-256 哈希值。您需要：

1. 删除旧 Token
2. 创建新 Token

### Q2: 如何增加 Token 的权限？

**A**: 无法修改现有 Token 的权限。需要：

1. 创建新 Token（包含所需权限）
2. 更新应用配置使用新 Token
3. 删除旧 Token

### Q3: Token 过期后会自动删除吗？

**A**: 不会。过期的 Token 仍保留在数据库中，状态为 `expired`，但无法用于认证。您可以手动删除以清理列表。

### Q4: 使用 PAT 时出现 403 错误？

**A**: 可能的原因：

1. **权限不足**：Token 没有所需权限
2. **Token 已删除**：状态为 `disabled`
3. **Token 已过期**：超过 `expires_at` 时间
4. **IP 限制**：请求 IP 不在白名单中

检查方法：

```bash
# 查看 Token 详情
curl -X GET http://localhost:8080/api/user/tokens/<token_id> \
  -H "Authorization: Bearer <jwt_token>"
```

### Q5: JWT 和 PAT 可以同时使用吗？

**A**: 是的。系统支持两种认证方式：

- Web 应用使用 JWT（短期，可刷新）
- API 集成使用 PAT（长期，特定权限）

选择依据：

- **需要自动刷新** → 使用 JWT
- **长期稳定访问** → 使用 PAT
- **限制权限范围** → 使用 PAT

### Q6: 如何批量管理 Token？

**A**: 目前需要逐个操作。建议：

1. 使用描述性名称和分类标签
2. 定期审查 Token 列表
3. 删除长期未使用的 Token

### Q7: PAT 支持通配符权限吗？

**A**: 支持。创建 Token 时可以使用通配符权限（如果用户拥有）：

```json
{
  "permissions": ["user:*:read"] // 用户域所有资源的读权限
}
```

但建议使用明确的权限列表以提高安全性。

## 相关文档

- [RBAC 权限系统](./identity-rbac.md) - 了解权限模型和三段式权限格式
- [认证授权](./identity-authentication.md) - JWT 认证流程
- [API 文档](/api/overview) - API 端点详细文档

## 技术实现

如果您是开发者，想了解 PAT 的技术实现细节：

**核心组件**：

- `internal/domain/pat/model.go` - PAT 领域模型
- `internal/infrastructure/auth/pat_service.go` - PAT 服务层
- `internal/infrastructure/auth/token_generator.go` - Token 生成器
- `internal/adapters/http/middleware/jwt.go` - 统一认证中间件
- `internal/adapters/http/handler/pat.go` - PAT HTTP 处理器

**数据库表**: `personal_access_tokens`

```sql
CREATE TABLE personal_access_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,  -- SHA-256 哈希
    token_prefix VARCHAR(20) NOT NULL,
    permissions JSONB NOT NULL,
    expires_at TIMESTAMP,
    last_used_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active',
    ip_whitelist JSONB,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```
