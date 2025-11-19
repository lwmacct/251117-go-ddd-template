# API 集成

本文档介绍如何在前端应用中集成后端 API，包括 Axios 配置、请求拦截器、错误处理等。

## Axios 客户端配置

### 创建客户端实例

**文件**: `src/api/client.ts`

```typescript
import axios from "axios";
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from "axios";
import { useAuthStore } from "@/stores/auth";
import router from "@/router";

// API 基础 URL
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

// 创建 Axios 实例
const client: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000, // 30 秒超时
  headers: {
    "Content-Type": "application/json",
  },
});

// 请求拦截器
client.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    const authStore = useAuthStore();

    // 自动添加 JWT Token
    if (authStore.accessToken) {
      config.headers = config.headers || {};
      config.headers.Authorization = `Bearer ${authStore.accessToken}`;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// 响应拦截器
client.interceptors.response.use(
  (response: AxiosResponse) => {
    return response.data; // 直接返回 data
  },
  async (error) => {
    const originalRequest = error.config;

    // 401 未授权：刷新 Token
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const authStore = useAuthStore();
        await authStore.refreshToken();

        // 重试原请求
        return client(originalRequest);
      } catch (refreshError) {
        // 刷新失败，跳转登录
        router.push("/login");
        return Promise.reject(refreshError);
      }
    }

    // 403 权限不足
    if (error.response?.status === 403) {
      // 显示权限不足提示
      console.error("权限不足");
    }

    return Promise.reject(error);
  },
);

export default client;
```

## API 接口封装

### 认证 API

**文件**: `src/api/auth.ts`

```typescript
import client from "./client";
import type { LoginRequest, LoginResponse, RegisterRequest } from "@/types/api";

export const authApi = {
  /**
   * 用户登录
   */
  login(data: LoginRequest): Promise<LoginResponse> {
    return client.post("/api/auth/login", data);
  },

  /**
   * 用户注册
   */
  register(data: RegisterRequest): Promise<void> {
    return client.post("/api/auth/register", data);
  },

  /**
   * 刷新 Token
   */
  refreshToken(refreshToken: string): Promise<LoginResponse> {
    return client.post("/api/auth/refresh", { refresh_token: refreshToken });
  },

  /**
   * 退出登录
   */
  logout(): Promise<void> {
    return client.post("/api/auth/logout");
  },
};
```

### 用户 API

**文件**: `src/api/users.ts`

```typescript
import client from "./client";
import type { User, UpdateProfileRequest } from "@/types/api";

export const userApi = {
  /**
   * 获取当前用户信息
   */
  getProfile(): Promise<User> {
    return client.get("/api/user/me");
  },

  /**
   * 更新个人资料
   */
  updateProfile(data: UpdateProfileRequest): Promise<User> {
    return client.put("/api/user/me", data);
  },

  /**
   * 修改密码
   */
  changePassword(oldPassword: string, newPassword: string): Promise<void> {
    return client.put("/api/user/me/password", {
      old_password: oldPassword,
      new_password: newPassword,
    });
  },

  /**
   * 删除账户
   */
  deleteAccount(): Promise<void> {
    return client.delete("/api/user/me");
  },
};
```

### Personal Access Token API

**文件**: `src/api/tokens.ts`

```typescript
import client from "./client";
import type { CreateTokenRequest, TokenResponse, TokenListItem } from "@/types/api";

export const tokenApi = {
  /**
   * 创建 Personal Access Token
   */
  create(data: CreateTokenRequest): Promise<TokenResponse> {
    return client.post("/api/user/tokens", data);
  },

  /**
   * 列出所有 Token
   */
  list(): Promise<TokenListItem[]> {
    return client.get("/api/user/tokens");
  },

  /**
   * 获取 Token 详情
   */
  get(id: number): Promise<TokenListItem> {
    return client.get(`/api/user/tokens/${id}`);
  },

  /**
   * 撤销 Token
   */
  revoke(id: number): Promise<void> {
    return client.delete(`/api/user/tokens/${id}`);
  },
};
```

## TypeScript 类型定义

### API 类型

**文件**: `src/types/api.ts`

```typescript
// ========== 认证相关 ==========

export interface LoginRequest {
  login: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  full_name: string;
}

export interface LoginResponse {
  message: string;
  data: {
    access_token: string;
    refresh_token: string;
    token_type: string;
    expires_in: number;
    user: User;
  };
}

// ========== 用户相关 ==========

export interface User {
  id: number;
  username: string;
  email: string;
  full_name: string;
  status: "active" | "inactive" | "banned";
  created_at: string;
  updated_at: string;
}

export interface UpdateProfileRequest {
  full_name?: string;
  email?: string;
}

// ========== Token 相关 ==========

export interface CreateTokenRequest {
  name: string;
  permissions: string[];
  expires_in?: number;
  ip_whitelist?: string[];
  description?: string;
}

export interface TokenResponse {
  token: string;
  id: number;
  name: string;
  token_prefix: string;
  permissions: string[];
  expires_at: string | null;
  created_at: string;
}

export interface TokenListItem {
  id: number;
  name: string;
  token_prefix: string;
  permissions: string[];
  expires_at: string | null;
  last_used_at: string | null;
  status: "active" | "revoked" | "expired";
  created_at: string;
}

// ========== 通用响应 ==========

export interface ApiResponse<T> {
  message: string;
  data: T;
}

export interface ApiError {
  error: string;
  details?: any;
}
```

## 在组件中使用

### Composition API 方式

```vue
<script setup lang="ts">
import { ref, onMounted } from "vue";
import { userApi } from "@/api/users";
import type { User } from "@/types/api";

const user = ref<User | null>(null);
const loading = ref(false);
const error = ref<string | null>(null);

// 获取用户信息
const fetchProfile = async () => {
  loading.value = true;
  error.value = null;

  try {
    user.value = await userApi.getProfile();
  } catch (err: any) {
    error.value = err.response?.data?.error || "获取用户信息失败";
  } finally {
    loading.value = false;
  }
};

// 更新个人资料
const updateProfile = async (data: UpdateProfileRequest) => {
  try {
    user.value = await userApi.updateProfile(data);
    // 显示成功提示
  } catch (err: any) {
    error.value = err.response?.data?.error || "更新失败";
  }
};

onMounted(() => {
  fetchProfile();
});
</script>

<template>
  <div>
    <v-progress-linear v-if="loading" indeterminate />
    <v-alert v-if="error" type="error">{{ error }}</v-alert>
    <div v-if="user">
      <h2>{{ user.full_name }}</h2>
      <p>{{ user.email }}</p>
    </div>
  </div>
</template>
```

### 在 Store 中使用

```typescript
// stores/user.ts
import { defineStore } from "pinia";
import { userApi } from "@/api/users";
import type { User } from "@/types/api";

export const useUserStore = defineStore("user", {
  state: () => ({
    profile: null as User | null,
    loading: false,
    error: null as string | null,
  }),

  actions: {
    async fetchProfile() {
      this.loading = true;
      this.error = null;

      try {
        this.profile = await userApi.getProfile();
      } catch (err: any) {
        this.error = err.response?.data?.error || "获取用户信息失败";
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async updateProfile(data: UpdateProfileRequest) {
      try {
        this.profile = await userApi.updateProfile(data);
      } catch (err: any) {
        this.error = err.response?.data?.error || "更新失败";
        throw err;
      }
    },
  },
});
```

## 错误处理

### 统一错误处理

**创建错误处理工具**:

```typescript
// utils/error-handler.ts
import type { AxiosError } from "axios";

export interface ApiErrorResponse {
  error: string;
  details?: any;
}

export class ApiError extends Error {
  constructor(
    public status: number,
    public data: ApiErrorResponse,
  ) {
    super(data.error);
  }
}

export const handleApiError = (error: unknown): string => {
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError<ApiErrorResponse>;

    // 网络错误
    if (!axiosError.response) {
      return "网络连接失败，请检查网络设置";
    }

    // HTTP 错误
    const { status, data } = axiosError.response;

    switch (status) {
      case 400:
        return data.error || "请求参数错误";
      case 401:
        return "未授权，请重新登录";
      case 403:
        return "权限不足";
      case 404:
        return "请求的资源不存在";
      case 500:
        return "服务器错误";
      default:
        return data.error || `请求失败 (${status})`;
    }
  }

  return "未知错误";
};
```

**在组件中使用**:

```typescript
import { handleApiError } from "@/utils/error-handler";

try {
  await userApi.updateProfile(data);
} catch (err) {
  const message = handleApiError(err);
  // 显示错误消息
  showErrorNotification(message);
}
```

## 请求取消

### 取消正在进行的请求

```typescript
import { ref } from "vue";
import axios from "axios";

const controller = ref<AbortController | null>(null);

const fetchData = async () => {
  // 取消之前的请求
  if (controller.value) {
    controller.value.abort();
  }

  // 创建新的 AbortController
  controller.value = new AbortController();

  try {
    const data = await userApi.getProfile({
      signal: controller.value.signal,
    });
    // 处理数据
  } catch (err) {
    if (axios.isCancel(err)) {
      console.log("请求已取消");
    }
  }
};

// 组件卸载时取消请求
onUnmounted(() => {
  if (controller.value) {
    controller.value.abort();
  }
});
```

## 环境配置

### 开发环境 vs 生产环境

**`.env.development`**:

```bash
VITE_API_BASE_URL=http://localhost:8080
```

**`.env.production`**:

```bash
VITE_API_BASE_URL=https://api.production.com
```

**在代码中使用**:

```typescript
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
```

## 最佳实践

### 1. 统一接口封装

```typescript
// ✓ 推荐
export const userApi = {
  getProfile: () => client.get("/api/user/me"),
};

// ✗ 避免
axios.get("http://localhost:8080/api/user/me");
```

### 2. 使用 TypeScript 类型

```typescript
// ✓ 推荐
getProfile(): Promise<User>

// ✗ 避免
getProfile(): Promise<any>
```

### 3. 错误处理

```typescript
// ✓ 推荐
try {
  await userApi.getProfile();
} catch (err) {
  handleApiError(err);
}

// ✗ 避免
userApi.getProfile(); // 忽略错误
```

### 4. 加载状态

```typescript
// ✓ 推荐
const loading = ref(false);
loading.value = true;
try {
  await fetchData();
} finally {
  loading.value = false;
}
```

## 相关文档

- [认证授权](/backend/authentication) - JWT 认证机制
- [Personal Access Token](/backend/pat) - PAT 使用指南
- [API 参考](/api/) - 后端 API 详细文档
<!-- TODO: 待完善的文档
- [状态管理](./state-management) - Pinia Store 使用
  -->

开始高效地集成 API 吧！ 🚀
