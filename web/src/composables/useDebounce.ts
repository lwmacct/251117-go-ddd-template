/**
 * 防抖工具 Composable
 * 用于搜索输入等需要减少请求频率的场景
 */
import { ref, watch, type Ref, onUnmounted } from "vue";

export interface UseDebounceOptions {
  /** 延迟时间（毫秒），默认 300ms */
  delay?: number;
  /** 是否立即执行第一次，默认 false */
  immediate?: boolean;
}

/**
 * 创建防抖值
 * @param source 源值
 * @param options 配置选项
 * @returns 防抖后的值
 */
export function useDebouncedRef<T>(
  source: Ref<T>,
  options: UseDebounceOptions = {}
): Ref<T> {
  const { delay = 300, immediate = false } = options;

  const debouncedValue = ref(source.value) as Ref<T>;
  let timeoutId: ReturnType<typeof setTimeout> | null = null;
  let isFirstCall = true;

  watch(
    source,
    (newValue) => {
      if (timeoutId) {
        clearTimeout(timeoutId);
      }

      if (immediate && isFirstCall) {
        debouncedValue.value = newValue;
        isFirstCall = false;
        return;
      }

      timeoutId = setTimeout(() => {
        debouncedValue.value = newValue;
      }, delay);
    },
    { immediate: false }
  );

  // 清理定时器
  onUnmounted(() => {
    if (timeoutId) {
      clearTimeout(timeoutId);
    }
  });

  return debouncedValue;
}

/**
 * 创建防抖函数
 * @param fn 要防抖的函数
 * @param options 配置选项
 * @returns 防抖后的函数和取消方法
 */
export function useDebounceFn<T extends (...args: unknown[]) => unknown>(
  fn: T,
  options: UseDebounceOptions = {}
): {
  /** 防抖后的函数 */
  debouncedFn: (...args: Parameters<T>) => void;
  /** 取消待执行的调用 */
  cancel: () => void;
  /** 立即执行 */
  flush: (...args: Parameters<T>) => ReturnType<T>;
  /** 是否有待执行的调用 */
  pending: Ref<boolean>;
} {
  const { delay = 300, immediate = false } = options;

  let timeoutId: ReturnType<typeof setTimeout> | null = null;
  let lastArgs: Parameters<T> | null = null;
  const pending = ref(false);
  let isFirstCall = true;

  const cancel = () => {
    if (timeoutId) {
      clearTimeout(timeoutId);
      timeoutId = null;
    }
    pending.value = false;
    lastArgs = null;
  };

  const flush = (...args: Parameters<T>): ReturnType<T> => {
    cancel();
    return fn(...args) as ReturnType<T>;
  };

  const debouncedFn = (...args: Parameters<T>) => {
    lastArgs = args;

    if (timeoutId) {
      clearTimeout(timeoutId);
    }

    if (immediate && isFirstCall) {
      isFirstCall = false;
      fn(...args);
      return;
    }

    pending.value = true;
    timeoutId = setTimeout(() => {
      if (lastArgs) {
        fn(...lastArgs);
      }
      pending.value = false;
      timeoutId = null;
    }, delay);
  };

  // 清理定时器
  onUnmounted(() => {
    cancel();
  });

  return {
    debouncedFn,
    cancel,
    flush,
    pending,
  };
}

/**
 * 创建搜索输入防抖
 * 专门用于搜索场景，提供搜索值和加载状态
 */
export function useSearchDebounce(options: UseDebounceOptions = {}) {
  const { delay = 300 } = options;

  const searchQuery = ref("");
  const debouncedQuery = useDebouncedRef(searchQuery, { delay });
  const isSearching = ref(false);

  // 监听防抖值变化，用于显示加载状态
  watch(searchQuery, () => {
    if (searchQuery.value !== debouncedQuery.value) {
      isSearching.value = true;
    }
  });

  watch(debouncedQuery, () => {
    isSearching.value = false;
  });

  const clear = () => {
    searchQuery.value = "";
  };

  return {
    /** 原始搜索值（用于输入绑定） */
    searchQuery,
    /** 防抖后的搜索值（用于 API 调用） */
    debouncedQuery,
    /** 是否正在等待防抖（可用于显示加载指示器） */
    isSearching,
    /** 清空搜索 */
    clear,
  };
}
