# 认证授权

本项目实现了完整的 JWT 认证授权系统，提供用户注册、登录、Token 刷新等功能。

<!--TOC-->

## Table of Contents

- [功能特性](#功能特性) `:27+9`
- [快速开始](#快速开始) `:36+18`
- [架构设计](#架构设计) `:54+20`
- [工作流程](#工作流程) `:74+23`
  - [注册流程](#注册流程) `:76+7`
  - [登录流程](#登录流程) `:83+7`
  - [Token 刷新流程](#token-刷新流程) `:90+7`
- [配置说明](#配置说明) `:97+14`
- [API 端点](#api-端点) `:111+18`
  - [公开端点 (无需认证)](#公开端点-无需认证) `:113+8`
  - [受保护端点 (需要认证)](#受保护端点-需要认证) `:121+8`
- [安全特性](#安全特性) `:129+11`
- [故障排查](#故障排查) `:140+26`
  - [Token 验证失败](#token-验证失败) `:142+10`
  - [用户无法登录](#用户无法登录) `:152+14`

<!--TOC-->

## 功能特性

- ✅ 用户注册 (用户名/邮箱唯一性验证)
- ✅ 用户登录 (支持用户名或邮箱登录)
- ✅ Token 刷新机制 (访问令牌 15 分钟，刷新令牌 7 天)
- ✅ JWT 认证中间件
- ✅ bcrypt 密码加密
- ✅ 用户状态检查 (仅 active 用户可登录)

## 快速开始

```bash
# 注册用户
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "email": "test@example.com", "password": "password123"}'

# 用户登录（login 可以是用户名或邮箱）
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login": "testuser", "password": "password123"}'

# 使用 Token 访问受保护端点
curl http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 架构设计

```
internal/
├── domain/auth/           # 认证领域服务接口
├── application/auth/      # Register/Login/Refresh Use Case
├── infrastructure/auth/   # JWT 管理器、服务实现
├── adapters/http/
│   ├── handler/auth.go    # HTTP Handler
│   └── middleware/jwt.go  # JWT 鉴权中间件
└── bootstrap/container.go # 依赖注入
```

**分层职责**：

- **Domain** 定义 `auth.Service` 接口（密码验证、Token 生成）
- **Infrastructure** 实现接口（bcrypt + JWT）
- **Application** 协调业务流程（验证 → 查用户 → 生成 Token）
- **Adapters** 仅做请求绑定和响应转换

## 工作流程

### 注册流程

```
用户提交注册信息 → 验证用户名/邮箱唯一 → bcrypt 加密密码
    → 创建用户记录 → 生成 JWT Token 对 → 返回 token 和用户信息
```

### 登录流程

```
提交用户名/邮箱和密码 → 查找用户 → 验证密码 (bcrypt)
    → 检查用户状态 (active) → 生成 JWT Token 对 → 返回 token
```

### Token 刷新流程

```
提交 refresh_token → 验证 token → 提取用户 ID
    → 查询用户并检查状态 → 生成新的 token 对
```

## 配置说明

```yaml
# config.yaml
jwt:
  secret: "your-secret-key-change-in-production"
  access_token_expiry: "15m" # 访问令牌过期时间
  refresh_token_expiry: "168h" # 刷新令牌过期时间 (7天)
```

**环境变量**：`APP_JWT_SECRET`、`APP_JWT_ACCESS_TOKEN_EXPIRY`、`APP_JWT_REFRESH_TOKEN_EXPIRY`

**重要**：生产环境必须使用强密钥（至少 32 字节随机字符串）。

## API 端点

### 公开端点 (无需认证)

| 方法 | 路径                 | 说明         |
| ---- | -------------------- | ------------ |
| POST | `/api/auth/register` | 注册新用户   |
| POST | `/api/auth/login`    | 用户登录     |
| POST | `/api/auth/refresh`  | 刷新访问令牌 |

### 受保护端点 (需要认证)

| 方法 | 路径                | 说明             |
| ---- | ------------------- | ---------------- |
| GET  | `/api/user/profile` | 获取当前用户信息 |

详细的 API 文档请通过 Swagger UI (`/swagger/index.html`) 查看。

## 安全特性

| 特性       | 说明                                     |
| ---------- | ---------------------------------------- |
| 密码加密   | bcrypt (cost=10)，响应中自动隐藏         |
| Token 签名 | HMAC-SHA256，包含过期时间                |
| 过期控制   | Access Token 15 分钟，Refresh Token 7 天 |
| 状态检查   | 仅 `active` 用户可登录                   |
| 唯一性约束 | 用户名、邮箱数据库层面强制唯一           |
| 错误处理   | 登录失败返回通用 "invalid credentials"   |

## 故障排查

### Token 验证失败

```bash
# 检查 token 内容
echo $ACCESS_TOKEN | cut -d'.' -f2 | base64 -d | jq '.'

# 确认 secret 配置
env | grep APP_JWT_SECRET
```

### 用户无法登录

- 检查用户状态（必须是 `active`）
- 确认用户名/邮箱正确
- 验证密码正确性
- 查看应用日志

---

**相关文档**：

- [RBAC 权限系统](./identity-rbac.md) - 权限模型
- [Personal Access Token](./identity-pat.md) - API Token
- Swagger UI (`/swagger/index.html`) - API 详细文档
