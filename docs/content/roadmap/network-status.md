# 网络状态检测 Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:28+4`
- [已实现功能](#已实现功能) `:32+25`
  - [useNetwork](#usenetwork) `:34+8`
  - [useOnline](#useonline) `:42+5`
  - [useNetworkSpeed](#usenetworkspeed) `:47+5`
  - [useNetworkBanner](#usenetworkbanner) `:52+5`
- [使用方式](#使用方式) `:57+54`
  - [基础用法](#基础用法) `:59+15`
  - [简化用法](#简化用法) `:74+11`
  - [网络速度测试](#网络速度测试) `:85+12`
  - [网络状态提示](#网络状态提示) `:97+14`
- [API](#api) `:111+14`
  - [useNetwork 返回值](#usenetwork-返回值) `:113+12`
- [代码位置](#代码位置) `:125+7`

<!--TOC-->

## 需求背景

需要检测用户网络状态，在离线或弱网时提供适当的用户体验。

## 已实现功能

### useNetwork

- 在线/离线状态检测
- Network Information API 支持
- 有效连接类型 (2g/3g/4g)
- 下行速度、RTT 信息
- 数据节省模式检测

### useOnline

- 简化的在线/离线检测
- 状态变化回调

### useNetworkSpeed

- 实际网络速度测试
- 下载速度计算 (Mbps)

### useNetworkBanner

- 网络状态提示数据
- 适用于显示离线/弱网提示条

## 使用方式

### 基础用法

```typescript
import { useNetwork } from "@/composables/useNetwork";

const { isOnline, effectiveType, downlink, isSlowConnection, connectionStatus } = useNetwork();

// 监听网络变化
watch(isOnline, (online) => {
  if (!online) {
    toast.warning("网络已断开");
  }
});
```

### 简化用法

```typescript
import { useOnline } from "@/composables/useNetwork";

const isOnline = useOnline({
  onOffline: () => toast.warning("网络已断开"),
  onOnline: () => toast.success("网络已恢复"),
});
```

### 网络速度测试

```typescript
import { useNetworkSpeed } from "@/composables/useNetwork";

const { speed, isLoading, test } = useNetworkSpeed();

// 执行测试
await test();
console.log(`下载速度: ${speed.value} Mbps`);
```

### 网络状态提示

```typescript
import { useNetworkBanner } from "@/composables/useNetwork";

const { shouldShowBanner, bannerType, bannerMessage } = useNetworkBanner();
```

```vue
<template>
  <v-banner v-if="shouldShowBanner" :type="bannerType" :text="bannerMessage" />
</template>
```

## API

### useNetwork 返回值

| 属性             | 类型                   | 说明                       |
| ---------------- | ---------------------- | -------------------------- |
| isOnline         | `Ref<boolean>`         | 是否在线                   |
| effectiveType    | `Ref<string>`          | 有效连接类型               |
| downlink         | `Ref<number>`          | 下行速度 (Mbps)            |
| rtt              | `Ref<number>`          | 往返时间 (ms)              |
| saveData         | `Ref<boolean>`         | 数据节省模式               |
| isSlowConnection | `ComputedRef<boolean>` | 是否慢速连接               |
| connectionStatus | `ComputedRef<string>`  | good/moderate/poor/offline |

## 代码位置

```
web/src/
└── composables/
    └── useNetwork.ts    # 网络状态检测
```
