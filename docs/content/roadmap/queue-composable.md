# Queue Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:31+4`
- [已实现功能](#已实现功能) `:35+17`
  - [数据结构](#数据结构) `:37+6`
  - [任务管理](#任务管理) `:43+4`
  - [应用场景](#应用场景) `:47+5`
- [使用方式](#使用方式) `:52+207`
  - [通用队列](#通用队列) `:54+30`
  - [栈](#栈) `:84+22`
  - [任务队列](#任务队列) `:106+47`
  - [通知队列](#通知队列) `:153+46`
  - [历史记录队列](#历史记录队列) `:199+36`
  - [环形缓冲区](#环形缓冲区) `:235+24`
- [API](#api) `:259+61`
  - [useQueue](#usequeue) `:261+21`
  - [useTaskQueue](#usetaskqueue) `:282+22`
  - [useNotificationQueue](#usenotificationqueue) `:304+16`
- [代码位置](#代码位置) `:320+7`

<!--TOC-->

## 需求背景

前端需要队列相关的工具函数，支持任务队列、通知队列、历史记录等高级功能。

## 已实现功能

### 数据结构

- `useQueue` - 通用队列（FIFO）
- `useStack` - 栈（LIFO）
- `useRingBuffer` - 环形缓冲区

### 任务管理

- `useTaskQueue` - 异步任务队列

### 应用场景

- `useNotificationQueue` - 通知消息队列
- `useHistoryQueue` - 历史记录队列

## 使用方式

### 通用队列

```typescript
import { useQueue } from "@/composables/useQueue";

const { items, enqueue, dequeue, peek, size, isEmpty, isFull } = useQueue<string>({
  maxSize: 100, // 最大容量
  prioritized: true, // 按优先级排序
});

// 入队
const id1 = enqueue("item1", 0); // 优先级 0
const id2 = enqueue("urgent", 10); // 优先级 10（会排在前面）

// 出队
const first = dequeue(); // 'urgent'（优先级最高）

// 查看队首（不移除）
const next = peek();

// 移除指定项
remove(id1);

// 查找
const found = find((item) => item.includes("item"));

// 检查是否包含
has(id1); // false
```

### 栈

```typescript
import { useStack } from "@/composables/useQueue";

const { items, push, pop, peek, size, isEmpty } = useStack<number>({
  maxSize: 50,
});

// 压栈
push(1);
push(2);
push(3);

// 弹栈（后进先出）
pop(); // 3
pop(); // 2

// 查看栈顶
peek(); // 1
```

### 任务队列

```typescript
import { useTaskQueue } from "@/composables/useQueue";

const { tasks, add, start, pause, resume, isRunning, isPaused, pendingCount, runningCount, completedCount, failedCount, retryFailed, waitAll } = useTaskQueue<string>({
  concurrency: 2, // 并发数
  interval: 100, // 任务间隔
  retries: 3, // 失败重试次数
  retryDelay: 1000, // 重试延迟
  autoStart: false, // 自动开始
});

// 添加任务
add(async () => {
  const res = await fetch("/api/data1");
  return res.json();
});

add(async () => {
  const res = await fetch("/api/data2");
  return res.json();
});

// 开始处理
start();

// 暂停
pause();

// 恢复
resume();

// 重试所有失败的任务
retryFailed();

// 等待所有任务完成
await waitAll();
console.log("All tasks completed!");

// 查看任务状态
console.log(`待处理: ${pendingCount.value}`);
console.log(`运行中: ${runningCount.value}`);
console.log(`已完成: ${completedCount.value}`);
console.log(`失败: ${failedCount.value}`);
```

### 通知队列

```typescript
import { useNotificationQueue } from "@/composables/useQueue";

const { notifications, add, remove, clear, info, success, warning, error } = useNotificationQueue({
  maxVisible: 5, // 最大显示数量
  defaultDuration: 3000, // 默认持续时间
  position: "top", // 新通知位置
});

// 快捷方法
info("这是一条信息");
success("操作成功！");
warning("请注意...");
error("发生错误");

// 自定义通知
add({
  type: "success",
  title: "上传完成",
  message: "文件已成功上传",
  duration: 5000,
  closable: true,
});

// 手动移除
const id = success("可以手动关闭");
remove(id);

// 清空所有
clear();
```

```vue
<template>
  <div class="notification-container">
    <div v-for="notification in notifications" :key="notification.id" :class="['notification', `notification-${notification.type}`]">
      <h4 v-if="notification.title">{{ notification.title }}</h4>
      <p>{{ notification.message }}</p>
      <button v-if="notification.closable" @click="remove(notification.id)">✕</button>
    </div>
  </div>
</template>
```

### 历史记录队列

```typescript
import { useHistoryQueue } from "@/composables/useQueue";

interface EditorState {
  content: string;
  cursor: number;
}

const { items, current, currentIndex, add, undo, redo, canUndo, canRedo, goto, clear } = useHistoryQueue<EditorState>(50);

// 添加状态
add({ content: "", cursor: 0 });
add({ content: "Hello", cursor: 5 });
add({ content: "Hello World", cursor: 11 });

// 撤销
if (canUndo.value) {
  const prevState = undo();
  console.log(prevState); // { content: 'Hello', cursor: 5 }
}

// 重做
if (canRedo.value) {
  const nextState = redo();
  console.log(nextState); // { content: 'Hello World', cursor: 11 }
}

// 跳转到指定历史
goto(0); // 跳到最早的状态

// 当前状态
console.log(current.value);
```

### 环形缓冲区

```typescript
import { useRingBuffer } from "@/composables/useQueue";

// 固定大小的缓冲区
const { items, push, shift, peek, peekLast, size, isFull, toArray } = useRingBuffer<number>(5);

// 添加数据
push(1, 2, 3, 4, 5); // 缓冲区已满

// 继续添加会移除最旧的
push(6); // 1 被移除，[2, 3, 4, 5, 6]
push(7, 8); // 2, 3 被移除，[4, 5, 6, 7, 8]

// 读取
shift(); // 4
peek(); // 5（第一个）
peekLast(); // 8（最后一个）

// 转为数组
const arr = toArray(); // [5, 6, 7, 8]
```

## API

### useQueue

| 选项        | 类型    | 默认值   | 说明             |
| ----------- | ------- | -------- | ---------------- |
| maxSize     | number  | Infinity | 最大容量         |
| prioritized | boolean | false    | 是否按优先级排序 |

| 返回值  | 类型                      | 说明       |
| ------- | ------------------------- | ---------- |
| items   | Ref\<QueueItem[]\>        | 队列项     |
| size    | ComputedRef\<number\>     | 队列长度   |
| isEmpty | ComputedRef\<boolean\>    | 是否为空   |
| isFull  | ComputedRef\<boolean\>    | 是否已满   |
| enqueue | `(data, priority?) => id` | 入队       |
| dequeue | `() => T`                 | 出队       |
| peek    | `() => T`                 | 查看队首   |
| clear   | `() => void`              | 清空       |
| remove  | `(id) => boolean`         | 移除指定项 |
| find    | `(predicate) => T`        | 查找       |
| has     | `(id) => boolean`         | 是否包含   |

### useTaskQueue

| 选项        | 类型    | 默认值 | 说明             |
| ----------- | ------- | ------ | ---------------- |
| concurrency | number  | 1      | 并发数           |
| interval    | number  | 0      | 任务间隔（毫秒） |
| retries     | number  | 0      | 重试次数         |
| retryDelay  | number  | 1000   | 重试延迟         |
| autoStart   | boolean | false  | 自动开始         |

| 返回值      | 类型            | 说明       |
| ----------- | --------------- | ---------- |
| tasks       | Ref\<Task[]\>   | 所有任务   |
| isRunning   | Ref\<boolean\>  | 是否运行中 |
| isPaused    | Ref\<boolean\>  | 是否暂停   |
| add         | `(fn) => id`    | 添加任务   |
| start       | `() => void`    | 开始处理   |
| pause       | `() => void`    | 暂停       |
| resume      | `() => void`    | 恢复       |
| retryFailed | `() => void`    | 重试失败   |
| waitAll     | `() => Promise` | 等待完成   |

### useNotificationQueue

| 选项            | 类型   | 默认值 | 说明         |
| --------------- | ------ | ------ | ------------ |
| maxVisible      | number | 5      | 最大显示数   |
| defaultDuration | number | 3000   | 默认持续时间 |
| position        | string | 'top'  | 新通知位置   |

| 返回值                     | 类型                   | 说明     |
| -------------------------- | ---------------------- | -------- |
| notifications              | Ref\<Notification[]\>  | 通知列表 |
| add                        | `(notification) => id` | 添加通知 |
| remove                     | `(id) => void`         | 移除通知 |
| clear                      | `() => void`           | 清空所有 |
| info/success/warning/error | `(msg, title?) => id`  | 快捷方法 |

## 代码位置

```
web/src/
└── composables/
    └── useQueue.ts    # Queue Composable
```
