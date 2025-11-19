# 快速开始

本指南帮助你快速搭建开发环境并运行前端应用。

## 环境要求

### Node.js 版本

本项目要求以下 Node.js 版本之一：

- **Node.js 20.19.0** 或更高版本（20.x 系列）
- **Node.js 22.12.0** 或更高版本（22.x 系列）

### 检查当前版本

```bash
node --version
# v22.12.0 (推荐) 或 v20.19.0+
```

### 安装 Node.js

如果需要安装或升级 Node.js：

**使用 nvm (推荐)**:
```bash
# 安装 nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

# 安装 Node.js 22
nvm install 22
nvm use 22

# 验证
node --version
```

**直接下载**:
访问 [Node.js 官网](https://nodejs.org/) 下载 LTS 版本。

## 项目设置

### 1. 进入前端目录

```bash
cd web
```

### 2. 安装依赖

```bash
npm install
```

**安装过程说明**:
- Vue 3.5 - 渐进式框架
- Vuetify 3 - Material Design 组件库
- Pinia - 状态管理
- Vue Router - 路由管理
- Axios - HTTP 客户端
- TypeScript - 类型检查
- Vite - 构建工具

### 3. 启动开发服务器

```bash
npm run dev
```

**输出示例**:
```
VITE v7.2.2  ready in 234 ms

➜  Local:   http://localhost:5173/
➜  Network: http://192.168.1.100:5173/
➜  press h + enter to show help
```

### 4. 访问应用

打开浏览器访问 `http://localhost:5173`

## 开发工作流

### 热模块替换 (HMR)

修改代码后，浏览器会自动刷新，无需手动重启服务器。

**示例**:
```vue
<!-- src/App.vue -->
<template>
  <div>
    <h1>Hello World</h1>  <!-- 修改后自动刷新 -->
  </div>
</template>
```

### TypeScript 类型检查

在开发过程中运行类型检查：

```bash
npm run type-check
```

### Vue DevTools

安装 [Vue DevTools](https://devtools.vuejs.org/) 浏览器扩展：

- Chrome: [Chrome Web Store](https://chrome.google.com/webstore/detail/vuejs-devtools/nhdogjmejiglipccpnnnanhbledajbpd)
- Firefox: [Firefox Add-ons](https://addons.mozilla.org/en-US/firefox/addon/vue-js-devtools/)

## 生产构建

### 构建应用

```bash
npm run build
```

**构建过程**:
1. TypeScript 类型检查
2. Vite 生产构建
3. 代码压缩和优化
4. 输出到 `dist/` 目录

**输出示例**:
```
vite v7.2.2 building for production...
✓ 234 modules transformed.
dist/index.html                   0.46 kB │ gzip:  0.30 kB
dist/assets/index-BfR5xN2K.css   156.78 kB │ gzip: 25.32 kB
dist/assets/index-C8fHLOt_.js    502.45 kB │ gzip: 165.23 kB
✓ built in 3.45s
```

### 预览构建产物

```bash
npm run preview
```

访问 `http://localhost:4173` 预览生产版本。

## 项目配置

### Vite 配置

**文件**: `vite.config.ts`

```typescript
import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),  // Vue DevTools 插件
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))  // @ 别名
    }
  },
  server: {
    port: 5173,
    host: true  // 允许外部访问
  }
})
```

### TypeScript 配置

**文件**: `tsconfig.json`

```json
{
  "extends": "@vue/tsconfig/tsconfig.dom.json",
  "include": ["env.d.ts", "src/**/*", "src/**/*.vue"],
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]  // @ 别名配置
    }
  }
}
```

## 开发工具推荐

### VS Code 扩展

**必备扩展**:
- [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) - Vue 3 语言支持
- [TypeScript Vue Plugin](https://marketplace.visualstudio.com/items?itemName=Vue.vscode-typescript-vue-plugin) - TS 支持

**推荐扩展**:
- [ESLint](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint) - 代码检查
- [Prettier](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode) - 代码格式化
- [Vue VSCode Snippets](https://marketplace.visualstudio.com/items?itemName=sdras.vue-vscode-snippets) - Vue 代码片段

### VS Code 设置

**`.vscode/settings.json`**:
```json
{
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  },
  "[vue]": {
    "editor.defaultFormatter": "Vue.volar"
  },
  "[typescript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  }
}
```

## 常用命令

| 命令 | 说明 | 使用场景 |
|------|------|---------|
| `npm install` | 安装依赖 | 首次克隆或依赖更新后 |
| `npm run dev` | 开发服务器 | 日常开发 |
| `npm run build` | 生产构建 | 部署前 |
| `npm run preview` | 预览构建 | 验证生产版本 |
| `npm run type-check` | 类型检查 | 提交代码前 |

## 目录说明

### 快速导航

```
web/
├── src/
│   ├── api/           # API 接口封装
│   ├── components/    # 通用组件
│   ├── pages/         # 页面组件
│   ├── router/        # 路由配置
│   ├── stores/        # Pinia Store
│   ├── types/         # TypeScript 类型
│   ├── utils/         # 工具函数
│   └── main.ts        # 应用入口
├── public/            # 静态资源
└── dist/              # 构建输出
```

详见 [项目结构](./project-structure.md)。

## 后端集成

### 启动后端服务

前端需要连接到后端 API，确保后端服务已启动：

```bash
# 在项目根目录
docker-compose up -d      # 启动数据库和 Redis
task go:run -- api        # 启动后端 API 服务
```

后端将在 `http://localhost:8080` 启动。

### API 配置

前端默认连接到 `http://localhost:8080`。

**修改 API 地址**:
```typescript
// src/api/client.ts
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
```

## 故障排查

### 端口被占用

**错误**:
```
Error: listen EADDRINUSE: address already in use :::5173
```

**解决方案**:
```bash
# 方案 1: 杀掉占用进程
lsof -ti:5173 | xargs kill -9

# 方案 2: 修改端口
# vite.config.ts
server: {
  port: 5174  // 使用其他端口
}
```

### 依赖安装失败

**错误**:
```
npm ERR! code ERESOLVE
```

**解决方案**:
```bash
# 清除缓存
rm -rf node_modules package-lock.json
npm cache clean --force

# 重新安装
npm install
```

### 类型错误

**错误**:
```
Property 'xxx' does not exist on type 'xxx'
```

**解决方案**:
```bash
# 重新生成类型声明
npm run type-check

# 重启 VS Code 的 TS 服务器
# VS Code: Ctrl+Shift+P → "TypeScript: Restart TS Server"
```

## 下一步

- 了解 [项目结构](./project-structure.md)
- 学习 [开发规范](./coding-standards.md)
- 开始 [API 集成](./api-integration.md)
- 探索 [Vuetify 组件](./vuetify.md)

开始你的前端开发之旅吧！ 🎉
