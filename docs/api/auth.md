# 认证接口

用户认证和授权相关的 API 接口。

## 用户注册

注册新用户账号。

### 请求

```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "full_name": "Test User"
}
```

### 请求参数

| 参数      | 类型   | 必填 | 描述                      |
| --------- | ------ | ---- | ------------------------- |
| username  | string | 是   | 用户名 (唯一)             |
| email     | string | 是   | 邮箱 (唯一)               |
| password  | string | 是   | 密码 (明文，服务端会加密) |
| full_name | string | 否   | 用户全名                  |

### 响应

```json
{
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "full_name": "Test User",
    "status": "active",
    "created_at": "2025-11-18T00:00:00Z",
    "updated_at": "2025-11-18T00:00:00Z"
  },
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "access_token_expires_at": "2025-11-18T00:15:00Z",
    "refresh_token_expires_at": "2025-11-25T00:00:00Z"
  }
}
```

### 错误响应

```json
{
  "error": "username already exists"
}
```

**可能的错误：**

- `username already exists` - 用户名已存在
- `email already exists` - 邮箱已存在
- `invalid request` - 请求参数无效

---

## 用户登录

使用用户名或邮箱登录。

### 请求

```http
POST /api/auth/login
Content-Type: application/json

{
  "login": "testuser",
  "password": "password123"
}
```

### 请求参数

| 参数     | 类型   | 必填 | 描述         |
| -------- | ------ | ---- | ------------ |
| login    | string | 是   | 用户名或邮箱 |
| password | string | 是   | 密码         |

### 响应

```json
{
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "full_name": "Test User",
    "status": "active",
    "created_at": "2025-11-18T00:00:00Z",
    "updated_at": "2025-11-18T00:00:00Z"
  },
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "access_token_expires_at": "2025-11-18T00:15:00Z",
    "refresh_token_expires_at": "2025-11-25T00:00:00Z"
  }
}
```

### 错误响应

```json
{
  "error": "invalid credentials"
}
```

**可能的错误：**

- `invalid credentials` - 用户名/邮箱或密码错误
- `user not found` - 用户不存在
- `user is not active` - 用户账号未激活

---

## 刷新令牌

使用刷新令牌获取新的访问令牌。

### 请求

```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### 请求参数

| 参数          | 类型   | 必填 | 描述     |
| ------------- | ------ | ---- | -------- |
| refresh_token | string | 是   | 刷新令牌 |

### 响应

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "access_token_expires_at": "2025-11-18T00:15:00Z",
  "refresh_token_expires_at": "2025-11-25T00:00:00Z"
}
```

### 错误响应

```json
{
  "error": "invalid or expired token"
}
```

**可能的错误：**

- `invalid or expired token` - 令牌无效或已过期
- `user not found` - 用户不存在
- `user is not active` - 用户账号未激活

---

## 获取当前用户信息

获取已登录用户的信息 (需要认证) 。

### 请求

```http
GET /api/auth/me
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### 响应

```json
{
  "id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "full_name": "Test User",
  "status": "active",
  "created_at": "2025-11-18T00:00:00Z",
  "updated_at": "2025-11-18T00:00:00Z"
}
```

### 错误响应

```json
{
  "error": "unauthorized"
}
```

**可能的错误：**

- `unauthorized` - 未提供有效的访问令牌
- `user not found` - 用户不存在

---

## 使用示例

### curl 示例

#### 注册

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

#### 登录

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "testuser",
    "password": "password123"
  }'
```

#### 刷新令牌

```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

#### 获取当前用户

```bash
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### JavaScript 示例

```javascript
// 注册
const register = async () => {
  const response = await fetch("http://localhost:8080/api/auth/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      username: "testuser",
      email: "test@example.com",
      password: "password123",
      full_name: "Test User",
    }),
  });
  return response.json();
};

// 登录
const login = async () => {
  const response = await fetch("http://localhost:8080/api/auth/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      login: "testuser",
      password: "password123",
    }),
  });
  return response.json();
};

// 获取当前用户
const getCurrentUser = async (accessToken) => {
  const response = await fetch("http://localhost:8080/api/auth/me", {
    headers: {
      Authorization: `Bearer ${accessToken}`,
    },
  });
  return response.json();
};
```

---

## Token 管理最佳实践

1. **存储 Token**：
   - 前端：存储在 `localStorage` 或 `sessionStorage`
   - 移动端：使用安全存储 (Keychain/KeyStore)

2. **Token 过期处理**：
   - 访问令牌过期时，使用刷新令牌自动获取新令牌
   - 刷新令牌过期时，需要重新登录

3. **安全建议**：
   - 使用 HTTPS
   - 不要在 URL 中传递 Token
   - 定期刷新 Token
   - 实现退出登录功能清除 Token

---

## 下一步

- 查看[用户接口文档](/api/users)
- 了解[认证授权机制](/backend/authentication)
- 学习[JWT 最佳实践](/backend/authentication#jwt-最佳实践)
