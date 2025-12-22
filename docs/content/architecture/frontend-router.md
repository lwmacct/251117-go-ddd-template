# 路由配置

前端使用 Vue Router 进行路由管理，采用 Hash 模式，支持路由守卫和权限控制。

<!--TOC-->

## Table of Contents

- [路由结构](#路由结构) `:26+10`
- [路由配置](#路由配置-1) `:36+58`
  - [认证路由 (auth.ts)](#认证路由-authts) `:38+21`
  - [管理后台路由 (admin.ts)](#管理后台路由-admints) `:59+19`
  - [用户中心路由 (user.ts)](#用户中心路由-userts) `:78+16`
- [路由守卫](#路由守卫) `:94+30`
- [路由元信息 (Meta)](#路由元信息-meta) `:124+8`
- [路由跳转](#路由跳转) `:132+27`
  - [编程式导航](#编程式导航) `:134+17`
  - [声明式导航](#声明式导航) `:151+8`
- [布局组件](#布局组件) `:159+14`
  - [AdminLayout](#adminlayout) `:161+8`
  - [UserLayout](#userlayout) `:169+4`
- [路由懒加载](#路由懒加载) `:173+10`

<!--TOC-->

## 路由结构

```
src/router/
├── index.ts       # 主路由入口 + 路由守卫
├── auth.ts        # 认证路由
├── admin.ts       # 管理后台路由
└── user.ts        # 用户中心路由
```

## 路由配置

### 认证路由 (`auth.ts`)

```typescript
// 公开路由，无需登录
export const authRoutes = {
  path: "/auth",
  children: [
    {
      path: "login",
      component: () => import("@/pages/auth/login/index.vue"),
      meta: { requiresAuth: false },
    },
    {
      path: "register",
      component: () => import("@/pages/auth/register/index.vue"),
      meta: { requiresAuth: false },
    },
  ],
};
```

### 管理后台路由 (`admin.ts`)

```typescript
// 需要登录，使用 AdminLayout 布局
export const adminRoutes = {
  path: "/admin",
  component: () => import("@/layout/AdminLayout.vue"),
  meta: { requiresAuth: true },
  children: [
    { path: "overview", name: "数据概览", icon: "mdi-speedometer" },
    { path: "users", name: "用户管理", icon: "mdi-account" },
    { path: "roles", name: "角色管理", icon: "mdi-account-group" },
    { path: "menus", name: "菜单管理", icon: "mdi-menu" },
    { path: "settings", name: "系统设置", icon: "mdi-cog" },
    { path: "auditlogs", name: "审计日志", icon: "mdi-file-document-outline" },
  ],
};
```

### 用户中心路由 (`user.ts`)

```typescript
// 需要登录，使用 UserLayout 布局
export const userRoutes = {
  path: "/user",
  component: () => import("@/layout/UserLayout.vue"),
  meta: { requiresAuth: true },
  children: [
    { path: "profile", name: "个人资料", icon: "mdi-account-circle" },
    { path: "security", name: "安全设置", icon: "mdi-shield-lock" },
    { path: "tokens", name: "访问令牌", icon: "mdi-key-variant" },
  ],
};
```

## 路由守卫

路由守卫在 `router/index.ts` 中实现，处理认证逻辑：

```typescript
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();

  // 检查路由是否需要认证
  if (to.meta.requiresAuth) {
    if (!authStore.hasToken) {
      // 未登录，重定向到登录页，保存目标路由
      next({
        path: "/auth/login",
        query: { redirect: to.fullPath },
      });
      return;
    }
  }

  // 已登录用户访问认证页面，重定向到管理后台
  if (to.path.startsWith("/auth") && authStore.hasToken) {
    next("/admin/overview");
    return;
  }

  next();
});
```

## 路由元信息 (Meta)

| 字段           | 类型      | 说明         |
| -------------- | --------- | ------------ |
| `requiresAuth` | `boolean` | 是否需要登录 |
| `title`        | `string`  | 页面标题     |
| `icon`         | `string`  | MDI 图标名称 |

## 路由跳转

### 编程式导航

```typescript
import { useRouter } from "vue-router";

const router = useRouter();

// 跳转到指定路由
router.push("/admin/users");

// 带参数跳转
router.push({ path: "/admin/users", query: { page: 1 } });

// 返回上一页
router.back();
```

### 声明式导航

```vue
<template>
  <router-link to="/admin/users">用户管理</router-link>
</template>
```

## 布局组件

### AdminLayout

管理后台布局，包含：

- 顶部导航栏 (AppBars)
- 左侧菜单 (Navigation)
- 主内容区 (router-view)

### UserLayout

用户中心布局，结构与 AdminLayout 类似，菜单项不同。

## 路由懒加载

所有页面组件都使用动态导入实现懒加载：

```typescript
// 懒加载语法
component: () => import("@/pages/admin/users/index.vue");
```

这样可以实现代码分割，按需加载页面，提升首屏加载速度。
