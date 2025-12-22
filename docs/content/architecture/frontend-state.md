# 状态管理

前端使用 Pinia 进行全局状态管理，主要用于认证状态的跨组件共享。

<!--TOC-->

## Table of Contents

- [Store 结构](#store-结构) `:27+9`
- [认证 Store](#认证-store) `:36+122`
  - [状态定义](#状态定义) `:40+47`
  - [初始化流程](#初始化流程) `:87+23`
  - [登录流程](#登录流程) `:110+33`
  - [登出流程](#登出流程) `:143+15`
- [在组件中使用](#在组件中使用) `:158+40`
  - [Options API](#options-api) `:160+19`
  - [Composition API](#composition-api) `:179+19`
- [Token 存储](#token-存储) `:198+35`
- [最佳实践](#最佳实践) `:233+25`
  - [1. 状态最小化](#1-状态最小化) `:235+4`
  - [2. 使用 Getters](#2-使用-getters) `:239+4`
  - [3. Action 异步处理](#3-action-异步处理) `:243+4`
  - [4. 错误处理](#4-错误处理) `:247+11`

<!--TOC-->

## Store 结构

```
src/stores/
├── index.ts       # 统一导出
├── auth.ts        # 认证状态
└── counter.ts     # 计数器示例
```

## 认证 Store

认证 Store 是最核心的状态管理模块，负责管理用户登录状态和 Token。

### 状态定义

```typescript
// src/stores/auth.ts
import { defineStore } from 'pinia'
import type { User } from '@/types/auth'

interface AuthState {
  currentUser: User | null    // 当前用户信息
  isLoading: boolean          // 加载状态
  error: string | null        // 错误信息
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    currentUser: null,
    isLoading: false,
    error: null
  }),

  getters: {
    // 是否已认证（有用户信息）
    isAuthenticated: (state) => !!state.currentUser,

    // 是否有 Token（localStorage 中）
    hasToken: () => !!getAccessToken()
  },

  actions: {
    // 初始化认证状态
    async initAuth() { ... },

    // 登录
    async login(credentials: LoginRequest) { ... },

    // 登出
    async logout() { ... },

    // 更新用户信息
    updateUser(user: User) { ... },

    // 清除错误
    clearError() { ... }
  }
})
```

### 初始化流程

应用启动时自动恢复认证状态：

```typescript
// src/main.ts
const app = createApp(App);
const pinia = createPinia();
app.use(pinia);

// 恢复认证状态（在挂载前）
const authStore = useAuthStore();
await authStore.initAuth();

app.mount("#app");
```

`initAuth()` 方法执行以下操作：

1. 检查 localStorage 中是否有 Token
2. 如果有 Token，调用 API 获取用户信息
3. 更新 `currentUser` 状态

### 登录流程

```typescript
async login(credentials: LoginRequest) {
  this.isLoading = true
  this.error = null

  try {
    // 1. 调用登录 API
    const result = await AuthAPI.login(credentials)

    // 2. 检查是否需要 2FA
    if (result.requiresTwoFactor) {
      return result  // 返回给组件处理 2FA
    }

    // 3. 保存 Token
    saveAccessToken(result.access_token)
    saveRefreshToken(result.refresh_token)

    // 4. 更新用户状态
    this.currentUser = result.user

    return { success: true }
  } catch (error) {
    this.error = formatAuthError(error)
    return { success: false, message: this.error }
  } finally {
    this.isLoading = false
  }
}
```

### 登出流程

```typescript
async logout() {
  // 1. 清除 Token
  clearAuthTokens()

  // 2. 清除用户状态
  this.currentUser = null

  // 3. 跳转到登录页
  router.push('/auth/login')
}
```

## 在组件中使用

### Options API

```vue
<script>
import { useAuthStore } from "@/stores/auth";

export default {
  computed: {
    authStore() {
      return useAuthStore();
    },
    isLoggedIn() {
      return this.authStore.isAuthenticated;
    },
  },
};
</script>
```

### Composition API

```vue
<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { storeToRefs } from "pinia";

const authStore = useAuthStore();

// 解构响应式状态
const { currentUser, isLoading, error } = storeToRefs(authStore);

// 直接使用 actions
const handleLogout = () => {
  authStore.logout();
};
</script>
```

## Token 存储

Token 存储在 localStorage 中，通过工具函数管理：

```typescript
// src/utils/auth/storage.ts

// 存储键名
const ACCESS_TOKEN_KEY = "access_token";
const REFRESH_TOKEN_KEY = "refresh_token";
const TOKEN_EXPIRY_KEY = "token_expiry";

// 保存访问令牌
export function saveAccessToken(token: string) {
  localStorage.setItem(ACCESS_TOKEN_KEY, token);
}

// 获取访问令牌
export function getAccessToken(): string | null {
  return localStorage.getItem(ACCESS_TOKEN_KEY);
}

// 清除所有 Token
export function clearAuthTokens() {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
  localStorage.removeItem(TOKEN_EXPIRY_KEY);
}

// 检查是否有 Token
export function hasAccessToken(): boolean {
  return !!getAccessToken();
}
```

## 最佳实践

### 1. 状态最小化

只在 Store 中存储真正需要跨组件共享的状态，组件内部状态使用 `ref`/`reactive`。

### 2. 使用 Getters

计算属性应该放在 getters 中，而不是在组件中重复计算。

### 3. Action 异步处理

所有异步操作（如 API 调用）都应该在 actions 中进行。

### 4. 错误处理

统一在 Store 中处理错误，组件只需要读取 `error` 状态。

```typescript
// Store 中
this.error = formatAuthError(error)

// 组件中
<v-alert v-if="error" type="error">{{ error }}</v-alert>
```
