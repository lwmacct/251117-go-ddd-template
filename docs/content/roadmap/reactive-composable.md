# Reactive Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:37+4`
- [已实现功能](#已实现功能) `:41+33`
  - [核心增强](#核心增强) `:43+7`
  - [对象操作](#对象操作) `:50+8`
  - [同步与控制](#同步与控制) `:58+6`
  - [验证与监听](#验证与监听) `:64+6`
  - [工具](#工具) `:70+4`
- [使用方式](#使用方式) `:74+236`
  - [带历史记录](#带历史记录) `:76+36`
  - [响应式表单](#响应式表单) `:112+30`
  - [带验证](#带验证) `:142+39`
  - [可重置](#可重置) `:181+22`
  - [对象选取与排除](#对象选取与排除) `:203+21`
  - [条件响应式](#条件响应式) `:224+15`
  - [同步响应式对象](#同步响应式对象) `:239+17`
  - [部分只读](#部分只读) `:256+14`
  - [深度监听](#深度监听) `:270+21`
  - [值转换](#值转换) `:291+19`
- [API](#api) `:310+42`
  - [reactiveHistory](#reactivehistory) `:312+19`
  - [reactiveForm](#reactiveform) `:331+12`
  - [reactiveValidated](#reactivevalidated) `:343+9`
- [代码位置](#代码位置) `:352+7`

<!--TOC-->

## 需求背景

前端需要增强的 reactive 工具函数，支持历史记录、表单管理、验证等高级功能。

## 已实现功能

### 核心增强

- `reactiveWithOptions` - 带选项的响应式对象
- `reactiveHistory` - 带历史记录（撤销/重做）
- `reactiveForm` - 响应式表单
- `reactiveResettable` - 可重置的响应式对象

### 对象操作

- `reactivePick` - 选取属性
- `reactiveOmit` - 排除属性
- `reactiveMerge` - 合并对象
- `reactiveExtend` - 扩展属性
- `reactiveDefault` - 带默认值

### 同步与控制

- `syncReactive` - 同步两个响应式对象
- `reactiveWhen` - 条件响应式
- `reactiveWithReadonly` - 部分只读

### 验证与监听

- `reactiveValidated` - 带验证
- `reactiveTransform` - 值转换
- `watchReactive` - 深度监听

### 工具

- `reactiveUtils` - 响应式工具集

## 使用方式

### 带历史记录

```typescript
import { reactiveHistory } from "@/composables/useReactive";

const { state, undo, redo, canUndo, canRedo, history, clear } = reactiveHistory(
  {
    name: "",
    age: 0,
  },
  {
    capacity: 50, // 最多保存 50 条历史
  },
);

state.name = "John";
state.age = 30;

console.log(history.value); // [{ name: '', age: 0 }, { name: 'John', age: 0 }]

// 撤销
if (canUndo.value) {
  undo();
  console.log(state.age); // 0
}

// 重做
if (canRedo.value) {
  redo();
  console.log(state.age); // 30
}

// 清空历史
clear();
```

### 响应式表单

```typescript
import { reactiveForm } from "@/composables/useReactive";

const { data, isDirty, isSubmitting, reset, getChanges, apply, submit } = reactiveForm({
  username: "",
  email: "",
  password: "",
});

// 编辑表单
data.username = "john";
data.email = "john@example.com";

console.log(isDirty.value); // true

// 获取变更的字段
const changes = getChanges();
// { username: 'john', email: 'john@example.com' }

// 提交表单
await submit(async (formData) => {
  await api.register(formData);
});

// 重置表单
reset();
```

### 带验证

```typescript
import { reactiveValidated } from "@/composables/useReactive";

const { data, errors, isValid, validate } = reactiveValidated(
  {
    email: "",
    password: "",
    age: 0,
  },
  {
    email: (v) => v.includes("@") || "请输入有效的邮箱地址",
    password: (v) => v.length >= 6 || "密码至少 6 个字符",
    age: (v) => v >= 18 || "年龄必须大于等于 18",
  },
);

data.email = "invalid";
data.password = "123";
data.age = 16;

console.log(errors.value);
// {
//   email: '请输入有效的邮箱地址',
//   password: '密码至少 6 个字符',
//   age: '年龄必须大于等于 18'
// }

console.log(isValid.value); // false

// 修正错误
data.email = "john@example.com";
data.password = "123456";
data.age = 20;

console.log(isValid.value); // true
```

### 可重置

```typescript
import { reactiveResettable } from "@/composables/useReactive";

const { state, reset, setInitial } = reactiveResettable({
  count: 0,
  name: "",
});

state.count = 10;
state.name = "John";

reset(); // state = { count: 0, name: '' }

state.count = 20;
setInitial(); // 设置当前值为新的初始值

state.count = 30;
reset(); // state = { count: 20, name: 'John' }
```

### 对象选取与排除

```typescript
import { reactivePick, reactiveOmit } from "@/composables/useReactive";

const user = reactive({
  id: 1,
  name: "John",
  email: "john@example.com",
  password: "secret",
});

// 选取属性
const picked = reactivePick(user, ["id", "name"]);
// { id: 1, name: 'John' }

// 排除属性
const safe = reactiveOmit(user, ["password"]);
// { id: 1, name: 'John', email: 'john@example.com' }
```

### 条件响应式

```typescript
import { reactiveWhen } from "@/composables/useReactive";

const isAdmin = ref(true);

const permissions = reactiveWhen(isAdmin, { read: true, write: true, delete: true, admin: true }, { read: true, write: false, delete: false, admin: false });

console.log(permissions.admin); // true

isAdmin.value = false;
console.log(permissions.admin); // false
```

### 同步响应式对象

```typescript
import { syncReactive } from "@/composables/useReactive";

const source = reactive({ count: 0 });
const target = reactive({ count: 0 });

const stop = syncReactive(source, target);

source.count = 10;
console.log(target.count); // 10

// 停止同步
stop();
```

### 部分只读

```typescript
import { reactiveWithReadonly } from "@/composables/useReactive";

const user = reactiveWithReadonly(
  { id: 1, name: "John", createdAt: "2024-01-01" },
  ["id", "createdAt"], // 这些字段只读
);

user.name = "Jane"; // 成功
user.id = 2; // 静默失败，id 仍为 1
```

### 深度监听

```typescript
import { watchReactive } from "@/composables/useReactive";

const state = reactive({
  user: {
    profile: {
      name: "John",
    },
  },
});

watchReactive(state, (path, newValue, oldValue) => {
  console.log(`${path} 从 ${oldValue} 变为 ${newValue}`);
});

state.user.profile.name = "Jane";
// 输出: 'user.profile.name 从 John 变为 Jane'
```

### 值转换

```typescript
import { reactiveTransform } from "@/composables/useReactive";

const raw = reactive({
  price: 100,
  createdAt: new Date(),
});

const formatted = reactiveTransform(raw, {
  price: (v) => `¥${v.toFixed(2)}`,
  createdAt: (v) => v.toLocaleDateString(),
});

console.log(formatted.price); // '¥100.00'
console.log(formatted.createdAt); // '2024/1/1'
```

## API

### reactiveHistory

| 选项     | 类型     | 默认值               | 说明         |
| -------- | -------- | -------------------- | ------------ |
| capacity | number   | 10                   | 历史记录容量 |
| deep     | boolean  | true                 | 深度监听     |
| clone    | Function | JSON.parse/stringify | 克隆函数     |

| 返回值  | 类型           | 说明         |
| ------- | -------------- | ------------ |
| state   | Reactive       | 响应式状态   |
| history | Ref\<T[]\>     | 历史记录     |
| canUndo | Ref\<boolean\> | 是否可撤销   |
| canRedo | Ref\<boolean\> | 是否可重做   |
| undo    | `() => void`   | 撤销         |
| redo    | `() => void`   | 重做         |
| clear   | `() => void`   | 清空历史     |
| commit  | `() => void`   | 手动提交快照 |

### reactiveForm

| 返回值       | 类型                   | 说明       |
| ------------ | ---------------------- | ---------- |
| data         | Reactive               | 表单数据   |
| isDirty      | Ref\<boolean\>         | 是否已修改 |
| isSubmitting | Ref\<boolean\>         | 是否提交中 |
| reset        | `() => void`           | 重置表单   |
| getChanges   | `() => Partial`        | 获取变更   |
| apply        | `(changes) => void`    | 应用变更   |
| submit       | `(handler) => Promise` | 提交表单   |

### reactiveValidated

| 返回值   | 类型            | 说明     |
| -------- | --------------- | -------- |
| data     | Reactive        | 表单数据 |
| errors   | Ref\<object\>   | 错误信息 |
| isValid  | Ref\<boolean\>  | 是否有效 |
| validate | `() => boolean` | 手动验证 |

## 代码位置

```
web/src/
└── composables/
    └── useReactive.ts    # Reactive Composable
```
