# 交叉观察器 Composable

<!--TOC-->

- [需求背景](#需求背景) `:27:30`
- [已实现功能](#已实现功能) `:31:32`
  - [useIntersectionObserver](#useintersectionobserver) `:33:39`
  - [useLazyLoad](#uselazyload) `:40:45`
  - [useInfiniteScroll](#useinfinitescroll) `:46:52`
  - [useAnimateOnScroll](#useanimateonscroll) `:53:58`
  - [useVisibility](#usevisibility) `:59:63`
- [使用方式](#使用方式) `:64:65`
  - [基础用法](#基础用法) `:66:75`
  - [懒加载图片](#懒加载图片) `:76:88`
  - [无限滚动](#无限滚动) `:89:114`
  - [滚动动画](#滚动动画) `:115:126`
- [API](#api) `:127:128`
  - [useIntersectionObserver 返回值](#useintersectionobserver-返回值) `:129:138`
- [代码位置](#代码位置) `:139:145`

<!--TOC-->

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

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
