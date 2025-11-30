/**
 * 窗口尺寸检测 Composable
 * 响应式获取窗口尺寸和断点信息
 */

import { ref, computed, onMounted, onUnmounted } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface WindowSize {
  width: number;
  height: number;
}

export interface UseWindowSizeOptions {
  /** 初始宽度 (SSR) */
  initialWidth?: number;
  /** 初始高度 (SSR) */
  initialHeight?: number;
  /** 监听 resize 事件，默认 true */
  listenResize?: boolean;
  /** 监听 orientationchange 事件，默认 true */
  listenOrientation?: boolean;
}

/** Vuetify 断点名称 */
export type Breakpoint = "xs" | "sm" | "md" | "lg" | "xl" | "xxl";

/** 断点阈值配置 (Vuetify 3 默认值) */
export const BREAKPOINTS = {
  xs: 0,
  sm: 600,
  md: 960,
  lg: 1280,
  xl: 1920,
  xxl: 2560,
} as const;

// ============================================================================
// 主函数
// ============================================================================

/**
 * 响应式窗口尺寸
 * @example
 * const { width, height, breakpoint, isMobile, isDesktop } = useWindowSize()
 */
export function useWindowSize(options: UseWindowSizeOptions = {}) {
  const {
    initialWidth = typeof window !== "undefined" ? window.innerWidth : 0,
    initialHeight = typeof window !== "undefined" ? window.innerHeight : 0,
    listenResize = true,
    listenOrientation = true,
  } = options;

  // 响应式尺寸
  const width = ref(initialWidth);
  const height = ref(initialHeight);

  // 更新尺寸
  const updateSize = () => {
    if (typeof window !== "undefined") {
      width.value = window.innerWidth;
      height.value = window.innerHeight;
    }
  };

  // 当前断点
  const breakpoint = computed<Breakpoint>(() => {
    const w = width.value;
    if (w >= BREAKPOINTS.xxl) return "xxl";
    if (w >= BREAKPOINTS.xl) return "xl";
    if (w >= BREAKPOINTS.lg) return "lg";
    if (w >= BREAKPOINTS.md) return "md";
    if (w >= BREAKPOINTS.sm) return "sm";
    return "xs";
  });

  // 断点判断 - 移动端 (xs, sm)
  const isMobile = computed(() => width.value < BREAKPOINTS.md);

  // 断点判断 - 平板 (md)
  const isTablet = computed(
    () => width.value >= BREAKPOINTS.md && width.value < BREAKPOINTS.lg
  );

  // 断点判断 - 桌面端 (lg+)
  const isDesktop = computed(() => width.value >= BREAKPOINTS.lg);

  // 断点判断 - 小屏 (xs)
  const isXs = computed(() => width.value < BREAKPOINTS.sm);

  // 断点判断 - sm 及以上
  const isSmAndUp = computed(() => width.value >= BREAKPOINTS.sm);

  // 断点判断 - md 及以上
  const isMdAndUp = computed(() => width.value >= BREAKPOINTS.md);

  // 断点判断 - lg 及以上
  const isLgAndUp = computed(() => width.value >= BREAKPOINTS.lg);

  // 断点判断 - xl 及以上
  const isXlAndUp = computed(() => width.value >= BREAKPOINTS.xl);

  // 横屏/竖屏
  const isLandscape = computed(() => width.value > height.value);
  const isPortrait = computed(() => height.value >= width.value);

  // 设备像素比
  const pixelRatio = ref(
    typeof window !== "undefined" ? window.devicePixelRatio : 1
  );

  // 生命周期
  onMounted(() => {
    updateSize();

    if (listenResize) {
      window.addEventListener("resize", updateSize);
    }

    if (listenOrientation) {
      window.addEventListener("orientationchange", updateSize);
    }
  });

  onUnmounted(() => {
    if (listenResize) {
      window.removeEventListener("resize", updateSize);
    }

    if (listenOrientation) {
      window.removeEventListener("orientationchange", updateSize);
    }
  });

  return {
    // 尺寸
    width,
    height,
    pixelRatio,

    // 断点
    breakpoint,

    // 设备类型判断
    isMobile,
    isTablet,
    isDesktop,

    // 断点判断
    isXs,
    isSmAndUp,
    isMdAndUp,
    isLgAndUp,
    isXlAndUp,

    // 方向判断
    isLandscape,
    isPortrait,

    // 方法
    updateSize,
  };
}

// ============================================================================
// 媒体查询 Composable
// ============================================================================

export interface UseMediaQueryOptions {
  /** 初始值 (SSR) */
  initialValue?: boolean;
}

/**
 * 响应式媒体查询
 * @example
 * const isDark = useMediaQuery("(prefers-color-scheme: dark)")
 * const isRetina = useMediaQuery("(min-resolution: 2dppx)")
 */
export function useMediaQuery(
  query: string,
  options: UseMediaQueryOptions = {}
) {
  const { initialValue = false } = options;

  const matches = ref(initialValue);

  let mediaQuery: MediaQueryList | null = null;

  const updateMatches = (e: MediaQueryListEvent | MediaQueryList) => {
    matches.value = e.matches;
  };

  onMounted(() => {
    if (typeof window !== "undefined" && "matchMedia" in window) {
      mediaQuery = window.matchMedia(query);
      matches.value = mediaQuery.matches;

      // 使用新的 addEventListener API
      if (mediaQuery.addEventListener) {
        mediaQuery.addEventListener("change", updateMatches);
      } else {
        // 兼容旧版浏览器
        mediaQuery.addListener(updateMatches);
      }
    }
  });

  onUnmounted(() => {
    if (mediaQuery) {
      if (mediaQuery.removeEventListener) {
        mediaQuery.removeEventListener("change", updateMatches);
      } else {
        mediaQuery.removeListener(updateMatches);
      }
    }
  });

  return matches;
}

// ============================================================================
// 预设媒体查询
// ============================================================================

/**
 * 检测是否偏好深色模式
 */
export function usePrefersDark() {
  return useMediaQuery("(prefers-color-scheme: dark)");
}

/**
 * 检测是否偏好减少动画
 */
export function usePrefersReducedMotion() {
  return useMediaQuery("(prefers-reduced-motion: reduce)");
}

/**
 * 检测是否为高对比度模式
 */
export function usePrefersContrast() {
  return useMediaQuery("(prefers-contrast: more)");
}

/**
 * 检测是否为 Retina 屏幕
 */
export function useIsRetina() {
  return useMediaQuery("(min-resolution: 2dppx)");
}

/**
 * 检测是否支持触摸
 */
export function useIsTouchDevice() {
  return useMediaQuery("(pointer: coarse)");
}

/**
 * 检测是否支持悬停
 */
export function useCanHover() {
  return useMediaQuery("(hover: hover)");
}

// ============================================================================
// 元素尺寸 Composable
// ============================================================================

export interface UseElementSizeOptions {
  /** 初始宽度 */
  initialWidth?: number;
  /** 初始高度 */
  initialHeight?: number;
  /** 观察的 box 类型 */
  box?: ResizeObserverBoxOptions;
}

/**
 * 响应式元素尺寸
 * @example
 * const el = ref<HTMLElement>()
 * const { width, height } = useElementSize(el)
 */
export function useElementSize(
  target: () => HTMLElement | null | undefined,
  options: UseElementSizeOptions = {}
) {
  const { initialWidth = 0, initialHeight = 0, box = "content-box" } = options;

  const width = ref(initialWidth);
  const height = ref(initialHeight);

  let observer: ResizeObserver | null = null;

  const observe = () => {
    const element = target();
    if (!element) return;

    // 初始尺寸
    const rect = element.getBoundingClientRect();
    width.value = rect.width;
    height.value = rect.height;

    // 创建观察器
    observer = new ResizeObserver((entries) => {
      const entry = entries[0];
      if (entry) {
        const boxSize =
          box === "border-box"
            ? entry.borderBoxSize
            : box === "content-box"
              ? entry.contentBoxSize
              : entry.devicePixelContentBoxSize;

        if (boxSize) {
          const size = Array.isArray(boxSize) ? boxSize[0] : boxSize;
          width.value = size.inlineSize;
          height.value = size.blockSize;
        } else {
          // 回退到 contentRect
          width.value = entry.contentRect.width;
          height.value = entry.contentRect.height;
        }
      }
    });

    observer.observe(element, { box });
  };

  onMounted(() => {
    observe();
  });

  onUnmounted(() => {
    if (observer) {
      observer.disconnect();
      observer = null;
    }
  });

  return {
    width,
    height,
  };
}
