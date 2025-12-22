# 防抖搜索

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:22+4`
- [已实现功能](#已实现功能) `:26+34`
  - [useDebounce Composable](#usedebounce-composable) `:28+8`
  - [使用方式](#使用方式) `:36+16`
  - [页面集成](#页面集成) `:52+8`
- [配置选项](#配置选项) `:60+7`
- [代码位置](#代码位置) `:67+8`
- [已集成页面](#已集成页面) `:75+4`

<!--TOC-->

## 需求背景

管理后台的搜索功能需要在用户输入时自动触发搜索，但频繁的 API 调用会影响性能。需要实现防抖机制，在用户停止输入后再发起请求。

## 已实现功能

### useDebounce Composable

提供三个主要功能：

1. **useDebouncedRef** - 创建防抖响应式值
2. **useDebounceFn** - 创建防抖函数
3. **useSearchDebounce** - 搜索场景专用

### 使用方式

```typescript
import { useDebouncedRef } from "@/composables/useDebounce";

// 在 composable 中
const searchQuery = ref("");
const debouncedSearchQuery = useDebouncedRef(searchQuery, { delay: 300 });

// 监听防抖值变化，自动触发搜索
watch(debouncedSearchQuery, () => {
  pagination.page = 1;
  fetchData();
});
```

### 页面集成

```vue
<template>
  <v-text-field v-model="searchQuery" prepend-inner-icon="mdi-magnify" label="搜索" clearable placeholder="输入后自动搜索..." />
</template>
```

## 配置选项

| 选项      | 类型    | 默认值 | 说明               |
| --------- | ------- | ------ | ------------------ |
| delay     | number  | 300    | 延迟时间（毫秒）   |
| immediate | boolean | false  | 是否立即执行第一次 |

## 代码位置

```
web/src/
└── composables/
    └── useDebounce.ts    # 防抖工具 Composable
```

## 已集成页面

- ✅ 用户管理 (`/admin/users`)
- ✅ 角色管理 (`/admin/roles`)
