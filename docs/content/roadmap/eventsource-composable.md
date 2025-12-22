# EventSource Composable

> **状态**: ✅ 已完成
> **优先级**: 高
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:31+4`
- [已实现功能](#已实现功能) `:35+16`
  - [连接管理](#连接管理) `:37+7`
  - [特性](#特性) `:44+7`
- [使用方式](#使用方式) `:51+117`
  - [基础用法](#基础用法) `:53+18`
  - [自动重连](#自动重连) `:71+19`
  - [监听命名事件](#监听命名事件) `:90+27`
  - [简化用法](#简化用法) `:117+16`
  - [多连接管理](#多连接管理) `:133+27`
  - [携带凭证](#携带凭证) `:160+8`
- [API](#api) `:168+46`
  - [useEventSource](#useeventsource) `:170+26`
  - [useEventSourceNamed 额外返回值](#useeventsourcenamed-额外返回值) `:196+9`
  - [AutoReconnect 配置](#autoreconnect-配置) `:205+9`
- [SSE vs WebSocket](#sse-vs-websocket) `:214+10`
- [代码位置](#代码位置) `:224+7`

<!--TOC-->

## 需求背景

前端需要响应式管理 Server-Sent Events (SSE) 连接，用于接收服务器推送的实时数据。

## 已实现功能

### 连接管理

- `useEventSource` - 基础 SSE 连接管理
- `useEventSourceNamed` - 监听命名事件
- `useServerSentEvents` - 简化的 SSE Hook
- `createEventSourceManager` - 多连接管理器

### 特性

- 自动重连（指数退避）
- 命名事件支持
- JSON 自动解析
- 事件 ID 追踪

## 使用方式

### 基础用法

```typescript
import { useEventSource } from "@/composables/useEventSource";

const { data, isConnected, status, open, close } = useEventSource("/api/events");

// 监听数据
watch(data, (newData) => {
  console.log("收到数据:", newData);
});

// 检查连接状态
if (isConnected.value) {
  console.log("SSE 已连接");
}
```

### 自动重连

```typescript
const { data, retryCount } = useEventSource("/api/events", {
  autoReconnect: {
    retries: 5, // 最大重连次数
    delay: 1000, // 初始延迟
    multiplier: 2, // 延迟递增因子
    maxDelay: 30000, // 最大延迟
  },
  onReconnect: (count) => {
    console.log(`正在重连... 第 ${count} 次`);
  },
  onFailed: () => {
    console.log("重连失败，已达最大次数");
  },
});
```

### 监听命名事件

```typescript
import { useEventSourceNamed } from "@/composables/useEventSource";

// 服务器发送: event: update\ndata: {...}\n\n
const { events, lastEvent, addEventListener, removeEventListener } = useEventSourceNamed("/api/events", {
  events: ["update", "delete", "create"],
});

// 获取特定事件的数据
const updateData = computed(() => events.value.get("update"));

// 监听最新事件
watch(lastEvent, (event) => {
  if (event) {
    console.log(`事件 ${event.name}:`, event.data);
  }
});

// 动态添加事件监听
addEventListener("custom-event");

// 移除事件监听
removeEventListener("update");
```

### 简化用法

```typescript
import { useServerSentEvents } from "@/composables/useEventSource";

// 直接获取响应式数据
const messages = useServerSentEvents<Message[]>("/api/messages");

// messages 会自动更新
watch(messages, (newMessages) => {
  if (newMessages) {
    displayMessages(newMessages);
  }
});
```

### 多连接管理

```typescript
import { createEventSourceManager } from "@/composables/useEventSource";

const sseManager = createEventSourceManager();

// 创建多个连接
sseManager.create("notifications", "/api/notifications", {
  autoReconnect: true,
});

sseManager.create("updates", "/api/updates");

// 获取并使用连接
const notifications = sseManager.get("notifications");
watch(notifications?.data, (data) => {
  showNotification(data);
});

// 关闭特定连接
sseManager.close("updates");

// 关闭所有连接
sseManager.closeAll();
```

### 携带凭证

```typescript
const { data } = useEventSource("/api/private-events", {
  withCredentials: true,
});
```

## API

### useEventSource

| 选项            | 类型                | 默认值 | 说明         |
| --------------- | ------------------- | ------ | ------------ |
| immediate       | boolean             | true   | 是否立即连接 |
| autoReconnect   | boolean \| object   | false  | 自动重连配置 |
| withCredentials | boolean             | false  | 是否携带凭证 |
| onOpen          | `(event) => void  ` | -      | 连接打开回调 |
| onMessage       | `(event) => void  ` | -      | 消息回调     |
| onError         | `(event) => void  ` | -      | 错误回调     |
| onReconnect     | `(retries) => void` | -      | 重连回调     |
| onFailed        | `() => void       ` | -      | 连接失败回调 |

| 返回值      | 类型                           | 说明             |
| ----------- | ------------------------------ | ---------------- |
| eventSource | Ref\<EventSource \| null\>     | EventSource 实例 |
| status      | Ref\<EventSourceStatus\>       | 连接状态         |
| isConnected | ComputedRef\<boolean\>         | 是否已连接       |
| data        | Ref\<T \| null\>               | 最后接收的数据   |
| event       | Ref\<string \| null\>          | 最后接收的事件   |
| lastEventId | Ref\<string \| null\>          | 最后的事件 ID    |
| error       | Ref\<Event \| null\>           | 错误信息         |
| retryCount  | Ref\<number\>                  | 重连次数         |
| open        | `() => void                  ` | 打开连接         |
| close       | `() => void                  ` | 关闭连接         |

### useEventSourceNamed 额外返回值

| 返回值              | 类型                            | 说明         |
| ------------------- | ------------------------------- | ------------ |
| events              | Ref\<Map\<string, T\>\>         | 所有事件数据 |
| lastEvent           | Ref\<NamedEventData \| null\>   | 最后命名事件 |
| addEventListener    | `(name: string) => void       ` | 添加事件监听 |
| removeEventListener | `(name: string) => void       ` | 移除事件监听 |

### AutoReconnect 配置

| 选项       | 类型   | 默认值 | 说明             |
| ---------- | ------ | ------ | ---------------- |
| retries    | number | 3      | 最大重连次数     |
| delay      | number | 1000   | 初始延迟（毫秒） |
| multiplier | number | 2      | 延迟递增因子     |
| maxDelay   | number | 30000  | 最大延迟（毫秒） |

## SSE vs WebSocket

| 特性     | SSE                   | WebSocket   |
| -------- | --------------------- | ----------- |
| 方向     | 单向（服务器→客户端） | 双向        |
| 协议     | HTTP                  | WS          |
| 重连     | 自动                  | 需手动实现  |
| 数据格式 | 文本                  | 文本/二进制 |
| 适用场景 | 实时通知、数据推送    | 聊天、游戏  |

## 代码位置

```
web/src/
└── composables/
    └── useEventSource.ts    # EventSource Composable
```
