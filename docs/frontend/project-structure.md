# 项目结构

本文档详细介绍前端项目的目录组织、文件命名规范和模块划分，所有示例均来自当前 `web/` 目录。

## 目录结构

```
web/
├── public/                     # 静态资源
│   └── favicon.ico
├── src/
│   ├── api/                    # API 客户端（按业务拆分）
│   │   ├── auth/
│   │   ├── user/
│   │   ├── admin/
│   │   └── index.ts
│   ├── global/                 # 全局样式/配置占位
│   ├── layout/                 # 布局组件 (AdminLayout.vue / UserLayout.vue)
│   ├── pages/                  # 页面（admin/auth/user 三大模块）
│   ├── router/                 # 路由定义（admin.ts、auth.ts、user.ts）
│   ├── stores/                 # Pinia Store
│   ├── types/                  # 类型定义（auth/admin/user 等子目录）
│   ├── utils/                  # 工具库（以 auth/token/storage 为核心）
│   ├── views/                  # 共享 UI 片段（AppBars、Navigation）
│   ├── App.vue
│   └── main.ts
├── env.d.ts
├── index.html
├── package.json
├── tsconfig*.json
└── vite.config.ts
```

## 核心目录说明

### `src/api/` - API 模块

接口按照业务归类并共享 `auth/client.ts` 中的 Axios 实例：

```
api/
├── auth/
│   ├── auth.ts          # 登录、注册、刷新 Token
│   ├── user.ts          # /api/auth/user/* 接口
│   ├── platformAuth.ts  # 带验证码的登录流程
│   └── client.ts        # Axios + 拦截器
├── user/
│   └── tokens.ts        # PAT API
├── admin/
│   └── ...              # 预留管理端接口
└── index.ts             # 对外导出
```

```ts
// src/api/auth/auth.ts
import { apiClient } from "./client";
import { saveAccessToken, saveRefreshToken } from "@/utils/auth";
import type { LoginRequest, AuthResponse, ApiResponse } from "@/types/auth";

export const login = async (req: LoginRequest): Promise<AuthResponse> => {
  const { data } = await apiClient.post<ApiResponse<AuthResponse>>("/login", req);
  if (data.data) {
    saveAccessToken(data.data.access_token);
    saveRefreshToken(data.data.refresh_token);
    return data.data;
  }
  throw new Error(data.error || "Login failed");
};
```

```ts
// src/api/user/tokens.ts
import { apiClient } from "../auth/client";
import type { PersonalAccessToken, CreateTokenRequest, CreateTokenResponse } from "@/types/user";
import type { ApiResponse } from "@/types/auth";

export const listTokens = async (): Promise<PersonalAccessToken[]> => {
  const { data } = await apiClient.get<ApiResponse<PersonalAccessToken[]>>("/user/tokens");
  if (data.data) {
    return data.data;
  }
  throw new Error(data.error || "获取 Token 列表失败");
};
```

### `src/layout/` - 布局组件

- `AdminLayout.vue`：包裹所有 `/admin/*` 页面，内置导航、侧栏等。
- `UserLayout.vue`：用户中心布局。

布局组件负责加载 `views/Navigation`、`views/AppBars` 等共享片段，并决定面包屑/权限粒度。

### `src/pages/` - 页面模块

以业务域划分目录，与后端的 DDD 模块一一对应：

```
pages/
├── admin/
│   ├── overview/
│   ├── roles/
│   ├── users/
│   ├── menus/
│   └── settings/
├── auth/
│   ├── login/
│   └── register/
└── user/
    ├── profile/
    ├── security/
    └── tokens/
```

- 每个目录包含一个 `index.vue`，并在路由中懒加载。
- 页面聚合 `stores/`、`api/`，不直接访问 HTTP 层。

### `src/router/` - 路由配置

路由同样模块化，`router/index.ts` 汇总三个子路由：

```ts
// src/router/index.ts
import { createRouter, createWebHashHistory } from "vue-router";
import { adminRoutes } from "./admin";
import { authRoutes } from "./auth";
import { userRoutes } from "./user";

export default createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [{ path: "/", redirect: "/auth/login" }, authRoutes, adminRoutes, userRoutes],
});
```

`router/admin.ts` 展示了完整的 Children 配置、`meta.requiresAuth`、懒加载写法，与后端 RBAC 路由保持一致。

### `src/stores/` - Pinia 状态

`useAuthStore` 负责登录/注册、2FA 会话、用户资料缓存。示例：

```ts
// src/stores/auth.ts
export const useAuthStore = defineStore("auth", () => {
  const currentUser = ref<User | null>(null);
  const isLoading = ref(false);
  const error = ref<string | null>(null);
  const isAuthenticated = computed(() => !!currentUser.value);

  async function initAuth() {
    if (!getAccessToken()) {
      currentUser.value = null;
      return;
    }
    try {
      isLoading.value = true;
      currentUser.value = await getCurrentUser();
    } finally {
      isLoading.value = false;
    }
  }

  return { currentUser, isAuthenticated, initAuth };
});
```

其他 Store（如 `counter.ts`）主要用于示例或局部状态。

### `src/types/` - 类型系统

类型按业务拆分（`auth/`、`admin/`、`user/`、`common/`），集中导出方便引入：

```ts
// src/types/auth/auth.ts
export interface LoginRequest {
  login: string;
  password: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
  user: User;
}
```

保持领域术语一致，方便与后端 DTO 对齐。

### `src/utils/` - 工具库

当前以认证工具为主：

- `auth/storage.ts`：localStorage 读写 token。
- `auth/token.ts`：解析 JWT、判断过期。
- `auth/validation.ts`：客户端校验规则。

```ts
// src/utils/auth/token.ts
export function parseJwtToken(token: string) {
  try {
    const [_, payload] = token.split(".");
    const base64 = payload.replace(/-/g, "+").replace(/_/g, "/");
    return JSON.parse(atob(base64));
  } catch (error) {
    console.error("Failed to parse JWT token:", error);
    return null;
  }
}
```

### `src/views/` - 共享视图片段

`AppBars`、`Navigation` 等目录存放跨页面复用的 UI 组合件，而不是零散的基础组件。它们通常被布局或页面直接引用，避免再引入额外的 `components/` 层。

---

通过上述结构，前端与后端 DDD 模块保持相同的域划分（auth/admin/user），API 与 Store 只处理与自身领域相关的逻辑，使得维护和协作更清晰。

```
types/
├── api.ts          # API 请求/响应类型
├── models.ts       # 数据模型
└── index.ts        # 类型导出
```

**示例**:

```typescript
// types/models.ts
export interface User {
  id: number;
  username: string;
  email: string;
}

export interface LoginRequest {
  login: string;
  password: string;
}
```

### `src/utils/` - 工具函数

**职责**: 通用工具函数

**常见工具**:

```
utils/
├── request.ts      # HTTP 请求封装
├── storage.ts      # LocalStorage/SessionStorage
├── format.ts       # 格式化函数
└── validate.ts     # 表单验证
```

## 文件命名规范

### Vue 组件

**格式**: `PascalCase.vue`

```
✓ UserProfile.vue
✓ DataTable.vue
✓ LoginForm.vue

✗ userProfile.vue
✗ user-profile.vue
✗ data_table.vue
```

### TypeScript 文件

**格式**: `kebab-case.ts` 或 `camelCase.ts`

```
✓ auth.ts
✓ user-api.ts
✓ requestClient.ts

✗ Auth.ts
✗ User-Api.ts
```

### 目录

**格式**: `kebab-case`

```
✓ api/
✓ user-management/
✓ auth-pages/

✗ API/
✗ UserManagement/
```

## 导入路径别名

### `@` 别名

配置在 `vite.config.ts` 和 `tsconfig.json`：

```typescript
// 使用 @ 别名
import { userApi } from "@/api/users";
import UserCard from "@/components/UserCard.vue";

// 等同于
import { userApi } from "../api/users";
import UserCard from "../components/UserCard.vue";
```

### 推荐的导入顺序

```typescript
// 1. Vue 核心
import { ref, computed } from "vue";

// 2. Vue 生态库
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";

// 3. 第三方库
import axios from "axios";

// 4. 项目内部
import { userApi } from "@/api/users";
import UserCard from "@/components/UserCard.vue";

// 5. 类型
import type { User } from "@/types/models";
```

## 代码组织原则

### 1. 单一职责

每个文件/组件只做一件事：

```
✓ UserCard.vue         # 显示用户信息卡片
✓ UserList.vue         # 显示用户列表
✓ UserForm.vue         # 用户表单

✗ User.vue             # 职责不明确
```

### 2. 组件大小

单个组件不超过 300 行，超过则拆分：

```vue
<!-- ✗ 太大 -->
<template>
  <!-- 500 行代码 -->
</template>

<!-- ✓ 拆分 -->
<!-- UserProfile.vue -->
<template>
  <UserHeader />
  <UserDetails />
  <UserActions />
</template>
```

### 3. 功能模块化

按功能模块组织相关文件：

```
views/
└── user-management/
    ├── UserList.vue
    ├── UserDetail.vue
    ├── UserEdit.vue
    └── components/
        ├── UserTable.vue
        └── UserFilter.vue
```

## 最佳实践

### 组件设计

**可复用组件** (`components/`):

- 无业务逻辑
- 通过 props 接收数据
- 通过 emit 触发事件

**页面组件** (`pages/`):

- 包含业务逻辑
- 调用 API
- 管理状态

### API 封装

```typescript
// ✓ 推荐：统一封装
export const userApi = {
  list: (params) => client.get("/api/users", { params }),
  get: (id) => client.get(`/api/users/${id}`),
  create: (data) => client.post("/api/users", data),
};

// ✗ 避免：直接调用
axios.get("/api/users");
```

### 类型定义

```typescript
// ✓ 推荐：定义类型
interface User {
  id: number;
  name: string;
}

const users = ref<User[]>([]);

// ✗ 避免：any
const users = ref<any>([]);
```

## 扩展阅读

- [API 集成](./api-integration) - API 使用指南

<!-- TODO: 待完善的文档
- [开发规范](./coding-standards) - 代码风格指南
- [组件规范](./components) - 组件设计原则
-->

熟悉项目结构后，开始愉快地编码吧！ 🎨
