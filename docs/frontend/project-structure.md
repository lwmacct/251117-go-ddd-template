# 项目结构

本文档详细介绍前端项目的目录组织、文件命名规范和模块划分。

## 目录结构

```
web/
├── public/                 # 静态资源（不经过 Vite 处理）
│   └── favicon.ico            # 网站图标
│
├── src/                    # 源代码目录
│   ├── api/                   # API 接口封装
│   │   ├── client.ts             # Axios 客户端配置
│   │   ├── auth.ts               # 认证相关接口
│   │   └── users.ts              # 用户相关接口
│   │
│   ├── components/            # 通用组件（可复用）
│   │   ├── common/               # 基础组件
│   │   └── business/             # 业务组件
│   │
│   ├── global/                # 全局配置和样式
│   │   ├── styles/               # 全局样式
│   │   └── plugins/              # 插件配置
│   │
│   ├── layout/                # 布局组件
│   │   ├── DefaultLayout.vue     # 默认布局
│   │   └── AuthLayout.vue        # 认证布局
│   │
│   ├── pages/                 # 页面组件
│   │   ├── Home.vue              # 首页
│   │   ├── Login.vue             # 登录页
│   │   └── Dashboard.vue         # 仪表板
│   │
│   ├── router/                # 路由配置
│   │   └── index.ts              # 路由定义
│   │
│   ├── stores/                # Pinia 状态管理
│   │   ├── auth.ts               # 认证 Store
│   │   └── user.ts               # 用户 Store
│   │
│   ├── types/                 # TypeScript 类型定义
│   │   ├── api.ts                # API 类型
│   │   └── models.ts             # 数据模型
│   │
│   ├── utils/                 # 工具函数
│   │   ├── request.ts            # 请求工具
│   │   └── storage.ts            # 存储工具
│   │
│   ├── views/                 # 视图组件
│   │   └── (按功能模块组织)
│   │
│   ├── App.vue                # 根组件
│   └── main.ts                # 应用入口
│
├── dist/                   # 构建输出目录
├── node_modules/           # 依赖包
│
├── index.html              # HTML 入口模板
├── vite.config.ts          # Vite 配置
├── tsconfig.json           # TypeScript 配置
├── tsconfig.app.json       # 应用 TS 配置
├── tsconfig.node.json      # Node TS 配置
├── package.json            # 项目配置
└── README.md               # 项目说明
```

## 核心目录说明

### `src/api/` - API 接口

**职责**: 封装所有后端 API 调用

**结构**:
```
api/
├── client.ts       # Axios 客户端配置、拦截器
├── auth.ts         # 认证接口：登录、注册、刷新 Token
├── users.ts        # 用户接口：CRUD、角色管理
└── index.ts        # 导出所有 API
```

**示例**:
```typescript
// api/users.ts
import client from './client'

export const userApi = {
  getProfile: () => client.get('/api/user/me'),
  updateProfile: (data) => client.put('/api/user/me', data)
}
```

### `src/components/` - 通用组件

**职责**: 可复用的 UI 组件

**分类**:
- `common/` - 基础组件（按钮、输入框、卡片）
- `business/` - 业务组件（用户卡片、数据表格）

**命名规范**:
```
PascalCase.vue

✓ UserCard.vue
✓ DataTable.vue
✗ userCard.vue
✗ data-table.vue
```

### `src/pages/` - 页面组件

**职责**: 路由对应的页面组件

**特点**:
- 每个页面对应一个路由
- 组合通用组件构建页面
- 处理页面级别的状态

**示例**:
```vue
<!-- pages/Dashboard.vue -->
<template>
  <v-container>
    <h1>Dashboard</h1>
    <UserCard :user="currentUser" />
  </v-container>
</template>
```

### `src/router/` - 路由配置

**职责**: 管理应用路由

**核心文件**: `index.ts`

```typescript
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('@/pages/Home.vue')
    },
    {
      path: '/login',
      component: () => import('@/pages/Login.vue')
    }
  ]
})

export default router
```

### `src/stores/` - 状态管理

**职责**: Pinia Store，管理全局状态

**结构**:
```
stores/
├── auth.ts         # 认证状态：token、用户信息
├── user.ts         # 用户状态：个人资料
└── index.ts        # Store 导出
```

**示例**:
```typescript
// stores/auth.ts
import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: null,
    user: null
  }),
  actions: {
    async login(credentials) {
      // 登录逻辑
    }
  }
})
```

### `src/types/` - 类型定义

**职责**: TypeScript 类型和接口定义

**结构**:
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
  id: number
  username: string
  email: string
}

export interface LoginRequest {
  login: string
  password: string
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
import { userApi } from '@/api/users'
import UserCard from '@/components/UserCard.vue'

// 等同于
import { userApi } from '../api/users'
import UserCard from '../components/UserCard.vue'
```

### 推荐的导入顺序

```typescript
// 1. Vue 核心
import { ref, computed } from 'vue'

// 2. Vue 生态库
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// 3. 第三方库
import axios from 'axios'

// 4. 项目内部
import { userApi } from '@/api/users'
import UserCard from '@/components/UserCard.vue'

// 5. 类型
import type { User } from '@/types/models'
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
  list: (params) => client.get('/api/users', { params }),
  get: (id) => client.get(`/api/users/${id}`),
  create: (data) => client.post('/api/users', data)
}

// ✗ 避免：直接调用
axios.get('/api/users')
```

### 类型定义

```typescript
// ✓ 推荐：定义类型
interface User {
  id: number
  name: string
}

const users = ref<User[]>([])

// ✗ 避免：any
const users = ref<any>([])
```

## 扩展阅读

- [开发规范](./coding-standards.md) - 代码风格指南
- [组件规范](./components.md) - 组件设计原则
- [API 集成](./api-integration.md) - API 使用指南

熟悉项目结构后，开始愉快地编码吧！ 🎨
