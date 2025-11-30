/**
 * Watch Composable
 * 提供增强的 watch 工具函数
 */

import {
  ref,
  watch,
  watchEffect,
  computed,
  nextTick,
  onMounted,
  type Ref,
  type WatchSource,
  type WatchOptions,
  type WatchStopHandle,
  type ComputedRef,
} from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface WatchOnceOptions extends WatchOptions {
  /** 是否立即执行 */
  immediate?: boolean;
}

export interface WatchDebouncedOptions extends WatchOptions {
  /** 防抖延迟（毫秒） */
  debounce?: number;
  /** 最大等待时间（毫秒） */
  maxWait?: number;
}

export interface WatchThrottledOptions extends WatchOptions {
  /** 节流间隔（毫秒） */
  throttle?: number;
  /** 是否在开始时执行 */
  leading?: boolean;
  /** 是否在结束时执行 */
  trailing?: boolean;
}

export interface WatchPausableReturn {
  /** 暂停监听 */
  pause: () => void;
  /** 恢复监听 */
  resume: () => void;
  /** 是否暂停中 */
  isActive: Ref<boolean>;
  /** 停止监听 */
  stop: WatchStopHandle;
}

export interface WatchIgnorableReturn {
  /** 忽略更新并执行函数 */
  ignoreUpdates: (updater: () => void) => void;
  /** 停止监听 */
  stop: WatchStopHandle;
}

export interface WatchTriggeredReturn<T> {
  /** 触发次数 */
  count: Ref<number>;
  /** 是否已触发 */
  isTriggered: ComputedRef<boolean>;
  /** 重置计数 */
  reset: () => void;
  /** 停止监听 */
  stop: WatchStopHandle;
}

// ============================================================================
// watchOnce - 只执行一次的监听
// ============================================================================

/**
 * 只执行一次的 watch
 * @example
 * // 等待数据加载完成
 * watchOnce(
 *   () => data.value,
 *   (newData) => {
 *     console.log('数据已加载:', newData)
 *   }
 * )
 */
export function watchOnce<T>(
  source: WatchSource<T>,
  callback: (value: T, oldValue: T | undefined) => void,
  options: WatchOnceOptions = {}
): WatchStopHandle {
  const stop = watch(
    source,
    (newValue, oldValue) => {
      callback(newValue, oldValue);
      nextTick(() => stop());
    },
    options
  );

  return stop;
}

// ============================================================================
// watchDebounced - 防抖的监听
// ============================================================================

/**
 * 防抖的 watch
 * @example
 * const searchQuery = ref('')
 *
 * watchDebounced(
 *   searchQuery,
 *   (query) => {
 *     performSearch(query)
 *   },
 *   { debounce: 500 }
 * )
 */
export function watchDebounced<T>(
  source: WatchSource<T>,
  callback: (value: T, oldValue: T | undefined) => void,
  options: WatchDebouncedOptions = {}
): WatchStopHandle {
  const { debounce = 250, maxWait, ...watchOptions } = options;

  let timer: ReturnType<typeof setTimeout> | null = null;
  let maxTimer: ReturnType<typeof setTimeout> | null = null;
  let lastCallTime = 0;

  const clear = () => {
    if (timer) {
      clearTimeout(timer);
      timer = null;
    }
    if (maxTimer) {
      clearTimeout(maxTimer);
      maxTimer = null;
    }
  };

  const stop = watch(
    source,
    (newValue, oldValue) => {
      const now = Date.now();

      const invokeCallback = () => {
        clear();
        callback(newValue, oldValue);
        lastCallTime = Date.now();
      };

      clear();

      // 设置防抖定时器
      timer = setTimeout(invokeCallback, debounce);

      // 设置最大等待定时器
      if (maxWait !== undefined && maxWait > 0) {
        const timeSinceLastCall = now - lastCallTime;
        if (timeSinceLastCall >= maxWait) {
          invokeCallback();
        } else {
          maxTimer = setTimeout(invokeCallback, maxWait - timeSinceLastCall);
        }
      }
    },
    watchOptions
  );

  return () => {
    clear();
    stop();
  };
}

// ============================================================================
// watchThrottled - 节流的监听
// ============================================================================

/**
 * 节流的 watch
 * @example
 * const scrollPosition = ref(0)
 *
 * watchThrottled(
 *   scrollPosition,
 *   (pos) => {
 *     updateVisibleItems(pos)
 *   },
 *   { throttle: 100 }
 * )
 */
export function watchThrottled<T>(
  source: WatchSource<T>,
  callback: (value: T, oldValue: T | undefined) => void,
  options: WatchThrottledOptions = {}
): WatchStopHandle {
  const {
    throttle = 100,
    leading = true,
    trailing = true,
    ...watchOptions
  } = options;

  let lastCallTime = 0;
  let timer: ReturnType<typeof setTimeout> | null = null;
  let lastValue: T;
  let lastOldValue: T | undefined;

  const invokeCallback = (newValue: T, oldValue: T | undefined) => {
    callback(newValue, oldValue);
    lastCallTime = Date.now();
  };

  const stop = watch(
    source,
    (newValue, oldValue) => {
      const now = Date.now();
      const timeSinceLastCall = now - lastCallTime;

      lastValue = newValue;
      lastOldValue = oldValue;

      if (timeSinceLastCall >= throttle) {
        if (leading || lastCallTime !== 0) {
          invokeCallback(newValue, oldValue);
        } else {
          lastCallTime = now;
        }

        if (timer) {
          clearTimeout(timer);
          timer = null;
        }
      } else if (trailing && !timer) {
        timer = setTimeout(() => {
          invokeCallback(lastValue, lastOldValue);
          timer = null;
        }, throttle - timeSinceLastCall);
      }
    },
    watchOptions
  );

  return () => {
    if (timer) {
      clearTimeout(timer);
    }
    stop();
  };
}

// ============================================================================
// watchPausable - 可暂停的监听
// ============================================================================

/**
 * 可暂停的 watch
 * @example
 * const { pause, resume, isActive, stop } = watchPausable(
 *   () => data.value,
 *   (newData) => {
 *     processData(newData)
 *   }
 * )
 *
 * // 暂停监听
 * pause()
 *
 * // 恢复监听
 * resume()
 */
export function watchPausable<T>(
  source: WatchSource<T>,
  callback: (value: T, oldValue: T | undefined) => void,
  options: WatchOptions = {}
): WatchPausableReturn {
  const isActive = ref(true);

  const stop = watch(
    source,
    (newValue, oldValue) => {
      if (isActive.value) {
        callback(newValue, oldValue);
      }
    },
    options
  );

  const pause = () => {
    isActive.value = false;
  };

  const resume = () => {
    isActive.value = true;
  };

  return {
    pause,
    resume,
    isActive,
    stop,
  };
}

// ============================================================================
// watchIgnorable - 可忽略的监听
// ============================================================================

/**
 * 可忽略更新的 watch
 * @example
 * const count = ref(0)
 *
 * const { ignoreUpdates, stop } = watchIgnorable(
 *   count,
 *   (value) => {
 *     console.log('count changed:', value)
 *   }
 * )
 *
 * // 这次更新不会触发回调
 * ignoreUpdates(() => {
 *   count.value = 100
 * })
 */
export function watchIgnorable<T>(
  source: WatchSource<T>,
  callback: (value: T, oldValue: T | undefined) => void,
  options: WatchOptions = {}
): WatchIgnorableReturn {
  let ignoring = false;

  const stop = watch(
    source,
    (newValue, oldValue) => {
      if (!ignoring) {
        callback(newValue, oldValue);
      }
    },
    options
  );

  const ignoreUpdates = (updater: () => void) => {
    ignoring = true;
    updater();
    nextTick(() => {
      ignoring = false;
    });
  };

  return {
    ignoreUpdates,
    stop,
  };
}

// ============================================================================
// watchTriggered - 触发计数的监听
// ============================================================================

/**
 * 带触发计数的 watch
 * @example
 * const { count, isTriggered, reset, stop } = watchTriggered(
 *   () => data.value,
 *   (newData) => {
 *     processData(newData)
 *   }
 * )
 *
 * console.log(count.value) // 触发次数
 * console.log(isTriggered.value) // 是否至少触发过一次
 * reset() // 重置计数
 */
export function watchTriggered<T>(
  source: WatchSource<T>,
  callback: (value: T, oldValue: T | undefined) => void,
  options: WatchOptions = {}
): WatchTriggeredReturn<T> {
  const count = ref(0);
  const isTriggered = computed(() => count.value > 0);

  const stop = watch(
    source,
    (newValue, oldValue) => {
      count.value++;
      callback(newValue, oldValue);
    },
    options
  );

  const reset = () => {
    count.value = 0;
  };

  return {
    count,
    isTriggered,
    reset,
    stop,
  };
}

// ============================================================================
// watchArray - 数组变化监听
// ============================================================================

export interface WatchArrayReturn<T> {
  /** 停止监听 */
  stop: WatchStopHandle;
}

/**
 * 监听数组变化（增加、删除、更新）
 * @example
 * const items = ref([{ id: 1 }, { id: 2 }])
 *
 * watchArray(items, {
 *   onAdd: (added) => console.log('Added:', added),
 *   onRemove: (removed) => console.log('Removed:', removed),
 *   onUpdate: (updated) => console.log('Updated:', updated)
 * })
 */
export function watchArray<T>(
  source: Ref<T[]>,
  callbacks: {
    onAdd?: (items: T[]) => void;
    onRemove?: (items: T[]) => void;
    onUpdate?: (newArray: T[], oldArray: T[]) => void;
  },
  options: WatchOptions = {}
): WatchArrayReturn<T> {
  const { onAdd, onRemove, onUpdate } = callbacks;

  let prevArray: T[] = [...source.value];

  const stop = watch(
    source,
    (newArray) => {
      const added: T[] = [];
      const removed: T[] = [];

      // 查找新增项
      for (const item of newArray) {
        if (!prevArray.includes(item)) {
          added.push(item);
        }
      }

      // 查找删除项
      for (const item of prevArray) {
        if (!newArray.includes(item)) {
          removed.push(item);
        }
      }

      if (added.length > 0) {
        onAdd?.(added);
      }

      if (removed.length > 0) {
        onRemove?.(removed);
      }

      if (added.length > 0 || removed.length > 0) {
        onUpdate?.(newArray, prevArray);
      }

      prevArray = [...newArray];
    },
    { deep: true, ...options }
  );

  return { stop };
}

// ============================================================================
// watchWithFilter - 带过滤条件的监听
// ============================================================================

/**
 * 带过滤条件的 watch
 * @example
 * watchWithFilter(
 *   () => count.value,
 *   (value) => {
 *     // 只在值为偶数时执行
 *     console.log('Even number:', value)
 *   },
 *   (value) => value % 2 === 0
 * )
 */
export function watchWithFilter<T>(
  source: WatchSource<T>,
  callback: (value: T, oldValue: T | undefined) => void,
  filter: (value: T, oldValue: T | undefined) => boolean,
  options: WatchOptions = {}
): WatchStopHandle {
  return watch(
    source,
    (newValue, oldValue) => {
      if (filter(newValue, oldValue)) {
        callback(newValue, oldValue);
      }
    },
    options
  );
}

// ============================================================================
// watchAtMost - 限制执行次数的监听
// ============================================================================

export interface WatchAtMostReturn {
  /** 已执行次数 */
  count: Ref<number>;
  /** 停止监听 */
  stop: WatchStopHandle;
}

/**
 * 限制执行次数的 watch
 * @example
 * // 最多执行 3 次
 * const { count, stop } = watchAtMost(
 *   () => data.value,
 *   (value) => {
 *     console.log('Executed:', value)
 *   },
 *   3
 * )
 */
export function watchAtMost<T>(
  source: WatchSource<T>,
  callback: (value: T, oldValue: T | undefined) => void,
  limit: number,
  options: WatchOptions = {}
): WatchAtMostReturn {
  const count = ref(0);

  const stop = watch(
    source,
    (newValue, oldValue) => {
      if (count.value < limit) {
        count.value++;
        callback(newValue, oldValue);

        if (count.value >= limit) {
          stop();
        }
      }
    },
    options
  );

  return { count, stop };
}

// ============================================================================
// whenever - 条件满足时执行
// ============================================================================

/**
 * 当条件为真时执行
 * @example
 * const isReady = ref(false)
 *
 * whenever(isReady, () => {
 *   console.log('Ready!')
 * })
 *
 * // 或者只执行一次
 * whenever(isReady, () => {
 *   console.log('Ready!')
 * }, { once: true })
 */
export function whenever<T>(
  source: WatchSource<T>,
  callback: (value: T) => void,
  options: WatchOptions & { once?: boolean } = {}
): WatchStopHandle {
  const { once = false, ...watchOptions } = options;

  if (once) {
    return watchOnce(
      source,
      (value) => {
        if (value) {
          callback(value);
        }
      },
      watchOptions
    );
  }

  return watch(
    source,
    (value, oldValue) => {
      if (value && !oldValue) {
        callback(value);
      }
    },
    watchOptions
  );
}

// ============================================================================
// until - 等待条件满足
// ============================================================================

export interface UntilReturn<T> {
  /** 等待条件为真 */
  toBe: (expected: T) => Promise<T>;
  /** 等待条件为真（布尔） */
  toBeTruthy: () => Promise<T>;
  /** 等待条件不为 null/undefined */
  toBeNotNull: () => Promise<NonNullable<T>>;
  /** 等待满足条件 */
  toMatch: (predicate: (value: T) => boolean) => Promise<T>;
  /** 设置超时 */
  timeout: (ms: number) => UntilReturn<T>;
}

/**
 * 等待条件满足
 * @example
 * const isReady = ref(false)
 *
 * // 等待 isReady 变为 true
 * await until(isReady).toBeTruthy()
 *
 * // 等待特定值
 * await until(status).toBe('completed')
 *
 * // 带超时
 * try {
 *   await until(data).timeout(5000).toBeNotNull()
 * } catch {
 *   console.log('Timeout!')
 * }
 */
export function until<T>(source: WatchSource<T>): UntilReturn<T> {
  let timeoutMs: number | null = null;

  const createPromise = (
    predicate: (value: T) => boolean
  ): Promise<T> => {
    return new Promise((resolve, reject) => {
      let timer: ReturnType<typeof setTimeout> | null = null;

      if (timeoutMs !== null) {
        timer = setTimeout(() => {
          stop();
          reject(new Error("Timeout"));
        }, timeoutMs);
      }

      const stop = watch(
        source,
        (value) => {
          if (predicate(value)) {
            if (timer) clearTimeout(timer);
            stop();
            resolve(value);
          }
        },
        { immediate: true }
      );
    });
  };

  const result: UntilReturn<T> = {
    toBe: (expected: T) => createPromise((value) => value === expected),
    toBeTruthy: () => createPromise((value) => !!value),
    toBeNotNull: () =>
      createPromise((value) => value !== null && value !== undefined) as Promise<
        NonNullable<T>
      >,
    toMatch: (predicate: (value: T) => boolean) => createPromise(predicate),
    timeout: (ms: number) => {
      timeoutMs = ms;
      return result;
    },
  };

  return result;
}

// ============================================================================
// useWatchArray - 响应式监听数组
// ============================================================================

/**
 * 创建可监听的数组
 * @example
 * const { array, push, pop, clear, onChange } = useWatchArray<number>([1, 2, 3])
 *
 * onChange((action, items) => {
 *   console.log(action, items) // 'push', [4] 或 'pop', [3] 等
 * })
 *
 * push(4)
 * pop()
 * clear()
 */
export function useWatchArray<T>(initialValue: T[] = []) {
  type ArrayAction = "push" | "pop" | "shift" | "unshift" | "splice" | "clear" | "set";

  const array = ref<T[]>([...initialValue]) as Ref<T[]>;
  const listeners = new Set<(action: ArrayAction, items: T[]) => void>();

  const notify = (action: ArrayAction, items: T[]) => {
    listeners.forEach((listener) => listener(action, items));
  };

  const push = (...items: T[]) => {
    array.value.push(...items);
    notify("push", items);
  };

  const pop = () => {
    const item = array.value.pop();
    if (item !== undefined) {
      notify("pop", [item]);
    }
    return item;
  };

  const shift = () => {
    const item = array.value.shift();
    if (item !== undefined) {
      notify("shift", [item]);
    }
    return item;
  };

  const unshift = (...items: T[]) => {
    array.value.unshift(...items);
    notify("unshift", items);
  };

  const splice = (start: number, deleteCount?: number, ...items: T[]) => {
    const removed = array.value.splice(start, deleteCount ?? 0, ...items);
    notify("splice", removed);
    return removed;
  };

  const clear = () => {
    const removed = [...array.value];
    array.value = [];
    notify("clear", removed);
  };

  const set = (index: number, value: T) => {
    const old = array.value[index];
    array.value[index] = value;
    notify("set", [old]);
  };

  const onChange = (callback: (action: ArrayAction, items: T[]) => void) => {
    listeners.add(callback);
    return () => listeners.delete(callback);
  };

  return {
    array,
    push,
    pop,
    shift,
    unshift,
    splice,
    clear,
    set,
    onChange,
  };
}
