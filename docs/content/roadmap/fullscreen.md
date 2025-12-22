# 全屏 API Composable

> **状态**: ✅ 已完成
> **优先级**: 低
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:26+4`
- [已实现功能](#已实现功能) `:30+19`
  - [useFullscreen](#usefullscreen) `:32+8`
  - [useDocumentFullscreen](#usedocumentfullscreen) `:40+4`
  - [useFullscreenButton](#usefullscreenbutton) `:44+5`
- [使用方式](#使用方式) `:49+47`
  - [元素全屏](#元素全屏) `:51+21`
  - [文档全屏](#文档全屏) `:72+8`
  - [全屏按钮](#全屏按钮) `:80+16`
- [API](#api) `:96+13`
  - [useFullscreen 返回值](#usefullscreen-返回值) `:98+11`
- [代码位置](#代码位置) `:109+7`

<!--TOC-->

## 需求背景

需要为视频播放器、数据可视化、演示模式等场景提供全屏功能支持。

## 已实现功能

### useFullscreen

- 元素/文档全屏控制
- 进入/退出/切换全屏
- 全屏状态检测
- 跨浏览器兼容
- 组件卸载自动退出

### useDocumentFullscreen

- 文档全屏简化版

### useFullscreenButton

- 全屏按钮数据
- 图标和提示文本

## 使用方式

### 元素全屏

```typescript
import { ref } from "vue";
import { useFullscreen } from "@/composables/useFullscreen";

const videoRef = ref<HTMLVideoElement>();
const { isFullscreen, toggle, enter, exit } = useFullscreen(videoRef);
```

```vue
<template>
  <div>
    <video ref="videoRef" src="video.mp4" />
    <v-btn @click="toggle">
      {{ isFullscreen ? "退出全屏" : "全屏播放" }}
    </v-btn>
  </div>
</template>
```

### 文档全屏

```typescript
import { useDocumentFullscreen } from "@/composables/useFullscreen";

const { isFullscreen, toggle } = useDocumentFullscreen();
```

### 全屏按钮

```typescript
import { useFullscreenButton } from "@/composables/useFullscreen";

const { icon, tooltip, toggle, isSupported } = useFullscreenButton();
```

```vue
<template>
  <v-btn v-if="isSupported" :icon="icon" @click="toggle">
    <v-tooltip activator="parent">{{ tooltip }}</v-tooltip>
  </v-btn>
</template>
```

## API

### useFullscreen 返回值

| 属性              | 类型                  | 说明         |
| ----------------- | --------------------- | ------------ |
| isFullscreen      | `Ref<boolean>`        | 是否全屏     |
| isSupported       | `Ref<boolean>`        | 是否支持全屏 |
| fullscreenElement | Ref<Element \| null>  | 当前全屏元素 |
| enter             | `() => Promise<void>` | 进入全屏     |
| exit              | `() => Promise<void>` | 退出全屏     |
| toggle            | `() => Promise<void>` | 切换全屏     |

## 代码位置

```
web/src/
└── composables/
    └── useFullscreen.ts    # 全屏 API
```
