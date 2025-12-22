# Title Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:35+4`
- [已实现功能](#已实现功能) `:39+23`
  - [标题管理](#标题管理) `:41+6`
  - [文档元素](#文档元素) `:47+5`
  - [可见性检测](#可见性检测) `:52+5`
  - [资源加载](#资源加载) `:57+5`
- [使用方式](#使用方式) `:62+153`
  - [基础标题管理](#基础标题管理) `:64+23`
  - [模板标题](#模板标题) `:87+18`
  - [Favicon 管理](#favicon-管理) `:105+16`
  - [文档可见性](#文档可见性) `:121+18`
  - [页面离开检测](#页面离开检测) `:139+17`
  - [Head 管理](#head-管理) `:156+16`
  - [动态脚本加载](#动态脚本加载) `:172+23`
  - [动态样式表加载](#动态样式表加载) `:195+20`
- [API](#api) `:215+55`
  - [useTitle](#usetitle) `:217+13`
  - [useTitleTemplate](#usetitletemplate) `:230+14`
  - [useDocumentVisibility](#usedocumentvisibility) `:244+8`
  - [useScript](#usescript) `:252+18`
- [代码位置](#代码位置) `:270+7`

<!--TOC-->

## 需求背景

前端需要响应式管理文档标题、Favicon、动态脚本和样式表加载等功能。

## 已实现功能

### 标题管理

- `useTitle` - 响应式文档标题管理
- `useTitleTemplate` - 带模板的标题管理
- `useDocumentTitle` - 简化的标题设置

### 文档元素

- `useFavicon` - Favicon 管理
- `useHead` - 简单的 Head 管理

### 可见性检测

- `usePageLeave` - 页面离开检测
- `useDocumentVisibility` - 文档可见性状态

### 资源加载

- `useScript` - 动态脚本加载
- `useStylesheet` - 动态样式表加载

## 使用方式

### 基础标题管理

```typescript
import { useTitle } from "@/composables/useTitle";

// 设置标题
const { title } = useTitle("My Page");

// 响应式更新
title.value = "New Title";

// 使用模板
const { title } = useTitle("Home", {
  template: "%s | My App",
});
// 结果: 'Home | My App'

// 组件卸载时恢复原标题
const { title } = useTitle("Temp Page", {
  restoreOnUnmount: true,
});
```

### 模板标题

```typescript
import { useTitleTemplate } from "@/composables/useTitle";

const { pageTitle, fullTitle, setPageTitle } = useTitleTemplate({
  siteName: "My App",
  separator: " - ",
  siteNamePosition: "suffix",
});

setPageTitle("Home");
console.log(fullTitle.value); // 'Home - My App'

setPageTitle("About");
console.log(fullTitle.value); // 'About - My App'
```

### Favicon 管理

```typescript
import { useFavicon } from "@/composables/useTitle";

const { favicon } = useFavicon("/favicon.ico");

// 动态切换
favicon.value = "/favicon-dark.ico";

// 组件卸载时恢复
const { favicon } = useFavicon("/temp-icon.ico", {
  restoreOnUnmount: true,
});
```

### 文档可见性

```typescript
import { useDocumentVisibility } from "@/composables/useTitle";

const { visibility, isVisible } = useDocumentVisibility();

watch(isVisible, (visible) => {
  if (visible) {
    // 恢复动画、重新连接
    resumeVideo();
  } else {
    // 暂停动画、断开连接
    pauseVideo();
  }
});
```

### 页面离开检测

```typescript
import { usePageLeave } from "@/composables/useTitle";

const { isLeft } = usePageLeave(() => {
  // 用户离开页面时触发
  saveDraft();
});

watch(isLeft, (left) => {
  if (left) {
    showRetentionPopup();
  }
});
```

### Head 管理

```typescript
import { useHead } from "@/composables/useTitle";

useHead({
  title: "My Page",
  meta: [
    { name: "description", content: "Page description" },
    { property: "og:title", content: "My Page" },
    { property: "og:description", content: "Page description" },
  ],
  link: [{ rel: "canonical", href: "https://example.com/page" }],
});
```

### 动态脚本加载

```typescript
import { useScript } from "@/composables/useTitle";

// 立即加载
const { isLoaded, error } = useScript("https://example.com/analytics.js", {
  immediate: true,
  async: true,
});

// 延迟加载
const { load, isLoading } = useScript("https://example.com/heavy.js", {
  immediate: false,
});

// 需要时加载
async function loadHeavyFeature() {
  await load();
  // 脚本已加载
}
```

### 动态样式表加载

```typescript
import { useStylesheet } from "@/composables/useTitle";

// 加载主题
const { isLoaded } = useStylesheet("/styles/theme.css");

// 暗色主题（媒体查询）
useStylesheet("/styles/dark.css", {
  media: "(prefers-color-scheme: dark)",
});

// 动态切换
const { load, unload } = useStylesheet("/styles/feature.css", {
  immediate: false,
  removeOnUnmount: true,
});
```

## API

### useTitle

| 选项             | 类型    | 默认值 | 说明                |
| ---------------- | ------- | ------ | ------------------- |
| template         | string  | -      | 标题模板（%s 占位） |
| restoreOnUnmount | boolean | false  | 卸载时恢复原标题    |
| observe          | boolean | false  | 观察外部标题变化    |

| 返回值      | 类型                   | 说明     |
| ----------- | ---------------------- | -------- |
| title       | Ref\<string\>          | 当前标题 |
| isSupported | ComputedRef\<boolean\> | 是否支持 |

### useTitleTemplate

| 选项             | 类型                 | 默认值   | 说明         |
| ---------------- | -------------------- | -------- | ------------ |
| separator        | string               | ' \| '   | 分隔符       |
| siteName         | string               | ''       | 网站名称     |
| siteNamePosition | 'prefix' \| 'suffix' | 'suffix' | 网站名称位置 |

| 返回值       | 类型                      | 说明         |
| ------------ | ------------------------- | ------------ |
| pageTitle    | Ref\<string\>             | 页面标题     |
| fullTitle    | ComputedRef\<string\>     | 完整标题     |
| setPageTitle | `(title: string) => void` | 设置页面标题 |

### useDocumentVisibility

| 返回值      | 类型                           | 说明     |
| ----------- | ------------------------------ | -------- |
| visibility  | Ref\<DocumentVisibilityState\> | 可见状态 |
| isVisible   | ComputedRef\<boolean\>         | 是否可见 |
| isSupported | ComputedRef\<boolean\>         | 是否支持 |

### useScript

| 选项            | 类型                           | 默认值 | 说明         |
| --------------- | ------------------------------ | ------ | ------------ |
| immediate       | boolean                        | true   | 是否立即加载 |
| async           | boolean                        | true   | 异步加载     |
| defer           | boolean                        | false  | 延迟加载     |
| crossOrigin     | 'anonymous'\|'use-credentials' | -      | 跨域设置     |
| removeOnUnmount | boolean                        | false  | 卸载时移除   |

| 返回值    | 类型                      | 说明         |
| --------- | ------------------------- | ------------ |
| isLoading | Ref\<boolean\>            | 是否正在加载 |
| isLoaded  | Ref\<boolean\>            | 是否已加载   |
| error     | Ref\<Error \| null\>      | 加载错误     |
| load      | `() => Promise          ` | 手动加载     |
| unload    | `() => void             ` | 移除脚本     |

## 代码位置

```
web/src/
└── composables/
    └── useTitle.ts    # Title Composable
```
