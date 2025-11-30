/**
 * Provide/Inject Composable
 * 提供增强的依赖注入工具函数
 */

import {
  provide,
  inject,
  ref,
  reactive,
  computed,
  readonly,
  watch,
  type InjectionKey,
  type Ref,
  type ComputedRef,
} from "vue";

// ============================================================================
// 类型定义
// ============================================================================

/**
 * 创建上下文返回值
 */
export interface CreateContextReturn<T> {
  /** 提供上下文 */
  provide: (value: T) => void;
  /** 注入上下文 */
  inject: (defaultValue?: T) => T;
  /** 注入键 */
  key: InjectionKey<T>;
}

/**
 * 创建状态上下文返回值
 */
export interface CreateStateContextReturn<T> {
  /** 提供状态上下文 */
  provide: (initialValue: T) => { state: Ref<T>; setState: (value: T) => void };
  /** 注入状态上下文 */
  inject: () => { state: Ref<T>; setState: (value: T) => void };
  /** 注入键 */
  key: InjectionKey<{ state: Ref<T>; setState: (value: T) => void }>;
}

/**
 * 创建响应式上下文返回值
 */
export interface CreateReactiveContextReturn<T extends object> {
  /** 提供响应式上下文 */
  provide: (initialValue: T) => T;
  /** 注入响应式上下文 */
  inject: () => T;
  /** 注入键 */
  key: InjectionKey<T>;
}

/**
 * 事件总线返回值
 */
export interface EventBusContext<T extends Record<string, unknown>> {
  /** 发布事件 */
  emit: <K extends keyof T>(event: K, payload: T[K]) => void;
  /** 订阅事件 */
  on: <K extends keyof T>(event: K, handler: (payload: T[K]) => void) => () => void;
  /** 取消订阅 */
  off: <K extends keyof T>(event: K, handler: (payload: T[K]) => void) => void;
}

/**
 * 主题上下文
 */
export interface ThemeContext {
  /** 当前主题 */
  theme: Ref<string>;
  /** 是否暗色模式 */
  isDark: ComputedRef<boolean>;
  /** 设置主题 */
  setTheme: (theme: string) => void;
  /** 切换暗色模式 */
  toggleDark: () => void;
}

/**
 * 国际化上下文
 */
export interface I18nContext {
  /** 当前语言 */
  locale: Ref<string>;
  /** 可用语言列表 */
  availableLocales: string[];
  /** 设置语言 */
  setLocale: (locale: string) => void;
  /** 翻译函数 */
  t: (key: string, params?: Record<string, string | number>) => string;
}

// ============================================================================
// 核心函数
// ============================================================================

/**
 * 创建类型安全的上下文
 *
 * @description 创建带有类型检查的 provide/inject 对
 *
 * @example
 * ```ts
 * // 定义上下文
 * const UserContext = createContext<User>('User')
 *
 * // 父组件提供
 * UserContext.provide({ id: 1, name: 'John' })
 *
 * // 子组件注入
 * const user = UserContext.inject()
 * ```
 */
export function createContext<T>(name: string, defaultValue?: T): CreateContextReturn<T> {
  const key = Symbol(name) as InjectionKey<T>;

  const provideContext = (value: T) => {
    provide(key, value);
  };

  const injectContext = (fallback?: T): T => {
    const value = inject(key, fallback ?? defaultValue);
    if (value === undefined) {
      throw new Error(`[${name}] 必须在提供者组件内使用`);
    }
    return value;
  };

  return {
    provide: provideContext,
    inject: injectContext,
    key,
  };
}

/**
 * 创建状态上下文
 *
 * @description 创建可变状态的上下文
 *
 * @example
 * ```ts
 * const CountContext = createStateContext<number>('Count')
 *
 * // 父组件
 * const { state, setState } = CountContext.provide(0)
 *
 * // 子组件
 * const { state, setState } = CountContext.inject()
 * setState(state.value + 1)
 * ```
 */
export function createStateContext<T>(name: string): CreateStateContextReturn<T> {
  const key = Symbol(name) as InjectionKey<{
    state: Ref<T>;
    setState: (value: T) => void;
  }>;

  const provideContext = (initialValue: T) => {
    const state = ref(initialValue) as Ref<T>;
    const setState = (value: T) => {
      state.value = value;
    };

    const context = { state, setState };
    provide(key, context);
    return context;
  };

  const injectContext = () => {
    const context = inject(key);
    if (!context) {
      throw new Error(`[${name}] 必须在提供者组件内使用`);
    }
    return context;
  };

  return {
    provide: provideContext,
    inject: injectContext,
    key,
  };
}

/**
 * 创建响应式上下文
 *
 * @description 创建响应式对象的上下文
 *
 * @example
 * ```ts
 * interface AppState {
 *   user: User | null
 *   settings: Settings
 * }
 *
 * const AppContext = createReactiveContext<AppState>('App')
 *
 * // 父组件
 * const state = AppContext.provide({ user: null, settings: {} })
 *
 * // 子组件
 * const state = AppContext.inject()
 * state.user = { id: 1, name: 'John' }
 * ```
 */
export function createReactiveContext<T extends object>(name: string): CreateReactiveContextReturn<T> {
  const key = Symbol(name) as InjectionKey<T>;

  const provideContext = (initialValue: T): T => {
    const state = reactive(initialValue) as T;
    provide(key, state);
    return state;
  };

  const injectContext = (): T => {
    const state = inject(key);
    if (!state) {
      throw new Error(`[${name}] 必须在提供者组件内使用`);
    }
    return state;
  };

  return {
    provide: provideContext,
    inject: injectContext,
    key,
  };
}

/**
 * 创建只读上下文
 *
 * @description 创建只读的上下文，子组件无法修改
 *
 * @example
 * ```ts
 * const ConfigContext = createReadonlyContext<Config>('Config')
 *
 * // 父组件
 * ConfigContext.provide({ apiUrl: '/api', timeout: 5000 })
 *
 * // 子组件
 * const config = ConfigContext.inject()
 * // config 是只读的
 * ```
 */
export function createReadonlyContext<T extends object>(name: string): CreateContextReturn<Readonly<T>> {
  const key = Symbol(name) as InjectionKey<Readonly<T>>;

  const provideContext = (value: T) => {
    provide(key, readonly(reactive(value)) as Readonly<T>);
  };

  const injectContext = (defaultValue?: T): Readonly<T> => {
    const value = inject(key, defaultValue ? readonly(reactive(defaultValue)) : undefined);
    if (value === undefined) {
      throw new Error(`[${name}] 必须在提供者组件内使用`);
    }
    return value;
  };

  return {
    provide: provideContext,
    inject: injectContext,
    key,
  };
}

/**
 * 创建事件总线上下文
 *
 * @description 创建组件间的事件通信上下文
 *
 * @example
 * ```ts
 * interface Events {
 *   'user:login': User
 *   'user:logout': void
 *   'notification': { message: string; type: string }
 * }
 *
 * const EventBus = createEventBusContext<Events>('EventBus')
 *
 * // 父组件
 * EventBus.provide()
 *
 * // 子组件 A - 发送
 * const bus = EventBus.inject()
 * bus.emit('user:login', { id: 1, name: 'John' })
 *
 * // 子组件 B - 监听
 * const bus = EventBus.inject()
 * bus.on('user:login', (user) => {
 *   console.log('User logged in:', user)
 * })
 * ```
 */
export function createEventBusContext<T extends Record<string, unknown>>(
  name: string
): {
  provide: () => EventBusContext<T>;
  inject: () => EventBusContext<T>;
  key: InjectionKey<EventBusContext<T>>;
} {
  const key = Symbol(name) as InjectionKey<EventBusContext<T>>;

  const provideContext = (): EventBusContext<T> => {
    const listeners = new Map<keyof T, Set<(payload: unknown) => void>>();

    const emit = <K extends keyof T>(event: K, payload: T[K]) => {
      const handlers = listeners.get(event);
      handlers?.forEach((handler) => handler(payload));
    };

    const on = <K extends keyof T>(event: K, handler: (payload: T[K]) => void): (() => void) => {
      if (!listeners.has(event)) {
        listeners.set(event, new Set());
      }
      listeners.get(event)!.add(handler as (payload: unknown) => void);

      return () => off(event, handler);
    };

    const off = <K extends keyof T>(event: K, handler: (payload: T[K]) => void) => {
      listeners.get(event)?.delete(handler as (payload: unknown) => void);
    };

    const context: EventBusContext<T> = { emit, on, off };
    provide(key, context);
    return context;
  };

  const injectContext = (): EventBusContext<T> => {
    const context = inject(key);
    if (!context) {
      throw new Error(`[${name}] 必须在提供者组件内使用`);
    }
    return context;
  };

  return {
    provide: provideContext,
    inject: injectContext,
    key,
  };
}

/**
 * 创建主题上下文
 *
 * @description 创建主题管理的上下文
 *
 * @example
 * ```ts
 * const ThemeProvider = createThemeContext('Theme')
 *
 * // 父组件
 * const theme = ThemeProvider.provide('light')
 *
 * // 子组件
 * const { theme, isDark, setTheme, toggleDark } = ThemeProvider.inject()
 * toggleDark()
 * ```
 */
export function createThemeContext(name: string): {
  provide: (initialTheme?: string) => ThemeContext;
  inject: () => ThemeContext;
  key: InjectionKey<ThemeContext>;
} {
  const key = Symbol(name) as InjectionKey<ThemeContext>;

  const provideContext = (initialTheme = "light"): ThemeContext => {
    const theme = ref(initialTheme);
    const isDark = computed(() => theme.value === "dark");

    const setTheme = (newTheme: string) => {
      theme.value = newTheme;
    };

    const toggleDark = () => {
      theme.value = isDark.value ? "light" : "dark";
    };

    const context: ThemeContext = { theme, isDark, setTheme, toggleDark };
    provide(key, context);
    return context;
  };

  const injectContext = (): ThemeContext => {
    const context = inject(key);
    if (!context) {
      throw new Error(`[${name}] 必须在提供者组件内使用`);
    }
    return context;
  };

  return {
    provide: provideContext,
    inject: injectContext,
    key,
  };
}

/**
 * 创建国际化上下文
 *
 * @description 创建简单的国际化上下文
 *
 * @example
 * ```ts
 * const messages = {
 *   en: { greeting: 'Hello, {name}!' },
 *   zh: { greeting: '你好, {name}!' }
 * }
 *
 * const I18n = createI18nContext('I18n', messages)
 *
 * // 父组件
 * const i18n = I18n.provide('zh')
 *
 * // 子组件
 * const { t, locale, setLocale } = I18n.inject()
 * console.log(t('greeting', { name: 'World' })) // 你好, World!
 * ```
 */
export function createI18nContext(
  name: string,
  messages: Record<string, Record<string, string>>
): {
  provide: (initialLocale?: string) => I18nContext;
  inject: () => I18nContext;
  key: InjectionKey<I18nContext>;
} {
  const key = Symbol(name) as InjectionKey<I18nContext>;
  const availableLocales = Object.keys(messages);

  const provideContext = (initialLocale = availableLocales[0]): I18nContext => {
    const locale = ref(initialLocale);

    const setLocale = (newLocale: string) => {
      if (availableLocales.includes(newLocale)) {
        locale.value = newLocale;
      }
    };

    const t = (key: string, params?: Record<string, string | number>): string => {
      const message = messages[locale.value]?.[key] ?? key;

      if (!params) return message;

      return message.replace(/\{(\w+)\}/g, (_, param) => String(params[param] ?? `{${param}}`));
    };

    const context: I18nContext = { locale, availableLocales, setLocale, t };
    provide(key, context);
    return context;
  };

  const injectContext = (): I18nContext => {
    const context = inject(key);
    if (!context) {
      throw new Error(`[${name}] 必须在提供者组件内使用`);
    }
    return context;
  };

  return {
    provide: provideContext,
    inject: injectContext,
    key,
  };
}

/**
 * 使用可选注入
 *
 * @description 注入可能不存在的值，返回 undefined 而不是报错
 *
 * @example
 * ```ts
 * const user = useOptionalInject(UserKey)
 * if (user) {
 *   console.log(user.name)
 * }
 * ```
 */
export function useOptionalInject<T>(key: InjectionKey<T> | string, defaultValue?: T): T | undefined {
  return inject(key, defaultValue);
}

/**
 * 使用必需注入
 *
 * @description 注入必需的值，不存在则报错
 *
 * @example
 * ```ts
 * const config = useRequiredInject(ConfigKey, 'Config')
 * // 如果 Config 未提供，会抛出错误
 * ```
 */
export function useRequiredInject<T>(key: InjectionKey<T> | string, name: string): T {
  const value = inject(key);
  if (value === undefined) {
    throw new Error(`[${name}] 必须在提供者组件内使用`);
  }
  return value;
}

/**
 * 创建存储上下文
 *
 * @description 创建带有持久化的状态上下文
 *
 * @example
 * ```ts
 * const SettingsContext = createStorageContext<Settings>('Settings', {
 *   storageKey: 'app-settings',
 *   storage: localStorage
 * })
 *
 * // 父组件
 * const settings = SettingsContext.provide({ theme: 'dark' })
 *
 * // 子组件
 * const settings = SettingsContext.inject()
 * settings.theme = 'light' // 自动保存到 localStorage
 * ```
 */
export function createStorageContext<T extends object>(
  name: string,
  options: {
    storageKey: string;
    storage?: Storage;
    defaultValue?: T;
  }
): CreateReactiveContextReturn<T> {
  const { storageKey, storage = localStorage, defaultValue } = options;
  const key = Symbol(name) as InjectionKey<T>;

  const loadFromStorage = (): T | null => {
    try {
      const stored = storage.getItem(storageKey);
      return stored ? JSON.parse(stored) : null;
    } catch {
      return null;
    }
  };

  const saveToStorage = (value: T) => {
    try {
      storage.setItem(storageKey, JSON.stringify(value));
    } catch {
      // 忽略存储错误
    }
  };

  const provideContext = (initialValue: T): T => {
    const stored = loadFromStorage();
    const state = reactive(stored ?? initialValue) as T;

    // 监听变化并保存
    watch(
      () => ({ ...state }),
      (newValue) => saveToStorage(newValue as T),
      { deep: true }
    );

    provide(key, state);
    return state;
  };

  const injectContext = (): T => {
    const state = inject(key);
    if (!state) {
      if (defaultValue) {
        return reactive(defaultValue) as T;
      }
      throw new Error(`[${name}] 必须在提供者组件内使用`);
    }
    return state;
  };

  return {
    provide: provideContext,
    inject: injectContext,
    key,
  };
}

/**
 * 创建工厂上下文
 *
 * @description 创建工厂函数的上下文，每次注入都创建新实例
 *
 * @example
 * ```ts
 * const LoggerContext = createFactoryContext('Logger', (name: string) => ({
 *   log: (msg: string) => console.log(`[${name}] ${msg}`),
 *   error: (msg: string) => console.error(`[${name}] ${msg}`)
 * }))
 *
 * // 父组件
 * LoggerContext.provide()
 *
 * // 子组件
 * const createLogger = LoggerContext.inject()
 * const logger = createLogger('MyComponent')
 * logger.log('Hello')
 * ```
 */
export function createFactoryContext<T, Args extends unknown[]>(
  name: string,
  factory: (...args: Args) => T
): {
  provide: () => void;
  inject: () => (...args: Args) => T;
  key: InjectionKey<(...args: Args) => T>;
} {
  const key = Symbol(name) as InjectionKey<(...args: Args) => T>;

  const provideContext = () => {
    provide(key, factory);
  };

  const injectContext = (): ((...args: Args) => T) => {
    const factoryFn = inject(key);
    if (!factoryFn) {
      throw new Error(`[${name}] 必须在提供者组件内使用`);
    }
    return factoryFn;
  };

  return {
    provide: provideContext,
    inject: injectContext,
    key,
  };
}
