# 焦点陷阱 Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:29+4`
- [已实现功能](#已实现功能) `:33+25`
  - [useFocusTrap](#usefocustrap) `:35+9`
  - [useFocusTrapWhenTrue](#usefocustrapwhentrue) `:44+5`
  - [useFocusReturn](#usefocusreturn) `:49+5`
  - [useAutoFocus](#useautofocus) `:54+4`
- [使用方式](#使用方式) `:58+63`
  - [基础用法](#基础用法) `:60+22`
  - [响应式用法](#响应式用法) `:82+12`
  - [配置选项](#配置选项) `:94+13`
  - [焦点返回](#焦点返回) `:107+14`
- [API](#api) `:121+24`
  - [useFocusTrap 选项](#usefocustrap-选项) `:123+11`
  - [useFocusTrap 返回值](#usefocustrap-返回值) `:134+11`
- [代码位置](#代码位置) `:145+7`

<!--TOC-->

## 需求背景

模态框等组件需要将键盘焦点限制在容器内，以符合无障碍访问 (A11y) 标准。

## 已实现功能

### useFocusTrap

- 焦点陷阱激活/停用
- Tab 键循环导航
- 自动聚焦首个元素
- 恢复之前焦点
- Escape 键退出
- 点击外部退出

### useFocusTrapWhenTrue

- 响应式焦点陷阱
- 根据 ref 值自动激活

### useFocusReturn

- 焦点返回管理
- 保存/恢复焦点

### useAutoFocus

- 自动聚焦元素

## 使用方式

### 基础用法

```typescript
import { ref } from "vue";
import { useFocusTrap } from "@/composables/useFocusTrap";

const modalRef = ref<HTMLElement>();
const { activate, deactivate, isActive } = useFocusTrap(modalRef);

// 打开模态框时激活焦点陷阱
const openModal = () => {
  modalVisible.value = true;
  activate();
};

// 关闭时停用
const closeModal = () => {
  deactivate();
  modalVisible.value = false;
};
```

### 响应式用法

```typescript
import { useFocusTrapWhenTrue } from "@/composables/useFocusTrap";

const modalRef = ref<HTMLElement>();
const isOpen = ref(false);

// 自动跟随 isOpen 状态
useFocusTrapWhenTrue(modalRef, isOpen);
```

### 配置选项

```typescript
const { activate, deactivate } = useFocusTrap(modalRef, {
  autoFocus: true, // 自动聚焦第一个元素
  restoreFocus: true, // 停用时恢复焦点
  initialFocus: ".my-input", // 初始聚焦的元素
  escapeDeactivates: true, // Escape 键停用
  clickOutsideDeactivates: true, // 点击外部停用
  onDeactivate: () => closeModal(),
});
```

### 焦点返回

```typescript
import { useFocusReturn } from "@/composables/useFocusTrap";

const { save, restore } = useFocusReturn();

// 打开模态框前
save();

// 关闭后恢复
restore();
```

## API

### useFocusTrap 选项

| 选项                    | 类型    | 默认值 | 说明               |
| ----------------------- | ------- | ------ | ------------------ |
| immediate               | boolean | false  | 是否立即激活       |
| autoFocus               | boolean | true   | 激活时自动聚焦     |
| restoreFocus            | boolean | true   | 停用时恢复焦点     |
| initialFocus            | string  | -      | 初始聚焦元素选择器 |
| escapeDeactivates       | boolean | false  | Escape 键停用      |
| clickOutsideDeactivates | boolean | false  | 点击外部停用       |

### useFocusTrap 返回值

| 属性       | 类型           | 说明         |
| ---------- | -------------- | ------------ |
| isActive   | `Ref<boolean>` | 是否激活     |
| activate   | `() => void`   | 激活         |
| deactivate | `() => void`   | 停用         |
| toggle     | `() => void`   | 切换         |
| focusFirst | `() => void`   | 聚焦首个元素 |
| focusLast  | `() => void`   | 聚焦最后元素 |

## 代码位置

```
web/src/
└── composables/
    └── useFocusTrap.ts    # 焦点陷阱
```
