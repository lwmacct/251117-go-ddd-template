# 限流工具

<!--TOC-->

- [需求背景](#需求背景) `:32:35`
- [已实现功能](#已实现功能) `:36:37`
  - [throttle](#throttle) `:38:43`
  - [debounce](#debounce) `:44:49`
  - [createRateLimiter](#createratelimiter) `:50:55`
  - [retry](#retry) `:56:61`
  - [withTimeout](#withtimeout) `:62:66`
  - [createDeduplicator](#creatededuplicator) `:67:71`
- [使用方式](#使用方式) `:72:73`
  - [节流](#节流) `:74:93`
  - [防抖](#防抖) `:94:111`
  - [限速器](#限速器) `:112:128`
  - [重试](#重试) `:129:143`
  - [超时](#超时) `:144:157`
  - [去重](#去重) `:158:170`
- [API](#api) `:171:172`
  - [throttle 选项](#throttle-选项) `:173:179`
  - [debounce 选项](#debounce-选项) `:180:187`
  - [retry 选项](#retry-选项) `:188:197`
- [代码位置](#代码位置) `:198:204`

<!--TOC-->

> **状态**: ✅ 已完成
> **优先级**: 高
> **完成日期**: 2024-11-30

## 需求背景

需要控制函数的执行频率，防止过度调用导致性能问题或 API 限制。

## 已实现功能

### throttle

- 节流函数
- 在指定时间内只执行一次
- 支持 leading/trailing 配置

### debounce

- 防抖函数
- 在停止调用后等待指定时间才执行
- 支持 maxWait 最大等待时间

### createRateLimiter

- 限速器
- 限制每秒/每分钟的调用次数
- 请求队列管理

### retry

- 重试函数
- 指数退避策略
- 自定义重试条件

### withTimeout

- 超时包装
- 为 Promise 添加超时限制

### createDeduplicator

- 去重执行器
- 相同 key 的调用只执行一次

## 使用方式

### 节流

```typescript
import { throttle } from "@/utils/throttle";

// 滚动事件节流
const handleScroll = throttle(() => {
  updateScrollPosition();
}, 200);

window.addEventListener("scroll", handleScroll);

// 带配置
const handleResize = throttle(() => recalculateLayout(), 300, { leading: true, trailing: false });

// 取消和立即执行
handleScroll.cancel(); // 取消待执行
handleScroll.flush(); // 立即执行
```

### 防抖

```typescript
import { debounce } from "@/utils/throttle";

// 搜索输入防抖
const handleSearch = debounce((query: string) => {
  search(query);
}, 300);

// 带最大等待时间
const handleInput = debounce(
  (value: string) => saveValue(value),
  500,
  { maxWait: 2000 }, // 最多等待 2 秒
);
```

### 限速器

```typescript
import { createRateLimiter } from "@/utils/throttle";

// 每秒最多 10 个请求
const limiter = createRateLimiter({
  maxRequests: 10,
  interval: 1000,
});

// 使用限速器发请求
for (const item of items) {
  await limiter.execute(() => api.process(item));
}
```

### 重试

```typescript
import { retry } from "@/utils/throttle";

const result = await retry(() => fetchData(), {
  maxRetries: 3,
  delay: 1000,
  factor: 2, // 指数退避
  onRetry: (err, attempt) => {
    console.log(`重试 ${attempt + 1} 次`);
  },
});
```

### 超时

```typescript
import { withTimeout } from "@/utils/throttle";

try {
  const result = await withTimeout(fetchData(), 5000, "请求超时");
} catch (err) {
  if (err.message === "请求超时") {
    // 处理超时
  }
}
```

### 去重

```typescript
import { createDeduplicator } from "@/utils/throttle";

const dedup = createDeduplicator();

// 同时发起的请求只执行一次
const user1 = await dedup.execute("user-1", () => fetchUser(1));
const user2 = await dedup.execute("user-1", () => fetchUser(1));
// user1 === user2，只发了一次请求
```

## API

### throttle 选项

| 选项     | 类型    | 默认值 | 说明             |
| -------- | ------- | ------ | ---------------- |
| leading  | boolean | true   | 是否在开始时调用 |
| trailing | boolean | true   | 是否在结束时调用 |

### debounce 选项

| 选项     | 类型    | 默认值 | 说明             |
| -------- | ------- | ------ | ---------------- |
| leading  | boolean | false  | 是否在开始时调用 |
| trailing | boolean | true   | 是否在结束时调用 |
| maxWait  | number  | -      | 最大等待时间     |

### retry 选项

| 选项        | 类型     | 默认值 | 说明         |
| ----------- | -------- | ------ | ------------ |
| maxRetries  | number   | 3      | 最大重试次数 |
| delay       | number   | 1000   | 重试延迟     |
| factor      | number   | 2      | 延迟增长因子 |
| shouldRetry | Function | -      | 是否重试判断 |
| onRetry     | Function | -      | 重试回调     |

## 代码位置

```
web/src/
└── utils/
    └── throttle.ts    # 限流工具
```
