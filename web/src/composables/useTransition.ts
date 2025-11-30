/**
 * Transition Composable
 * 提供过渡和动画相关的工具函数
 */

import {
  ref,
  computed,
  watch,
  nextTick,
  onMounted,
  onUnmounted,
  type Ref,
  type ComputedRef,
} from "vue";

// ============================================================================
// 类型定义
// ============================================================================

/**
 * 过渡状态
 */
export type TransitionState = "idle" | "enter" | "leave";

/**
 * 动画缓动函数类型
 */
export type EasingFunction =
  | "linear"
  | "ease"
  | "ease-in"
  | "ease-out"
  | "ease-in-out"
  | "cubic-bezier"
  | string;

/**
 * 过渡配置
 */
export interface TransitionConfig {
  /** 过渡持续时间（毫秒） */
  duration?: number;
  /** 缓动函数 */
  easing?: EasingFunction;
  /** 延迟时间（毫秒） */
  delay?: number;
  /** 进入前回调 */
  onBeforeEnter?: () => void;
  /** 进入后回调 */
  onAfterEnter?: () => void;
  /** 离开前回调 */
  onBeforeLeave?: () => void;
  /** 离开后回调 */
  onAfterLeave?: () => void;
}

/**
 * 过渡返回值
 */
export interface UseTransitionReturn {
  /** 是否可见 */
  isVisible: Ref<boolean>;
  /** 当前状态 */
  state: Ref<TransitionState>;
  /** 是否正在过渡 */
  isTransitioning: ComputedRef<boolean>;
  /** 显示 */
  show: () => Promise<void>;
  /** 隐藏 */
  hide: () => Promise<void>;
  /** 切换 */
  toggle: () => Promise<void>;
  /** 过渡类名 */
  transitionClass: ComputedRef<string>;
  /** 过渡样式 */
  transitionStyle: ComputedRef<Record<string, string>>;
}

/**
 * 动画配置
 */
export interface AnimationConfig {
  /** 动画名称 */
  name: string;
  /** 持续时间（毫秒） */
  duration?: number;
  /** 缓动函数 */
  easing?: EasingFunction;
  /** 延迟时间（毫秒） */
  delay?: number;
  /** 迭代次数 */
  iterations?: number | "infinite";
  /** 动画方向 */
  direction?: "normal" | "reverse" | "alternate" | "alternate-reverse";
  /** 填充模式 */
  fillMode?: "none" | "forwards" | "backwards" | "both";
}

/**
 * 动画返回值
 */
export interface UseAnimationReturn {
  /** 是否正在播放 */
  isPlaying: Ref<boolean>;
  /** 是否已暂停 */
  isPaused: Ref<boolean>;
  /** 播放 */
  play: () => void;
  /** 暂停 */
  pause: () => void;
  /** 停止 */
  stop: () => void;
  /** 重新开始 */
  restart: () => void;
  /** 动画样式 */
  animationStyle: ComputedRef<Record<string, string>>;
}

/**
 * 淡入淡出配置
 */
export interface FadeConfig {
  /** 持续时间 */
  duration?: number;
  /** 缓动函数 */
  easing?: EasingFunction;
}

/**
 * 滑动配置
 */
export interface SlideConfig extends FadeConfig {
  /** 滑动方向 */
  direction?: "up" | "down" | "left" | "right";
  /** 滑动距离 */
  distance?: string;
}

// ============================================================================
// 核心函数
// ============================================================================

/**
 * 使用过渡
 *
 * @description 创建可控的过渡效果
 *
 * @example
 * ```ts
 * const { isVisible, show, hide, toggle, state, transitionClass } = useTransition({
 *   duration: 300,
 *   easing: 'ease-out'
 * })
 *
 * // 显示
 * await show()
 *
 * // 隐藏
 * await hide()
 *
 * // 在模板中使用
 * // <div v-show="isVisible" :class="transitionClass">Content</div>
 * ```
 */
export function useTransition(
  config: TransitionConfig = {}
): UseTransitionReturn {
  const {
    duration = 300,
    easing = "ease",
    delay = 0,
    onBeforeEnter,
    onAfterEnter,
    onBeforeLeave,
    onAfterLeave,
  } = config;

  const isVisible = ref(false);
  const state = ref<TransitionState>("idle");

  const isTransitioning = computed(
    () => state.value === "enter" || state.value === "leave"
  );

  const show = async (): Promise<void> => {
    if (isVisible.value && state.value === "idle") return;

    onBeforeEnter?.();
    state.value = "enter";
    isVisible.value = true;

    await nextTick();

    return new Promise((resolve) => {
      setTimeout(() => {
        state.value = "idle";
        onAfterEnter?.();
        resolve();
      }, duration + delay);
    });
  };

  const hide = async (): Promise<void> => {
    if (!isVisible.value && state.value === "idle") return;

    onBeforeLeave?.();
    state.value = "leave";

    return new Promise((resolve) => {
      setTimeout(() => {
        isVisible.value = false;
        state.value = "idle";
        onAfterLeave?.();
        resolve();
      }, duration + delay);
    });
  };

  const toggle = async (): Promise<void> => {
    if (isVisible.value) {
      await hide();
    } else {
      await show();
    }
  };

  const transitionClass = computed(() => {
    switch (state.value) {
      case "enter":
        return "transition-enter";
      case "leave":
        return "transition-leave";
      default:
        return "";
    }
  });

  const transitionStyle = computed(() => ({
    transition: `all ${duration}ms ${easing}`,
    transitionDelay: delay > 0 ? `${delay}ms` : undefined,
  }));

  return {
    isVisible,
    state,
    isTransitioning,
    show,
    hide,
    toggle,
    transitionClass,
    transitionStyle,
  };
}

/**
 * 使用淡入淡出
 *
 * @description 创建淡入淡出效果
 *
 * @example
 * ```ts
 * const { isVisible, opacity, show, hide, toggle, style } = useFade({
 *   duration: 300
 * })
 *
 * // <div v-show="isVisible" :style="style">Content</div>
 * ```
 */
export function useFade(config: FadeConfig = {}): {
  isVisible: Ref<boolean>;
  opacity: Ref<number>;
  show: () => Promise<void>;
  hide: () => Promise<void>;
  toggle: () => Promise<void>;
  style: ComputedRef<Record<string, string>>;
} {
  const { duration = 300, easing = "ease" } = config;

  const isVisible = ref(false);
  const opacity = ref(0);

  const style = computed(() => ({
    opacity: String(opacity.value),
    transition: `opacity ${duration}ms ${easing}`,
  }));

  const show = async (): Promise<void> => {
    isVisible.value = true;
    await nextTick();
    opacity.value = 1;

    return new Promise((resolve) => {
      setTimeout(resolve, duration);
    });
  };

  const hide = async (): Promise<void> => {
    opacity.value = 0;

    return new Promise((resolve) => {
      setTimeout(() => {
        isVisible.value = false;
        resolve();
      }, duration);
    });
  };

  const toggle = () => (isVisible.value ? hide() : show());

  return {
    isVisible,
    opacity,
    show,
    hide,
    toggle,
    style,
  };
}

/**
 * 使用滑动
 *
 * @description 创建滑动效果
 *
 * @example
 * ```ts
 * const { isVisible, show, hide, style } = useSlide({
 *   direction: 'up',
 *   distance: '20px',
 *   duration: 300
 * })
 * ```
 */
export function useSlide(config: SlideConfig = {}): {
  isVisible: Ref<boolean>;
  show: () => Promise<void>;
  hide: () => Promise<void>;
  toggle: () => Promise<void>;
  style: ComputedRef<Record<string, string>>;
} {
  const {
    duration = 300,
    easing = "ease",
    direction = "up",
    distance = "20px",
  } = config;

  const isVisible = ref(false);
  const progress = ref(0);

  const getTransform = (p: number) => {
    const d = `calc(${distance} * ${1 - p})`;
    switch (direction) {
      case "up":
        return `translateY(${d})`;
      case "down":
        return `translateY(calc(-1 * ${d}))`;
      case "left":
        return `translateX(${d})`;
      case "right":
        return `translateX(calc(-1 * ${d}))`;
      default:
        return "none";
    }
  };

  const style = computed(() => ({
    opacity: String(progress.value),
    transform: getTransform(progress.value),
    transition: `all ${duration}ms ${easing}`,
  }));

  const show = async (): Promise<void> => {
    isVisible.value = true;
    await nextTick();
    progress.value = 1;

    return new Promise((resolve) => {
      setTimeout(resolve, duration);
    });
  };

  const hide = async (): Promise<void> => {
    progress.value = 0;

    return new Promise((resolve) => {
      setTimeout(() => {
        isVisible.value = false;
        resolve();
      }, duration);
    });
  };

  const toggle = () => (isVisible.value ? hide() : show());

  return {
    isVisible,
    show,
    hide,
    toggle,
    style,
  };
}

/**
 * 使用缩放
 *
 * @description 创建缩放效果
 *
 * @example
 * ```ts
 * const { isVisible, show, hide, style } = useScale({
 *   fromScale: 0.8,
 *   duration: 200
 * })
 * ```
 */
export function useScale(
  config: FadeConfig & { fromScale?: number } = {}
): {
  isVisible: Ref<boolean>;
  scale: Ref<number>;
  show: () => Promise<void>;
  hide: () => Promise<void>;
  toggle: () => Promise<void>;
  style: ComputedRef<Record<string, string>>;
} {
  const { duration = 200, easing = "ease", fromScale = 0.9 } = config;

  const isVisible = ref(false);
  const scale = ref(fromScale);
  const opacity = ref(0);

  const style = computed(() => ({
    opacity: String(opacity.value),
    transform: `scale(${scale.value})`,
    transition: `all ${duration}ms ${easing}`,
  }));

  const show = async (): Promise<void> => {
    isVisible.value = true;
    await nextTick();
    scale.value = 1;
    opacity.value = 1;

    return new Promise((resolve) => {
      setTimeout(resolve, duration);
    });
  };

  const hide = async (): Promise<void> => {
    scale.value = fromScale;
    opacity.value = 0;

    return new Promise((resolve) => {
      setTimeout(() => {
        isVisible.value = false;
        resolve();
      }, duration);
    });
  };

  const toggle = () => (isVisible.value ? hide() : show());

  return {
    isVisible,
    scale,
    show,
    hide,
    toggle,
    style,
  };
}

/**
 * 使用动画
 *
 * @description 创建可控的 CSS 动画
 *
 * @example
 * ```ts
 * const { isPlaying, play, pause, stop, restart, animationStyle } = useAnimation({
 *   name: 'bounce',
 *   duration: 1000,
 *   iterations: 'infinite'
 * })
 *
 * // <div :style="animationStyle">Animated</div>
 * ```
 */
export function useAnimation(config: AnimationConfig): UseAnimationReturn {
  const {
    name,
    duration = 1000,
    easing = "ease",
    delay = 0,
    iterations = 1,
    direction = "normal",
    fillMode = "none",
  } = config;

  const isPlaying = ref(false);
  const isPaused = ref(false);

  const animationStyle = computed(() => {
    if (!isPlaying.value) {
      return {};
    }

    return {
      animationName: name,
      animationDuration: `${duration}ms`,
      animationTimingFunction: easing,
      animationDelay: `${delay}ms`,
      animationIterationCount: iterations === "infinite" ? "infinite" : String(iterations),
      animationDirection: direction,
      animationFillMode: fillMode,
      animationPlayState: isPaused.value ? "paused" : "running",
    };
  });

  const play = () => {
    isPlaying.value = true;
    isPaused.value = false;
  };

  const pause = () => {
    isPaused.value = true;
  };

  const stop = () => {
    isPlaying.value = false;
    isPaused.value = false;
  };

  const restart = () => {
    stop();
    nextTick(() => {
      play();
    });
  };

  return {
    isPlaying,
    isPaused,
    play,
    pause,
    stop,
    restart,
    animationStyle,
  };
}

/**
 * 使用过渡组
 *
 * @description 管理列表项的过渡效果
 *
 * @example
 * ```ts
 * const { items, addItem, removeItem, getItemStyle } = useTransitionGroup({
 *   duration: 300
 * })
 *
 * addItem({ id: 1, data: 'Item 1' })
 *
 * // <div v-for="item in items" :key="item.id" :style="getItemStyle(item.id)">
 * //   {{ item.data }}
 * // </div>
 * ```
 */
export function useTransitionGroup<T extends { id: string | number }>(
  config: TransitionConfig = {}
): {
  items: Ref<T[]>;
  addItem: (item: T) => void;
  removeItem: (id: string | number) => void;
  getItemStyle: (id: string | number) => Record<string, string>;
} {
  const { duration = 300, easing = "ease" } = config;

  const items = ref<T[]>([]) as Ref<T[]>;
  const enteringIds = ref<Set<string | number>>(new Set());
  const leavingIds = ref<Set<string | number>>(new Set());

  const addItem = (item: T) => {
    enteringIds.value.add(item.id);
    items.value = [...items.value, item];

    setTimeout(() => {
      enteringIds.value.delete(item.id);
    }, duration);
  };

  const removeItem = (id: string | number) => {
    leavingIds.value.add(id);

    setTimeout(() => {
      items.value = items.value.filter((item) => item.id !== id);
      leavingIds.value.delete(id);
    }, duration);
  };

  const getItemStyle = (id: string | number): Record<string, string> => {
    const isEntering = enteringIds.value.has(id);
    const isLeaving = leavingIds.value.has(id);

    return {
      transition: `all ${duration}ms ${easing}`,
      opacity: isEntering || isLeaving ? "0" : "1",
      transform: isEntering
        ? "translateY(-10px)"
        : isLeaving
          ? "translateY(10px)"
          : "none",
    };
  };

  return {
    items,
    addItem,
    removeItem,
    getItemStyle,
  };
}

/**
 * 使用数值过渡
 *
 * @description 创建数值的平滑过渡效果
 *
 * @example
 * ```ts
 * const { value, set, tweenedValue } = useNumberTransition(0, {
 *   duration: 500
 * })
 *
 * set(100) // value 将从 0 平滑过渡到 100
 *
 * // <span>{{ Math.round(tweenedValue) }}</span>
 * ```
 */
export function useNumberTransition(
  initial: number,
  config: { duration?: number; easing?: (t: number) => number } = {}
): {
  value: Ref<number>;
  tweenedValue: Ref<number>;
  set: (target: number) => void;
  isAnimating: Ref<boolean>;
} {
  const { duration = 500, easing = (t: number) => t } = config;

  const value = ref(initial);
  const tweenedValue = ref(initial);
  const isAnimating = ref(false);

  let animationFrame: number | null = null;

  const set = (target: number) => {
    if (animationFrame) {
      cancelAnimationFrame(animationFrame);
    }

    const startValue = tweenedValue.value;
    const startTime = Date.now();
    value.value = target;
    isAnimating.value = true;

    const animate = () => {
      const elapsed = Date.now() - startTime;
      const progress = Math.min(elapsed / duration, 1);
      const easedProgress = easing(progress);

      tweenedValue.value = startValue + (target - startValue) * easedProgress;

      if (progress < 1) {
        animationFrame = requestAnimationFrame(animate);
      } else {
        isAnimating.value = false;
        animationFrame = null;
      }
    };

    animationFrame = requestAnimationFrame(animate);
  };

  onUnmounted(() => {
    if (animationFrame) {
      cancelAnimationFrame(animationFrame);
    }
  });

  return {
    value,
    tweenedValue,
    set,
    isAnimating,
  };
}

/**
 * 使用抖动效果
 *
 * @description 创建抖动动画效果（用于错误提示等）
 *
 * @example
 * ```ts
 * const { shake, isShaking, style } = useShake()
 *
 * // 触发抖动
 * shake()
 *
 * // <div :style="style">Content</div>
 * ```
 */
export function useShake(config: { duration?: number; intensity?: number } = {}): {
  shake: () => void;
  isShaking: Ref<boolean>;
  style: ComputedRef<Record<string, string>>;
} {
  const { duration = 500, intensity = 10 } = config;

  const isShaking = ref(false);
  const offset = ref(0);

  let animationFrame: number | null = null;

  const shake = () => {
    if (isShaking.value) return;

    isShaking.value = true;
    const startTime = Date.now();

    const animate = () => {
      const elapsed = Date.now() - startTime;
      const progress = elapsed / duration;

      if (progress < 1) {
        // 衰减的正弦波
        const decay = 1 - progress;
        const wave = Math.sin(progress * Math.PI * 8);
        offset.value = wave * intensity * decay;
        animationFrame = requestAnimationFrame(animate);
      } else {
        offset.value = 0;
        isShaking.value = false;
        animationFrame = null;
      }
    };

    animationFrame = requestAnimationFrame(animate);
  };

  const style = computed(() => ({
    transform: `translateX(${offset.value}px)`,
  }));

  onUnmounted(() => {
    if (animationFrame) {
      cancelAnimationFrame(animationFrame);
    }
  });

  return {
    shake,
    isShaking,
    style,
  };
}

/**
 * 使用脉冲效果
 *
 * @description 创建脉冲动画效果
 *
 * @example
 * ```ts
 * const { pulse, isPulsing, stop, style } = usePulse({
 *   scale: 1.1,
 *   duration: 300
 * })
 *
 * // 单次脉冲
 * pulse()
 *
 * // <div :style="style">Content</div>
 * ```
 */
export function usePulse(
  config: { scale?: number; duration?: number } = {}
): {
  pulse: () => void;
  isPulsing: Ref<boolean>;
  stop: () => void;
  style: ComputedRef<Record<string, string>>;
} {
  const { scale = 1.1, duration = 300 } = config;

  const isPulsing = ref(false);
  const currentScale = ref(1);

  let timeout: ReturnType<typeof setTimeout> | null = null;

  const pulse = () => {
    if (isPulsing.value) return;

    isPulsing.value = true;
    currentScale.value = scale;

    timeout = setTimeout(() => {
      currentScale.value = 1;

      timeout = setTimeout(() => {
        isPulsing.value = false;
      }, duration);
    }, duration);
  };

  const stop = () => {
    if (timeout) {
      clearTimeout(timeout);
      timeout = null;
    }
    currentScale.value = 1;
    isPulsing.value = false;
  };

  const style = computed(() => ({
    transform: `scale(${currentScale.value})`,
    transition: `transform ${duration}ms ease`,
  }));

  onUnmounted(stop);

  return {
    pulse,
    isPulsing,
    stop,
    style,
  };
}

/**
 * 使用打字机效果
 *
 * @description 创建打字机文本动画
 *
 * @example
 * ```ts
 * const { text, start, pause, reset, isTyping } = useTypewriter('Hello World', {
 *   speed: 50
 * })
 *
 * start()
 *
 * // <span>{{ text }}</span>
 * ```
 */
export function useTypewriter(
  fullText: string,
  config: { speed?: number; delay?: number } = {}
): {
  text: Ref<string>;
  start: () => void;
  pause: () => void;
  reset: () => void;
  isTyping: Ref<boolean>;
  isComplete: ComputedRef<boolean>;
} {
  const { speed = 50, delay = 0 } = config;

  const text = ref("");
  const isTyping = ref(false);
  const currentIndex = ref(0);

  let timer: ReturnType<typeof setTimeout> | null = null;

  const isComplete = computed(() => text.value === fullText);

  const typeNextChar = () => {
    if (currentIndex.value < fullText.length) {
      text.value += fullText[currentIndex.value];
      currentIndex.value++;
      timer = setTimeout(typeNextChar, speed);
    } else {
      isTyping.value = false;
    }
  };

  const start = () => {
    if (isComplete.value) return;

    isTyping.value = true;
    timer = setTimeout(typeNextChar, delay);
  };

  const pause = () => {
    if (timer) {
      clearTimeout(timer);
      timer = null;
    }
    isTyping.value = false;
  };

  const reset = () => {
    pause();
    text.value = "";
    currentIndex.value = 0;
  };

  onUnmounted(pause);

  return {
    text,
    start,
    pause,
    reset,
    isTyping,
    isComplete,
  };
}
