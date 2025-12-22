# 类型工具函数

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:34+4`
- [已实现功能](#已实现功能) `:38+81`
  - [基础类型检查](#基础类型检查) `:40+13`
  - [复杂类型检查](#复杂类型检查) `:53+12`
  - [特殊类型检查](#特殊类型检查) `:65+9`
  - [DOM 类型检查](#dom-类型检查) `:74+9`
  - [字符串检查](#字符串检查) `:83+7`
  - [类型转换](#类型转换) `:90+9`
  - [断言函数](#断言函数) `:99+6`
  - [类型守卫组合](#类型守卫组合) `:105+6`
  - [实用工具](#实用工具) `:111+8`
- [使用方式](#使用方式) `:119+94`
  - [类型守卫](#类型守卫) `:121+18`
  - [类型转换](#类型转换-1) `:139+15`
  - [断言函数](#断言函数-1) `:154+29`
  - [类型守卫组合](#类型守卫组合-1) `:183+13`
  - [获取类型信息](#获取类型信息) `:196+17`
- [API](#api) `:213+15`
  - [主要函数](#主要函数) `:215+13`
- [代码位置](#代码位置) `:228+7`

<!--TOC-->

## 需求背景

前端需要运行时类型检查和类型安全的转换，用于表单验证、API 响应处理等场景。

## 已实现功能

### 基础类型检查

- `isString` - 检查字符串
- `isNumber` - 检查数字（排除 NaN）
- `isFiniteNumber` - 检查有限数字
- `isInteger` - 检查整数
- `isBoolean` - 检查布尔值
- `isNull` - 检查 null
- `isUndefined` - 检查 undefined
- `isNullish` - 检查 null 或 undefined
- `isDefined` - 检查已定义
- `hasValue` - 检查有值

### 复杂类型检查

- `isObject` - 检查对象
- `isPlainObject` - 检查普通对象
- `isArray` - 检查数组
- `isArrayOf` - 检查指定类型数组
- `isNonEmptyArray` - 检查非空数组
- `isFunction` - 检查函数
- `isAsyncFunction` - 检查异步函数
- `isSymbol` - 检查 Symbol
- `isBigInt` - 检查 BigInt

### 特殊类型检查

- `isDate` - 检查 Date
- `isRegExp` - 检查正则表达式
- `isError` - 检查 Error
- `isPromise` - 检查 Promise
- `isMap` - 检查 Map
- `isSet` - 检查 Set

### DOM 类型检查

- `isElement` - 检查 DOM 元素
- `isHTMLElement` - 检查 HTML 元素
- `isNode` - 检查 Node
- `isBlob` - 检查 Blob
- `isFile` - 检查 File
- `isFormData` - 检查 FormData

### 字符串检查

- `isEmptyString` - 检查空字符串
- `isBlankString` - 检查空白字符串
- `isNonEmptyString` - 检查非空字符串
- `isNonBlankString` - 检查非空白字符串

### 类型转换

- `toString` - 安全转字符串
- `toNumber` - 安全转数字
- `toInteger` - 安全转整数
- `toBoolean` - 安全转布尔
- `toArray` - 安全转数组
- `toDate` - 安全转日期

### 断言函数

- `assertDefined` - 断言非空
- `assert` - 断言条件
- `assertNever` - 穷尽检查

### 类型守卫组合

- `unionGuard` - 联合类型守卫
- `intersectionGuard` - 交叉类型守卫
- `notGuard` - 否定类型守卫

### 实用工具

- `getType` - 获取类型名称
- `isSameType` - 比较类型
- `isPrimitive` - 检查原始类型
- `isIterable` - 检查可迭代
- `isJsonSerializable` - 检查可序列化

## 使用方式

### 类型守卫

```typescript
import { isString, isNumber, isArray, isArrayOf } from "@/utils/type";

// 基础检查
if (isString(value)) {
  // value 的类型已收窄为 string
  console.log(value.toUpperCase());
}

// 数组类型检查
if (isArrayOf(data, isNumber)) {
  // data 的类型已收窄为 number[]
  const sum = data.reduce((a, b) => a + b, 0);
}
```

### 类型转换

```typescript
import { toNumber, toBoolean, toArray, toDate } from "@/utils/type";

// 安全转换
const num = toNumber("123.5"); // 123.5
const bool = toBoolean("true"); // true
const arr = toArray("single"); // ['single']
const date = toDate("2024-01-01"); // Date

// 带默认值
const safeNum = toNumber("invalid", -1); // -1
```

### 断言函数

```typescript
import { assertDefined, assert, assertNever } from "@/utils/type";

// 非空断言
function processUser(user: User | null) {
  assertDefined(user, "User is required");
  // 此后 user 类型为 User
  console.log(user.name);
}

// 穷尽检查
type Status = "pending" | "success" | "error";

function handleStatus(status: Status) {
  switch (status) {
    case "pending":
      return "Loading...";
    case "success":
      return "Done!";
    case "error":
      return "Failed!";
    default:
      assertNever(status); // 如果漏掉 case，编译报错
  }
}
```

### 类型守卫组合

```typescript
import { unionGuard, isString, isNumber } from "@/utils/type";

// 创建联合类型守卫
const isStringOrNumber = unionGuard(isString, isNumber);

if (isStringOrNumber(value)) {
  // value 是 string 或 number
}
```

### 获取类型信息

```typescript
import { getType, isPrimitive, isJsonSerializable } from "@/utils/type";

getType(null); // 'null'
getType([]); // 'array'
getType({}); // 'object'
getType(new Date()); // 'date'

isPrimitive("hello"); // true
isPrimitive({}); // false

isJsonSerializable({ a: 1 }); // true
isJsonSerializable(() => {}); // false
```

## API

### 主要函数

| 函数                     | 说明                |
| ------------------------ | ------------------- |
| isString(value)          | 检查字符串          |
| isNumber(value)          | 检查数字            |
| isObject(value)          | 检查对象            |
| isArray(value)           | 检查数组            |
| isNullish(value)         | 检查 null/undefined |
| toNumber(value, default) | 安全转数字          |
| assertDefined(value)     | 断言非空            |
| getType(value)           | 获取类型名称        |

## 代码位置

```
web/src/
└── utils/
    └── type.ts    # 类型工具函数
```
