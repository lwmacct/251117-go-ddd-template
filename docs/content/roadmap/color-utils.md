# 颜色工具函数

<!--TOC-->

- [需求背景](#需求背景) `:29:32`
- [已实现功能](#已实现功能) `:33:34`
  - [解析函数](#解析函数) `:35:41`
  - [转换函数](#转换函数) `:42:49`
  - [格式化函数](#格式化函数) `:50:54`
  - [颜色操作](#颜色操作) `:55:66`
  - [颜色分析](#颜色分析) `:67:74`
  - [颜色生成](#颜色生成) `:75:80`
  - [其他](#其他) `:81:85`
- [使用方式](#使用方式) `:86:87`
  - [颜色转换](#颜色转换) `:88:101`
  - [颜色调整](#颜色调整) `:102:117`
  - [颜色分析](#颜色分析-1) `:118:134`
  - [生成调色板](#生成调色板) `:135:148`
- [API](#api) `:149:150`
  - [主要函数](#主要函数) `:151:163`
- [代码位置](#代码位置) `:164:170`

<!--TOC-->

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

## 需求背景

前端需要处理颜色转换、调整和生成，用于主题定制、动态样式等场景。

## 已实现功能

### 解析函数

- `parseHex` - 解析十六进制颜色
- `parseRgb` - 解析 RGB/RGBA 字符串
- `parseHsl` - 解析 HSL/HSLA 字符串
- `parseColor` - 解析任意格式颜色

### 转换函数

- `rgbToHex` - RGB 转十六进制
- `rgbToHsl` - RGB 转 HSL
- `hslToRgb` - HSL 转 RGB
- `rgbToHsv` - RGB 转 HSV
- `hsvToRgb` - HSV 转 RGB

### 格式化函数

- `formatRgb` - 格式化为 RGB 字符串
- `formatHsl` - 格式化为 HSL 字符串

### 颜色操作

- `lighten` - 调亮颜色
- `darken` - 调暗颜色
- `saturate` - 增加饱和度
- `desaturate` - 降低饱和度
- `setAlpha` - 设置透明度
- `invert` - 反转颜色
- `grayscale` - 转为灰度
- `mix` - 混合两种颜色
- `complement` - 获取补色

### 颜色分析

- `getLuminance` - 计算亮度
- `getContrast` - 计算对比度
- `isDark` - 判断是否为深色
- `isLight` - 判断是否为浅色
- `getTextColor` - 获取适合的文本颜色

### 颜色生成

- `randomColor` - 生成随机颜色
- `generateGradient` - 生成渐变色数组
- `generatePalette` - 生成调色板

### 其他

- `getNamedColor` - 获取命名颜色
- `isValidColor` - 检查颜色有效性

## 使用方式

### 颜色转换

```typescript
import { parseColor, rgbToHex, rgbToHsl } from "@/utils/color";

// 解析任意格式
const rgba = parseColor("#ff0000");
// { r: 255, g: 0, b: 0, a: 1 }

// 转换格式
rgbToHex({ r: 255, g: 0, b: 0 }); // '#ff0000'
rgbToHsl({ r: 255, g: 0, b: 0 }); // { h: 0, s: 100, l: 50 }
```

### 颜色调整

```typescript
import { lighten, darken, setAlpha, mix } from "@/utils/color";

// 调亮/调暗
lighten("#ff0000", 20); // '#ff6666'
darken("#ff0000", 20); // '#990000'

// 设置透明度
setAlpha("#ff0000", 0.5); // '#ff000080'

// 混合颜色
mix("#ff0000", "#0000ff", 0.5); // '#800080'
```

### 颜色分析

```typescript
import { isDark, getTextColor, getContrast } from "@/utils/color";

// 判断深浅
isDark("#000000"); // true
isDark("#ffffff"); // false

// 获取适合的文本颜色
getTextColor("#000000"); // '#ffffff'
getTextColor("#ffffff"); // '#000000'

// 计算对比度
getContrast("#ffffff", "#000000"); // 21
```

### 生成调色板

```typescript
import { generatePalette, generateGradient } from "@/utils/color";

// 生成完整调色板
const palette = generatePalette("#3498db");
// { 50: '#e8f4fc', 100: '#d1e9f9', ..., 900: '#1a4d6e' }

// 生成渐变色
const gradient = generateGradient("#ff0000", "#0000ff", 5);
// ['#ff0000', '#bf003f', '#7f007f', '#3f00bf', '#0000ff']
```

## API

### 主要函数

| 函数                       | 说明             |
| -------------------------- | ---------------- |
| parseColor(color)          | 解析任意格式颜色 |
| rgbToHex(rgb)              | RGB 转十六进制   |
| lighten(color, amount)     | 调亮颜色         |
| darken(color, amount)      | 调暗颜色         |
| mix(color1, color2, ratio) | 混合颜色         |
| isDark(color)              | 判断是否为深色   |
| getTextColor(bgColor)      | 获取适合文本颜色 |
| generatePalette(baseColor) | 生成调色板       |

## 代码位置

```
web/src/
└── utils/
    └── color.ts    # 颜色工具函数
```
