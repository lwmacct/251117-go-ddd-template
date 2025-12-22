# 日期时间格式化工具

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:21+4`
- [已实现功能](#已实现功能) `:25+15`
  - [格式化函数](#格式化函数) `:27+8`
  - [辅助函数](#辅助函数) `:35+5`
- [使用方式](#使用方式) `:40+21`
- [配置选项](#配置选项) `:61+11`
- [代码位置](#代码位置) `:72+7`

<!--TOC-->

## 需求背景

项目中多处需要格式化日期时间，但格式不统一，代码重复。需要统一的日期时间格式化工具。

## 已实现功能

### 格式化函数

- `formatDateTime` - 格式化日期时间
- `formatDate` - 格式化日期（不含时间）
- `formatTime` - 格式化时间（不含日期）
- `formatRelativeTime` - 相对时间（如"5 分钟前"）
- `formatSmart` - 智能格式化

### 辅助函数

- `isToday` - 判断是否是今天
- `isYesterday` - 判断是否是昨天

## 使用方式

```typescript
import { formatDateTime, formatRelativeTime, formatSmart } from "@/utils/date";

// 基本格式化
formatDateTime("2024-11-30T10:30:00Z"); // "2024/11/30 18:30"

// 不含时间
formatDate("2024-11-30T10:30:00Z"); // "2024/11/30"

// 相对时间
formatRelativeTime("2024-11-30T10:25:00Z"); // "5 分钟前"

// 智能格式化
formatSmart(new Date()); // "18:30"（今天）
formatSmart(yesterday); // "昨天 18:30"
formatSmart(lastWeek); // "11/23 18:30"
formatSmart(lastYear); // "2023/11/30"
```

## 配置选项

```typescript
interface DateFormatOptions {
  showTime?: boolean; // 是否显示时间
  showSeconds?: boolean; // 是否显示秒
  fallback?: string; // 空值显示文本
  locale?: string; // 语言
}
```

## 代码位置

```
web/src/
└── utils/
    └── date.ts    # 日期格式化工具
```
