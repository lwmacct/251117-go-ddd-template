# Component Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:41+4`
- [已实现功能](#已实现功能) `:45+37`
  - [插槽工具](#插槽工具) `:47+6`
  - [属性工具](#属性工具) `:53+6`
  - [组件信息](#组件信息) `:59+7`
  - [组件引用](#组件引用) `:66+5`
  - [暴露工具](#暴露工具) `:71+5`
  - [其他工具](#其他工具) `:76+6`
- [使用方式](#使用方式) `:82+275`
  - [插槽信息](#插槽信息) `:84+29`
  - [条件插槽](#条件插槽) `:113+17`
  - [插槽传递](#插槽传递) `:130+15`
  - [增强属性](#增强属性) `:145+23`
  - [类名合并](#类名合并) `:168+15`
  - [样式合并](#样式合并) `:183+17`
  - [组件信息](#组件信息-1) `:200+18`
  - [单组件引用](#单组件引用) `:218+29`
  - [多组件引用](#多组件引用) `:247+31`
  - [暴露方法](#暴露方法) `:278+37`
  - [强制更新](#强制更新) `:315+17`
  - [类型安全的事件](#类型安全的事件) `:332+25`
- [API](#api) `:357+42`
  - [useSlotsInfo](#useslotsinfo) `:359+10`
  - [useAttrsEnhanced](#useattrsenhanced) `:369+11`
  - [useComponentRef](#usecomponentref) `:380+8`
  - [useComponentRefs](#usecomponentrefs) `:388+11`
- [代码位置](#代码位置) `:399+7`

<!--TOC-->

## 需求背景

前端需要组件相关的工具函数，支持插槽处理、属性增强、组件引用、暴露方法等高级功能。

## 已实现功能

### 插槽工具

- `useSlotsInfo` - 获取插槽详细信息
- `useConditionalSlot` - 条件插槽渲染
- `useSlotPass` - 插槽传递

### 属性工具

- `useAttrsEnhanced` - 增强的属性处理
- `useClassMerge` - 类名合并
- `useStyleMerge` - 样式合并

### 组件信息

- `useComponentInfo` - 获取组件详细信息
- `useParentComponent` - 获取父组件
- `useRootComponent` - 获取根组件
- `useComponentType` - 获取组件类型信息

### 组件引用

- `useComponentRef` - 单组件引用
- `useComponentRefs` - 多组件引用

### 暴露工具

- `useExposeMethod` - 暴露方法
- `useExposeState` - 暴露状态和方法

### 其他工具

- `useForceUpdate` - 强制更新
- `useComponentEmit` - 类型安全的事件发射
- `useComponentProxy` - 获取组件代理

## 使用方式

### 插槽信息

```typescript
import { useSlotsInfo } from "@/composables/useComponent";

const { hasSlot, isEmpty, getSlotNames, slotInfo, slots } = useSlotsInfo();

// 检查插槽是否存在
if (hasSlot("header")) {
  // 渲染 header 插槽
}

// 检查插槽是否为空
if (!isEmpty("content")) {
  // content 插槽有内容
}

// 获取所有插槽名称
console.log(getSlotNames()); // ['default', 'header', 'footer']

// 获取详细信息
console.log(slotInfo.value);
// [
//   { name: 'default', exists: true, isEmpty: false },
//   { name: 'header', exists: true, isEmpty: true },
//   ...
// ]
```

### 条件插槽

```typescript
import { useConditionalSlot } from '@/composables/useComponent'

// 如果插槽不存在，使用默认内容
const renderHeader = useConditionalSlot('header', () => h('div', 'Default Header'))

// 在渲染函数中使用
render() {
  return h('div', [
    renderHeader(),
    // ...
  ])
}
```

### 插槽传递

```typescript
import { useSlotPass } from '@/composables/useComponent'

// 获取要传递给子组件的插槽
const passSlots = useSlotPass(['header', 'footer'])

// 或传递所有插槽
const allSlots = useSlotPass()

// 在模板中使用
<ChildComponent v-bind="passSlots" />
```

### 增强属性

```typescript
import { useAttrsEnhanced } from "@/composables/useComponent";

const { attrs, hasAttr, getAttr, attrsWithout, attrsOnly, attrNames } = useAttrsEnhanced();

// 检查属性
if (hasAttr("disabled")) {
  // ...
}

// 获取属性值（带默认值）
const className = getAttr("class", "default-class");

// 排除某些属性
const restAttrs = attrsWithout(["class", "style"]);
// 传递给内部元素: <div v-bind="restAttrs" />

// 只保留某些属性
const inputAttrs = attrsOnly(["type", "placeholder", "disabled"]);
```

### 类名合并

```typescript
import { useClassMerge } from '@/composables/useComponent'

// 基础类名 + 父组件传入的 class
const mergedClass = useClassMerge('btn', 'btn-primary')

// 如果父组件传入 class="custom large"
// 结果: "btn btn-primary custom large"

// 在模板中使用
<button :class="mergedClass">Click</button>
```

### 样式合并

```typescript
import { useStyleMerge } from '@/composables/useComponent'

// 基础样式 + 父组件传入的 style
const mergedStyle = useStyleMerge({
  color: 'red',
  fontSize: '14px'
})

// 如果父组件传入 style="background: blue"
// 结果: { color: 'red', fontSize: '14px', background: 'blue' }

<div :style="mergedStyle">Content</div>
```

### 组件信息

```typescript
import { useComponentInfo, useParentComponent, useRootComponent } from "@/composables/useComponent";

const { name, uid, isMounted, parent, root } = useComponentInfo();

console.log(`组件 ${name} (UID: ${uid})`);
console.log(`已挂载: ${isMounted.value}`);

// 获取父组件
const parentComponent = useParentComponent();
console.log("父组件:", parentComponent?.type.name);

// 获取根组件
const rootComponent = useRootComponent();
```

### 单组件引用

```typescript
import { useComponentRef } from "@/composables/useComponent";

interface FormInstance {
  validate: () => Promise<boolean>;
  reset: () => void;
}

const { ref: formRef, isLoaded, whenLoaded } = useComponentRef<FormInstance>();

// 在模板中: <Form ref="formRef" />

// 检查是否加载
if (isLoaded.value) {
  formRef.value?.validate();
}

// 等待组件加载
async function submit() {
  const form = await whenLoaded();
  const valid = await form.validate();
  if (valid) {
    // 提交表单
  }
}
```

### 多组件引用

```typescript
import { useComponentRefs } from '@/composables/useComponent'

interface InputInstance {
  focus: () => void
  clear: () => void
}

const { refs, setRef, getRef, getAllRefs, clearAll } = useComponentRefs<InputInstance>()

// 在模板中:
<Input
  v-for="field in fields"
  :key="field.id"
  :ref="el => setRef(field.id, el)"
/>

// 获取特定输入框
function focusField(fieldId: string) {
  const input = getRef(fieldId)
  input?.focus()
}

// 清除所有输入框
function clearAll() {
  getAllRefs().forEach(input => input.clear())
}
```

### 暴露方法

```typescript
import { useExposeMethod, useExposeState } from "@/composables/useComponent";

// 方式 1: 只暴露方法
const count = ref(0);

useExposeMethod({
  increment: () => count.value++,
  decrement: () => count.value--,
  reset: () => (count.value = 0),
  getValue: () => count.value,
});

// 父组件可以调用:
// childRef.value.increment()
// childRef.value.getValue()

// 方式 2: 暴露状态和方法
const name = ref("");
const loading = ref(false);

useExposeState(
  { name, loading },
  {
    setName: (n: string) => (name.value = n),
    startLoading: () => (loading.value = true),
    stopLoading: () => (loading.value = false),
  },
);

// 父组件可以访问:
// childRef.value.name.value
// childRef.value.setName('test')
```

### 强制更新

```typescript
import { useForceUpdate } from '@/composables/useComponent'

const { forceUpdate, updateKey } = useForceUpdate()

// 在模板中使用 key 触发完全重渲染
<HeavyComponent :key="updateKey" />

// 手动触发重渲染
function handleReset() {
  // ... 重置逻辑
  forceUpdate()
}
```

### 类型安全的事件

```typescript
import { useComponentEmit } from "@/composables/useComponent";

interface Events {
  "update:value": [value: string];
  change: [newValue: string, oldValue: string];
  submit: [];
  error: [error: Error];
}

const emit = useComponentEmit<Events>();

// 类型安全的事件发射
emit("update:value", "new value");
emit("change", "new", "old");
emit("submit");
emit("error", new Error("Something went wrong"));

// 错误示例（TypeScript 会报错）:
// emit('update:value', 123) // 类型错误
// emit('change', 'only one arg') // 参数数量错误
```

## API

### useSlotsInfo

| 返回值       | 类型                      | 说明             |
| ------------ | ------------------------- | ---------------- |
| hasSlot      | `(name) => boolean`       | 检查插槽是否存在 |
| isEmpty      | `(name) => boolean`       | 检查插槽是否为空 |
| getSlotNames | `() => string[]`          | 获取所有插槽名称 |
| slotInfo     | ComputedRef\<SlotInfo[]\> | 插槽详细信息     |
| slots        | Slots                     | 原始插槽对象     |

### useAttrsEnhanced

| 返回值       | 类型                    | 说明             |
| ------------ | ----------------------- | ---------------- |
| attrs        | object                  | 原始属性对象     |
| hasAttr      | `(name) => boolean`     | 检查属性是否存在 |
| getAttr      | `(name, default?) => T` | 获取属性值       |
| attrsWithout | `(names) => object`     | 排除指定属性     |
| attrsOnly    | `(names) => object`     | 只包含指定属性   |
| attrNames    | ComputedRef\<string[]\> | 属性名列表       |

### useComponentRef

| 返回值     | 类型                   | 说明         |
| ---------- | ---------------------- | ------------ |
| ref        | Ref\<T \| null\>       | 组件引用     |
| isLoaded   | ComputedRef\<boolean\> | 是否已加载   |
| whenLoaded | `() => Promise\<T\>`   | 等待加载完成 |

### useComponentRefs

| 返回值     | 类型                | 说明         |
| ---------- | ------------------- | ------------ |
| refs       | Ref\<Map\>          | 引用映射     |
| setRef     | `(key, el) => void` | 设置引用     |
| getRef     | `(key) => T`        | 获取引用     |
| getAllRefs | `() => T[]`         | 获取所有引用 |
| clearRef   | `(key) => void`     | 清除引用     |
| clearAll   | `() => void`        | 清除所有引用 |

## 代码位置

```
web/src/
└── composables/
    └── useComponent.ts    # Component Composable
```
