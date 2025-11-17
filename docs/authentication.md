# JWT 认证实现说明

## 概述

本项目已实现完整的 JWT 认证系统，包括：
- 用户注册
- 用户登录
- Token 刷新
- JWT 认证中间件
- 受保护的 API 端点

## 快速开始

### 1. 确保数据库和 Redis 运行

```bash
docker-compose up -d
```

### 2. 启动应用

```bash
task go:run -- api
```

### 3. 测试认证功能

#### 注册新用户

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'
```

响应示例：
```json
{
  "message": "user registered successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900,
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "full_name": "Test User",
      "status": "active",
      "created_at": "2025-01-18T00:00:00Z"
    }
  }
}
```

#### 用户登录

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "testuser",
    "password": "password123"
  }'
```

**注意**：`login` 字段可以是用户名或邮箱。

#### 刷新 Token

```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

#### 获取当前用户信息（需要认证）

```bash
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### 访问受保护的端点

```bash
# 获取用户列表（需要认证）
curl -X GET "http://localhost:8080/api/users?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 获取用户详情（需要认证）
curl -X GET http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 更新用户（需要认证）
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"full_name": "Updated Name"}'
```

## API 端点

### 公开端点（无需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/auth/register | 注册新用户 |
| POST | /api/auth/login | 用户登录 |
| POST | /api/auth/refresh | 刷新访问令牌 |
| GET | /health | 健康检查 |
| POST/GET/DELETE | /api/cache/* | 缓存操作（示例） |

### 受保护端点（需要认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/auth/me | 获取当前用户信息 |
| GET | /api/users | 获取用户列表 |
| GET | /api/users/:id | 获取用户详情 |
| PUT | /api/users/:id | 更新用户 |
| DELETE | /api/users/:id | 删除用户 |

## 代码结构

```
internal/
├── infrastructure/
│   └── auth/
│       ├── jwt.go         # JWT token 生成和验证
│       └── service.go     # 认证服务（注册、登录、刷新）
├── adapters/
│   └── http/
│       ├── handler/
│       │   └── auth.go    # 认证 HTTP 处理器
│       └── middleware/
│           └── jwt.go     # JWT 认证中间件
└── bootstrap/
    └── container.go       # 依赖注入（包含 JWT 管理器和认证服务）
```

## JWT 配置

### 配置文件（config.yaml）

```yaml
jwt:
  secret: "your-secret-key-change-in-production"  # 生产环境必须修改！
  access_token_expiry: "15m"   # 访问令牌过期时间（15分钟）
  refresh_token_expiry: "168h" # 刷新令牌过期时间（7天）
```

### 环境变量

```bash
export APP_JWT_SECRET="your-secret-key"
export APP_JWT_ACCESS_TOKEN_EXPIRY="15m"
export APP_JWT_REFRESH_TOKEN_EXPIRY="168h"
```

**重要**：生产环境必须使用强密钥！建议使用至少 32 字节的随机字符串。

## JWT Claims 结构

```go
type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}
```

## 工作流程

### 1. 注册流程

```
用户提交注册信息
    ↓
验证用户名和邮箱是否已存在
    ↓
密码 bcrypt 加密
    ↓
创建用户记录
    ↓
生成 JWT token 对（access + refresh）
    ↓
返回 token 和用户信息
```

### 2. 登录流程

```
用户提交用户名/邮箱和密码
    ↓
查找用户（支持用户名或邮箱）
    ↓
验证密码（bcrypt.CompareHashAndPassword）
    ↓
检查用户状态（必须是 active）
    ↓
生成 JWT token 对
    ↓
返回 token 和用户信息
```

### 3. Token 验证流程（中间件）

```
从请求头提取 Authorization: Bearer <token>
    ↓
验证 token 格式
    ↓
解析和验证 JWT token
    ↓
提取用户信息（UserID, Username, Email）
    ↓
将用户信息存入 Gin Context
    ↓
继续处理请求
```

### 4. Token 刷新流程

```
用户提交 refresh_token
    ↓
验证 refresh_token
    ↓
从 token 提取用户 ID
    ↓
查询用户信息
    ↓
检查用户状态
    ↓
生成新的 token 对
    ↓
返回新的 token
```

## 安全特性

1. **密码加密**：使用 bcrypt 加密存储
2. **Token 签名**：使用 HMAC-SHA256 签名
3. **Token 过期**：访问令牌15分钟，刷新令牌7天
4. **用户状态检查**：只有 active 用户可以登录
5. **唯一性约束**：用户名和邮箱唯一
6. **参数验证**：使用 Gin binding 验证

## 注意事项

1. **JWT Secret**：
   - 开发环境可以使用默认值
   - **生产环境必须修改为强密钥**
   - 建议使用环境变量设置

2. **Token 过期时间**：
   - 访问令牌（access_token）：短期，用于 API 调用
   - 刷新令牌（refresh_token）：长期，用于获取新的访问令牌
   - 根据安全需求调整时间

3. **密码强度**：
   - 当前最小长度：6 字符
   - 建议在生产环境增加密码复杂度要求

4. **错误处理**：
   - 登录失败返回通用错误"invalid credentials"
   - 避免泄露用户是否存在的信息

## 使用示例

### 完整的认证流程示例

```bash
# 1. 注册用户
REGISTER_RESP=$(curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "email": "john@example.com",
    "password": "password123",
    "full_name": "John Doe"
  }')

echo "Register Response:"
echo $REGISTER_RESP | jq '.'

# 2. 提取 access_token
ACCESS_TOKEN=$(echo $REGISTER_RESP | jq -r '.data.access_token')
REFRESH_TOKEN=$(echo $REGISTER_RESP | jq -r '.data.refresh_token')

echo "Access Token: $ACCESS_TOKEN"

# 3. 使用 token 访问受保护端点
curl -s http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq '.'

# 4. 获取用户列表
curl -s "http://localhost:8080/api/users?page=1&limit=10" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq '.'

# 5. Token 过期后刷新
curl -s -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}" | jq '.'
```

## 扩展建议

1. **添加功能**：
   - 密码重置
   - 邮箱验证
   - 双因素认证（2FA）
   - 社交登录（OAuth）
   - Remember me 功能

2. **安全增强**：
   - Token 黑名单（用于注销）
   - IP 白名单/黑名单
   - 登录尝试限制
   - 设备管理
   - 审计日志

3. **性能优化**：
   - Token 缓存到 Redis
   - 用户会话管理
   - 并发登录控制

4. **监控和日志**：
   - 登录成功/失败日志
   - Token 生成/验证统计
   - 异常登录检测
