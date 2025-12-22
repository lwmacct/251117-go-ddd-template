# SWR Composable

> **状态**: ✅ 已完成
> **优先级**: 高
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:32+4`
- [已实现功能](#已实现功能) `:36+16`
  - [数据获取](#数据获取) `:38+6`
  - [缓存管理](#缓存管理) `:44+8`
- [SWR 原理](#swr-原理) `:52+8`
- [使用方式](#使用方式) `:60+155`
  - [基础用法](#基础用法) `:62+13`
  - [自动重新验证](#自动重新验证) `:75+18`
  - [手动更新数据](#手动更新数据) `:93+16`
  - [错误处理和重试](#错误处理和重试) `:109+22`
  - [数据修改](#数据修改) `:131+35`
  - [无限加载](#无限加载) `:166+28`
  - [缓存管理](#缓存管理-1) `:194+21`
- [API](#api) `:215+59`
  - [useSWR](#useswr) `:217+27`
  - [useSWRMutation](#useswrmutation) `:244+17`
  - [useSWRInfinite](#useswrinfinite) `:261+13`
- [代码位置](#代码位置) `:274+7`

<!--TOC-->

## 需求背景

前端需要实现 Stale-While-Revalidate (SWR) 数据获取策略，提供更好的用户体验和数据一致性。

## 已实现功能

### 数据获取

- `useSWR` - 基础 SWR 数据获取
- `useSWRMutation` - 用于修改数据的 SWR
- `useSWRInfinite` - 无限加载的 SWR

### 缓存管理

- `clearSWRCache` - 清除所有缓存
- `deleteSWRCache` - 清除特定缓存
- `getSWRCache` - 获取缓存
- `setSWRCache` - 设置缓存
- `revalidateSWR` - 全局重新验证

## SWR 原理

```
1. 显示缓存数据（stale）→ 用户立即看到内容
2. 后台重新获取（revalidate）→ 更新最新数据
3. 更新缓存和 UI → 保持数据新鲜
```

## 使用方式

### 基础用法

```typescript
import { useSWR } from "@/composables/useSWR";

const { data, error, isLoading, isValidating, revalidate, mutate } = useSWR("users", () => fetch("/api/users").then((r) => r.json()));

// data: 用户列表数据
// error: 错误信息
// isLoading: 首次加载中
// isValidating: 后台验证中
```

### 自动重新验证

```typescript
const { data } = useSWR("users", fetcher, {
  // 每 30 秒自动刷新
  revalidateInterval: 30000,

  // 页面聚焦时刷新
  revalidateOnFocus: true,

  // 网络重连时刷新
  revalidateOnReconnect: true,

  // 挂载时刷新
  revalidateOnMount: true,
});
```

### 手动更新数据

```typescript
const { data, mutate, revalidate } = useSWR("users", fetcher);

// 乐观更新（立即更新 UI，不重新获取）
mutate([...data.value, newUser]);

// 乐观更新 + 重新获取
mutate([...data.value, newUser]);
await revalidate();

// 使用函数更新
mutate((prev) => prev?.filter((u) => u.id !== deletedId) ?? []);
```

### 错误处理和重试

```typescript
const { data, error } = useSWR("users", fetcher, {
  // 错误重试 3 次
  errorRetryCount: 3,

  // 重试间隔递增
  errorRetryInterval: 5000,

  // 错误回调
  onError: (error) => {
    console.error("获取失败:", error);
  },

  // 成功回调
  onSuccess: (data) => {
    console.log("获取成功:", data);
  },
});
```

### 数据修改

```typescript
import { useSWRMutation } from "@/composables/useSWR";

const { trigger, isMutating, error } = useSWRMutation(
  "users",
  (user: CreateUserDto) =>
    fetch("/api/users", {
      method: "POST",
      body: JSON.stringify(user),
    }).then((r) => r.json()),
  {
    // 乐观更新
    optimisticData: (current, newUser) => [...(current ?? []), newUser],

    // 错误时回滚
    rollbackOnError: true,

    onSuccess: (data) => {
      console.log("创建成功:", data);
    },
  },
);

// 触发创建
const handleCreate = async () => {
  try {
    await trigger({ name: "John", email: "john@example.com" });
  } catch (error) {
    // 处理错误
  }
};
```

### 无限加载

```typescript
import { useSWRInfinite } from "@/composables/useSWR";

const { data, loadMore, hasMore, isLoadingMore, reset } = useSWRInfinite(
  (pageIndex) => `users-page-${pageIndex}`,
  (pageIndex) => fetch(`/api/users?page=${pageIndex}&size=10`).then((r) => r.json()),
  {
    pageSize: 10,
    hasMore: (pageData, pageIndex) => pageData.length === 10,
  },
);

// data 包含所有页的数据
console.log(data.value); // 所有用户

// 加载更多
const handleLoadMore = async () => {
  if (hasMore.value) {
    await loadMore();
  }
};

// 重置
reset();
```

### 缓存管理

```typescript
import { getSWRCache, setSWRCache, deleteSWRCache, clearSWRCache, revalidateSWR } from "@/composables/useSWR";

// 获取缓存
const users = getSWRCache<User[]>("users");

// 设置缓存（会通知所有订阅者）
setSWRCache("users", updatedUsers);

// 删除缓存
deleteSWRCache("users");

// 清除所有缓存
clearSWRCache();

// 触发重新验证
revalidateSWR("users");
```

## API

### useSWR

| 选项                  | 类型    | 默认值 | 说明           |
| --------------------- | ------- | ------ | -------------- |
| initialData           | T       | -      | 初始数据       |
| immediate             | boolean | true   | 是否立即获取   |
| revalidateInterval    | number  | 0      | 重新验证间隔   |
| revalidateOnFocus     | boolean | true   | 聚焦时重新验证 |
| revalidateOnReconnect | boolean | true   | 重连时重新验证 |
| revalidateOnMount     | boolean | true   | 挂载时重新验证 |
| dedupingInterval      | number  | 2000   | 去重间隔       |
| errorRetryCount       | number  | 3      | 错误重试次数   |
| errorRetryInterval    | number  | 5000   | 错误重试间隔   |
| keepPreviousData      | boolean | false  | 保持之前的数据 |
| onSuccess             | func    | -      | 成功回调       |
| onError               | func    | -      | 错误回调       |

| 返回值       | 类型                      | 说明         |
| ------------ | ------------------------- | ------------ |
| data         | Ref\<T \| null\>          | 数据         |
| error        | Ref\<Error \| null\>      | 错误         |
| status       | Ref\<SWRStatus\>          | 状态         |
| isValidating | Ref\<boolean\>            | 是否正在验证 |
| isLoading    | ComputedRef\<boolean\>    | 是否首次加载 |
| mutate       | `(data?) => Promise     ` | 更新数据     |
| revalidate   | `() => Promise          ` | 重新获取     |

### useSWRMutation

| 选项            | 类型    | 默认值 | 说明         |
| --------------- | ------- | ------ | ------------ |
| optimisticData  | func    | -      | 乐观更新函数 |
| rollbackOnError | boolean | true   | 错误时回滚   |
| onSuccess       | func    | -      | 成功回调     |
| onError         | func    | -      | 错误回调     |

| 返回值     | 类型                   | 说明         |
| ---------- | ---------------------- | ------------ |
| data       | Ref\<T \| null\>       | 数据         |
| error      | Ref\<Error \| null\>   | 错误         |
| isMutating | Ref\<boolean\>         | 是否正在执行 |
| trigger    | `(arg) => Promise    ` | 触发修改     |
| reset      | `() => void          ` | 重置状态     |

### useSWRInfinite

| 返回值        | 类型             | 说明             |
| ------------- | ---------------- | ---------------- |
| data          | Ref\<T[]\>       | 所有页数据       |
| error         | Ref\<Error\>     | 错误             |
| isValidating  | Ref\<boolean\>   | 是否正在验证     |
| isLoadingMore | Ref\<boolean\>   | 是否正在加载更多 |
| hasMore       | Ref\<boolean\>   | 是否有更多       |
| size          | Ref\<number\>    | 当前页数         |
| loadMore      | `() => Promise ` | 加载更多         |
| reset         | `() => void    ` | 重置             |

## 代码位置

```
web/src/
└── composables/
    └── useSWR.ts    # SWR Composable
```
