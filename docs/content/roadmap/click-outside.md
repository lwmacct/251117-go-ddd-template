# 点击外部检测 Composable

<!--TOC-->

- [需求背景](#需求背景) `:25+4`
- [已实现功能](#已实现功能) `:29+19`
  - [useClickOutside](#useclickoutside) `:31+7`
  - [useClickOutsideToggle](#useclickoutsidetoggle) `:38+5`
  - [vClickOutside 指令](#vclickoutside-指令) `:43+5`
- [使用方式](#使用方式) `:48+63`
  - [Composable 用法](#composable-用法) `:50+14`
  - [可切换用法](#可切换用法) `:64+18`
  - [指令用法](#指令用法) `:82+18`
  - [忽略特定元素](#忽略特定元素) `:100+11`
- [API](#api) `:111+12`
  - [useClickOutside 选项](#useclickoutside-选项) `:113+10`
- [代码位置](#代码位置) `:123+7`

<!--TOC-->

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

## 需求背景

下拉菜单、弹出框等组件需要在点击外部时关闭，需要统一的点击外部检测方案。

## 已实现功能

### useClickOutside

- 检测点击是否在指定元素外部
- 支持多个目标元素
- 支持忽略特定元素
- 可配置事件类型

### useClickOutsideToggle

- 带开关状态的点击外部检测
- 适用于下拉菜单、弹出框

### vClickOutside 指令

- Vue 指令形式
- 简化模板使用

## 使用方式

### Composable 用法

```typescript
import { ref } from "vue";
import { useClickOutside } from "@/composables/useClickOutside";

const menuRef = ref<HTMLElement>();

useClickOutside(menuRef, () => {
  console.log("点击了菜单外部");
  closeMenu();
});
```

### 可切换用法

```typescript
import { useClickOutsideToggle } from "@/composables/useClickOutside";

const dropdownRef = ref<HTMLElement>();
const { isOpen, toggle, close } = useClickOutsideToggle(dropdownRef);
```

```vue
<template>
  <div ref="dropdownRef">
    <button @click="toggle">切换菜单</button>
    <div v-if="isOpen" class="dropdown-menu">菜单内容</div>
  </div>
</template>
```

### 指令用法

```typescript
// main.ts
import { vClickOutside } from "@/composables/useClickOutside";
app.directive("click-outside", vClickOutside);
```

```vue
<template>
  <!-- 简单用法 -->
  <div v-click-outside="handleClose">...</div>

  <!-- 带忽略选项 -->
  <div v-click-outside="{ handler: handleClose, ignore: ['.ignore-me'] }">...</div>
</template>
```

### 忽略特定元素

```typescript
const triggerRef = ref<HTMLElement>();
const menuRef = ref<HTMLElement>();

useClickOutside(menuRef, closeMenu, {
  ignore: [triggerRef, ".modal-overlay"],
});
```

## API

### useClickOutside 选项

| 选项             | 类型    | 默认值        | 说明         |
| ---------------- | ------- | ------------- | ------------ |
| immediate        | boolean | true          | 是否立即激活 |
| event            | string  | "pointerdown" | 事件类型     |
| detectRightClick | boolean | true          | 是否检测右键 |
| ignore           | array   | []            | 忽略的元素   |
| capture          | boolean | true          | 是否捕获阶段 |

## 代码位置

```
web/src/
└── composables/
    └── useClickOutside.ts    # 点击外部检测
```
