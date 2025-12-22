# 数字格式化工具

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:20+4`
- [已实现功能](#已实现功能) `:24+16`
  - [格式化函数](#格式化函数) `:26+9`
  - [辅助函数](#辅助函数) `:35+5`
- [使用方式](#使用方式) `:40+12`
- [代码位置](#代码位置) `:52+7`

<!--TOC-->

## 需求背景

项目中多处需要格式化数字（货币、百分比、文件大小等），需要统一的数字格式化工具。

## 已实现功能

### 格式化函数

- `formatNumber` - 添加千分位分隔符
- `formatCurrency` - 格式化货币
- `formatPercent` - 格式化百分比
- `formatCompact` - 大数字缩写（万、亿）
- `formatFileSize` - 文件大小
- `formatDuration` - 持续时间

### 辅助函数

- `toOrdinal` - 序数词（第N）
- `clamp` - 数字范围限制

## 使用方式

```typescript
import { formatNumber, formatCurrency, formatPercent, formatCompact, formatFileSize } from "@/utils/number";

formatNumber(1234567); // "1,234,567"
formatCurrency(99.9); // "¥99.90"
formatPercent(0.156, { multiply: true }); // "15.6%"
formatCompact(12345678); // "1234.6万"
formatFileSize(1073741824); // "1.00 GB"
```

## 代码位置

```
web/src/
└── utils/
    └── number.ts    # 数字格式化工具
```
