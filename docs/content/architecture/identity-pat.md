# Personal Access Token (PAT)

PAT 是一种用于 API 认证的长期凭证，作为 JWT Token 的替代方案。

<!--TOC-->

## Table of Contents

- [适用场景](#适用场景) `:29+7`
- [PAT vs JWT](#pat-vs-jwt) `:36+11`
- [Token 格式](#token-格式) `:47+12`
- [API 端点](#api-端点) `:59+36`
  - [创建 Token](#创建-token) `:68+27`
- [使用方法](#使用方法) `:95+8`
- [最佳实践](#最佳实践) `:103+33`
  - [权限控制](#权限控制) `:105+5`
  - [有效期设置](#有效期设置) `:110+8`
  - [安全存储](#安全存储) `:118+13`
  - [Token 轮换](#token-轮换) `:131+5`
- [Token 状态](#token-状态) `:136+8`
- [常见问题](#常见问题) `:144+31`
  - [Q1: Token 创建后忘记保存，如何找回？](#q1-token-创建后忘记保存如何找回) `:146+4`
  - [Q2: 如何修改 Token 权限？](#q2-如何修改-token-权限) `:150+4`
  - [Q3: 使用 PAT 时出现 403 错误？](#q3-使用-pat-时出现-403-错误) `:154+8`
  - [Q4: JWT 和 PAT 可以同时使用吗？](#q4-jwt-和-pat-可以同时使用吗) `:162+13`

<!--TOC-->

## 适用场景

- **API 集成**：第三方应用调用 API
- **CLI 工具**：命令行工具自动化操作
- **自动化脚本**：定时任务、部署脚本
- **Webhook 回调**：接收和处理 Webhook 事件

## PAT vs JWT

| 特性    | JWT Token     | Personal Access Token   |
| ------- | ------------- | ----------------------- |
| 用途    | Web/移动应用  | API 集成、CLI、脚本     |
| 有效期  | 短期（1小时） | 长期（7/30/90天或永久） |
| 刷新    | Refresh Token | 无需刷新                |
| 权限    | 用户全部权限  | 用户权限的子集          |
| 删除    | 不支持        | 支持即时删除            |
| IP 限制 | 不支持        | 支持 IP 白名单          |

## Token 格式

**格式**：`pat_<5位前缀>_<32位随机字符>`

**示例**：`pat_2Kj9X_aB3cD4eF5gH6iJ7kL8mN9oP0qR1sT2uV3wX4yZ`

**安全特性**：

- 完整 token 仅在创建时显示一次
- 数据库存储 SHA-256 哈希值
- 前缀用于快速识别

## API 端点

| 方法   | 路径                   | 说明            | 权限                 |
| ------ | ---------------------- | --------------- | -------------------- |
| POST   | `/api/user/tokens`     | 创建 Token      | `user:tokens:create` |
| GET    | `/api/user/tokens`     | 查看 Token 列表 | `user:tokens:read`   |
| GET    | `/api/user/tokens/:id` | 查看 Token 详情 | `user:tokens:read`   |
| DELETE | `/api/user/tokens/:id` | 删除 Token      | `user:tokens:delete` |

### 创建 Token

```bash
curl -X POST http://localhost:8080/api/user/tokens \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My API Integration",
    "permissions": ["user:profile:read", "api:cache:read"],
    "expires_in": 90,
    "ip_whitelist": ["192.168.1.100"],
    "description": "Token for production API"
  }'
```

**请求参数**：

| 字段           | 类型     | 必填 | 说明                               |
| -------------- | -------- | ---- | ---------------------------------- |
| `name`         | string   | ✓    | Token 名称（1-100 字符）           |
| `permissions`  | string[] | ✓    | 权限列表（必须是用户已有权限子集） |
| `expires_in`   | int      | ✗    | 有效期天数（7/30/90 或 null=永久） |
| `ip_whitelist` | string[] | ✗    | IP 白名单                          |
| `description`  | string   | ✗    | 描述信息                           |

**重要**：⚠️ 完整 token 仅显示一次，请立即保存。

## 使用方法

```bash
# 使用 PAT 访问 API（系统自动识别 pat_ 前缀）
curl http://localhost:8080/api/user/me \
  -H "Authorization: Bearer pat_2Kj9X_..."
```

## 最佳实践

### 权限控制

- ✓ 只授予完成任务所需的最小权限
- ✗ 避免授予所有用户权限

### 有效期设置

| 场景     | 建议有效期 |
| -------- | ---------- |
| 测试环境 | 7 天       |
| 生产环境 | 90 天      |
| 永久     | 需定期审查 |

### 安全存储

**推荐**：

- ✓ 环境变量
- ✓ 密钥管理服务（AWS Secrets Manager, Vault）
- ✓ CI/CD 密钥存储

**避免**：

- ✗ 硬编码在代码中
- ✗ 提交到版本控制系统

### Token 轮换

- 每 90 天轮换生产环境 Token
- 删除旧 Token 前先部署新 Token

## Token 状态

| 状态       | 说明     | 可用性                 |
| ---------- | -------- | ---------------------- |
| `active`   | 正常激活 | ✓ 可用                 |
| `disabled` | 已禁用   | ✗ 不可用（可重新启用） |
| `expired`  | 已过期   | ✗ 不可用               |

## 常见问题

### Q1: Token 创建后忘记保存，如何找回？

无法找回。Token 明文仅创建时显示一次，数据库只存储哈希值。需删除旧 Token 并重新创建。

### Q2: 如何修改 Token 权限？

无法修改。需创建新 Token（包含所需权限）→ 更新应用配置 → 删除旧 Token。

### Q3: 使用 PAT 时出现 403 错误？

可能原因：

1. Token 没有所需权限
2. Token 已删除或过期
3. 请求 IP 不在白名单中

### Q4: JWT 和 PAT 可以同时使用吗？

可以。系统自动识别认证方式：

- **需要自动刷新** → 使用 JWT
- **长期稳定访问** → 使用 PAT

---

**相关文档**：

- [RBAC 权限系统](./identity-rbac.md) - 权限模型和三段式格式
- [认证授权](./identity-authentication.md) - JWT 认证流程
