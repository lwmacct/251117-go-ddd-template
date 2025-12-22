# 字符串工具函数

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:24+4`
- [已实现功能](#已实现功能) `:28+48`
  - [截断与填充](#截断与填充) `:30+5`
  - [大小写转换](#大小写转换) `:35+10`
  - [URL 与 Slug](#url-与-slug) `:45+4`
  - [搜索与匹配](#搜索与匹配) `:49+7`
  - [格式化（脱敏）](#格式化脱敏) `:56+8`
  - [其他工具](#其他工具) `:64+12`
- [使用方式](#使用方式) `:76+26`
- [代码位置](#代码位置) `:102+7`

<!--TOC-->

## 需求背景

项目中多处需要处理字符串（截断、格式化、大小写转换、脱敏等），需要统一的字符串处理工具。

## 已实现功能

### 截断与填充

- `truncate` - 截断字符串（支持首/中/尾三种位置）
- `padStart` / `padEnd` - 左/右填充

### 大小写转换

- `capitalize` - 首字母大写
- `titleCase` - 每个单词首字母大写
- `camelCase` - 驼峰命名
- `pascalCase` - 帕斯卡命名（大驼峰）
- `snakeCase` - 蛇形命名
- `kebabCase` - 短横线命名
- `constantCase` - 常量命名

### URL 与 Slug

- `slugify` - 转换为 URL 友好的 slug

### 搜索与匹配

- `containsIgnoreCase` - 不区分大小写包含检查
- `highlight` - 高亮搜索关键词
- `escapeRegExp` - 转义正则特殊字符
- `fuzzyMatch` - 模糊匹配

### 格式化（脱敏）

- `maskPhone` - 手机号脱敏
- `maskEmail` - 邮箱脱敏
- `maskIdCard` - 身份证脱敏
- `formatBankCard` - 银行卡格式化
- `maskBankCard` - 银行卡脱敏

### 其他工具

- `stripHtml` - 移除 HTML 标签
- `escapeHtml` / `unescapeHtml` - HTML 转义
- `randomString` - 生成随机字符串
- `byteLength` - 计算 UTF-8 字节长度
- `isBlank` / `isNotBlank` - 空白检查
- `template` - 字符串模板替换
- `countOccurrences` - 统计子串出现次数
- `reverse` - 反转字符串
- `normalizeSpaces` - 移除重复空格

## 使用方式

```typescript
import { truncate, camelCase, slugify, highlight, maskPhone, template } from "@/utils/string";

// 截断字符串
truncate("Hello World", { length: 8 }); // "Hello..."
truncate("Hello World", { length: 8, position: "start" }); // "...World"

// 大小写转换
camelCase("hello-world"); // "helloWorld"
snakeCase("helloWorld"); // "hello_world"

// URL slug
slugify("Hello World!"); // "hello-world"

// 高亮搜索词
highlight("Hello World", "world"); // "Hello <mark>World</mark>"

// 手机号脱敏
maskPhone("13812345678"); // "138****5678"

// 模板替换
template("Hello, {name}!", { name: "World" }); // "Hello, World!"
```

## 代码位置

```
web/src/
└── utils/
    └── string.ts    # 字符串工具函数
```
