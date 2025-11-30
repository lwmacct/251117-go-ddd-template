/**
 * Component Composable
 * 提供组件相关的工具函数
 */

import {
  ref,
  computed,
  watch,
  useSlots,
  useAttrs,
  getCurrentInstance,
  nextTick,
  type Ref,
  type ComputedRef,
  type Slots,
  type SetupContext,
  type ComponentInternalInstance,
} from "vue";

// ============================================================================
// 类型定义
// ============================================================================

/**
 * 插槽信息
 */
export interface SlotInfo {
  /** 插槽名称 */
  name: string;
  /** 是否存在 */
  exists: boolean;
  /** 是否为空 */
  isEmpty: boolean;
}

/**
 * 组件信息
 */
export interface ComponentInfo {
  /** 组件名称 */
  name: string | undefined;
  /** 组件 UID */
  uid: number;
  /** 是否已挂载 */
  isMounted: Ref<boolean>;
  /** 父组件 */
  parent: ComponentInternalInstance | null;
  /** 根组件 */
  root: ComponentInternalInstance;
}

/**
 * 组件引用返回值
 */
export interface ComponentRefReturn<T> {
  /** 组件引用 */
  ref: Ref<T | null>;
  /** 是否已加载 */
  isLoaded: ComputedRef<boolean>;
  /** 等待加载 */
  whenLoaded: () => Promise<T>;
}

/**
 * 暴露方法选项
 */
export interface ExposeOptions<T extends object> {
  /** 暴露的方法 */
  methods: T;
  /** 是否只读 */
  readonly?: boolean;
}

// ============================================================================
// 插槽工具
// ============================================================================

/**
 * 使用插槽信息
 *
 * @description 获取组件插槽的详细信息
 *
 * @example
 * ```ts
 * const { hasSlot, isEmpty, getSlotNames, slotInfo } = useSlotsInfo()
 *
 * if (hasSlot('header')) {
 *   // 渲染 header 插槽
 * }
 *
 * if (!isEmpty('content')) {
 *   // content 插槽有内容
 * }
 * ```
 */
export function useSlotsInfo(): {
  /** 检查插槽是否存在 */
  hasSlot: (name: string) => boolean;
  /** 检查插槽是否为空 */
  isEmpty: (name: string) => boolean;
  /** 获取所有插槽名称 */
  getSlotNames: () => string[];
  /** 插槽详细信息 */
  slotInfo: ComputedRef<SlotInfo[]>;
  /** 原始插槽对象 */
  slots: Slots;
} {
  const slots = useSlots();

  const hasSlot = (name: string): boolean => {
    return !!slots[name];
  };

  const isEmpty = (name: string): boolean => {
    const slot = slots[name];
    if (!slot) return true;

    const content = slot();
    if (!content || content.length === 0) return true;

    // 检查是否只有注释或空文本节点
    return content.every((vnode) => {
      if (vnode.type === Comment) return true;
      if (typeof vnode.children === "string" && !vnode.children.trim()) return true;
      return false;
    });
  };

  const getSlotNames = (): string[] => {
    return Object.keys(slots);
  };

  const slotInfo = computed<SlotInfo[]>(() => {
    return getSlotNames().map((name) => ({
      name,
      exists: hasSlot(name),
      isEmpty: isEmpty(name),
    }));
  });

  return {
    hasSlot,
    isEmpty,
    getSlotNames,
    slotInfo,
    slots,
  };
}

/**
 * 使用条件插槽
 *
 * @description 根据插槽是否存在返回不同的渲染内容
 *
 * @example
 * ```ts
 * const renderHeader = useConditionalSlot('header', () => h('div', 'Default Header'))
 * ```
 */
export function useConditionalSlot(name: string, fallback?: () => unknown): () => unknown {
  const slots = useSlots();

  return () => {
    const slot = slots[name];
    if (slot) {
      return slot();
    }
    return fallback?.();
  };
}

/**
 * 使用插槽传递
 *
 * @description 获取用于传递给子组件的插槽
 *
 * @example
 * ```ts
 * const passSlots = useSlotPass(['header', 'footer'])
 * // 在模板中: <ChildComponent v-bind="passSlots" />
 * ```
 */
export function useSlotPass(slotNames?: string[]): Record<string, unknown> {
  const slots = useSlots();

  const names = slotNames ?? Object.keys(slots);

  return names.reduce(
    (acc, name) => {
      if (slots[name]) {
        acc[name] = slots[name];
      }
      return acc;
    },
    {} as Record<string, unknown>
  );
}

// ============================================================================
// 属性工具
// ============================================================================

/**
 * 使用属性
 *
 * @description 获取组件属性的增强功能
 *
 * @example
 * ```ts
 * const { attrs, hasAttr, getAttr, attrsWithout } = useAttrsEnhanced()
 *
 * if (hasAttr('disabled')) {
 *   // ...
 * }
 *
 * const className = getAttr('class', 'default-class')
 *
 * // 排除某些属性
 * const restAttrs = attrsWithout(['class', 'style'])
 * ```
 */
export function useAttrsEnhanced(): {
  /** 原始属性 */
  attrs: ReturnType<typeof useAttrs>;
  /** 检查属性是否存在 */
  hasAttr: (name: string) => boolean;
  /** 获取属性值 */
  getAttr: <T>(name: string, defaultValue?: T) => T | undefined;
  /** 排除指定属性 */
  attrsWithout: (names: string[]) => Record<string, unknown>;
  /** 只包含指定属性 */
  attrsOnly: (names: string[]) => Record<string, unknown>;
  /** 属性名列表 */
  attrNames: ComputedRef<string[]>;
} {
  const attrs = useAttrs();

  const hasAttr = (name: string): boolean => {
    return name in attrs;
  };

  const getAttr = <T>(name: string, defaultValue?: T): T | undefined => {
    const value = attrs[name];
    return value !== undefined ? (value as T) : defaultValue;
  };

  const attrsWithout = (names: string[]): Record<string, unknown> => {
    const result: Record<string, unknown> = {};
    for (const key in attrs) {
      if (!names.includes(key)) {
        result[key] = attrs[key];
      }
    }
    return result;
  };

  const attrsOnly = (names: string[]): Record<string, unknown> => {
    const result: Record<string, unknown> = {};
    for (const name of names) {
      if (name in attrs) {
        result[name] = attrs[name];
      }
    }
    return result;
  };

  const attrNames = computed(() => Object.keys(attrs));

  return {
    attrs,
    hasAttr,
    getAttr,
    attrsWithout,
    attrsOnly,
    attrNames,
  };
}

/**
 * 使用类名合并
 *
 * @description 合并组件接收的 class 属性
 *
 * @example
 * ```ts
 * const mergedClass = useClassMerge('btn', 'btn-primary')
 * // 如果父组件传入 class="custom"，结果为 "btn btn-primary custom"
 * ```
 */
export function useClassMerge(...baseClasses: string[]): ComputedRef<string> {
  const attrs = useAttrs();

  return computed(() => {
    const attrClass = attrs.class as string | string[] | Record<string, boolean> | undefined;

    const classes = [...baseClasses];

    if (typeof attrClass === "string") {
      classes.push(attrClass);
    } else if (Array.isArray(attrClass)) {
      classes.push(...attrClass);
    } else if (attrClass && typeof attrClass === "object") {
      for (const [key, value] of Object.entries(attrClass)) {
        if (value) classes.push(key);
      }
    }

    return classes.filter(Boolean).join(" ");
  });
}

/**
 * 使用样式合并
 *
 * @description 合并组件接收的 style 属性
 *
 * @example
 * ```ts
 * const mergedStyle = useStyleMerge({ color: 'red' })
 * ```
 */
export function useStyleMerge(
  baseStyle: Record<string, string | number>
): ComputedRef<Record<string, string | number>> {
  const attrs = useAttrs();

  return computed(() => {
    const attrStyle = attrs.style as Record<string, string | number> | string | undefined;

    let styleObj: Record<string, string | number> = { ...baseStyle };

    if (typeof attrStyle === "string") {
      // 解析 style 字符串
      const pairs = attrStyle.split(";").filter(Boolean);
      for (const pair of pairs) {
        const [key, value] = pair.split(":").map((s) => s.trim());
        if (key && value) {
          styleObj[key] = value;
        }
      }
    } else if (attrStyle && typeof attrStyle === "object") {
      styleObj = { ...styleObj, ...attrStyle };
    }

    return styleObj;
  });
}

// ============================================================================
// 组件信息
// ============================================================================

/**
 * 使用组件信息
 *
 * @description 获取当前组件的详细信息
 *
 * @example
 * ```ts
 * const { name, uid, isMounted, parent, root } = useComponentInfo()
 *
 * console.log(`组件 ${name} (${uid}) 已挂载: ${isMounted.value}`)
 * ```
 */
export function useComponentInfo(): ComponentInfo {
  const instance = getCurrentInstance();
  const isMounted = ref(false);

  if (instance) {
    // 使用 nextTick 确保挂载状态正确
    nextTick(() => {
      isMounted.value = true;
    });
  }

  return {
    name: instance?.type.name ?? instance?.type.__name,
    uid: instance?.uid ?? 0,
    isMounted,
    parent: instance?.parent ?? null,
    root: instance?.root ?? (instance as ComponentInternalInstance),
  };
}

/**
 * 使用父组件
 *
 * @description 获取父组件实例
 *
 * @example
 * ```ts
 * const parent = useParentComponent()
 * if (parent) {
 *   console.log('Parent component:', parent.type.name)
 * }
 * ```
 */
export function useParentComponent(): ComponentInternalInstance | null {
  const instance = getCurrentInstance();
  return instance?.parent ?? null;
}

/**
 * 使用根组件
 *
 * @description 获取根组件实例
 *
 * @example
 * ```ts
 * const root = useRootComponent()
 * ```
 */
export function useRootComponent(): ComponentInternalInstance | null {
  const instance = getCurrentInstance();
  return instance?.root ?? null;
}

// ============================================================================
// 组件引用
// ============================================================================

/**
 * 使用组件引用
 *
 * @description 创建类型安全的组件引用
 *
 * @example
 * ```ts
 * const { ref: formRef, isLoaded, whenLoaded } = useComponentRef<FormInstance>()
 *
 * // 在模板中: <Form ref="formRef" />
 *
 * // 等待组件加载
 * const form = await whenLoaded()
 * form.validate()
 * ```
 */
export function useComponentRef<T>(): ComponentRefReturn<T> {
  const componentRef = ref<T | null>(null);

  const isLoaded = computed(() => componentRef.value !== null);

  const whenLoaded = (): Promise<T> => {
    return new Promise((resolve) => {
      if (componentRef.value) {
        resolve(componentRef.value);
        return;
      }

      const stop = watch(
        componentRef,
        (value) => {
          if (value) {
            stop();
            resolve(value);
          }
        },
        { immediate: true }
      );
    });
  };

  return {
    ref: componentRef as Ref<T | null>,
    isLoaded,
    whenLoaded,
  };
}

/**
 * 使用多组件引用
 *
 * @description 创建多个组件引用的集合
 *
 * @example
 * ```ts
 * const { refs, setRef, getRef, getAllRefs } = useComponentRefs<InputInstance>()
 *
 * // 在模板中: <Input v-for="item in items" :ref="el => setRef(item.id, el)" />
 *
 * // 获取特定引用
 * const input = getRef('input-1')
 * input?.focus()
 * ```
 */
export function useComponentRefs<T>(): {
  /** 引用映射 */
  refs: Ref<Map<string | number, T>>;
  /** 设置引用 */
  setRef: (key: string | number, el: T | null) => void;
  /** 获取引用 */
  getRef: (key: string | number) => T | undefined;
  /** 获取所有引用 */
  getAllRefs: () => T[];
  /** 清除引用 */
  clearRef: (key: string | number) => void;
  /** 清除所有引用 */
  clearAll: () => void;
} {
  const refs = ref(new Map<string | number, T>()) as Ref<Map<string | number, T>>;

  const setRef = (key: string | number, el: T | null) => {
    if (el) {
      refs.value.set(key, el);
    } else {
      refs.value.delete(key);
    }
  };

  const getRef = (key: string | number): T | undefined => {
    return refs.value.get(key);
  };

  const getAllRefs = (): T[] => {
    return Array.from(refs.value.values());
  };

  const clearRef = (key: string | number) => {
    refs.value.delete(key);
  };

  const clearAll = () => {
    refs.value.clear();
  };

  return {
    refs,
    setRef,
    getRef,
    getAllRefs,
    clearRef,
    clearAll,
  };
}

// ============================================================================
// 暴露工具
// ============================================================================

/**
 * 使用方法暴露
 *
 * @description 便捷地暴露组件方法
 *
 * @example
 * ```ts
 * const state = ref(0)
 *
 * useExposeMethod({
 *   increment: () => state.value++,
 *   reset: () => state.value = 0,
 *   getValue: () => state.value
 * })
 *
 * // 父组件可以通过 ref 调用这些方法
 * ```
 */
export function useExposeMethod<T extends Record<string, (...args: unknown[]) => unknown>>(methods: T): void {
  const instance = getCurrentInstance();
  if (instance) {
    // 直接赋值到 exposed
    Object.assign(instance.exposed ?? {}, methods);
    instance.exposed = instance.exposed ?? methods;
  }
}

/**
 * 使用状态暴露
 *
 * @description 暴露组件的状态和方法
 *
 * @example
 * ```ts
 * const count = ref(0)
 * const name = ref('test')
 *
 * useExposeState({
 *   count,
 *   name
 * }, {
 *   increment: () => count.value++,
 *   setName: (n: string) => name.value = n
 * })
 * ```
 */
export function useExposeState<
  S extends Record<string, Ref<unknown>>,
  M extends Record<string, (...args: unknown[]) => unknown>,
>(state: S, methods?: M): void {
  const instance = getCurrentInstance();
  if (instance) {
    const exposed = {
      ...state,
      ...(methods ?? {}),
    };
    Object.assign(instance.exposed ?? {}, exposed);
    instance.exposed = instance.exposed ?? exposed;
  }
}

// ============================================================================
// 组件工具
// ============================================================================

/**
 * 使用强制更新
 *
 * @description 强制组件重新渲染
 *
 * @example
 * ```ts
 * const { forceUpdate, updateKey } = useForceUpdate()
 *
 * // 在模板中使用 key 触发重渲染
 * // <Component :key="updateKey" />
 *
 * // 手动触发
 * forceUpdate()
 * ```
 */
export function useForceUpdate(): {
  /** 强制更新 */
  forceUpdate: () => void;
  /** 更新键 */
  updateKey: Ref<number>;
} {
  const updateKey = ref(0);

  const forceUpdate = () => {
    updateKey.value++;
  };

  return {
    forceUpdate,
    updateKey,
  };
}

/**
 * 使用组件事件
 *
 * @description 创建类型安全的事件发射器
 *
 * @example
 * ```ts
 * interface Events {
 *   'update:value': [value: string]
 *   'change': [newValue: string, oldValue: string]
 *   'submit': []
 * }
 *
 * const emit = useComponentEmit<Events>()
 *
 * emit('update:value', 'new value')
 * emit('change', 'new', 'old')
 * emit('submit')
 * ```
 */
export function useComponentEmit<T extends Record<string, unknown[]>>(): <K extends keyof T>(
  event: K,
  ...args: T[K]
) => void {
  const instance = getCurrentInstance();

  return <K extends keyof T>(event: K, ...args: T[K]) => {
    instance?.emit(event as string, ...args);
  };
}

/**
 * 使用组件代理
 *
 * @description 获取组件的代理对象
 *
 * @example
 * ```ts
 * const proxy = useComponentProxy()
 * console.log(proxy?.$el) // DOM 元素
 * console.log(proxy?.$props) // props
 * ```
 */
export function useComponentProxy(): ComponentInternalInstance["proxy"] {
  const instance = getCurrentInstance();
  return instance?.proxy ?? null;
}

/**
 * 使用组件类型
 *
 * @description 获取组件的类型信息
 *
 * @example
 * ```ts
 * const { name, props, emits } = useComponentType()
 * console.log('Component name:', name)
 * console.log('Props:', props)
 * ```
 */
export function useComponentType(): {
  name: string | undefined;
  props: Record<string, unknown> | undefined;
  emits: string[] | undefined;
} {
  const instance = getCurrentInstance();
  const type = instance?.type as Record<string, unknown> | undefined;

  return {
    name: (type?.name as string) ?? (type?.__name as string),
    props: type?.props as Record<string, unknown> | undefined,
    emits: type?.emits as string[] | undefined,
  };
}

/**
 * 使用自定义属性
 *
 * @description 向组件添加自定义属性
 *
 * @example
 * ```ts
 * useCustomProperties({
 *   version: '1.0.0',
 *   author: 'Me'
 * })
 * ```
 */
export function useCustomProperties(properties: Record<string, unknown>): void {
  const instance = getCurrentInstance();
  if (instance) {
    Object.assign(instance.appContext.config.globalProperties, properties);
  }
}
