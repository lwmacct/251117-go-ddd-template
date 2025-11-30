/**
 * 异步状态 Composable
 * 管理异步操作的 loading、error、data 状态
 */

import { ref, shallowRef, computed, type Ref, type ComputedRef } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UseAsyncStateOptions<T> {
  /** 初始数据 */
  initialData?: T;
  /** 是否立即执行，默认 false */
  immediate?: boolean;
  /** 是否在执行前重置状态，默认 true */
  resetOnExecute?: boolean;
  /** 是否使用 shallowRef，默认 false */
  shallow?: boolean;
  /** 错误处理器 */
  onError?: (error: Error) => void;
  /** 成功处理器 */
  onSuccess?: (data: T) => void;
  /** 完成处理器（无论成功或失败） */
  onFinally?: () => void;
  /** 延迟显示 loading，默认 0 */
  delay?: number;
}

export interface UseAsyncStateReturn<T, P extends unknown[]> {
  /** 数据 */
  data: Ref<T | undefined>;
  /** 是否正在加载 */
  isLoading: Ref<boolean>;
  /** 是否已完成（成功或失败） */
  isFinished: Ref<boolean>;
  /** 是否成功 */
  isSuccess: Ref<boolean>;
  /** 错误信息 */
  error: Ref<Error | null>;
  /** 执行异步操作 */
  execute: (...args: P) => Promise<T>;
  /** 重置状态 */
  reset: () => void;
  /** 状态 */
  state: ComputedRef<"idle" | "loading" | "success" | "error">;
}

// ============================================================================
// 主函数
// ============================================================================

/**
 * 异步状态管理
 * @example
 * const { data, isLoading, error, execute } = useAsyncState(
 *   (id: number) => fetchUser(id),
 *   { initialData: null }
 * )
 *
 * // 执行
 * await execute(1)
 */
export function useAsyncState<T, P extends unknown[] = []>(
  fn: (...args: P) => Promise<T>,
  options: UseAsyncStateOptions<T> = {}
): UseAsyncStateReturn<T, P> {
  const {
    initialData,
    immediate = false,
    resetOnExecute = true,
    shallow = false,
    onError,
    onSuccess,
    onFinally,
    delay = 0,
  } = options;

  // 状态
  const data = (shallow ? shallowRef : ref)(initialData) as Ref<T | undefined>;
  const isLoading = ref(false);
  const isFinished = ref(false);
  const isSuccess = ref(false);
  const error = ref<Error | null>(null);

  // 计算状态
  const state = computed<"idle" | "loading" | "success" | "error">(() => {
    if (isLoading.value) return "loading";
    if (error.value) return "error";
    if (isSuccess.value) return "success";
    return "idle";
  });

  // 重置状态
  const reset = () => {
    data.value = initialData;
    isLoading.value = false;
    isFinished.value = false;
    isSuccess.value = false;
    error.value = null;
  };

  // 延迟定时器
  let delayTimer: ReturnType<typeof setTimeout> | null = null;

  // 执行异步操作
  const execute = async (...args: P): Promise<T> => {
    // 重置状态
    if (resetOnExecute) {
      error.value = null;
      isSuccess.value = false;
    }

    // 延迟显示 loading
    if (delay > 0) {
      delayTimer = setTimeout(() => {
        isLoading.value = true;
      }, delay);
    } else {
      isLoading.value = true;
    }

    isFinished.value = false;

    try {
      const result = await fn(...args);
      data.value = result;
      isSuccess.value = true;
      onSuccess?.(result);
      return result;
    } catch (e) {
      const err = e instanceof Error ? e : new Error(String(e));
      error.value = err;
      onError?.(err);
      throw err;
    } finally {
      // 清除延迟定时器
      if (delayTimer) {
        clearTimeout(delayTimer);
        delayTimer = null;
      }
      isLoading.value = false;
      isFinished.value = true;
      onFinally?.();
    }
  };

  // 立即执行
  if (immediate) {
    execute(...([] as unknown as P));
  }

  return {
    data,
    isLoading,
    isFinished,
    isSuccess,
    error,
    execute,
    reset,
    state,
  };
}

// ============================================================================
// 可重试的异步状态
// ============================================================================

export interface UseAsyncRetryOptions<T> extends UseAsyncStateOptions<T> {
  /** 最大重试次数，默认 3 */
  maxRetries?: number;
  /** 重试延迟（毫秒），默认 1000 */
  retryDelay?: number;
  /** 重试延迟增长因子，默认 2 */
  retryDelayFactor?: number;
  /** 是否重试的判断函数 */
  shouldRetry?: (error: Error, retryCount: number) => boolean;
}

/**
 * 可重试的异步状态
 * @example
 * const { data, execute, retryCount } = useAsyncRetry(
 *   () => fetchData(),
 *   { maxRetries: 3 }
 * )
 */
export function useAsyncRetry<T, P extends unknown[] = []>(
  fn: (...args: P) => Promise<T>,
  options: UseAsyncRetryOptions<T> = {}
) {
  const {
    maxRetries = 3,
    retryDelay = 1000,
    retryDelayFactor = 2,
    shouldRetry = () => true,
    ...asyncOptions
  } = options;

  const retryCount = ref(0);

  const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

  const wrappedFn = async (...args: P): Promise<T> => {
    retryCount.value = 0;
    let lastError: Error | null = null;
    let currentDelay = retryDelay;

    while (retryCount.value <= maxRetries) {
      try {
        return await fn(...args);
      } catch (e) {
        const err = e instanceof Error ? e : new Error(String(e));
        lastError = err;

        if (retryCount.value < maxRetries && shouldRetry(err, retryCount.value)) {
          retryCount.value++;
          await sleep(currentDelay);
          currentDelay *= retryDelayFactor;
        } else {
          break;
        }
      }
    }

    throw lastError;
  };

  const asyncState = useAsyncState(wrappedFn, asyncOptions);

  return {
    ...asyncState,
    retryCount,
  };
}

// ============================================================================
// 轮询异步状态
// ============================================================================

export interface UsePollingOptions<T> extends UseAsyncStateOptions<T> {
  /** 轮询间隔（毫秒），默认 5000 */
  interval?: number;
  /** 是否立即开始轮询，默认 false */
  immediate?: boolean;
  /** 是否在页面不可见时暂停，默认 true */
  pauseOnHidden?: boolean;
}

/**
 * 轮询异步状态
 * @example
 * const { data, isPolling, start, stop } = usePolling(
 *   () => fetchStatus(),
 *   { interval: 3000 }
 * )
 */
export function usePolling<T, P extends unknown[] = []>(
  fn: (...args: P) => Promise<T>,
  options: UsePollingOptions<T> = {}
) {
  const { interval = 5000, immediate = false, pauseOnHidden = true, ...asyncOptions } = options;

  const isPolling = ref(false);
  let timer: ReturnType<typeof setInterval> | null = null;
  let savedArgs: P | null = null;

  const asyncState = useAsyncState(fn, asyncOptions);

  const stop = () => {
    isPolling.value = false;
    if (timer) {
      clearInterval(timer);
      timer = null;
    }
  };

  const start = (...args: P) => {
    savedArgs = args;
    stop();
    isPolling.value = true;

    // 立即执行一次
    asyncState.execute(...args);

    // 设置轮询
    timer = setInterval(() => {
      if (savedArgs) {
        asyncState.execute(...savedArgs);
      }
    }, interval);
  };

  // 页面可见性变化处理
  if (pauseOnHidden && typeof document !== "undefined") {
    const handleVisibilityChange = () => {
      if (document.hidden) {
        if (timer) {
          clearInterval(timer);
          timer = null;
        }
      } else if (isPolling.value && savedArgs) {
        // 恢复轮询
        asyncState.execute(...savedArgs);
        timer = setInterval(() => {
          if (savedArgs) {
            asyncState.execute(...savedArgs);
          }
        }, interval);
      }
    };

    document.addEventListener("visibilitychange", handleVisibilityChange);
  }

  // 立即开始
  if (immediate) {
    start(...([] as unknown as P));
  }

  return {
    ...asyncState,
    isPolling,
    start,
    stop,
  };
}

// ============================================================================
// 简单的异步 ref
// ============================================================================

/**
 * 异步计算的 ref
 * @example
 * const user = useAsyncRef(() => fetchUser(userId.value), {
 *   watch: [userId]
 * })
 */
export function useAsyncRef<T>(fn: () => Promise<T>, options: UseAsyncStateOptions<T> & { watch?: unknown[] } = {}) {
  const { watch: watchSources, immediate = true, ...asyncOptions } = options;

  const asyncState = useAsyncState(fn, { ...asyncOptions, immediate });

  // 监听依赖变化
  if (watchSources && watchSources.length > 0) {
    const { watch } = await import("vue");
    watch(watchSources, () => {
      asyncState.execute();
    });
  }

  return asyncState;
}

// ============================================================================
// Promise 队列
// ============================================================================

/**
 * Promise 队列，按顺序执行
 * @example
 * const queue = usePromiseQueue()
 * queue.add(() => fetchData1())
 * queue.add(() => fetchData2())
 */
export function usePromiseQueue() {
  const pending = ref<Promise<unknown>[]>([]);
  const isProcessing = ref(false);

  let chain = Promise.resolve();

  const add = <T>(fn: () => Promise<T>): Promise<T> => {
    return new Promise<T>((resolve, reject) => {
      chain = chain.then(async () => {
        isProcessing.value = true;
        try {
          const result = await fn();
          resolve(result);
        } catch (e) {
          reject(e);
        } finally {
          isProcessing.value = pending.value.length > 0;
        }
      });
    });
  };

  const clear = () => {
    chain = Promise.resolve();
    pending.value = [];
    isProcessing.value = false;
  };

  return {
    add,
    clear,
    isProcessing,
    pending,
  };
}
