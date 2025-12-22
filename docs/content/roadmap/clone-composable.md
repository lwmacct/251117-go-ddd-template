# Clone Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:33+4`
- [已实现功能](#已实现功能) `:37+22`
  - [克隆工具](#克隆工具) `:39+5`
  - [响应式克隆](#响应式克隆) `:44+6`
  - [状态管理](#状态管理) `:50+5`
  - [缓存](#缓存) `:55+4`
- [使用方式](#使用方式) `:59+179`
  - [深克隆](#深克隆) `:61+20`
  - [响应式克隆](#响应式克隆-1) `:81+23`
  - [表单编辑场景](#表单编辑场景) `:104+32`
  - [脏状态检测](#脏状态检测) `:136+30`
  - [状态快照](#状态快照) `:166+29`
  - [延迟同步](#延迟同步) `:195+21`
  - [函数记忆化](#函数记忆化) `:216+22`
- [API](#api) `:238+42`
  - [useCloned](#usecloned) `:240+16`
  - [useDirtyState](#usedirtystate) `:256+11`
  - [useSnapshot](#usesnapshot) `:267+13`
- [代码位置](#代码位置) `:280+7`

<!--TOC-->

## 需求背景

前端需要响应式对象的克隆、状态同步、脏状态检测等功能，用于表单编辑、撤销重做等场景。

## 已实现功能

### 克隆工具

- `deepClone` - 深克隆对象
- `structuredClonePolyfill` - structuredClone 兼容实现

### 响应式克隆

- `useCloned` - 克隆响应式值
- `useManualClone` - 手动控制的克隆
- `useSyncedRef` - 延迟同步的引用

### 状态管理

- `useDirtyState` - 脏状态检测
- `useSnapshot` - 状态快照

### 缓存

- `useMemoize` - 函数记忆化

## 使用方式

### 深克隆

```typescript
import { deepClone, structuredClonePolyfill } from "@/composables/useClone";

const original = {
  name: "John",
  nested: { value: 1 },
  arr: [1, 2, 3],
  date: new Date(),
  map: new Map([["key", "value"]]),
};

const cloned = deepClone(original);
// cloned 是完全独立的副本

// 或使用 structuredClone 兼容版本
const cloned2 = structuredClonePolyfill(original);
```

### 响应式克隆

```typescript
import { useCloned } from "@/composables/useClone";

const source = ref({ name: "John", age: 25 });

const { cloned, sync, reset, isModified } = useCloned(source);

// 修改克隆值不影响源值
cloned.value.name = "Jane";
console.log(source.value.name); // 'John'
console.log(isModified.value); // true

// 同步克隆值到源值
sync();
console.log(source.value.name); // 'Jane'

// 重置为源值
reset();
console.log(cloned.value.name); // 'Jane'
```

### 表单编辑场景

```typescript
import { useManualClone } from "@/composables/useClone";

const user = ref({
  name: "John",
  email: "john@example.com",
  role: "user",
});

const { cloned, apply, reset, isModified } = useManualClone(user);

// 用户在表单中编辑 cloned
cloned.value.name = "Jane";
cloned.value.email = "jane@example.com";

// 保存时应用更改
const handleSave = async () => {
  apply(); // 将 cloned 应用到 user
  await saveToServer(user.value);
};

// 取消时重置
const handleCancel = () => {
  reset(); // 重置 cloned 为 user 的值
};

// 离开提示
const canLeave = computed(() => !isModified.value);
```

### 脏状态检测

```typescript
import { useDirtyState } from "@/composables/useClone";

const { state, isDirty, dirtyFields, markClean, reset, getChanges } = useDirtyState({
  name: "John",
  email: "john@example.com",
  age: 25,
});

// 修改状态
state.value.name = "Jane";
state.value.age = 26;

console.log(isDirty.value); // true
console.log(dirtyFields.value); // ['name', 'age']

// 获取变更
const changes = getChanges();
// { name: 'Jane', age: 26 }

// 保存后标记为干净
await saveToServer(state.value);
markClean();

// 或重置为初始状态
reset();
```

### 状态快照

```typescript
import { useSnapshot } from "@/composables/useClone";

const { state, takeSnapshot, restoreSnapshot, snapshotCount } = useSnapshot({
  items: [],
});

// 修改前创建快照
takeSnapshot();
state.value.items.push({ id: 1, name: "Item 1" });

// 再次修改前创建快照
takeSnapshot();
state.value.items.push({ id: 2, name: "Item 2" });

console.log(snapshotCount.value); // 3

// 恢复到特定快照
restoreSnapshot(1); // 恢复到只有一个 item 的状态

// 恢复到上一个快照
restorePrevious();

// 清除所有快照
clearSnapshots();
```

### 延迟同步

```typescript
import { useSyncedRef } from "@/composables/useClone";

const searchQuery = ref("");

const { local, isSyncing, forceSync } = useSyncedRef(searchQuery, {
  delay: 500, // 500ms 延迟
});

// 用户输入时修改 local
const onInput = (value: string) => {
  local.value = value;
  // 500ms 后自动同步到 searchQuery
};

// 立即同步
forceSync();
```

### 函数记忆化

```typescript
import { useMemoize } from "@/composables/useClone";

const expensiveCalculation = (a: number, b: number) => {
  console.log("计算中...");
  return a * b;
};

const { memoized, clear, size, has } = useMemoize(expensiveCalculation);

memoized(2, 3); // 输出 "计算中..."，返回 6
memoized(2, 3); // 从缓存返回 6，不输出
memoized(3, 4); // 输出 "计算中..."，返回 12

console.log(size.value); // 2
console.log(has(2, 3)); // true

clear(); // 清除所有缓存
```

## API

### useCloned

| 选项      | 类型                | 默认值    | 说明           |
| --------- | ------------------- | --------- | -------------- |
| deep      | boolean             | true      | 是否深克隆     |
| immediate | boolean             | true      | 是否立即克隆   |
| clone     | `(source: T) => T ` | deepClone | 自定义克隆函数 |
| manual    | boolean             | false     | 是否手动同步   |

| 返回值     | 类型                     | 说明       |
| ---------- | ------------------------ | ---------- |
| cloned     | Ref\<T\>                 | 克隆的值   |
| sync       | `() => void            ` | 同步到源值 |
| reset      | `() => void            ` | 重置为源值 |
| isModified | ComputedRef\<boolean\>   | 是否已修改 |

### useDirtyState

| 返回值      | 类型                     | 说明           |
| ----------- | ------------------------ | -------------- |
| state       | Ref\<T\>                 | 当前值         |
| isDirty     | ComputedRef\<boolean\>   | 是否为脏状态   |
| dirtyFields | ComputedRef\<string[]\>  | 脏字段列表     |
| markClean   | `() => void            ` | 标记为干净     |
| reset       | `() => void            ` | 重置为初始状态 |
| getChanges  | `() => Partial\<T\>    ` | 获取更改       |

### useSnapshot

| 返回值          | 类型                      | 说明         |
| --------------- | ------------------------- | ------------ |
| state           | Ref\<T\>                  | 当前值       |
| snapshotIndex   | Ref\<number\>             | 当前快照索引 |
| snapshotCount   | ComputedRef\<number\>     | 快照数量     |
| takeSnapshot    | `() => void            `  | 创建快照     |
| restoreSnapshot | `(index: number) => void` | 恢复到快照   |
| restorePrevious | `() => void            `  | 恢复到上一个 |
| clearSnapshots  | `() => void            `  | 清除所有快照 |
| getSnapshots    | `() => T[]             `  | 获取所有快照 |

## 代码位置

```
web/src/
└── composables/
    └── useClone.ts    # Clone Composable
```
