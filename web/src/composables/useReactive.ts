/**
 * Reactive Composable
 * 提供增强的 reactive 工具函数
 */

import {
  reactive,
  readonly,
  toRaw,
  isReactive,
  isReadonly,
  isProxy,
  markRaw,
  toRef,
  toRefs,
  shallowReactive,
  watch,
  computed,
  ref,
  type UnwrapNestedRefs,
  type Ref,
  type DeepReadonly,
} from "vue";

// ============================================================================
// 类型定义
// ============================================================================

/**
 * 响应式选项
 */
export interface ReactiveWithOptionsConfig<T extends object> {
  /** 是否深度响应式 */
  deep?: boolean;
  /** 是否只读 */
  readonly?: boolean;
  /** 初始化回调 */
  onInit?: (state: T) => void;
  /** 变化回调 */
  onChange?: (newValue: T, oldValue: T) => void;
}

/**
 * 响应式历史选项
 */
export interface ReactiveHistoryOptions {
  /** 历史记录容量 */
  capacity?: number;
  /** 是否深度克隆 */
  deep?: boolean;
  /** 自定义克隆函数 */
  clone?: <T>(value: T) => T;
}

/**
 * 响应式历史返回值
 */
export interface ReactiveHistoryReturn<T extends object> {
  /** 响应式状态 */
  state: UnwrapNestedRefs<T>;
  /** 历史记录 */
  history: Ref<T[]>;
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
  /** 提交快照 */
  commit: () => void;
}

/**
 * 响应式表单返回值
 */
export interface ReactiveFormReturn<T extends object> {
  /** 表单数据 */
  data: UnwrapNestedRefs<T>;
  /** 是否已修改 */
  isDirty: Ref<boolean>;
  /** 是否正在提交 */
  isSubmitting: Ref<boolean>;
  /** 重置表单 */
  reset: () => void;
  /** 获取变更 */
  getChanges: () => Partial<T>;
  /** 应用变更 */
  apply: (changes: Partial<T>) => void;
  /** 提交表单 */
  submit: (handler: (data: T) => Promise<void>) => Promise<void>;
}

/**
 * 响应式 Pick 类型
 */
export type ReactivePick<T extends object, K extends keyof T> = {
  [P in K]: T[P];
};

/**
 * 响应式 Omit 类型
 */
export type ReactiveOmit<T extends object, K extends keyof T> = {
  [P in Exclude<keyof T, K>]: T[P];
};

// ============================================================================
// 核心函数
// ============================================================================

/**
 * 带选项的响应式对象
 *
 * @description 创建带有额外选项的响应式对象
 *
 * @example
 * ```ts
 * const state = reactiveWithOptions({ count: 0 }, {
 *   readonly: false,
 *   onChange: (newVal, oldVal) => console.log('changed')
 * })
 * ```
 */
export function reactiveWithOptions<T extends object>(
  initial: T,
  options: ReactiveWithOptionsConfig<T> = {}
): UnwrapNestedRefs<T> {
  const { deep = true, readonly: isReadonly = false, onInit, onChange } = options;

  const state = deep ? reactive(initial) : shallowReactive(initial);

  if (onInit) {
    onInit(state as T);
  }

  if (onChange) {
    watch(
      () => ({ ...state }),
      (newValue, oldValue) => {
        onChange(newValue as T, oldValue as T);
      },
      { deep }
    );
  }

  return isReadonly ? (readonly(state) as UnwrapNestedRefs<T>) : (state as UnwrapNestedRefs<T>);
}

/**
 * 带历史记录的响应式对象
 *
 * @description 创建支持撤销/重做的响应式对象
 *
 * @example
 * ```ts
 * const { state, undo, redo, canUndo, canRedo } = reactiveHistory({
 *   name: '',
 *   age: 0
 * }, { capacity: 50 })
 *
 * state.name = 'John'
 * state.age = 30
 * undo() // 撤销 age 变更
 * ```
 */
export function reactiveHistory<T extends object>(
  initial: T,
  options: ReactiveHistoryOptions = {}
): ReactiveHistoryReturn<T> {
  const {
    capacity = 10,
    deep = true,
    clone = (v) => JSON.parse(JSON.stringify(v)),
  } = options;

  const state = reactive(clone(initial)) as UnwrapNestedRefs<T>;
  const history = ref<T[]>([]) as Ref<T[]>;
  const future = ref<T[]>([]) as Ref<T[]>;

  const canUndo = computed(() => history.value.length > 0);
  const canRedo = computed(() => future.value.length > 0);

  // 提交当前快照到历史
  const commit = () => {
    const snapshot = clone(toRaw(state) as T);
    history.value = [...history.value.slice(-(capacity - 1)), snapshot];
    future.value = [];
  };

  // 监听变化，自动记录历史
  let skipWatch = false;
  watch(
    () => ({ ...state }),
    (_, oldValue) => {
      if (skipWatch) return;
      const snapshot = clone(oldValue as T);
      history.value = [...history.value.slice(-(capacity - 1)), snapshot];
      future.value = [];
    },
    { deep }
  );

  const undo = () => {
    if (!canUndo.value) return;

    const current = clone(toRaw(state) as T);
    future.value = [current, ...future.value];

    const previous = history.value.pop();
    if (previous) {
      skipWatch = true;
      Object.assign(state, previous);
      skipWatch = false;
    }
  };

  const redo = () => {
    if (!canRedo.value) return;

    const current = clone(toRaw(state) as T);
    history.value = [...history.value, current];

    const next = future.value.shift();
    if (next) {
      skipWatch = true;
      Object.assign(state, next);
      skipWatch = false;
    }
  };

  const clear = () => {
    history.value = [];
    future.value = [];
  };

  return {
    state,
    history,
    canUndo,
    canRedo,
    undo,
    redo,
    clear,
    commit,
  };
}

/**
 * 响应式表单
 *
 * @description 创建带有表单功能的响应式对象
 *
 * @example
 * ```ts
 * const { data, isDirty, reset, getChanges, submit } = reactiveForm({
 *   username: '',
 *   email: ''
 * })
 *
 * data.username = 'john'
 * console.log(isDirty.value) // true
 *
 * const changes = getChanges() // { username: 'john' }
 *
 * await submit(async (formData) => {
 *   await api.save(formData)
 * })
 * ```
 */
export function reactiveForm<T extends object>(
  initial: T
): ReactiveFormReturn<T> {
  const initialSnapshot = JSON.parse(JSON.stringify(initial));
  const data = reactive(JSON.parse(JSON.stringify(initial))) as UnwrapNestedRefs<T>;
  const isSubmitting = ref(false);

  const isDirty = computed(() => {
    const current = toRaw(data);
    return JSON.stringify(current) !== JSON.stringify(initialSnapshot);
  });

  const reset = () => {
    Object.assign(data, JSON.parse(JSON.stringify(initialSnapshot)));
  };

  const getChanges = (): Partial<T> => {
    const current = toRaw(data) as T;
    const changes: Partial<T> = {};

    for (const key in current) {
      if (JSON.stringify(current[key]) !== JSON.stringify(initialSnapshot[key])) {
        changes[key] = current[key];
      }
    }

    return changes;
  };

  const apply = (changes: Partial<T>) => {
    Object.assign(data, changes);
  };

  const submit = async (handler: (formData: T) => Promise<void>) => {
    isSubmitting.value = true;
    try {
      await handler(toRaw(data) as T);
      // 提交成功后更新初始快照
      Object.assign(initialSnapshot, JSON.parse(JSON.stringify(toRaw(data))));
    } finally {
      isSubmitting.value = false;
    }
  };

  return {
    data,
    isDirty,
    isSubmitting,
    reset,
    getChanges,
    apply,
    submit,
  };
}

/**
 * 响应式选取
 *
 * @description 从响应式对象中选取指定属性
 *
 * @example
 * ```ts
 * const user = reactive({ id: 1, name: 'John', email: 'john@example.com' })
 * const picked = reactivePick(user, ['id', 'name'])
 * // picked = reactive({ id: 1, name: 'John' })
 * ```
 */
export function reactivePick<T extends object, K extends keyof T>(
  source: T,
  keys: K[]
): UnwrapNestedRefs<Pick<T, K>> {
  const picked = {} as Pick<T, K>;
  for (const key of keys) {
    picked[key] = source[key];
  }
  return reactive(picked) as UnwrapNestedRefs<Pick<T, K>>;
}

/**
 * 响应式排除
 *
 * @description 从响应式对象中排除指定属性
 *
 * @example
 * ```ts
 * const user = reactive({ id: 1, name: 'John', password: 'secret' })
 * const safe = reactiveOmit(user, ['password'])
 * // safe = reactive({ id: 1, name: 'John' })
 * ```
 */
export function reactiveOmit<T extends object, K extends keyof T>(
  source: T,
  keys: K[]
): UnwrapNestedRefs<Omit<T, K>> {
  const omitted = { ...source };
  for (const key of keys) {
    delete omitted[key];
  }
  return reactive(omitted as Omit<T, K>) as UnwrapNestedRefs<Omit<T, K>>;
}

/**
 * 响应式合并
 *
 * @description 合并多个对象为响应式对象
 *
 * @example
 * ```ts
 * const merged = reactiveMerge(
 *   { a: 1 },
 *   { b: 2 },
 *   { c: 3 }
 * )
 * // merged = reactive({ a: 1, b: 2, c: 3 })
 * ```
 */
export function reactiveMerge<T extends object[]>(
  ...sources: T
): UnwrapNestedRefs<T[number]> {
  const merged = Object.assign({}, ...sources);
  return reactive(merged) as UnwrapNestedRefs<T[number]>;
}

/**
 * 响应式默认值
 *
 * @description 创建带默认值的响应式对象
 *
 * @example
 * ```ts
 * const state = reactiveDefault({ name: '', age: 0 })
 * state.name = null // 自动回退到 ''
 * ```
 */
export function reactiveDefault<T extends object>(
  defaults: T
): UnwrapNestedRefs<T> {
  const state = reactive(JSON.parse(JSON.stringify(defaults))) as UnwrapNestedRefs<T>;

  return new Proxy(state, {
    set(target, prop, value) {
      const key = prop as keyof T;
      if (value === null || value === undefined) {
        (target as T)[key] = defaults[key];
      } else {
        (target as T)[key] = value;
      }
      return true;
    },
  }) as UnwrapNestedRefs<T>;
}

/**
 * 响应式扩展
 *
 * @description 扩展响应式对象的属性
 *
 * @example
 * ```ts
 * const base = reactive({ name: 'John' })
 * const extended = reactiveExtend(base, { age: 30 })
 * // extended 包含 name 和 age
 * ```
 */
export function reactiveExtend<T extends object, E extends object>(
  source: T,
  extension: E
): UnwrapNestedRefs<T & E> {
  const extended = Object.assign({}, toRaw(source), extension);
  return reactive(extended) as UnwrapNestedRefs<T & E>;
}

/**
 * 同步响应式对象
 *
 * @description 保持两个响应式对象同步
 *
 * @example
 * ```ts
 * const source = reactive({ count: 0 })
 * const target = reactive({ count: 0 })
 *
 * const stop = syncReactive(source, target)
 * source.count = 1 // target.count 也变为 1
 * ```
 */
export function syncReactive<T extends object>(
  source: T,
  target: T,
  options: { immediate?: boolean; deep?: boolean } = {}
): () => void {
  const { immediate = true, deep = true } = options;

  return watch(
    () => ({ ...source }),
    (newValue) => {
      Object.assign(target, newValue);
    },
    { immediate, deep }
  );
}

/**
 * 响应式条件
 *
 * @description 根据条件创建响应式对象
 *
 * @example
 * ```ts
 * const isAdmin = ref(true)
 * const permissions = reactiveWhen(
 *   isAdmin,
 *   { read: true, write: true, delete: true },
 *   { read: true, write: false, delete: false }
 * )
 * ```
 */
export function reactiveWhen<T extends object, F extends object>(
  condition: Ref<boolean>,
  truthy: T,
  falsy: F
): UnwrapNestedRefs<T | F> {
  const state = reactive(condition.value ? truthy : falsy) as UnwrapNestedRefs<T | F>;

  watch(condition, (value) => {
    const newState = value ? truthy : falsy;
    // 清除旧属性
    for (const key in state) {
      delete (state as Record<string, unknown>)[key];
    }
    // 设置新属性
    Object.assign(state, newState);
  });

  return state;
}

/**
 * 响应式转换
 *
 * @description 对响应式对象的值进行转换
 *
 * @example
 * ```ts
 * const raw = reactive({ price: 100 })
 * const formatted = reactiveTransform(raw, {
 *   price: (v) => `$${v.toFixed(2)}`
 * })
 * console.log(formatted.price) // '$100.00'
 * ```
 */
export function reactiveTransform<T extends object, R extends object>(
  source: T,
  transforms: { [K in keyof T]?: (value: T[K]) => R[K] }
): UnwrapNestedRefs<R> {
  const result = {} as R;

  for (const key in source) {
    const transform = transforms[key];
    if (transform) {
      (result as Record<string, unknown>)[key] = computed(() => transform(source[key])).value;
    } else {
      (result as Record<string, unknown>)[key] = source[key];
    }
  }

  return reactive(result) as UnwrapNestedRefs<R>;
}

/**
 * 响应式验证
 *
 * @description 创建带验证的响应式对象
 *
 * @example
 * ```ts
 * const { data, errors, isValid, validate } = reactiveValidated(
 *   { email: '' },
 *   {
 *     email: (v) => v.includes('@') || '无效的邮箱'
 *   }
 * )
 *
 * data.email = 'invalid'
 * console.log(errors.value) // { email: '无效的邮箱' }
 * console.log(isValid.value) // false
 * ```
 */
export function reactiveValidated<T extends object>(
  initial: T,
  validators: { [K in keyof T]?: (value: T[K]) => true | string }
): {
  data: UnwrapNestedRefs<T>;
  errors: Ref<Partial<Record<keyof T, string>>>;
  isValid: Ref<boolean>;
  validate: () => boolean;
} {
  const data = reactive(JSON.parse(JSON.stringify(initial))) as UnwrapNestedRefs<T>;
  const errors = ref<Partial<Record<keyof T, string>>>({});

  const validate = (): boolean => {
    const newErrors: Partial<Record<keyof T, string>> = {};
    let valid = true;

    for (const key in validators) {
      const validator = validators[key];
      if (validator) {
        const result = validator((data as T)[key]);
        if (result !== true) {
          newErrors[key] = result;
          valid = false;
        }
      }
    }

    errors.value = newErrors;
    return valid;
  };

  // 自动验证
  watch(
    () => ({ ...data }),
    () => validate(),
    { deep: true }
  );

  const isValid = computed(() => Object.keys(errors.value).length === 0);

  return {
    data,
    errors,
    isValid,
    validate,
  };
}

/**
 * 响应式重置
 *
 * @description 创建可重置的响应式对象
 *
 * @example
 * ```ts
 * const { state, reset, setInitial } = reactiveResettable({ count: 0 })
 *
 * state.count = 10
 * reset() // state.count = 0
 *
 * state.count = 20
 * setInitial() // 设置新的初始值为 { count: 20 }
 * state.count = 30
 * reset() // state.count = 20
 * ```
 */
export function reactiveResettable<T extends object>(
  initial: T
): {
  state: UnwrapNestedRefs<T>;
  reset: () => void;
  setInitial: (value?: T) => void;
} {
  let initialSnapshot = JSON.parse(JSON.stringify(initial));
  const state = reactive(JSON.parse(JSON.stringify(initial))) as UnwrapNestedRefs<T>;

  const reset = () => {
    Object.assign(state, JSON.parse(JSON.stringify(initialSnapshot)));
  };

  const setInitial = (value?: T) => {
    initialSnapshot = JSON.parse(JSON.stringify(value ?? toRaw(state)));
  };

  return {
    state,
    reset,
    setInitial,
  };
}

/**
 * 响应式只读字段
 *
 * @description 创建部分字段只读的响应式对象
 *
 * @example
 * ```ts
 * const user = reactiveWithReadonly(
 *   { id: 1, name: 'John' },
 *   ['id'] // id 字段只读
 * )
 *
 * user.name = 'Jane' // 成功
 * user.id = 2 // 静默失败
 * ```
 */
export function reactiveWithReadonly<T extends object, K extends keyof T>(
  initial: T,
  readonlyKeys: K[]
): UnwrapNestedRefs<T> {
  const state = reactive(JSON.parse(JSON.stringify(initial))) as UnwrapNestedRefs<T>;
  const readonlySet = new Set(readonlyKeys);

  return new Proxy(state, {
    set(target, prop, value) {
      if (readonlySet.has(prop as K)) {
        return true; // 静默失败
      }
      (target as T)[prop as keyof T] = value;
      return true;
    },
  }) as UnwrapNestedRefs<T>;
}

/**
 * 响应式深度监听
 *
 * @description 监听响应式对象的所有属性变化
 *
 * @example
 * ```ts
 * const state = reactive({ user: { name: 'John' } })
 *
 * watchReactive(state, (path, newValue, oldValue) => {
 *   console.log(`${path} changed from ${oldValue} to ${newValue}`)
 * })
 *
 * state.user.name = 'Jane'
 * // 输出: 'user.name changed from John to Jane'
 * ```
 */
export function watchReactive<T extends object>(
  source: T,
  callback: (path: string, newValue: unknown, oldValue: unknown) => void,
  options: { deep?: boolean } = {}
): () => void {
  const { deep = true } = options;

  return watch(
    () => JSON.stringify(source),
    (newVal, oldVal) => {
      const newObj = JSON.parse(newVal);
      const oldObj = JSON.parse(oldVal);

      const findChanges = (obj1: Record<string, unknown>, obj2: Record<string, unknown>, prefix = "") => {
        for (const key in obj1) {
          const path = prefix ? `${prefix}.${key}` : key;
          const val1 = obj1[key];
          const val2 = obj2[key];

          if (typeof val1 === "object" && val1 !== null && typeof val2 === "object" && val2 !== null) {
            findChanges(val1 as Record<string, unknown>, val2 as Record<string, unknown>, path);
          } else if (val1 !== val2) {
            callback(path, val1, val2);
          }
        }
      };

      findChanges(newObj, oldObj);
    },
    { deep }
  );
}

/**
 * 响应式工具
 *
 * @description 提供响应式相关的工具函数
 */
export const reactiveUtils = {
  /** 获取原始对象 */
  toRaw,
  /** 检查是否响应式 */
  isReactive,
  /** 检查是否只读 */
  isReadonly,
  /** 检查是否代理 */
  isProxy,
  /** 标记为非响应式 */
  markRaw,
  /** 转换为 ref */
  toRef,
  /** 转换为 refs */
  toRefs,
};
