# VModel Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:35+4`
- [已实现功能](#已实现功能) `:39+17`
  - [v-model 处理](#v-model-处理) `:41+8`
  - [高级功能](#高级功能) `:49+7`
- [使用方式](#使用方式) `:56+139`
  - [基础 v-model](#基础-v-model) `:58+21`
  - [多个 v-model](#多个-v-model) `:79+16`
  - [代理模式（表单编辑）](#代理模式表单编辑) `:95+25`
  - [受控/非受控组件](#受控非受控组件) `:120+19`
  - [防抖 v-model](#防抖-v-model) `:139+13`
  - [节流 v-model](#节流-v-model) `:152+13`
  - [切换值](#切换值) `:165+14`
  - [循环列表](#循环列表) `:179+16`
- [API](#api) `:195+60`
  - [useVModel](#usevmodel) `:197+10`
  - [useProxyModel](#useproxymodel) `:207+9`
  - [useControlled](#usecontrolled) `:216+13`
  - [useDebouncedVModel / useThrottledVModel](#usedebouncedvmodel-usethrottledvmodel) `:229+7`
  - [useToggle](#usetoggle) `:236+9`
  - [useCycleList](#usecyclelist) `:245+10`
- [代码位置](#代码位置) `:255+7`

<!--TOC-->

## 需求背景

前端需要简化 v-model 双向绑定的处理，支持防抖、节流、代理模式等高级功能。

## 已实现功能

### v-model 处理

- `useVModel` - 基础 v-model 双向绑定
- `useVModels` - 多个 v-model 处理
- `useModelValue` - 简化的 modelValue 处理
- `useProxyModel` - 代理 model（本地修改）
- `useControlled` - 受控/非受控组件模式

### 高级功能

- `useDebouncedVModel` - 防抖的 v-model
- `useThrottledVModel` - 节流的 v-model
- `useToggle` - 切换值
- `useCycleList` - 循环列表

## 使用方式

### 基础 v-model

```typescript
// 父组件
<MyInput v-model="name" />

// 子组件
<script setup lang="ts">
import { useVModel } from '@/composables/useVModel'

const props = defineProps<{ modelValue: string }>()
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

const model = useVModel(props, 'modelValue', emit)
</script>

<template>
  <input v-model="model" />
</template>
```

### 多个 v-model

```typescript
// 父组件
<UserForm v-model:firstName="first" v-model:lastName="last" />

// 子组件
const props = defineProps<{
  firstName: string
  lastName: string
}>()
const emit = defineEmits(['update:firstName', 'update:lastName'])

const { firstName, lastName } = useVModels(props, emit)
```

### 代理模式（表单编辑）

```typescript
const props = defineProps<{ modelValue: User }>();
const emit = defineEmits<{ "update:modelValue": [value: User] }>();

const { proxy, isModified, sync, reset } = useProxyModel(props, "modelValue", emit);

// 本地修改不会立即触发更新
proxy.value.name = "Jane";
proxy.value.email = "jane@example.com";

console.log(isModified.value); // true

// 保存时同步
const handleSave = () => {
  sync(); // 触发 update:modelValue
};

// 取消时重置
const handleCancel = () => {
  reset(); // 恢复原值
};
```

### 受控/非受控组件

```typescript
// 支持两种使用方式
// 受控: <MyInput v-model="value" />
// 非受控: <MyInput :defaultValue="initialValue" />

const props = defineProps<{
  modelValue?: string;
  defaultValue?: string;
}>();
const emit = defineEmits<{ "update:modelValue": [value: string] }>();

const { value, isControlled, setValue } = useControlled(props, "modelValue", emit, { defaultValue: props.defaultValue || "" });

// 无论哪种模式，都使用相同的 API
setValue("new value");
```

### 防抖 v-model

```typescript
const props = defineProps<{ modelValue: string }>();
const emit = defineEmits<{ "update:modelValue": [value: string] }>();

const model = useDebouncedVModel(props, "modelValue", emit, {
  debounce: 300, // 300ms 防抖
});

// 输入时会防抖触发更新
```

### 节流 v-model

```typescript
const props = defineProps<{ modelValue: number }>();
const emit = defineEmits<{ "update:modelValue": [value: number] }>();

const model = useThrottledVModel(props, "modelValue", emit, {
  throttle: 100, // 100ms 节流
});

// 频繁更新时会节流
```

### 切换值

```typescript
import { useToggle } from "@/composables/useVModel";

const { value, toggle, setTrue, setFalse } = useToggle(false);

toggle(); // true
toggle(); // false
setTrue(); // true
setFalse(); // false
toggle(true); // 强制设置为 true
```

### 循环列表

```typescript
import { useCycleList } from "@/composables/useVModel";

const themes = ["light", "dark", "auto"];
const { value, next, prev, go, index } = useCycleList(themes);

console.log(value.value); // 'light'
next(); // 'dark'
next(); // 'auto'
next(); // 'light' (循环)
prev(); // 'auto'
go(1); // 'dark'
```

## API

### useVModel

| 选项         | 类型                 | 默认值 | 说明           |
| ------------ | -------------------- | ------ | -------------- |
| deep         | boolean              | false  | 是否深度监听   |
| defaultValue | T                    | -      | 默认值         |
| onChange     | `(value: T) => void` | -      | 值变化回调     |
| passive      | boolean              | false  | 被动模式       |
| clone        | `(value: T) => T  `  | -      | 自定义克隆函数 |

### useProxyModel

| 返回值     | 类型             | 说明       |
| ---------- | ---------------- | ---------- |
| proxy      | Ref\<T\>         | 代理值     |
| isModified | Ref\<boolean\>   | 是否已修改 |
| sync       | `() => void    ` | 同步到源值 |
| reset      | `() => void    ` | 重置为源值 |

### useControlled

| 选项         | 类型                 | 说明       |
| ------------ | -------------------- | ---------- |
| defaultValue | T                    | 默认值     |
| onChange     | `(value: T) => void` | 值变化回调 |

| 返回值       | 类型                 | 说明     |
| ------------ | -------------------- | -------- |
| value        | Ref\<T\>             | 当前值   |
| isControlled | boolean              | 是否受控 |
| setValue     | `(value: T) => void` | 设置值   |

### useDebouncedVModel / useThrottledVModel

| 选项     | 类型   | 默认值 | 说明             |
| -------- | ------ | ------ | ---------------- |
| debounce | number | 300    | 防抖延迟（毫秒） |
| throttle | number | 100    | 节流间隔（毫秒） |

### useToggle

| 返回值   | 类型                        | 说明     |
| -------- | --------------------------- | -------- |
| value    | Ref\<boolean\>              | 当前值   |
| toggle   | `(value?: boolean) => void` | 切换     |
| setTrue  | `() => void               ` | 设置为真 |
| setFalse | `() => void               ` | 设置为假 |

### useCycleList

| 返回值 | 类型                      | 说明           |
| ------ | ------------------------- | -------------- |
| value  | Ref\<T\>                  | 当前值         |
| index  | Ref\<number\>             | 当前索引       |
| next   | `() => void            `  | 下一个         |
| prev   | `() => void            `  | 上一个         |
| go     | `(index: number) => void` | 跳转到指定索引 |

## 代码位置

```
web/src/
└── composables/
    └── useVModel.ts    # VModel Composable
```
