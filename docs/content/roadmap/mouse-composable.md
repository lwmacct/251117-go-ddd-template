# Mouse Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:33+4`
- [已实现功能](#已实现功能) `:37+20`
  - [位置跟踪](#位置跟踪) `:39+5`
  - [状态跟踪](#状态跟踪) `:44+5`
  - [光标控制](#光标控制) `:49+4`
  - [拖放](#拖放) `:53+4`
- [使用方式](#使用方式) `:57+100`
  - [鼠标位置](#鼠标位置) `:59+16`
  - [元素内位置](#元素内位置) `:75+13`
  - [鼠标按下状态](#鼠标按下状态) `:88+15`
  - [悬停状态](#悬停状态) `:103+15`
  - [光标控制](#光标控制-1) `:118+17`
  - [拖放区域](#拖放区域) `:135+22`
- [API](#api) `:157+47`
  - [useMouse](#usemouse) `:159+16`
  - [useMouseInElement](#usemouseinelement) `:175+12`
  - [useHover](#usehover) `:187+7`
  - [useDropZone](#usedropzone) `:194+10`
- [代码位置](#代码位置) `:204+7`

<!--TOC-->

## 需求背景

前端需要跟踪鼠标位置、状态和交互，用于自定义光标、拖放、悬停效果等场景。

## 已实现功能

### 位置跟踪

- `useMouse` - 跟踪全局鼠标位置
- `useMouseInElement` - 跟踪元素内相对位置

### 状态跟踪

- `useMousePressed` - 跟踪鼠标按下状态
- `useHover` - 跟踪元素悬停状态

### 光标控制

- `useCursor` - 控制光标样式

### 拖放

- `useDropZone` - 创建拖放区域

## 使用方式

### 鼠标位置

```typescript
import { useMouse } from "@/composables/useMouse";

// 跟踪全局鼠标位置
const { x, y, position, sourceType } = useMouse();

// 使用不同坐标类型
const { x: clientX, y: clientY } = useMouse({ type: "client" });

// 指定目标元素
const target = ref<HTMLElement | null>(null);
const { x, y } = useMouse({ target });
```

### 元素内位置

```typescript
import { useMouseInElement } from "@/composables/useMouse";

const target = ref<HTMLElement | null>(null);
const { x, y, isOutside, elementWidth, elementHeight } = useMouseInElement(target);

// 计算百分比位置
const percentX = computed(() => (x.value / elementWidth.value) * 100);
const percentY = computed(() => (y.value / elementHeight.value) * 100);
```

### 鼠标按下状态

```typescript
import { useMousePressed } from "@/composables/useMouse";

const { pressed, button } = useMousePressed();

// button: 0=左键, 1=中键, 2=右键
watch(pressed, (isPressed) => {
  if (isPressed && button.value === 0) {
    console.log("左键按下");
  }
});
```

### 悬停状态

```typescript
import { useHover } from "@/composables/useMouse";

const target = ref<HTMLElement | null>(null);
const { isHovered } = useHover(target);

// 带延迟的悬停
const { isHovered: delayedHover } = useHover(target, {
  delayEnter: 300, // 进入延迟 300ms
  delayLeave: 100, // 离开延迟 100ms
});
```

### 光标控制

```typescript
import { useCursor } from "@/composables/useMouse";

const { cursor, setCursor, resetCursor } = useCursor();

// 拖动时改变光标
function onDragStart() {
  setCursor("grabbing");
}

function onDragEnd() {
  resetCursor();
}
```

### 拖放区域

```typescript
import { useDropZone } from "@/composables/useMouse";

const target = ref<HTMLElement | null>(null);

const { isOverDropZone, files } = useDropZone(target, {
  accept: [".jpg", ".png", "image/*"],
  onDrop: (files) => {
    console.log("Dropped files:", files);
    uploadFiles(files);
  },
});

// 视觉反馈
const dropZoneClass = computed(() => ({
  "drop-zone": true,
  "drop-zone--active": isOverDropZone.value,
}));
```

## API

### useMouse

| 选项         | 类型                       | 说明         |
| ------------ | -------------------------- | ------------ |
| target       | Window \| HTMLElement      | 目标元素     |
| type         | 'page'\|'client'\|'screen' | 坐标类型     |
| touch        | boolean                    | 是否支持触摸 |
| initialValue | MousePosition              | 初始位置     |

| 返回值     | 类型                    | 说明       |
| ---------- | ----------------------- | ---------- |
| x          | `Ref<number>          ` | X 坐标     |
| y          | `Ref<number>          ` | Y 坐标     |
| position   | `Ref<MousePosition>   ` | 位置对象   |
| sourceType | Ref<'mouse'\|'touch'>   | 输入源类型 |

### useMouseInElement

| 返回值        | 说明              |
| ------------- | ----------------- |
| x             | 相对元素的 X 坐标 |
| y             | 相对元素的 Y 坐标 |
| isOutside     | 是否在元素外      |
| elementX      | 元素左边距        |
| elementY      | 元素上边距        |
| elementWidth  | 元素宽度          |
| elementHeight | 元素高度          |

### useHover

| 选项       | 说明             |
| ---------- | ---------------- |
| delayEnter | 进入延迟（毫秒） |
| delayLeave | 离开延迟（毫秒） |

### useDropZone

| 选项           | 说明             |
| -------------- | ---------------- |
| accept         | 接受的文件类型   |
| preventDefault | 是否阻止默认行为 |
| onDragEnter    | 拖入回调         |
| onDragLeave    | 拖出回调         |
| onDrop         | 拖放回调         |

## 代码位置

```
web/src/
└── composables/
    └── useMouse.ts    # Mouse Composable
```
