/**
 * Fetch Composable
 * 提供 HTTP 请求的响应式管理
 */

import {
  ref,
  computed,
  watch,
  onMounted,
  onUnmounted,
  type Ref,
  type ComputedRef,
  type MaybeRef,
} from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export type FetchStatus = "idle" | "pending" | "success" | "error";

export interface UseFetchOptions<T = unknown> extends RequestInit {
  /** 是否立即请求 */
  immediate?: boolean;
  /** 请求超时（毫秒） */
  timeout?: number;
  /** 重试次数 */
  retry?: number;
  /** 重试延迟（毫秒） */
  retryDelay?: number;
  /** 请求前回调 */
  beforeFetch?: (
    ctx: { url: string; options: RequestInit }
  ) => { url: string; options: RequestInit } | false | void;
  /** 请求后回调 */
  afterFetch?: (ctx: { data: T; response: Response }) => T;
  /** 错误回调 */
  onFetchError?: (ctx: {
    error: Error;
    data: unknown;
    response: Response | null;
  }) => void;
  /** 响应类型 */
  responseType?: "json" | "text" | "blob" | "arrayBuffer" | "formData";
  /** 是否在组件卸载时取消请求 */
  abortOnUnmount?: boolean;
  /** 缓存键 */
  cacheKey?: string;
  /** 缓存时间（毫秒） */
  cacheTime?: number;
  /** 请求防抖时间（毫秒） */
  debounce?: number;
}

export interface UseFetchReturn<T = unknown> {
  /** 响应数据 */
  data: Ref<T | null>;
  /** 请求状态 */
  status: Ref<FetchStatus>;
  /** 是否正在加载 */
  isLoading: ComputedRef<boolean>;
  /** 是否完成 */
  isFinished: ComputedRef<boolean>;
  /** 错误信息 */
  error: Ref<Error | null>;
  /** 原始响应 */
  response: Ref<Response | null>;
  /** 状态码 */
  statusCode: Ref<number | null>;
  /** 执行请求 */
  execute: (throwError?: boolean) => Promise<T>;
  /** 取消请求 */
  abort: () => void;
  /** 重试请求 */
  retry: () => Promise<T>;
  /** 是否可以取消 */
  canAbort: ComputedRef<boolean>;
  /** 是否已取消 */
  aborted: Ref<boolean>;
}

// 请求缓存
const fetchCache = new Map<
  string,
  { data: unknown; timestamp: number; promise?: Promise<unknown> }
>();

// ============================================================================
// useFetch - HTTP 请求
// ============================================================================

/**
 * 响应式 HTTP 请求
 * @example
 * const { data, isLoading, error, execute, abort } = useFetch<User[]>(
 *   '/api/users',
 *   { immediate: true }
 * )
 *
 * // 手动执行
 * await execute()
 *
 * // 取消请求
 * abort()
 */
export function useFetch<T = unknown>(
  url: MaybeRef<string>,
  options: UseFetchOptions<T> = {}
): UseFetchReturn<T> {
  const {
    immediate = true,
    timeout = 30000,
    retry = 0,
    retryDelay = 1000,
    beforeFetch,
    afterFetch,
    onFetchError,
    responseType = "json",
    abortOnUnmount = true,
    cacheKey,
    cacheTime = 0,
    debounce = 0,
    ...fetchOptions
  } = options;

  const data = ref<T | null>(null) as Ref<T | null>;
  const status = ref<FetchStatus>("idle");
  const error = ref<Error | null>(null);
  const response = ref<Response | null>(null);
  const statusCode = ref<number | null>(null);
  const aborted = ref(false);

  const isLoading = computed(() => status.value === "pending");
  const isFinished = computed(
    () => status.value === "success" || status.value === "error"
  );
  const canAbort = computed(() => isLoading.value);

  let abortController: AbortController | null = null;
  let debounceTimer: ReturnType<typeof setTimeout> | null = null;
  let retryCount = 0;

  const getUrl = () => {
    return typeof url === "object" && "value" in url ? url.value : url;
  };

  const getCacheKey = () => {
    return cacheKey || `${getUrl()}-${JSON.stringify(fetchOptions)}`;
  };

  // 从缓存获取
  const getFromCache = (): T | null => {
    if (cacheTime <= 0) return null;

    const cached = fetchCache.get(getCacheKey());
    if (cached && Date.now() - cached.timestamp < cacheTime) {
      return cached.data as T;
    }
    return null;
  };

  // 设置缓存
  const setCache = (value: T) => {
    if (cacheTime > 0) {
      fetchCache.set(getCacheKey(), {
        data: value,
        timestamp: Date.now(),
      });
    }
  };

  const abort = () => {
    if (abortController) {
      abortController.abort();
      aborted.value = true;
    }
    if (debounceTimer) {
      clearTimeout(debounceTimer);
    }
  };

  const execute = async (throwError = false): Promise<T> => {
    // 检查缓存
    const cached = getFromCache();
    if (cached !== null) {
      data.value = cached;
      status.value = "success";
      return cached;
    }

    // 取消之前的请求
    abort();
    aborted.value = false;

    // 创建新的 AbortController
    abortController = new AbortController();
    const signal = abortController.signal;

    // 超时处理
    let timeoutId: ReturnType<typeof setTimeout> | null = null;
    if (timeout > 0) {
      timeoutId = setTimeout(() => {
        abort();
        error.value = new Error(`Request timeout after ${timeout}ms`);
        status.value = "error";
      }, timeout);
    }

    const doFetch = async (): Promise<T> => {
      status.value = "pending";
      error.value = null;

      let requestUrl = getUrl();
      let requestOptions: RequestInit = {
        ...fetchOptions,
        signal,
      };

      // 请求前拦截
      if (beforeFetch) {
        const result = beforeFetch({ url: requestUrl, options: requestOptions });
        if (result === false) {
          throw new Error("Request aborted by beforeFetch");
        }
        if (result) {
          requestUrl = result.url;
          requestOptions = result.options;
        }
      }

      try {
        const res = await fetch(requestUrl, requestOptions);
        response.value = res;
        statusCode.value = res.status;

        if (!res.ok) {
          throw new Error(`HTTP error! status: ${res.status}`);
        }

        let responseData: T;

        switch (responseType) {
          case "text":
            responseData = (await res.text()) as T;
            break;
          case "blob":
            responseData = (await res.blob()) as T;
            break;
          case "arrayBuffer":
            responseData = (await res.arrayBuffer()) as T;
            break;
          case "formData":
            responseData = (await res.formData()) as T;
            break;
          default:
            responseData = await res.json();
        }

        // 请求后处理
        if (afterFetch) {
          responseData = afterFetch({ data: responseData, response: res });
        }

        data.value = responseData;
        status.value = "success";
        setCache(responseData);

        return responseData;
      } catch (err) {
        // 处理取消
        if (signal.aborted) {
          throw new Error("Request aborted");
        }

        const fetchError =
          err instanceof Error ? err : new Error(String(err));

        // 重试
        if (retryCount < retry) {
          retryCount++;
          await new Promise((resolve) =>
            setTimeout(resolve, retryDelay * retryCount)
          );
          return doFetch();
        }

        error.value = fetchError;
        status.value = "error";

        onFetchError?.({
          error: fetchError,
          data: data.value,
          response: response.value,
        });

        if (throwError) {
          throw fetchError;
        }

        return data.value as T;
      } finally {
        if (timeoutId) {
          clearTimeout(timeoutId);
        }
        retryCount = 0;
      }
    };

    // 防抖
    if (debounce > 0) {
      return new Promise((resolve, reject) => {
        debounceTimer = setTimeout(() => {
          doFetch().then(resolve).catch(reject);
        }, debounce);
      });
    }

    return doFetch();
  };

  const retryFn = () => {
    retryCount = 0;
    return execute();
  };

  // 监听 URL 变化
  if (typeof url === "object" && "value" in url) {
    watch(url, () => {
      if (status.value !== "idle") {
        execute();
      }
    });
  }

  // 自动执行
  if (immediate) {
    onMounted(() => {
      execute();
    });
  }

  // 组件卸载时取消请求
  if (abortOnUnmount) {
    onUnmounted(() => {
      abort();
    });
  }

  return {
    data,
    status,
    isLoading,
    isFinished,
    error,
    response,
    statusCode,
    execute,
    abort,
    retry: retryFn,
    canAbort,
    aborted,
  };
}

// ============================================================================
// createFetch - 创建带默认配置的 fetch
// ============================================================================

export interface CreateFetchOptions {
  /** 基础 URL */
  baseUrl?: string;
  /** 默认请求配置 */
  options?: UseFetchOptions;
  /** 请求拦截器 */
  interceptors?: {
    request?: (ctx: {
      url: string;
      options: RequestInit;
    }) => { url: string; options: RequestInit } | void;
    response?: (ctx: { data: unknown; response: Response }) => unknown;
    error?: (ctx: {
      error: Error;
      data: unknown;
      response: Response | null;
    }) => void;
  };
}

export interface CreateFetchReturn {
  /** 发起 GET 请求 */
  get: <T = unknown>(
    url: string,
    options?: UseFetchOptions<T>
  ) => UseFetchReturn<T>;
  /** 发起 POST 请求 */
  post: <T = unknown>(
    url: string,
    body?: unknown,
    options?: UseFetchOptions<T>
  ) => UseFetchReturn<T>;
  /** 发起 PUT 请求 */
  put: <T = unknown>(
    url: string,
    body?: unknown,
    options?: UseFetchOptions<T>
  ) => UseFetchReturn<T>;
  /** 发起 PATCH 请求 */
  patch: <T = unknown>(
    url: string,
    body?: unknown,
    options?: UseFetchOptions<T>
  ) => UseFetchReturn<T>;
  /** 发起 DELETE 请求 */
  delete: <T = unknown>(
    url: string,
    options?: UseFetchOptions<T>
  ) => UseFetchReturn<T>;
}

/**
 * 创建带默认配置的 fetch 实例
 * @example
 * const api = createFetch({
 *   baseUrl: '/api',
 *   options: {
 *     headers: {
 *       'Authorization': `Bearer ${token}`
 *     }
 *   },
 *   interceptors: {
 *     request: ({ url, options }) => {
 *       // 添加认证头
 *       return { url, options }
 *     },
 *     response: ({ data }) => {
 *       // 处理响应
 *       return data
 *     },
 *     error: ({ error }) => {
 *       // 处理错误
 *       console.error(error)
 *     }
 *   }
 * })
 *
 * // 使用
 * const { data } = api.get<User[]>('/users')
 * const { data } = api.post<User>('/users', { name: 'John' })
 */
export function createFetch(config: CreateFetchOptions = {}): CreateFetchReturn {
  const { baseUrl = "", options: defaultOptions = {}, interceptors = {} } = config;

  const createRequest = <T>(
    method: string,
    url: string,
    body?: unknown,
    options?: UseFetchOptions<T>
  ): UseFetchReturn<T> => {
    const fullUrl = baseUrl + url;

    const mergedOptions: UseFetchOptions<T> = {
      ...defaultOptions,
      ...options,
      method,
      beforeFetch: (ctx) => {
        // 应用默认拦截器
        if (interceptors.request) {
          const result = interceptors.request(ctx);
          if (result) {
            ctx = result;
          }
        }
        // 应用自定义拦截器
        if (options?.beforeFetch) {
          return options.beforeFetch(ctx);
        }
        return ctx;
      },
      afterFetch: (ctx) => {
        let data = ctx.data as T;
        // 应用默认拦截器
        if (interceptors.response) {
          data = interceptors.response(ctx) as T;
        }
        // 应用自定义拦截器
        if (options?.afterFetch) {
          data = options.afterFetch({ ...ctx, data });
        }
        return data;
      },
      onFetchError: (ctx) => {
        // 应用默认错误处理
        if (interceptors.error) {
          interceptors.error(ctx);
        }
        // 应用自定义错误处理
        if (options?.onFetchError) {
          options.onFetchError(ctx);
        }
      },
    };

    // 添加 body
    if (body !== undefined) {
      if (body instanceof FormData) {
        mergedOptions.body = body;
      } else {
        mergedOptions.body = JSON.stringify(body);
        mergedOptions.headers = {
          "Content-Type": "application/json",
          ...(mergedOptions.headers as Record<string, string>),
        };
      }
    }

    return useFetch<T>(fullUrl, mergedOptions);
  };

  return {
    get: <T>(url: string, options?: UseFetchOptions<T>) =>
      createRequest<T>("GET", url, undefined, options),
    post: <T>(url: string, body?: unknown, options?: UseFetchOptions<T>) =>
      createRequest<T>("POST", url, body, options),
    put: <T>(url: string, body?: unknown, options?: UseFetchOptions<T>) =>
      createRequest<T>("PUT", url, body, options),
    patch: <T>(url: string, body?: unknown, options?: UseFetchOptions<T>) =>
      createRequest<T>("PATCH", url, body, options),
    delete: <T>(url: string, options?: UseFetchOptions<T>) =>
      createRequest<T>("DELETE", url, undefined, options),
  };
}

// ============================================================================
// useLazyFetch - 懒加载请求
// ============================================================================

/**
 * 懒加载请求（不立即执行）
 * @example
 * const { execute, data } = useLazyFetch<User>('/api/user')
 *
 * // 需要时手动执行
 * const handleClick = async () => {
 *   await execute()
 * }
 */
export function useLazyFetch<T = unknown>(
  url: MaybeRef<string>,
  options: Omit<UseFetchOptions<T>, "immediate"> = {}
): UseFetchReturn<T> {
  return useFetch<T>(url, { ...options, immediate: false });
}

// ============================================================================
// usePaginatedFetch - 分页请求
// ============================================================================

export interface UsePaginatedFetchOptions<T> extends UseFetchOptions<T> {
  /** 每页数量 */
  pageSize?: number;
  /** 页码参数名 */
  pageParam?: string;
  /** 每页数量参数名 */
  pageSizeParam?: string;
}

export interface UsePaginatedFetchReturn<T>
  extends Omit<UseFetchReturn<T>, "data"> {
  /** 当前页数据 */
  data: Ref<T | null>;
  /** 当前页码 */
  page: Ref<number>;
  /** 每页数量 */
  pageSize: Ref<number>;
  /** 是否有下一页 */
  hasMore: Ref<boolean>;
  /** 加载下一页 */
  loadMore: () => Promise<T | null>;
  /** 刷新 */
  refresh: () => Promise<T | null>;
}

/**
 * 分页请求
 * @example
 * const { data, page, hasMore, loadMore, refresh } = usePaginatedFetch<User[]>(
 *   '/api/users',
 *   { pageSize: 10 }
 * )
 *
 * // 加载更多
 * const handleLoadMore = async () => {
 *   await loadMore()
 * }
 *
 * // 刷新
 * const handleRefresh = async () => {
 *   await refresh()
 * }
 */
export function usePaginatedFetch<T = unknown>(
  baseUrl: MaybeRef<string>,
  options: UsePaginatedFetchOptions<T> = {}
): UsePaginatedFetchReturn<T> {
  const {
    pageSize: initialPageSize = 10,
    pageParam = "page",
    pageSizeParam = "pageSize",
    ...fetchOptions
  } = options;

  const page = ref(1);
  const pageSize = ref(initialPageSize);
  const hasMore = ref(true);

  const getUrl = () => {
    const base =
      typeof baseUrl === "object" && "value" in baseUrl
        ? baseUrl.value
        : baseUrl;
    const separator = base.includes("?") ? "&" : "?";
    return `${base}${separator}${pageParam}=${page.value}&${pageSizeParam}=${pageSize.value}`;
  };

  const url = computed(getUrl);

  const fetchReturn = useFetch<T>(url, {
    ...fetchOptions,
    immediate: false,
  });

  // 检查是否有更多数据
  watch(fetchReturn.data, (newData) => {
    if (Array.isArray(newData)) {
      hasMore.value = newData.length >= pageSize.value;
    }
  });

  const loadMore = async (): Promise<T | null> => {
    if (!hasMore.value) return null;
    page.value++;
    return fetchReturn.execute();
  };

  const refresh = async (): Promise<T | null> => {
    page.value = 1;
    hasMore.value = true;
    return fetchReturn.execute();
  };

  // 初始加载
  if (options.immediate !== false) {
    onMounted(() => {
      fetchReturn.execute();
    });
  }

  return {
    ...fetchReturn,
    page,
    pageSize,
    hasMore,
    loadMore,
    refresh,
  };
}

// ============================================================================
// useInfiniteFetch - 无限加载请求
// ============================================================================

export interface UseInfiniteFetchOptions<T> extends UsePaginatedFetchOptions<T> {
  /** 数据合并函数 */
  merge?: (prev: T[], next: T[]) => T[];
}

export interface UseInfiniteFetchReturn<T>
  extends Omit<UsePaginatedFetchReturn<T[]>, "data"> {
  /** 所有数据 */
  data: Ref<T[]>;
  /** 重置数据 */
  reset: () => void;
}

/**
 * 无限加载请求（累积数据）
 * @example
 * const { data, loadMore, hasMore, isLoading, reset } = useInfiniteFetch<User>(
 *   '/api/users',
 *   { pageSize: 20 }
 * )
 *
 * // data 包含所有已加载的数据
 * // 加载更多时数据会累积
 */
export function useInfiniteFetch<T = unknown>(
  baseUrl: MaybeRef<string>,
  options: UseInfiniteFetchOptions<T> = {}
): UseInfiniteFetchReturn<T> {
  const { merge = (prev, next) => [...prev, ...next], ...paginatedOptions } =
    options;

  const allData = ref<T[]>([]) as Ref<T[]>;

  const paginatedFetch = usePaginatedFetch<T[]>(baseUrl, {
    ...paginatedOptions,
    immediate: false,
    afterFetch: (ctx) => {
      const newData = ctx.data;
      if (Array.isArray(newData)) {
        allData.value = merge(allData.value, newData);
      }
      return newData;
    },
  });

  const reset = () => {
    allData.value = [];
    paginatedFetch.page.value = 1;
    paginatedFetch.hasMore.value = true;
  };

  // 初始加载
  if (options.immediate !== false) {
    onMounted(() => {
      paginatedFetch.execute();
    });
  }

  return {
    ...paginatedFetch,
    data: allData,
    reset,
  };
}

// ============================================================================
// 清除缓存
// ============================================================================

/**
 * 清除所有 fetch 缓存
 */
export function clearFetchCache(): void {
  fetchCache.clear();
}

/**
 * 清除特定缓存
 */
export function deleteFetchCache(key: string): boolean {
  return fetchCache.delete(key);
}

/**
 * 获取缓存大小
 */
export function getFetchCacheSize(): number {
  return fetchCache.size;
}
