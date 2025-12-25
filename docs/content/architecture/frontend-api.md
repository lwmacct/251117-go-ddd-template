# API 层设计

前端 API 层基于 Axios 封装，按业务模块组织，实现了统一的请求/响应处理、Token 认证和自动刷新机制。

<!--TOC-->

## Table of Contents

- [目录结构](#目录结构) `:29+26`
- [Axios 客户端配置](#axios-客户端配置) `:55+73`
  - [基础配置](#基础配置) `:57+13`
  - [请求拦截器](#请求拦截器) `:70+17`
  - [响应拦截器](#响应拦截器) `:87+41`
- [API 模块示例](#api-模块示例) `:128+82`
  - [认证 API](#认证-api) `:130+37`
  - [管理端用户 API](#管理端用户-api) `:167+43`
- [响应类型](#响应类型) `:210+61`
  - [统一响应格式](#统一响应格式) `:212+31`
  - [分页响应处理](#分页响应处理) `:243+28`
- [使用方式](#使用方式) `:271+73`
  - [在 Composable 中使用](#在-composable-中使用) `:273+52`
  - [在组件中使用](#在组件中使用) `:325+19`
- [错误处理](#错误处理) `:344+36`
  - [统一错误格式化](#统一错误格式化) `:346+34`
- [开发代理配置](#开发代理配置) `:380+22`

<!--TOC-->

## 目录结构

```
src/api/
├── index.ts                 # 统一导出
├── auth/                    # 认证相关 API
│   ├── index.ts
│   ├── client.ts            # Axios 实例 + 拦截器
│   ├── platformAuth.ts      # 认证 API (登录/注册/2FA)
│   └── user.ts              # 用户 API (个人信息)
├── admin/                   # 管理端 API
│   ├── index.ts
│   ├── users.ts             # 用户管理
│   ├── roles.ts             # 角色管理
│   ├── permissions.ts       # 权限管理
│   ├── menus.ts             # 菜单管理
│   ├── settings.ts          # 系统设置
│   ├── auditlogs.ts         # 审计日志
│   └── overview.ts          # 系统概览
├── user/                    # 用户端 API
│   ├── index.ts
│   └── tokens.ts            # Personal Access Token
└── helpers/
    └── pagination.ts        # 分页响应处理
```

## Axios 客户端配置

### 基础配置

```typescript
// src/api/auth/client.ts
import axios from "axios";

const apiClient = axios.create({
  timeout: 10000, // 10 秒超时
});

export default apiClient;
```

### 请求拦截器

自动添加 Authorization 头：

```typescript
apiClient.interceptors.request.use(
  (config) => {
    const token = getAccessToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error),
);
```

### 响应拦截器

处理 401 错误和 Token 刷新：

```typescript
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // 401 错误且未重试过
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        // 尝试刷新 Token
        const refreshToken = getRefreshToken();
        const response = await axios.post("/api/auth/refresh", {
          refresh_token: refreshToken,
        });

        // 保存新 Token
        const { access_token } = response.data.data;
        saveAccessToken(access_token);

        // 重试原请求
        originalRequest.headers.Authorization = `Bearer ${access_token}`;
        return apiClient(originalRequest);
      } catch (refreshError) {
        // 刷新失败，清除 Token 并跳转登录
        clearAuthTokens();
        window.location.href = "/#/auth/login";
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  },
);
```

## API 模块示例

### 认证 API

```typescript
// src/api/auth/platformAuth.ts
import apiClient from "./client";
import type { LoginRequest, AuthResponse, CaptchaData } from "@/types/auth";

export class AuthAPI {
  // 获取验证码
  static async getCaptcha(): Promise<CaptchaData> {
    const response = await apiClient.get("/api/auth/captcha");
    return response.data.data;
  }

  // 登录
  static async login(req: LoginRequest): Promise<AuthResponse> {
    const response = await apiClient.post("/api/auth/login", req);
    return response.data.data;
  }

  // 2FA 验证
  static async verify2FA(sessionToken: string, code: string) {
    const response = await apiClient.post("/api/auth/2fa/verify", {
      session_token: sessionToken,
      code,
    });
    return response.data.data;
  }

  // 获取 2FA 状态
  static async get2FAStatus() {
    const response = await apiClient.get("/api/auth/2fa/status");
    return response.data.data;
  }
}
```

### 管理端用户 API

```typescript
// src/api/admin/users.ts
import apiClient from "../auth/client";
import type { AdminUser, CreateUserRequest, UpdateUserRequest } from "@/types/admin";
import type { PaginationParams, PaginatedResponse } from "@/types/common";

// 获取用户列表
export async function listUsers(params: PaginationParams): Promise<PaginatedResponse<AdminUser>> {
  const response = await apiClient.get("/api/admin/users", { params });
  return normalizeListResponse(response.data, []);
}

// 获取单个用户
export async function getUser(id: number): Promise<AdminUser> {
  const response = await apiClient.get(`/api/admin/users/${id}`);
  return response.data.data;
}

// 创建用户
export async function createUser(data: CreateUserRequest): Promise<AdminUser> {
  const response = await apiClient.post("/api/admin/users", data);
  return response.data.data;
}

// 更新用户
export async function updateUser(id: number, data: UpdateUserRequest): Promise<AdminUser> {
  const response = await apiClient.put(`/api/admin/users/${id}`, data);
  return response.data.data;
}

// 删除用户
export async function deleteUser(id: number): Promise<void> {
  await apiClient.delete(`/api/admin/users/${id}`);
}

// 分配角色
export async function assignRoles(id: number, roleIds: number[]): Promise<void> {
  await apiClient.put(`/api/admin/users/${id}/roles`, { role_ids: roleIds });
}
```

## 响应类型

### 统一响应格式

```typescript
// src/types/response/index.ts

// 单个数据响应
interface ApiResponse<T = any> {
  code: number;
  message: string;
  data?: T;
  error?: any;
}

// 列表数据响应
interface ListApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
  meta?: PaginationMeta;
}

// 分页元数据
interface PaginationMeta {
  total: number;
  page: number;
  per_page: number;
  total_pages?: number;
  has_more?: boolean;
}
```

### 分页响应处理

```typescript
// src/api/helpers/pagination.ts

// 标准化列表响应
export function normalizeListResponse<T>(payload: any, fallback: T[]): PaginatedResponse<T> {
  // 处理不同格式的后端响应
  const data = payload.data ?? fallback;
  const meta = payload.meta ?? {
    page: 1,
    limit: 10,
    total: data.length,
    total_pages: 1,
  };

  return {
    data,
    pagination: {
      page: meta.page,
      limit: meta.per_page ?? meta.limit,
      total: meta.total,
      total_pages: meta.total_pages ?? Math.ceil(meta.total / meta.per_page),
    },
  };
}
```

## 使用方式

### 在 Composable 中使用

```typescript
// src/pages/admin/users/composables/useAdminUsers.ts
import { ref, reactive } from "vue";
import { listUsers, createUser, deleteUser } from "@/api/admin/users";
import type { AdminUser } from "@/types/admin";

export function useAdminUsers() {
  const users = ref<AdminUser[]>([]);
  const loading = ref(false);
  const pagination = reactive({
    page: 1,
    limit: 10,
    total: 0,
  });

  async function fetchUsers() {
    loading.value = true;
    try {
      const response = await listUsers({
        page: pagination.page,
        limit: pagination.limit,
      });
      users.value = response.data;
      pagination.total = response.pagination.total;
    } finally {
      loading.value = false;
    }
  }

  async function handleCreate(data: CreateUserRequest) {
    await createUser(data);
    await fetchUsers(); // 刷新列表
  }

  async function handleDelete(id: number) {
    await deleteUser(id);
    await fetchUsers();
  }

  return {
    users,
    loading,
    pagination,
    fetchUsers,
    handleCreate,
    handleDelete,
  };
}
```

### 在组件中使用

```vue
<script setup lang="ts">
import { onMounted } from "vue";
import { useAdminUsers } from "./composables/useAdminUsers";

const { users, loading, pagination, fetchUsers } = useAdminUsers();

onMounted(() => {
  fetchUsers();
});
</script>

<template>
  <v-data-table :items="users" :loading="loading" :items-per-page="pagination.limit" />
</template>
```

## 错误处理

### 统一错误格式化

```typescript
// src/utils/auth/error.ts
export function formatAuthError(error: any): string {
  // 优先使用服务器返回的错误信息
  if (error.response?.data?.error) {
    return error.response.data.error;
  }
  if (error.response?.data?.message) {
    return error.response.data.message;
  }

  // HTTP 状态码错误
  if (error.response?.status) {
    const statusMessages: Record<number, string> = {
      400: "请求参数错误",
      401: "未授权，请重新登录",
      403: "权限不足",
      404: "资源不存在",
      500: "服务器内部错误",
    };
    return statusMessages[error.response.status] ?? "请求失败";
  }

  // 网络错误
  if (error.message === "Network Error") {
    return "网络连接失败，请检查网络";
  }

  return "未知错误";
}
```

## 开发代理配置

Vite 开发服务器配置 API 代理：

```typescript
// vite.config.ts
export default defineConfig({
  server: {
    port: 40013,
    proxy: {
      "/api": {
        target: "http://localhost:40012",
        changeOrigin: true,
      },
      "/swagger": {
        target: "http://localhost:40012",
        changeOrigin: true,
      },
    },
  },
});
```
