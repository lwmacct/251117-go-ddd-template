# Fetch Composable

> **状态**: ✅ 已完成
> **优先级**: 高
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:31+4`
- [已实现功能](#已实现功能) `:35+16`
  - [请求管理](#请求管理) `:37+8`
  - [缓存](#缓存) `:45+6`
- [使用方式](#使用方式) `:51+183`
  - [基础用法](#基础用法) `:53+22`
  - [请求配置](#请求配置) `:75+34`
  - [缓存](#缓存-1) `:109+12`
  - [创建 API 实例](#创建-api-实例) `:121+38`
  - [懒加载请求](#懒加载请求) `:159+13`
  - [分页请求](#分页请求) `:172+24`
  - [无限加载](#无限加载) `:196+20`
  - [POST 请求](#post-请求) `:216+18`
- [API](#api) `:234+50`
  - [useFetch](#usefetch) `:236+32`
  - [createFetch](#createfetch) `:268+16`
- [代码位置](#代码位置) `:284+7`

<!--TOC-->

## 需求背景

前端需要响应式管理 HTTP 请求，支持缓存、重试、分页、无限加载等功能。

## 已实现功能

### 请求管理

- `useFetch` - 基础 HTTP 请求
- `createFetch` - 创建带默认配置的 fetch 实例
- `useLazyFetch` - 懒加载请求
- `usePaginatedFetch` - 分页请求
- `useInfiniteFetch` - 无限加载请求

### 缓存

- `clearFetchCache` - 清除所有缓存
- `deleteFetchCache` - 清除特定缓存
- `getFetchCacheSize` - 获取缓存大小

## 使用方式

### 基础用法

```typescript
import { useFetch } from "@/composables/useFetch";

const { data, isLoading, error, execute, abort, retry } = useFetch<User[]>("/api/users");

// 监听数据
watch(data, (users) => {
  console.log("用户列表:", users);
});

// 手动执行
await execute();

// 取消请求
abort();

// 重试
await retry();
```

### 请求配置

```typescript
const { data } = useFetch<User>("/api/user/1", {
  immediate: true, // 立即请求
  timeout: 10000, // 10秒超时
  retry: 3, // 重试3次
  retryDelay: 1000, // 重试延迟1秒
  responseType: "json", // 响应类型

  // 请求前拦截
  beforeFetch: ({ url, options }) => {
    options.headers = {
      ...options.headers,
      Authorization: `Bearer ${token}`,
    };
    return { url, options };
  },

  // 响应后处理
  afterFetch: ({ data, response }) => {
    console.log("状态码:", response.status);
    return data;
  },

  // 错误处理
  onFetchError: ({ error, response }) => {
    if (response?.status === 401) {
      router.push("/login");
    }
  },
});
```

### 缓存

```typescript
const { data, execute } = useFetch<Config>("/api/config", {
  cacheKey: "app-config",
  cacheTime: 60000, // 缓存1分钟
});

// 第二次请求会使用缓存
await execute();
```

### 创建 API 实例

```typescript
import { createFetch } from "@/composables/useFetch";

const api = createFetch({
  baseUrl: "/api",
  options: {
    timeout: 10000,
  },
  interceptors: {
    request: ({ url, options }) => {
      options.headers = {
        ...options.headers,
        Authorization: `Bearer ${getToken()}`,
      };
      return { url, options };
    },
    response: ({ data }) => {
      // 统一处理响应
      return data.data; // 假设后端返回 { code, data, message }
    },
    error: ({ error, response }) => {
      if (response?.status === 401) {
        logout();
      }
    },
  },
});

// 使用
const { data: users } = api.get<User[]>("/users");
const { data: user } = api.post<User>("/users", { name: "John" });
const { data } = api.put<User>("/users/1", { name: "Jane" });
const { data } = api.patch<User>("/users/1", { age: 26 });
const { data } = api.delete<void>("/users/1");
```

### 懒加载请求

```typescript
import { useLazyFetch } from "@/composables/useFetch";

const { execute, data, isLoading } = useLazyFetch<UserDetails>("/api/user/1");

// 需要时手动执行
const handleClick = async () => {
  await execute();
};
```

### 分页请求

```typescript
import { usePaginatedFetch } from "@/composables/useFetch";

const { data, page, pageSize, hasMore, loadMore, refresh, isLoading } = usePaginatedFetch<User[]>("/api/users", {
  pageSize: 10,
  pageParam: "page",
  pageSizeParam: "size",
});

// 加载更多
const handleLoadMore = async () => {
  if (hasMore.value) {
    await loadMore();
  }
};

// 刷新
const handleRefresh = async () => {
  await refresh();
};
```

### 无限加载

```typescript
import { useInfiniteFetch } from "@/composables/useFetch";

const { data, loadMore, hasMore, isLoading, reset } = useInfiniteFetch<User>("/api/users", {
  pageSize: 20,
  merge: (prev, next) => [...prev, ...next], // 自定义合并
});

// data 包含所有已加载的数据
console.log(data.value); // 所有用户

// 加载更多时数据会累积
await loadMore();

// 重置
reset();
```

### POST 请求

```typescript
const { data, execute, isLoading } = useFetch<User>("/api/users", {
  method: "POST",
  body: JSON.stringify({ name: "John", email: "john@example.com" }),
  headers: {
    "Content-Type": "application/json",
  },
  immediate: false,
});

// 提交表单
const handleSubmit = async () => {
  await execute();
};
```

## API

### useFetch

| 选项           | 类型     | 默认值 | 说明             |
| -------------- | -------- | ------ | ---------------- |
| immediate      | boolean  | true   | 是否立即请求     |
| timeout        | number   | 30000  | 超时时间（毫秒） |
| retry          | number   | 0      | 重试次数         |
| retryDelay     | number   | 1000   | 重试延迟（毫秒） |
| responseType   | string   | 'json' | 响应类型         |
| cacheKey       | string   | -      | 缓存键           |
| cacheTime      | number   | 0      | 缓存时间（毫秒） |
| debounce       | number   | 0      | 防抖时间（毫秒） |
| abortOnUnmount | boolean  | true   | 卸载时取消请求   |
| beforeFetch    | function | -      | 请求前拦截       |
| afterFetch     | function | -      | 响应后处理       |
| onFetchError   | function | -      | 错误回调         |

| 返回值     | 类型                          | 说明         |
| ---------- | ----------------------------- | ------------ |
| data       | Ref\<T \| null\>              | 响应数据     |
| status     | Ref\<FetchStatus\>            | 请求状态     |
| isLoading  | ComputedRef\<boolean\>        | 是否正在加载 |
| isFinished | ComputedRef\<boolean\>        | 是否完成     |
| error      | Ref\<Error \| null\>          | 错误信息     |
| response   | Ref\<Response \| null\>       | 原始响应     |
| statusCode | Ref\<number \| null\>         | 状态码       |
| execute    | `(throwError?) => Promise   ` | 执行请求     |
| abort      | `() => void                 ` | 取消请求     |
| retry      | `() => Promise              ` | 重试请求     |
| canAbort   | ComputedRef\<boolean\>        | 是否可取消   |
| aborted    | Ref\<boolean\>                | 是否已取消   |

### createFetch

| 选项         | 类型   | 说明         |
| ------------ | ------ | ------------ |
| baseUrl      | string | 基础 URL     |
| options      | object | 默认请求配置 |
| interceptors | object | 拦截器配置   |

| 方法   | 说明        |
| ------ | ----------- |
| get    | GET 请求    |
| post   | POST 请求   |
| put    | PUT 请求    |
| patch  | PATCH 请求  |
| delete | DELETE 请求 |

## 代码位置

```
web/src/
└── composables/
    └── useFetch.ts    # Fetch Composable
```
