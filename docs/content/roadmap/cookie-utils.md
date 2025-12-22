# Cookie 工具函数

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:33+4`
- [已实现功能](#已实现功能) `:37+46`
  - [基础操作](#基础操作) `:39+7`
  - [批量操作](#批量操作) `:46+8`
  - [JSON Cookie](#json-cookie) `:54+5`
  - [解析工具](#解析工具) `:59+5`
  - [工具函数](#工具函数) `:64+8`
  - [Cookie 管理器](#cookie-管理器) `:72+4`
  - [预设配置](#预设配置) `:76+7`
- [使用方式](#使用方式) `:83+102`
  - [基础操作](#基础操作-1) `:85+27`
  - [JSON Cookie](#json-cookie-1) `:112+18`
  - [批量操作](#批量操作-1) `:130+22`
  - [Cookie 管理器](#cookie-管理器-1) `:152+21`
  - [使用预设配置](#使用预设配置) `:173+12`
- [API](#api) `:185+25`
  - [Cookie 选项](#cookie-选项) `:187+11`
  - [主要函数](#主要函数) `:198+12`
- [代码位置](#代码位置) `:210+7`

<!--TOC-->

## 需求背景

前端需要统一的 Cookie 操作工具，用于 Token 存储、用户偏好设置等场景。

## 已实现功能

### 基础操作

- `getCookie` - 获取 Cookie
- `setCookie` - 设置 Cookie
- `removeCookie` - 删除 Cookie
- `hasCookie` - 检查是否存在

### 批量操作

- `getAllCookies` - 获取所有 Cookie
- `getCookies` - 获取多个 Cookie
- `setCookies` - 设置多个 Cookie
- `removeCookies` - 删除多个 Cookie
- `clearAllCookies` - 清除所有 Cookie

### JSON Cookie

- `getJsonCookie` - 获取 JSON Cookie
- `setJsonCookie` - 设置 JSON Cookie

### 解析工具

- `parseCookieString` - 解析 Cookie 字符串
- `serializeCookie` - 序列化 Cookie

### 工具函数

- `getCookieCount` - 获取数量
- `getCookieNames` - 获取所有名称
- `areCookiesEnabled` - 检查是否启用
- `getCookiesSize` - 获取总大小
- `getCookiesRemainingSpace` - 获取剩余空间

### Cookie 管理器

- `createCookieManager` - 创建管理器实例

### 预设配置

- `SESSION_COOKIE` - 会话 Cookie
- `PERSISTENT_COOKIE` - 持久 Cookie（7 天）
- `SECURE_COOKIE` - 安全 Cookie
- `CROSS_SITE_COOKIE` - 跨站 Cookie

## 使用方式

### 基础操作

```typescript
import { getCookie, setCookie, removeCookie } from "@/utils/cookie";

// 获取 Cookie
const token = getCookie("token"); // 'abc123' or null

// 设置 Cookie
setCookie("token", "abc123");

// 设置带选项的 Cookie
setCookie("token", "abc123", {
  expires: 3600, // 1 小时后过期
  secure: true,
  sameSite: "Strict",
});

// 使用 Date 对象设置过期时间
setCookie("token", "abc123", {
  expires: new Date("2024-12-31"),
});

// 删除 Cookie
removeCookie("token");
```

### JSON Cookie

```typescript
import { getJsonCookie, setJsonCookie } from "@/utils/cookie";

interface User {
  id: number;
  name: string;
}

// 存储对象
setJsonCookie<User>("user", { id: 1, name: "John" });

// 读取对象
const user = getJsonCookie<User>("user");
// { id: 1, name: 'John' }
```

### 批量操作

```typescript
import { getAllCookies, setCookies, clearAllCookies } from "@/utils/cookie";

// 获取所有
const allCookies = getAllCookies();
// { token: 'abc', user: '...', theme: 'dark' }

// 批量设置
setCookies(
  {
    token: "abc",
    theme: "dark",
  },
  { expires: 86400 },
);

// 清除所有
clearAllCookies();
```

### Cookie 管理器

```typescript
import { createCookieManager } from "@/utils/cookie";

// 创建带默认选项的管理器
const cookies = createCookieManager({
  path: "/",
  secure: true,
  sameSite: "Strict",
});

// 使用管理器
cookies.set("token", "abc123");
cookies.get("token"); // 'abc123'
cookies.setJson("user", { id: 1 });
cookies.getJson<User>("user"); // { id: 1 }
cookies.remove("token");
cookies.clear();
```

### 使用预设配置

```typescript
import { setCookie, SESSION_COOKIE, SECURE_COOKIE } from "@/utils/cookie";

// 会话 Cookie（浏览器关闭后过期）
setCookie("session", "temp", SESSION_COOKIE);

// 安全 Cookie
setCookie("token", "secret", SECURE_COOKIE);
```

## API

### Cookie 选项

| 选项     | 类型                        | 说明                 |
| -------- | --------------------------- | -------------------- |
| expires  | number \| Date              | 过期时间（秒或日期） |
| maxAge   | number                      | 最大存活时间（秒）   |
| path     | string                      | 有效路径             |
| domain   | string                      | 有效域名             |
| secure   | boolean                     | 是否仅 HTTPS         |
| sameSite | 'Strict' \| 'Lax' \| 'None' | SameSite 策略        |

### 主要函数

| 函数                         | 说明             |
| ---------------------------- | ---------------- |
| getCookie(name)              | 获取 Cookie      |
| setCookie(name, value, opts) | 设置 Cookie      |
| removeCookie(name)           | 删除 Cookie      |
| getJsonCookie(name)          | 获取 JSON Cookie |
| setJsonCookie(name, value)   | 设置 JSON Cookie |
| getAllCookies()              | 获取所有 Cookie  |
| createCookieManager(opts)    | 创建管理器       |

## 代码位置

```
web/src/
└── utils/
    └── cookie.ts    # Cookie 工具函数
```
