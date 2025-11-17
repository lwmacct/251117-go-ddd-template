# API 参考

本文档提供完整的 REST API 接口参考。

## 基础信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **认证方式**: JWT Bearer Token

## 接口概览

### 公开接口（无需认证）

| 方法   | 路径                 | 描述     |
| ------ | -------------------- | -------- |
| GET    | `/health`            | 健康检查 |
| POST   | `/api/auth/register` | 用户注册 |
| POST   | `/api/auth/login`    | 用户登录 |
| POST   | `/api/auth/refresh`  | 刷新令牌 |
| POST   | `/api/cache/:key`    | 设置缓存 |
| GET    | `/api/cache/:key`    | 获取缓存 |
| DELETE | `/api/cache/:key`    | 删除缓存 |

### 受保护接口（需要 JWT）

| 方法   | 路径             | 描述             |
| ------ | ---------------- | ---------------- |
| GET    | `/api/auth/me`   | 获取当前用户信息 |
| GET    | `/api/users`     | 获取用户列表     |
| GET    | `/api/users/:id` | 获取用户详情     |
| PUT    | `/api/users/:id` | 更新用户         |
| DELETE | `/api/users/:id` | 删除用户         |

## 认证接口

详见 [认证接口文档](/api/auth)

## 用户接口

详见 [用户接口文档](/api/users)

## 通用响应格式

### 成功响应

```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

### 错误响应

```json
{
  "error": "错误描述信息"
}
```

## HTTP 状态码

| 状态码 | 说明                       |
| ------ | -------------------------- |
| 200    | 请求成功                   |
| 201    | 资源创建成功               |
| 400    | 请求参数错误               |
| 401    | 未认证或认证失败           |
| 403    | 权限不足                   |
| 404    | 资源不存在                 |
| 409    | 资源冲突（如用户名已存在） |
| 500    | 服务器内部错误             |

## 认证机制

受保护的接口需要在请求头中包含 JWT Token：

```bash
Authorization: Bearer <access_token>
```

### 获取 Token

1. 注册或登录获取 Token
2. 在后续请求中使用 Token
3. Token 过期后使用刷新令牌获取新 Token

### Token 有效期

- **访问令牌（Access Token）**: 15 分钟
- **刷新令牌（Refresh Token）**: 7 天

## 分页

列表接口支持分页，使用以下查询参数：

| 参数      | 类型 | 默认值 | 描述              |
| --------- | ---- | ------ | ----------------- |
| page      | int  | 1      | 页码（从 1 开始） |
| page_size | int  | 10     | 每页数量          |

### 示例

```bash
GET /api/users?page=2&page_size=20
```

### 响应格式

```json
{
  "items": [...],
  "total": 100,
  "page": 2,
  "page_size": 20
}
```

## 错误代码

| 错误代码          | 描述           |
| ----------------- | -------------- |
| `INVALID_REQUEST` | 请求参数无效   |
| `UNAUTHORIZED`    | 未认证         |
| `FORBIDDEN`       | 权限不足       |
| `NOT_FOUND`       | 资源不存在     |
| `CONFLICT`        | 资源冲突       |
| `INTERNAL_ERROR`  | 服务器内部错误 |

## Postman 集合

你可以导入以下 Postman 集合进行 API 测试：

[下载 Postman 集合](./postman_collection.json)

## 速率限制

当前版本暂未实施速率限制，生产环境建议添加。

## 下一步

- 查看[认证接口详情](/api/auth)
- 查看[用户接口详情](/api/users)
- 了解[认证授权机制](/guide/authentication)
