# 滚动锁定 Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:30+4`
- [已实现功能](#已实现功能) `:34+29`
  - [useScrollLock](#usescrolllock) `:36+7`
  - [useScrollLockWhenTrue](#usescrolllockwhentrue) `:43+5`
  - [useElementScrollLock](#useelementscrolllock) `:48+5`
  - [useScrollPosition](#usescrollposition) `:53+5`
  - [useScrollDirection](#usescrolldirection) `:58+5`
- [使用方式](#使用方式) `:63+70`
  - [基础用法](#基础用法) `:65+20`
  - [响应式锁定](#响应式锁定) `:85+16`
  - [滚动方向检测](#滚动方向检测) `:101+15`
  - [滚动位置控制](#滚动位置控制) `:116+17`
- [API](#api) `:133+19`
  - [useScrollLock 返回值](#usescrolllock-返回值) `:135+9`
  - [useScrollDirection 返回值](#usescrolldirection-返回值) `:144+8`
- [代码位置](#代码位置) `:152+7`

<!--TOC-->

## 需求背景

模态框、抽屉等组件需要在打开时锁定页面滚动，关闭时恢复滚动。

## 已实现功能

### useScrollLock

- 页面滚动锁定/解锁
- 支持嵌套锁定
- 保持滚动条占位防抖动
- 恢复滚动位置

### useScrollLockWhenTrue

- 响应式滚动锁定
- 根据 ref 值自动锁定

### useElementScrollLock

- 元素滚动锁定
- 锁定指定容器而非整个页面

### useScrollPosition

- 滚动位置保存/恢复
- 滚动到顶部/底部

### useScrollDirection

- 滚动方向检测
- 向上/向下滚动判断

## 使用方式

### 基础用法

```typescript
import { useScrollLock } from "@/composables/useScrollLock";

const { isLocked, lock, unlock } = useScrollLock();

// 显示模态框时锁定
const showModal = () => {
  lock();
  modalVisible.value = true;
};

// 关闭模态框时解锁
const hideModal = () => {
  unlock();
  modalVisible.value = false;
};
```

### 响应式锁定

```typescript
import { ref } from "vue";
import { useScrollLockWhenTrue } from "@/composables/useScrollLock";

const isModalOpen = ref(false);
useScrollLockWhenTrue(isModalOpen);

// 打开模态框时自动锁定滚动
isModalOpen.value = true;

// 关闭时自动解锁
isModalOpen.value = false;
```

### 滚动方向检测

```typescript
import { useScrollDirection } from "@/composables/useScrollLock";

const { direction, isScrollingUp, isScrollingDown } = useScrollDirection();

// 向下滚动时隐藏导航栏
watch(isScrollingDown, (scrollingDown) => {
  if (scrollingDown) {
    hideNavbar();
  }
});
```

### 滚动位置控制

```typescript
import { useScrollPosition } from "@/composables/useScrollLock";

const { save, restore, scrollToTop, scrollToBottom } = useScrollPosition();

// 保存当前位置
save();

// 跳转到其他位置后恢复
restore();

// 滚动到顶部
scrollToTop();
```

## API

### useScrollLock 返回值

| 属性     | 类型           | 说明       |
| -------- | -------------- | ---------- |
| isLocked | `Ref<boolean>` | 是否已锁定 |
| lock     | `() => void`   | 锁定滚动   |
| unlock   | `() => void`   | 解锁滚动   |
| toggle   | `() => void`   | 切换锁定   |

### useScrollDirection 返回值

| 属性            | 类型                            | 说明         |
| --------------- | ------------------------------- | ------------ |
| direction       | `Ref<'up' \| 'down' \| 'none'>` | 滚动方向     |
| isScrollingUp   | `Ref<boolean>`                  | 是否向上滚动 |
| isScrollingDown | `Ref<boolean>`                  | 是否向下滚动 |

## 代码位置

```
web/src/
└── composables/
    └── useScrollLock.ts    # 滚动锁定
```
