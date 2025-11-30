/**
 * SWR (Stale-While-Revalidate) Composable
 * 提供数据缓存和重新验证的策略
 */

import {
  ref,
  computed,
  watch,
  onMounted,
  onUnmounted,
  type Ref,
  type ComputedRef,
} from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export type SWRStatus = "idle" | "loading" | "validating" | "success" | "error";

export interface UseSWROptions<T> {
  /** 初始数据 */
  initialData?: T;
  /** 是否立即获取 */
  immediate?: boolean;
  /** 重新验证间隔（毫秒），0 表示禁用 */
  revalidateInterval?: number;
  /** 是否在聚焦时重新验证 */
  revalidateOnFocus?: boolean;
  /** 是否在重新连接时重新验证 */
  revalidateOnReconnect?: boolean;
  /** 是否在挂载时重新验证 */
  revalidateOnMount?: boolean;
  /** 去重间隔（毫秒） */
  dedupingInterval?: number;
  /** 错误重试次数 */
  errorRetryCount?: number;
  /** 错误重试间隔（毫秒） */
  errorRetryInterval?: number;
  /** 数据变化回调 */
  onSuccess?: (data: T) => void;
  /** 错误回调 */
  onError?: (error: Error) => void;
  /** 比较函数 */
  compare?: (a: T | null, b: T | null) => boolean;
  /** 是否保持之前的数据（加载时） */
  keepPreviousData?: boolean;
}

export interface UseSWRReturn<T> {
  /** 数据 */
  data: Ref<T | null>;
  /** 错误 */
  error: Ref<Error | null>;
  /** 状态 */
  status: Ref<SWRStatus>;
  /** 是否正在验证 */
  isValidating: Ref<boolean>;
  /** 是否正在加载（首次） */
  isLoading: ComputedRef<boolean>;
  /** 重新验证 */
  mutate: (data?: T | ((prev: T | null) => T)) => Promise<void>;
  /** 重新获取 */
  revalidate: () => Promise<void>;
}

// 全局缓存
const swrCache = new Map<
  string,
  {
    data: unknown;
    error: Error | null;
    timestamp: number;
    promise?: Promise<unknown>;
  }
>();

// 订阅者
const subscribers = new Map<string, Set<() => void>>();

// ============================================================================
// useSWR - Stale-While-Revalidate
// ============================================================================

/**
 * 实现 SWR 数据获取策略
 * @example
 * const { data, error, isLoading, revalidate, mutate } = useSWR(
 *   'users',
 *   () => fetch('/api/users').then(r => r.json()),
 *   {
 *     revalidateInterval: 30000,
 *     revalidateOnFocus: true
 *   }
 * )
 */
export function useSWR<T>(
  key: string,
  fetcher: () => Promise<T>,
  options: UseSWROptions<T> = {}
): UseSWRReturn<T> {
  const {
    initialData,
    immediate = true,
    revalidateInterval = 0,
    revalidateOnFocus = true,
    revalidateOnReconnect = true,
    revalidateOnMount = true,
    dedupingInterval = 2000,
    errorRetryCount = 3,
    errorRetryInterval = 5000,
    onSuccess,
    onError,
    compare = (a, b) => JSON.stringify(a) === JSON.stringify(b),
    keepPreviousData = false,
  } = options;

  const data = ref<T | null>(initialData ?? null) as Ref<T | null>;
  const error = ref<Error | null>(null);
  const status = ref<SWRStatus>("idle");
  const isValidating = ref(false);

  const isLoading = computed(
    () => status.value === "loading" && data.value === null
  );

  let retryCount = 0;
  let intervalTimer: ReturnType<typeof setInterval> | null = null;

  // 从缓存加载
  const loadFromCache = () => {
    const cached = swrCache.get(key);
    if (cached) {
      if (cached.data !== undefined) {
        data.value = cached.data as T;
      }
      if (cached.error) {
        error.value = cached.error;
      }
    }
  };

  // 更新缓存
  const updateCache = (newData: T | null, newError: Error | null = null) => {
    swrCache.set(key, {
      data: newData,
      error: newError,
      timestamp: Date.now(),
    });

    // 通知订阅者
    const subs = subscribers.get(key);
    if (subs) {
      subs.forEach((callback) => callback());
    }
  };

  // 订阅更新
  const subscribe = (callback: () => void) => {
    if (!subscribers.has(key)) {
      subscribers.set(key, new Set());
    }
    subscribers.get(key)!.add(callback);

    return () => {
      subscribers.get(key)?.delete(callback);
    };
  };

  // 获取数据
  const fetchData = async () => {
    const cached = swrCache.get(key);
    const now = Date.now();

    // 去重：如果最近刚获取过，跳过
    if (
      cached &&
      cached.promise &&
      now - cached.timestamp < dedupingInterval
    ) {
      return cached.promise;
    }

    // 如果有缓存数据，先显示缓存
    if (cached?.data !== undefined) {
      data.value = cached.data as T;
    }

    isValidating.value = true;
    if (data.value === null && !keepPreviousData) {
      status.value = "loading";
    } else {
      status.value = "validating";
    }

    try {
      const promise = fetcher();

      // 保存 promise 用于去重
      swrCache.set(key, {
        ...cached,
        data: cached?.data,
        error: null,
        timestamp: now,
        promise,
      });

      const result = await promise;

      // 检查数据是否变化
      if (!compare(data.value, result)) {
        data.value = result;
        onSuccess?.(result);
      }

      error.value = null;
      status.value = "success";
      retryCount = 0;

      updateCache(result);

      return result;
    } catch (err) {
      const fetchError =
        err instanceof Error ? err : new Error(String(err));
      error.value = fetchError;
      status.value = "error";
      onError?.(fetchError);

      updateCache(data.value, fetchError);

      // 错误重试
      if (retryCount < errorRetryCount) {
        retryCount++;
        setTimeout(() => {
          fetchData();
        }, errorRetryInterval * retryCount);
      }

      throw fetchError;
    } finally {
      isValidating.value = false;
    }
  };

  // 重新验证
  const revalidate = async () => {
    try {
      await fetchData();
    } catch {
      // 错误已处理
    }
  };

  // 手动更新数据
  const mutate = async (
    newData?: T | ((prev: T | null) => T)
  ): Promise<void> => {
    if (newData === undefined) {
      // 不传参数则重新获取
      await revalidate();
      return;
    }

    // 乐观更新
    const resolvedData =
      typeof newData === "function"
        ? (newData as (prev: T | null) => T)(data.value)
        : newData;

    data.value = resolvedData;
    updateCache(resolvedData);

    // 可选：重新验证
    // await revalidate();
  };

  // 处理焦点事件
  const handleFocus = () => {
    if (revalidateOnFocus && document.visibilityState === "visible") {
      revalidate();
    }
  };

  // 处理网络重连
  const handleOnline = () => {
    if (revalidateOnReconnect) {
      revalidate();
    }
  };

  // 订阅其他组件的更新
  const handleCacheUpdate = () => {
    loadFromCache();
  };

  onMounted(() => {
    // 加载缓存
    loadFromCache();

    // 订阅缓存更新
    const unsubscribe = subscribe(handleCacheUpdate);

    // 初始获取
    if (immediate && revalidateOnMount) {
      revalidate();
    }

    // 设置自动重新验证
    if (revalidateInterval > 0) {
      intervalTimer = setInterval(revalidate, revalidateInterval);
    }

    // 监听焦点
    if (revalidateOnFocus) {
      document.addEventListener("visibilitychange", handleFocus);
      window.addEventListener("focus", handleFocus);
    }

    // 监听网络
    if (revalidateOnReconnect) {
      window.addEventListener("online", handleOnline);
    }

    // 清理
    onUnmounted(() => {
      unsubscribe();

      if (intervalTimer) {
        clearInterval(intervalTimer);
      }

      if (revalidateOnFocus) {
        document.removeEventListener("visibilitychange", handleFocus);
        window.removeEventListener("focus", handleFocus);
      }

      if (revalidateOnReconnect) {
        window.removeEventListener("online", handleOnline);
      }
    });
  });

  return {
    data,
    error,
    status,
    isValidating,
    isLoading,
    mutate,
    revalidate,
  };
}

// ============================================================================
// useSWRMutation - 用于修改数据的 SWR
// ============================================================================

export interface UseSWRMutationOptions<T, A> {
  /** 乐观更新函数 */
  optimisticData?: (currentData: T | null, arg: A) => T;
  /** 回滚函数 */
  rollbackOnError?: boolean;
  /** 成功回调 */
  onSuccess?: (data: T, arg: A) => void;
  /** 错误回调 */
  onError?: (error: Error, arg: A) => void;
}

export interface UseSWRMutationReturn<T, A> {
  /** 数据 */
  data: Ref<T | null>;
  /** 错误 */
  error: Ref<Error | null>;
  /** 是否正在执行 */
  isMutating: Ref<boolean>;
  /** 触发修改 */
  trigger: (arg: A) => Promise<T>;
  /** 重置状态 */
  reset: () => void;
}

/**
 * 用于修改数据的 SWR
 * @example
 * const { trigger, isMutating, error } = useSWRMutation(
 *   'users',
 *   (user: User) => fetch('/api/users', {
 *     method: 'POST',
 *     body: JSON.stringify(user)
 *   }).then(r => r.json()),
 *   {
 *     onSuccess: (data) => {
 *       // 更新列表缓存
 *       mutate('users')
 *     }
 *   }
 * )
 *
 * // 触发创建
 * await trigger({ name: 'John' })
 */
export function useSWRMutation<T, A = void>(
  key: string,
  mutator: (arg: A) => Promise<T>,
  options: UseSWRMutationOptions<T, A> = {}
): UseSWRMutationReturn<T, A> {
  const { optimisticData, rollbackOnError = true, onSuccess, onError } = options;

  const data = ref<T | null>(null) as Ref<T | null>;
  const error = ref<Error | null>(null);
  const isMutating = ref(false);

  let previousData: T | null = null;

  const trigger = async (arg: A): Promise<T> => {
    const cached = swrCache.get(key);
    previousData = cached?.data as T | null;

    error.value = null;
    isMutating.value = true;

    // 乐观更新
    if (optimisticData && previousData !== null) {
      const optimistic = optimisticData(previousData, arg);
      swrCache.set(key, {
        data: optimistic,
        error: null,
        timestamp: Date.now(),
      });

      // 通知订阅者
      const subs = subscribers.get(key);
      if (subs) {
        subs.forEach((callback) => callback());
      }
    }

    try {
      const result = await mutator(arg);
      data.value = result;

      // 更新缓存
      swrCache.set(key, {
        data: result,
        error: null,
        timestamp: Date.now(),
      });

      // 通知订阅者
      const subs = subscribers.get(key);
      if (subs) {
        subs.forEach((callback) => callback());
      }

      onSuccess?.(result, arg);
      return result;
    } catch (err) {
      const mutationError =
        err instanceof Error ? err : new Error(String(err));
      error.value = mutationError;

      // 回滚
      if (rollbackOnError && previousData !== null) {
        swrCache.set(key, {
          data: previousData,
          error: null,
          timestamp: Date.now(),
        });

        // 通知订阅者
        const subs = subscribers.get(key);
        if (subs) {
          subs.forEach((callback) => callback());
        }
      }

      onError?.(mutationError, arg);
      throw mutationError;
    } finally {
      isMutating.value = false;
    }
  };

  const reset = () => {
    data.value = null;
    error.value = null;
    isMutating.value = false;
  };

  return {
    data,
    error,
    isMutating,
    trigger,
    reset,
  };
}

// ============================================================================
// useSWRInfinite - 无限加载的 SWR
// ============================================================================

export interface UseSWRInfiniteOptions<T> extends Omit<UseSWROptions<T>, "initialData"> {
  /** 初始页数 */
  initialSize?: number;
  /** 每页数据量 */
  pageSize?: number;
  /** 判断是否还有更多数据 */
  hasMore?: (data: T[], pageIndex: number) => boolean;
}

export interface UseSWRInfiniteReturn<T> {
  /** 所有页的数据 */
  data: Ref<T[]>;
  /** 错误 */
  error: Ref<Error | null>;
  /** 是否正在验证 */
  isValidating: Ref<boolean>;
  /** 是否正在加载更多 */
  isLoadingMore: Ref<boolean>;
  /** 是否有更多数据 */
  hasMore: Ref<boolean>;
  /** 当前页数 */
  size: Ref<number>;
  /** 设置页数 */
  setSize: (size: number | ((prev: number) => number)) => void;
  /** 加载更多 */
  loadMore: () => Promise<void>;
  /** 重新验证 */
  revalidate: () => Promise<void>;
  /** 重置 */
  reset: () => void;
}

/**
 * 无限加载的 SWR
 * @example
 * const { data, loadMore, hasMore, isLoadingMore } = useSWRInfinite(
 *   (pageIndex) => `users-page-${pageIndex}`,
 *   (pageIndex) => fetch(`/api/users?page=${pageIndex}`).then(r => r.json()),
 *   {
 *     pageSize: 10,
 *     hasMore: (data, pageIndex) => data.length === 10
 *   }
 * )
 */
export function useSWRInfinite<T>(
  getKey: (pageIndex: number) => string,
  fetcher: (pageIndex: number) => Promise<T[]>,
  options: UseSWRInfiniteOptions<T> = {}
): UseSWRInfiniteReturn<T> {
  const {
    initialSize = 1,
    pageSize = 10,
    hasMore: hasMoreFn = (data, _) => data.length === pageSize,
    ...swrOptions
  } = options;

  const data = ref<T[]>([]) as Ref<T[]>;
  const error = ref<Error | null>(null);
  const isValidating = ref(false);
  const isLoadingMore = ref(false);
  const hasMore = ref(true);
  const size = ref(initialSize);

  // 获取所有页的数据
  const fetchAllPages = async () => {
    isValidating.value = true;
    error.value = null;

    try {
      const pages: T[][] = [];

      for (let i = 0; i < size.value; i++) {
        const key = getKey(i);
        const cached = swrCache.get(key);

        if (cached?.data && !swrOptions.revalidateOnMount) {
          pages.push(cached.data as T[]);
        } else {
          const pageData = await fetcher(i);
          pages.push(pageData);
          swrCache.set(key, {
            data: pageData,
            error: null,
            timestamp: Date.now(),
          });
        }
      }

      const allData = pages.flat();
      data.value = allData;

      // 判断是否还有更多
      const lastPage = pages[pages.length - 1] || [];
      hasMore.value = hasMoreFn(lastPage, size.value - 1);
    } catch (err) {
      error.value =
        err instanceof Error ? err : new Error(String(err));
    } finally {
      isValidating.value = false;
    }
  };

  const setSize = (newSize: number | ((prev: number) => number)) => {
    size.value =
      typeof newSize === "function" ? newSize(size.value) : newSize;
    fetchAllPages();
  };

  const loadMore = async () => {
    if (!hasMore.value || isLoadingMore.value) return;

    isLoadingMore.value = true;

    try {
      const pageIndex = size.value;
      const key = getKey(pageIndex);
      const pageData = await fetcher(pageIndex);

      swrCache.set(key, {
        data: pageData,
        error: null,
        timestamp: Date.now(),
      });

      data.value = [...data.value, ...pageData];
      size.value++;

      hasMore.value = hasMoreFn(pageData, pageIndex);
    } catch (err) {
      error.value =
        err instanceof Error ? err : new Error(String(err));
    } finally {
      isLoadingMore.value = false;
    }
  };

  const revalidate = async () => {
    await fetchAllPages();
  };

  const reset = () => {
    data.value = [];
    error.value = null;
    size.value = initialSize;
    hasMore.value = true;
    fetchAllPages();
  };

  onMounted(() => {
    fetchAllPages();
  });

  return {
    data,
    error,
    isValidating,
    isLoadingMore,
    hasMore,
    size,
    setSize,
    loadMore,
    revalidate,
    reset,
  };
}

// ============================================================================
// 缓存管理
// ============================================================================

/**
 * 清除所有 SWR 缓存
 */
export function clearSWRCache(): void {
  swrCache.clear();
}

/**
 * 清除特定缓存
 */
export function deleteSWRCache(key: string): boolean {
  return swrCache.delete(key);
}

/**
 * 获取缓存
 */
export function getSWRCache<T>(key: string): T | null {
  const cached = swrCache.get(key);
  return (cached?.data as T) ?? null;
}

/**
 * 设置缓存
 */
export function setSWRCache<T>(key: string, data: T): void {
  swrCache.set(key, {
    data,
    error: null,
    timestamp: Date.now(),
  });

  // 通知订阅者
  const subs = subscribers.get(key);
  if (subs) {
    subs.forEach((callback) => callback());
  }
}

/**
 * 全局重新验证
 */
export function revalidateSWR(key: string): void {
  const subs = subscribers.get(key);
  if (subs) {
    subs.forEach((callback) => callback());
  }
}
