/**
 * Computed Composable
 * 提供增强的 computed 工具函数
 */

import { computed, ref, watch, shallowRef, triggerRef, type Ref, type ComputedRef, type WatchSource } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

/**
 * 异步 computed 选项
 */
export interface ComputedAsyncOptions<T> {
  /** 初始值 */
  initialValue?: T;
  /** 是否懒加载（首次不执行） */
  lazy?: boolean;
  /** 错误处理函数 */
  onError?: (error: Error) => void;
  /** 是否在依赖变化时重新计算 */
  shallow?: boolean;
  /** 防抖延迟（毫秒） */
  debounce?: number;
}

/**
 * 异步 computed 返回值
 */
export interface ComputedAsyncReturn<T> {
  /** 计算结果 */
  state: Ref<T>;
  /** 是否正在加载 */
  isLoading: Ref<boolean>;
  /** 错误信息 */
  error: Ref<Error | null>;
  /** 手动重新计算 */
  execute: () => Promise<void>;
}

/**
 * 可控 computed 选项
 */
export interface ComputedWithControlOptions {
  /** 是否立即评估 */
  immediate?: boolean;
}

/**
 * 可控 computed 返回值
 */
export interface ComputedWithControlReturn<T> {
  /** 计算结果 */
  state: ComputedRef<T>;
  /** 手动触发重新计算 */
  trigger: () => void;
  /** 暂停自动更新 */
  pause: () => void;
  /** 恢复自动更新 */
  resume: () => void;
  /** 是否暂停中 */
  isPaused: Ref<boolean>;
}

/**
 * 防抖 computed 选项
 */
export interface ComputedDebouncedOptions {
  /** 防抖延迟（毫秒） */
  debounce?: number;
  /** 最大等待时间 */
  maxWait?: number;
}

/**
 * 节流 computed 选项
 */
export interface ComputedThrottledOptions {
  /** 节流间隔（毫秒） */
  throttle?: number;
  /** 是否在开始时执行 */
  leading?: boolean;
  /** 是否在结束时执行 */
  trailing?: boolean;
}

// ============================================================================
// 核心函数
// ============================================================================

/**
 * 立即求值的 computed（无缓存）
 *
 * @description 每次访问都会重新计算，适用于依赖外部状态的场景
 *
 * @example
 * ```ts
 * const now = computedEager(() => Date.now())
 * console.log(now.value) // 每次都是最新时间
 * ```
 */
export function computedEager<T>(getter: () => T): Readonly<Ref<T>> {
  const result = shallowRef<T>(getter());

  // 创建一个 effect 来追踪依赖
  watch(
    getter,
    (value) => {
      result.value = value;
    },
    { flush: "sync" }
  );

  return result;
}

/**
 * 异步 computed
 *
 * @description 支持异步 getter，自动处理 loading 和 error 状态
 *
 * @example
 * ```ts
 * const userId = ref(1)
 * const { state: user, isLoading, error } = computedAsync(
 *   async () => {
 *     const res = await fetch(`/api/users/${userId.value}`)
 *     return res.json()
 *   },
 *   { initialValue: null }
 * )
 * ```
 */
export function computedAsync<T>(
  getter: () => Promise<T>,
  options: ComputedAsyncOptions<T> = {}
): ComputedAsyncReturn<T> {
  const { initialValue = undefined as T, lazy = false, onError, shallow = false, debounce = 0 } = options;

  const state = shallow ? shallowRef<T>(initialValue) : ref<T>(initialValue);
  const isLoading = ref(false);
  const error = ref<Error | null>(null);

  let timer: ReturnType<typeof setTimeout> | null = null;

  const execute = async () => {
    isLoading.value = true;
    error.value = null;

    try {
      const result = await getter();
      state.value = result;
    } catch (e) {
      const err = e instanceof Error ? e : new Error(String(e));
      error.value = err;
      onError?.(err);
    } finally {
      isLoading.value = false;
    }
  };

  const debouncedExecute = () => {
    if (timer) {
      clearTimeout(timer);
    }
    if (debounce > 0) {
      timer = setTimeout(execute, debounce);
    } else {
      execute();
    }
  };

  // 监听 getter 中的响应式依赖
  watch(
    () => {
      try {
        // 触发 getter 来追踪依赖，但不等待结果
        getter();
      } catch {
        // 忽略同步错误
      }
    },
    () => {
      debouncedExecute();
    },
    { immediate: !lazy }
  );

  return {
    state: state as Ref<T>,
    isLoading,
    error,
    execute,
  };
}

/**
 * 可控的 computed
 *
 * @description 支持暂停/恢复自动更新，以及手动触发重新计算
 *
 * @example
 * ```ts
 * const count = ref(0)
 * const { state, pause, resume, trigger } = computedWithControl(
 *   () => count.value * 2
 * )
 *
 * pause() // 暂停自动更新
 * count.value = 5 // state 不会更新
 * trigger() // 手动触发更新
 * resume() // 恢复自动更新
 * ```
 */
export function computedWithControl<T>(
  getter: () => T,
  options: ComputedWithControlOptions = {}
): ComputedWithControlReturn<T> {
  const { immediate = true } = options;

  const isPaused = ref(false);
  const trigger$ = ref(0);

  const state = computed(() => {
    // 触发依赖追踪
    void trigger$.value;
    return getter();
  });

  // 使用 shallowRef 存储实际值
  const internalValue = shallowRef<T>(immediate ? getter() : (undefined as T));

  watch(
    state,
    (value) => {
      if (!isPaused.value) {
        internalValue.value = value;
      }
    },
    { immediate }
  );

  const trigger = () => {
    trigger$.value++;
    internalValue.value = getter();
  };

  const pause = () => {
    isPaused.value = true;
  };

  const resume = () => {
    isPaused.value = false;
    internalValue.value = getter();
  };

  return {
    state: computed(() => internalValue.value),
    trigger,
    pause,
    resume,
    isPaused,
  };
}

/**
 * 从响应式对象中选取属性
 *
 * @description 创建一个只包含指定属性的 computed 对象
 *
 * @example
 * ```ts
 * const user = reactive({ id: 1, name: 'John', email: 'john@example.com' })
 * const basicInfo = computedPick(user, ['id', 'name'])
 * // basicInfo.value = { id: 1, name: 'John' }
 * ```
 */
export function computedPick<T extends object, K extends keyof T>(source: T, keys: K[]): ComputedRef<Pick<T, K>> {
  return computed(() => {
    const result = {} as Pick<T, K>;
    for (const key of keys) {
      result[key] = source[key];
    }
    return result;
  });
}

/**
 * 从响应式对象中排除属性
 *
 * @description 创建一个排除指定属性的 computed 对象
 *
 * @example
 * ```ts
 * const user = reactive({ id: 1, name: 'John', password: 'secret' })
 * const safeUser = computedOmit(user, ['password'])
 * // safeUser.value = { id: 1, name: 'John' }
 * ```
 */
export function computedOmit<T extends object, K extends keyof T>(source: T, keys: K[]): ComputedRef<Omit<T, K>> {
  return computed(() => {
    const result = { ...source };
    for (const key of keys) {
      delete result[key];
    }
    return result as Omit<T, K>;
  });
}

/**
 * 将 getter 函数转换为 computed
 *
 * @description 包装普通函数为 computed，支持传入参数
 *
 * @example
 * ```ts
 * const items = ref([1, 2, 3, 4, 5])
 * const doubledItems = toComputed(() => items.value.map(x => x * 2))
 * ```
 */
export function toComputed<T>(getter: () => T): ComputedRef<T> {
  return computed(getter);
}

/**
 * 从多个源创建 computed
 *
 * @description 组合多个响应式源创建新的 computed
 *
 * @example
 * ```ts
 * const firstName = ref('John')
 * const lastName = ref('Doe')
 * const fullName = computedFrom(
 *   [firstName, lastName],
 *   ([first, last]) => `${first} ${last}`
 * )
 * ```
 */
export function computedFrom<T extends readonly unknown[], R>(
  sources: [...{ [K in keyof T]: WatchSource<T[K]> | T[K] }],
  getter: (values: T) => R
): ComputedRef<R> {
  return computed(() => {
    const values = sources.map((source) => {
      if (typeof source === "function") {
        return (source as () => unknown)();
      }
      if (source && typeof source === "object" && "value" in source) {
        return (source as Ref<unknown>).value;
      }
      return source;
    }) as unknown as T;
    return getter(values);
  });
}

/**
 * 防抖 computed
 *
 * @description 当依赖变化时，延迟更新 computed 值
 *
 * @example
 * ```ts
 * const searchQuery = ref('')
 * const debouncedQuery = computedDebounced(
 *   () => searchQuery.value.trim(),
 *   { debounce: 300 }
 * )
 * ```
 */
export function computedDebounced<T>(getter: () => T, options: ComputedDebouncedOptions = {}): Readonly<Ref<T>> {
  const { debounce = 250, maxWait } = options;

  const result = shallowRef<T>(getter());
  let timer: ReturnType<typeof setTimeout> | null = null;
  let maxTimer: ReturnType<typeof setTimeout> | null = null;
  let lastUpdate = Date.now();

  const update = () => {
    result.value = getter();
    lastUpdate = Date.now();
    if (maxTimer) {
      clearTimeout(maxTimer);
      maxTimer = null;
    }
  };

  watch(
    getter,
    () => {
      if (timer) {
        clearTimeout(timer);
      }

      // 设置 maxWait 定时器
      if (maxWait && !maxTimer) {
        maxTimer = setTimeout(() => {
          if (timer) {
            clearTimeout(timer);
            timer = null;
          }
          update();
        }, maxWait);
      }

      timer = setTimeout(update, debounce);
    },
    { flush: "sync" }
  );

  return result;
}

/**
 * 节流 computed
 *
 * @description 限制 computed 更新频率
 *
 * @example
 * ```ts
 * const scrollY = ref(0)
 * const throttledScrollY = computedThrottled(
 *   () => scrollY.value,
 *   { throttle: 100 }
 * )
 * ```
 */
export function computedThrottled<T>(getter: () => T, options: ComputedThrottledOptions = {}): Readonly<Ref<T>> {
  const { throttle = 100, leading = true, trailing = true } = options;

  const result = shallowRef<T>(getter());
  let timer: ReturnType<typeof setTimeout> | null = null;
  let lastExec = 0;
  let pendingValue: T | undefined;

  const update = (value: T) => {
    result.value = value;
    lastExec = Date.now();
  };

  watch(
    getter,
    (value) => {
      const now = Date.now();
      const elapsed = now - lastExec;

      if (elapsed >= throttle) {
        if (timer) {
          clearTimeout(timer);
          timer = null;
        }
        if (leading || lastExec > 0) {
          update(value);
        } else {
          lastExec = now;
        }
      } else {
        pendingValue = value;
        if (!timer && trailing) {
          timer = setTimeout(() => {
            timer = null;
            if (pendingValue !== undefined) {
              update(pendingValue);
              pendingValue = undefined;
            }
          }, throttle - elapsed);
        }
      }
    },
    { flush: "sync" }
  );

  return result;
}

/**
 * 带历史记录的 computed
 *
 * @description 保存 computed 值的历史记录
 *
 * @example
 * ```ts
 * const count = ref(0)
 * const { current, history, canUndo, undo } = computedWithHistory(
 *   () => count.value * 2,
 *   { capacity: 10 }
 * )
 * ```
 */
export function computedWithHistory<T>(
  getter: () => T,
  options: { capacity?: number } = {}
): {
  current: ComputedRef<T>;
  history: Readonly<Ref<T[]>>;
  canUndo: ComputedRef<boolean>;
  undo: () => T | undefined;
  clear: () => void;
} {
  const { capacity = 10 } = options;

  const current = computed(getter);
  const history = ref<T[]>([]) as Ref<T[]>;

  watch(
    current,
    (value, oldValue) => {
      if (oldValue !== undefined) {
        history.value = [...history.value.slice(-(capacity - 1)), oldValue];
      }
    },
    { immediate: true }
  );

  const canUndo = computed(() => history.value.length > 0);

  const undo = () => {
    return history.value.pop();
  };

  const clear = () => {
    history.value = [];
  };

  return {
    current,
    history,
    canUndo,
    undo,
    clear,
  };
}

/**
 * 条件 computed
 *
 * @description 根据条件返回不同的 computed 值
 *
 * @example
 * ```ts
 * const isAdmin = ref(false)
 * const permissions = computedIf(
 *   isAdmin,
 *   () => ['read', 'write', 'delete'],
 *   () => ['read']
 * )
 * ```
 */
export function computedIf<T, F>(
  condition: Ref<boolean> | (() => boolean),
  truthy: () => T,
  falsy: () => F
): ComputedRef<T | F> {
  return computed(() => {
    const cond = typeof condition === "function" ? condition() : condition.value;
    return cond ? truthy() : falsy();
  });
}

/**
 * 可写的 computed（带转换）
 *
 * @description 创建可写的 computed，支持自定义 get/set 转换
 *
 * @example
 * ```ts
 * const rawValue = ref('hello')
 * const upperValue = writableComputed(
 *   () => rawValue.value.toUpperCase(),
 *   (value) => { rawValue.value = value.toLowerCase() }
 * )
 * upperValue.value = 'WORLD' // rawValue 变为 'world'
 * ```
 */
export function writableComputed<T>(getter: () => T, setter: (value: T) => void): Ref<T> {
  return computed({
    get: getter,
    set: setter,
  });
}

/**
 * 计算数组的派生值
 *
 * @description 提供数组相关的计算属性
 *
 * @example
 * ```ts
 * const numbers = ref([1, 2, 3, 4, 5])
 * const { sum, avg, min, max, count } = computedArray(numbers)
 * ```
 */
export function computedArray<T>(source: Ref<T[]>): {
  sum: ComputedRef<number>;
  avg: ComputedRef<number>;
  min: ComputedRef<T | undefined>;
  max: ComputedRef<T | undefined>;
  count: ComputedRef<number>;
  first: ComputedRef<T | undefined>;
  last: ComputedRef<T | undefined>;
  isEmpty: ComputedRef<boolean>;
  unique: ComputedRef<T[]>;
  sorted: ComputedRef<T[]>;
  reversed: ComputedRef<T[]>;
} {
  const sum = computed(() => {
    const arr = source.value;
    if (arr.length === 0) return 0;
    return arr.reduce((acc, val) => acc + (val as number), 0);
  });

  const avg = computed(() => {
    const arr = source.value;
    if (arr.length === 0) return 0;
    return sum.value / arr.length;
  });

  const min = computed(() => {
    const arr = source.value;
    if (arr.length === 0) return undefined;
    return arr.reduce((a, b) => (a < b ? a : b));
  });

  const max = computed(() => {
    const arr = source.value;
    if (arr.length === 0) return undefined;
    return arr.reduce((a, b) => (a > b ? a : b));
  });

  const count = computed(() => source.value.length);

  const first = computed(() => source.value[0]);

  const last = computed(() => source.value[source.value.length - 1]);

  const isEmpty = computed(() => source.value.length === 0);

  const unique = computed(() => [...new Set(source.value)]);

  const sorted = computed(() => [...source.value].sort());

  const reversed = computed(() => [...source.value].reverse());

  return {
    sum,
    avg,
    min,
    max,
    count,
    first,
    last,
    isEmpty,
    unique,
    sorted,
    reversed,
  };
}

/**
 * 延迟更新的 computed
 *
 * @description computed 值在指定延迟后才更新
 *
 * @example
 * ```ts
 * const query = ref('')
 * const delayedQuery = computedDelayed(() => query.value, 500)
 * ```
 */
export function computedDelayed<T>(getter: () => T, delay: number): Readonly<Ref<T>> {
  const result = shallowRef<T>(getter());
  let timer: ReturnType<typeof setTimeout> | null = null;

  watch(
    getter,
    (value) => {
      if (timer) {
        clearTimeout(timer);
      }
      timer = setTimeout(() => {
        result.value = value;
      }, delay);
    },
    { flush: "sync" }
  );

  return result;
}

/**
 * 带默认值的 computed
 *
 * @description 当计算结果为 null 或 undefined 时返回默认值
 *
 * @example
 * ```ts
 * const data = ref<string | null>(null)
 * const safeData = computedDefault(() => data.value, 'default')
 * ```
 */
export function computedDefault<T>(getter: () => T | null | undefined, defaultValue: T): ComputedRef<T> {
  return computed(() => getter() ?? defaultValue);
}

/**
 * 组合多个 computed 为对象
 *
 * @description 将多个 computed 组合为单个响应式对象
 *
 * @example
 * ```ts
 * const firstName = ref('John')
 * const lastName = ref('Doe')
 * const user = computedObject({
 *   fullName: () => `${firstName.value} ${lastName.value}`,
 *   initials: () => `${firstName.value[0]}${lastName.value[0]}`
 * })
 * ```
 */
export function computedObject<T extends Record<string, () => unknown>>(
  getters: T
): ComputedRef<{ [K in keyof T]: ReturnType<T[K]> }> {
  return computed(() => {
    const result = {} as { [K in keyof T]: ReturnType<T[K]> };
    for (const key in getters) {
      result[key] = getters[key]() as ReturnType<T[typeof key]>;
    }
    return result;
  });
}

/**
 * 映射 computed
 *
 * @description 对数组中每个元素应用转换函数
 *
 * @example
 * ```ts
 * const users = ref([{ name: 'John' }, { name: 'Jane' }])
 * const names = computedMap(users, user => user.name)
 * ```
 */
export function computedMap<T, R>(source: Ref<T[]>, mapper: (item: T, index: number) => R): ComputedRef<R[]> {
  return computed(() => source.value.map(mapper));
}

/**
 * 过滤 computed
 *
 * @description 对数组应用过滤条件
 *
 * @example
 * ```ts
 * const numbers = ref([1, 2, 3, 4, 5])
 * const evens = computedFilter(numbers, n => n % 2 === 0)
 * ```
 */
export function computedFilter<T>(source: Ref<T[]>, predicate: (item: T, index: number) => boolean): ComputedRef<T[]> {
  return computed(() => source.value.filter(predicate));
}

/**
 * 查找 computed
 *
 * @description 在数组中查找满足条件的第一个元素
 *
 * @example
 * ```ts
 * const users = ref([{ id: 1 }, { id: 2 }])
 * const user = computedFind(users, u => u.id === 2)
 * ```
 */
export function computedFind<T>(
  source: Ref<T[]>,
  predicate: (item: T, index: number) => boolean
): ComputedRef<T | undefined> {
  return computed(() => source.value.find(predicate));
}

/**
 * 分组 computed
 *
 * @description 按键函数对数组进行分组
 *
 * @example
 * ```ts
 * const items = ref([
 *   { type: 'a', value: 1 },
 *   { type: 'b', value: 2 },
 *   { type: 'a', value: 3 }
 * ])
 * const grouped = computedGroupBy(items, item => item.type)
 * // { a: [...], b: [...] }
 * ```
 */
export function computedGroupBy<T, K extends string | number>(
  source: Ref<T[]>,
  keyFn: (item: T) => K
): ComputedRef<Record<K, T[]>> {
  return computed(() => {
    const result = {} as Record<K, T[]>;
    for (const item of source.value) {
      const key = keyFn(item);
      if (!result[key]) {
        result[key] = [];
      }
      result[key].push(item);
    }
    return result;
  });
}

/**
 * 排序 computed
 *
 * @description 对数组进行排序
 *
 * @example
 * ```ts
 * const users = ref([{ name: 'John' }, { name: 'Alice' }])
 * const sorted = computedSort(users, (a, b) => a.name.localeCompare(b.name))
 * ```
 */
export function computedSort<T>(source: Ref<T[]>, compareFn?: (a: T, b: T) => number): ComputedRef<T[]> {
  return computed(() => [...source.value].sort(compareFn));
}
