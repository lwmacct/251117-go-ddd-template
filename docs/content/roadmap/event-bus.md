# 事件总线 Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:31+4`
- [已实现功能](#已实现功能) `:35+30`
  - [createEventBus](#createeventbus) `:37+7`
  - [useEventBus](#useeventbus) `:44+6`
  - [useEventListener](#useeventlistener) `:50+5`
  - [useEventValue](#useeventvalue) `:55+5`
  - [appEventBus](#appeventbus) `:60+5`
- [使用方式](#使用方式) `:65+74`
  - [定义事件类型](#定义事件类型) `:67+14`
  - [创建事件总线](#创建事件总线) `:81+9`
  - [在组件中使用](#在组件中使用) `:90+21`
  - [响应式事件值](#响应式事件值) `:111+17`
  - [一次性订阅](#一次性订阅) `:128+11`
- [API](#api) `:139+22`
  - [createEventBus 返回值](#createeventbus-返回值) `:141+11`
  - [预定义事件类型](#预定义事件类型) `:152+9`
- [代码位置](#代码位置) `:161+7`

<!--TOC-->

## 需求背景

需要为非父子组件间提供通信机制，实现解耦的事件驱动架构。

## 已实现功能

### createEventBus

- 创建类型安全的事件总线
- 订阅/取消订阅
- 触发事件
- 一次性订阅

### useEventBus

- Vue Composable 封装
- 组件卸载自动清理
- 类型安全

### useEventListener

- 特定事件监听
- 简化单事件订阅

### useEventValue

- 响应式事件值
- 自动更新为最新 payload

### appEventBus

- 预定义应用事件总线
- 常用事件类型定义

## 使用方式

### 定义事件类型

```typescript
// types/events.ts
export interface AppEvents {
  "user:login": { userId: string; username: string };
  "user:logout": void;
  "notification:show": {
    message: string;
    type: "success" | "error";
  };
}
```

### 创建事件总线

```typescript
import { createEventBus } from "@/composables/useEventBus";
import type { AppEvents } from "@/types/events";

export const appEventBus = createEventBus<AppEvents>();
```

### 在组件中使用

```typescript
import { useEventBus } from "@/composables/useEventBus";
import type { AppEvents } from "@/types/events";

// 组件 A：发送事件
const { emit } = useEventBus<AppEvents>();

function handleLogin() {
  emit("user:login", { userId: "123", username: "John" });
}

// 组件 B：接收事件（自动清理）
const { on } = useEventBus<AppEvents>();

on("user:login", (data) => {
  console.log("用户登录:", data.username);
});
```

### 响应式事件值

```typescript
import { useEventValue } from "@/composables/useEventBus";

// 最新的通知会自动更新到 value
const { value: notification } = useEventValue<AppEvents, "notification:show">("notification:show");
```

```vue
<template>
  <v-alert v-if="notification" :type="notification.type">
    {{ notification.message }}
  </v-alert>
</template>
```

### 一次性订阅

```typescript
const { once } = useEventBus<AppEvents>();

// 只触发一次
once("user:login", (data) => {
  showWelcomeMessage(data.username);
});
```

## API

### createEventBus 返回值

| 方法          | 类型                              | 说明       |
| ------------- | --------------------------------- | ---------- |
| on            | `(event, handler) => unsubscribe` | 订阅事件   |
| off           | `(event, handler?) => void`       | 取消订阅   |
| emit          | `(event, payload) => void`        | 触发事件   |
| once          | `(event, handler) => unsubscribe` | 一次性订阅 |
| clear         | `() => void`                      | 清除所有   |
| listenerCount | `(event) => number`               | 监听器数量 |

### 预定义事件类型

| 事件              | Payload              | 说明     |
| ----------------- | -------------------- | -------- |
| user:login        | { userId, username } | 用户登录 |
| user:logout       | void                 | 用户登出 |
| theme:change      | { theme }            | 主题变化 |
| notification:show | { message, type }    | 显示通知 |

## 代码位置

```
web/src/
└── composables/
    └── useEventBus.ts    # 事件总线
```
