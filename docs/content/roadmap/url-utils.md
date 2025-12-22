# URL 工具函数

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:31+4`
- [已实现功能](#已实现功能) `:35+47`
  - [URL 解析](#url-解析) `:37+7`
  - [查询参数](#查询参数) `:44+9`
  - [URL 操作](#url-操作) `:53+10`
  - [URL 编码](#url-编码) `:63+7`
  - [URL 构建](#url-构建) `:70+5`
  - [特殊 URL](#特殊-url) `:75+7`
- [使用方式](#使用方式) `:82+83`
  - [URL 解析](#url-解析-1) `:84+22`
  - [查询参数](#查询参数-1) `:106+18`
  - [路径操作](#路径操作) `:124+18`
  - [URL 构建器](#url-构建器) `:142+9`
  - [Data URL](#data-url) `:151+14`
- [API](#api) `:165+14`
  - [主要函数](#主要函数) `:167+12`
- [代码位置](#代码位置) `:179+7`

<!--TOC-->

## 需求背景

前端需要处理 URL 解析、查询参数操作、路径拼接等场景。

## 已实现功能

### URL 解析

- `parseURL` - 解析 URL
- `isValidURL` - 检查有效性
- `isAbsoluteURL` - 检查是否为绝对 URL
- `isRelativeURL` - 检查是否为相对 URL

### 查询参数

- `parseQuery` - 解析查询字符串
- `buildQuery` - 构建查询字符串
- `getQueryParam` - 获取查询参数
- `setQueryParam` - 设置查询参数
- `removeQueryParam` - 删除查询参数
- `mergeQueryParams` - 合并查询参数

### URL 操作

- `joinURL` - 连接 URL 路径
- `normalizeURL` - 规范化 URL
- `getBasePath` - 获取基础路径
- `getFileName` - 获取文件名
- `getExtension` - 获取扩展名
- `setHash` - 设置 hash
- `removeHash` - 移除 hash

### URL 编码

- `encodeURLComponent` - 编码组件
- `decodeURLComponent` - 解码组件
- `encodeURL` - 编码 URL
- `decodeURL` - 解码 URL

### URL 构建

- `buildURL` - 构建 URL
- `createURLBuilder` - 创建链式构建器

### 特殊 URL

- `createDataURL` - 创建 Data URL
- `parseDataURL` - 解析 Data URL
- `createBlobURL` - 创建 Blob URL
- `revokeBlobURL` - 释放 Blob URL

## 使用方式

### URL 解析

```typescript
import { parseURL, isValidURL, isAbsoluteURL } from "@/utils/url";

// 解析 URL
const parsed = parseURL("https://example.com:8080/path?a=1#hash");
// {
//   protocol: 'https:',
//   hostname: 'example.com',
//   port: '8080',
//   pathname: '/path',
//   search: '?a=1',
//   hash: '#hash',
//   ...
// }

// 验证
isValidURL("https://example.com"); // true
isAbsoluteURL("/path/to/page"); // false
```

### 查询参数

```typescript
import { parseQuery, buildQuery, setQueryParam } from "@/utils/url";

// 解析
parseQuery("a=1&b=2&a=3");
// { a: ['1', '3'], b: '2' }

// 构建
buildQuery({ a: "1", b: ["2", "3"] });
// 'a=1&b=2&b=3'

// 修改 URL 参数
setQueryParam("https://example.com?a=1", "b", "2");
// 'https://example.com?a=1&b=2'
```

### 路径操作

```typescript
import { joinURL, getBasePath, getFileName } from "@/utils/url";

// 连接路径
joinURL("https://example.com", "api", "users");
// 'https://example.com/api/users'

// 获取基础路径
getBasePath("https://example.com/path/to/file.pdf");
// 'https://example.com/path/to'

// 获取文件名
getFileName("https://example.com/path/to/file.pdf");
// 'file.pdf'
```

### URL 构建器

```typescript
import { createURLBuilder } from "@/utils/url";

const url = createURLBuilder("https://example.com").setPath("/api/users").setQuery({ page: "1", limit: "10" }).setHash("top").toString();
// 'https://example.com/api/users?page=1&limit=10#top'
```

### Data URL

```typescript
import { createDataURL, parseDataURL } from "@/utils/url";

// 创建 Data URL
const dataURL = createDataURL("Hello World", "text/plain");
// 'data:text/plain;base64,SGVsbG8gV29ybGQ='

// 解析 Data URL
const parsed = parseDataURL(dataURL);
// { mimeType: 'text/plain', data: 'Hello World' }
```

## API

### 主要函数

| 函数                           | 说明            |
| ------------------------------ | --------------- |
| parseURL(url)                  | 解析 URL        |
| parseQuery(query)              | 解析查询字符串  |
| buildQuery(params)             | 构建查询字符串  |
| joinURL(...parts)              | 连接 URL 路径   |
| setQueryParam(url, key, value) | 设置查询参数    |
| createURLBuilder(base)         | 创建 URL 构建器 |
| createDataURL(data, type)      | 创建 Data URL   |

## 代码位置

```
web/src/
└── utils/
    └── url.ts    # URL 工具函数
```
