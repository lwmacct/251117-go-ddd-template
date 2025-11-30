/**
 * BroadcastChannel Composable
 * 提供跨标签页通信的响应式管理
 */

import { ref, computed, onMounted, onUnmounted, type Ref, type ComputedRef } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UseBroadcastChannelOptions<T = unknown> {
  /** 消息回调 */
  onMessage?: (event: MessageEvent<T>) => void;
  /** 错误回调 */
  onError?: (event: MessageEvent) => void;
}

export interface UseBroadcastChannelReturn<T = unknown, P = T> {
  /** BroadcastChannel 实例 */
  channel: Ref<BroadcastChannel | null>;
  /** 是否支持 */
  isSupported: ComputedRef<boolean>;
  /** 是否已关闭 */
  isClosed: Ref<boolean>;
  /** 最后接收的数据 */
  data: Ref<T | null>;
  /** 发送消息 */
  post: (data: P) => void;
  /** 关闭频道 */
  close: () => void;
  /** 错误信息 */
  error: Ref<Event | null>;
}

// ============================================================================
// useBroadcastChannel - 广播频道通信
// ============================================================================

/**
 * 广播频道通信（跨标签页）
 * @example
 * // 在多个标签页中使用相同的频道名
 * const { data, post, isSupported } = useBroadcastChannel<Message>('my-channel')
 *
 * // 发送消息
 * post({ type: 'update', payload: { id: 1 } })
 *
 * // 监听消息
 * watch(data, (newData) => {
 *   console.log('收到消息:', newData)
 * })
 */
export function useBroadcastChannel<T = unknown, P = T>(
  name: string,
  options: UseBroadcastChannelOptions<T> = {}
): UseBroadcastChannelReturn<T, P> {
  const { onMessage, onError } = options;

  const isSupported = computed(() => typeof BroadcastChannel !== "undefined");

  const channel = ref<BroadcastChannel | null>(null);
  const isClosed = ref(false);
  const data = ref<T | null>(null) as Ref<T | null>;
  const error = ref<Event | null>(null);

  const post = (message: P) => {
    if (channel.value && !isClosed.value) {
      channel.value.postMessage(message);
    }
  };

  const close = () => {
    if (channel.value) {
      channel.value.close();
      channel.value = null;
      isClosed.value = true;
    }
  };

  onMounted(() => {
    if (!isSupported.value) return;

    const bc = new BroadcastChannel(name);

    bc.onmessage = (event: MessageEvent<T>) => {
      data.value = event.data;
      onMessage?.(event);
    };

    bc.onmessageerror = (event: MessageEvent) => {
      error.value = event;
      onError?.(event);
    };

    channel.value = bc;
  });

  onUnmounted(() => {
    close();
  });

  return {
    channel,
    isSupported,
    isClosed,
    data,
    post,
    close,
    error,
  };
}

// ============================================================================
// useBroadcastChannelJSON - JSON 格式的广播频道
// ============================================================================

export interface UseBroadcastChannelJSONOptions<T> extends Omit<UseBroadcastChannelOptions<T>, "onMessage"> {
  /** 消息回调 */
  onMessage?: (data: T) => void;
}

/**
 * JSON 格式的广播频道
 * @example
 * interface SyncMessage {
 *   type: 'sync' | 'update'
 *   data: unknown
 * }
 *
 * const { data, post } = useBroadcastChannelJSON<SyncMessage>('sync-channel')
 *
 * // 发送 JSON 消息
 * post({ type: 'sync', data: { foo: 'bar' } })
 */
export function useBroadcastChannelJSON<T = unknown>(
  name: string,
  options: UseBroadcastChannelJSONOptions<T> = {}
): UseBroadcastChannelReturn<T, T> {
  const { onMessage, ...restOptions } = options;

  return useBroadcastChannel<T, T>(name, {
    ...restOptions,
    onMessage: (event) => {
      onMessage?.(event.data);
    },
  });
}

// ============================================================================
// useTabSync - 标签页同步状态
// ============================================================================

export interface UseTabSyncOptions<T> {
  /** 初始状态 */
  initialState: T;
  /** 是否立即同步 */
  immediate?: boolean;
  /** 状态变化回调 */
  onSync?: (state: T, source: "local" | "remote") => void;
}

export interface UseTabSyncReturn<T> {
  /** 同步状态 */
  state: Ref<T>;
  /** 是否支持 */
  isSupported: ComputedRef<boolean>;
  /** 更新状态并广播 */
  setState: (newState: T) => void;
  /** 合并状态并广播 */
  mergeState: (partialState: Partial<T>) => void;
  /** 请求其他标签页同步状态 */
  requestSync: () => void;
}

interface TabSyncMessage<T> {
  type: "state" | "request" | "response";
  state?: T;
  timestamp: number;
}

/**
 * 标签页状态同步
 * @example
 * const { state, setState, mergeState } = useTabSync<UserState>('user-state', {
 *   initialState: { name: '', theme: 'light' }
 * })
 *
 * // 更新状态（会同步到其他标签页）
 * setState({ name: 'John', theme: 'dark' })
 *
 * // 部分更新
 * mergeState({ theme: 'dark' })
 */
export function useTabSync<T extends object>(channelName: string, options: UseTabSyncOptions<T>): UseTabSyncReturn<T> {
  const { initialState, immediate = true, onSync } = options;

  const isSupported = computed(() => typeof BroadcastChannel !== "undefined");

  const state = ref<T>(initialState) as Ref<T>;
  let channel: BroadcastChannel | null = null;
  let lastTimestamp = 0;

  const broadcast = (message: TabSyncMessage<T>) => {
    if (channel) {
      channel.postMessage(message);
    }
  };

  const setState = (newState: T) => {
    state.value = newState;
    lastTimestamp = Date.now();
    broadcast({
      type: "state",
      state: newState,
      timestamp: lastTimestamp,
    });
    onSync?.(newState, "local");
  };

  const mergeState = (partialState: Partial<T>) => {
    const newState = { ...state.value, ...partialState };
    setState(newState);
  };

  const requestSync = () => {
    broadcast({
      type: "request",
      timestamp: Date.now(),
    });
  };

  onMounted(() => {
    if (!isSupported.value) return;

    channel = new BroadcastChannel(channelName);

    channel.onmessage = (event: MessageEvent<TabSyncMessage<T>>) => {
      const message = event.data;

      switch (message.type) {
        case "state":
          // 只接受比当前更新的状态
          if (message.state && message.timestamp > lastTimestamp) {
            state.value = message.state;
            lastTimestamp = message.timestamp;
            onSync?.(message.state, "remote");
          }
          break;

        case "request":
          // 响应同步请求
          broadcast({
            type: "response",
            state: state.value,
            timestamp: lastTimestamp,
          });
          break;

        case "response":
          // 处理同步响应（接受最新的状态）
          if (message.state && message.timestamp > lastTimestamp) {
            state.value = message.state;
            lastTimestamp = message.timestamp;
            onSync?.(message.state, "remote");
          }
          break;
      }
    };

    // 立即请求同步
    if (immediate) {
      requestSync();
    }
  });

  onUnmounted(() => {
    if (channel) {
      channel.close();
      channel = null;
    }
  });

  return {
    state,
    isSupported,
    setState,
    mergeState,
    requestSync,
  };
}

// ============================================================================
// useTabLeader - 标签页 Leader 选举
// ============================================================================

export interface UseTabLeaderOptions {
  /** 心跳间隔（毫秒） */
  heartbeatInterval?: number;
  /** 心跳超时（毫秒） */
  heartbeatTimeout?: number;
  /** 成为 Leader 回调 */
  onBecomeLeader?: () => void;
  /** 失去 Leader 回调 */
  onLoseLeadership?: () => void;
}

export interface UseTabLeaderReturn {
  /** 是否是 Leader */
  isLeader: Ref<boolean>;
  /** 是否支持 */
  isSupported: ComputedRef<boolean>;
  /** 主动放弃 Leader */
  resign: () => void;
  /** 尝试成为 Leader */
  elect: () => void;
}

interface LeaderMessage {
  type: "heartbeat" | "elect" | "resign";
  id: string;
  timestamp: number;
}

/**
 * 标签页 Leader 选举
 * @example
 * const { isLeader, resign, elect } = useTabLeader('app-leader', {
 *   onBecomeLeader: () => {
 *     // 成为 Leader 后执行的操作（如启动后台任务）
 *     startBackgroundSync()
 *   },
 *   onLoseLeadership: () => {
 *     // 失去 Leader 后执行的操作
 *     stopBackgroundSync()
 *   }
 * })
 *
 * // 只在 Leader 标签页执行
 * if (isLeader.value) {
 *   performExpensiveOperation()
 * }
 */
export function useTabLeader(channelName: string, options: UseTabLeaderOptions = {}): UseTabLeaderReturn {
  const { heartbeatInterval = 1000, heartbeatTimeout = 3000, onBecomeLeader, onLoseLeadership } = options;

  const isSupported = computed(() => typeof BroadcastChannel !== "undefined");
  const isLeader = ref(false);

  // 生成唯一 ID
  const tabId =
    typeof crypto !== "undefined" ? crypto.randomUUID() : `tab-${Date.now()}-${Math.random().toString(36).slice(2)}`;

  let channel: BroadcastChannel | null = null;
  let heartbeatTimer: ReturnType<typeof setInterval> | null = null;
  let lastLeaderHeartbeat = 0;
  let leaderCheckTimer: ReturnType<typeof setInterval> | null = null;

  const broadcast = (message: LeaderMessage) => {
    if (channel) {
      channel.postMessage(message);
    }
  };

  const sendHeartbeat = () => {
    if (isLeader.value) {
      broadcast({
        type: "heartbeat",
        id: tabId,
        timestamp: Date.now(),
      });
    }
  };

  const becomeLeader = () => {
    if (!isLeader.value) {
      isLeader.value = true;
      onBecomeLeader?.();

      // 开始发送心跳
      heartbeatTimer = setInterval(sendHeartbeat, heartbeatInterval);
      sendHeartbeat();
    }
  };

  const loseLeadership = () => {
    if (isLeader.value) {
      isLeader.value = false;
      onLoseLeadership?.();

      if (heartbeatTimer) {
        clearInterval(heartbeatTimer);
        heartbeatTimer = null;
      }
    }
  };

  const resign = () => {
    if (isLeader.value) {
      broadcast({
        type: "resign",
        id: tabId,
        timestamp: Date.now(),
      });
      loseLeadership();
    }
  };

  const elect = () => {
    // 发送选举消息
    broadcast({
      type: "elect",
      id: tabId,
      timestamp: Date.now(),
    });

    // 等待一小段时间看是否有 Leader 响应
    setTimeout(() => {
      if (Date.now() - lastLeaderHeartbeat > heartbeatTimeout) {
        becomeLeader();
      }
    }, 100);
  };

  const checkLeader = () => {
    // 如果超时没有收到心跳，尝试成为 Leader
    if (Date.now() - lastLeaderHeartbeat > heartbeatTimeout) {
      if (!isLeader.value) {
        elect();
      }
    }
  };

  onMounted(() => {
    if (!isSupported.value) return;

    channel = new BroadcastChannel(channelName);

    channel.onmessage = (event: MessageEvent<LeaderMessage>) => {
      const message = event.data;

      switch (message.type) {
        case "heartbeat":
          lastLeaderHeartbeat = message.timestamp;
          // 如果收到其他标签页的心跳，且自己是 Leader，需要比较
          if (isLeader.value && message.id !== tabId) {
            // ID 较小的保持 Leader（简单的选举策略）
            if (message.id < tabId) {
              loseLeadership();
            }
          }
          break;

        case "elect":
          // 如果自己是 Leader，发送心跳确认
          if (isLeader.value) {
            sendHeartbeat();
          }
          break;

        case "resign":
          // Leader 放弃，尝试选举
          if (message.id !== tabId) {
            elect();
          }
          break;
      }
    };

    // 初始化时尝试选举
    lastLeaderHeartbeat = 0;
    elect();

    // 定期检查 Leader 状态
    leaderCheckTimer = setInterval(checkLeader, heartbeatInterval);
  });

  onUnmounted(() => {
    resign();

    if (heartbeatTimer) {
      clearInterval(heartbeatTimer);
    }

    if (leaderCheckTimer) {
      clearInterval(leaderCheckTimer);
    }

    if (channel) {
      channel.close();
      channel = null;
    }
  });

  return {
    isLeader,
    isSupported,
    resign,
    elect,
  };
}

// ============================================================================
// useTabMessenger - 标签页消息传递
// ============================================================================

export type MessageHandler<T = unknown> = (data: T, source: string) => void;

export interface UseTabMessengerReturn<T = unknown> {
  /** 是否支持 */
  isSupported: ComputedRef<boolean>;
  /** 当前标签页 ID */
  tabId: string;
  /** 发送消息给所有标签页 */
  broadcast: (type: string, data: T) => void;
  /** 注册消息处理函数 */
  on: (type: string, handler: MessageHandler<T>) => void;
  /** 移除消息处理函数 */
  off: (type: string, handler?: MessageHandler<T>) => void;
  /** 发送一次性消息并等待响应 */
  request: <R = unknown>(type: string, data: T, timeout?: number) => Promise<R>;
}

interface MessengerMessage<T = unknown> {
  type: string;
  data: T;
  source: string;
  requestId?: string;
}

/**
 * 标签页消息传递
 * @example
 * const messenger = useTabMessenger<unknown>('app-messenger')
 *
 * // 注册处理函数
 * messenger.on('userLogout', () => {
 *   // 在所有标签页处理登出
 *   router.push('/login')
 * })
 *
 * // 广播消息
 * messenger.broadcast('userLogout', { userId: 123 })
 *
 * // 请求-响应模式
 * const result = await messenger.request('getData', { key: 'user' })
 */
export function useTabMessenger<T = unknown>(channelName: string): UseTabMessengerReturn<T> {
  const isSupported = computed(() => typeof BroadcastChannel !== "undefined");

  const tabId =
    typeof crypto !== "undefined" ? crypto.randomUUID() : `tab-${Date.now()}-${Math.random().toString(36).slice(2)}`;

  const handlers = new Map<string, Set<MessageHandler<T>>>();
  const pendingRequests = new Map<
    string,
    {
      resolve: (value: unknown) => void;
      reject: (reason: unknown) => void;
      timer: ReturnType<typeof setTimeout>;
    }
  >();

  let channel: BroadcastChannel | null = null;

  const broadcast = (type: string, data: T) => {
    if (channel) {
      const message: MessengerMessage<T> = {
        type,
        data,
        source: tabId,
      };
      channel.postMessage(message);
    }
  };

  const on = (type: string, handler: MessageHandler<T>) => {
    if (!handlers.has(type)) {
      handlers.set(type, new Set());
    }
    handlers.get(type)!.add(handler);
  };

  const off = (type: string, handler?: MessageHandler<T>) => {
    if (handler) {
      handlers.get(type)?.delete(handler);
    } else {
      handlers.delete(type);
    }
  };

  const request = <R = unknown>(type: string, data: T, timeout = 5000): Promise<R> => {
    return new Promise((resolve, reject) => {
      const requestId = `${tabId}-${Date.now()}-${Math.random().toString(36).slice(2)}`;

      const timer = setTimeout(() => {
        pendingRequests.delete(requestId);
        reject(new Error(`Request timeout: ${type}`));
      }, timeout);

      pendingRequests.set(requestId, {
        resolve: resolve as (value: unknown) => void,
        reject,
        timer,
      });

      if (channel) {
        const message: MessengerMessage<T> = {
          type,
          data,
          source: tabId,
          requestId,
        };
        channel.postMessage(message);
      }
    });
  };

  onMounted(() => {
    if (!isSupported.value) return;

    channel = new BroadcastChannel(channelName);

    channel.onmessage = (event: MessageEvent<MessengerMessage<T>>) => {
      const message = event.data;

      // 忽略自己发送的消息
      if (message.source === tabId) return;

      // 处理响应
      if (message.type.endsWith(":response") && message.requestId) {
        const pending = pendingRequests.get(message.requestId);
        if (pending) {
          clearTimeout(pending.timer);
          pendingRequests.delete(message.requestId);
          pending.resolve(message.data);
        }
        return;
      }

      // 调用处理函数
      const typeHandlers = handlers.get(message.type);
      if (typeHandlers) {
        typeHandlers.forEach((handler) => {
          handler(message.data, message.source);
        });
      }
    };
  });

  onUnmounted(() => {
    // 清理所有 pending requests
    pendingRequests.forEach(({ reject, timer }) => {
      clearTimeout(timer);
      reject(new Error("Channel closed"));
    });
    pendingRequests.clear();

    if (channel) {
      channel.close();
      channel = null;
    }
  });

  return {
    isSupported,
    tabId,
    broadcast,
    on,
    off,
    request,
  };
}
