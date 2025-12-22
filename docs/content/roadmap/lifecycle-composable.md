# Lifecycle Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:42+4`
- [已实现功能](#已实现功能) `:46+43`
  - [挂载状态](#挂载状态) `:48+7`
  - [生命周期包装](#生命周期包装) `:55+7`
  - [延迟执行](#延迟执行) `:62+6`
  - [条件与控制](#条件与控制) `:68+6`
  - [追踪与调试](#追踪与调试) `:74+7`
  - [工具函数](#工具函数) `:81+8`
- [使用方式](#使用方式) `:89+247`
  - [安全的异步操作](#安全的异步操作) `:91+21`
  - [完整挂载状态](#完整挂载状态) `:112+12`
  - [带清理的挂载](#带清理的挂载) `:124+17`
  - [异步等待挂载](#异步等待挂载) `:141+14`
  - [延迟执行](#延迟执行-1) `:155+21`
  - [条件挂载](#条件挂载) `:176+23`
  - [生命周期追踪](#生命周期追踪) `:199+19`
  - [清理函数管理](#清理函数管理) `:218+26`
  - [渲染计数](#渲染计数) `:244+11`
  - [错误捕获](#错误捕获) `:255+19`
  - [keep-alive 支持](#keep-alive-支持) `:274+26`
  - [组件存活时间](#组件存活时间) `:300+11`
  - [组合生命周期钩子](#组合生命周期钩子) `:311+25`
- [API](#api) `:336+37`
  - [useMountedState](#usemountedstate) `:338+8`
  - [useLifecycleTracker](#uselifecycletracker) `:346+13`
  - [useAsyncMounted](#useasyncmounted) `:359+7`
  - [useErrorCapture](#useerrorcapture) `:366+7`
- [代码位置](#代码位置) `:373+7`

<!--TOC-->

## 需求背景

前端需要增强的生命周期钩子工具函数，支持安全的异步操作、清理函数管理、状态追踪等高级功能。

## 已实现功能

### 挂载状态

- `useMountedState` - 完整挂载状态（挂载中/已挂载/已卸载）
- `useSafeMounted` - 安全的挂载状态检查
- `useMounted` - 带清理函数的挂载
- `useAsyncMounted` - 异步等待挂载完成

### 生命周期包装

- `useUnmounted` - 卸载时执行
- `useUpdated` - 更新时执行
- `useActivated` - 激活时执行（keep-alive）
- `useDeactivated` - 停用时执行（keep-alive）

### 延迟执行

- `useMountedDelay` - 延迟执行
- `useMountedNextFrame` - 下一帧执行
- `useMountedNextTick` - 下一个 tick 执行

### 条件与控制

- `useMountedWhen` - 条件满足时执行
- `useMountedOnce` - 只执行一次
- `useMountedOrActivated` - 挂载或激活时执行

### 追踪与调试

- `useLifecycleTracker` - 生命周期追踪器
- `useRenderCount` - 渲染计数
- `useRenderTracking` - 渲染追踪
- `useComponentAliveTime` - 组件存活时间

### 工具函数

- `useCleanup` - 清理函数管理
- `useInstance` - 获取组件实例
- `useErrorCapture` - 错误捕获
- `useComponentVisible` - 组件可见状态
- `useLifecycle` - 组合生命周期钩子

## 使用方式

### 安全的异步操作

```typescript
import { useSafeMounted } from "@/composables/useLifecycle";

const data = ref(null);
const isMounted = useSafeMounted();

async function loadData() {
  const result = await api.fetchData();

  // 检查组件是否仍然挂载
  if (!isMounted.value) return;

  // 安全更新状态
  data.value = result;
}

onMounted(loadData);
```

### 完整挂载状态

```typescript
import { useMountedState } from '@/composables/useLifecycle'

const { isMounted, isMounting, isUnmounted } = useMountedState()

// 在模板中
<div v-if="isMounting">加载中...</div>
<div v-else-if="isMounted">内容</div>
```

### 带清理的挂载

```typescript
import { useMounted } from "@/composables/useLifecycle";

useMounted(() => {
  // 设置事件监听
  const handler = (e) => console.log(e);
  window.addEventListener("resize", handler);

  // 返回清理函数，会在卸载时自动调用
  return () => {
    window.removeEventListener("resize", handler);
  };
});
```

### 异步等待挂载

```typescript
import { useAsyncMounted } from "@/composables/useLifecycle";

async function setup() {
  // 等待组件挂载完成
  await useAsyncMounted({ timeout: 5000 });

  // 此时 DOM 已准备就绪
  initializeChart();
}
```

### 延迟执行

```typescript
import { useMountedDelay, useMountedNextFrame, useMountedNextTick } from "@/composables/useLifecycle";

// 延迟 500ms 执行
useMountedDelay(() => {
  showWelcomeMessage();
}, 500);

// 下一帧执行（适合 DOM 测量）
useMountedNextFrame(() => {
  measureElementSize();
});

// 下一个 Vue tick 执行
useMountedNextTick(() => {
  scrollToBottom();
});
```

### 条件挂载

```typescript
import { useMountedWhen } from "@/composables/useLifecycle";

const isReady = ref(false);

useMountedWhen(isReady, () => {
  console.log("条件满足且已挂载");
  initializeFeature();

  return () => {
    cleanupFeature();
  };
});

// 稍后设置为 true 触发回调
onMounted(async () => {
  await loadDependencies();
  isReady.value = true;
});
```

### 生命周期追踪

```typescript
import { useLifecycleTracker } from "@/composables/useLifecycle";

const { history, currentPhase, isActive } = useLifecycleTracker({
  log: import.meta.env.DEV, // 开发环境输出日志
  prefix: "UserProfile",
});

// 查看历史
console.log(history.value);
// [
//   { event: 'beforeMount', timestamp: 1699999000000 },
//   { event: 'mounted', timestamp: 1699999000010 },
//   ...
// ]
```

### 清理函数管理

```typescript
import { useCleanup } from "@/composables/useLifecycle";

const { onCleanup, cleanup } = useCleanup();

onMounted(() => {
  const timer1 = setInterval(() => {}, 1000);
  const timer2 = setTimeout(() => {}, 5000);
  const subscription = observable.subscribe();

  // 注册多个清理函数
  onCleanup(() => clearInterval(timer1));
  onCleanup(() => clearTimeout(timer2));
  onCleanup(() => subscription.unsubscribe());

  // 所有清理函数会在卸载时自动执行
});

// 也可以手动触发清理
function handleReset() {
  cleanup();
}
```

### 渲染计数

```typescript
import { useRenderCount } from '@/composables/useLifecycle'

const { count, reset } = useRenderCount()

// 在开发模式显示
<div v-if="isDev">渲染次数: {{ count }}</div>
```

### 错误捕获

```typescript
import { useErrorCapture } from '@/composables/useLifecycle'

const { error, clearError } = useErrorCapture((err, instance, info) => {
  // 报告错误
  errorReporter.report(err, { componentInfo: info })

  // 返回 false 让错误继续传播
  // 返回 true 或不返回则阻止传播
  return true
})

// 显示错误 UI
<ErrorBoundary v-if="error" :error="error" @dismiss="clearError" />
<Content v-else />
```

### keep-alive 支持

```typescript
import { useMountedOrActivated, useComponentVisible } from "@/composables/useLifecycle";

// 每次显示时刷新数据
useMountedOrActivated(() => {
  refreshData();

  return () => {
    cancelRequests();
  };
});

// 追踪可见状态
const isVisible = useComponentVisible();

watch(isVisible, (visible) => {
  if (visible) {
    startPolling();
  } else {
    stopPolling();
  }
});
```

### 组件存活时间

```typescript
import { useComponentAliveTime } from '@/composables/useLifecycle'

const { aliveTime, mountedAt } = useComponentAliveTime()

// 显示组件已运行时间
<span>运行时间: {{ Math.floor(aliveTime / 1000) }}秒</span>
```

### 组合生命周期钩子

```typescript
import { useLifecycle } from "@/composables/useLifecycle";

useLifecycle({
  onMounted: () => {
    console.log("挂载完成");
    initData();
  },
  onUnmounted: () => {
    console.log("即将卸载");
    cleanup();
  },
  onActivated: () => {
    console.log("激活");
    refresh();
  },
  onDeactivated: () => {
    console.log("停用");
    pause();
  },
});
```

## API

### useMountedState

| 返回值      | 类型           | 说明         |
| ----------- | -------------- | ------------ |
| isMounted   | Ref\<boolean\> | 是否已挂载   |
| isMounting  | Ref\<boolean\> | 是否正在挂载 |
| isUnmounted | Ref\<boolean\> | 是否已卸载   |

### useLifecycleTracker

| 选项   | 类型    | 默认值      | 说明         |
| ------ | ------- | ----------- | ------------ |
| log    | boolean | false       | 是否输出日志 |
| prefix | string  | 'Component' | 日志前缀     |

| 返回值       | 类型           | 说明     |
| ------------ | -------------- | -------- |
| history      | Ref\<Array\>   | 事件历史 |
| currentPhase | Ref\<string\>  | 当前阶段 |
| isActive     | Ref\<boolean\> | 是否活跃 |

### useAsyncMounted

| 选项      | 类型     | 默认值 | 说明     |
| --------- | -------- | ------ | -------- |
| timeout   | number   | -      | 超时时间 |
| onTimeout | Function | -      | 超时回调 |

### useErrorCapture

| 返回值     | 类型                 | 说明       |
| ---------- | -------------------- | ---------- |
| error      | Ref\<Error \| null\> | 捕获的错误 |
| clearError | `() => void`         | 清除错误   |

## 代码位置

```
web/src/
└── composables/
    └── useLifecycle.ts    # Lifecycle Composable
```
