---
paths:
  - "src/layout/**/*.vue"
---

# Layout 规范

## 核心概念

Layout 定义页面的**整体框架结构**，通过 `<router-view>` 渲染子页面。

```
src/layout/
├── BaseLayout.vue    # 基础布局（AppBars + Navigation + router-view）
├── AdminLayout.vue   # 继承 BaseLayout，配置 admin 菜单
├── UserLayout.vue    # 继承 BaseLayout，配置 user 菜单
└── AuthLayout.vue    # 独立布局（左右分屏）
```

## 两种模式

### 继承模式（推荐）

```vue
<!-- AdminLayout.vue -->
<script setup lang="ts">
import BaseLayout from "./BaseLayout.vue";
import { useMenus } from "@/composables/useMenus";

const { adminMenus } = useMenus();
</script>

<template>
  <BaseLayout :menu-items="adminMenus" />
</template>
```

### 独立模式

完全独立的布局结构（如 AuthLayout 的左右分屏设计）。

## 职责划分

| 组件          | 职责                   |
| ------------- | ---------------------- |
| `*Layout`     | 定义框架结构、配置菜单 |
| `router-view` | 渲染子页面内容         |
| `@/views/`    | 布局内的可复用视图组件 |

## 禁止事项

- ❌ Layout 中不包含业务逻辑
- ❌ Layout 中不直接调用业务 API
