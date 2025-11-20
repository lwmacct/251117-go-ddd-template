# 用户管理 API 文档

## 概述

本文档描述了用户管理功能的前后端 API 对应关系和使用说明。

## API 端点

### 1. 创建用户

**端点**: `POST /api/admin/users`
**权限**: `admin:users:create`
**请求体**:
```json
{
  "username": "string (必填)",
  "email": "string (必填)",
  "password": "string (必填)",
  "full_name": "string (可选)",
  "status": "active | inactive (可选, 默认: active)",
  "role_ids": [1, 2, 3] // 可选，同时分配角色
}
```

**响应**:
```json
{
  "code": 201,
  "message": "user created successfully",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "full_name": "Test User",
      "status": "active",
      "roles": [...],
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### 2. 获取用户列表

**端点**: `GET /api/admin/users`
**权限**: `admin:users:read`
**查询参数**:
- `page`: 页码（默认: 1）
- `limit`: 每页数量（默认: 20，最大: 100）
- `search`: 搜索关键词（支持用户名、邮箱、全名模糊匹配）

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com",
        "full_name": "Administrator",
        "status": "active",
        "roles": [...],
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  },
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### 3. 获取用户详情

**端点**: `GET /api/admin/users/:id`
**权限**: `admin:users:read`

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "full_name": "Administrator",
    "status": "active",
    "roles": [
      {
        "id": 1,
        "name": "admin",
        "display_name": "管理员",
        "permissions": [...]
      }
    ],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 4. 更新用户

**端点**: `PUT /api/admin/users/:id`
**权限**: `admin:users:update`
**请求体**:
```json
{
  "email": "string (可选)",
  "full_name": "string (可选)",
  "avatar": "string (可选)",
  "bio": "string (可选)",
  "status": "active | inactive | banned (可选)"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "user updated successfully",
  "data": {
    "user": {
      "id": 1,
      "username": "admin",
      "email": "newemail@example.com",
      "full_name": "Updated Name",
      "status": "active",
      "roles": [...],
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  }
}
```

### 5. 删除用户

**端点**: `DELETE /api/admin/users/:id`
**权限**: `admin:users:delete`

**响应**:
```json
{
  "code": 200,
  "message": "user deleted successfully",
  "data": null
}
```

### 6. 分配角色

**端点**: `PUT /api/admin/users/:id/roles`
**权限**: `admin:users:update`
**请求体**:
```json
{
  "role_ids": [1, 2, 3]
}
```

**响应**:
```json
{
  "code": 200,
  "message": "roles assigned successfully",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "full_name": "Administrator",
    "status": "active",
    "roles": [
      {
        "id": 1,
        "name": "admin",
        "display_name": "管理员"
      },
      {
        "id": 2,
        "name": "user",
        "display_name": "普通用户"
      }
    ],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

## 前端实现

### API 客户端

位于 `web/src/api/admin/users.ts`：

```typescript
import { apiClient } from "../auth/client";
import type { AdminUser, CreateUserRequest, UpdateUserRequest, AssignRolesRequest } from "@/types/admin";

// 获取用户列表
export const listUsers = async (params: Partial<PaginationParams>): Promise<PaginatedResponse<AdminUser>>

// 获取用户详情
export const getUser = async (id: number): Promise<AdminUser>

// 创建用户
export const createUser = async (params: CreateUserRequest): Promise<AdminUser>

// 更新用户
export const updateUser = async (id: number, params: UpdateUserRequest): Promise<AdminUser>

// 删除用户
export const deleteUser = async (id: number): Promise<void>

// 分配角色
export const assignRoles = async (id: number, params: AssignRolesRequest): Promise<AdminUser>
```

### Composable

位于 `web/src/pages/admin/users/composables/useAdminUsers.ts`，提供：

- `fetchUsers()`: 获取用户列表
- `createUser()`: 创建用户
- `updateUser()`: 更新用户
- `deleteUser()`: 删除用户
- `assignRoles()`: 分配角色
- `searchUsers()`: 搜索用户
- `changePage()`: 翻页

### 组件

- **UserDialog.vue**: 用户创建/编辑对话框
- **RoleSelector.vue**: 角色选择器

## 后端架构

### 目录结构

```
internal/
├── adapters/http/handler/
│   └── admin_user.go              # HTTP Handler（请求/响应处理）
├── application/user/
│   ├── command/
│   │   ├── create_user.go         # 创建用户 Command
│   │   ├── create_user_handler.go
│   │   ├── update_user.go         # 更新用户 Command
│   │   ├── update_user_handler.go
│   │   ├── delete_user.go         # 删除用户 Command
│   │   ├── delete_user_handler.go
│   │   ├── assign_roles.go        # 分配角色 Command
│   │   └── assign_roles_handler.go
│   ├── query/
│   │   ├── get_user.go            # 获取用户 Query
│   │   ├── get_user_handler.go
│   │   ├── list_users.go          # 获取列表 Query
│   │   └── list_users_handler.go
│   └── dto.go                     # DTO 定义
├── domain/user/
│   ├── entity_user.go             # 用户实体
│   ├── command_repository.go      # 写仓储接口
│   ├── query_repository.go        # 读仓储接口
│   └── errors.go                  # 领域错误
└── infrastructure/persistence/
    ├── user_model.go              # GORM Model + 映射
    ├── user_command_repository.go # 写仓储实现
    └── user_query_repository.go   # 读仓储实现
```

### 关键实现

#### 1. 响应格式一致性

所有 API 使用统一的响应格式（`internal/adapters/http/response`）：
- `response.OK()` - 200 成功
- `response.Created()` - 201 创建成功
- `response.List()` - 200 列表（含分页）
- `response.BadRequest()` - 400 错误请求
- `response.NotFound()` - 404 未找到
- `response.InternalError()` - 500 内部错误

#### 2. 搜索功能实现

- **Handler 层**: 从 query 参数获取 `search` 关键词
- **Application 层**: 根据关键词选择 `Search()` 或 `List()` 方法
- **Infrastructure 层**: 使用 LIKE 查询（支持用户名、邮箱、全名）

```go
// QueryRepository 接口
Search(ctx context.Context, keyword string, offset, limit int) ([]*User, error)
CountBySearch(ctx context.Context, keyword string) (int64, error)

// 实现
WHERE username LIKE '%keyword%' OR email LIKE '%keyword%' OR full_name LIKE '%keyword%'
```

#### 3. 角色分配优化

`AssignRoles` API 现在返回更新后的用户对象（包含角色信息），而不仅仅是成功消息：

```go
// 分配角色后获取完整用户信息
updatedUser, err := h.getUserHandler.Handle(ctx, userQuery.GetUserQuery{
    UserID:    uint(id),
    WithRoles: true,
})
response.OK(c, "roles assigned successfully", updatedUser)
```

## 测试

### 测试脚本

位于 `testing/test_admin_users.py`：

```bash
# 使用 uv 运行测试
uv run testing/test_admin_users.py
```

测试覆盖：
1. ✅ 登录获取 Token
2. ✅ 获取用户列表
3. ✅ 创建用户
4. ✅ 获取用户详情
5. ✅ 更新用户
6. ✅ 搜索用户
7. ✅ 分配角色
8. ✅ 删除用户

## 近期修改

### 2024-11-20

1. **修复响应格式不匹配**
   - AssignRoles 现在返回更新后的用户对象（包含角色）
   - 与前端期望的响应格式完全对应

2. **添加搜索功能**
   - Handler 层支持 `search` 查询参数
   - Application 层根据关键词选择搜索或列表方法
   - Infrastructure 层实现 `CountBySearch()` 方法

3. **修复 Swagger 文档**
   - AssignRoles 路由从 POST 改为 PUT（与实际路由一致）
   - 添加 search 参数文档

4. **前端完整实现**
   - 用户列表页面（支持分页、搜索）
   - 用户创建/编辑对话框
   - 角色选择器
   - 所有 CRUD 操作

## 注意事项

1. **权限控制**: 所有 admin 路由都需要 `admin` 角色和相应的权限
2. **用户名不可修改**: 创建后用户名不能更改（编辑对话框中禁用）
3. **密码更新**: 编辑时留空密码则不修改密码
4. **角色分配**: 会覆盖现有角色（不是追加）
5. **搜索**: 支持用户名、邮箱、全名的模糊匹配

## 相关文档

- DDD + CQRS 架构 - 参考 `docs/architecture/` 目录
- API 开发指南 - 参考 Swagger 文档 `/swagger/index.html`
- 权限系统 - 基于 RBAC 模型的三段式权限控制
