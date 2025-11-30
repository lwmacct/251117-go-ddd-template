/**
 * VModel Composable
 * 提供 v-model 双向绑定的工具函数
 */

import { ref, computed, watch, getCurrentInstance, type Ref, type WritableComputedRef, type ModelRef } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UseVModelOptions<T> {
  /** 是否深度监听 */
  deep?: boolean;
  /** 默认值 */
  defaultValue?: T;
  /** 值变化时的回调 */
  onChange?: (value: T) => void;
  /** 是否使用被动模式（不触发更新事件） */
  passive?: boolean;
  /** 自定义克隆函数 */
  clone?: (value: T) => T;
}

// ============================================================================
// useVModel - 双向绑定
// ============================================================================

/**
 * 简化 v-model 双向绑定
 * @example
 * // 父组件
 * <MyInput v-model="name" />
 *
 * // 子组件
 * const props = defineProps<{ modelValue: string }>()
 * const emit = defineEmits<{ 'update:modelValue': [value: string] }>()
 *
 * const model = useVModel(props, 'modelValue', emit)
 * // 直接使用 model.value 进行双向绑定
 */
export function useVModel<P extends object, K extends keyof P>(
  props: P,
  key: K,
  emit?: (name: string, ...args: unknown[]) => void,
  options: UseVModelOptions<P[K]> = {}
): Ref<P[K]> {
  const { deep = false, defaultValue, onChange, passive = false, clone } = options;

  // 获取当前组件实例
  const vm = getCurrentInstance();
  const emitFn = emit || vm?.emit;

  if (!emitFn) {
    console.warn("useVModel: emit function not found");
  }

  // 事件名
  const eventName = `update:${String(key)}`;

  // 获取初始值
  const getValue = () => {
    const value = props[key];
    if (value === undefined && defaultValue !== undefined) {
      return defaultValue;
    }
    return value;
  };

  // 被动模式：使用本地 ref
  if (passive) {
    const local = ref(getValue()) as Ref<P[K]>;

    watch(
      () => props[key],
      (newValue) => {
        if (newValue !== local.value) {
          local.value = clone ? clone(newValue) : newValue;
        }
      },
      { deep }
    );

    return local;
  }

  // 主动模式：使用 computed
  return computed({
    get() {
      return getValue();
    },
    set(value) {
      if (clone) {
        value = clone(value);
      }
      emitFn?.(eventName, value);
      onChange?.(value);
    },
  }) as WritableComputedRef<P[K]>;
}

// ============================================================================
// useVModels - 多个 v-model
// ============================================================================

type UseVModelsReturn<P extends object, K extends keyof P> = {
  [Key in K]: Ref<P[Key]>;
};

/**
 * 处理多个 v-model
 * @example
 * const props = defineProps<{
 *   firstName: string
 *   lastName: string
 *   age: number
 * }>()
 * const emit = defineEmits(['update:firstName', 'update:lastName', 'update:age'])
 *
 * const { firstName, lastName, age } = useVModels(props, emit)
 */
export function useVModels<P extends object, K extends keyof P = keyof P>(
  props: P,
  emit?: (name: string, ...args: unknown[]) => void,
  keys?: K[],
  options?: UseVModelOptions<P[K]>
): UseVModelsReturn<P, K> {
  const result = {} as UseVModelsReturn<P, K>;
  const keysToProcess = keys || (Object.keys(props) as K[]);

  for (const key of keysToProcess) {
    result[key] = useVModel(props, key, emit, options) as Ref<P[typeof key]>;
  }

  return result;
}

// ============================================================================
// useModelValue - 简化的 modelValue 处理
// ============================================================================

/**
 * 简化的 modelValue 处理
 * @example
 * const props = defineProps<{ modelValue: string }>()
 * const emit = defineEmits<{ 'update:modelValue': [value: string] }>()
 *
 * const model = useModelValue(props, emit)
 */
export function useModelValue<T>(
  props: { modelValue: T },
  emit?: (name: "update:modelValue", value: T) => void,
  options?: UseVModelOptions<T>
): Ref<T> {
  return useVModel(props, "modelValue", emit as (name: string, ...args: unknown[]) => void, options);
}

// ============================================================================
// useProxyModel - 代理 model（支持本地修改）
// ============================================================================

export interface UseProxyModelOptions<T> {
  /** 是否深度监听 */
  deep?: boolean;
  /** 自定义克隆函数 */
  clone?: (value: T) => T;
  /** 值变化时的回调 */
  onChange?: (value: T) => void;
}

export interface UseProxyModelReturn<T> {
  /** 代理值 */
  proxy: Ref<T>;
  /** 是否已修改 */
  isModified: Ref<boolean>;
  /** 同步代理值到源值 */
  sync: () => void;
  /** 重置为源值 */
  reset: () => void;
}

/**
 * 代理 model（允许本地修改，手动同步）
 * @example
 * const props = defineProps<{ modelValue: User }>()
 * const emit = defineEmits<{ 'update:modelValue': [value: User] }>()
 *
 * const { proxy, isModified, sync, reset } = useProxyModel(props, 'modelValue', emit)
 *
 * // 本地修改
 * proxy.value.name = 'Jane'
 *
 * // 保存时同步
 * const handleSave = () => {
 *   sync()
 * }
 *
 * // 取消时重置
 * const handleCancel = () => {
 *   reset()
 * }
 */
export function useProxyModel<P extends object, K extends keyof P>(
  props: P,
  key: K,
  emit?: (name: string, ...args: unknown[]) => void,
  options: UseProxyModelOptions<P[K]> = {}
): UseProxyModelReturn<P[K]> {
  const { deep = true, clone = (v) => JSON.parse(JSON.stringify(v)), onChange } = options;

  const vm = getCurrentInstance();
  const emitFn = emit || vm?.emit;
  const eventName = `update:${String(key)}`;

  const original = ref(clone(props[key])) as Ref<P[K]>;
  const proxy = ref(clone(props[key])) as Ref<P[K]>;
  const isModified = ref(false);

  // 监听 props 变化
  watch(
    () => props[key],
    (newValue) => {
      const cloned = clone(newValue);
      original.value = cloned;
      if (!isModified.value) {
        proxy.value = clone(cloned);
      }
    },
    { deep }
  );

  // 监听 proxy 变化
  watch(
    proxy,
    (newValue) => {
      isModified.value = JSON.stringify(newValue) !== JSON.stringify(original.value);
      onChange?.(newValue);
    },
    { deep }
  );

  const sync = () => {
    emitFn?.(eventName, clone(proxy.value));
    original.value = clone(proxy.value);
    isModified.value = false;
  };

  const reset = () => {
    proxy.value = clone(original.value);
    isModified.value = false;
  };

  return {
    proxy,
    isModified,
    sync,
    reset,
  };
}

// ============================================================================
// useControlled - 受控/非受控组件
// ============================================================================

export interface UseControlledOptions<T> {
  /** 默认值（非受控模式） */
  defaultValue: T;
  /** 值变化回调 */
  onChange?: (value: T) => void;
}

export interface UseControlledReturn<T> {
  /** 当前值 */
  value: Ref<T>;
  /** 是否受控 */
  isControlled: boolean;
  /** 设置值 */
  setValue: (value: T) => void;
}

/**
 * 支持受控/非受控模式
 * @example
 * const props = defineProps<{
 *   value?: string
 *   defaultValue?: string
 * }>()
 * const emit = defineEmits<{ 'update:value': [value: string] }>()
 *
 * const { value, isControlled, setValue } = useControlled(props, 'value', emit, {
 *   defaultValue: ''
 * })
 *
 * // 无论受控还是非受控，都使用 value 和 setValue
 */
export function useControlled<P extends object, K extends keyof P>(
  props: P,
  key: K,
  emit?: (name: string, ...args: unknown[]) => void,
  options: UseControlledOptions<NonNullable<P[K]>> = {} as UseControlledOptions<NonNullable<P[K]>>
): UseControlledReturn<NonNullable<P[K]>> {
  const { defaultValue, onChange } = options;

  const vm = getCurrentInstance();
  const emitFn = emit || vm?.emit;
  const eventName = `update:${String(key)}`;

  // 判断是否受控
  const isControlled = props[key] !== undefined;

  // 非受控模式使用本地状态
  const localValue = ref(isControlled ? props[key] : defaultValue) as Ref<NonNullable<P[K]>>;

  // 受控模式监听 props 变化
  if (isControlled) {
    watch(
      () => props[key],
      (newValue) => {
        if (newValue !== undefined) {
          localValue.value = newValue as NonNullable<P[K]>;
        }
      }
    );
  }

  const setValue = (newValue: NonNullable<P[K]>) => {
    if (isControlled) {
      emitFn?.(eventName, newValue);
    } else {
      localValue.value = newValue;
    }
    onChange?.(newValue);
  };

  return {
    value: localValue,
    isControlled,
    setValue,
  };
}

// ============================================================================
// useDebouncedVModel - 防抖的 v-model
// ============================================================================

export interface UseDebouncedVModelOptions<T> extends UseVModelOptions<T> {
  /** 防抖延迟（毫秒） */
  debounce?: number;
}

/**
 * 防抖的 v-model
 * @example
 * const props = defineProps<{ modelValue: string }>()
 * const emit = defineEmits<{ 'update:modelValue': [value: string] }>()
 *
 * const model = useDebouncedVModel(props, 'modelValue', emit, {
 *   debounce: 300
 * })
 *
 * // 输入时会防抖触发更新
 */
export function useDebouncedVModel<P extends object, K extends keyof P>(
  props: P,
  key: K,
  emit?: (name: string, ...args: unknown[]) => void,
  options: UseDebouncedVModelOptions<P[K]> = {}
): Ref<P[K]> {
  const { debounce = 300, onChange, ...restOptions } = options;

  const vm = getCurrentInstance();
  const emitFn = emit || vm?.emit;
  const eventName = `update:${String(key)}`;

  // 本地值
  const local = ref(props[key]) as Ref<P[K]>;
  let timer: ReturnType<typeof setTimeout> | null = null;

  // 监听 props 变化
  watch(
    () => props[key],
    (newValue) => {
      local.value = newValue;
    }
  );

  // 监听本地值变化，防抖更新
  watch(
    local,
    (newValue) => {
      if (timer) {
        clearTimeout(timer);
      }

      timer = setTimeout(() => {
        emitFn?.(eventName, newValue);
        onChange?.(newValue);
      }, debounce);
    },
    { deep: restOptions.deep }
  );

  return local;
}

// ============================================================================
// useThrottledVModel - 节流的 v-model
// ============================================================================

export interface UseThrottledVModelOptions<T> extends UseVModelOptions<T> {
  /** 节流间隔（毫秒） */
  throttle?: number;
}

/**
 * 节流的 v-model
 * @example
 * const props = defineProps<{ modelValue: number }>()
 * const emit = defineEmits<{ 'update:modelValue': [value: number] }>()
 *
 * const model = useThrottledVModel(props, 'modelValue', emit, {
 *   throttle: 100
 * })
 *
 * // 频繁更新时会节流
 */
export function useThrottledVModel<P extends object, K extends keyof P>(
  props: P,
  key: K,
  emit?: (name: string, ...args: unknown[]) => void,
  options: UseThrottledVModelOptions<P[K]> = {}
): Ref<P[K]> {
  const { throttle = 100, onChange, ...restOptions } = options;

  const vm = getCurrentInstance();
  const emitFn = emit || vm?.emit;
  const eventName = `update:${String(key)}`;

  // 本地值
  const local = ref(props[key]) as Ref<P[K]>;
  let lastUpdate = 0;
  let timer: ReturnType<typeof setTimeout> | null = null;

  // 监听 props 变化
  watch(
    () => props[key],
    (newValue) => {
      local.value = newValue;
    }
  );

  // 监听本地值变化，节流更新
  watch(
    local,
    (newValue) => {
      const now = Date.now();

      if (now - lastUpdate >= throttle) {
        lastUpdate = now;
        emitFn?.(eventName, newValue);
        onChange?.(newValue);
      } else {
        // 确保最后一次更新被发送
        if (timer) {
          clearTimeout(timer);
        }
        timer = setTimeout(
          () => {
            lastUpdate = Date.now();
            emitFn?.(eventName, newValue);
            onChange?.(newValue);
          },
          throttle - (now - lastUpdate)
        );
      }
    },
    { deep: restOptions.deep }
  );

  return local;
}

// ============================================================================
// useToggle - 切换值
// ============================================================================

export interface UseToggleOptions {
  /** 真值 */
  truthyValue?: unknown;
  /** 假值 */
  falsyValue?: unknown;
}

export interface UseToggleReturn {
  /** 当前值 */
  value: Ref<boolean>;
  /** 切换 */
  toggle: (value?: boolean) => void;
  /** 设置为真 */
  setTrue: () => void;
  /** 设置为假 */
  setFalse: () => void;
}

/**
 * 切换值
 * @example
 * const { value, toggle, setTrue, setFalse } = useToggle(false)
 *
 * toggle() // true
 * toggle() // false
 * setTrue() // true
 * setFalse() // false
 */
export function useToggle(initialValue: boolean = false, options: UseToggleOptions = {}): UseToggleReturn {
  const { truthyValue = true, falsyValue = false } = options;

  const value = ref(initialValue);

  const toggle = (newValue?: boolean) => {
    if (typeof newValue === "boolean") {
      value.value = newValue;
    } else {
      value.value = !value.value;
    }
  };

  const setTrue = () => {
    value.value = true;
  };

  const setFalse = () => {
    value.value = false;
  };

  return {
    value,
    toggle,
    setTrue,
    setFalse,
  };
}

// ============================================================================
// useCycleList - 循环列表
// ============================================================================

export interface UseCycleListReturn<T> {
  /** 当前值 */
  value: Ref<T>;
  /** 当前索引 */
  index: Ref<number>;
  /** 下一个 */
  next: () => void;
  /** 上一个 */
  prev: () => void;
  /** 跳转到指定索引 */
  go: (index: number) => void;
}

/**
 * 循环列表
 * @example
 * const themes = ['light', 'dark', 'auto']
 * const { value, next, prev, go, index } = useCycleList(themes)
 *
 * next() // 'dark'
 * next() // 'auto'
 * next() // 'light'
 * prev() // 'auto'
 */
export function useCycleList<T>(list: T[], initialIndex: number = 0): UseCycleListReturn<T> {
  const index = ref(initialIndex);
  const value = computed({
    get: () => list[index.value],
    set: (newValue) => {
      const idx = list.indexOf(newValue);
      if (idx !== -1) {
        index.value = idx;
      }
    },
  });

  const next = () => {
    index.value = (index.value + 1) % list.length;
  };

  const prev = () => {
    index.value = (index.value - 1 + list.length) % list.length;
  };

  const go = (i: number) => {
    index.value = ((i % list.length) + list.length) % list.length;
  };

  return {
    value: value as Ref<T>,
    index,
    next,
    prev,
    go,
  };
}
