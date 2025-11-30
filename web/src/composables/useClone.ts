/**
 * Clone Composable
 * 提供响应式对象的克隆和同步功能
 */

import {
  ref,
  watch,
  computed,
  toRaw,
  isRef,
  unref,
  type Ref,
  type UnwrapRef,
  type ComputedRef,
  type WatchStopHandle,
} from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UseClonedOptions<T> {
  /** 是否深克隆 */
  deep?: boolean;
  /** 是否立即克隆 */
  immediate?: boolean;
  /** 自定义克隆函数 */
  clone?: (source: T) => T;
  /** 是否手动同步 */
  manual?: boolean;
}

export interface UseClonedReturn<T> {
  /** 克隆的值 */
  cloned: Ref<T>;
  /** 同步克隆值到源值 */
  sync: () => void;
  /** 重新从源值克隆 */
  reset: () => void;
  /** 是否已修改 */
  isModified: ComputedRef<boolean>;
}

// ============================================================================
// 深克隆工具函数
// ============================================================================

/**
 * 深克隆对象
 */
export function deepClone<T>(source: T): T {
  if (source === null || typeof source !== "object") {
    return source;
  }

  // 处理 Date
  if (source instanceof Date) {
    return new Date(source.getTime()) as T;
  }

  // 处理 RegExp
  if (source instanceof RegExp) {
    return new RegExp(source.source, source.flags) as T;
  }

  // 处理 Map
  if (source instanceof Map) {
    const result = new Map();
    source.forEach((value, key) => {
      result.set(deepClone(key), deepClone(value));
    });
    return result as T;
  }

  // 处理 Set
  if (source instanceof Set) {
    const result = new Set();
    source.forEach((value) => {
      result.add(deepClone(value));
    });
    return result as T;
  }

  // 处理数组
  if (Array.isArray(source)) {
    return source.map((item) => deepClone(item)) as T;
  }

  // 处理普通对象
  const result: Record<string, unknown> = {};
  for (const key in source) {
    if (Object.prototype.hasOwnProperty.call(source, key)) {
      result[key] = deepClone((source as Record<string, unknown>)[key]);
    }
  }
  return result as T;
}

/**
 * 使用 structuredClone 或回退到 JSON 方式克隆
 */
export function structuredClonePolyfill<T>(source: T): T {
  if (typeof structuredClone === "function") {
    return structuredClone(source);
  }
  // 回退到 JSON 方式
  return JSON.parse(JSON.stringify(source));
}

// ============================================================================
// useCloned - 克隆响应式值
// ============================================================================

/**
 * 克隆响应式值
 * @example
 * const source = ref({ name: 'John', age: 25 })
 * const { cloned, sync, reset, isModified } = useCloned(source)
 *
 * // 修改克隆值
 * cloned.value.name = 'Jane'
 *
 * // 检查是否修改
 * console.log(isModified.value) // true
 *
 * // 同步到源值
 * sync()
 *
 * // 重置为源值
 * reset()
 */
export function useCloned<T>(
  source: Ref<T> | T,
  options: UseClonedOptions<T> = {}
): UseClonedReturn<T> {
  const { deep = true, immediate = true, clone = deepClone, manual = false } =
    options;

  const getSource = () => {
    const value = isRef(source) ? source.value : source;
    return toRaw(value);
  };

  const cloned = ref<T>(immediate ? clone(getSource()) : getSource()) as Ref<T>;

  // 监听源值变化
  let stopWatch: WatchStopHandle | null = null;

  if (!manual && isRef(source)) {
    stopWatch = watch(
      source,
      (newValue) => {
        cloned.value = clone(toRaw(newValue)) as UnwrapRef<T>;
      },
      { deep }
    );
  }

  // 同步克隆值到源值
  const sync = () => {
    if (isRef(source)) {
      (source as Ref<T>).value = clone(toRaw(cloned.value));
    }
  };

  // 重置为源值
  const reset = () => {
    cloned.value = clone(getSource()) as UnwrapRef<T>;
  };

  // 检查是否已修改
  const isModified = computed(() => {
    return JSON.stringify(toRaw(cloned.value)) !== JSON.stringify(getSource());
  });

  return {
    cloned,
    sync,
    reset,
    isModified,
  };
}

// ============================================================================
// useManualClone - 手动控制的克隆
// ============================================================================

export interface UseManualCloneReturn<T> {
  /** 克隆的值 */
  cloned: Ref<T>;
  /** 从源值克隆 */
  clone: () => void;
  /** 应用克隆值到源值 */
  apply: () => void;
  /** 重置 */
  reset: () => void;
  /** 是否已修改 */
  isModified: ComputedRef<boolean>;
}

/**
 * 手动控制的克隆（适用于表单编辑场景）
 * @example
 * const user = ref({ name: 'John', email: 'john@example.com' })
 * const { cloned, apply, reset, isModified } = useManualClone(user)
 *
 * // 用户在表单中编辑 cloned
 * cloned.value.name = 'Jane'
 *
 * // 保存时应用更改
 * const handleSave = () => {
 *   apply()
 *   saveToServer(user.value)
 * }
 *
 * // 取消时重置
 * const handleCancel = () => {
 *   reset()
 * }
 */
export function useManualClone<T>(
  source: Ref<T>,
  cloneFn: (source: T) => T = deepClone
): UseManualCloneReturn<T> {
  const cloned = ref<T>(cloneFn(toRaw(source.value))) as Ref<T>;

  const clone = () => {
    cloned.value = cloneFn(toRaw(source.value)) as UnwrapRef<T>;
  };

  const apply = () => {
    source.value = cloneFn(toRaw(cloned.value));
  };

  const reset = () => {
    cloned.value = cloneFn(toRaw(source.value)) as UnwrapRef<T>;
  };

  const isModified = computed(() => {
    return (
      JSON.stringify(toRaw(cloned.value)) !==
      JSON.stringify(toRaw(source.value))
    );
  });

  return {
    cloned,
    clone,
    apply,
    reset,
    isModified,
  };
}

// ============================================================================
// useDirtyState - 脏状态检测
// ============================================================================

export interface UseDirtyStateReturn<T> {
  /** 当前值 */
  state: Ref<T>;
  /** 是否为脏状态 */
  isDirty: ComputedRef<boolean>;
  /** 脏字段列表 */
  dirtyFields: ComputedRef<string[]>;
  /** 标记为干净 */
  markClean: () => void;
  /** 重置为初始状态 */
  reset: () => void;
  /** 获取更改 */
  getChanges: () => Partial<T>;
}

/**
 * 脏状态检测
 * @example
 * const { state, isDirty, dirtyFields, markClean, reset } = useDirtyState({
 *   name: 'John',
 *   email: 'john@example.com'
 * })
 *
 * state.value.name = 'Jane'
 * console.log(isDirty.value) // true
 * console.log(dirtyFields.value) // ['name']
 *
 * // 保存后标记为干净
 * markClean()
 *
 * // 或重置为初始状态
 * reset()
 */
export function useDirtyState<T extends object>(
  initialState: T
): UseDirtyStateReturn<T> {
  const cleanState = ref<T>(deepClone(initialState)) as Ref<T>;
  const state = ref<T>(deepClone(initialState)) as Ref<T>;

  const isDirty = computed(() => {
    return (
      JSON.stringify(toRaw(state.value)) !==
      JSON.stringify(toRaw(cleanState.value))
    );
  });

  const dirtyFields = computed(() => {
    const dirty: string[] = [];
    const current = toRaw(state.value) as Record<string, unknown>;
    const clean = toRaw(cleanState.value) as Record<string, unknown>;

    for (const key in current) {
      if (JSON.stringify(current[key]) !== JSON.stringify(clean[key])) {
        dirty.push(key);
      }
    }

    return dirty;
  });

  const markClean = () => {
    cleanState.value = deepClone(toRaw(state.value)) as UnwrapRef<T>;
  };

  const reset = () => {
    state.value = deepClone(toRaw(cleanState.value)) as UnwrapRef<T>;
  };

  const getChanges = (): Partial<T> => {
    const changes: Partial<T> = {};
    const current = toRaw(state.value) as Record<string, unknown>;
    const clean = toRaw(cleanState.value) as Record<string, unknown>;

    for (const key in current) {
      if (JSON.stringify(current[key]) !== JSON.stringify(clean[key])) {
        (changes as Record<string, unknown>)[key] = current[key];
      }
    }

    return changes;
  };

  return {
    state,
    isDirty,
    dirtyFields,
    markClean,
    reset,
    getChanges,
  };
}

// ============================================================================
// useSnapshot - 状态快照
// ============================================================================

export interface UseSnapshotReturn<T> {
  /** 当前值 */
  state: Ref<T>;
  /** 当前快照索引 */
  snapshotIndex: Ref<number>;
  /** 快照数量 */
  snapshotCount: ComputedRef<number>;
  /** 创建快照 */
  takeSnapshot: () => void;
  /** 恢复到快照 */
  restoreSnapshot: (index: number) => void;
  /** 恢复到上一个快照 */
  restorePrevious: () => void;
  /** 清除所有快照 */
  clearSnapshots: () => void;
  /** 获取所有快照 */
  getSnapshots: () => T[];
}

/**
 * 状态快照
 * @example
 * const { state, takeSnapshot, restoreSnapshot, snapshotCount } = useSnapshot({
 *   items: []
 * })
 *
 * // 修改前创建快照
 * takeSnapshot()
 * state.value.items.push({ id: 1 })
 *
 * // 修改前再次创建快照
 * takeSnapshot()
 * state.value.items.push({ id: 2 })
 *
 * // 恢复到特定快照
 * restoreSnapshot(0) // 恢复到第一个快照
 */
export function useSnapshot<T>(initialState: T): UseSnapshotReturn<T> {
  const state = ref<T>(deepClone(initialState)) as Ref<T>;
  const snapshots = ref<T[]>([deepClone(initialState)]) as Ref<T[]>;
  const snapshotIndex = ref(0);

  const snapshotCount = computed(() => snapshots.value.length);

  const takeSnapshot = () => {
    snapshots.value.push(deepClone(toRaw(state.value)));
    snapshotIndex.value = snapshots.value.length - 1;
  };

  const restoreSnapshot = (index: number) => {
    if (index >= 0 && index < snapshots.value.length) {
      state.value = deepClone(snapshots.value[index]) as UnwrapRef<T>;
      snapshotIndex.value = index;
    }
  };

  const restorePrevious = () => {
    if (snapshotIndex.value > 0) {
      restoreSnapshot(snapshotIndex.value - 1);
    }
  };

  const clearSnapshots = () => {
    snapshots.value = [deepClone(toRaw(state.value))];
    snapshotIndex.value = 0;
  };

  const getSnapshots = () => {
    return snapshots.value.map((s) => deepClone(s));
  };

  return {
    state,
    snapshotIndex,
    snapshotCount,
    takeSnapshot,
    restoreSnapshot,
    restorePrevious,
    clearSnapshots,
    getSnapshots,
  };
}

// ============================================================================
// useSyncedRef - 同步的引用
// ============================================================================

export interface UseSyncedRefOptions {
  /** 同步延迟（毫秒） */
  delay?: number;
  /** 是否深度同步 */
  deep?: boolean;
}

export interface UseSyncedRefReturn<T> {
  /** 本地值 */
  local: Ref<T>;
  /** 是否正在同步 */
  isSyncing: Ref<boolean>;
  /** 强制同步 */
  forceSync: () => void;
}

/**
 * 同步的引用（延迟同步到外部值）
 * @example
 * const external = ref('Hello')
 * const { local, isSyncing } = useSyncedRef(external, { delay: 500 })
 *
 * // 本地修改
 * local.value = 'World'
 * // 500ms 后自动同步到 external
 */
export function useSyncedRef<T>(
  source: Ref<T>,
  options: UseSyncedRefOptions = {}
): UseSyncedRefReturn<T> {
  const { delay = 0, deep = true } = options;

  const local = ref<T>(unref(source)) as Ref<T>;
  const isSyncing = ref(false);

  let syncTimer: ReturnType<typeof setTimeout> | null = null;

  // 监听本地变化，同步到源
  watch(
    local,
    (newValue) => {
      if (syncTimer) {
        clearTimeout(syncTimer);
      }

      if (delay > 0) {
        isSyncing.value = true;
        syncTimer = setTimeout(() => {
          source.value = deepClone(toRaw(newValue));
          isSyncing.value = false;
        }, delay);
      } else {
        source.value = deepClone(toRaw(newValue));
      }
    },
    { deep }
  );

  // 监听源变化，同步到本地
  watch(
    source,
    (newValue) => {
      if (!isSyncing.value) {
        local.value = deepClone(toRaw(newValue)) as UnwrapRef<T>;
      }
    },
    { deep }
  );

  const forceSync = () => {
    if (syncTimer) {
      clearTimeout(syncTimer);
    }
    source.value = deepClone(toRaw(local.value));
    isSyncing.value = false;
  };

  return {
    local,
    isSyncing,
    forceSync,
  };
}

// ============================================================================
// useMemoize - 记忆化
// ============================================================================

export interface UseMemoizeReturn<T extends (...args: unknown[]) => unknown> {
  /** 记忆化函数 */
  memoized: T;
  /** 清除缓存 */
  clear: () => void;
  /** 删除特定缓存 */
  delete: (...args: Parameters<T>) => boolean;
  /** 检查是否有缓存 */
  has: (...args: Parameters<T>) => boolean;
  /** 缓存大小 */
  size: ComputedRef<number>;
}

/**
 * 记忆化函数
 * @example
 * const expensiveCalculation = (a: number, b: number) => {
 *   // 耗时计算
 *   return a * b
 * }
 *
 * const { memoized, clear, size } = useMemoize(expensiveCalculation)
 *
 * memoized(2, 3) // 计算并缓存
 * memoized(2, 3) // 从缓存返回
 *
 * clear() // 清除所有缓存
 */
export function useMemoize<T extends (...args: unknown[]) => unknown>(
  fn: T,
  keyResolver?: (...args: Parameters<T>) => string
): UseMemoizeReturn<T> {
  const cache = ref(new Map<string, ReturnType<T>>());

  const getKey = (...args: Parameters<T>): string => {
    if (keyResolver) {
      return keyResolver(...args);
    }
    return JSON.stringify(args);
  };

  const memoized = ((...args: Parameters<T>): ReturnType<T> => {
    const key = getKey(...args);

    if (cache.value.has(key)) {
      return cache.value.get(key)!;
    }

    const result = fn(...args) as ReturnType<T>;
    cache.value.set(key, result);
    return result;
  }) as T;

  const clear = () => {
    cache.value.clear();
  };

  const del = (...args: Parameters<T>): boolean => {
    const key = getKey(...args);
    return cache.value.delete(key);
  };

  const has = (...args: Parameters<T>): boolean => {
    const key = getKey(...args);
    return cache.value.has(key);
  };

  const size = computed(() => cache.value.size);

  return {
    memoized,
    clear,
    delete: del,
    has,
    size,
  };
}
