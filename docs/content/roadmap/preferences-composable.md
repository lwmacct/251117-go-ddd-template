# Preferences Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:34+4`
- [已实现功能](#已实现功能) `:38+16`
  - [偏好检测](#偏好检测) `:40+9`
  - [主题控制](#主题控制) `:49+5`
- [使用方式](#使用方式) `:54+115`
  - [检测深色模式偏好](#检测深色模式偏好) `:56+13`
  - [深色模式控制](#深色模式控制) `:69+23`
  - [颜色模式管理](#颜色模式管理) `:92+26`
  - [语言偏好](#语言偏好) `:118+14`
  - [减少动画偏好](#减少动画偏好) `:132+11`
  - [对比度偏好](#对比度偏好) `:143+13`
  - [透明度偏好](#透明度偏好) `:156+13`
- [API](#api) `:169+68`
  - [usePreferredDark](#usepreferreddark) `:171+7`
  - [useDark](#usedark) `:178+20`
  - [useColorMode](#usecolormode) `:198+17`
  - [usePreferredLanguage](#usepreferredlanguage) `:215+8`
  - [usePreferredReducedMotion](#usepreferredreducedmotion) `:223+7`
  - [usePreferredContrast](#usepreferredcontrast) `:230+7`
- [代码位置](#代码位置) `:237+7`

<!--TOC-->

## 需求背景

前端需要响应式检测和管理用户偏好设置，包括深色模式、语言、动画偏好等系统级设置。

## 已实现功能

### 偏好检测

- `usePreferredDark` - 检测深色模式偏好
- `usePreferredLanguage` - 获取语言偏好
- `usePreferredReducedMotion` - 检测减少动画偏好
- `usePreferredContrast` - 检测对比度偏好
- `usePreferredColorScheme` - 检测颜色方案偏好
- `usePreferredTransparency` - 检测透明度偏好

### 主题控制

- `useDark` - 深色模式控制（支持持久化）
- `useColorMode` - 多模式颜色主题管理

## 使用方式

### 检测深色模式偏好

```typescript
import { usePreferredDark } from "@/composables/usePreferences";

const { isDark, isSupported } = usePreferredDark();

// 响应式监听系统深色模式
watch(isDark, (dark) => {
  console.log("系统深色模式:", dark);
});
```

### 深色模式控制

```typescript
import { useDark } from "@/composables/usePreferences";

const { isDark, toggle, setDark, setLight, systemIsDark } = useDark({
  initialValue: "auto", // 'auto' | true | false
  storageKey: "theme-dark",
  selector: "html",
  attribute: "class",
  valueDark: "dark",
  valueLight: "",
  onChanged: (dark) => console.log("主题变更:", dark),
});

// 切换深色模式
toggle();

// 直接设置
setDark();
setLight();
```

### 颜色模式管理

```typescript
import { useColorMode } from "@/composables/usePreferences";

const { mode, resolvedMode, setMode, cycle } = useColorMode({
  modes: ["light", "dark", "auto", "sepia"],
  initialValue: "auto",
  storageKey: "color-mode",
  selector: "html",
  attribute: "data-theme",
});

// 获取当前模式
console.log(mode.value); // 'auto'

// 获取实际应用的模式
console.log(resolvedMode.value); // 'light' 或 'dark'

// 切换模式
cycle(); // light -> dark -> auto -> sepia -> light

// 设置特定模式
setMode("dark");
```

### 语言偏好

```typescript
import { usePreferredLanguage } from "@/composables/usePreferences";

const { language, languages, isSupported } = usePreferredLanguage();

// 获取首选语言
console.log(language.value); // 'zh-CN'

// 获取所有语言
console.log(languages.value); // ['zh-CN', 'en-US', 'ja']
```

### 减少动画偏好

```typescript
import { usePreferredReducedMotion } from "@/composables/usePreferences";

const { isReduced, isSupported } = usePreferredReducedMotion();

// 根据用户偏好调整动画
const transition = computed(() => (isReduced.value ? "none" : "all 0.3s ease"));
```

### 对比度偏好

```typescript
import { usePreferredContrast } from "@/composables/usePreferences";

const { contrast, isSupported } = usePreferredContrast();

// contrast.value: 'more' | 'less' | 'custom' | 'no-preference'
if (contrast.value === "more") {
  // 使用高对比度样式
}
```

### 透明度偏好

```typescript
import { usePreferredTransparency } from "@/composables/usePreferences";

const { isReduced, isSupported } = usePreferredTransparency();

// 如果用户偏好减少透明度
if (isReduced.value) {
  // 使用不透明背景
}
```

## API

### usePreferredDark

| 返回值      | 类型                   | 说明             |
| ----------- | ---------------------- | ---------------- |
| isDark      | Ref\<boolean\>         | 是否偏好深色模式 |
| isSupported | ComputedRef\<boolean\> | 是否支持         |

### useDark

| 选项         | 类型                        | 默认值         | 说明       |
| ------------ | --------------------------- | -------------- | ---------- |
| initialValue | boolean \| 'auto'           | 'auto'         | 初始值     |
| storageKey   | string                      | 'vue-use-dark' | 存储键名   |
| selector     | string                      | 'html'         | 目标选择器 |
| attribute    | string                      | 'class'        | 属性名     |
| valueDark    | string                      | 'dark'         | 深色值     |
| valueLight   | string                      | ''             | 浅色值     |
| onChanged    | `(isDark: boolean) => void` | -              | 变化回调   |

| 返回值       | 类型             | 说明           |
| ------------ | ---------------- | -------------- |
| isDark       | Ref\<boolean\>   | 是否为深色模式 |
| toggle       | `() => void    ` | 切换模式       |
| setDark      | `() => void    ` | 设置为深色     |
| setLight     | `() => void    ` | 设置为浅色     |
| systemIsDark | Ref\<boolean\>   | 系统偏好       |

### useColorMode

| 选项         | 类型        | 默认值                  | 说明       |
| ------------ | ----------- | ----------------------- | ---------- |
| initialValue | ColorMode   | 'auto'                  | 初始模式   |
| storageKey   | string      | 'vue-use-color-mode'    | 存储键名   |
| modes        | ColorMode[] | ['light','dark','auto'] | 可用模式   |
| selector     | string      | 'html'                  | 目标选择器 |
| attribute    | string      | 'data-theme'            | 属性名     |

| 返回值       | 类型                              | 说明         |
| ------------ | --------------------------------- | ------------ |
| mode         | Ref\<ColorMode\>                  | 当前模式     |
| resolvedMode | ComputedRef\<'light'\|'dark'\>    | 实际应用模式 |
| setMode      | `(mode: ColorMode) => void      ` | 设置模式     |
| cycle        | `() => void                     ` | 循环切换     |

### usePreferredLanguage

| 返回值      | 类型                     | 说明     |
| ----------- | ------------------------ | -------- |
| language    | Ref\<string\>            | 首选语言 |
| languages   | Ref\<readonly string[]\> | 语言列表 |
| isSupported | ComputedRef\<boolean\>   | 是否支持 |

### usePreferredReducedMotion

| 返回值      | 类型                   | 说明             |
| ----------- | ---------------------- | ---------------- |
| isReduced   | Ref\<boolean\>         | 是否偏好减少动画 |
| isSupported | ComputedRef\<boolean\> | 是否支持         |

### usePreferredContrast

| 返回值      | 类型                                             | 说明       |
| ----------- | ------------------------------------------------ | ---------- |
| contrast    | Ref\<'more'\|'less'\|'custom'\|'no-preference'\> | 对比度偏好 |
| isSupported | ComputedRef\<boolean\>                           | 是否支持   |

## 代码位置

```
web/src/
└── composables/
    └── usePreferences.ts    # Preferences Composable
```
