# Permission Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:32+4`
- [已实现功能](#已实现功能) `:36+17`
  - [通用权限](#通用权限) `:38+6`
  - [专用权限](#专用权限) `:44+9`
- [使用方式](#使用方式) `:53+106`
  - [通用权限查询](#通用权限查询) `:55+17`
  - [通知权限](#通知权限) `:72+21`
  - [摄像头权限](#摄像头权限) `:93+21`
  - [麦克风权限](#麦克风权限) `:114+17`
  - [屏幕唤醒锁](#屏幕唤醒锁) `:131+19`
  - [批量查询权限](#批量查询权限) `:150+9`
- [API](#api) `:159+44`
  - [usePermission](#usepermission) `:161+11`
  - [useNotificationPermission](#usenotificationpermission) `:172+11`
  - [useCameraPermission / useMicrophonePermission](#usecamerapermission-usemicrophonepermission) `:183+11`
  - [useScreenWakeLock](#usescreenwakelock) `:194+9`
- [支持的权限类型](#支持的权限类型) `:203+12`
- [代码位置](#代码位置) `:215+7`

<!--TOC-->

## 需求背景

前端需要管理各种浏览器权限（通知、摄像头、麦克风等）的请求和状态查询。

## 已实现功能

### 通用权限

- `usePermission` - 通用权限查询
- `queryPermissions` - 批量查询权限
- `isPermissionsApiSupported` - 检查 API 支持

### 专用权限

- `useNotificationPermission` - 通知权限
- `useClipboardPermission` - 剪贴板权限
- `useCameraPermission` - 摄像头权限
- `useMicrophonePermission` - 麦克风权限
- `useGeolocationPermission` - 地理位置权限
- `useScreenWakeLock` - 屏幕唤醒锁

## 使用方式

### 通用权限查询

```typescript
import { usePermission } from "@/composables/usePermission";

const { state, isGranted, isDenied, isPrompt, query } = usePermission({
  name: "notifications",
});

// 检查权限状态
if (isGranted.value) {
  sendNotification();
} else if (isDenied.value) {
  showPermissionDeniedMessage();
}
```

### 通知权限

```typescript
import { useNotificationPermission } from "@/composables/usePermission";

const { isGranted, request, notify } = useNotificationPermission();

async function sendNotification() {
  if (!isGranted.value) {
    await request();
  }

  if (isGranted.value) {
    notify("新消息", {
      body: "您有一条新消息",
      icon: "/icon.png",
    });
  }
}
```

### 摄像头权限

```typescript
import { useCameraPermission } from "@/composables/usePermission";

const { isGranted, request, stream, stop } = useCameraPermission();

async function startCamera() {
  const mediaStream = await request();

  if (mediaStream) {
    videoElement.value.srcObject = mediaStream;
  }
}

function stopCamera() {
  stop();
  videoElement.value.srcObject = null;
}
```

### 麦克风权限

```typescript
import { useMicrophonePermission } from "@/composables/usePermission";

const { isGranted, request, stream } = useMicrophonePermission();

async function startRecording() {
  const audioStream = await request();

  if (audioStream) {
    const mediaRecorder = new MediaRecorder(audioStream);
    mediaRecorder.start();
  }
}
```

### 屏幕唤醒锁

```typescript
import { useScreenWakeLock } from "@/composables/usePermission";

const { isActive, isSupported, request, release } = useScreenWakeLock();

// 视频播放时保持屏幕常亮
async function onVideoPlay() {
  if (isSupported.value) {
    await request();
  }
}

function onVideoPause() {
  release();
}
```

### 批量查询权限

```typescript
import { queryPermissions } from "@/composables/usePermission";

const permissions = await queryPermissions(["camera", "microphone", "notifications"]);
// { camera: 'granted', microphone: 'prompt', notifications: 'denied' }
```

## API

### usePermission

| 返回值      | 类型                              | 说明             |
| ----------- | --------------------------------- | ---------------- |
| state       | `Ref<PermissionState>           ` | 权限状态         |
| isGranted   | `Ref<boolean>                   ` | 是否已授权       |
| isDenied    | `Ref<boolean>                   ` | 是否被拒绝       |
| isPrompt    | `Ref<boolean>                   ` | 是否等待用户选择 |
| isSupported | `Ref<boolean>                   ` | 是否支持此权限   |
| query       | `() => Promise<PermissionState> ` | 刷新权限状态     |

### useNotificationPermission

| 返回值      | 说明         |
| ----------- | ------------ |
| permission  | 权限状态     |
| isGranted   | 是否已授权   |
| isDenied    | 是否被拒绝   |
| isSupported | 是否支持通知 |
| request     | 请求权限     |
| notify      | 发送通知     |

### useCameraPermission / useMicrophonePermission

| 返回值      | 说明             |
| ----------- | ---------------- |
| state       | 权限状态         |
| isGranted   | 是否已授权       |
| isSupported | 是否支持         |
| request     | 请求权限并获取流 |
| stop        | 停止所有轨道     |
| stream      | 当前媒体流       |

### useScreenWakeLock

| 返回值      | 说明       |
| ----------- | ---------- |
| isActive    | 是否已锁定 |
| isSupported | 是否支持   |
| request     | 请求锁定   |
| release     | 释放锁定   |

## 支持的权限类型

| 权限名称         | 说明       |
| ---------------- | ---------- |
| geolocation      | 地理位置   |
| notifications    | 通知       |
| camera           | 摄像头     |
| microphone       | 麦克风     |
| clipboard-read   | 剪贴板读取 |
| clipboard-write  | 剪贴板写入 |
| screen-wake-lock | 屏幕唤醒锁 |

## 代码位置

```
web/src/
└── composables/
    └── usePermission.ts    # Permission Composable
```
