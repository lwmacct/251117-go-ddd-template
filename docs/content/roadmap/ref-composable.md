# Ref Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:40+4`
- [已实现功能](#已实现功能) `:44+31`
  - [值处理](#值处理) `:46+8`
  - [历史与状态](#历史与状态) `:54+6`
  - [响应控制](#响应控制) `:60+6`
  - [数据结构](#数据结构) `:66+9`
- [使用方式](#使用方式) `:75+259`
  - [防抖 ref](#防抖-ref) `:77+13`
  - [节流 ref](#节流-ref) `:90+19`
  - [历史记录（撤销/重做）](#历史记录撤销重做) `:109+40`
  - [自动重置](#自动重置) `:149+17`
  - [可控 ref](#可控-ref) `:166+24`
  - [同步 ref](#同步-ref) `:190+20`
  - [模板 ref](#模板-ref) `:210+15`
  - [计数器](#计数器) `:225+21`
  - [布尔开关](#布尔开关) `:246+17`
  - [对象操作](#对象操作) `:263+23`
  - [数组操作](#数组操作) `:286+15`
  - [集合操作](#集合操作) `:301+16`
  - [Map 操作](#map-操作) `:317+17`
- [API](#api) `:334+54`
  - [refDebounced](#refdebounced) `:336+7`
  - [refThrottled](#refthrottled) `:343+8`
  - [refHistory](#refhistory) `:351+22`
  - [useCounter](#usecounter) `:373+15`
- [代码位置](#代码位置) `:388+7`

<!--TOC-->

## 需求背景

前端需要增强的 ref 工具函数，支持防抖、节流、历史记录、自动重置等高级功能。

## 已实现功能

### 值处理

- `refDefault` - 带默认值的 ref
- `refDebounced` - 防抖 ref
- `refThrottled` - 节流 ref
- `refAutoReset` - 自动重置 ref
- `refLocked` - 可锁定 ref

### 历史与状态

- `refHistory` - 带历史记录（撤销/重做）
- `usePrevious` - 获取上一个值
- `useLatest` - 获取最新值

### 响应控制

- `refWithControl` - 可控 ref（拦截 get/set）
- `syncRefs` - 同步两个 ref
- `templateRef` - 模板 ref 辅助

### 数据结构

- `useCounter` - 计数器
- `useBoolean` - 布尔开关
- `useObject` - 对象操作
- `useArray` - 数组操作
- `useSet` - 集合操作
- `useMap` - 映射操作

## 使用方式

### 防抖 ref

```typescript
import { refDebounced } from '@/composables/useRef'

const text = ref('')
const debouncedText = refDebounced(text, { delay: 300 })

// 在输入框中使用
<input v-model="text" />
// debouncedText 在停止输入 300ms 后才更新
```

### 节流 ref

```typescript
import { refThrottled } from "@/composables/useRef";

const scrollY = ref(0);
const throttledScrollY = refThrottled(scrollY, {
  delay: 100,
  leading: true,
  trailing: true,
});

// 滚动监听中使用
window.addEventListener("scroll", () => {
  scrollY.value = window.scrollY;
});
// throttledScrollY 每 100ms 最多更新一次
```

### 历史记录（撤销/重做）

```typescript
import { refHistory } from "@/composables/useRef";

const {
  value: content,
  undo,
  redo,
  canUndo,
  canRedo,
  history,
  clear,
  pause,
  resume,
} = refHistory("", {
  capacity: 50,
  deep: false,
  debounce: 500, // 防抖，减少历史记录数量
});

// 编辑器中使用
content.value = "第一行";
content.value = "第二行";
console.log(history.value); // ['', '第一行']

// 撤销
if (canUndo.value) undo();
console.log(content.value); // '第一行'

// 重做
if (canRedo.value) redo();
console.log(content.value); // '第二行'

// 编辑时暂停记录
pause();
// ... 编辑操作
resume();
```

### 自动重置

```typescript
import { refAutoReset } from "@/composables/useRef";

// 通知消息，3秒后自动清空
const notification = refAutoReset("", { delay: 3000 });

notification.value = "保存成功！";
// 3秒后 notification.value 自动变为 ''

// 错误提示
const error = refAutoReset<string | null>(null, { delay: 5000 });
error.value = "网络错误";
// 5秒后自动清除
```

### 可控 ref

```typescript
import { refWithControl } from "@/composables/useRef";

const { value, pause, resume, silentSet, peek } = refWithControl(0, {
  onGet: (val) => val * 2, // 读取时翻倍
  onSet: (newVal) => Math.max(0, newVal), // 写入时确保非负
  onBeforeSet: (newVal, oldVal) => newVal !== oldVal, // 值变化时才更新
});

console.log(value.value); // 0 (实际值 0 * 2)
value.value = -5; // 被 onSet 处理为 0
console.log(peek()); // 0 (原始值)

// 暂停响应
pause();
value.value = 10; // 无效
resume();

// 静默设置（不触发响应）
silentSet(100);
```

### 同步 ref

```typescript
import { syncRefs } from "@/composables/useRef";

const source = ref(0);
const target = ref(0);

const stop = syncRefs(source, target, {
  immediate: true,
  direction: "both", // 或 'ltr', 'rtl'
});

source.value = 1; // target.value 也变为 1
target.value = 2; // source.value 也变为 2

// 停止同步
stop();
```

### 模板 ref

```typescript
import { templateRef } from '@/composables/useRef'

const { ref: inputRef, isMounted, onMounted } = templateRef<HTMLInputElement>()

onMounted((el) => {
  el.focus() // 元素挂载后自动聚焦
})

// 模板中
<input ref="inputRef" />
```

### 计数器

```typescript
import { useCounter } from "@/composables/useRef";

const { count, inc, dec, reset, set } = useCounter(0, {
  min: 0,
  max: 100,
});

inc(); // count = 1
inc(5); // count = 6
dec(2); // count = 4
set(50); // count = 50
reset(); // count = 0

// 超出范围自动限制
inc(200); // count = 100 (max)
dec(200); // count = 0 (min)
```

### 布尔开关

```typescript
import { useBoolean } from '@/composables/useRef'

const { value: isOpen, toggle, setTrue, setFalse } = useBoolean(false)

toggle() // true
toggle() // false
setTrue() // true
setFalse() // false

// 模态框控制
<button @click="setTrue">打开</button>
<Modal v-if="isOpen" @close="setFalse" />
```

### 对象操作

```typescript
import { useObject } from "@/composables/useRef";

const {
  state: user,
  set,
  merge,
  reset,
  patch,
} = useObject({
  name: "",
  age: 0,
  email: "",
});

set({ name: "John", age: 30, email: "john@example.com" });
merge({ name: "Jane" }); // 只更新 name
patch("age", 31); // 只更新 age
reset(); // 重置为初始值
```

### 数组操作

```typescript
import { useArray } from "@/composables/useRef";

const { array: items, push, pop, remove, clear, insert, update } = useArray([1, 2, 3]);

push(4, 5); // [1, 2, 3, 4, 5]
pop(); // [1, 2, 3, 4]
remove(1); // [1, 3, 4] (移除索引 1)
insert(1, 2); // [1, 2, 3, 4]
update(0, 0); // [0, 2, 3, 4]
clear(); // []
```

### 集合操作

```typescript
import { useSet } from "@/composables/useRef";

const { set: selected, add, remove, has, toggle, clear, values } = useSet<number>();

add(1);
add(2);
has(1); // true
toggle(1); // 移除 1
toggle(3); // 添加 3
values(); // [2, 3]
clear();
```

### Map 操作

```typescript
import { useMap } from "@/composables/useRef";

const { map: cache, set, get, remove, has, keys, values, entries } = useMap<string, number>();

set("a", 1);
set("b", 2);
get("a"); // 1
has("a"); // true
keys(); // ['a', 'b']
values(); // [1, 2]
entries(); // [['a', 1], ['b', 2]]
remove("a");
```

## API

### refDebounced

| 选项    | 类型   | 默认值 | 说明             |
| ------- | ------ | ------ | ---------------- |
| delay   | number | 250    | 防抖延迟（毫秒） |
| maxWait | number | -      | 最大等待时间     |

### refThrottled

| 选项     | 类型    | 默认值 | 说明             |
| -------- | ------- | ------ | ---------------- |
| delay    | number  | 100    | 节流间隔（毫秒） |
| leading  | boolean | true   | 开始时触发       |
| trailing | boolean | true   | 结束时触发       |

### refHistory

| 选项     | 类型     | 默认值 | 说明           |
| -------- | -------- | ------ | -------------- |
| capacity | number   | 10     | 历史记录容量   |
| deep     | boolean  | false  | 深度克隆       |
| clone    | Function | -      | 自定义克隆函数 |
| debounce | number   | 0      | 记录防抖延迟   |

| 返回值  | 类型           | 说明       |
| ------- | -------------- | ---------- |
| value   | Ref\<T\>       | 当前值     |
| history | Ref\<T[]\>     | 历史记录   |
| future  | Ref\<T[]\>     | 未来记录   |
| canUndo | Ref\<boolean\> | 是否可撤销 |
| canRedo | Ref\<boolean\> | 是否可重做 |
| undo    | `() => void`   | 撤销       |
| redo    | `() => void`   | 重做       |
| clear   | `() => void`   | 清空历史   |
| pause   | `() => void`   | 暂停记录   |
| resume  | `() => void`   | 恢复记录   |

### useCounter

| 选项 | 类型   | 默认值    | 说明   |
| ---- | ------ | --------- | ------ |
| min  | number | -Infinity | 最小值 |
| max  | number | Infinity  | 最大值 |

| 返回值 | 类型               | 说明     |
| ------ | ------------------ | -------- |
| count  | Ref\<number\>      | 当前计数 |
| inc    | `(delta?) => void` | 增加     |
| dec    | `(delta?) => void` | 减少     |
| reset  | `() => void`       | 重置     |
| set    | `(value) => void`  | 设置值   |

## 代码位置

```
web/src/
└── composables/
    └── useRef.ts    # Ref Composable
```
