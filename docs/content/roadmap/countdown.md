# 倒计时 Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:28+4`
- [已实现功能](#已实现功能) `:32+25`
  - [useCountdown](#usecountdown) `:34+7`
  - [useStopwatch](#usestopwatch) `:41+5`
  - [useVerificationCode](#useverificationcode) `:46+6`
  - [useTargetDateCountdown](#usetargetdatecountdown) `:52+5`
- [使用方式](#使用方式) `:57+79`
  - [基础倒计时](#基础倒计时) `:59+18`
  - [验证码倒计时](#验证码倒计时) `:77+23`
  - [秒表](#秒表) `:100+13`
  - [活动倒计时](#活动倒计时) `:113+23`
- [API](#api) `:136+17`
  - [useCountdown 返回值](#usecountdown-返回值) `:138+15`
- [代码位置](#代码位置) `:153+7`

<!--TOC-->

## 需求背景

需要为验证码发送、活动截止等场景提供倒计时功能。

## 已实现功能

### useCountdown

- 倒计时功能
- 开始/暂停/重置
- 格式化输出
- 结束回调

### useStopwatch

- 正计时（秒表）
- 开始/暂停/重置

### useVerificationCode

- 验证码发送倒计时
- 发送状态管理
- 按钮文本生成

### useTargetDateCountdown

- 目标日期倒计时
- 自动计算剩余时间

## 使用方式

### 基础倒计时

```typescript
import { useCountdown } from "@/composables/useCountdown";

const { remaining, formatted, start, pause, reset, isRunning, isFinished } = useCountdown({
  seconds: 60,
  onEnd: () => {
    toast.info("倒计时结束");
  },
});

// 开始倒计时
start();

// 显示: "01:00", "00:59", ...
```

### 验证码倒计时

```typescript
import { useVerificationCode } from "@/composables/useCountdown";

const { buttonText, isDisabled, send } = useVerificationCode({
  seconds: 60,
  onSend: async () => {
    await api.sendVerificationCode(phone);
  },
  onSendSuccess: () => toast.success("验证码已发送"),
  onSendError: (err) => toast.error(err.message),
});
```

```vue
<template>
  <v-btn :disabled="isDisabled" @click="send">
    {{ buttonText }}
  </v-btn>
</template>
```

### 秒表

```typescript
import { useStopwatch } from "@/composables/useCountdown";

const { elapsed, formatted, start, pause, reset, isRunning } = useStopwatch();

// 开始计时
start();

// 显示: "00:00:01", "00:00:02", ...
```

### 活动倒计时

```typescript
import { useTargetDateCountdown } from "@/composables/useCountdown";

const { days, hours, minutes, seconds, isFinished } = useTargetDateCountdown({
  targetDate: "2024-12-31T23:59:59",
  onEnd: () => toast.info("活动已结束"),
});
```

```vue
<template>
  <div v-if="!isFinished" class="countdown">
    <span>{{ days }}天</span>
    <span>{{ hours }}小时</span>
    <span>{{ minutes }}分</span>
    <span>{{ seconds }}秒</span>
  </div>
  <div v-else>活动已结束</div>
</template>
```

## API

### useCountdown 返回值

| 属性                       | 类型                  | 说明       |
| -------------------------- | --------------------- | ---------- |
| remaining                  | `Ref<number>`         | 剩余秒数   |
| isRunning                  | `Ref<boolean>`        | 是否运行中 |
| isFinished                 | `Ref<boolean>`        | 是否已结束 |
| formatted                  | `ComputedRef<string>` | 格式化时间 |
| days/hours/minutes/seconds | `ComputedRef<number>` | 时间部分   |
| start                      | `() => void`          | 开始       |
| pause                      | `() => void`          | 暂停       |
| stop                       | `() => void`          | 停止并重置 |
| reset                      | `(seconds?) => void`  | 重置       |
| restart                    | `() => void`          | 重新开始   |

## 代码位置

```
web/src/
└── composables/
    └── useCountdown.ts    # 倒计时
```
