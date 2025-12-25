# 前端架构概述

本项目前端采用 Vue 3 + TypeScript + Vuetify 技术栈，遵循模块化设计原则，与后端 DDD 架构保持一致的领域划分。

<!--TOC-->

## Table of Contents

- [技术栈](#技术栈) `:24+13`
- [项目结构](#项目结构) `:37+36`
- [架构分层](#架构分层) `:73+28`
- [模块划分](#模块划分) `:101+22`
  - [1. Auth 模块 (认证)](#1-auth-模块-认证) `:105+6`
  - [2. Admin 模块 (管理后台)](#2-admin-模块-管理后台) `:111+6`
  - [3. User 模块 (用户中心)](#3-user-模块-用户中心) `:117+6`
- [开发端口](#开发端口) `:123+7`
- [开发工具链](#开发工具链) `:130+31`
  - [自动导入](#自动导入) `:143+18`
- [快速开始](#快速开始) `:161+28`
- [下一步](#下一步) `:189+7`

<!--TOC-->

## 技术栈

| 类别        | 技术       | 版本  | 说明                   |
| ----------- | ---------- | ----- | ---------------------- |
| 框架        | Vue        | ^3.5  | 响应式前端框架         |
| 构建工具    | Vite       | ^7.2  | 下一代前端构建工具     |
| 语言        | TypeScript | ~5.9  | 类型安全的 JavaScript  |
| UI 组件库   | Vuetify    | ^3.10 | Material Design 组件库 |
| 状态管理    | Pinia      | ^3.0  | Vue 官方状态管理       |
| 路由        | Vue Router | ^4.6  | Vue 官方路由           |
| HTTP 客户端 | Axios      | ^1.13 | Promise 风格 HTTP 库   |
| 图标        | MDI        | ^7.4  | Material Design Icons  |

## 项目结构

```
src/
├── api/                    # API 请求层
│   ├── admin/              # 管理员 API
│   ├── auth/               # 认证 API（含 Axios 客户端）
│   ├── user/               # 用户 API
│   └── helpers/            # 工具函数
├── types/                  # TypeScript 类型定义
│   ├── admin/              # 管理端类型
│   ├── auth/               # 认证类型
│   ├── user/               # 用户类型
│   ├── response/           # 响应类型
│   └── common/             # 通用类型
├── composables/            # 组合式函数（自动导入）
├── components/             # 共享组件（自动注册）
├── stores/                 # Pinia 状态管理
├── router/                 # Vue Router 路由
├── utils/                  # 工具函数
├── layout/                 # 布局组件
├── views/                  # 共享视图组件
├── pages/                  # 页面组件
│   ├── admin/              # 管理后台页面
│   ├── auth/               # 认证页面
│   └── user/               # 用户中心页面
├── App.vue                 # 根组件
└── main.ts                 # 应用入口

# 根目录配置文件
├── public/                 # 静态资源
├── index.html              # 入口 HTML
├── vite.config.ts          # Vite 配置
└── tsconfig.json           # TypeScript 配置
```

## 架构分层

前端采用清晰的分层架构，各层职责明确：

```
┌─────────────────────────────────────────────────────────────┐
│                      Pages (页面组件)                         │
│  - 负责页面布局和 UI 展示                                      │
│  - 调用 Composables 处理业务逻辑                               │
├─────────────────────────────────────────────────────────────┤
│                  Composables (组合式函数)                      │
│  - 封装页面级业务逻辑                                          │
│  - 管理组件状态和副作用                                        │
├─────────────────────────────────────────────────────────────┤
│                    Stores (状态管理)                          │
│  - 全局状态管理 (Pinia)                                       │
│  - 跨组件状态共享                                             │
├─────────────────────────────────────────────────────────────┤
│                      API (请求层)                             │
│  - 封装 HTTP 请求                                            │
│  - 处理认证和错误                                             │
├─────────────────────────────────────────────────────────────┤
│                    Types (类型定义)                           │
│  - TypeScript 类型和接口                                      │
│  - 请求/响应数据结构                                          │
└─────────────────────────────────────────────────────────────┘
```

## 模块划分

前端按业务领域划分为三大模块，与后端保持一致：

### 1. Auth 模块 (认证)

- **页面**: 登录、注册
- **功能**: 用户认证、验证码、双因素认证
- **路由前缀**: `/auth/*`

### 2. Admin 模块 (管理后台)

- **页面**: 概览、用户管理、角色管理、菜单管理、系统设置、审计日志
- **功能**: 系统管理、权限配置、日志查看
- **路由前缀**: `/admin/*`

### 3. User 模块 (用户中心)

- **页面**: 个人资料、安全设置、访问令牌
- **功能**: 个人信息管理、密码修改、2FA 设置、PAT 管理
- **路由前缀**: `/user/*`

## 开发端口

| 服务           | 端口  | 说明            |
| -------------- | ----- | --------------- |
| 前端开发服务器 | 40013 | Vite dev server |
| 后端 API       | 40012 | Go HTTP server  |

## 开发工具链

项目配置了完整的代码质量工具链：

| 工具                     | 用途             | 命令                 |
| ------------------------ | ---------------- | -------------------- |
| **ESLint 9.x**           | 代码检查         | `npm run lint`       |
| **Prettier**             | 代码格式化       | `npm run format`     |
| **Vitest**               | 单元测试         | `npm run test`       |
| **vue-tsc**              | 类型检查         | `npm run type-check` |
| **vite-plugin-vuetify**  | Vuetify 按需导入 | 自动                 |
| **unplugin-auto-import** | Vue API 自动导入 | 自动                 |

### 自动导入

项目配置了自动导入功能，无需手动导入常用 API：

```vue
<script setup lang="ts">
// 无需手动导入，以下 API 自动可用：
// - Vue: ref, computed, watch, onMounted, etc.
// - Vue Router: useRouter, useRoute
// - Pinia: defineStore, storeToRefs
// - 所有 composables

const count = ref(0);
const router = useRouter();
const { copy } = useClipboard();
</script>
```

## 快速开始

```bash
# 进入前端目录
cd web

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 代码检查
npm run lint

# 代码格式化
npm run format

# 运行测试
npm run test

# 类型检查
npm run type-check

# 生产构建
npm run build
```

## 下一步

- [代码规范](./frontend-code-style) - 了解 ESLint、Prettier、测试规范
- [路由配置](./frontend-router) - 了解路由系统和守卫
- [状态管理](./frontend-state) - 了解 Pinia 状态管理
- [API 层设计](./frontend-api) - 了解 API 请求封装
- [组件开发](./frontend-components) - 了解页面和组件结构
