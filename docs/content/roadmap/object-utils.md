# 对象工具函数

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:31+4`
- [已实现功能](#已实现功能) `:35+45`
  - [深度操作](#深度操作) `:37+6`
  - [路径操作](#路径操作) `:43+7`
  - [筛选与转换](#筛选与转换) `:50+8`
  - [类型检查](#类型检查) `:58+6`
  - [遍历](#遍历) `:64+5`
  - [差异比较](#差异比较) `:69+4`
  - [其他](#其他) `:73+7`
- [使用方式](#使用方式) `:80+75`
  - [深度操作](#深度操作-1) `:82+16`
  - [路径操作](#路径操作-1) `:98+22`
  - [筛选与转换](#筛选与转换-1) `:120+22`
  - [差异比较](#差异比较-1) `:142+13`
- [API](#api) `:155+15`
  - [主要函数](#主要函数) `:157+13`
- [代码位置](#代码位置) `:170+7`

<!--TOC-->

## 需求背景

项目中大量使用对象操作，需要统一的对象处理工具函数。

## 已实现功能

### 深度操作

- `deepClone` - 深拷贝
- `deepMerge` - 深度合并
- `deepEqual` - 深度比较

### 路径操作

- `get` - 获取嵌套属性
- `set` - 设置嵌套属性
- `has` - 检查属性存在
- `unset` - 删除嵌套属性

### 筛选与转换

- `pick` - 选取指定键
- `omit` - 排除指定键
- `filterObject` - 按条件筛选
- `mapObject` - 转换键值
- `invert` - 键值反转

### 类型检查

- `isObject` - 是否为对象
- `isEmpty` - 是否为空对象
- `isPlainObject` - 是否为普通对象

### 遍历

- `forEachObject` - 遍历对象
- `deepForEach` - 深度遍历

### 差异比较

- `diff` - 比较两个对象差异

### 其他

- `compact` - 移除 null/undefined
- `defaults` - 带默认值创建
- `entries` - 对象转数组
- `fromEntries` - 数组转对象

## 使用方式

### 深度操作

```typescript
import { deepClone, deepMerge, deepEqual } from "@/utils/object";

// 深拷贝
const copy = deepClone({ a: { b: 1 } });

// 深度合并
const merged = deepMerge({ a: { b: 1 } }, { a: { c: 2 } });
// { a: { b: 1, c: 2 } }

// 深度比较
deepEqual({ a: 1 }, { a: 1 }); // true
```

### 路径操作

```typescript
import { get, set, has, unset } from "@/utils/object";

const obj = { a: { b: { c: 1 } } };

// 获取嵌套值
get(obj, "a.b.c"); // 1
get(obj, "a.b.d", "default"); // "default"

// 设置嵌套值
const newObj = set({}, "a.b.c", 1);
// { a: { b: { c: 1 } } }

// 检查属性
has(obj, "a.b.c"); // true

// 删除属性
unset(obj, "a.b.c"); // { a: { b: {} } }
```

### 筛选与转换

```typescript
import { pick, omit, filterObject, mapObject } from "@/utils/object";

const user = { id: 1, name: "John", password: "secret" };

// 选取
pick(user, ["id", "name"]); // { id: 1, name: "John" }

// 排除
omit(user, ["password"]); // { id: 1, name: "John" }

// 过滤
filterObject({ a: 1, b: null, c: 2 }, (v) => v !== null);
// { a: 1, c: 2 }

// 转换
mapObject({ a: 1, b: 2 }, (v) => v * 2);
// { a: 2, b: 4 }
```

### 差异比较

```typescript
import { diff } from "@/utils/object";

const changes = diff({ a: 1, b: 2 }, { a: 1, c: 3 });
// {
//   added: { c: 3 },
//   removed: { b: 2 },
//   changed: {}
// }
```

## API

### 主要函数

| 函数                          | 说明       |
| ----------------------------- | ---------- |
| deepClone(obj)                | 深拷贝     |
| deepMerge(target, ...sources) | 深度合并   |
| deepEqual(a, b)               | 深度比较   |
| get(obj, path, default?)      | 获取嵌套值 |
| set(obj, path, value)         | 设置嵌套值 |
| pick(obj, keys)               | 选取指定键 |
| omit(obj, keys)               | 排除指定键 |
| diff(prev, next)              | 比较差异   |

## 代码位置

```
web/src/
└── utils/
    └── object.ts    # 对象工具函数
```
