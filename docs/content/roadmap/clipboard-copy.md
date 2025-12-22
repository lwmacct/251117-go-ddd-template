# 复制到剪贴板功能

> **状态**: ✅ 已完成
> **优先级**: 低
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:24+4`
- [已实现功能](#已实现功能) `:28+15`
  - [功能特性](#功能特性) `:30+7`
  - [技术特性](#技术特性) `:37+6`
- [组件接口](#组件接口) `:43+24`
  - [CopyButton 组件](#copybutton-组件) `:45+12`
  - [useClipboard Composable](#useclipboard-composable) `:57+10`
- [集成位置](#集成位置) `:67+4`
- [代码位置](#代码位置) `:71+12`
- [使用示例](#使用示例) `:83+14`

<!--TOC-->

## 需求背景

用户需要快速复制表格中的 ID、邮箱等信息，手动选择文本复制体验不佳。

## 已实现功能

### 功能特性

- 一键复制文本到剪贴板
- 复制成功后显示视觉反馈（图标变绿）
- 2 秒后自动恢复原状态
- 支持工具提示显示复制状态

### 技术特性

- 渐进增强：优先使用 Clipboard API
- 降级支持：旧浏览器使用 execCommand
- 安全上下文检测

## 组件接口

### CopyButton 组件

```vue
<CopyButton
  :text="要复制的文本"
  :size="按钮大小"           <!-- x-small | small | default -->
  :icon-only="仅图标模式"    <!-- 默认 true -->
  :success-text="成功提示"   <!-- 默认 '已复制' -->
  :color="按钮颜色"
/>
```

### useClipboard Composable

```typescript
const { copied, error, copy } = useClipboard({
  successDuration: 2000, // 成功状态持续时间
});

await copy("要复制的文本");
```

## 集成位置

- **用户管理表格** - ID 列添加复制按钮

## 代码位置

```
web/src/
├── composables/
│   └── useClipboard.ts          # 剪贴板操作 Composable
├── components/
│   └── CopyButton.vue           # 可复用复制按钮组件
└── pages/admin/users/
    └── index.vue                # 用户表格集成
```

## 使用示例

```vue
<script setup>
import CopyButton from "@/components/CopyButton.vue";
</script>

<template>
  <div class="d-flex align-center">
    <span>{{ userId }}</span>
    <CopyButton :text="String(userId)" size="x-small" />
  </div>
</template>
```
