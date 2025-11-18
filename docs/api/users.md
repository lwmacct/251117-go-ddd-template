# 用户接口

用户管理相关的 API 接口 (需要 JWT 认证) 。

## 获取用户列表

获取所有用户的分页列表。

### 请求

```http
GET /api/users?page=1&page_size=10
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### 查询参数

| 参数      | 类型 | 必填 | 默认值 | 描述             |
| --------- | ---- | ---- | ------ | ---------------- |
| page      | int  | 否   | 1      | 页码 (从 1 开始) |
| page_size | int  | 否   | 10     | 每页数量         |

### 响应

```json
{
  "items": [
    {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "full_name": "Test User",
      "status": "active",
      "created_at": "2025-11-18T00:00:00Z",
      "updated_at": "2025-11-18T00:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 10
}
```

---

## 获取用户详情

根据 ID 获取单个用户的详细信息。

### 请求

```http
GET /api/users/:id
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### 路径参数

| 参数 | 类型 | 描述    |
| ---- | ---- | ------- |
| id   | int  | 用户 ID |

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
  "error": "user not found"
}
```

---

## 更新用户

更新用户信息。

### 请求

```http
PUT /api/users/:id
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "email": "newemail@example.com",
  "full_name": "New Name",
  "status": "active"
}
```

### 路径参数

| 参数 | 类型 | 描述    |
| ---- | ---- | ------- |
| id   | int  | 用户 ID |

### 请求体参数

| 参数      | 类型   | 必填 | 描述                                 |
| --------- | ------ | ---- | ------------------------------------ |
| email     | string | 否   | 新邮箱                               |
| full_name | string | 否   | 新全名                               |
| status    | string | 否   | 用户状态 (active/inactive/suspended) |

**注意：**

- `username` 不可修改
- `password` 需要通过专门的修改密码接口 (未实现)
- 只能更新提供的字段

### 响应

```json
{
  "id": 1,
  "username": "testuser",
  "email": "newemail@example.com",
  "full_name": "New Name",
  "status": "active",
  "created_at": "2025-11-18T00:00:00Z",
  "updated_at": "2025-11-18T00:15:00Z"
}
```

### 错误响应

```json
{
  "error": "user not found"
}
```

**可能的错误：**

- `user not found` - 用户不存在
- `email already exists` - 邮箱已被其他用户使用
- `invalid request` - 请求参数无效

---

## 删除用户

软删除用户 (标记为已删除，不真正从数据库删除) 。

### 请求

```http
DELETE /api/users/:id
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### 路径参数

| 参数 | 类型 | 描述    |
| ---- | ---- | ------- |
| id   | int  | 用户 ID |

### 响应

```http
204 No Content
```

### 错误响应

```json
{
  "error": "user not found"
}
```

**注意：**

- 这是软删除，用户数据会保留在数据库中
- 已删除的用户无法登录
- 已删除的用户不会出现在用户列表中

---

## 用户状态说明

| 状态      | 描述                           |
| --------- | ------------------------------ |
| active    | 激活状态，可以正常登录和使用   |
| inactive  | 未激活状态，可能需要邮箱验证等 |
| suspended | 暂停状态，已被管理员封禁       |

---

## 使用示例

### curl 示例

#### 获取用户列表

```bash
curl http://localhost:8080/api/users?page=1&page_size=20 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### 获取用户详情

```bash
curl http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### 更新用户

```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newemail@example.com",
    "full_name": "New Name"
  }'
```

#### 删除用户

```bash
curl -X DELETE http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### JavaScript 示例

```javascript
// 获取用户列表
const getUsers = async (accessToken, page = 1, pageSize = 10) => {
  const response = await fetch(`http://localhost:8080/api/users?page=${page}&page_size=${pageSize}`, {
    headers: {
      Authorization: `Bearer ${accessToken}`,
    },
  });
  return response.json();
};

// 获取用户详情
const getUser = async (accessToken, userId) => {
  const response = await fetch(`http://localhost:8080/api/users/${userId}`, {
    headers: {
      Authorization: `Bearer ${accessToken}`,
    },
  });
  return response.json();
};

// 更新用户
const updateUser = async (accessToken, userId, data) => {
  const response = await fetch(`http://localhost:8080/api/users/${userId}`, {
    method: "PUT",
    headers: {
      Authorization: `Bearer ${accessToken}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });
  return response.json();
};

// 删除用户
const deleteUser = async (accessToken, userId) => {
  const response = await fetch(`http://localhost:8080/api/users/${userId}`, {
    method: "DELETE",
    headers: {
      Authorization: `Bearer ${accessToken}`,
    },
  });
  return response.status === 204;
};
```

### Python 示例

```python
import requests

class UserAPI:
    def __init__(self, base_url, access_token):
        self.base_url = base_url
        self.headers = {
            'Authorization': f'Bearer {access_token}',
            'Content-Type': 'application/json'
        }

    def get_users(self, page=1, page_size=10):
        """获取用户列表"""
        response = requests.get(
            f'{self.base_url}/api/users',
            params={'page': page, 'page_size': page_size},
            headers=self.headers
        )
        return response.json()

    def get_user(self, user_id):
        """获取用户详情"""
        response = requests.get(
            f'{self.base_url}/api/users/{user_id}',
            headers=self.headers
        )
        return response.json()

    def update_user(self, user_id, data):
        """更新用户"""
        response = requests.put(
            f'{self.base_url}/api/users/{user_id}',
            json=data,
            headers=self.headers
        )
        return response.json()

    def delete_user(self, user_id):
        """删除用户"""
        response = requests.delete(
            f'{self.base_url}/api/users/{user_id}',
            headers=self.headers
        )
        return response.status_code == 204

# 使用示例
api = UserAPI('http://localhost:8080', 'YOUR_ACCESS_TOKEN')

# 获取用户列表
users = api.get_users(page=1, page_size=20)

# 获取用户详情
user = api.get_user(1)

# 更新用户
updated = api.update_user(1, {
    'email': 'newemail@example.com',
    'full_name': 'New Name'
})

# 删除用户
success = api.delete_user(1)
```

---

## 权限说明

当前版本所有已认证用户都可以：

- 查看所有用户列表
- 查看任意用户详情
- 更新任意用户信息
- 删除任意用户

**生产环境建议：**

- 实现基于角色的访问控制 (RBAC)
- 普通用户只能查看和修改自己的信息
- 管理员可以管理所有用户
- 添加审计日志记录敏感操作

---

## 下一步

- 查看[认证接口文档](/api/auth)
- 了解[认证授权机制](/guide/authentication)
- 学习[项目架构](/guide/architecture)
