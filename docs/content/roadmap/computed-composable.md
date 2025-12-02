# Computed Composable

<!--TOC-->

- [需求背景](#需求背景) `:36:39`
- [已实现功能](#已实现功能) `:40:41`
  - [核心增强](#核心增强) `:42:48`
  - [对象操作](#对象操作) `:49:55`
  - [性能优化](#性能优化) `:56:61`
  - [数组操作](#数组操作) `:62:70`
  - [其他](#其他) `:71:77`
- [使用方式](#使用方式) `:78:79`
  - [异步 computed](#异步-computed) `:80:108`
  - [防抖 computed](#防抖-computed) `:109:122`
  - [节流 computed](#节流-computed) `:123:137`
  - [可控 computed](#可控-computed) `:138:157`
  - [对象属性选取](#对象属性选取) `:158:178`
  - [从多个源创建](#从多个源创建) `:179:193`
  - [数组操作](#数组操作-1) `:194:223`
  - [带历史记录](#带历史记录) `:224:242`
  - [条件 computed](#条件-computed) `:243:261`
  - [可写 computed](#可写-computed) `:262:280`
- [API](#api) `:281:282`
  - [computedAsync](#computedasync) `:283:299`
  - [computedDebounced](#computeddebounced) `:300:306`
  - [computedThrottled](#computedthrottled) `:307:314`
  - [computedArray](#computedarray) `:315:330`
- [代码位置](#代码位置) `:331:337`

<!--TOC-->

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

## 需求背景

前端需要增强的 computed 工具函数，支持异步计算、防抖、节流、数组操作等高级功能。

## 已实现功能

### 核心增强

- `computedEager` - 立即求值（无缓存）
- `computedAsync` - 异步 computed
- `computedWithControl` - 可控 computed（暂停/恢复）
- `writableComputed` - 可写 computed

### 对象操作

- `computedPick` - 选取对象属性
- `computedOmit` - 排除对象属性
- `computedObject` - 组合为对象
- `computedFrom` - 从多个源创建

### 性能优化

- `computedDebounced` - 防抖 computed
- `computedThrottled` - 节流 computed
- `computedDelayed` - 延迟更新

### 数组操作

- `computedArray` - 数组派生值（sum/avg/min/max）
- `computedMap` - 映射转换
- `computedFilter` - 过滤
- `computedFind` - 查找
- `computedGroupBy` - 分组
- `computedSort` - 排序

### 其他

- `toComputed` - 转换为 computed
- `computedWithHistory` - 带历史记录
- `computedIf` - 条件计算
- `computedDefault` - 带默认值

## 使用方式

### 异步 computed

```typescript
import { computedAsync } from "@/composables/useComputed";

const userId = ref(1);

const {
  state: user,
  isLoading,
  error,
  execute,
} = computedAsync(
  async () => {
    const res = await fetch(`/api/users/${userId.value}`);
    return res.json();
  },
  {
    initialValue: null,
    lazy: false,
    onError: (err) => console.error(err),
    debounce: 300,
  },
);

// 手动重新获取
execute();
```

### 防抖 computed

```typescript
import { computedDebounced } from "@/composables/useComputed";

const searchQuery = ref("");

// 用户停止输入 300ms 后才更新
const debouncedQuery = computedDebounced(() => searchQuery.value.trim(), {
  debounce: 300,
  maxWait: 1000, // 最大等待 1 秒
});
```

### 节流 computed

```typescript
import { computedThrottled } from "@/composables/useComputed";

const scrollY = ref(0);

// 每 100ms 最多更新一次
const throttledScrollY = computedThrottled(() => scrollY.value, {
  throttle: 100,
  leading: true,
  trailing: true,
});
```

### 可控 computed

```typescript
import { computedWithControl } from "@/composables/useComputed";

const count = ref(0);

const { state, pause, resume, trigger, isPaused } = computedWithControl(() => count.value * 2);

// 编辑模式时暂停自动更新
pause();
count.value = 5; // state 不会更新

// 手动触发更新
trigger();

// 完成编辑后恢复
resume();
```

### 对象属性选取

```typescript
import { computedPick, computedOmit } from "@/composables/useComputed";

const user = reactive({
  id: 1,
  name: "John",
  email: "john@example.com",
  password: "secret",
});

// 只选取 id 和 name
const basicInfo = computedPick(user, ["id", "name"]);
// { id: 1, name: 'John' }

// 排除密码字段
const safeUser = computedOmit(user, ["password"]);
// { id: 1, name: 'John', email: 'john@example.com' }
```

### 从多个源创建

```typescript
import { computedFrom } from "@/composables/useComputed";

const firstName = ref("John");
const lastName = ref("Doe");
const age = ref(30);

const userInfo = computedFrom([firstName, lastName, age], ([first, last, userAge]) => ({
  fullName: `${first} ${last}`,
  isAdult: userAge >= 18,
}));
```

### 数组操作

```typescript
import { computedArray, computedMap, computedFilter, computedGroupBy } from "@/composables/useComputed";

const numbers = ref([1, 2, 3, 4, 5]);

// 获取数组派生值
const { sum, avg, min, max, count, isEmpty, unique, sorted, reversed } = computedArray(numbers);
console.log(sum.value); // 15
console.log(avg.value); // 3

// 映射转换
const doubled = computedMap(numbers, (n) => n * 2);
// [2, 4, 6, 8, 10]

// 过滤
const evens = computedFilter(numbers, (n) => n % 2 === 0);
// [2, 4]

// 分组
const items = ref([
  { type: "fruit", name: "apple" },
  { type: "vegetable", name: "carrot" },
  { type: "fruit", name: "banana" },
]);
const grouped = computedGroupBy(items, (item) => item.type);
// { fruit: [...], vegetable: [...] }
```

### 带历史记录

```typescript
import { computedWithHistory } from "@/composables/useComputed";

const count = ref(0);

const { current, history, canUndo, undo, clear } = computedWithHistory(() => count.value * 2, { capacity: 10 });

count.value = 1; // current = 2, history = [0]
count.value = 2; // current = 4, history = [0, 2]

if (canUndo.value) {
  const prev = undo(); // 返回 2
}

clear(); // 清空历史
```

### 条件 computed

```typescript
import { computedIf, computedDefault } from "@/composables/useComputed";

const isAdmin = ref(false);

// 根据条件返回不同值
const permissions = computedIf(
  isAdmin,
  () => ["read", "write", "delete"],
  () => ["read"],
);

// 带默认值
const data = ref<string | null>(null);
const safeData = computedDefault(() => data.value, "默认值");
```

### 可写 computed

```typescript
import { writableComputed } from "@/composables/useComputed";

const rawValue = ref("hello");

// 自定义 get/set 转换
const upperValue = writableComputed(
  () => rawValue.value.toUpperCase(),
  (value) => {
    rawValue.value = value.toLowerCase();
  },
);

console.log(upperValue.value); // 'HELLO'
upperValue.value = "WORLD"; // rawValue 变为 'world'
```

## API

### computedAsync

| 选项         | 类型     | 默认值    | 说明                |
| ------------ | -------- | --------- | ------------------- |
| initialValue | T        | undefined | 初始值              |
| lazy         | boolean  | false     | 是否懒加载          |
| onError      | Function | -         | 错误处理函数        |
| shallow      | boolean  | false     | 是否使用 shallowRef |
| debounce     | number   | 0         | 防抖延迟            |

| 返回值    | 类型                 | 说明       |
| --------- | -------------------- | ---------- |
| state     | Ref\<T\>             | 计算结果   |
| isLoading | Ref\<boolean\>       | 是否加载中 |
| error     | Ref\<Error \| null\> | 错误信息   |
| execute   | `() => Promise`      | 手动执行   |

### computedDebounced

| 选项     | 类型   | 默认值 | 说明         |
| -------- | ------ | ------ | ------------ |
| debounce | number | 250    | 防抖延迟     |
| maxWait  | number | -      | 最大等待时间 |

### computedThrottled

| 选项     | 类型    | 默认值 | 说明       |
| -------- | ------- | ------ | ---------- |
| throttle | number  | 100    | 节流间隔   |
| leading  | boolean | true   | 开始时执行 |
| trailing | boolean | true   | 结束时执行 |

### computedArray

| 返回值   | 类型                   | 说明         |
| -------- | ---------------------- | ------------ |
| sum      | ComputedRef\<number\>  | 数组求和     |
| avg      | ComputedRef\<number\>  | 平均值       |
| min      | ComputedRef\<T\>       | 最小值       |
| max      | ComputedRef\<T\>       | 最大值       |
| count    | ComputedRef\<number\>  | 元素数量     |
| first    | ComputedRef\<T\>       | 第一个元素   |
| last     | ComputedRef\<T\>       | 最后一个元素 |
| isEmpty  | ComputedRef\<boolean\> | 是否为空     |
| unique   | ComputedRef\<T[]\>     | 去重数组     |
| sorted   | ComputedRef\<T[]\>     | 排序后数组   |
| reversed | ComputedRef\<T[]\>     | 反转后数组   |

## 代码位置

```
web/src/
└── composables/
    └── useComputed.ts    # Computed Composable
```
