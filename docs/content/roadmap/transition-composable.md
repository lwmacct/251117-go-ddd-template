# Transition Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:36+4`
- [已实现功能](#已实现功能) `:40+21`
  - [基础过渡](#基础过渡) `:42+7`
  - [CSS 动画](#css-动画) `:49+4`
  - [高级效果](#高级效果) `:53+8`
- [使用方式](#使用方式) `:61+285`
  - [基础过渡](#基础过渡-1) `:63+34`
  - [淡入淡出](#淡入淡出) `:97+21`
  - [滑动效果](#滑动效果) `:118+29`
  - [缩放效果](#缩放效果) `:147+17`
  - [CSS 动画](#css-动画-1) `:164+41`
  - [列表过渡](#列表过渡) `:205+32`
  - [数值过渡](#数值过渡) `:237+26`
  - [抖动效果](#抖动效果) `:263+25`
  - [脉冲效果](#脉冲效果) `:288+22`
  - [打字机效果](#打字机效果) `:310+36`
- [API](#api) `:346+41`
  - [useTransition](#usetransition) `:348+12`
  - [useSlide](#useslide) `:360+9`
  - [useAnimation](#useanimation) `:369+11`
  - [useNumberTransition](#usenumbertransition) `:380+7`
- [代码位置](#代码位置) `:387+7`

<!--TOC-->

## 需求背景

前端需要过渡和动画相关的工具函数，支持淡入淡出、滑动、缩放等效果的编程控制。

## 已实现功能

### 基础过渡

- `useTransition` - 通用过渡效果
- `useFade` - 淡入淡出
- `useSlide` - 滑动效果
- `useScale` - 缩放效果

### CSS 动画

- `useAnimation` - CSS 动画控制

### 高级效果

- `useTransitionGroup` - 列表过渡
- `useNumberTransition` - 数值平滑过渡
- `useShake` - 抖动效果
- `usePulse` - 脉冲效果
- `useTypewriter` - 打字机效果

## 使用方式

### 基础过渡

```typescript
import { useTransition } from "@/composables/useTransition";

const { isVisible, show, hide, toggle, state, transitionStyle } = useTransition({
  duration: 300,
  easing: "ease-out",
  onBeforeEnter: () => console.log("开始进入"),
  onAfterEnter: () => console.log("进入完成"),
  onBeforeLeave: () => console.log("开始离开"),
  onAfterLeave: () => console.log("离开完成"),
});

// 显示（返回 Promise）
await show();

// 隐藏
await hide();

// 切换
await toggle();

// 检查状态
console.log(state.value); // 'idle' | 'enter' | 'leave'
```

```vue
<template>
  <button @click="toggle">Toggle</button>
  <div v-show="isVisible" :style="transitionStyle">Content</div>
</template>
```

### 淡入淡出

```typescript
import { useFade } from "@/composables/useTransition";

const { isVisible, opacity, show, hide, toggle, style } = useFade({
  duration: 300,
  easing: "ease",
});

// 使用
await show(); // opacity: 0 -> 1
await hide(); // opacity: 1 -> 0
```

```vue
<template>
  <div v-show="isVisible" :style="style">Fade Content</div>
</template>
```

### 滑动效果

```typescript
import { useSlide } from "@/composables/useTransition";

// 从上方滑入
const slideUp = useSlide({
  direction: "up",
  distance: "20px",
  duration: 300,
});

// 从右侧滑入
const slideRight = useSlide({
  direction: "right",
  distance: "100%",
  duration: 500,
});

await slideUp.show();
await slideRight.show();
```

```vue
<template>
  <div v-show="slideUp.isVisible.value" :style="slideUp.style.value">Slide Up Content</div>
</template>
```

### 缩放效果

```typescript
import { useScale } from "@/composables/useTransition";

const { isVisible, scale, show, hide, style } = useScale({
  fromScale: 0.8, // 从 80% 缩放到 100%
  duration: 200,
  easing: "ease-out",
});

// 模态框打开效果
async function openModal() {
  await show();
}
```

### CSS 动画

```typescript
import { useAnimation } from "@/composables/useTransition";

const { isPlaying, isPaused, play, pause, stop, restart, animationStyle } = useAnimation({
  name: "bounce", // CSS @keyframes 名称
  duration: 1000,
  easing: "ease",
  iterations: "infinite",
  direction: "alternate",
});

// 控制动画
play();
pause();
stop();
restart();
```

```css
@keyframes bounce {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-20px);
  }
}
```

```vue
<template>
  <div :style="animationStyle">Bouncing</div>
  <button @click="isPlaying ? pause() : play()">
    {{ isPlaying ? (isPaused ? "Resume" : "Pause") : "Play" }}
  </button>
</template>
```

### 列表过渡

```typescript
import { useTransitionGroup } from "@/composables/useTransition";

interface Item {
  id: number;
  text: string;
}

const { items, addItem, removeItem, getItemStyle } = useTransitionGroup<Item>({
  duration: 300,
});

// 添加项目（带进入动画）
addItem({ id: 1, text: "Item 1" });

// 移除项目（带离开动画）
removeItem(1);
```

```vue
<template>
  <ul>
    <li v-for="item in items" :key="item.id" :style="getItemStyle(item.id)">
      {{ item.text }}
      <button @click="removeItem(item.id)">删除</button>
    </li>
  </ul>
</template>
```

### 数值过渡

```typescript
import { useNumberTransition } from "@/composables/useTransition";

const { value, tweenedValue, set, isAnimating } = useNumberTransition(0, {
  duration: 500,
  easing: (t) => t * t, // 缓动函数
});

// 设置目标值，数值会平滑过渡
set(100);

// 用于显示
// tweenedValue 会从 0 平滑变化到 100
```

```vue
<template>
  <div>
    <span>{{ Math.round(tweenedValue) }}</span>
    <button @click="set(1000)">Set to 1000</button>
  </div>
</template>
```

### 抖动效果

```typescript
import { useShake } from "@/composables/useTransition";

const { shake, isShaking, style } = useShake({
  duration: 500,
  intensity: 10, // 抖动强度（像素）
});

// 表单验证失败时触发抖动
function onValidationError() {
  shake();
}
```

```vue
<template>
  <form :style="style" @submit.prevent="validate">
    <input v-model="email" type="email" />
    <button type="submit">Submit</button>
  </form>
</template>
```

### 脉冲效果

```typescript
import { usePulse } from "@/composables/useTransition";

const { pulse, isPulsing, stop, style } = usePulse({
  scale: 1.1, // 放大到 110%
  duration: 300,
});

// 点击时脉冲
function onClick() {
  pulse();
}
```

```vue
<template>
  <button :style="style" @click="onClick">Click Me</button>
</template>
```

### 打字机效果

```typescript
import { useTypewriter } from "@/composables/useTransition";

const { text, start, pause, reset, isTyping, isComplete } = useTypewriter("Hello, World! Welcome to our application.", {
  speed: 50, // 每个字符的间隔（毫秒）
  delay: 500, // 开始延迟
});

onMounted(() => {
  start();
});
```

```vue
<template>
  <p>
    {{ text }}
    <span v-if="isTyping" class="cursor">|</span>
  </p>
  <div>
    <button
      @click="
        reset();
        start();
      "
    >
      重新开始
    </button>
    <button v-if="isTyping" @click="pause">暂停</button>
    <button v-else-if="!isComplete" @click="start">继续</button>
  </div>
</template>
```

## API

### useTransition

| 选项          | 类型     | 默认值 | 说明             |
| ------------- | -------- | ------ | ---------------- |
| duration      | number   | 300    | 持续时间（毫秒） |
| easing        | string   | 'ease' | 缓动函数         |
| delay         | number   | 0      | 延迟时间         |
| onBeforeEnter | Function | -      | 进入前回调       |
| onAfterEnter  | Function | -      | 进入后回调       |
| onBeforeLeave | Function | -      | 离开前回调       |
| onAfterLeave  | Function | -      | 离开后回调       |

### useSlide

| 选项      | 类型   | 默认值 | 说明                       |
| --------- | ------ | ------ | -------------------------- |
| direction | string | 'up'   | 方向（up/down/left/right） |
| distance  | string | '20px' | 滑动距离                   |
| duration  | number | 300    | 持续时间                   |
| easing    | string | 'ease' | 缓动函数                   |

### useAnimation

| 选项       | 类型               | 默认值   | 说明         |
| ---------- | ------------------ | -------- | ------------ |
| name       | string             | -        | CSS 动画名称 |
| duration   | number             | 1000     | 持续时间     |
| easing     | string             | 'ease'   | 缓动函数     |
| iterations | number\|'infinite' | 1        | 迭代次数     |
| direction  | string             | 'normal' | 动画方向     |
| fillMode   | string             | 'none'   | 填充模式     |

### useNumberTransition

| 选项     | 类型     | 默认值 | 说明     |
| -------- | -------- | ------ | -------- |
| duration | number   | 500    | 过渡时间 |
| easing   | Function | linear | 缓动函数 |

## 代码位置

```
web/src/
└── composables/
    └── useTransition.ts    # Transition Composable
```
