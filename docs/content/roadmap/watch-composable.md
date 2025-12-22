# Watch Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:39+4`
- [已实现功能](#已实现功能) `:43+29`
  - [执行控制](#执行控制) `:45+7`
  - [性能优化](#性能优化) `:52+5`
  - [条件监听](#条件监听) `:57+6`
  - [数组监听](#数组监听) `:63+5`
  - [其他](#其他) `:68+4`
- [使用方式](#使用方式) `:72+222`
  - [只执行一次](#只执行一次) `:74+17`
  - [防抖监听](#防抖监听) `:91+20`
  - [节流监听](#节流监听) `:111+21`
  - [可暂停监听](#可暂停监听) `:132+23`
  - [可忽略更新](#可忽略更新) `:155+17`
  - [条件执行](#条件执行) `:172+23`
  - [等待条件](#等待条件) `:195+27`
  - [带过滤的监听](#带过滤的监听) `:222+17`
  - [限制执行次数](#限制执行次数) `:239+17`
  - [数组变化监听](#数组变化监听) `:256+20`
  - [响应式数组](#响应式数组) `:276+18`
- [API](#api) `:294+36`
  - [watchDebounced](#watchdebounced) `:296+7`
  - [watchThrottled](#watchthrottled) `:303+8`
  - [watchPausable](#watchpausable) `:311+9`
  - [until](#until) `:320+10`
- [代码位置](#代码位置) `:330+7`

<!--TOC-->

## 需求背景

前端需要增强的 watch 工具函数，支持防抖、节流、暂停、条件执行等高级功能。

## 已实现功能

### 执行控制

- `watchOnce` - 只执行一次
- `watchAtMost` - 限制执行次数
- `watchPausable` - 可暂停/恢复
- `watchIgnorable` - 可忽略更新

### 性能优化

- `watchDebounced` - 防抖监听
- `watchThrottled` - 节流监听

### 条件监听

- `watchWithFilter` - 带过滤条件
- `whenever` - 条件为真时执行
- `until` - 等待条件满足

### 数组监听

- `watchArray` - 监听数组变化
- `useWatchArray` - 响应式数组操作

### 其他

- `watchTriggered` - 触发计数

## 使用方式

### 只执行一次

```typescript
import { watchOnce } from "@/composables/useWatch";

const data = ref(null);

// 数据加载完成后只执行一次
watchOnce(
  () => data.value,
  (newData) => {
    console.log("数据已加载:", newData);
    initializeApp();
  },
);
```

### 防抖监听

```typescript
import { watchDebounced } from "@/composables/useWatch";

const searchQuery = ref("");

watchDebounced(
  searchQuery,
  (query) => {
    // 用户停止输入 500ms 后执行搜索
    performSearch(query);
  },
  {
    debounce: 500,
    maxWait: 2000, // 最大等待 2 秒
  },
);
```

### 节流监听

```typescript
import { watchThrottled } from "@/composables/useWatch";

const scrollPosition = ref(0);

watchThrottled(
  scrollPosition,
  (pos) => {
    // 每 100ms 最多执行一次
    updateVisibleItems(pos);
  },
  {
    throttle: 100,
    leading: true,
    trailing: true,
  },
);
```

### 可暂停监听

```typescript
import { watchPausable } from "@/composables/useWatch";

const { pause, resume, isActive, stop } = watchPausable(
  () => data.value,
  (newData) => {
    processData(newData);
  },
);

// 编辑模式时暂停
const startEdit = () => {
  pause();
};

// 完成编辑后恢复
const finishEdit = () => {
  resume();
};
```

### 可忽略更新

```typescript
import { watchIgnorable } from "@/composables/useWatch";

const count = ref(0);

const { ignoreUpdates, stop } = watchIgnorable(count, (value) => {
  console.log("count changed:", value);
});

// 这次更新不会触发回调
ignoreUpdates(() => {
  count.value = 100;
});
```

### 条件执行

```typescript
import { whenever } from "@/composables/useWatch";

const isReady = ref(false);

// 当 isReady 变为 true 时执行
whenever(isReady, () => {
  console.log("Ready!");
  startApplication();
});

// 只执行一次
whenever(
  isReady,
  () => {
    console.log("First ready!");
  },
  { once: true },
);
```

### 等待条件

```typescript
import { until } from "@/composables/useWatch";

const data = ref(null);
const status = ref("loading");

// 等待数据加载
await until(data).toBeNotNull();
console.log("数据已加载:", data.value);

// 等待特定状态
await until(status).toBe("completed");
console.log("处理完成");

// 带超时
try {
  await until(data).timeout(5000).toBeNotNull();
} catch {
  console.log("加载超时");
}

// 等待满足条件
await until(count).toMatch((v) => v > 10);
```

### 带过滤的监听

```typescript
import { watchWithFilter } from "@/composables/useWatch";

const count = ref(0);

watchWithFilter(
  () => count.value,
  (value) => {
    // 只在值为偶数时执行
    console.log("Even number:", value);
  },
  (value) => value % 2 === 0,
);
```

### 限制执行次数

```typescript
import { watchAtMost } from "@/composables/useWatch";

// 最多执行 3 次
const { count, stop } = watchAtMost(
  () => data.value,
  (value) => {
    console.log("Executed:", value);
  },
  3,
);

console.log(count.value); // 当前执行次数
```

### 数组变化监听

```typescript
import { watchArray } from "@/composables/useWatch";

const items = ref([{ id: 1 }, { id: 2 }]);

watchArray(items, {
  onAdd: (added) => {
    console.log("新增:", added);
  },
  onRemove: (removed) => {
    console.log("删除:", removed);
  },
  onUpdate: (newArray, oldArray) => {
    console.log("更新:", newArray);
  },
});
```

### 响应式数组

```typescript
import { useWatchArray } from "@/composables/useWatch";

const { array, push, pop, clear, onChange } = useWatchArray<number>([1, 2, 3]);

// 监听变化
onChange((action, items) => {
  console.log(action, items);
  // 'push', [4] 或 'pop', [3] 等
});

push(4);
pop();
clear();
```

## API

### watchDebounced

| 选项     | 类型   | 默认值 | 说明             |
| -------- | ------ | ------ | ---------------- |
| debounce | number | 250    | 防抖延迟（毫秒） |
| maxWait  | number | -      | 最大等待时间     |

### watchThrottled

| 选项     | 类型    | 默认值 | 说明             |
| -------- | ------- | ------ | ---------------- |
| throttle | number  | 100    | 节流间隔（毫秒） |
| leading  | boolean | true   | 是否在开始时执行 |
| trailing | boolean | true   | 是否在结束时执行 |

### watchPausable

| 返回值   | 类型             | 说明     |
| -------- | ---------------- | -------- |
| pause    | `() => void    ` | 暂停监听 |
| resume   | `() => void    ` | 恢复监听 |
| isActive | Ref\<boolean\>   | 是否活跃 |
| stop     | WatchStopHandle  | 停止监听 |

### until

| 方法        | 说明                    |
| ----------- | ----------------------- |
| toBe(v)     | 等待值等于 v            |
| toBeTruthy  | 等待值为真              |
| toBeNotNull | 等待值非 null/undefined |
| toMatch(fn) | 等待满足条件            |
| timeout(ms) | 设置超时                |

## 代码位置

```
web/src/
└── composables/
    └── useWatch.ts    # Watch Composable
```
