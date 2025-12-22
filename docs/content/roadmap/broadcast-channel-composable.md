# BroadcastChannel Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:30+4`
- [已实现功能](#已实现功能) `:34+10`
  - [通信管理](#通信管理) `:36+8`
- [使用方式](#使用方式) `:44+144`
  - [基础用法](#基础用法) `:46+20`
  - [状态同步](#状态同步) `:66+36`
  - [Leader 选举](#leader-选举) `:102+34`
  - [消息传递](#消息传递) `:136+32`
  - [JSON 格式广播](#json-格式广播) `:168+20`
- [API](#api) `:188+57`
  - [useBroadcastChannel](#usebroadcastchannel) `:190+12`
  - [useTabSync](#usetabsync) `:202+16`
  - [useTabLeader](#usetableader) `:218+16`
  - [useTabMessenger](#usetabmessenger) `:234+11`
- [使用场景](#使用场景) `:245+10`
- [代码位置](#代码位置) `:255+7`

<!--TOC-->

## 需求背景

前端需要实现跨标签页通信，用于状态同步、Leader 选举等场景。

## 已实现功能

### 通信管理

- `useBroadcastChannel` - 基础广播频道通信
- `useBroadcastChannelJSON` - JSON 格式广播
- `useTabSync` - 标签页状态同步
- `useTabLeader` - 标签页 Leader 选举
- `useTabMessenger` - 标签页消息传递

## 使用方式

### 基础用法

```typescript
import { useBroadcastChannel } from "@/composables/useBroadcastChannel";

// 在多个标签页中使用相同的频道名
const { data, post, isSupported, close } = useBroadcastChannel<Message>("my-channel");

// 发送消息
post({ type: "update", payload: { id: 1 } });

// 监听消息
watch(data, (newData) => {
  console.log("收到消息:", newData);
});

// 关闭频道
close();
```

### 状态同步

```typescript
import { useTabSync } from "@/composables/useBroadcastChannel";

interface UserState {
  name: string;
  theme: "light" | "dark";
  lastActivity: number;
}

const { state, setState, mergeState, requestSync } = useTabSync<UserState>("user-state", {
  initialState: {
    name: "",
    theme: "light",
    lastActivity: Date.now(),
  },
  onSync: (state, source) => {
    console.log(`状态更新来源: ${source}`, state);
  },
});

// 完整更新状态（同步到所有标签页）
setState({
  name: "John",
  theme: "dark",
  lastActivity: Date.now(),
});

// 部分更新状态
mergeState({ theme: "dark" });

// 请求其他标签页同步状态
requestSync();
```

### Leader 选举

```typescript
import { useTabLeader } from "@/composables/useBroadcastChannel";

const { isLeader, resign, elect } = useTabLeader("app-leader", {
  heartbeatInterval: 1000, // 心跳间隔
  heartbeatTimeout: 3000, // 心跳超时

  onBecomeLeader: () => {
    // 成为 Leader 后执行的操作
    console.log("成为 Leader");
    startBackgroundSync();
  },

  onLoseLeadership: () => {
    // 失去 Leader 后执行的操作
    console.log("失去 Leader");
    stopBackgroundSync();
  },
});

// 只在 Leader 标签页执行昂贵操作
if (isLeader.value) {
  performExpensiveOperation();
}

// 主动放弃 Leader
resign();

// 尝试成为 Leader
elect();
```

### 消息传递

```typescript
import { useTabMessenger } from "@/composables/useBroadcastChannel";

const messenger = useTabMessenger<unknown>("app-messenger");

// 注册消息处理函数
messenger.on("userLogout", (data, source) => {
  console.log(`标签页 ${source} 触发登出`);
  router.push("/login");
});

messenger.on("themeChange", (data) => {
  applyTheme(data.theme);
});

// 广播消息
messenger.broadcast("userLogout", { userId: 123 });

// 移除处理函数
messenger.off("themeChange");

// 请求-响应模式
try {
  const result = await messenger.request("getData", { key: "user" }, 5000);
  console.log("响应:", result);
} catch (error) {
  console.log("请求超时或失败");
}
```

### JSON 格式广播

```typescript
import { useBroadcastChannelJSON } from "@/composables/useBroadcastChannel";

interface SyncMessage {
  type: "sync" | "update";
  data: unknown;
}

const { data, post } = useBroadcastChannelJSON<SyncMessage>("sync-channel", {
  onMessage: (data) => {
    console.log("收到:", data.type, data.data);
  },
});

// 发送 JSON 消息
post({ type: "sync", data: { foo: "bar" } });
```

## API

### useBroadcastChannel

| 返回值      | 类型                             | 说明           |
| ----------- | -------------------------------- | -------------- |
| channel     | Ref\<BroadcastChannel \| null\>  | 频道实例       |
| isSupported | ComputedRef\<boolean\>           | 是否支持       |
| isClosed    | Ref\<boolean\>                   | 是否已关闭     |
| data        | Ref\<T \| null\>                 | 最后接收的数据 |
| error       | Ref\<Event \| null\>             | 错误信息       |
| post        | `(data: P) => void             ` | 发送消息       |
| close       | `() => void                    ` | 关闭频道       |

### useTabSync

| 选项         | 类型    | 默认值 | 说明         |
| ------------ | ------- | ------ | ------------ |
| initialState | T       | -      | 初始状态     |
| immediate    | boolean | true   | 是否立即同步 |
| onSync       | func    | -      | 状态变化回调 |

| 返回值      | 类型                              | 说明           |
| ----------- | --------------------------------- | -------------- |
| state       | Ref\<T\>                          | 同步状态       |
| isSupported | ComputedRef\<boolean\>            | 是否支持       |
| setState    | `(newState: T) => void  `         | 更新状态并广播 |
| mergeState  | `(partial: Partial\<T\>) => void` | 合并状态       |
| requestSync | `() => void             `         | 请求同步       |

### useTabLeader

| 选项              | 类型   | 默认值 | 说明             |
| ----------------- | ------ | ------ | ---------------- |
| heartbeatInterval | number | 1000   | 心跳间隔（毫秒） |
| heartbeatTimeout  | number | 3000   | 心跳超时（毫秒） |
| onBecomeLeader    | func   | -      | 成为 Leader 回调 |
| onLoseLeadership  | func   | -      | 失去 Leader 回调 |

| 返回值      | 类型                     | 说明            |
| ----------- | ------------------------ | --------------- |
| isLeader    | Ref\<boolean\>           | 是否是 Leader   |
| isSupported | ComputedRef\<boolean\>   | 是否支持        |
| resign      | `() => void            ` | 放弃 Leader     |
| elect       | `() => void            ` | 尝试成为 Leader |

### useTabMessenger

| 返回值      | 类型                                | 说明          |
| ----------- | ----------------------------------- | ------------- |
| isSupported | ComputedRef\<boolean\>              | 是否支持      |
| tabId       | string                              | 当前标签页 ID |
| broadcast   | `(type, data) => void          `    | 广播消息      |
| on          | `(type, handler) => void       `    | 注册处理函数  |
| off         | `(type, handler?) => void      `    | 移除处理函数  |
| request     | `(type, data, timeout?) => Promise` | 请求响应      |

## 使用场景

| 场景       | 推荐方案        | 说明                 |
| ---------- | --------------- | -------------------- |
| 用户登出   | useTabMessenger | 广播登出事件         |
| 主题同步   | useTabSync      | 同步主题状态         |
| 购物车同步 | useTabSync      | 同步购物车数据       |
| 后台任务   | useTabLeader    | 只在一个标签页执行   |
| WebSocket  | useTabLeader    | 只在 Leader 维护连接 |

## 代码位置

```
web/src/
└── composables/
    └── useBroadcastChannel.ts    # BroadcastChannel Composable
```
