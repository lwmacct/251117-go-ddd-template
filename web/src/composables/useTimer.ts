/**
 * Timer Composable
 * 提供定时器、间隔器和延迟执行功能
 */

import { ref, onUnmounted, type Ref, computed, watch } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UseTimeoutOptions {
  /** 是否立即启动 */
  immediate?: boolean;
  /** 回调函数 */
  callback?: () => void;
}

export interface UseTimeoutReturn {
  /** 是否准备就绪（超时已触发） */
  ready: Ref<boolean>;
  /** 是否正在运行 */
  isPending: Ref<boolean>;
  /** 启动定时器 */
  start: () => void;
  /** 停止定时器 */
  stop: () => void;
}

export interface UseIntervalOptions {
  /** 是否立即启动 */
  immediate?: boolean;
  /** 是否立即执行回调 */
  immediateCallback?: boolean;
}

export interface UseIntervalReturn {
  /** 当前计数 */
  counter: Ref<number>;
  /** 是否正在运行 */
  isActive: Ref<boolean>;
  /** 暂停 */
  pause: () => void;
  /** 恢复 */
  resume: () => void;
  /** 重置计数 */
  reset: () => void;
}

export interface UseTimestampOptions {
  /** 更新间隔（毫秒） */
  interval?: number;
  /** 是否立即启动 */
  immediate?: boolean;
  /** 时间偏移（毫秒） */
  offset?: number;
}

export interface UseTimestampReturn {
  /** 当前时间戳 */
  timestamp: Ref<number>;
  /** 是否正在运行 */
  isActive: Ref<boolean>;
  /** 暂停 */
  pause: () => void;
  /** 恢复 */
  resume: () => void;
}

// ============================================================================
// useTimeout - 超时定时器
// ============================================================================

/**
 * 超时定时器
 * @example
 * const { ready, start, stop } = useTimeout(3000)
 * start()
 * watch(ready, (val) => { if (val) console.log('Timeout!') })
 */
export function useTimeout(ms: number | Ref<number>, options: UseTimeoutOptions = {}): UseTimeoutReturn {
  const { immediate = false, callback } = options;

  const ready = ref(false);
  const isPending = ref(false);
  let timer: ReturnType<typeof setTimeout> | null = null;

  const stop = () => {
    if (timer) {
      clearTimeout(timer);
      timer = null;
    }
    isPending.value = false;
  };

  const start = () => {
    stop();
    ready.value = false;
    isPending.value = true;

    const delay = typeof ms === "number" ? ms : ms.value;

    timer = setTimeout(() => {
      ready.value = true;
      isPending.value = false;
      callback?.();
    }, delay);
  };

  if (immediate) {
    start();
  }

  onUnmounted(stop);

  return {
    ready,
    isPending,
    start,
    stop,
  };
}

/**
 * 可控的超时定时器（响应式延迟）
 * @example
 * const delay = ref(1000)
 * const { ready, start } = useTimeoutFn(() => console.log('Done'), delay)
 */
export function useTimeoutFn<T extends (...args: unknown[]) => unknown>(
  fn: T,
  ms: number | Ref<number>,
  options: { immediate?: boolean } = {}
): UseTimeoutReturn & { call: (...args: Parameters<T>) => void } {
  const { immediate = false } = options;

  const ready = ref(false);
  const isPending = ref(false);
  let timer: ReturnType<typeof setTimeout> | null = null;
  let pendingArgs: Parameters<T> | null = null;

  const stop = () => {
    if (timer) {
      clearTimeout(timer);
      timer = null;
    }
    isPending.value = false;
    pendingArgs = null;
  };

  const start = () => {
    stop();
    ready.value = false;
    isPending.value = true;

    const delay = typeof ms === "number" ? ms : ms.value;

    timer = setTimeout(() => {
      ready.value = true;
      isPending.value = false;
      if (pendingArgs) {
        fn(...pendingArgs);
      } else {
        fn();
      }
    }, delay);
  };

  const call = (...args: Parameters<T>) => {
    pendingArgs = args;
    start();
  };

  if (immediate) {
    start();
  }

  onUnmounted(stop);

  return {
    ready,
    isPending,
    start,
    stop,
    call,
  };
}

// ============================================================================
// useInterval - 间隔定时器
// ============================================================================

/**
 * 间隔定时器
 * @example
 * const { counter, pause, resume, reset } = useInterval(1000)
 * // counter 每秒加 1
 */
export function useInterval(ms: number | Ref<number>, options: UseIntervalOptions = {}): UseIntervalReturn {
  const { immediate = true, immediateCallback = false } = options;

  const counter = ref(0);
  const isActive = ref(false);
  let timer: ReturnType<typeof setInterval> | null = null;

  const pause = () => {
    if (timer) {
      clearInterval(timer);
      timer = null;
    }
    isActive.value = false;
  };

  const resume = () => {
    if (isActive.value) return;

    isActive.value = true;
    const delay = typeof ms === "number" ? ms : ms.value;

    timer = setInterval(() => {
      counter.value++;
    }, delay);
  };

  const reset = () => {
    counter.value = 0;
  };

  if (immediate) {
    resume();
  }

  if (immediateCallback) {
    counter.value++;
  }

  onUnmounted(pause);

  return {
    counter,
    isActive,
    pause,
    resume,
    reset,
  };
}

/**
 * 带回调的间隔定时器
 * @example
 * useIntervalFn(() => {
 *   console.log('Tick!')
 * }, 1000)
 */
export function useIntervalFn<T extends (...args: unknown[]) => unknown>(
  fn: T,
  ms: number | Ref<number>,
  options: UseIntervalOptions = {}
): UseIntervalReturn {
  const { immediate = true, immediateCallback = false } = options;

  const counter = ref(0);
  const isActive = ref(false);
  let timer: ReturnType<typeof setInterval> | null = null;

  const pause = () => {
    if (timer) {
      clearInterval(timer);
      timer = null;
    }
    isActive.value = false;
  };

  const resume = () => {
    if (isActive.value) return;

    isActive.value = true;
    const delay = typeof ms === "number" ? ms : ms.value;

    timer = setInterval(() => {
      counter.value++;
      fn();
    }, delay);
  };

  const reset = () => {
    counter.value = 0;
  };

  if (immediate) {
    resume();
  }

  if (immediateCallback) {
    fn();
  }

  onUnmounted(pause);

  return {
    counter,
    isActive,
    pause,
    resume,
    reset,
  };
}

// ============================================================================
// useTimestamp - 实时时间戳
// ============================================================================

/**
 * 实时时间戳
 * @example
 * const { timestamp } = useTimestamp()
 * const formattedTime = computed(() => new Date(timestamp.value).toLocaleTimeString())
 */
export function useTimestamp(options: UseTimestampOptions = {}): UseTimestampReturn {
  const { interval = 1000, immediate = true, offset = 0 } = options;

  const timestamp = ref(Date.now() + offset);
  const isActive = ref(false);
  let timer: ReturnType<typeof setInterval> | null = null;

  const update = () => {
    timestamp.value = Date.now() + offset;
  };

  const pause = () => {
    if (timer) {
      clearInterval(timer);
      timer = null;
    }
    isActive.value = false;
  };

  const resume = () => {
    if (isActive.value) return;

    update();
    isActive.value = true;
    timer = setInterval(update, interval);
  };

  if (immediate) {
    resume();
  }

  onUnmounted(pause);

  return {
    timestamp,
    isActive,
    pause,
    resume,
  };
}

/**
 * 实时日期时间
 * @example
 * const { now, date, time } = useNow()
 */
export function useNow(options: UseTimestampOptions = {}) {
  const { timestamp, isActive, pause, resume } = useTimestamp(options);

  const now = computed(() => new Date(timestamp.value));

  const date = computed(() => {
    const d = now.value;
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
  });

  const time = computed(() => {
    const d = now.value;
    return `${String(d.getHours()).padStart(2, "0")}:${String(d.getMinutes()).padStart(2, "0")}:${String(d.getSeconds()).padStart(2, "0")}`;
  });

  return {
    now,
    date,
    time,
    timestamp,
    isActive,
    pause,
    resume,
  };
}

// ============================================================================
// useRafFn - requestAnimationFrame
// ============================================================================

export interface UseRafFnOptions {
  /** 是否立即启动 */
  immediate?: boolean;
}

export interface UseRafFnReturn {
  /** 是否正在运行 */
  isActive: Ref<boolean>;
  /** 暂停 */
  pause: () => void;
  /** 恢复 */
  resume: () => void;
}

/**
 * requestAnimationFrame 循环
 * @example
 * useRafFn(() => {
 *   // 每帧执行
 *   updateAnimation()
 * })
 */
export function useRafFn(fn: (timestamp: number) => void, options: UseRafFnOptions = {}): UseRafFnReturn {
  const { immediate = true } = options;

  const isActive = ref(false);
  let rafId: number | null = null;

  const loop = (timestamp: number) => {
    if (!isActive.value) return;

    fn(timestamp);
    rafId = requestAnimationFrame(loop);
  };

  const pause = () => {
    if (rafId !== null) {
      cancelAnimationFrame(rafId);
      rafId = null;
    }
    isActive.value = false;
  };

  const resume = () => {
    if (isActive.value) return;

    isActive.value = true;
    rafId = requestAnimationFrame(loop);
  };

  if (immediate) {
    resume();
  }

  onUnmounted(pause);

  return {
    isActive,
    pause,
    resume,
  };
}

// ============================================================================
// useDateFormat - 日期格式化
// ============================================================================

export interface UseDateFormatOptions {
  /** 更新间隔（毫秒），0 表示不自动更新 */
  updateInterval?: number;
}

/**
 * 响应式日期格式化
 * @example
 * const { formatted } = useDateFormat('YYYY-MM-DD HH:mm:ss')
 */
export function useDateFormat(
  format: string,
  date?: Date | Ref<Date> | number | Ref<number>,
  options: UseDateFormatOptions = {}
) {
  const { updateInterval = 0 } = options;

  const getDate = () => {
    if (!date) return new Date();
    const d = typeof date === "number" || date instanceof Date ? date : date.value;
    return d instanceof Date ? d : new Date(d);
  };

  const formatDate = (d: Date, fmt: string): string => {
    const tokens: Record<string, () => string> = {
      YYYY: () => String(d.getFullYear()),
      YY: () => String(d.getFullYear()).slice(-2),
      MM: () => String(d.getMonth() + 1).padStart(2, "0"),
      M: () => String(d.getMonth() + 1),
      DD: () => String(d.getDate()).padStart(2, "0"),
      D: () => String(d.getDate()),
      HH: () => String(d.getHours()).padStart(2, "0"),
      H: () => String(d.getHours()),
      hh: () => String(d.getHours() % 12 || 12).padStart(2, "0"),
      h: () => String(d.getHours() % 12 || 12),
      mm: () => String(d.getMinutes()).padStart(2, "0"),
      m: () => String(d.getMinutes()),
      ss: () => String(d.getSeconds()).padStart(2, "0"),
      s: () => String(d.getSeconds()),
      SSS: () => String(d.getMilliseconds()).padStart(3, "0"),
      A: () => (d.getHours() < 12 ? "AM" : "PM"),
      a: () => (d.getHours() < 12 ? "am" : "pm"),
    };

    let result = fmt;
    for (const [token, fn] of Object.entries(tokens)) {
      result = result.replace(new RegExp(token, "g"), fn());
    }

    return result;
  };

  const formatted = ref(formatDate(getDate(), format));

  const update = () => {
    formatted.value = formatDate(getDate(), format);
  };

  let timer: ReturnType<typeof setInterval> | null = null;

  if (updateInterval > 0) {
    timer = setInterval(update, updateInterval);
  }

  // 监听响应式日期变化
  if (date && typeof date !== "number" && !(date instanceof Date)) {
    watch(date, update);
  }

  onUnmounted(() => {
    if (timer) {
      clearInterval(timer);
    }
  });

  return {
    formatted,
    update,
  };
}

// ============================================================================
// useIdleCallback - 空闲回调
// ============================================================================

export interface UseIdleCallbackOptions {
  /** 超时时间（毫秒） */
  timeout?: number;
}

export interface UseIdleCallbackReturn {
  /** 是否正在运行 */
  isSupported: boolean;
  /** 取消回调 */
  cancel: () => void;
}

/**
 * 空闲时执行回调
 * @example
 * useIdleCallback(() => {
 *   // 浏览器空闲时执行
 *   heavyComputation()
 * })
 */
export function useIdleCallback(fn: IdleRequestCallback, options: UseIdleCallbackOptions = {}): UseIdleCallbackReturn {
  const { timeout } = options;

  const isSupported = typeof window !== "undefined" && "requestIdleCallback" in window;

  let idleId: number | null = null;

  const cancel = () => {
    if (idleId !== null && isSupported) {
      (window as Window).cancelIdleCallback(idleId);
      idleId = null;
    }
  };

  if (isSupported) {
    idleId = (window as Window).requestIdleCallback(fn, timeout ? { timeout } : undefined);
  } else {
    // 回退到 setTimeout
    const timeoutId = setTimeout(() => {
      fn({
        didTimeout: false,
        timeRemaining: () => 0,
      } as IdleDeadline);
    }, 1);

    onUnmounted(() => clearTimeout(timeoutId));
  }

  onUnmounted(cancel);

  return {
    isSupported,
    cancel,
  };
}

// ============================================================================
// useScheduler - 任务调度器
// ============================================================================

export interface ScheduledTask {
  id: string;
  fn: () => void;
  interval: number;
  isActive: boolean;
}

export interface UseSchedulerReturn {
  /** 所有任务 */
  tasks: Ref<ScheduledTask[]>;
  /** 添加任务 */
  addTask: (id: string, fn: () => void, interval: number) => void;
  /** 移除任务 */
  removeTask: (id: string) => void;
  /** 暂停任务 */
  pauseTask: (id: string) => void;
  /** 恢复任务 */
  resumeTask: (id: string) => void;
  /** 暂停所有任务 */
  pauseAll: () => void;
  /** 恢复所有任务 */
  resumeAll: () => void;
  /** 清除所有任务 */
  clearAll: () => void;
}

/**
 * 任务调度器
 * @example
 * const { addTask, pauseTask } = useScheduler()
 * addTask('sync', () => syncData(), 5000)
 * addTask('ping', () => ping(), 30000)
 */
export function useScheduler(): UseSchedulerReturn {
  const tasks = ref<ScheduledTask[]>([]);
  const timers = new Map<string, ReturnType<typeof setInterval>>();

  const startTimer = (task: ScheduledTask) => {
    if (timers.has(task.id)) {
      clearInterval(timers.get(task.id)!);
    }

    const timer = setInterval(() => {
      if (task.isActive) {
        task.fn();
      }
    }, task.interval);

    timers.set(task.id, timer);
  };

  const addTask = (id: string, fn: () => void, interval: number) => {
    const existingIndex = tasks.value.findIndex((t) => t.id === id);

    if (existingIndex !== -1) {
      // 更新现有任务
      removeTask(id);
    }

    const task: ScheduledTask = {
      id,
      fn,
      interval,
      isActive: true,
    };

    tasks.value.push(task);
    startTimer(task);
  };

  const removeTask = (id: string) => {
    const timer = timers.get(id);
    if (timer) {
      clearInterval(timer);
      timers.delete(id);
    }

    const index = tasks.value.findIndex((t) => t.id === id);
    if (index !== -1) {
      tasks.value.splice(index, 1);
    }
  };

  const pauseTask = (id: string) => {
    const task = tasks.value.find((t) => t.id === id);
    if (task) {
      task.isActive = false;
    }
  };

  const resumeTask = (id: string) => {
    const task = tasks.value.find((t) => t.id === id);
    if (task) {
      task.isActive = true;
    }
  };

  const pauseAll = () => {
    tasks.value.forEach((task) => {
      task.isActive = false;
    });
  };

  const resumeAll = () => {
    tasks.value.forEach((task) => {
      task.isActive = true;
    });
  };

  const clearAll = () => {
    timers.forEach((timer) => clearInterval(timer));
    timers.clear();
    tasks.value = [];
  };

  onUnmounted(clearAll);

  return {
    tasks,
    addTask,
    removeTask,
    pauseTask,
    resumeTask,
    pauseAll,
    resumeAll,
    clearAll,
  };
}
