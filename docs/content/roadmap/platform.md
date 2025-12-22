# 平台检测工具

> **状态**: ✅ 已完成
> **优先级**: 低
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:29+4`
- [已实现功能](#已实现功能) `:33+26`
  - [浏览器检测](#浏览器检测) `:35+5`
  - [操作系统检测](#操作系统检测) `:40+5`
  - [设备检测](#设备检测) `:45+5`
  - [特性检测](#特性检测) `:50+9`
- [使用方式](#使用方式) `:59+76`
  - [获取完整信息](#获取完整信息) `:61+11`
  - [便捷判断函数](#便捷判断函数) `:72+23`
  - [特性检测](#特性检测-1) `:95+22`
  - [在 Vue 中使用](#在-vue-中使用) `:117+18`
- [API](#api) `:135+23`
  - [getPlatformInfo 返回值](#getplatforminfo-返回值) `:137+12`
  - [便捷函数](#便捷函数) `:149+9`
- [代码位置](#代码位置) `:158+7`

<!--TOC-->

## 需求背景

需要检测用户的浏览器、操作系统、设备类型等信息，用于条件渲染、功能适配、统计分析等。

## 已实现功能

### 浏览器检测

- Chrome / Firefox / Safari / Edge / IE
- 版本号解析

### 操作系统检测

- Windows / macOS / iOS / Android / Linux
- 版本号解析

### 设备检测

- 设备类型：mobile / tablet / desktop
- 设备厂商和型号

### 特性检测

- 触摸支持
- PWA 独立模式
- WebGL 支持
- WebP / AVIF 支持
- 暗色模式偏好
- 减少动画偏好

## 使用方式

### 获取完整信息

```typescript
import { getPlatformInfo } from "@/utils/platform";

const info = getPlatformInfo();
console.log(info.browser.name); // "Chrome"
console.log(info.os.name); // "macOS"
console.log(info.device.type); // "desktop"
```

### 便捷判断函数

```typescript
import { isChrome, isSafari, isIOS, isAndroid, isMobile, isDesktop, isTouchDevice } from "@/utils/platform";

// 条件渲染
if (isMobile()) {
  showMobileUI();
} else {
  showDesktopUI();
}

// Safari 特殊处理
if (isSafari()) {
  applyWebKitFix();
}

// 触摸设备
if (isTouchDevice()) {
  enableTouchGestures();
}
```

### 特性检测

```typescript
import { supportsWebP, supportsWebGL, prefersDarkMode, prefersReducedMotion } from "@/utils/platform";

// 异步检测
const hasWebP = await supportsWebP();
if (hasWebP) {
  loadWebPImages();
}

// 暗色模式
if (prefersDarkMode()) {
  setDarkTheme();
}

// 减少动画
if (prefersReducedMotion()) {
  disableAnimations();
}
```

### 在 Vue 中使用

```vue
<template>
  <div :class="{ 'mobile-layout': isMobileDevice }">
    <MobileNav v-if="isMobileDevice" />
    <DesktopNav v-else />
  </div>
</template>

<script setup>
import { computed } from "vue";
import { isMobile } from "@/utils/platform";

const isMobileDevice = computed(() => isMobile());
</script>
```

## API

### getPlatformInfo 返回值

| 属性          | 类型        | 说明            |
| ------------- | ----------- | --------------- |
| browser       | BrowserInfo | 浏览器信息      |
| os            | OSInfo      | 操作系统信息    |
| device        | DeviceInfo  | 设备信息        |
| isTouch       | boolean     | 是否触摸设备    |
| isStandalone  | boolean     | 是否 PWA 模式   |
| language      | string      | 用户语言        |
| cookieEnabled | boolean     | Cookie 是否启用 |

### 便捷函数

| 函数                                     | 说明         |
| ---------------------------------------- | ------------ |
| isChrome / isFirefox / isSafari / isEdge | 浏览器判断   |
| isWindows / isMacOS / isIOS / isAndroid  | 操作系统判断 |
| isMobile / isTablet / isDesktop          | 设备类型判断 |
| isTouchDevice / isStandalone             | 特性判断     |

## 代码位置

```
web/src/
└── utils/
    └── platform.ts    # 平台检测工具
```
