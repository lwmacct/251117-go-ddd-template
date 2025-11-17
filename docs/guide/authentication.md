# 认证授权

本项目实现了完整的 JWT 认证授权系统，提供用户注册、登录、Token 刷新等功能。

## 功能特性

- ✅ 用户注册（用户名/邮箱唯一性验证）
- ✅ 用户登录（支持用户名或邮箱登录）
- ✅ Token 刷新机制（访问令牌 15 分钟，刷新令牌 7 天）
- ✅ JWT 认证中间件
- ✅ bcrypt 密码加密
- ✅ 用户状态检查（仅 active 用户可登录）
- ✅ 受保护的 API 端点

## 快速开始

### 1. 确保服务运行

```bash
# 启动数据库和 Redis
docker-compose up -d

# 启动应用
task go:run -- api
```

### 2. 注册用户

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

**响应示例：**

```json
{
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "full_name": "Test User",
    "status": "active",
    "created_at": "2025-11-18T00:00:00Z"
  },
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "access_token_expires_at": "2025-11-18T00:15:00Z",
    "refresh_token_expires_at": "2025-11-25T00:00:00Z"
  }
}
```

### 3. 用户登录

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "testuser",
    "password": "password123"
  }'
```

**注意：** `login` 字段可以是用户名或邮箱。

### 4. 使用 Token 访问受保护端点

```bash
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 架构设计

### 代码结构

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
    └── container.go       # 依赖注入
```

### JWT Claims 结构

```go
type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}
```

## 工作流程

### 注册流程

```
用户提交注册信息
    ↓
验证用户名和邮箱是否已存在
    ↓
密码 bcrypt 加密
    ↓
创建用户记录（状态：active）
    ↓
生成 JWT token 对（access + refresh）
    ↓
返回 token 和用户信息
```

### 登录流程

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

### Token 验证流程

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

### Token 刷新流程

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

## 配置说明

### 配置文件

在 `config.yaml` 中配置 JWT 参数：

```yaml
jwt:
  secret: "your-secret-key-change-in-production"
  access_token_expiry: "15m" # 访问令牌过期时间
  refresh_token_expiry: "168h" # 刷新令牌过期时间（7天）
```

### 环境变量

```bash
export APP_JWT_SECRET="your-secret-key"
export APP_JWT_ACCESS_TOKEN_EXPIRY="15m"
export APP_JWT_REFRESH_TOKEN_EXPIRY="168h"
```

**重要提示：** 生产环境必须使用强密钥！建议使用至少 32 字节的随机字符串。

生成强密钥：

```bash
# 使用 openssl 生成随机密钥
openssl rand -base64 32
```

## API 端点

### 公开端点（无需认证）

| 方法 | 路径                 | 说明         |
| ---- | -------------------- | ------------ |
| POST | `/api/auth/register` | 注册新用户   |
| POST | `/api/auth/login`    | 用户登录     |
| POST | `/api/auth/refresh`  | 刷新访问令牌 |

### 受保护端点（需要认证）

| 方法   | 路径             | 说明             |
| ------ | ---------------- | ---------------- |
| GET    | `/api/auth/me`   | 获取当前用户信息 |
| GET    | `/api/users`     | 获取用户列表     |
| GET    | `/api/users/:id` | 获取用户详情     |
| PUT    | `/api/users/:id` | 更新用户         |
| DELETE | `/api/users/:id` | 删除用户         |

详细的 API 文档请参考 [认证接口](/api/auth) 和 [用户接口](/api/users)。

## 安全特性

1. **密码加密**

   - 使用 bcrypt 加密存储密码
   - 成本因子：10（默认）
   - 密码字段在响应中自动隐藏

2. **Token 签名**

   - 使用 HMAC-SHA256 算法签名
   - Secret 密钥从配置读取
   - Token 包含过期时间

3. **Token 过期控制**

   - 访问令牌：15 分钟（短期）
   - 刷新令牌：7 天（长期）
   - 可通过配置调整

4. **用户状态检查**

   - 只有 `active` 状态的用户可以登录
   - 支持 `inactive`、`suspended` 等状态

5. **唯一性约束**

   - 用户名唯一
   - 邮箱唯一
   - 数据库层面强制约束

6. **参数验证**

   - 使用 Gin binding 验证请求参数
   - 邮箱格式验证
   - 密码最小长度验证

7. **错误处理**
   - 登录失败返回通用错误 "invalid credentials"
   - 避免泄露用户是否存在的信息

## 使用示例

### 完整的认证流程

```bash
#!/bin/bash

# 1. 注册用户
echo "=== 注册用户 ==="
REGISTER_RESP=$(curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "email": "john@example.com",
    "password": "password123",
    "full_name": "John Doe"
  }')

echo $REGISTER_RESP | jq '.'

# 2. 提取 token
ACCESS_TOKEN=$(echo $REGISTER_RESP | jq -r '.tokens.access_token')
REFRESH_TOKEN=$(echo $REGISTER_RESP | jq -r '.tokens.refresh_token')

echo "Access Token: $ACCESS_TOKEN"

# 3. 获取当前用户信息
echo -e "\n=== 获取当前用户 ==="
curl -s http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq '.'

# 4. 获取用户列表
echo -e "\n=== 获取用户列表 ==="
curl -s "http://localhost:8080/api/users?page=1&page_size=10" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq '.'

# 5. Token 过期后刷新
echo -e "\n=== 刷新 Token ==="
curl -s -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}" | jq '.'
```

### 在代码中使用认证服务

```go
// 从容器获取认证服务
authService := container.AuthService

// 注册用户
user, tokens, err := authService.Register(ctx, &auth.RegisterRequest{
    Username: "testuser",
    Email:    "test@example.com",
    Password: "password123",
    FullName: "Test User",
})

// 登录
user, tokens, err = authService.Login(ctx, &auth.LoginRequest{
    Login:    "testuser",
    Password: "password123",
})

// 刷新 Token
tokens, err = authService.RefreshToken(ctx, refreshToken)
```

## JWT 最佳实践

### 1. Token 存储

**前端应用：**

- 推荐使用 `localStorage` 或 `sessionStorage`
- 避免在 Cookie 中存储（如果不需要）

**移动应用：**

- iOS: 使用 Keychain
- Android: 使用 KeyStore

### 2. Token 刷新策略

```javascript
// 前端示例：自动刷新 Token
async function apiCall(url, options) {
  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        Authorization: `Bearer ${localStorage.getItem("access_token")}`,
        ...options.headers,
      },
    });

    if (response.status === 401) {
      // Token 过期，尝试刷新
      const newTokens = await refreshToken();
      localStorage.setItem("access_token", newTokens.access_token);

      // 重试原请求
      return fetch(url, {
        ...options,
        headers: {
          Authorization: `Bearer ${newTokens.access_token}`,
          ...options.headers,
        },
      });
    }

    return response;
  } catch (error) {
    console.error("API call failed:", error);
    throw error;
  }
}
```

### 3. 安全建议

- ✅ 使用 HTTPS
- ✅ 不要在 URL 中传递 Token
- ✅ 定期刷新 Token
- ✅ 实现退出登录功能清除 Token
- ✅ 生产环境使用强密钥
- ✅ 监控异常登录行为

### 4. Token 黑名单（可选）

对于需要支持"注销"功能的场景，可以实现 Token 黑名单：

```go
// 使用 Redis 存储已注销的 Token
func (s *Service) Logout(ctx context.Context, token string) error {
    // 解析 token 获取过期时间
    claims, _ := s.jwtManager.ValidateToken(token)
    expiration := time.Until(claims.ExpiresAt.Time)

    // 将 token 加入黑名单（设置 TTL 为剩余有效期）
    return s.redis.Set(ctx, "blacklist:"+token, "1", expiration).Err()
}

// 在中间件中检查黑名单
func (m *JWTMiddleware) checkBlacklist(token string) bool {
    exists, _ := m.redis.Exists(ctx, "blacklist:"+token).Result()
    return exists > 0
}
```

## 扩展功能建议

### 1. 邮箱验证

```go
// 注册后发送验证邮件
func (s *Service) Register(...) {
    user.Status = "inactive"
    user.VerificationToken = generateToken()
    // 发送验证邮件
    sendVerificationEmail(user.Email, user.VerificationToken)
}

// 验证邮箱端点
func (h *Handler) VerifyEmail(c *gin.Context) {
    token := c.Query("token")
    // 验证 token 并激活用户
}
```

### 2. 密码重置

```go
// 忘记密码
func (s *Service) ForgotPassword(email string) {
    user := findUserByEmail(email)
    resetToken := generateToken()
    // 发送重置邮件
}

// 重置密码
func (s *Service) ResetPassword(token, newPassword string) {
    // 验证 token 并更新密码
}
```

### 3. 双因素认证（2FA）

```go
// 启用 2FA
func (s *Service) Enable2FA(userID uint) (qrCode string) {
    secret := generateTOTPSecret()
    // 生成 QR 码
}

// 登录时验证 2FA
func (s *Service) Verify2FA(code string) bool {
    // 验证 TOTP 代码
}
```

### 4. OAuth 社交登录

```go
// GitHub OAuth
func (h *Handler) LoginWithGitHub(c *gin.Context) {
    // OAuth 回调处理
}
```

## 故障排查

### Token 验证失败

```bash
# 检查 token 格式
echo $ACCESS_TOKEN | cut -d'.' -f2 | base64 -d | jq '.'

# 确认 secret 配置正确
env | grep APP_JWT_SECRET
```

### 密码验证失败

- 确认密码加密正确
- 检查密码最小长度要求
- 查看应用日志

### 用户无法登录

- 检查用户状态（必须是 `active`）
- 确认用户名/邮箱正确
- 验证密码正确性

## 下一步

- 查看 [API 认证接口文档](/api/auth)
- 了解 [用户管理接口](/api/users)
- 学习 [配置系统](/guide/configuration)
