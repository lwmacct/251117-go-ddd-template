# 异步状态 Composable

> **状态**: ✅ 已完成
> **优先级**: 高
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:28+4`
- [已实现功能](#已实现功能) `:32+27`
  - [useAsyncState](#useasyncstate) `:34+7`
  - [useAsyncRetry](#useasyncretry) `:41+6`
  - [usePolling](#usepolling) `:47+6`
  - [usePromiseQueue](#usepromisequeue) `:53+6`
- [使用方式](#使用方式) `:59+61`
  - [基础用法](#基础用法) `:61+18`
  - [自动重试](#自动重试) `:79+12`
  - [轮询](#轮询) `:91+17`
  - [Promise 队列](#promise-队列) `:108+12`
- [API](#api) `:120+15`
  - [useAsyncState 返回值](#useasyncstate-返回值) `:122+13`
- [代码位置](#代码位置) `:135+7`

<!--TOC-->

## 需求背景

项目中大量异步操作需要管理 loading、error、data 状态，需要统一的异步状态管理方案。

## 已实现功能

### useAsyncState

- 基础异步状态管理
- loading/error/data 状态
- 延迟显示 loading（防抖）
- 成功/失败回调

### useAsyncRetry

- 自动重试机制
- 可配置重试次数和延迟
- 指数退避策略

### usePolling

- 轮询请求
- 页面不可见时自动暂停
- 手动开始/停止

### usePromiseQueue

- Promise 队列
- 按顺序执行
- 并发控制

## 使用方式

### 基础用法

```typescript
import { useAsyncState } from "@/composables/useAsync";

const { data, isLoading, error, execute } = useAsyncState((id: number) => api.getUser(id), { initialData: null });

// 执行请求
await execute(1);

// 模板中使用
// <template>
//   <v-progress-circular v-if="isLoading" />
//   <v-alert v-else-if="error" type="error">{{ error.message }}</v-alert>
//   <div v-else>{{ data?.name }}</div>
// </template>
```

### 自动重试

```typescript
import { useAsyncRetry } from "@/composables/useAsync";

const { data, execute, retryCount } = useAsyncRetry(() => fetchUnstableApi(), {
  maxRetries: 3,
  retryDelay: 1000,
  retryDelayFactor: 2, // 指数退避
});
```

### 轮询

```typescript
import { usePolling } from "@/composables/useAsync";

const { data, isPolling, start, stop } = usePolling(() => fetchStatus(), {
  interval: 5000,
  pauseOnHidden: true, // 页面不可见时暂停
});

// 开始轮询
start();

// 组件卸载时停止
onUnmounted(() => stop());
```

### Promise 队列

```typescript
import { usePromiseQueue } from "@/composables/useAsync";

const queue = usePromiseQueue();

// 按顺序执行，即使第二个先完成也会等待第一个
await queue.add(() => slowRequest());
await queue.add(() => fastRequest());
```

## API

### useAsyncState 返回值

| 属性       | 类型               | 说明                             |
| ---------- | ------------------ | -------------------------------- |
| data       | `Ref<T>`           | 数据                             |
| isLoading  | `Ref<boolean>`     | 是否加载中                       |
| isFinished | `Ref<boolean>`     | 是否已完成                       |
| isSuccess  | `Ref<boolean>`     | 是否成功                         |
| error      | Ref<Error \| null> | 错误信息                         |
| execute    | Function           | 执行函数                         |
| reset      | Function           | 重置状态                         |
| state      | ComputedRef        | 状态: idle/loading/success/error |

## 代码位置

```
web/src/
└── composables/
    └── useAsync.ts    # 异步状态管理
```
