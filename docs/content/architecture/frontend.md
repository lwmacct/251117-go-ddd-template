# 前端架构

Vue 3 + TypeScript + Vuetify 技术栈，与后端 DDD 架构保持一致的领域划分。

<!--TOC-->

## Table of Contents

- [技术栈](#技术栈) `:21+10`
- [项目结构](#项目结构) `:31+16`
- [模块划分](#模块划分) `:47+8`
- [路由系统](#路由系统) `:55+18`
- [状态管理](#状态管理) `:73+12`
- [API 层](#api-层) `:85+18`
- [组件开发](#组件开发) `:103+22`
- [代码规范](#代码规范) `:125+19`
- [开发命令](#开发命令) `:144+8`

<!--TOC-->

## 技术栈

| 类别 | 技术       | 说明            |
| ---- | ---------- | --------------- |
| 框架 | Vue 3      | 响应式前端框架  |
| 构建 | Vite       | 下一代构建工具  |
| UI   | Vuetify 3  | Material Design |
| 状态 | Pinia      | 状态管理        |
| 路由 | Vue Router | Hash 模式       |

## 项目结构

```
src/
├── api/            # API 请求层
├── stores/         # Pinia 状态
├── router/         # 路由配置
├── pages/          # 页面组件
│   ├── admin/      # 管理后台
│   ├── auth/       # 认证页面
│   └── user/       # 用户中心
├── composables/    # 组合式函数
├── components/     # 共享组件
└── layouts/        # 布局组件
```

## 模块划分

| 模块  | 路由       | 功能                 |
| ----- | ---------- | -------------------- |
| Auth  | `/auth/*`  | 登录、注册、2FA      |
| Admin | `/admin/*` | 用户、角色、菜单管理 |
| User  | `/user/*`  | 个人资料、安全设置   |

## 路由系统

**配置文件**: `src/router/index.ts`

| 元信息                | 说明     |
| --------------------- | -------- |
| `requiresAuth: true`  | 需要登录 |
| `requiresAuth: false` | 公开页面 |

**路由守卫**: 在 `router/index.ts` 中实现认证逻辑

**布局组件**:

| 组件        | 位置                           | 用途         |
| ----------- | ------------------------------ | ------------ |
| AdminLayout | `src/layouts/AdminLayout.vue`  | 管理后台布局 |
| Sidebar     | `src/views/common/Sidebar.vue` | 导航菜单     |

## 状态管理

核心 Store: `src/stores/auth.ts`

| 状态              | 说明         |
| ----------------- | ------------ |
| `currentUser`     | 当前用户信息 |
| `isAuthenticated` | 是否已登录   |
| `hasToken`        | 是否有 Token |

**Token 存储**: `src/utils/auth/storage.ts`（localStorage）

## API 层

| 目录             | 用途       |
| ---------------- | ---------- |
| `src/api/auth/`  | 认证 API   |
| `src/api/admin/` | 管理端 API |
| `src/api/user/`  | 用户端 API |

**核心机制**:

| 机制           | 位置                   | 说明                 |
| -------------- | ---------------------- | -------------------- |
| Token 自动附加 | `client.ts` 请求拦截器 | 从 localStorage 读取 |
| Token 自动刷新 | `client.ts` 响应拦截器 | 401 时刷新后重试     |
| 错误统一处理   | `errors.ts`            | 格式化错误消息       |

> 类型使用规范见 `.claude/rules/frontend-api.md`

## 组件开发

每个页面包含独立的组件、Composables 和类型定义：

```
src/pages/{module}/{feature}/
├── index.vue           # 页面入口
├── components/         # 页面私有组件
├── composables/        # 页面逻辑
└── types.ts           # 页面类型（可选）
```

**Composable 模式**:

| 类型   | 命名           | 返回                             |
| ------ | -------------- | -------------------------------- |
| 列表   | `useXxxList`   | `{ items, loading, pagination }` |
| 表单   | `useXxxForm`   | `{ form, submit, reset }`        |
| 对话框 | `useXxxDialog` | `{ open, close, visible }`       |

**参考实现**: `src/pages/admin/users/composables/useUserList.ts`

## 代码规范

| 工具       | 用途       |
| ---------- | ---------- |
| ESLint     | 代码检查   |
| Prettier   | 代码格式化 |
| TypeScript | 类型检查   |

**文件命名**:

- 组件: `PascalCase.vue`
- Composables: `use*.ts`
- 工具函数: `camelCase.ts`

**TypeScript 规范**:

- 禁止重复定义后端 DTO
- 前端派生类型放在 `src/api/types.ts`

## 开发命令

```bash
pnpm dev          # 启动开发服务器 (40013)
pnpm build        # 生产构建
pnpm lint         # ESLint 检查
pnpm vue-tsc      # 类型检查
```
