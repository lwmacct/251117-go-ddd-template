/**
 * WebSocket Composable
 * 提供 WebSocket 连接的响应式管理
 */

import { ref, computed, watch, onMounted, onUnmounted, type Ref, type ComputedRef } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export type WebSocketStatus = "CONNECTING" | "OPEN" | "CLOSING" | "CLOSED";

export interface UseWebSocketOptions {
  /** 是否自动连接 */
  immediate?: boolean;
  /** 是否自动重连 */
  autoReconnect?:
    | boolean
    | {
        /** 最大重连次数 */
        retries?: number;
        /** 重连延迟（毫秒） */
        delay?: number;
        /** 延迟递增因子 */
        multiplier?: number;
        /** 最大延迟 */
        maxDelay?: number;
        /** 是否在页面可见时重连 */
        onVisibilityChange?: boolean;
      };
  /** 心跳配置 */
  heartbeat?:
    | boolean
    | {
        /** 心跳消息 */
        message?: string | ArrayBuffer | Blob;
        /** 心跳间隔（毫秒） */
        interval?: number;
        /** 心跳超时（毫秒） */
        timeout?: number;
      };
  /** 子协议 */
  protocols?: string | string[];
  /** 连接打开回调 */
  onOpen?: (ws: WebSocket, event: Event) => void;
  /** 消息回调 */
  onMessage?: (ws: WebSocket, event: MessageEvent) => void;
  /** 连接关闭回调 */
  onClose?: (ws: WebSocket, event: CloseEvent) => void;
  /** 错误回调 */
  onError?: (ws: WebSocket, event: Event) => void;
  /** 重连回调 */
  onReconnect?: (retries: number) => void;
  /** 连接失败回调（超过最大重连次数） */
  onFailed?: () => void;
}

export interface UseWebSocketReturn<T = unknown> {
  /** WebSocket 实例 */
  ws: Ref<WebSocket | null>;
  /** 连接状态 */
  status: Ref<WebSocketStatus>;
  /** 是否已连接 */
  isConnected: ComputedRef<boolean>;
  /** 最后接收的数据 */
  data: Ref<T | null>;
  /** 错误信息 */
  error: Ref<Event | null>;
  /** 打开连接 */
  open: () => void;
  /** 关闭连接 */
  close: (code?: number, reason?: string) => void;
  /** 发送消息 */
  send: (data: string | ArrayBuffer | Blob, useBuffer?: boolean) => boolean;
  /** 重连次数 */
  retryCount: Ref<number>;
}

// ============================================================================
// useWebSocket - WebSocket 连接管理
// ============================================================================

/**
 * WebSocket 连接管理
 * @example
 * const { data, send, isConnected, open, close } = useWebSocket('ws://localhost:8080', {
 *   autoReconnect: true,
 *   heartbeat: {
 *     message: 'ping',
 *     interval: 30000
 *   },
 *   onMessage: (ws, event) => {
 *     console.log('收到消息:', event.data)
 *   }
 * })
 *
 * // 发送消息
 * send('Hello Server')
 *
 * // 发送 JSON
 * send(JSON.stringify({ type: 'chat', message: 'Hi' }))
 */
export function useWebSocket<T = unknown>(
  url: string | Ref<string>,
  options: UseWebSocketOptions = {}
): UseWebSocketReturn<T> {
  const {
    immediate = true,
    autoReconnect = false,
    heartbeat = false,
    protocols,
    onOpen,
    onMessage,
    onClose,
    onError,
    onReconnect,
    onFailed,
  } = options;

  // 解析自动重连配置
  const autoReconnectOptions =
    typeof autoReconnect === "object"
      ? {
          retries: autoReconnect.retries ?? 3,
          delay: autoReconnect.delay ?? 1000,
          multiplier: autoReconnect.multiplier ?? 2,
          maxDelay: autoReconnect.maxDelay ?? 30000,
          onVisibilityChange: autoReconnect.onVisibilityChange ?? true,
        }
      : autoReconnect
        ? { retries: 3, delay: 1000, multiplier: 2, maxDelay: 30000, onVisibilityChange: true }
        : null;

  // 解析心跳配置
  const heartbeatOptions =
    typeof heartbeat === "object"
      ? {
          message: heartbeat.message ?? "ping",
          interval: heartbeat.interval ?? 30000,
          timeout: heartbeat.timeout ?? 10000,
        }
      : heartbeat
        ? { message: "ping", interval: 30000, timeout: 10000 }
        : null;

  // 响应式状态
  const ws = ref<WebSocket | null>(null);
  const status = ref<WebSocketStatus>("CLOSED");
  const data = ref<T | null>(null) as Ref<T | null>;
  const error = ref<Event | null>(null);
  const retryCount = ref(0);

  const isConnected = computed(() => status.value === "OPEN");

  // 消息缓冲区（连接未建立时缓存消息）
  const messageBuffer: Array<string | ArrayBuffer | Blob> = [];

  // 定时器
  let heartbeatTimer: ReturnType<typeof setInterval> | null = null;
  let heartbeatTimeoutTimer: ReturnType<typeof setTimeout> | null = null;
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null;

  // 是否手动关闭
  let manualClose = false;

  // 获取 URL
  const getUrl = () => {
    return typeof url === "object" && "value" in url ? url.value : url;
  };

  // 清除心跳定时器
  const clearHeartbeat = () => {
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer);
      heartbeatTimer = null;
    }
    if (heartbeatTimeoutTimer) {
      clearTimeout(heartbeatTimeoutTimer);
      heartbeatTimeoutTimer = null;
    }
  };

  // 启动心跳
  const startHeartbeat = () => {
    if (!heartbeatOptions || !ws.value) return;

    clearHeartbeat();

    heartbeatTimer = setInterval(() => {
      if (ws.value?.readyState === WebSocket.OPEN) {
        ws.value.send(
          typeof heartbeatOptions.message === "string" ? heartbeatOptions.message : heartbeatOptions.message
        );

        // 设置超时检测
        heartbeatTimeoutTimer = setTimeout(() => {
          // 心跳超时，关闭连接触发重连
          ws.value?.close();
        }, heartbeatOptions.timeout);
      }
    }, heartbeatOptions.interval);
  };

  // 重置心跳超时
  const resetHeartbeatTimeout = () => {
    if (heartbeatTimeoutTimer) {
      clearTimeout(heartbeatTimeoutTimer);
      heartbeatTimeoutTimer = null;
    }
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
      autoReconnectOptions.delay * Math.pow(autoReconnectOptions.multiplier, retryCount.value),
      autoReconnectOptions.maxDelay
    );

    reconnectTimer = setTimeout(() => {
      retryCount.value++;
      onReconnect?.(retryCount.value);
      open();
    }, delay);
  };

  // 刷新消息缓冲区
  const flushBuffer = () => {
    if (ws.value?.readyState === WebSocket.OPEN) {
      while (messageBuffer.length > 0) {
        const msg = messageBuffer.shift();
        if (msg) {
          ws.value.send(msg);
        }
      }
    }
  };

  // 打开连接
  const open = () => {
    if (typeof WebSocket === "undefined") {
      console.warn("WebSocket is not supported");
      return;
    }

    // 清理现有连接
    if (ws.value) {
      ws.value.close();
    }

    if (reconnectTimer) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }

    manualClose = false;
    status.value = "CONNECTING";

    const socket = protocols ? new WebSocket(getUrl(), protocols) : new WebSocket(getUrl());

    socket.onopen = (event) => {
      status.value = "OPEN";
      error.value = null;
      retryCount.value = 0;

      // 刷新缓冲区
      flushBuffer();

      // 启动心跳
      startHeartbeat();

      onOpen?.(socket, event);
    };

    socket.onmessage = (event) => {
      // 重置心跳超时
      resetHeartbeatTimeout();

      // 解析数据
      try {
        data.value = JSON.parse(event.data);
      } catch {
        data.value = event.data;
      }

      onMessage?.(socket, event);
    };

    socket.onclose = (event) => {
      status.value = "CLOSED";
      clearHeartbeat();

      onClose?.(socket, event);

      // 自动重连
      if (!manualClose && autoReconnectOptions) {
        reconnect();
      }
    };

    socket.onerror = (event) => {
      error.value = event;
      onError?.(socket, event);
    };

    ws.value = socket;
  };

  // 关闭连接
  const close = (code = 1000, reason?: string) => {
    manualClose = true;
    clearHeartbeat();

    if (reconnectTimer) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }

    if (ws.value) {
      status.value = "CLOSING";
      ws.value.close(code, reason);
      ws.value = null;
    }
  };

  // 发送消息
  const send = (message: string | ArrayBuffer | Blob, useBuffer = true): boolean => {
    if (ws.value?.readyState === WebSocket.OPEN) {
      ws.value.send(message);
      return true;
    }

    if (useBuffer) {
      messageBuffer.push(message);
    }

    return false;
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

  // 页面可见性变化时重连
  if (autoReconnectOptions?.onVisibilityChange) {
    const handleVisibilityChange = () => {
      if (document.visibilityState === "visible" && !isConnected.value && !manualClose) {
        open();
      }
    };

    onMounted(() => {
      document.addEventListener("visibilitychange", handleVisibilityChange);
    });

    onUnmounted(() => {
      document.removeEventListener("visibilitychange", handleVisibilityChange);
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
    ws,
    status,
    isConnected,
    data,
    error,
    open,
    close,
    send,
    retryCount,
  };
}

// ============================================================================
// useWebSocketJSON - JSON 格式的 WebSocket
// ============================================================================

export interface UseWebSocketJSONOptions<T> extends Omit<UseWebSocketOptions, "onMessage"> {
  /** 消息回调 */
  onMessage?: (ws: WebSocket, data: T) => void;
}

export interface UseWebSocketJSONReturn<T, S = T> extends Omit<UseWebSocketReturn<T>, "send"> {
  /** 发送 JSON 消息 */
  send: (data: S, useBuffer?: boolean) => boolean;
}

/**
 * JSON 格式的 WebSocket
 * @example
 * interface Message {
 *   type: string
 *   payload: unknown
 * }
 *
 * const { data, send } = useWebSocketJSON<Message>('ws://localhost:8080')
 *
 * // 发送 JSON 对象
 * send({ type: 'chat', payload: { message: 'Hello' } })
 */
export function useWebSocketJSON<T = unknown, S = T>(
  url: string | Ref<string>,
  options: UseWebSocketJSONOptions<T> = {}
): UseWebSocketJSONReturn<T, S> {
  const { onMessage, ...restOptions } = options;

  const baseWs = useWebSocket<T>(url, {
    ...restOptions,
    onMessage: (ws, event) => {
      try {
        const parsed = JSON.parse(event.data);
        onMessage?.(ws, parsed);
      } catch {
        onMessage?.(ws, event.data);
      }
    },
  });

  const send = (data: S, useBuffer = true): boolean => {
    return baseWs.send(JSON.stringify(data), useBuffer);
  };

  return {
    ...baseWs,
    send,
  };
}

// ============================================================================
// useWebSocketBinary - 二进制 WebSocket
// ============================================================================

export interface UseWebSocketBinaryOptions extends Omit<UseWebSocketOptions, "onMessage"> {
  /** 二进制类型 */
  binaryType?: "blob" | "arraybuffer";
  /** 消息回调 */
  onMessage?: (ws: WebSocket, data: ArrayBuffer | Blob) => void;
}

export interface UseWebSocketBinaryReturn extends Omit<UseWebSocketReturn<ArrayBuffer | Blob>, "data"> {
  /** 最后接收的二进制数据 */
  data: Ref<ArrayBuffer | Blob | null>;
}

/**
 * 二进制 WebSocket
 * @example
 * const { data, send } = useWebSocketBinary('ws://localhost:8080', {
 *   binaryType: 'arraybuffer'
 * })
 *
 * // 发送二进制数据
 * const buffer = new ArrayBuffer(8)
 * send(buffer)
 */
export function useWebSocketBinary(
  url: string | Ref<string>,
  options: UseWebSocketBinaryOptions = {}
): UseWebSocketBinaryReturn {
  const { binaryType = "arraybuffer", onMessage, ...restOptions } = options;

  const data = ref<ArrayBuffer | Blob | null>(null);

  const baseWs = useWebSocket<ArrayBuffer | Blob>(url, {
    ...restOptions,
    onOpen: (ws, event) => {
      ws.binaryType = binaryType;
      options.onOpen?.(ws, event);
    },
    onMessage: (ws, event) => {
      data.value = event.data;
      onMessage?.(ws, event.data);
    },
  });

  return {
    ...baseWs,
    data,
  };
}

// ============================================================================
// createWebSocketManager - WebSocket 管理器工厂
// ============================================================================

export interface WebSocketManager {
  /** 连接列表 */
  connections: Map<string, UseWebSocketReturn>;
  /** 创建连接 */
  create: (name: string, url: string, options?: UseWebSocketOptions) => UseWebSocketReturn;
  /** 获取连接 */
  get: (name: string) => UseWebSocketReturn | undefined;
  /** 关闭连接 */
  close: (name: string) => void;
  /** 关闭所有连接 */
  closeAll: () => void;
  /** 广播消息 */
  broadcast: (message: string | ArrayBuffer | Blob) => void;
}

/**
 * 创建 WebSocket 管理器
 * @example
 * const wsManager = createWebSocketManager()
 *
 * // 创建多个连接
 * wsManager.create('chat', 'ws://localhost:8080/chat')
 * wsManager.create('notifications', 'ws://localhost:8080/notifications')
 *
 * // 获取连接
 * const chat = wsManager.get('chat')
 * chat?.send('Hello')
 *
 * // 广播消息
 * wsManager.broadcast('ping')
 *
 * // 关闭所有连接
 * wsManager.closeAll()
 */
export function createWebSocketManager(): WebSocketManager {
  const connections = new Map<string, UseWebSocketReturn>();

  const create = (name: string, url: string, options?: UseWebSocketOptions): UseWebSocketReturn => {
    // 关闭现有连接
    if (connections.has(name)) {
      connections.get(name)?.close();
    }

    const ws = useWebSocket(url, options);
    connections.set(name, ws);
    return ws;
  };

  const get = (name: string): UseWebSocketReturn | undefined => {
    return connections.get(name);
  };

  const close = (name: string): void => {
    const ws = connections.get(name);
    if (ws) {
      ws.close();
      connections.delete(name);
    }
  };

  const closeAll = (): void => {
    connections.forEach((ws) => ws.close());
    connections.clear();
  };

  const broadcast = (message: string | ArrayBuffer | Blob): void => {
    connections.forEach((ws) => {
      if (ws.isConnected.value) {
        ws.send(message);
      }
    });
  };

  return {
    connections,
    create,
    get,
    close,
    closeAll,
    broadcast,
  };
}
