# Media Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:31+4`
- [已实现功能](#已实现功能) `:35+18`
  - [播放控制](#播放控制) `:37+5`
  - [高级功能](#高级功能) `:42+7`
  - [工具](#工具) `:49+4`
- [使用方式](#使用方式) `:53+244`
  - [通用媒体控制](#通用媒体控制) `:55+61`
  - [音频播放器](#音频播放器) `:116+18`
  - [音频可视化](#音频可视化) `:134+39`
  - [录音机](#录音机) `:173+43`
  - [屏幕共享](#屏幕共享) `:216+40`
  - [画中画](#画中画) `:256+30`
  - [时长格式化](#时长格式化) `:286+11`
- [API](#api) `:297+44`
  - [useMedia](#usemedia) `:299+24`
  - [useRecorder](#userecorder) `:323+18`
- [代码位置](#代码位置) `:341+7`

<!--TOC-->

## 需求背景

前端需要媒体（音频/视频）相关的工具函数，支持播放控制、可视化、录制等高级功能。

## 已实现功能

### 播放控制

- `useMedia` - 通用媒体控制器
- `useAudio` - 音频播放器

### 高级功能

- `useAudioVisualizer` - 音频可视化
- `useRecorder` - 音频/视频录制
- `useScreenShare` - 屏幕共享
- `usePictureInPicture` - 画中画模式

### 工具

- `formatDuration` - 时长格式化

## 使用方式

### 通用媒体控制

```typescript
import { useMedia } from "@/composables/useMedia";

const { mediaRef, state, isPlaying, isPaused, currentTime, duration, progress, volume, isMuted, playbackRate, play, pause, stop, toggle, toggleMute, seek, seekPercent, forward, backward, setSource } = useMedia({
  autoplay: false,
  loop: false,
  muted: false,
  volume: 0.8,
  playbackRate: 1,
  preload: "metadata",
});

// 播放控制
await play();
pause();
await toggle();
stop();

// 跳转
seek(30); // 跳到 30 秒
seekPercent(0.5); // 跳到 50%

// 快进/快退
forward(10); // 快进 10 秒
backward(10); // 快退 10 秒

// 音量控制
volume.value = 0.5;
toggleMute();

// 播放速度
playbackRate.value = 1.5;

// 更换源
setSource("/new-video.mp4");
```

```vue
<template>
  <video ref="mediaRef" :src="videoUrl" />

  <div class="controls">
    <button @click="toggle">
      {{ isPlaying ? "暂停" : "播放" }}
    </button>

    <input type="range" :value="progress * 100" @input="seekPercent($event.target.value / 100)" />

    <span>{{ formatDuration(currentTime) }} / {{ formatDuration(duration) }}</span>

    <button @click="toggleMute">
      {{ isMuted ? "取消静音" : "静音" }}
    </button>

    <input type="range" :value="volume * 100" @input="volume = $event.target.value / 100" />
  </div>
</template>
```

### 音频播放器

```typescript
import { useAudio } from "@/composables/useMedia";

// 创建音频播放器
const { play, pause, isPlaying, currentTime, duration, audio } = useAudio("/audio.mp3", {
  autoplay: false,
  volume: 0.8,
});

// 播放
await play();

// 暂停
pause();
```

### 音频可视化

```typescript
import { useAudioVisualizer } from "@/composables/useMedia";

const { connect, disconnect, frequencyData, timeDomainData, isActive } = useAudioVisualizer({
  fftSize: 256,
  smoothingTimeConstant: 0.8,
});

// 连接到音频元素
const audioElement = document.querySelector("audio");
connect(audioElement);

// 使用频率数据绘制可视化
function draw() {
  if (!isActive.value) return;

  const canvas = canvasRef.value;
  const ctx = canvas.getContext("2d");

  // 清除画布
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  // 绘制频谱
  const barWidth = canvas.width / frequencyData.value.length;
  frequencyData.value.forEach((value, index) => {
    const barHeight = (value / 255) * canvas.height;
    ctx.fillStyle = `hsl(${(index / frequencyData.value.length) * 360}, 100%, 50%)`;
    ctx.fillRect(index * barWidth, canvas.height - barHeight, barWidth - 1, barHeight);
  });

  requestAnimationFrame(draw);
}

// 断开连接
disconnect();
```

### 录音机

```typescript
import { useRecorder } from "@/composables/useMedia";

const { isRecording, isPaused, duration, data, dataUrl, start, pause, resume, stop, clear, error } = useRecorder({
  mimeType: "audio/webm",
  audioBitsPerSecond: 128000,
});

// 开始录音
async function startRecording() {
  try {
    await start();
    console.log("录音开始");
  } catch (e) {
    console.error("无法访问麦克风:", e);
  }
}

// 暂停/恢复
function togglePause() {
  if (isPaused.value) {
    resume();
  } else {
    pause();
  }
}

// 停止并获取数据
async function stopRecording() {
  const blob = await stop();
  console.log("录音完成:", blob.size, "bytes");

  // 播放录音
  const audio = new Audio(dataUrl.value);
  audio.play();
}

// 清除数据
clear();
```

### 屏幕共享

```typescript
import { useScreenShare } from "@/composables/useMedia";

const { stream, isSharing, error, start, stop } = useScreenShare();

// 开始共享
async function startSharing() {
  try {
    const mediaStream = await start({
      video: { cursor: "always" },
      audio: true,
    });

    // 显示到视频元素
    videoElement.srcObject = mediaStream;
  } catch (e) {
    if (e.name === "NotAllowedError") {
      console.log("用户取消了屏幕共享");
    }
  }
}

// 停止共享
stop();
```

```vue
<template>
  <div>
    <video ref="videoRef" autoplay />

    <button @click="isSharing ? stop() : startSharing()">
      {{ isSharing ? "停止共享" : "开始共享" }}
    </button>
  </div>
</template>
```

### 画中画

```typescript
import { usePictureInPicture } from "@/composables/useMedia";

const videoRef = ref<HTMLVideoElement | null>(null);
const { isSupported, isPip, enter, exit, toggle } = usePictureInPicture(videoRef);

// 进入画中画
if (isSupported.value) {
  await enter();
}

// 退出画中画
await exit();

// 切换
await toggle();
```

```vue
<template>
  <video ref="videoRef" :src="videoUrl" />

  <button v-if="isSupported" @click="toggle">
    {{ isPip ? "退出画中画" : "进入画中画" }}
  </button>
</template>
```

### 时长格式化

```typescript
import { formatDuration } from "@/composables/useMedia";

formatDuration(65); // '1:05'
formatDuration(125); // '2:05'
formatDuration(3725); // '1:02:05'
formatDuration(7325); // '2:02:05'
```

## API

### useMedia

| 选项         | 类型    | 默认值     | 说明       |
| ------------ | ------- | ---------- | ---------- |
| autoplay     | boolean | false      | 自动播放   |
| loop         | boolean | false      | 循环播放   |
| muted        | boolean | false      | 静音       |
| volume       | number  | 1          | 音量 (0-1) |
| playbackRate | number  | 1          | 播放速度   |
| preload      | string  | 'metadata' | 预加载方式 |

| 返回值      | 类型             | 说明           |
| ----------- | ---------------- | -------------- |
| mediaRef    | Ref              | 媒体元素引用   |
| state       | Ref              | 当前状态       |
| isPlaying   | ComputedRef      | 是否播放中     |
| currentTime | Ref              | 当前时间       |
| duration    | Ref              | 总时长         |
| progress    | ComputedRef      | 播放进度 (0-1) |
| volume      | Ref              | 音量           |
| play        | `() => Promise`  | 播放           |
| pause       | `() => void`     | 暂停           |
| seek        | `(time) => void` | 跳转           |

### useRecorder

| 选项               | 类型   | 默认值       | 说明       |
| ------------------ | ------ | ------------ | ---------- |
| mimeType           | string | 'audio/webm' | MIME 类型  |
| audioBitsPerSecond | number | 128000       | 音频比特率 |
| videoBitsPerSecond | number | -            | 视频比特率 |

| 返回值      | 类型            | 说明            |
| ----------- | --------------- | --------------- |
| isRecording | Ref             | 是否录制中      |
| isPaused    | Ref             | 是否暂停        |
| duration    | Ref             | 录制时长        |
| data        | Ref             | 录制数据 (Blob) |
| dataUrl     | Ref             | 数据 URL        |
| start       | `() => Promise` | 开始录制        |
| stop        | `() => Promise` | 停止录制        |

## 代码位置

```
web/src/
└── composables/
    └── useMedia.ts    # Media Composable
```
