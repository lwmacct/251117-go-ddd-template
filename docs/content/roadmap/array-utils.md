# 数组工具函数

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:32+4`
- [已实现功能](#已实现功能) `:36+51`
  - [分组与分块](#分组与分块) `:38+6`
  - [去重与交集](#去重与交集) `:44+7`
  - [查找与索引](#查找与索引) `:51+6`
  - [排序](#排序) `:57+5`
  - [变换](#变换) `:62+6`
  - [移动与交换](#移动与交换) `:68+5`
  - [聚合](#聚合) `:73+7`
  - [树形结构](#树形结构) `:80+7`
- [使用方式](#使用方式) `:87+80`
  - [分块与分组](#分块与分组) `:89+21`
  - [去重与集合操作](#去重与集合操作) `:110+21`
  - [排序](#排序-1) `:131+16`
  - [树形操作](#树形操作) `:147+20`
- [API](#api) `:167+15`
  - [主要函数](#主要函数) `:169+13`
- [代码位置](#代码位置) `:182+7`

<!--TOC-->

## 需求背景

项目中大量使用数组操作，需要统一的数组处理工具函数。

## 已实现功能

### 分组与分块

- `chunk` - 数组分块
- `groupBy` - 按字段分组
- `partition` - 按条件分为两组

### 去重与交集

- `unique` - 数组去重
- `intersection` - 数组交集
- `difference` - 数组差集
- `union` - 数组并集

### 查找与索引

- `findIndex` - 查找索引
- `findLastIndex` - 查找最后索引
- `findAllIndices` - 查找所有索引

### 排序

- `sortBy` - 按字段排序
- `sortByMultiple` - 多字段排序

### 变换

- `flatten` - 打平数组
- `shuffle` - 数组洗牌
- `sample` - 随机取样

### 移动与交换

- `move` - 移动元素
- `swap` - 交换元素

### 聚合

- `sum` - 求和
- `average` - 平均值
- `maxBy` - 最大值项
- `minBy` - 最小值项

### 树形结构

- `arrayToTree` - 数组转树
- `treeToArray` - 树转数组
- `findInTree` - 在树中查找
- `filterTree` - 过滤树

## 使用方式

### 分块与分组

```typescript
import { chunk, groupBy, partition } from "@/utils/array";

// 分块
chunk([1, 2, 3, 4, 5], 2); // [[1, 2], [3, 4], [5]]

// 分组
const users = [
  { name: "Alice", role: "admin" },
  { name: "Bob", role: "user" },
];
groupBy(users, "role");
// { admin: [{ name: 'Alice', role: 'admin' }], user: [...] }

// 分区
partition([1, 2, 3, 4, 5], (n) => n % 2 === 0);
// [[2, 4], [1, 3, 5]]
```

### 去重与集合操作

```typescript
import { unique, intersection, difference, union } from "@/utils/array";

// 简单去重
unique([1, 2, 2, 3]); // [1, 2, 3]

// 按字段去重
unique([{ id: 1 }, { id: 1 }], "id"); // [{ id: 1 }]

// 交集
intersection([1, 2, 3], [2, 3, 4]); // [2, 3]

// 差集
difference([1, 2, 3], [2, 3, 4]); // [1]

// 并集
union([1, 2], [2, 3], [3, 4]); // [1, 2, 3, 4]
```

### 排序

```typescript
import { sortBy, sortByMultiple } from "@/utils/array";

// 单字段排序
sortBy(users, "name"); // 按 name 升序
sortBy(users, "age", "desc"); // 按 age 降序

// 多字段排序
sortByMultiple(users, [
  { key: "role", order: "asc" },
  { key: "name", order: "asc" },
]);
```

### 树形操作

```typescript
import { arrayToTree, treeToArray, findInTree } from "@/utils/array";

// 数组转树
const items = [
  { id: 1, parentId: null, name: "Root" },
  { id: 2, parentId: 1, name: "Child 1" },
  { id: 3, parentId: 1, name: "Child 2" },
];
const tree = arrayToTree(items);

// 在树中查找
const node = findInTree(tree, (n) => n.id === 3);

// 树转数组
const flat = treeToArray(tree);
```

## API

### 主要函数

| 函数                        | 说明 |
| --------------------------- | ---- |
| chunk(arr, size)            | 分块 |
| groupBy(arr, key)           | 分组 |
| unique(arr, key?)           | 去重 |
| sortBy(arr, key, order?)    | 排序 |
| flatten(arr, depth?)        | 打平 |
| shuffle(arr)                | 洗牌 |
| sum(arr, getter?)           | 求和 |
| arrayToTree(items, options) | 转树 |

## 代码位置

```
web/src/
└── utils/
    └── array.ts    # 数组工具函数
```
