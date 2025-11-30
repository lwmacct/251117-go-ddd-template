/**
 * 事件总线 Composable
 * 提供组件间通信功能
 */

import { ref, onUnmounted, type Ref } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export type EventHandler<T = unknown> = (payload: T) => void;
export type UnsubscribeFn = () => void;

export interface EventBus<Events extends Record<string, unknown>> {
  /** 订阅事件 */
  on: <K extends keyof Events>(
    event: K,
    handler: EventHandler<Events[K]>
  ) => UnsubscribeFn;
  /** 取消订阅 */
  off: <K extends keyof Events>(
    event: K,
    handler?: EventHandler<Events[K]>
  ) => void;
  /** 触发事件 */
  emit: <K extends keyof Events>(event: K, payload: Events[K]) => void;
  /** 订阅一次 */
  once: <K extends keyof Events>(
    event: K,
    handler: EventHandler<Events[K]>
  ) => UnsubscribeFn;
  /** 清除所有事件 */
  clear: () => void;
  /** 获取事件监听器数量 */
  listenerCount: <K extends keyof Events>(event: K) => number;
}

// ============================================================================
// 全局事件总线
// ============================================================================

// 存储事件处理器
const globalListeners = new Map<string, Set<EventHandler>>();

/**
 * 创建事件总线
 * @example
 * // 定义事件类型
 * interface AppEvents {
 *   'user:login': { userId: string }
 *   'user:logout': void
 *   'notification': { message: string; type: 'info' | 'error' }
 * }
 *
 * const bus = createEventBus<AppEvents>()
 *
 * // 订阅
 * bus.on('user:login', ({ userId }) => {
 *   console.log('用户登录:', userId)
 * })
 *
 * // 触发
 * bus.emit('user:login', { userId: '123' })
 */
export function createEventBus<
  Events extends Record<string, unknown> = Record<string, unknown>
>(): EventBus<Events> {
  const listeners = new Map<keyof Events, Set<EventHandler>>();

  // 订阅事件
  const on = <K extends keyof Events>(
    event: K,
    handler: EventHandler<Events[K]>
  ): UnsubscribeFn => {
    if (!listeners.has(event)) {
      listeners.set(event, new Set());
    }
    listeners.get(event)!.add(handler as EventHandler);

    return () => off(event, handler);
  };

  // 取消订阅
  const off = <K extends keyof Events>(
    event: K,
    handler?: EventHandler<Events[K]>
  ) => {
    if (!listeners.has(event)) return;

    if (handler) {
      listeners.get(event)!.delete(handler as EventHandler);
    } else {
      listeners.delete(event);
    }
  };

  // 触发事件
  const emit = <K extends keyof Events>(event: K, payload: Events[K]) => {
    if (!listeners.has(event)) return;

    listeners.get(event)!.forEach((handler) => {
      handler(payload);
    });
  };

  // 订阅一次
  const once = <K extends keyof Events>(
    event: K,
    handler: EventHandler<Events[K]>
  ): UnsubscribeFn => {
    const wrappedHandler: EventHandler<Events[K]> = (payload) => {
      off(event, wrappedHandler);
      handler(payload);
    };
    return on(event, wrappedHandler);
  };

  // 清除所有
  const clear = () => {
    listeners.clear();
  };

  // 获取监听器数量
  const listenerCount = <K extends keyof Events>(event: K): number => {
    return listeners.get(event)?.size ?? 0;
  };

  return {
    on,
    off,
    emit,
    once,
    clear,
    listenerCount,
  };
}

// ============================================================================
// 全局事件总线单例
// ============================================================================

// 默认全局事件类型
export interface GlobalEvents {
  [key: string]: unknown;
}

// 全局事件总线实例
let globalEventBus: EventBus<GlobalEvents> | null = null;

/**
 * 获取全局事件总线
 */
export function getGlobalEventBus(): EventBus<GlobalEvents> {
  if (!globalEventBus) {
    globalEventBus = createEventBus<GlobalEvents>();
  }
  return globalEventBus;
}

// ============================================================================
// Vue Composable
// ============================================================================

/**
 * 使用事件总线（自动清理）
 * @example
 * // 在组件中使用
 * const { on, emit } = useEventBus<AppEvents>()
 *
 * // 订阅（组件卸载时自动取消）
 * on('user:login', (data) => {
 *   console.log('用户登录:', data)
 * })
 *
 * // 触发
 * emit('user:login', { userId: '123' })
 */
export function useEventBus<
  Events extends Record<string, unknown> = GlobalEvents
>(bus?: EventBus<Events>) {
  const eventBus = bus ?? (getGlobalEventBus() as unknown as EventBus<Events>);
  const unsubscribers: UnsubscribeFn[] = [];

  // 订阅（自动清理）
  const on = <K extends keyof Events>(
    event: K,
    handler: EventHandler<Events[K]>
  ) => {
    const unsubscribe = eventBus.on(event, handler);
    unsubscribers.push(unsubscribe);
    return unsubscribe;
  };

  // 订阅一次（自动清理）
  const once = <K extends keyof Events>(
    event: K,
    handler: EventHandler<Events[K]>
  ) => {
    const unsubscribe = eventBus.once(event, handler);
    unsubscribers.push(unsubscribe);
    return unsubscribe;
  };

  // 组件卸载时清理
  onUnmounted(() => {
    unsubscribers.forEach((unsubscribe) => unsubscribe());
    unsubscribers.length = 0;
  });

  return {
    on,
    once,
    off: eventBus.off,
    emit: eventBus.emit,
    clear: eventBus.clear,
    listenerCount: eventBus.listenerCount,
  };
}

// ============================================================================
// 特定事件的 Hook
// ============================================================================

/**
 * 监听特定事件
 * @example
 * useEventListener('user:login', (data) => {
 *   console.log('用户登录:', data)
 * })
 */
export function useEventListener<
  Events extends Record<string, unknown>,
  K extends keyof Events
>(
  event: K,
  handler: EventHandler<Events[K]>,
  bus?: EventBus<Events>
): UnsubscribeFn {
  const { on } = useEventBus(bus);
  return on(event, handler);
}

// ============================================================================
// 响应式事件值
// ============================================================================

/**
 * 响应式事件值
 * 订阅事件并将最新值存储在 ref 中
 * @example
 * const { value, reset } = useEventValue<AppEvents, 'notification'>('notification')
 *
 * // value 会自动更新为最新的 notification 事件 payload
 */
export function useEventValue<
  Events extends Record<string, unknown>,
  K extends keyof Events
>(
  event: K,
  initialValue?: Events[K],
  bus?: EventBus<Events>
): { value: Ref<Events[K] | undefined>; reset: () => void } {
  const value = ref(initialValue) as Ref<Events[K] | undefined>;
  const { on } = useEventBus(bus);

  on(event, (payload) => {
    value.value = payload;
  });

  const reset = () => {
    value.value = initialValue;
  };

  return { value, reset };
}

// ============================================================================
// 预定义的应用事件
// ============================================================================

/**
 * 应用事件类型示例
 */
export interface AppEvents {
  /** 用户登录 */
  "user:login": { userId: string; username: string };
  /** 用户登出 */
  "user:logout": void;
  /** 主题变化 */
  "theme:change": { theme: "light" | "dark" };
  /** 语言变化 */
  "locale:change": { locale: string };
  /** 通知 */
  "notification:show": {
    message: string;
    type: "success" | "info" | "warning" | "error";
    duration?: number;
  };
  /** 刷新数据 */
  "data:refresh": { scope?: string };
  /** 错误发生 */
  "error:occurred": { error: Error; context?: string };
}

/**
 * 应用事件总线
 */
export const appEventBus = createEventBus<AppEvents>();

/**
 * 使用应用事件总线
 */
export function useAppEventBus() {
  return useEventBus(appEventBus);
}
