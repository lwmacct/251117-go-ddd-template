# 交叉观察器 Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:29+4`
- [已实现功能](#已实现功能) `:33+33`
  - [useIntersectionObserver](#useintersectionobserver) `:35+7`
  - [useLazyLoad](#uselazyload) `:42+6`
  - [useInfiniteScroll](#useinfinitescroll) `:48+7`
  - [useAnimateOnScroll](#useanimateonscroll) `:55+6`
  - [useVisibility](#usevisibility) `:61+5`
- [使用方式](#使用方式) `:66+63`
  - [基础用法](#基础用法) `:68+10`
  - [懒加载图片](#懒加载图片) `:78+13`
  - [无限滚动](#无限滚动) `:91+26`
  - [滚动动画](#滚动动画) `:117+12`
- [API](#api) `:129+12`
  - [useIntersectionObserver 返回值](#useintersectionobserver-返回值) `:131+10`
- [代码位置](#代码位置) `:141+7`

<!--TOC-->

## 需求背景

需要检测元素可见性以实现懒加载、无限滚动、进入动画等常见功能。

## 已实现功能

### useIntersectionObserver

- 基础交叉观察器封装
- 可见性状态追踪
- 交叉比例追踪
- 自定义回调支持

### useLazyLoad

- 懒加载实现
- 提前加载（预加载边距）
- 单次触发选项

### useInfiniteScroll

- 无限滚动实现
- 加载状态管理
- 还有更多数据检测
- 自动触发加载

### useAnimateOnScroll

- 滚动进入动画
- 自定义动画类
- 单次/重复触发

### useVisibility

- 可见性监控
- 进入/离开回调

## 使用方式

### 基础用法

```typescript
import { ref } from "vue";
import { useIntersectionObserver } from "@/composables/useIntersectionObserver";

const target = ref<HTMLElement>();
const { isVisible, hasBeenVisible } = useIntersectionObserver(target);
```

### 懒加载图片

```typescript
import { useLazyLoad } from "@/composables/useIntersectionObserver";

const imageRef = ref<HTMLImageElement>();
const { shouldLoad } = useLazyLoad(imageRef);

// 模板中:
// <img v-if="shouldLoad" :src="imageSrc" ref="imageRef" />
// <div v-else ref="imageRef" class="placeholder" />
```

### 无限滚动

```typescript
import { useInfiniteScroll } from "@/composables/useIntersectionObserver";

const sentinelRef = ref<HTMLElement>();
const loading = ref(false);
const hasMore = ref(true);

useInfiniteScroll(sentinelRef, {
  loading,
  hasMore,
  onLoadMore: async () => {
    loading.value = true;
    await fetchMoreData();
    loading.value = false;
  },
});

// 模板中:
// <div v-for="item in items">...</div>
// <div ref="sentinelRef" v-show="hasMore">
//   <v-progress-circular v-if="loading" />
// </div>
```

### 滚动动画

```typescript
import { useAnimateOnScroll } from "@/composables/useIntersectionObserver";

const cardRef = ref<HTMLElement>();
useAnimateOnScroll(cardRef, {
  animationClass: "fade-in-up",
  threshold: 0.2,
});
```

## API

### useIntersectionObserver 返回值

| 属性              | 类型           | 说明           |
| ----------------- | -------------- | -------------- |
| isVisible         | `Ref<boolean>` | 当前是否可见   |
| hasBeenVisible    | `Ref<boolean>` | 是否曾经可见过 |
| intersectionRatio | `Ref<number>`  | 交叉比例 (0-1) |
| observe           | `() => void`   | 开始观察       |
| unobserve         | `() => void`   | 停止观察       |

## 代码位置

```
web/src/
└── composables/
    └── useIntersectionObserver.ts    # 交叉观察器
```
