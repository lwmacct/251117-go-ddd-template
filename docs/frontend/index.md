# 前端文档

欢迎来到 Go DDD Template 前端文档。本文档介绍基于 Vue 3 的现代化前端应用。

## 技术栈

### 核心框架

- **[Vue 3.5](https://vuejs.org/)** - 渐进式 JavaScript 框架
- **[TypeScript 5.9](https://www.typescriptlang.org/)** - 类型安全的 JavaScript 超集
- **[Vite 7](https://vitejs.dev/)** - 下一代前端构建工具
- **[Vuetify 3](https://vuetifyjs.com/)** - Material Design 组件库

### 状态管理 & 路由

- **[Pinia 3](https://pinia.vuejs.org/)** - Vue 官方状态管理库
- **[Vue Router 4](https://router.vuejs.org/)** - Vue 官方路由管理器

### 工具库

- **[Axios](https://axios-http.com/)** - HTTP 客户端
- **[MDI Font](https://pictogrammers.com/library/mdi/)** - Material Design Icons

## 项目结构

```
web/
├── public/              # 静态资源
│   └── favicon.ico
├── src/
│   ├── api/            # API 接口封装
│   ├── global/         # 全局配置和样式
│   ├── layout/         # 布局组件
│   ├── pages/          # 页面组件
│   ├── router/         # 路由配置
│   ├── stores/         # Pinia 状态管理
│   ├── types/          # TypeScript 类型定义
│   ├── utils/          # 工具函数
│   ├── views/          # 视图组件
│   ├── App.vue         # 根组件
│   └── main.ts         # 应用入口
├── index.html          # HTML 模板
├── vite.config.ts      # Vite 配置
└── package.json        # 项目配置
```

## 快速开始

### 环境要求

- **Node.js**: ^20.19.0 || >=22.12.0
- **npm**: 最新版本

### 安装依赖

```bash
cd web
npm install
```

### 开发模式

```bash
npm run dev
```

应用将在 `http://localhost:5173` 启动。

### 生产构建

```bash
npm run build
```

构建产物将输出到 `dist/` 目录。

### 类型检查

```bash
npm run type-check
```

## 文档导航

### 基础

- **[快速开始](./getting-started)** - 环境搭建、项目启动
- **[项目结构](./project-structure)** - 目录组织、文件命名规范
- **[API 集成](./api-integration)** - Axios 封装、请求拦截器

<!-- TODO: 待完善的文档
- **[开发规范](./coding-standards)** - 代码风格、最佳实践

### 核心概念

- **[路由管理](./routing)** - Vue Router 配置、路由守卫
- **[状态管理](./state-management)** - Pinia Store 使用指南

### 组件开发

- **[组件规范](./components)** - 组件设计原则、最佳实践
- **[Vuetify 使用](./vuetify)** - Material Design 组件库
- **[布局系统](./layouts)** - 应用布局、响应式设计

### 进阶主题

- **[认证授权](./authentication)** - JWT 集成、权限控制
- **[类型系统](./typescript)** - TypeScript 类型定义
- **[性能优化](./performance)** - 构建优化、懒加载

### 部署上线

- **[构建部署](./deployment)** - 生产构建、部署策略
- **[环境配置](./environment)** - 环境变量、多环境配置
-->

## 特性亮点

### ✅ 现代化技术栈

- Vue 3 Composition API
- TypeScript 类型安全
- Vite 极速构建
- Pinia 轻量级状态管理

### ✅ Material Design

- Vuetify 3 组件库
- 响应式布局
- 主题定制
- 丰富的图标库

### ✅ 开发体验

- 热模块替换 (HMR)
- TypeScript 智能提示
- Vue DevTools 支持
- 快速刷新

### ✅ 生产就绪

- 代码分割
- Tree Shaking
- 压缩优化
- 浏览器兼容

## 开发工作流

### 1. 创建新页面

```bash
# 创建页面组件
src/pages/MyPage.vue

# 配置路由
src/router/index.ts
```

### 2. API 集成

```bash
# 定义 API 接口
src/api/my-api.ts

# 创建 Store
src/stores/my-store.ts
```

### 3. 组件开发

```bash
# 创建组件
src/components/MyComponent.vue

# 使用组件
<template>
  <MyComponent />
</template>
```

## 配置文件

### Vite 配置 (`vite.config.ts`)

```typescript
import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
```

### TypeScript 配置 (`tsconfig.json`)

```json
{
  "extends": "@vue/tsconfig/tsconfig.dom.json",
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  }
}
```

## 常用命令

| 命令 | 说明 |
|------|------|
| `npm run dev` | 启动开发服务器 |
| `npm run build` | 生产构建 |
| `npm run preview` | 预览构建产物 |
| `npm run type-check` | TypeScript 类型检查 |

## 相关资源

### 官方文档

- [Vue 3 文档](https://vuejs.org/)
- [Vite 文档](https://vitejs.dev/)
- [Vuetify 文档](https://vuetifyjs.com/)
- [Pinia 文档](https://pinia.vuejs.org/)
- [Vue Router 文档](https://router.vuejs.org/)

### 后端集成

- [API 文档](/api/) - 后端 API 接口文档
- [认证授权](/architecture/authentication) - 认证机制说明
- [RBAC 权限](/architecture/rbac) - 权限系统详解

## 贡献指南

欢迎为前端项目贡献代码！请遵循以下规范：

1. **代码风格**: 遵循 Vue 3 和 TypeScript 最佳实践
2. **组件设计**: 单一职责、可复用、可测试
3. **类型定义**: 为所有函数和组件添加类型
4. **文档注释**: 为复杂逻辑添加注释

<!-- 查看 [开发规范](./coding-standards) 了解详细信息。 -->

## 常见问题

### Q: 如何添加新的 API 接口？

**A**: 在 `src/api/` 目录下创建接口文件，使用 Axios 封装。详见 [API 集成](./api-integration)。

<!-- TODO: 待完善的文档
### Q: 如何管理应用状态？

**A**: 使用 Pinia Store。详见 [状态管理](./state-management)。

### Q: 如何配置路由？

**A**: 在 `src/router/index.ts` 中配置。详见 [路由管理](./routing)。

### Q: 如何使用 Vuetify 组件？

**A**: 导入并使用 Vuetify 提供的组件。详见 [Vuetify 使用](./vuetify)。
-->

## 下一步

- 阅读 [快速开始](./getting-started) 开始开发
- 了解 [项目结构](./project-structure) 熟悉代码组织
- 学习 [API 集成](./api-integration) 了解后端对接

开始构建出色的前端应用吧！ 🚀
