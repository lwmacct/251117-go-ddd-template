/**
 * EventSource Composable
 * 提供 Server-Sent Events (SSE) 的响应式管理
 */

import {
  ref,
  computed,
  watch,
  onMounted,
  onUnmounted,
  type Ref,
  type ComputedRef,
} from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export type EventSourceStatus = "CONNECTING" | "OPEN" | "CLOSED";

export interface UseEventSourceOptions {
  /** 是否自动连接 */
  immediate?: boolean;
  /** 是否自动重连 */
  autoReconnect?: boolean | {
    /** 最大重连次数 */
    retries?: number;
    /** 重连延迟（毫秒） */
    delay?: number;
    /** 延迟递增因子 */
    multiplier?: number;
    /** 最大延迟 */
    maxDelay?: number;
  };
  /** 是否携带凭证 */
  withCredentials?: boolean;
  /** 连接打开回调 */
  onOpen?: (event: Event) => void;
  /** 消息回调 */
  onMessage?: (event: MessageEvent) => void;
  /** 错误回调 */
  onError?: (event: Event) => void;
  /** 重连回调 */
  onReconnect?: (retries: number) => void;
  /** 连接失败回调 */
  onFailed?: () => void;
}

export interface UseEventSourceReturn<T = unknown> {
  /** EventSource 实例 */
  eventSource: Ref<EventSource | null>;
  /** 连接状态 */
  status: Ref<EventSourceStatus>;
  /** 是否已连接 */
  isConnected: ComputedRef<boolean>;
  /** 最后接收的数据 */
  data: Ref<T | null>;
  /** 最后接收的事件 */
  event: Ref<string | null>;
  /** 最后的事件 ID */
  lastEventId: Ref<string | null>;
  /** 错误信息 */
  error: Ref<Event | null>;
  /** 打开连接 */
  open: () => void;
  /** 关闭连接 */
  close: () => void;
  /** 重连次数 */
  retryCount: Ref<number>;
}

// ============================================================================
// useEventSource - SSE 连接管理
// ============================================================================

/**
 * Server-Sent Events 连接管理
 * @example
 * const { data, isConnected, open, close } = useEventSource('/api/events', {
 *   autoReconnect: true,
 *   onMessage: (event) => {
 *     console.log('收到事件:', event.data)
 *   }
 * })
 *
 * // 监听数据
 * watch(data, (newData) => {
 *   console.log('新数据:', newData)
 * })
 */
export function useEventSource<T = unknown>(
  url: string | Ref<string>,
  options: UseEventSourceOptions = {}
): UseEventSourceReturn<T> {
  const {
    immediate = true,
    autoReconnect = false,
    withCredentials = false,
    onOpen,
    onMessage,
    onError,
    onReconnect,
    onFailed,
  } = options;

  // 解析自动重连配置
  const autoReconnectOptions = typeof autoReconnect === "object"
    ? {
        retries: autoReconnect.retries ?? 3,
        delay: autoReconnect.delay ?? 1000,
        multiplier: autoReconnect.multiplier ?? 2,
        maxDelay: autoReconnect.maxDelay ?? 30000,
      }
    : autoReconnect
    ? { retries: 3, delay: 1000, multiplier: 2, maxDelay: 30000 }
    : null;

  // 响应式状态
  const eventSource = ref<EventSource | null>(null);
  const status = ref<EventSourceStatus>("CLOSED");
  const data = ref<T | null>(null) as Ref<T | null>;
  const event = ref<string | null>(null);
  const lastEventId = ref<string | null>(null);
  const error = ref<Event | null>(null);
  const retryCount = ref(0);

  const isConnected = computed(() => status.value === "OPEN");

  // 定时器
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null;

  // 是否手动关闭
  let manualClose = false;

  // 获取 URL
  const getUrl = () => {
    return typeof url === "object" && "value" in url ? url.value : url;
  };

  // 重连
  const reconnect = () => {
    if (!autoReconnectOptions) return;
    if (retryCount.value >= autoReconnectOptions.retries) {
      onFailed?.();
      return;
    }

    // 计算延迟
    const delay = Math.min(
      autoReconnectOptions.delay *
        Math.pow(autoReconnectOptions.multiplier, retryCount.value),
      autoReconnectOptions.maxDelay
    );

    reconnectTimer = setTimeout(() => {
      retryCount.value++;
      onReconnect?.(retryCount.value);
      open();
    }, delay);
  };

  // 打开连接
  const open = () => {
    if (typeof EventSource === "undefined") {
      console.warn("EventSource is not supported");
      return;
    }

    // 清理现有连接
    if (eventSource.value) {
      eventSource.value.close();
    }

    if (reconnectTimer) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }

    manualClose = false;
    status.value = "CONNECTING";

    const es = new EventSource(getUrl(), { withCredentials });

    es.onopen = (evt) => {
      status.value = "OPEN";
      error.value = null;
      retryCount.value = 0;
      onOpen?.(evt);
    };

    es.onmessage = (evt) => {
      event.value = evt.type;
      lastEventId.value = evt.lastEventId;

      // 尝试解析 JSON
      try {
        data.value = JSON.parse(evt.data);
      } catch {
        data.value = evt.data;
      }

      onMessage?.(evt);
    };

    es.onerror = (evt) => {
      error.value = evt;
      status.value = "CLOSED";
      onError?.(evt);

      // 自动重连
      if (!manualClose && autoReconnectOptions) {
        reconnect();
      }
    };

    eventSource.value = es;
  };

  // 关闭连接
  const close = () => {
    manualClose = true;

    if (reconnectTimer) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }

    if (eventSource.value) {
      eventSource.value.close();
      eventSource.value = null;
      status.value = "CLOSED";
    }
  };

  // 监听 URL 变化
  if (typeof url === "object" && "value" in url) {
    watch(url, () => {
      if (isConnected.value) {
        close();
        open();
      }
    });
  }

  // 自动连接
  if (immediate) {
    onMounted(open);
  }

  // 清理
  onUnmounted(() => {
    close();
  });

  return {
    eventSource,
    status,
    isConnected,
    data,
    event,
    lastEventId,
    error,
    open,
    close,
    retryCount,
  };
}

// ============================================================================
// useEventSourceNamed - 监听命名事件
// ============================================================================

export interface UseEventSourceNamedOptions extends UseEventSourceOptions {
  /** 要监听的事件名称列表 */
  events?: string[];
}

export interface NamedEventData<T = unknown> {
  /** 事件名称 */
  name: string;
  /** 事件数据 */
  data: T;
  /** 事件 ID */
  id: string | null;
}

export interface UseEventSourceNamedReturn<T = unknown>
  extends Omit<UseEventSourceReturn<T>, "data" | "event"> {
  /** 所有事件数据 */
  events: Ref<Map<string, T>>;
  /** 最后接收的命名事件 */
  lastEvent: Ref<NamedEventData<T> | null>;
  /** 添加事件监听 */
  addEventListener: (eventName: string) => void;
  /** 移除事件监听 */
  removeEventListener: (eventName: string) => void;
}

/**
 * 监听 SSE 命名事件
 * @example
 * const { events, lastEvent, addEventListener } = useEventSourceNamed('/api/events', {
 *   events: ['update', 'delete', 'create']
 * })
 *
 * // 获取特定事件的数据
 * const updateData = computed(() => events.value.get('update'))
 *
 * // 动态添加事件监听
 * addEventListener('custom-event')
 */
export function useEventSourceNamed<T = unknown>(
  url: string | Ref<string>,
  options: UseEventSourceNamedOptions = {}
): UseEventSourceNamedReturn<T> {
  const { events: eventNames = [], ...restOptions } = options;

  const events = ref(new Map<string, T>()) as Ref<Map<string, T>>;
  const lastEvent = ref<NamedEventData<T> | null>(null);

  // 事件处理函数映射
  const eventHandlers = new Map<string, (e: MessageEvent) => void>();

  const baseEventSource = useEventSource<T>(url, {
    ...restOptions,
    immediate: false,
  });

  // 创建事件处理函数
  const createHandler = (eventName: string) => {
    return (e: MessageEvent) => {
      let parsedData: T;
      try {
        parsedData = JSON.parse(e.data);
      } catch {
        parsedData = e.data;
      }

      events.value.set(eventName, parsedData);
      lastEvent.value = {
        name: eventName,
        data: parsedData,
        id: e.lastEventId,
      };
    };
  };

  // 添加事件监听
  const addEventListener = (eventName: string) => {
    if (!baseEventSource.eventSource.value) return;
    if (eventHandlers.has(eventName)) return;

    const handler = createHandler(eventName);
    eventHandlers.set(eventName, handler);
    baseEventSource.eventSource.value.addEventListener(eventName, handler);
  };

  // 移除事件监听
  const removeEventListener = (eventName: string) => {
    if (!baseEventSource.eventSource.value) return;

    const handler = eventHandlers.get(eventName);
    if (handler) {
      baseEventSource.eventSource.value.removeEventListener(eventName, handler);
      eventHandlers.delete(eventName);
      events.value.delete(eventName);
    }
  };

  // 重写 open 以添加事件监听
  const originalOpen = baseEventSource.open;
  baseEventSource.open = () => {
    originalOpen();

    // 等待连接打开后添加监听
    const checkConnection = setInterval(() => {
      if (baseEventSource.eventSource.value?.readyState === EventSource.OPEN) {
        clearInterval(checkConnection);
        eventNames.forEach(addEventListener);
      }
    }, 50);

    // 超时处理
    setTimeout(() => clearInterval(checkConnection), 5000);
  };

  // 自动连接
  if (restOptions.immediate !== false) {
    onMounted(() => {
      baseEventSource.open();
    });
  }

  return {
    ...baseEventSource,
    events,
    lastEvent,
    addEventListener,
    removeEventListener,
  };
}

// ============================================================================
// useServerSentEvents - 简化的 SSE Hook
// ============================================================================

/**
 * 简化的 SSE Hook
 * @example
 * const messages = useServerSentEvents<string[]>('/api/messages')
 *
 * // messages 是响应式的数据数组
 */
export function useServerSentEvents<T = unknown>(
  url: string | Ref<string>,
  options: UseEventSourceOptions = {}
): Ref<T | null> {
  const { data } = useEventSource<T>(url, options);
  return data;
}

// ============================================================================
// createEventSourceManager - SSE 管理器工厂
// ============================================================================

export interface EventSourceManager {
  /** 连接列表 */
  connections: Map<string, UseEventSourceReturn>;
  /** 创建连接 */
  create: (
    name: string,
    url: string,
    options?: UseEventSourceOptions
  ) => UseEventSourceReturn;
  /** 获取连接 */
  get: (name: string) => UseEventSourceReturn | undefined;
  /** 关闭连接 */
  close: (name: string) => void;
  /** 关闭所有连接 */
  closeAll: () => void;
}

/**
 * 创建 SSE 管理器
 * @example
 * const sseManager = createEventSourceManager()
 *
 * // 创建多个连接
 * sseManager.create('notifications', '/api/notifications')
 * sseManager.create('updates', '/api/updates')
 *
 * // 获取连接
 * const notifications = sseManager.get('notifications')
 *
 * // 关闭所有连接
 * sseManager.closeAll()
 */
export function createEventSourceManager(): EventSourceManager {
  const connections = new Map<string, UseEventSourceReturn>();

  const create = (
    name: string,
    url: string,
    options?: UseEventSourceOptions
  ): UseEventSourceReturn => {
    // 关闭现有连接
    if (connections.has(name)) {
      connections.get(name)?.close();
    }

    const es = useEventSource(url, options);
    connections.set(name, es);
    return es;
  };

  const get = (name: string): UseEventSourceReturn | undefined => {
    return connections.get(name);
  };

  const close = (name: string): void => {
    const es = connections.get(name);
    if (es) {
      es.close();
      connections.delete(name);
    }
  };

  const closeAll = (): void => {
    connections.forEach((es) => es.close());
    connections.clear();
  };

  return {
    connections,
    create,
    get,
    close,
    closeAll,
  };
}
