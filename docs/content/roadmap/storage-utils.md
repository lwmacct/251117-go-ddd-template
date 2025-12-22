# 存储工具 Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:21+4`
- [已实现功能](#已实现功能) `:25+9`
  - [useStorage Composable](#usestorage-composable) `:27+7`
- [使用方式](#使用方式) `:34+25`
- [API](#api) `:59+12`
  - [useStorage 返回值](#usestorage-返回值) `:61+10`
- [代码位置](#代码位置) `:71+7`

<!--TOC-->

## 需求背景

需要响应式地使用 localStorage/sessionStorage，支持过期时间和类型安全。

## 已实现功能

### useStorage Composable

- 响应式存储值
- 支持过期时间
- 自动 JSON 序列化
- 类型安全

## 使用方式

```typescript
import { useLocalStorage, useSessionStorage } from "@/composables/useStorage";

// 基础用法
const { value, set, remove } = useLocalStorage<User>("user");

// 带过期时间（1小时）
const { value: token } = useLocalStorage<string>("token", {
  expires: 60 * 60 * 1000,
});

// 带默认值
const { value: settings } = useLocalStorage("settings", {
  defaultValue: { theme: "light" },
});

// 检查过期
const { isExpired } = useLocalStorage("cache");
if (isExpired()) {
  // 重新获取数据
}
```

## API

### useStorage 返回值

| 属性/方法 | 类型                                   | 说明         |
| --------- | -------------------------------------- | ------------ |
| value     | Ref<T \| null>                         | 响应式存储值 |
| set       | `(value: T, expires?: number) => void` | 设置值       |
| get       | () => T \| null                        | 获取值       |
| remove    | `() => void`                           | 删除值       |
| isExpired | `() => boolean`                        | 检查是否过期 |

## 代码位置

```
web/src/
└── composables/
    └── useStorage.ts    # 存储工具
```
