# Timer Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:35+4`
- [已实现功能](#已实现功能) `:39+30`
  - [超时定时器](#超时定时器) `:41+5`
  - [间隔定时器](#间隔定时器) `:46+5`
  - [时间戳](#时间戳) `:51+5`
  - [动画帧](#动画帧) `:56+4`
  - [日期格式化](#日期格式化) `:60+4`
  - [高级功能](#高级功能) `:64+5`
- [使用方式](#使用方式) `:69+112`
  - [超时定时器](#超时定时器-1) `:71+19`
  - [间隔定时器](#间隔定时器-1) `:90+20`
  - [实时时间](#实时时间) `:110+11`
  - [requestAnimationFrame](#requestanimationframe) `:121+12`
  - [日期格式化](#日期格式化-1) `:133+16`
  - [任务调度器](#任务调度器) `:149+18`
  - [空闲回调](#空闲回调) `:167+14`
- [API](#api) `:181+34`
  - [useTimeout](#usetimeout) `:183+9`
  - [useInterval](#useinterval) `:192+10`
  - [useScheduler](#usescheduler) `:202+13`
- [代码位置](#代码位置) `:215+7`

<!--TOC-->

## 需求背景

前端需要处理各种定时器场景：超时、轮询、动画帧、任务调度等。

## 已实现功能

### 超时定时器

- `useTimeout` - 基础超时定时器
- `useTimeoutFn` - 带回调的超时定时器

### 间隔定时器

- `useInterval` - 基础间隔定时器
- `useIntervalFn` - 带回调的间隔定时器

### 时间戳

- `useTimestamp` - 实时时间戳
- `useNow` - 实时日期时间

### 动画帧

- `useRafFn` - requestAnimationFrame 循环

### 日期格式化

- `useDateFormat` - 响应式日期格式化

### 高级功能

- `useIdleCallback` - 空闲时执行回调
- `useScheduler` - 任务调度器

## 使用方式

### 超时定时器

```typescript
import { useTimeout, useTimeoutFn } from "@/composables/useTimer";

// 基础用法
const { ready, start, stop, isPending } = useTimeout(3000);
start();
watch(ready, (val) => {
  if (val) console.log("Timeout!");
});

// 带回调
const { start: delayedStart } = useTimeoutFn(() => {
  console.log("3 seconds later...");
}, 3000);
delayedStart();
```

### 间隔定时器

```typescript
import { useInterval, useIntervalFn } from "@/composables/useTimer";

// 计数器
const { counter, pause, resume, reset } = useInterval(1000);
// counter 每秒加 1

// 带回调
useIntervalFn(
  () => {
    console.log("Tick!");
    fetchData();
  },
  5000,
  { immediate: true },
);
```

### 实时时间

```typescript
import { useNow } from "@/composables/useTimer";

const { now, date, time, pause, resume } = useNow();
// now: Date 对象
// date: '2024-11-30'
// time: '10:30:45'
```

### requestAnimationFrame

```typescript
import { useRafFn } from "@/composables/useTimer";

const position = ref(0);

const { pause, resume } = useRafFn((timestamp) => {
  position.value = Math.sin(timestamp / 1000) * 100;
});
```

### 日期格式化

```typescript
import { useDateFormat } from "@/composables/useTimer";

// 实时更新的格式化时间
const { formatted } = useDateFormat("YYYY-MM-DD HH:mm:ss", undefined, {
  updateInterval: 1000,
});
// formatted: '2024-11-30 10:30:45'

// 格式化指定日期
const birthday = ref(new Date("1990-01-15"));
const { formatted: birthdayStr } = useDateFormat("YYYY年MM月DD日", birthday);
```

### 任务调度器

```typescript
import { useScheduler } from "@/composables/useTimer";

const { addTask, pauseTask, removeTask, pauseAll } = useScheduler();

// 添加定时任务
addTask("sync", () => syncData(), 5000); // 每 5 秒同步
addTask("ping", () => ping(), 30000); // 每 30 秒 ping

// 暂停特定任务
pauseTask("sync");

// 暂停所有任务
pauseAll();
```

### 空闲回调

```typescript
import { useIdleCallback } from "@/composables/useTimer";

// 浏览器空闲时执行低优先级任务
const { cancel } = useIdleCallback(
  () => {
    heavyComputation();
  },
  { timeout: 2000 },
);
```

## API

### useTimeout

| 返回值    | 类型            | 说明         |
| --------- | --------------- | ------------ |
| ready     | `Ref<boolean> ` | 是否已超时   |
| isPending | `Ref<boolean> ` | 是否正在等待 |
| start     | `() => void   ` | 启动定时器   |
| stop      | `() => void   ` | 停止定时器   |

### useInterval

| 返回值   | 类型            | 说明         |
| -------- | --------------- | ------------ |
| counter  | `Ref<number>  ` | 触发次数     |
| isActive | `Ref<boolean> ` | 是否正在运行 |
| pause    | `() => void   ` | 暂停         |
| resume   | `() => void   ` | 恢复         |
| reset    | `() => void   ` | 重置计数     |

### useScheduler

| 返回值     | 说明         |
| ---------- | ------------ |
| tasks      | 所有任务列表 |
| addTask    | 添加任务     |
| removeTask | 移除任务     |
| pauseTask  | 暂停任务     |
| resumeTask | 恢复任务     |
| pauseAll   | 暂停所有     |
| resumeAll  | 恢复所有     |
| clearAll   | 清除所有     |

## 代码位置

```
web/src/
└── composables/
    └── useTimer.ts    # Timer Composable
```
