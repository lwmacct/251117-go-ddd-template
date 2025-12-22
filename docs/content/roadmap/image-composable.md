# Image Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:30+4`
- [已实现功能](#已实现功能) `:34+14`
  - [加载管理](#加载管理) `:36+7`
  - [处理工具](#处理工具) `:43+5`
- [使用方式](#使用方式) `:48+138`
  - [基础加载](#基础加载) `:50+31`
  - [图片预加载](#图片预加载) `:81+16`
  - [懒加载](#懒加载) `:97+17`
  - [渐进式加载](#渐进式加载) `:114+20`
  - [图片验证](#图片验证) `:134+27`
  - [图片压缩](#图片压缩) `:161+25`
- [API](#api) `:186+49`
  - [useImage](#useimage) `:188+25`
  - [validateImage](#validateimage) `:213+13`
  - [useImageCompression](#useimagecompression) `:226+9`
- [代码位置](#代码位置) `:235+7`

<!--TOC-->

## 需求背景

前端需要响应式管理图片加载、预加载、懒加载、验证和压缩等功能。

## 已实现功能

### 加载管理

- `useImage` - 基础图片加载
- `useImagePreload` - 图片预加载
- `useLazyImage` - 懒加载图片
- `useProgressiveImage` - 渐进式加载

### 处理工具

- `validateImage` - 图片验证
- `useImageCompression` - 图片压缩

## 使用方式

### 基础加载

```typescript
import { useImage } from "@/composables/useImage";

const { isLoading, isReady, error, width, height, aspectRatio } = useImage("/image.jpg");

// 监听加载状态
watch(isReady, (ready) => {
  if (ready) {
    console.log(`图片尺寸: ${width.value} x ${height.value}`);
  }
});

// 带选项
const { load, abort } = useImage("/large-image.jpg", {
  immediate: false, // 不立即加载
  timeout: 10000, // 10秒超时
  fallback: "/fallback.jpg", // 失败时回退
  crossOrigin: "anonymous",
  onLoad: (img) => console.log("加载成功"),
  onError: (err) => console.log("加载失败"),
});

// 手动加载
await load();

// 取消加载
abort();
```

### 图片预加载

```typescript
import { useImagePreload } from "@/composables/useImage";

const images = ["/img1.jpg", "/img2.jpg", "/img3.jpg"];

const { preload, progress, isLoading, loaded, total, failed } = useImagePreload(images);

// 开始预加载
await preload();

console.log(`已加载 ${loaded.value}/${total.value}`);
console.log(`失败: ${failed.value}`);
```

### 懒加载

```typescript
import { useLazyImage } from "@/composables/useImage";

const { targetRef, isVisible, isReady, image } = useLazyImage("/large-image.jpg", {
  threshold: 0.1, // 10% 可见时加载
  rootMargin: "50px", // 提前 50px 开始加载
});

// 在模板中使用
// <div ref="targetRef">
//   <img v-if="isReady" :src="image?.src" />
//   <div v-else class="placeholder">加载中...</div>
// </div>
```

### 渐进式加载

```typescript
import { useProgressiveImage } from "@/composables/useImage";

const { currentSrc, stage, isReady } = useProgressiveImage({
  placeholder: "/placeholder.svg", // 占位图
  thumbnail: "/thumb.jpg", // 缩略图
  full: "/full.jpg", // 原图
});

// stage 变化: 'placeholder' -> 'thumbnail' -> 'full'
watch(stage, (s) => {
  console.log(`当前阶段: ${s}`);
});

// 在模板中使用
// <img :src="currentSrc" :class="{ blur: stage !== 'full' }" />
```

### 图片验证

```typescript
import { validateImage } from "@/composables/useImage";

const handleFileSelect = async (file: File) => {
  const result = await validateImage(file, {
    maxWidth: 1920,
    maxHeight: 1080,
    minWidth: 100,
    minHeight: 100,
    maxSize: 5 * 1024 * 1024, // 5MB
    allowedTypes: ["image/jpeg", "image/png", "image/webp"],
    aspectRatio: 16 / 9,
    aspectRatioTolerance: 0.1,
  });

  if (!result.valid) {
    console.log("验证失败:", result.errors);
    return;
  }

  console.log("图片信息:", result.info);
  // { width, height, type, size, aspectRatio }
};
```

### 图片压缩

```typescript
import { useImageCompression } from "@/composables/useImage";

const { compress, compressMultiple, isCompressing, progress } = useImageCompression({
  maxWidth: 1920,
  maxHeight: 1080,
  quality: 0.8,
  type: "image/jpeg",
});

// 压缩单张图片
const compressedBlob = await compress(file);
console.log(`压缩后: ${(compressedBlob.size / 1024).toFixed(2)}KB`);

// 压缩多张图片
const compressedBlobs = await compressMultiple(files);

// 显示进度
watch(progress, (p) => {
  console.log(`压缩进度: ${(p * 100).toFixed(0)}%`);
});
```

## API

### useImage

| 选项        | 类型                   | 默认值 | 说明             |
| ----------- | ---------------------- | ------ | ---------------- |
| immediate   | boolean                | true   | 是否立即加载     |
| delay       | number                 | 0      | 加载延迟（毫秒） |
| fallback    | string                 | -      | 回退图片         |
| placeholder | string                 | -      | 占位图片         |
| crossOrigin | string                 | -      | 跨域设置         |
| timeout     | number                 | -      | 超时时间（毫秒） |
| onLoad      | `(img) => void       ` | -      | 加载成功回调     |
| onError     | `(err) => void       ` | -      | 加载失败回调     |

| 返回值      | 类型                            | 说明         |
| ----------- | ------------------------------- | ------------ |
| image       | Ref\<HTMLImageElement \| null\> | 图片元素     |
| isLoading   | Ref\<boolean\>                  | 是否正在加载 |
| isReady     | Ref\<boolean\>                  | 是否加载完成 |
| error       | Ref\<Error \| null\>            | 加载错误     |
| width       | Ref\<number\>                   | 图片宽度     |
| height      | Ref\<number\>                   | 图片高度     |
| aspectRatio | ComputedRef\<number\>           | 宽高比       |
| load        | `() => Promise                ` | 手动加载     |
| abort       | `() => void                   ` | 取消加载     |

### validateImage

| 选项                 | 类型     | 说明                  |
| -------------------- | -------- | --------------------- |
| maxWidth             | number   | 最大宽度              |
| maxHeight            | number   | 最大高度              |
| minWidth             | number   | 最小宽度              |
| minHeight            | number   | 最小高度              |
| maxSize              | number   | 最大文件大小（字节）  |
| allowedTypes         | string[] | 允许的 MIME 类型      |
| aspectRatio          | number   | 固定宽高比            |
| aspectRatioTolerance | number   | 宽高比容差（默认0.1） |

### useImageCompression

| 选项      | 类型   | 默认值       | 说明           |
| --------- | ------ | ------------ | -------------- |
| maxWidth  | number | 1920         | 最大宽度       |
| maxHeight | number | 1080         | 最大高度       |
| quality   | number | 0.8          | 压缩质量 (0-1) |
| type      | string | 'image/jpeg' | 输出类型       |

## 代码位置

```
web/src/
└── composables/
    └── useImage.ts    # Image Composable
```
