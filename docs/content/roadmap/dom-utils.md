# DOM 工具函数

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:35+4`
- [已实现功能](#已实现功能) `:39+81`
  - [元素查询](#元素查询) `:41+7`
  - [类名操作](#类名操作) `:48+8`
  - [样式操作](#样式操作) `:56+8`
  - [属性操作](#属性操作) `:64+9`
  - [尺寸和位置](#尺寸和位置) `:73+8`
  - [滚动操作](#滚动操作) `:81+8`
  - [焦点操作](#焦点操作) `:89+7`
  - [元素操作](#元素操作) `:96+10`
  - [事件工具](#事件工具) `:106+7`
  - [可见性](#可见性) `:113+7`
- [使用方式](#使用方式) `:120+82`
  - [元素查询](#元素查询-1) `:122+13`
  - [类名和样式](#类名和样式) `:135+18`
  - [事件处理](#事件处理) `:153+20`
  - [滚动操作](#滚动操作-1) `:173+17`
  - [元素创建](#元素创建) `:190+12`
- [API](#api) `:202+14`
  - [主要函数](#主要函数) `:204+12`
- [代码位置](#代码位置) `:216+7`

<!--TOC-->

## 需求背景

前端需要统一的 DOM 操作工具，用于元素查询、样式操作、事件处理等场景。

## 已实现功能

### 元素查询

- `getElement` - 获取单个元素
- `getElements` - 获取所有匹配元素
- `elementExists` - 检查元素是否存在
- `waitForElement` - 等待元素出现

### 类名操作

- `addClass` - 添加类名
- `removeClass` - 移除类名
- `toggleClass` - 切换类名
- `hasClass` - 检查类名
- `replaceClass` - 替换类名

### 样式操作

- `getStyle` - 获取计算样式
- `setStyle` - 设置样式
- `removeStyle` - 移除样式
- `getCSSVariable` - 获取 CSS 变量
- `setCSSVariable` - 设置 CSS 变量

### 属性操作

- `getAttribute` - 获取属性
- `setAttribute` - 设置属性
- `removeAttribute` - 移除属性
- `hasAttribute` - 检查属性
- `getDataAttribute` - 获取 data 属性
- `setDataAttribute` - 设置 data 属性

### 尺寸和位置

- `getRect` - 获取边界矩形
- `getSize` - 获取尺寸
- `getOffset` - 获取文档位置
- `getWindowSize` - 获取窗口尺寸
- `getScrollPosition` - 获取滚动位置

### 滚动操作

- `scrollTo` - 滚动到位置/元素
- `scrollToTop` - 滚动到顶部
- `scrollToBottom` - 滚动到底部
- `isInViewport` - 检查是否在视口内
- `isPartiallyInViewport` - 检查是否部分可见

### 焦点操作

- `focus` - 设置焦点
- `blur` - 移除焦点
- `getActiveElement` - 获取焦点元素
- `hasFocus` - 检查焦点状态

### 元素操作

- `createElement` - 创建元素
- `removeElement` - 移除元素
- `cloneElement` - 克隆元素
- `insertBefore` - 前插入
- `insertAfter` - 后插入
- `wrap` - 包裹元素
- `unwrap` - 解除包裹

### 事件工具

- `on` - 添加事件监听（返回清理函数）
- `off` - 移除事件监听
- `once` - 一次性事件
- `trigger` - 触发事件

### 可见性

- `show` - 显示元素
- `hide` - 隐藏元素
- `toggle` - 切换显示
- `isVisible` - 检查是否可见

## 使用方式

### 元素查询

```typescript
import { getElement, waitForElement } from "@/utils/dom";

// 获取元素（支持选择器或元素）
const el = getElement("#my-id");
const el2 = getElement(document.body);

// 等待动态元素
const dynamicEl = await waitForElement("#dynamic-element", 5000);
```

### 类名和样式

```typescript
import { addClass, toggleClass, setStyle } from "@/utils/dom";

// 类名操作
addClass("#box", "active", "visible");
toggleClass("#box", "expanded");

// 样式操作
setStyle("#box", "backgroundColor", "red");
setStyle("#box", {
  backgroundColor: "red",
  fontSize: "16px",
  padding: 20, // 自动添加 px
});
```

### 事件处理

```typescript
import { on, once, trigger } from "@/utils/dom";

// 添加事件（返回清理函数）
const cleanup = on("#button", "click", (e) => {
  console.log("Clicked!");
});

// 组件卸载时清理
onUnmounted(() => cleanup());

// 一次性事件
once("#button", "click", handleClick);

// 触发自定义事件
trigger("#element", "my-event", { detail: { foo: "bar" } });
```

### 滚动操作

```typescript
import { scrollTo, scrollToTop, isInViewport } from "@/utils/dom";

// 滚动到元素
scrollTo("#section", { behavior: "smooth" });

// 滚动到顶部
scrollToTop();

// 检查可见性
if (isInViewport("#lazy-image")) {
  loadImage();
}
```

### 元素创建

```typescript
import { createElement, wrap } from "@/utils/dom";

// 创建元素
const div = createElement("div", { className: "box", id: "my-box" }, "Hello");

// 包裹元素
wrap("#content", "div", { className: "wrapper" });
```

## API

### 主要函数

| 函数                          | 说明             |
| ----------------------------- | ---------------- |
| getElement(target)            | 获取元素         |
| addClass(target, ...classes)  | 添加类名         |
| setStyle(target, prop, value) | 设置样式         |
| on(target, event, handler)    | 添加事件监听     |
| scrollTo(target, options)     | 滚动到位置/元素  |
| isInViewport(target)          | 检查是否在视口内 |
| createElement(tag, attrs)     | 创建元素         |

## 代码位置

```
web/src/
└── utils/
    └── dom.ts    # DOM 工具函数
```
