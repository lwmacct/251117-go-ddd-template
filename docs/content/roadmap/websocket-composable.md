# WebSocket Composable

> **状态**: ✅ 已完成
> **优先级**: 高
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:31+4`
- [已实现功能](#已实现功能) `:35+16`
  - [连接管理](#连接管理) `:37+7`
  - [特性](#特性) `:44+7`
- [使用方式](#使用方式) `:51+146`
  - [基础用法](#基础用法) `:53+21`
  - [自动重连](#自动重连) `:74+20`
  - [心跳检测](#心跳检测) `:94+12`
  - [JSON WebSocket](#json-websocket) `:106+25`
  - [二进制 WebSocket](#二进制-websocket) `:131+21`
  - [多连接管理](#多连接管理) `:152+30`
  - [消息缓冲](#消息缓冲) `:182+15`
- [API](#api) `:197+47`
  - [useWebSocket](#usewebsocket) `:199+27`
  - [AutoReconnect 配置](#autoreconnect-配置) `:226+10`
  - [Heartbeat 配置](#heartbeat-配置) `:236+8`
- [代码位置](#代码位置) `:244+7`

<!--TOC-->

## 需求背景

前端需要响应式管理 WebSocket 连接，支持自动重连、心跳检测、消息缓冲等功能，用于实时通信场景。

## 已实现功能

### 连接管理

- `useWebSocket` - 基础 WebSocket 连接管理
- `useWebSocketJSON` - JSON 格式 WebSocket
- `useWebSocketBinary` - 二进制 WebSocket
- `createWebSocketManager` - 多连接管理器

### 特性

- 自动重连（指数退避）
- 心跳检测
- 消息缓冲（连接前发送的消息自动缓存）
- 页面可见性变化时自动重连

## 使用方式

### 基础用法

```typescript
import { useWebSocket } from "@/composables/useWebSocket";

const { data, send, isConnected, open, close, status } = useWebSocket("ws://localhost:8080/ws");

// 发送消息
send("Hello Server");

// 监听数据
watch(data, (newData) => {
  console.log("收到消息:", newData);
});

// 检查连接状态
if (isConnected.value) {
  send("Connected message");
}
```

### 自动重连

```typescript
const { data, retryCount } = useWebSocket("ws://localhost:8080/ws", {
  autoReconnect: {
    retries: 5, // 最大重连次数
    delay: 1000, // 初始延迟
    multiplier: 2, // 延迟递增因子
    maxDelay: 30000, // 最大延迟
    onVisibilityChange: true, // 页面可见时重连
  },
  onReconnect: (count) => {
    console.log(`正在重连... 第 ${count} 次`);
  },
  onFailed: () => {
    console.log("重连失败，已达最大次数");
  },
});
```

### 心跳检测

```typescript
const { isConnected } = useWebSocket("ws://localhost:8080/ws", {
  heartbeat: {
    message: "ping", // 心跳消息
    interval: 30000, // 心跳间隔（30秒）
    timeout: 10000, // 心跳超时（10秒）
  },
});
```

### JSON WebSocket

```typescript
import { useWebSocketJSON } from "@/composables/useWebSocket";

interface ChatMessage {
  type: "chat" | "system";
  content: string;
  timestamp: number;
}

const { data, send, isConnected } = useWebSocketJSON<ChatMessage>("ws://localhost:8080/chat", {
  onMessage: (ws, data) => {
    console.log("收到消息:", data.type, data.content);
  },
});

// 发送 JSON 对象（自动序列化）
send({
  type: "chat",
  content: "Hello!",
  timestamp: Date.now(),
});
```

### 二进制 WebSocket

```typescript
import { useWebSocketBinary } from "@/composables/useWebSocket";

const { data, send } = useWebSocketBinary("ws://localhost:8080/binary", {
  binaryType: "arraybuffer",
  onMessage: (ws, buffer) => {
    // 处理二进制数据
    const view = new DataView(buffer as ArrayBuffer);
    console.log("收到数据:", view.getInt32(0));
  },
});

// 发送二进制数据
const buffer = new ArrayBuffer(8);
const view = new DataView(buffer);
view.setInt32(0, 12345);
send(buffer);
```

### 多连接管理

```typescript
import { createWebSocketManager } from "@/composables/useWebSocket";

const wsManager = createWebSocketManager();

// 创建多个连接
wsManager.create("chat", "ws://localhost:8080/chat", {
  autoReconnect: true,
});

wsManager.create("notifications", "ws://localhost:8080/notifications", {
  heartbeat: true,
});

// 获取并使用连接
const chat = wsManager.get("chat");
chat?.send("Hello from chat");

// 广播消息到所有连接
wsManager.broadcast("ping");

// 关闭特定连接
wsManager.close("chat");

// 关闭所有连接
wsManager.closeAll();
```

### 消息缓冲

```typescript
const { send, isConnected } = useWebSocket("ws://localhost:8080/ws", {
  immediate: false, // 不立即连接
});

// 连接前发送的消息会被缓冲
send("Message 1"); // 缓冲
send("Message 2"); // 缓冲

// 连接后自动发送缓冲的消息
open();
```

## API

### useWebSocket

| 选项          | 类型                       | 默认值 | 说明         |
| ------------- | -------------------------- | ------ | ------------ |
| immediate     | boolean                    | true   | 是否立即连接 |
| autoReconnect | boolean \| object          | false  | 自动重连配置 |
| heartbeat     | boolean \| object          | false  | 心跳配置     |
| protocols     | string \| string[]         | -      | 子协议       |
| onOpen        | `(ws, event) => void     ` | -      | 连接打开回调 |
| onMessage     | `(ws, event) => void     ` | -      | 消息回调     |
| onClose       | `(ws, event) => void     ` | -      | 连接关闭回调 |
| onError       | `(ws, event) => void     ` | -      | 错误回调     |
| onReconnect   | `(retries) => void       ` | -      | 重连回调     |
| onFailed      | `() => void              ` | -      | 连接失败回调 |

| 返回值      | 类型                            | 说明           |
| ----------- | ------------------------------- | -------------- |
| ws          | Ref\<WebSocket \| null\>        | WebSocket 实例 |
| status      | Ref\<WebSocketStatus\>          | 连接状态       |
| isConnected | ComputedRef\<boolean\>          | 是否已连接     |
| data        | Ref\<T \| null\>                | 最后接收的数据 |
| error       | Ref\<Event \| null\>            | 错误信息       |
| retryCount  | Ref\<number\>                   | 重连次数       |
| open        | `() => void               `     | 打开连接       |
| close       | `(code?, reason?) => void `     | 关闭连接       |
| send        | `(data, useBuffer?) => boolean` | 发送消息       |

### AutoReconnect 配置

| 选项               | 类型    | 默认值 | 说明             |
| ------------------ | ------- | ------ | ---------------- |
| retries            | number  | 3      | 最大重连次数     |
| delay              | number  | 1000   | 初始延迟（毫秒） |
| multiplier         | number  | 2      | 延迟递增因子     |
| maxDelay           | number  | 30000  | 最大延迟（毫秒） |
| onVisibilityChange | boolean | true   | 页面可见时重连   |

### Heartbeat 配置

| 选项     | 类型                          | 默认值 | 说明             |
| -------- | ----------------------------- | ------ | ---------------- |
| message  | string \| ArrayBuffer \| Blob | 'ping' | 心跳消息         |
| interval | number                        | 30000  | 心跳间隔（毫秒） |
| timeout  | number                        | 10000  | 心跳超时（毫秒） |

## 代码位置

```
web/src/
└── composables/
    └── useWebSocket.ts    # WebSocket Composable
```
