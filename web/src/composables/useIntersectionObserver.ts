/**
 * 交叉观察器 Composable
 * 用于检测元素可见性、懒加载、无限滚动等场景
 */

import { ref, watch, onMounted, onUnmounted, type Ref } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UseIntersectionObserverOptions {
  /** 根元素，默认为视口 */
  root?: Element | Document | null;
  /** 根元素的边距 */
  rootMargin?: string;
  /** 触发回调的阈值 (0-1) */
  threshold?: number | number[];
  /** 是否立即开始观察，默认 true */
  immediate?: boolean;
}

export interface UseIntersectionObserverReturn {
  /** 是否可见 */
  isVisible: Ref<boolean>;
  /** 是否曾经可见过 */
  hasBeenVisible: Ref<boolean>;
  /** 交叉比例 (0-1) */
  intersectionRatio: Ref<number>;
  /** 开始观察 */
  observe: () => void;
  /** 停止观察 */
  unobserve: () => void;
}

// ============================================================================
// 主函数
// ============================================================================

/**
 * 交叉观察器
 * @example
 * const target = ref<HTMLElement>()
 * const { isVisible, hasBeenVisible } = useIntersectionObserver(target)
 */
export function useIntersectionObserver(
  target: Ref<Element | null | undefined>,
  callback?: (
    entries: IntersectionObserverEntry[],
    observer: IntersectionObserver
  ) => void,
  options: UseIntersectionObserverOptions = {}
): UseIntersectionObserverReturn {
  const {
    root = null,
    rootMargin = "0px",
    threshold = 0,
    immediate = true,
  } = options;

  const isVisible = ref(false);
  const hasBeenVisible = ref(false);
  const intersectionRatio = ref(0);

  let observer: IntersectionObserver | null = null;

  const observe = () => {
    const element = target.value;
    if (!element || typeof IntersectionObserver === "undefined") return;

    // 清理旧的观察器
    unobserve();

    observer = new IntersectionObserver(
      (entries, obs) => {
        const entry = entries[0];
        if (entry) {
          isVisible.value = entry.isIntersecting;
          intersectionRatio.value = entry.intersectionRatio;

          if (entry.isIntersecting) {
            hasBeenVisible.value = true;
          }
        }

        // 调用自定义回调
        callback?.(entries, obs);
      },
      {
        root,
        rootMargin,
        threshold,
      }
    );

    observer.observe(element);
  };

  const unobserve = () => {
    if (observer) {
      observer.disconnect();
      observer = null;
    }
  };

  // 监听目标元素变化
  watch(
    target,
    (newTarget) => {
      if (newTarget && immediate) {
        observe();
      } else {
        unobserve();
      }
    },
    { immediate }
  );

  onMounted(() => {
    if (immediate && target.value) {
      observe();
    }
  });

  onUnmounted(() => {
    unobserve();
  });

  return {
    isVisible,
    hasBeenVisible,
    intersectionRatio,
    observe,
    unobserve,
  };
}

// ============================================================================
// 懒加载 Composable
// ============================================================================

export interface UseLazyLoadOptions extends UseIntersectionObserverOptions {
  /** 加载后是否停止观察，默认 true */
  once?: boolean;
}

/**
 * 懒加载
 * @example
 * const imageRef = ref<HTMLImageElement>()
 * const { shouldLoad } = useLazyLoad(imageRef)
 * // <img v-if="shouldLoad" :src="imageSrc" />
 */
export function useLazyLoad(
  target: Ref<Element | null | undefined>,
  options: UseLazyLoadOptions = {}
) {
  const { once = true, ...observerOptions } = options;

  const shouldLoad = ref(false);

  const { isVisible, hasBeenVisible, observe, unobserve } =
    useIntersectionObserver(
      target,
      (entries) => {
        const entry = entries[0];
        if (entry?.isIntersecting) {
          shouldLoad.value = true;

          if (once) {
            unobserve();
          }
        }
      },
      {
        rootMargin: "50px", // 提前 50px 开始加载
        ...observerOptions,
      }
    );

  return {
    shouldLoad,
    isVisible,
    hasBeenVisible,
    observe,
    unobserve,
  };
}

// ============================================================================
// 无限滚动 Composable
// ============================================================================

export interface UseInfiniteScrollOptions extends UseIntersectionObserverOptions {
  /** 是否正在加载 */
  loading?: Ref<boolean>;
  /** 是否还有更多数据 */
  hasMore?: Ref<boolean>;
  /** 加载更多的回调 */
  onLoadMore?: () => void | Promise<void>;
}

/**
 * 无限滚动
 * @example
 * const sentinelRef = ref<HTMLElement>()
 * const loading = ref(false)
 * const hasMore = ref(true)
 *
 * useInfiniteScroll(sentinelRef, {
 *   loading,
 *   hasMore,
 *   onLoadMore: async () => {
 *     loading.value = true
 *     await fetchMoreData()
 *     loading.value = false
 *   }
 * })
 */
export function useInfiniteScroll(
  target: Ref<Element | null | undefined>,
  options: UseInfiniteScrollOptions = {}
) {
  const {
    loading = ref(false),
    hasMore = ref(true),
    onLoadMore,
    ...observerOptions
  } = options;

  const { isVisible, observe, unobserve } = useIntersectionObserver(
    target,
    async (entries) => {
      const entry = entries[0];

      if (
        entry?.isIntersecting &&
        !loading.value &&
        hasMore.value &&
        onLoadMore
      ) {
        await onLoadMore();
      }
    },
    {
      rootMargin: "100px", // 提前 100px 触发
      ...observerOptions,
    }
  );

  return {
    isVisible,
    loading,
    hasMore,
    observe,
    unobserve,
  };
}

// ============================================================================
// 元素进入视口动画 Composable
// ============================================================================

export interface UseAnimateOnScrollOptions extends UseIntersectionObserverOptions {
  /** 动画类名 */
  animationClass?: string;
  /** 是否只触发一次，默认 true */
  once?: boolean;
  /** 触发阈值，默认 0.1 */
  threshold?: number;
}

/**
 * 元素进入视口时添加动画类
 * @example
 * const el = ref<HTMLElement>()
 * const { isAnimated } = useAnimateOnScroll(el, {
 *   animationClass: 'fade-in'
 * })
 */
export function useAnimateOnScroll(
  target: Ref<Element | null | undefined>,
  options: UseAnimateOnScrollOptions = {}
) {
  const {
    animationClass = "animate-in",
    once = true,
    threshold = 0.1,
    ...observerOptions
  } = options;

  const isAnimated = ref(false);

  const { isVisible, hasBeenVisible, observe, unobserve } =
    useIntersectionObserver(
      target,
      (entries) => {
        const entry = entries[0];
        const element = target.value;

        if (entry?.isIntersecting && element) {
          isAnimated.value = true;
          element.classList.add(animationClass);

          if (once) {
            unobserve();
          }
        } else if (!once && element) {
          isAnimated.value = false;
          element.classList.remove(animationClass);
        }
      },
      {
        threshold,
        ...observerOptions,
      }
    );

  return {
    isAnimated,
    isVisible,
    hasBeenVisible,
    observe,
    unobserve,
  };
}

// ============================================================================
// 监控元素可见性 Composable
// ============================================================================

export interface UseVisibilityOptions extends UseIntersectionObserverOptions {
  /** 进入视口回调 */
  onEnter?: () => void;
  /** 离开视口回调 */
  onLeave?: () => void;
}

/**
 * 监控元素可见性
 * @example
 * const el = ref<HTMLElement>()
 * useVisibility(el, {
 *   onEnter: () => console.log('进入视口'),
 *   onLeave: () => console.log('离开视口')
 * })
 */
export function useVisibility(
  target: Ref<Element | null | undefined>,
  options: UseVisibilityOptions = {}
) {
  const { onEnter, onLeave, ...observerOptions } = options;

  let wasVisible = false;

  const { isVisible, hasBeenVisible, intersectionRatio, observe, unobserve } =
    useIntersectionObserver(
      target,
      (entries) => {
        const entry = entries[0];
        if (!entry) return;

        if (entry.isIntersecting && !wasVisible) {
          wasVisible = true;
          onEnter?.();
        } else if (!entry.isIntersecting && wasVisible) {
          wasVisible = false;
          onLeave?.();
        }
      },
      observerOptions
    );

  return {
    isVisible,
    hasBeenVisible,
    intersectionRatio,
    observe,
    unobserve,
  };
}
