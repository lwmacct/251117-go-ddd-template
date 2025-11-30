# 窗口尺寸检测 Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

## 需求背景

需要响应式获取窗口尺寸和断点信息，以实现响应式布局和条件渲染。

## 已实现功能

### useWindowSize

- 响应式窗口宽高
- 当前断点判断（xs/sm/md/lg/xl/xxl）
- 设备类型判断（移动端/平板/桌面）
- 横屏/竖屏检测
- 设备像素比

### useMediaQuery

- 自定义媒体查询
- 响应式匹配结果

### 预设媒体查询

- `usePrefersDark` - 深色模式偏好
- `usePrefersReducedMotion` - 减少动画偏好
- `useIsRetina` - Retina 屏幕检测
- `useIsTouchDevice` - 触摸设备检测
- `useCanHover` - 悬停支持检测

### useElementSize

- 响应式元素尺寸（使用 ResizeObserver）

## 使用方式

### 窗口尺寸

```typescript
import { useWindowSize } from "@/composables/useWindowSize";

const { width, height, breakpoint, isMobile, isDesktop } = useWindowSize();

// 在模板中使用
// <div v-if="isMobile">移动端视图</div>
// <div v-else>桌面端视图</div>
```

### 媒体查询

```typescript
import { useMediaQuery, usePrefersDark } from "@/composables/useWindowSize";

// 自定义媒体查询
const isWideScreen = useMediaQuery("(min-width: 1600px)");

// 深色模式偏好
const prefersDark = usePrefersDark();
```

### 元素尺寸

```typescript
import { ref } from "vue";
import { useElementSize } from "@/composables/useWindowSize";

const containerRef = ref<HTMLElement>();
const { width, height } = useElementSize(() => containerRef.value);
```

## API

### useWindowSize 返回值

| 属性 | 类型 | 说明 |
|------|------|------|
| width | `Ref<number>` | 窗口宽度 |
| height | `Ref<number>` | 窗口高度 |
| breakpoint | `ComputedRef<Breakpoint>` | 当前断点 |
| isMobile | `ComputedRef<boolean>` | 是否移动端 (< md) |
| isTablet | `ComputedRef<boolean>` | 是否平板 (md) |
| isDesktop | `ComputedRef<boolean>` | 是否桌面端 (>= lg) |
| isLandscape | `ComputedRef<boolean>` | 是否横屏 |
| isPortrait | `ComputedRef<boolean>` | 是否竖屏 |

## 代码位置

```
web/src/
└── composables/
    └── useWindowSize.ts    # 窗口尺寸检测
```
