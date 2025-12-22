/**
 * 限流工具
 * 提供节流、防抖、限速等功能
 */

// ============================================================================
// 类型定义
// ============================================================================

export interface ThrottleOptions {
  /** 是否在开始时立即调用，默认 true */
  leading?: boolean;
  /** 是否在结束时调用，默认 true */
  trailing?: boolean;
}

export interface DebounceOptions {
  /** 最大等待时间 */
  maxWait?: number;
  /** 是否在开始时调用 */
  leading?: boolean;
  /** 是否在结束时调用，默认 true */
  trailing?: boolean;
}

export interface ThrottledFunction<T extends (...args: unknown[]) => unknown> {
  (...args: Parameters<T>): ReturnType<T> | undefined;
  /** 取消待执行的调用 */
  cancel: () => void;
  /** 立即执行 */
  flush: () => ReturnType<T> | undefined;
  /** 检查是否有待执行的调用 */
  pending: () => boolean;
}

// ============================================================================
// 节流
// ============================================================================

/**
 * 节流函数
 * 在指定时间内只执行一次
 * @example
 * const throttledScroll = throttle(() => {
 *   console.log('scroll')
 * }, 200)
 *
 * window.addEventListener('scroll', throttledScroll)
 */
export function throttle<T extends (...args: unknown[]) => unknown>(
  fn: T,
  wait: number,
  options: ThrottleOptions = {}
): ThrottledFunction<T> {
  const { leading = true, trailing = true } = options;

  let lastCallTime: number | null = null;
  let lastInvokeTime = 0;
  let timerId: ReturnType<typeof setTimeout> | null = null;
  let lastArgs: Parameters<T> | null = null;
  let lastThis: unknown = null;
  let result: ReturnType<T> | undefined;

  // 获取当前时间
  const now = () => Date.now();

  // 调用函数
  const invokeFunc = (time: number): ReturnType<T> | undefined => {
    const args = lastArgs!;
    const thisArg = lastThis;

    lastArgs = null;
    lastThis = null;
    lastInvokeTime = time;
    result = fn.apply(thisArg, args) as ReturnType<T>;
    return result;
  };

  // 计算剩余等待时间
  const remainingWait = (time: number): number => {
    const timeSinceLastCall = time - (lastCallTime ?? 0);
    const _timeSinceLastInvoke = time - lastInvokeTime;
    const timeWaiting = wait - timeSinceLastCall;

    return Math.max(0, timeWaiting);
  };

  // 是否应该调用
  const shouldInvoke = (time: number): boolean => {
    const timeSinceLastCall = time - (lastCallTime ?? 0);
    const timeSinceLastInvoke = time - lastInvokeTime;

    return lastCallTime === null || timeSinceLastCall >= wait || timeSinceLastCall < 0 || timeSinceLastInvoke >= wait;
  };

  // 定时器到期处理
  const timerExpired = (): ReturnType<T> | undefined => {
    const time = now();
    if (shouldInvoke(time)) {
      return trailingEdge(time);
    }
    timerId = setTimeout(timerExpired, remainingWait(time));
    return undefined;
  };

  // 开始边缘调用
  const leadingEdge = (time: number): ReturnType<T> | undefined => {
    lastInvokeTime = time;
    timerId = setTimeout(timerExpired, wait);
    return leading ? invokeFunc(time) : result;
  };

  // 结束边缘调用
  const trailingEdge = (time: number): ReturnType<T> | undefined => {
    timerId = null;

    if (trailing && lastArgs) {
      return invokeFunc(time);
    }
    lastArgs = null;
    lastThis = null;
    return result;
  };

  // 取消
  const cancel = () => {
    if (timerId !== null) {
      clearTimeout(timerId);
    }
    lastInvokeTime = 0;
    lastArgs = null;
    lastCallTime = null;
    lastThis = null;
    timerId = null;
  };

  // 立即执行
  const flush = (): ReturnType<T> | undefined => {
    if (timerId === null) {
      return result;
    }
    return trailingEdge(now());
  };

  // 检查是否有待执行
  const pending = (): boolean => {
    return timerId !== null;
  };

  // 主函数
  const throttled = function (this: unknown, ...args: Parameters<T>): ReturnType<T> | undefined {
    const time = now();
    const isInvoking = shouldInvoke(time);

    lastArgs = args;
    lastThis = this;
    lastCallTime = time;

    if (isInvoking) {
      if (timerId === null) {
        return leadingEdge(lastCallTime);
      }
    }

    if (timerId === null) {
      timerId = setTimeout(timerExpired, wait);
    }

    return result;
  } as ThrottledFunction<T>;

  throttled.cancel = cancel;
  throttled.flush = flush;
  throttled.pending = pending;

  return throttled;
}

// ============================================================================
// 防抖
// ============================================================================

/**
 * 防抖函数
 * 在停止调用后等待指定时间才执行
 * @example
 * const debouncedSearch = debounce((query: string) => {
 *   search(query)
 * }, 300)
 *
 * input.addEventListener('input', (e) => {
 *   debouncedSearch(e.target.value)
 * })
 */
export function debounce<T extends (...args: unknown[]) => unknown>(
  fn: T,
  wait: number,
  options: DebounceOptions = {}
): ThrottledFunction<T> {
  const { maxWait, leading = false, trailing = true } = options;

  let lastArgs: Parameters<T> | null = null;
  let lastThis: unknown = null;
  let result: ReturnType<T> | undefined;
  let timerId: ReturnType<typeof setTimeout> | null = null;
  let lastCallTime: number | undefined;
  let lastInvokeTime = 0;
  const maxing = maxWait !== undefined;
  let maxTimerId: ReturnType<typeof setTimeout> | null = null;

  const now = () => Date.now();

  const invokeFunc = (time: number): ReturnType<T> | undefined => {
    const args = lastArgs!;
    const thisArg = lastThis;

    lastArgs = null;
    lastThis = null;
    lastInvokeTime = time;
    result = fn.apply(thisArg, args) as ReturnType<T>;
    return result;
  };

  const startTimer = (pendingFunc: () => void, waitTime: number) => {
    return setTimeout(pendingFunc, waitTime);
  };

  const cancelTimer = (id: ReturnType<typeof setTimeout> | null) => {
    if (id !== null) {
      clearTimeout(id);
    }
  };

  const leadingEdge = (time: number): ReturnType<T> | undefined => {
    lastInvokeTime = time;
    timerId = startTimer(timerExpired, wait);

    if (maxing) {
      maxTimerId = startTimer(() => {
        if (lastArgs) {
          invokeFunc(now());
        }
      }, maxWait!);
    }

    return leading ? invokeFunc(time) : result;
  };

  const remainingWait = (time: number): number => {
    const timeSinceLastCall = time - (lastCallTime ?? 0);
    const timeSinceLastInvoke = time - lastInvokeTime;
    const timeWaiting = wait - timeSinceLastCall;

    return maxing ? Math.min(timeWaiting, maxWait! - timeSinceLastInvoke) : timeWaiting;
  };

  const shouldInvoke = (time: number): boolean => {
    const timeSinceLastCall = time - (lastCallTime ?? 0);
    const timeSinceLastInvoke = time - lastInvokeTime;

    return (
      lastCallTime === undefined ||
      timeSinceLastCall >= wait ||
      timeSinceLastCall < 0 ||
      (maxing && timeSinceLastInvoke >= maxWait!)
    );
  };

  const timerExpired = (): ReturnType<T> | undefined => {
    const time = now();
    if (shouldInvoke(time)) {
      return trailingEdge(time);
    }
    timerId = startTimer(timerExpired, remainingWait(time));
    return undefined;
  };

  const trailingEdge = (time: number): ReturnType<T> | undefined => {
    timerId = null;
    cancelTimer(maxTimerId);
    maxTimerId = null;

    if (trailing && lastArgs) {
      return invokeFunc(time);
    }

    lastArgs = null;
    lastThis = null;
    return result;
  };

  const cancel = () => {
    cancelTimer(timerId);
    cancelTimer(maxTimerId);
    lastInvokeTime = 0;
    lastArgs = null;
    lastCallTime = undefined;
    lastThis = null;
    timerId = null;
    maxTimerId = null;
  };

  const flush = (): ReturnType<T> | undefined => {
    if (timerId === null) {
      return result;
    }
    return trailingEdge(now());
  };

  const pending = (): boolean => {
    return timerId !== null;
  };

  const debounced = function (this: unknown, ...args: Parameters<T>): ReturnType<T> | undefined {
    const time = now();
    const isInvoking = shouldInvoke(time);

    lastArgs = args;
    lastThis = this;
    lastCallTime = time;

    if (isInvoking) {
      if (timerId === null) {
        return leadingEdge(lastCallTime);
      }
      if (maxing) {
        cancelTimer(timerId);
        timerId = startTimer(timerExpired, wait);
        return invokeFunc(lastCallTime);
      }
    }

    if (timerId === null) {
      timerId = startTimer(timerExpired, wait);
    }

    return result;
  } as ThrottledFunction<T>;

  debounced.cancel = cancel;
  debounced.flush = flush;
  debounced.pending = pending;

  return debounced;
}

// ============================================================================
// 限速器
// ============================================================================

/**
 * 限速器
 * 限制每秒/每分钟的调用次数
 * @example
 * const limiter = createRateLimiter({
 *   maxRequests: 10,
 *   interval: 1000 // 每秒最多 10 次
 * })
 *
 * limiter.execute(() => callApi())
 */
export function createRateLimiter(options: { maxRequests: number; interval: number }) {
  const { maxRequests, interval } = options;
  const queue: Array<{
    fn: () => Promise<unknown>;
    resolve: (value: unknown) => void;
    reject: (error: unknown) => void;
  }> = [];
  let currentRequests = 0;
  let isProcessing = false;

  const processQueue = async () => {
    if (isProcessing) return;
    isProcessing = true;

    while (queue.length > 0) {
      if (currentRequests >= maxRequests) {
        await new Promise((resolve) => setTimeout(resolve, interval));
        currentRequests = 0;
      }

      const item = queue.shift();
      if (item) {
        currentRequests++;
        try {
          const result = await item.fn();
          item.resolve(result);
        } catch (error) {
          item.reject(error);
        }
      }
    }

    isProcessing = false;
  };

  const execute = <T>(fn: () => Promise<T>): Promise<T> => {
    return new Promise((resolve, reject) => {
      queue.push({ fn, resolve: resolve as (value: unknown) => void, reject });
      processQueue();
    });
  };

  const clear = () => {
    queue.length = 0;
  };

  return {
    execute,
    clear,
    get pending() {
      return queue.length;
    },
  };
}

// ============================================================================
// 重试
// ============================================================================

export interface RetryOptions {
  /** 最大重试次数，默认 3 */
  maxRetries?: number;
  /** 重试延迟（毫秒），默认 1000 */
  delay?: number;
  /** 延迟增长因子（指数退避），默认 2 */
  factor?: number;
  /** 是否重试的判断函数 */
  shouldRetry?: (error: unknown, attempt: number) => boolean;
  /** 重试前回调 */
  onRetry?: (error: unknown, attempt: number) => void;
}

/**
 * 重试函数
 * @example
 * const result = await retry(
 *   () => fetchData(),
 *   { maxRetries: 3, delay: 1000 }
 * )
 */
export async function retry<T>(fn: () => Promise<T>, options: RetryOptions = {}): Promise<T> {
  const { maxRetries = 3, delay = 1000, factor = 2, shouldRetry = () => true, onRetry } = options;

  let lastError: unknown;
  let currentDelay = delay;

  for (let attempt = 0; attempt <= maxRetries; attempt++) {
    try {
      return await fn();
    } catch (error) {
      lastError = error;

      if (attempt < maxRetries && shouldRetry(error, attempt)) {
        onRetry?.(error, attempt);
        await new Promise((resolve) => setTimeout(resolve, currentDelay));
        currentDelay *= factor;
      }
    }
  }

  throw lastError;
}

// ============================================================================
// 超时
// ============================================================================

/**
 * 超时包装
 * @example
 * const result = await withTimeout(
 *   fetchData(),
 *   5000,
 *   '请求超时'
 * )
 */
export async function withTimeout<T>(
  promise: Promise<T>,
  timeout: number,
  message = "Operation timed out"
): Promise<T> {
  let timeoutId: ReturnType<typeof setTimeout>;

  const timeoutPromise = new Promise<never>((_, reject) => {
    timeoutId = setTimeout(() => {
      reject(new Error(message));
    }, timeout);
  });

  try {
    return await Promise.race([promise, timeoutPromise]);
  } finally {
    clearTimeout(timeoutId!);
  }
}

// ============================================================================
// 去重执行
// ============================================================================

/**
 * 创建去重执行器
 * 相同 key 的调用只执行一次
 * @example
 * const dedup = createDeduplicator()
 *
 * // 这两次调用只会执行一次 fetchUser(1)
 * dedup.execute('user-1', () => fetchUser(1))
 * dedup.execute('user-1', () => fetchUser(1))
 */
export function createDeduplicator() {
  const pending = new Map<string, Promise<unknown>>();

  const execute = <T>(key: string, fn: () => Promise<T>): Promise<T> => {
    if (pending.has(key)) {
      return pending.get(key) as Promise<T>;
    }

    const promise = fn().finally(() => {
      pending.delete(key);
    });

    pending.set(key, promise);
    return promise;
  };

  const clear = () => {
    pending.clear();
  };

  return {
    execute,
    clear,
    get size() {
      return pending.size;
    },
  };
}
