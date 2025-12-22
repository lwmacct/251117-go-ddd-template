# ID 生成工具

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:34+4`
- [已实现功能](#已实现功能) `:38+46`
  - [UUID 生成](#uuid-生成) `:40+6`
  - [NanoID 风格](#nanoid-风格) `:46+9`
  - [时间戳 ID](#时间戳-id) `:55+6`
  - [前缀 ID](#前缀-id) `:61+5`
  - [序列 ID](#序列-id) `:66+5`
  - [高级 ID](#高级-id) `:71+6`
  - [实用工具](#实用工具) `:77+7`
- [使用方式](#使用方式) `:84+109`
  - [UUID](#uuid) `:86+15`
  - [NanoID 风格](#nanoid-风格-1) `:101+18`
  - [前缀 ID](#前缀-id-1) `:119+15`
  - [序列 ID](#序列-id-1) `:134+16`
  - [雪花 ID](#雪花-id) `:150+14`
  - [ULID](#ulid) `:164+13`
  - [ID 工厂](#id-工厂) `:177+16`
- [API](#api) `:193+15`
  - [主要函数](#主要函数) `:195+13`
- [代码位置](#代码位置) `:208+7`

<!--TOC-->

## 需求背景

项目中需要生成各种类型的唯一标识符，用于用户 ID、订单号、会话令牌等场景。

## 已实现功能

### UUID 生成

- `uuid` - 标准 UUID v4
- `shortUuid` - 无连字符 UUID
- `isValidUuid` - UUID 格式验证

### NanoID 风格

- `nanoid` - URL 安全的短 ID
- `customId` - 自定义字母表 ID
- `alphanumericId` - 字母数字 ID
- `numericId` - 纯数字 ID
- `alphabeticId` - 纯字母 ID
- `hexId` - 十六进制 ID

### 时间戳 ID

- `timestampId` - 时间戳 + 随机
- `sortableId` - 可排序 ID
- `extractTimestamp` - 提取时间戳

### 前缀 ID

- `prefixedId` - 带前缀的 ID
- `createPrefixedIdGenerator` - 创建前缀 ID 生成器

### 序列 ID

- `createSequence` - 序列生成器
- `createFormattedSequence` - 格式化序列

### 高级 ID

- `createSnowflake` - 雪花 ID 生成器
- `ulid` - ULID 生成
- `extractUlidTimestamp` - 提取 ULID 时间戳

### 实用工具

- `uniqueDomId` - DOM 元素 ID
- `createIdFactory` - ID 工厂
- `generateIds` - 批量生成
- `ensureUniqueId` - 确保唯一

## 使用方式

### UUID

```typescript
import { uuid, shortUuid, isValidUuid } from "@/utils/id";

// 标准 UUID
uuid(); // '550e8400-e29b-41d4-a716-446655440000'

// 短 UUID
shortUuid(); // '550e8400e29b41d4a716446655440000'

// 验证
isValidUuid("550e8400-e29b-41d4-a716-446655440000"); // true
```

### NanoID 风格

```typescript
import { nanoid, alphanumericId, numericId } from "@/utils/id";

// 默认 21 字符
nanoid(); // 'V1StGXR8_Z5jdHi6B-myT'

// 指定长度
nanoid(10); // 'IRFa-VaY2b'

// 字母数字
alphanumericId(8); // 'a1B2c3D4'

// 纯数字（验证码等）
numericId(6); // '123456'
```

### 前缀 ID

```typescript
import { prefixedId, createPrefixedIdGenerator } from "@/utils/id";

// 单次生成
prefixedId("user"); // 'user_a1b2c3d4e5f6'
prefixedId("order", 16); // 'order_a1b2c3d4e5f67890'

// 创建生成器
const userIdGen = createPrefixedIdGenerator("user");
userIdGen(); // 'user_a1b2c3d4e5f6'
userIdGen(); // 'user_x7y8z9w0a1b2'
```

### 序列 ID

```typescript
import { createSequence, createFormattedSequence } from "@/utils/id";

// 简单序列
const seq = createSequence();
seq(); // 1
seq(); // 2

// 格式化序列
const orderSeq = createFormattedSequence("ORD", 6);
orderSeq(); // 'ORD000001'
orderSeq(); // 'ORD000002'
```

### 雪花 ID

```typescript
import { createSnowflake } from "@/utils/id";

const snowflake = createSnowflake({
  machineId: 1,
  datacenterId: 1,
});

snowflake(); // '7159558526853120001'
snowflake(); // '7159558526853120002'
```

### ULID

```typescript
import { ulid, extractUlidTimestamp } from "@/utils/id";

// 生成 ULID
const id = ulid(); // '01ARZ3NDEKTSV4RRFFQ69G5FAV'

// 提取时间戳
const timestamp = extractUlidTimestamp(id);
const date = new Date(timestamp);
```

### ID 工厂

```typescript
import { createIdFactory, prefixedId, nanoid } from "@/utils/id";

const ids = createIdFactory({
  user: () => prefixedId("usr"),
  order: () => prefixedId("ord"),
  session: () => nanoid(),
});

ids.user(); // 'usr_a1b2c3d4e5f6'
ids.order(); // 'ord_x7y8z9w0a1b2'
ids.session(); // 'V1StGXR8_Z5jdHi6B-myT'
```

## API

### 主要函数

| 函数                      | 说明           |
| ------------------------- | -------------- |
| uuid()                    | UUID v4        |
| nanoid(size?)             | NanoID 风格 ID |
| prefixedId(prefix, size?) | 带前缀 ID      |
| createSequence(start?)    | 序列生成器     |
| createSnowflake(config?)  | 雪花 ID 生成器 |
| ulid()                    | ULID           |
| generateIds(count, gen?)  | 批量生成       |
| ensureUniqueId(set, gen?) | 确保唯一       |

## 代码位置

```
web/src/
└── utils/
    └── id.ts    # ID 生成工具
```
