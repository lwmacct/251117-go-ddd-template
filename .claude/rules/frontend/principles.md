---
paths:
  - "src/**/*.vue"
  - "src/**/*.ts"
  - "src/**/*.css"
---

# 前端核心原则

## Vuetify 主题系统

所有样式应使用 Vuetify 主题系统，确保亮/暗模式自适应。

### 方式一：CSS 工具类（推荐）

Vuetify 自动生成主题感知的 CSS 类，优先使用：

```vue
<!-- 背景色 -->
<div class="bg-primary">主色背景</div>
<div class="bg-surface">表面背景</div>
<div class="bg-surface-light">浅表面背景</div>

<!-- 文字颜色 -->
<div class="text-primary">主色文字</div>
<div class="text-high-emphasis">高强调文字 (87%)</div>
<div class="text-medium-emphasis">中强调文字 (60%)</div>
<div class="text-disabled">禁用文字 (38%)</div>

<!-- 边框颜色 -->
<div class="border-primary">主色边框</div>
```

### 方式二：CSS 变量

自定义样式时使用 Vuetify CSS 变量：

```css
.custom-element {
  /* 颜色 */
  background: rgb(var(--v-theme-surface));
  color: rgb(var(--v-theme-primary));
  border-color: rgb(var(--v-theme-error));

  /* 带透明度（遵循 Material Design 规范）*/
  color: rgba(var(--v-theme-on-surface), 0.87); /* 高强调 */
  color: rgba(var(--v-theme-on-surface), 0.6); /* 中强调 */
  color: rgba(var(--v-theme-on-surface), 0.38); /* 禁用/低强调 */
  border-color: rgba(var(--v-theme-on-surface), 0.12); /* 边框 */
}
```

### 透明度规范 (Material Design)

| 用途   | 透明度 | Vuetify 变量                  |
| ------ | ------ | ----------------------------- |
| 高强调 | 0.87   | `--v-high-emphasis-opacity`   |
| 中强调 | 0.60   | `--v-medium-emphasis-opacity` |
| 禁用   | 0.38   | `--v-disabled-opacity`        |
| 边框   | 0.12   | `--v-border-opacity`          |
| 悬停   | 0.04   | `--v-hover-opacity`           |
| 聚焦   | 0.12   | `--v-focus-opacity`           |

### 组件属性优先

优先使用 Vuetify 组件属性而非自定义 CSS：

```vue
<!-- ✅ 正确：使用组件属性 -->
<v-btn color="primary" variant="elevated" />
<v-card color="surface" />
<v-alert type="error" />

<!-- ❌ 禁止：内联样式硬编码 -->
<v-btn style="background: #1976d2" />
```

## @mdi/font 图标规范

统一使用 Material Design Icons，通过 `@mdi/font` 包提供。

### 基础用法

```vue
<!-- 方式一：icon 属性 -->
<v-icon icon="mdi-home" />

<!-- 方式二：插槽内容 -->
<v-icon>mdi-account</v-icon>

<!-- 带属性 -->
<v-icon icon="mdi-heart" size="large" color="error" />
<v-icon icon="mdi-check" size="x-small" color="success" />
```

### SVG 导入（Tree-shaking 优化）

大型项目推荐按需导入，减少包体积：

```vue
<template>
  <v-icon :icon="mdiAccount" />
</template>

<script setup>
import { mdiAccount } from "@mdi/js";
</script>
```

### 图标别名

使用 Vuetify 内置别名保持一致性：

```vue
<!-- 内置别名 -->
<v-icon icon="$success" />
<!-- mdi-check-circle -->
<v-icon icon="$info" />
<!-- mdi-information -->
<v-icon icon="$warning" />
<!-- mdi-alert -->
<v-icon icon="$error" />
<!-- mdi-alert-circle -->
```

### 常用图标速查

| 场景      | 图标                                      |
| --------- | ----------------------------------------- |
| 用户/账户 | `mdi-account`                             |
| 登录      | `mdi-login`                               |
| 登出      | `mdi-logout`                              |
| 设置      | `mdi-cog`                                 |
| 编辑      | `mdi-pencil`                              |
| 删除      | `mdi-delete`                              |
| 添加      | `mdi-plus`                                |
| 搜索      | `mdi-magnify`                             |
| 主题切换  | `mdi-weather-sunny` / `mdi-weather-night` |
| 成功      | `mdi-check-circle`                        |
| 警告      | `mdi-alert`                               |
| 错误      | `mdi-close-circle`                        |
| 刷新      | `mdi-refresh`                             |
| 菜单      | `mdi-menu`                                |

## 禁止事项

- ❌ 硬编码颜色值（`#ffffff`, `rgb(0,0,0)` 等）
- ❌ 使用非 @mdi/font 的图标库（Font Awesome, Heroicons 等）
- ❌ 内联 SVG 图标（使用 `@mdi/js` 代替）
- ❌ 使用 `!important` 覆盖 Vuetify 样式
- ❌ 在组件中定义与主题冲突的固定颜色
- ❌ 忽略透明度规范直接使用纯色文字
