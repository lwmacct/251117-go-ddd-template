/**
 * Ref Composable
 * 提供增强的 ref 工具函数
 */

import {
  ref,
  shallowRef,
  customRef,
  watch,
  watchEffect,
  computed,
  onMounted,
  onBeforeUnmount,
  type Ref,
  type ShallowRef,
  type UnwrapRef,
  type ComponentPublicInstance,
} from "vue";

// ============================================================================
// 类型定义
// ============================================================================

/**
 * 防抖 ref 选项
 */
export interface RefDebouncedOptions {
  /** 防抖延迟（毫秒） */
  delay?: number;
  /** 最大等待时间 */
  maxWait?: number;
}

/**
 * 节流 ref 选项
 */
export interface RefThrottledOptions {
  /** 节流间隔（毫秒） */
  delay?: number;
  /** 是否在开始时触发 */
  leading?: boolean;
  /** 是否在结束时触发 */
  trailing?: boolean;
}

/**
 * 历史记录 ref 选项
 */
export interface RefHistoryOptions<T> {
  /** 历史记录容量 */
  capacity?: number;
  /** 是否深度克隆 */
  deep?: boolean;
  /** 克隆函数 */
  clone?: (value: T) => T;
  /** 防抖延迟 */
  debounce?: number;
}

/**
 * 历史记录 ref 返回值
 */
export interface RefHistoryReturn<T> {
  /** 当前值 */
  value: Ref<T>;
  /** 历史记录 */
  history: Ref<T[]>;
  /** 未来记录（撤销后的值） */
  future: Ref<T[]>;
  /** 是否可撤销 */
  canUndo: Ref<boolean>;
  /** 是否可重做 */
  canRedo: Ref<boolean>;
  /** 撤销 */
  undo: () => void;
  /** 重做 */
  redo: () => void;
  /** 清空历史 */
  clear: () => void;
  /** 暂停记录 */
  pause: () => void;
  /** 恢复记录 */
  resume: () => void;
  /** 是否暂停中 */
  isPaused: Ref<boolean>;
  /** 提交当前值到历史（手动模式） */
  commit: () => void;
}

/**
 * 自动重置 ref 选项
 */
export interface RefAutoResetOptions {
  /** 重置延迟（毫秒） */
  delay?: number;
}

/**
 * 可控 ref 选项
 */
export interface RefWithControlOptions<T> {
  /** 获取值时的拦截器 */
  onGet?: (value: T) => T;
  /** 设置值时的拦截器 */
  onSet?: (newValue: T, oldValue: T) => T;
  /** 设置前验证 */
  onBeforeSet?: (newValue: T, oldValue: T) => boolean;
}

/**
 * 可控 ref 返回值
 */
export interface RefWithControlReturn<T> {
  /** ref 值 */
  value: Ref<T>;
  /** 暂停响应 */
  pause: () => void;
  /** 恢复响应 */
  resume: () => void;
  /** 是否暂停中 */
  isPaused: Ref<boolean>;
  /** 静默设置（不触发响应） */
  silentSet: (value: T) => void;
  /** 获取原始值 */
  peek: () => T;
}

/**
 * 模板 ref 返回值
 */
export interface TemplateRefReturn<T> {
  /** ref 值 */
  ref: Ref<T | null>;
  /** 是否已挂载 */
  isMounted: Ref<boolean>;
  /** 等待挂载 */
  onMounted: (callback: (el: T) => void) => void;
}

// ============================================================================
// 核心函数
// ============================================================================

/**
 * 带默认值的 ref
 *
 * @description 当值为 null 或 undefined 时返回默认值
 *
 * @example
 * ```ts
 * const value = refDefault<string>(null, 'default')
 * console.log(value.value) // 'default'
 * value.value = 'hello'
 * console.log(value.value) // 'hello'
 * value.value = null
 * console.log(value.value) // 'default'
 * ```
 */
export function refDefault<T>(source: Ref<T | null | undefined>, defaultValue: T): Ref<T> {
  return computed({
    get: () => source.value ?? defaultValue,
    set: (value) => {
      source.value = value;
    },
  }) as Ref<T>;
}

/**
 * 防抖 ref
 *
 * @description 值变化后延迟更新
 *
 * @example
 * ```ts
 * const text = ref('')
 * const debouncedText = refDebounced(text, { delay: 300 })
 *
 * text.value = 'hello' // debouncedText 在 300ms 后变为 'hello'
 * ```
 */
export function refDebounced<T>(source: Ref<T>, options: RefDebouncedOptions = {}): Readonly<Ref<T>> {
  const { delay = 250, maxWait } = options;

  return customRef<T>((track, trigger) => {
    let value = source.value;
    let timer: ReturnType<typeof setTimeout> | null = null;
    let maxTimer: ReturnType<typeof setTimeout> | null = null;

    const update = () => {
      value = source.value;
      trigger();
      if (maxTimer) {
        clearTimeout(maxTimer);
        maxTimer = null;
      }
    };

    watch(source, () => {
      if (timer) clearTimeout(timer);

      if (maxWait && !maxTimer) {
        maxTimer = setTimeout(() => {
          if (timer) clearTimeout(timer);
          update();
        }, maxWait);
      }

      timer = setTimeout(update, delay);
    });

    return {
      get() {
        track();
        return value;
      },
      set() {
        // 只读
      },
    };
  });
}

/**
 * 节流 ref
 *
 * @description 限制值更新频率
 *
 * @example
 * ```ts
 * const scrollY = ref(0)
 * const throttledScrollY = refThrottled(scrollY, { delay: 100 })
 * ```
 */
export function refThrottled<T>(source: Ref<T>, options: RefThrottledOptions = {}): Readonly<Ref<T>> {
  const { delay = 100, leading = true, trailing = true } = options;

  return customRef<T>((track, trigger) => {
    let value = source.value;
    let lastExec = 0;
    let timer: ReturnType<typeof setTimeout> | null = null;

    const update = (newValue: T) => {
      value = newValue;
      lastExec = Date.now();
      trigger();
    };

    watch(source, (newValue) => {
      const now = Date.now();
      const elapsed = now - lastExec;

      if (elapsed >= delay) {
        if (timer) {
          clearTimeout(timer);
          timer = null;
        }
        if (leading || lastExec > 0) {
          update(newValue);
        } else {
          lastExec = now;
        }
      } else if (!timer && trailing) {
        timer = setTimeout(() => {
          timer = null;
          update(source.value);
        }, delay - elapsed);
      }
    });

    return {
      get() {
        track();
        return value;
      },
      set() {
        // 只读
      },
    };
  });
}

/**
 * 带历史记录的 ref
 *
 * @description 支持撤销/重做操作
 *
 * @example
 * ```ts
 * const { value, undo, redo, canUndo, canRedo, history } = refHistory(0, {
 *   capacity: 10
 * })
 *
 * value.value = 1
 * value.value = 2
 * console.log(history.value) // [0, 1]
 *
 * undo()
 * console.log(value.value) // 1
 *
 * redo()
 * console.log(value.value) // 2
 * ```
 */
export function refHistory<T>(initialValue: T, options: RefHistoryOptions<T> = {}): RefHistoryReturn<T> {
  const {
    capacity = 10,
    deep = false,
    clone = (v: T) => (deep ? JSON.parse(JSON.stringify(v)) : v),
    debounce = 0,
  } = options;

  const value = ref<T>(initialValue) as Ref<T>;
  const history = ref<T[]>([]) as Ref<T[]>;
  const future = ref<T[]>([]) as Ref<T[]>;
  const isPaused = ref(false);

  let timer: ReturnType<typeof setTimeout> | null = null;

  const commit = () => {
    if (isPaused.value) return;

    const cloned = clone(value.value);
    history.value = [...history.value.slice(-(capacity - 1)), cloned];
    future.value = [];
  };

  // 监听值变化
  watch(
    value,
    (newValue, oldValue) => {
      if (isPaused.value) return;

      if (timer) clearTimeout(timer);

      const record = () => {
        const cloned = clone(oldValue as T);
        history.value = [...history.value.slice(-(capacity - 1)), cloned];
        future.value = [];
      };

      if (debounce > 0) {
        timer = setTimeout(record, debounce);
      } else {
        record();
      }
    },
    { deep }
  );

  const canUndo = computed(() => history.value.length > 0);
  const canRedo = computed(() => future.value.length > 0);

  const undo = () => {
    if (!canUndo.value) return;

    const current = clone(value.value);
    future.value = [current, ...future.value];

    const previous = history.value.pop();
    if (previous !== undefined) {
      isPaused.value = true;
      value.value = previous;
      isPaused.value = false;
    }
  };

  const redo = () => {
    if (!canRedo.value) return;

    const current = clone(value.value);
    history.value = [...history.value, current];

    const next = future.value.shift();
    if (next !== undefined) {
      isPaused.value = true;
      value.value = next;
      isPaused.value = false;
    }
  };

  const clear = () => {
    history.value = [];
    future.value = [];
  };

  const pause = () => {
    isPaused.value = true;
  };

  const resume = () => {
    isPaused.value = false;
  };

  return {
    value,
    history,
    future,
    canUndo,
    canRedo,
    undo,
    redo,
    clear,
    pause,
    resume,
    isPaused,
    commit,
  };
}

/**
 * 自动重置 ref
 *
 * @description 值变化后自动重置为初始值
 *
 * @example
 * ```ts
 * const notification = refAutoReset('', { delay: 3000 })
 *
 * notification.value = '保存成功' // 3秒后自动清空
 * ```
 */
export function refAutoReset<T>(defaultValue: T, options: RefAutoResetOptions = {}): Ref<T> {
  const { delay = 1000 } = options;

  const value = ref<T>(defaultValue) as Ref<T>;
  let timer: ReturnType<typeof setTimeout> | null = null;

  watch(value, (newValue) => {
    if (newValue === defaultValue) return;

    if (timer) clearTimeout(timer);
    timer = setTimeout(() => {
      value.value = defaultValue;
    }, delay);
  });

  return value;
}

/**
 * 同步两个 ref
 *
 * @description 保持两个 ref 值同步
 *
 * @example
 * ```ts
 * const a = ref(0)
 * const b = ref(0)
 *
 * syncRefs(a, b)
 *
 * a.value = 1 // b.value 也变为 1
 * b.value = 2 // a.value 也变为 2
 * ```
 */
export function syncRefs<T>(
  source: Ref<T>,
  target: Ref<T>,
  options: { immediate?: boolean; direction?: "ltr" | "rtl" | "both" } = {}
): () => void {
  const { immediate = true, direction = "both" } = options;

  const stops: (() => void)[] = [];

  if (direction === "ltr" || direction === "both") {
    stops.push(
      watch(
        source,
        (value) => {
          target.value = value;
        },
        { immediate }
      )
    );
  }

  if (direction === "rtl" || direction === "both") {
    stops.push(
      watch(
        target,
        (value) => {
          source.value = value;
        },
        { immediate: direction === "rtl" && immediate }
      )
    );
  }

  return () => stops.forEach((stop) => stop());
}

/**
 * 可控 ref
 *
 * @description 支持拦截 get/set 操作
 *
 * @example
 * ```ts
 * const { value, pause, resume, silentSet, peek } = refWithControl(0, {
 *   onSet: (newValue) => Math.max(0, newValue), // 确保非负
 *   onGet: (value) => value * 2 // 读取时翻倍
 * })
 * ```
 */
export function refWithControl<T>(initialValue: T, options: RefWithControlOptions<T> = {}): RefWithControlReturn<T> {
  const { onGet, onSet, onBeforeSet } = options;

  let internalValue = initialValue;
  const isPaused = ref(false);

  const value = customRef<T>((track, trigger) => ({
    get() {
      track();
      const val = internalValue;
      return onGet ? onGet(val) : val;
    },
    set(newValue) {
      if (isPaused.value) return;

      if (onBeforeSet && !onBeforeSet(newValue, internalValue)) {
        return;
      }

      const processedValue = onSet ? onSet(newValue, internalValue) : newValue;
      internalValue = processedValue;
      trigger();
    },
  }));

  const pause = () => {
    isPaused.value = true;
  };

  const resume = () => {
    isPaused.value = false;
  };

  const silentSet = (newValue: T) => {
    internalValue = newValue;
  };

  const peek = () => internalValue;

  return {
    value,
    pause,
    resume,
    isPaused,
    silentSet,
    peek,
  };
}

/**
 * 模板 ref
 *
 * @description 用于获取模板中的 DOM 元素或组件实例
 *
 * @example
 * ```ts
 * const { ref: inputRef, isMounted, onMounted } = templateRef<HTMLInputElement>()
 *
 * onMounted((el) => {
 *   el.focus()
 * })
 *
 * // 在模板中: <input ref="inputRef" />
 * ```
 */
export function templateRef<T extends HTMLElement | ComponentPublicInstance>(): TemplateRefReturn<T> {
  const elementRef = ref<T | null>(null);
  const isMounted = ref(false);
  const callbacks: ((el: T) => void)[] = [];

  watch(elementRef, (el) => {
    if (el) {
      isMounted.value = true;
      callbacks.forEach((cb) => cb(el));
    } else {
      isMounted.value = false;
    }
  });

  return {
    ref: elementRef,
    isMounted,
    onMounted: (callback) => {
      if (elementRef.value) {
        callback(elementRef.value);
      } else {
        callbacks.push(callback);
      }
    },
  };
}

/**
 * 获取上一个值
 *
 * @description 保存 ref 的上一个值
 *
 * @example
 * ```ts
 * const count = ref(0)
 * const prevCount = usePrevious(count)
 *
 * count.value = 1
 * console.log(prevCount.value) // 0
 *
 * count.value = 2
 * console.log(prevCount.value) // 1
 * ```
 */
export function usePrevious<T>(source: Ref<T>): Readonly<Ref<T | undefined>> {
  const previous = ref<T | undefined>(undefined);

  watch(
    source,
    (_, oldValue) => {
      previous.value = oldValue;
    },
    { flush: "sync" }
  );

  return previous as Readonly<Ref<T | undefined>>;
}

/**
 * 获取最新值
 *
 * @description 始终返回最新值的函数，不触发依赖追踪
 *
 * @example
 * ```ts
 * const count = ref(0)
 * const getLatest = useLatest(count)
 *
 * // 在 effect 中使用不会被追踪
 * watchEffect(() => {
 *   console.log(getLatest()) // 不会触发重新执行
 * })
 * ```
 */
export function useLatest<T>(source: Ref<T>): () => T {
  const latest = shallowRef(source.value);

  watch(source, (value) => {
    latest.value = value;
  });

  return () => latest.value;
}

/**
 * 锁定 ref
 *
 * @description 创建一个可锁定的 ref，锁定后无法修改
 *
 * @example
 * ```ts
 * const { value, lock, unlock, isLocked } = refLocked(0)
 *
 * value.value = 1 // 成功
 * lock()
 * value.value = 2 // 无效，值仍为 1
 * unlock()
 * value.value = 2 // 成功
 * ```
 */
export function refLocked<T>(initialValue: T): {
  value: Ref<T>;
  lock: () => void;
  unlock: () => void;
  isLocked: Ref<boolean>;
} {
  const isLocked = ref(false);
  let internalValue = initialValue;

  const value = customRef<T>((track, trigger) => ({
    get() {
      track();
      return internalValue;
    },
    set(newValue) {
      if (isLocked.value) return;
      internalValue = newValue;
      trigger();
    },
  }));

  return {
    value,
    lock: () => {
      isLocked.value = true;
    },
    unlock: () => {
      isLocked.value = false;
    },
    isLocked,
  };
}

/**
 * 计数器 ref
 *
 * @description 创建带有 inc/dec/reset 方法的计数器
 *
 * @example
 * ```ts
 * const { count, inc, dec, reset, set } = useCounter(0, {
 *   min: 0,
 *   max: 10
 * })
 *
 * inc() // count = 1
 * inc(5) // count = 6
 * dec(2) // count = 4
 * reset() // count = 0
 * ```
 */
export function useCounter(
  initialValue = 0,
  options: { min?: number; max?: number } = {}
): {
  count: Ref<number>;
  inc: (delta?: number) => void;
  dec: (delta?: number) => void;
  reset: () => void;
  set: (value: number) => void;
} {
  const { min = -Infinity, max = Infinity } = options;

  const count = ref(initialValue);

  const clamp = (value: number) => Math.min(max, Math.max(min, value));

  return {
    count,
    inc: (delta = 1) => {
      count.value = clamp(count.value + delta);
    },
    dec: (delta = 1) => {
      count.value = clamp(count.value - delta);
    },
    reset: () => {
      count.value = initialValue;
    },
    set: (value: number) => {
      count.value = clamp(value);
    },
  };
}

/**
 * 布尔 ref 切换
 *
 * @description 创建带有 toggle/setTrue/setFalse 方法的布尔 ref
 *
 * @example
 * ```ts
 * const { value, toggle, setTrue, setFalse } = useBoolean(false)
 *
 * toggle() // true
 * toggle() // false
 * setTrue() // true
 * setFalse() // false
 * ```
 */
export function useBoolean(initialValue = false): {
  value: Ref<boolean>;
  toggle: () => void;
  setTrue: () => void;
  setFalse: () => void;
} {
  const value = ref(initialValue);

  return {
    value,
    toggle: () => {
      value.value = !value.value;
    },
    setTrue: () => {
      value.value = true;
    },
    setFalse: () => {
      value.value = false;
    },
  };
}

/**
 * 对象 ref
 *
 * @description 创建带有便捷方法的对象 ref
 *
 * @example
 * ```ts
 * const { state, set, merge, reset, patch } = useObject({ name: '', age: 0 })
 *
 * set({ name: 'John', age: 30 })
 * merge({ name: 'Jane' }) // { name: 'Jane', age: 30 }
 * patch('name', 'Bob') // { name: 'Bob', age: 30 }
 * reset() // { name: '', age: 0 }
 * ```
 */
export function useObject<T extends Record<string, unknown>>(
  initialValue: T
): {
  state: Ref<T>;
  set: (value: T) => void;
  merge: (partial: Partial<T>) => void;
  reset: () => void;
  patch: <K extends keyof T>(key: K, value: T[K]) => void;
} {
  const state = ref<T>({ ...initialValue }) as Ref<T>;

  return {
    state,
    set: (value) => {
      state.value = value;
    },
    merge: (partial) => {
      state.value = { ...state.value, ...partial };
    },
    reset: () => {
      state.value = { ...initialValue };
    },
    patch: (key, value) => {
      state.value = { ...state.value, [key]: value };
    },
  };
}

/**
 * 数组 ref
 *
 * @description 创建带有便捷方法的数组 ref
 *
 * @example
 * ```ts
 * const { array, push, pop, shift, unshift, remove, clear, set } = useArray([1, 2, 3])
 *
 * push(4) // [1, 2, 3, 4]
 * pop() // [1, 2, 3]
 * remove(1) // [1, 3]
 * clear() // []
 * ```
 */
export function useArray<T>(initialValue: T[] = []): {
  array: Ref<T[]>;
  push: (...items: T[]) => void;
  pop: () => T | undefined;
  shift: () => T | undefined;
  unshift: (...items: T[]) => void;
  remove: (index: number) => void;
  removeItem: (item: T) => void;
  clear: () => void;
  set: (value: T[]) => void;
  insert: (index: number, item: T) => void;
  update: (index: number, item: T) => void;
} {
  const array = ref<T[]>([...initialValue]) as Ref<T[]>;

  return {
    array,
    push: (...items) => {
      array.value = [...array.value, ...items];
    },
    pop: () => {
      const item = array.value[array.value.length - 1];
      array.value = array.value.slice(0, -1);
      return item;
    },
    shift: () => {
      const item = array.value[0];
      array.value = array.value.slice(1);
      return item;
    },
    unshift: (...items) => {
      array.value = [...items, ...array.value];
    },
    remove: (index) => {
      array.value = [...array.value.slice(0, index), ...array.value.slice(index + 1)];
    },
    removeItem: (item) => {
      array.value = array.value.filter((i) => i !== item);
    },
    clear: () => {
      array.value = [];
    },
    set: (value) => {
      array.value = value;
    },
    insert: (index, item) => {
      array.value = [...array.value.slice(0, index), item, ...array.value.slice(index)];
    },
    update: (index, item) => {
      array.value = array.value.map((v, i) => (i === index ? item : v));
    },
  };
}

/**
 * 集合 ref
 *
 * @description 创建带有便捷方法的 Set ref
 *
 * @example
 * ```ts
 * const { set, add, remove, has, clear, toggle } = useSet([1, 2, 3])
 *
 * add(4) // Set(1, 2, 3, 4)
 * remove(2) // Set(1, 3, 4)
 * toggle(1) // Set(3, 4)
 * has(3) // true
 * ```
 */
export function useSet<T>(initialValue: Iterable<T> = []): {
  set: Ref<Set<T>>;
  add: (value: T) => void;
  remove: (value: T) => void;
  has: (value: T) => boolean;
  clear: () => void;
  toggle: (value: T) => void;
  values: () => T[];
} {
  const set = ref(new Set(initialValue)) as Ref<Set<T>>;

  return {
    set,
    add: (value) => {
      set.value = new Set([...set.value, value]);
    },
    remove: (value) => {
      const newSet = new Set(set.value);
      newSet.delete(value);
      set.value = newSet;
    },
    has: (value) => set.value.has(value),
    clear: () => {
      set.value = new Set();
    },
    toggle: (value) => {
      if (set.value.has(value)) {
        const newSet = new Set(set.value);
        newSet.delete(value);
        set.value = newSet;
      } else {
        set.value = new Set([...set.value, value]);
      }
    },
    values: () => [...set.value],
  };
}

/**
 * Map ref
 *
 * @description 创建带有便捷方法的 Map ref
 *
 * @example
 * ```ts
 * const { map, set, get, remove, has, clear } = useMap([['a', 1]])
 *
 * set('b', 2)
 * get('a') // 1
 * remove('a')
 * has('a') // false
 * ```
 */
export function useMap<K, V>(
  initialValue: Iterable<[K, V]> = []
): {
  map: Ref<Map<K, V>>;
  set: (key: K, value: V) => void;
  get: (key: K) => V | undefined;
  remove: (key: K) => void;
  has: (key: K) => boolean;
  clear: () => void;
  keys: () => K[];
  values: () => V[];
  entries: () => [K, V][];
} {
  const map = ref(new Map(initialValue)) as Ref<Map<K, V>>;

  return {
    map,
    set: (key, value) => {
      const newMap = new Map(map.value);
      newMap.set(key, value);
      map.value = newMap;
    },
    get: (key) => map.value.get(key),
    remove: (key) => {
      const newMap = new Map(map.value);
      newMap.delete(key);
      map.value = newMap;
    },
    has: (key) => map.value.has(key),
    clear: () => {
      map.value = new Map();
    },
    keys: () => [...map.value.keys()],
    values: () => [...map.value.values()],
    entries: () => [...map.value.entries()],
  };
}
