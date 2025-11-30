/**
 * Queue Composable
 * 提供队列相关的工具函数
 */

import { ref, computed, watch, onUnmounted, type Ref, type ComputedRef } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

/**
 * 队列项
 */
export interface QueueItem<T> {
  /** 唯一标识 */
  id: string | number;
  /** 数据 */
  data: T;
  /** 优先级（数字越大优先级越高） */
  priority?: number;
  /** 创建时间 */
  createdAt: number;
}

/**
 * 队列配置
 */
export interface QueueConfig {
  /** 最大容量 */
  maxSize?: number;
  /** 是否按优先级排序 */
  prioritized?: boolean;
}

/**
 * 队列返回值
 */
export interface UseQueueReturn<T> {
  /** 队列项 */
  items: Ref<QueueItem<T>[]>;
  /** 队列长度 */
  size: ComputedRef<number>;
  /** 是否为空 */
  isEmpty: ComputedRef<boolean>;
  /** 是否已满 */
  isFull: ComputedRef<boolean>;
  /** 入队 */
  enqueue: (data: T, priority?: number) => string | number;
  /** 出队 */
  dequeue: () => T | undefined;
  /** 查看队首 */
  peek: () => T | undefined;
  /** 清空队列 */
  clear: () => void;
  /** 移除指定项 */
  remove: (id: string | number) => boolean;
  /** 查找项 */
  find: (predicate: (item: T) => boolean) => T | undefined;
  /** 包含检查 */
  has: (id: string | number) => boolean;
}

/**
 * 任务队列配置
 */
export interface TaskQueueConfig {
  /** 并发数 */
  concurrency?: number;
  /** 任务间隔（毫秒） */
  interval?: number;
  /** 失败重试次数 */
  retries?: number;
  /** 重试延迟（毫秒） */
  retryDelay?: number;
  /** 自动开始 */
  autoStart?: boolean;
}

/**
 * 任务状态
 */
export type TaskStatus = "pending" | "running" | "completed" | "failed";

/**
 * 任务
 */
export interface Task<T = unknown> {
  /** 唯一标识 */
  id: string | number;
  /** 任务函数 */
  fn: () => Promise<T>;
  /** 状态 */
  status: TaskStatus;
  /** 结果 */
  result?: T;
  /** 错误 */
  error?: Error;
  /** 重试次数 */
  retryCount: number;
  /** 创建时间 */
  createdAt: number;
  /** 完成时间 */
  completedAt?: number;
}

/**
 * 任务队列返回值
 */
export interface UseTaskQueueReturn<T> {
  /** 所有任务 */
  tasks: Ref<Task<T>[]>;
  /** 待处理任务数 */
  pendingCount: ComputedRef<number>;
  /** 运行中任务数 */
  runningCount: ComputedRef<number>;
  /** 已完成任务数 */
  completedCount: ComputedRef<number>;
  /** 失败任务数 */
  failedCount: ComputedRef<number>;
  /** 是否正在运行 */
  isRunning: Ref<boolean>;
  /** 是否暂停 */
  isPaused: Ref<boolean>;
  /** 添加任务 */
  add: (fn: () => Promise<T>) => string | number;
  /** 开始处理 */
  start: () => void;
  /** 暂停处理 */
  pause: () => void;
  /** 恢复处理 */
  resume: () => void;
  /** 清空队列 */
  clear: () => void;
  /** 重试失败任务 */
  retryFailed: () => void;
  /** 移除任务 */
  remove: (id: string | number) => boolean;
  /** 等待所有任务完成 */
  waitAll: () => Promise<void>;
}

/**
 * 通知队列配置
 */
export interface NotificationQueueConfig {
  /** 最大显示数量 */
  maxVisible?: number;
  /** 默认持续时间（毫秒） */
  defaultDuration?: number;
  /** 新通知位置 */
  position?: "top" | "bottom";
}

/**
 * 通知项
 */
export interface Notification {
  /** 唯一标识 */
  id: string | number;
  /** 类型 */
  type: "info" | "success" | "warning" | "error";
  /** 标题 */
  title?: string;
  /** 消息 */
  message: string;
  /** 持续时间 */
  duration?: number;
  /** 是否可关闭 */
  closable?: boolean;
  /** 创建时间 */
  createdAt: number;
}

/**
 * 通知队列返回值
 */
export interface UseNotificationQueueReturn {
  /** 可见通知 */
  notifications: Ref<Notification[]>;
  /** 添加通知 */
  add: (notification: Omit<Notification, "id" | "createdAt">) => string | number;
  /** 移除通知 */
  remove: (id: string | number) => void;
  /** 清空所有 */
  clear: () => void;
  /** 快捷方法 */
  info: (message: string, title?: string) => string | number;
  success: (message: string, title?: string) => string | number;
  warning: (message: string, title?: string) => string | number;
  error: (message: string, title?: string) => string | number;
}

// ============================================================================
// 工具函数
// ============================================================================

let idCounter = 0;
const generateId = (): number => ++idCounter;

// ============================================================================
// 核心函数
// ============================================================================

/**
 * 使用队列
 *
 * @description 创建通用队列
 *
 * @example
 * ```ts
 * const { items, enqueue, dequeue, peek, size, isEmpty } = useQueue<string>({
 *   maxSize: 100
 * })
 *
 * enqueue('item1')
 * enqueue('item2')
 *
 * const first = dequeue() // 'item1'
 * const next = peek() // 'item2' (不移除)
 * ```
 */
export function useQueue<T>(config: QueueConfig = {}): UseQueueReturn<T> {
  const { maxSize = Infinity, prioritized = false } = config;

  const items = ref<QueueItem<T>[]>([]) as Ref<QueueItem<T>[]>;

  const size = computed(() => items.value.length);
  const isEmpty = computed(() => items.value.length === 0);
  const isFull = computed(() => items.value.length >= maxSize);

  const sortByPriority = () => {
    if (prioritized) {
      items.value.sort((a, b) => (b.priority ?? 0) - (a.priority ?? 0));
    }
  };

  const enqueue = (data: T, priority = 0): string | number => {
    if (isFull.value) {
      throw new Error("Queue is full");
    }

    const id = generateId();
    const item: QueueItem<T> = {
      id,
      data,
      priority,
      createdAt: Date.now(),
    };

    items.value = [...items.value, item];
    sortByPriority();

    return id;
  };

  const dequeue = (): T | undefined => {
    if (isEmpty.value) return undefined;

    const [first, ...rest] = items.value;
    items.value = rest;
    return first.data;
  };

  const peek = (): T | undefined => {
    return items.value[0]?.data;
  };

  const clear = () => {
    items.value = [];
  };

  const remove = (id: string | number): boolean => {
    const index = items.value.findIndex((item) => item.id === id);
    if (index === -1) return false;

    items.value = [...items.value.slice(0, index), ...items.value.slice(index + 1)];
    return true;
  };

  const find = (predicate: (item: T) => boolean): T | undefined => {
    const item = items.value.find((i) => predicate(i.data));
    return item?.data;
  };

  const has = (id: string | number): boolean => {
    return items.value.some((item) => item.id === id);
  };

  return {
    items,
    size,
    isEmpty,
    isFull,
    enqueue,
    dequeue,
    peek,
    clear,
    remove,
    find,
    has,
  };
}

/**
 * 使用栈
 *
 * @description 创建后进先出（LIFO）栈
 *
 * @example
 * ```ts
 * const { push, pop, peek, size, isEmpty } = useStack<number>()
 *
 * push(1)
 * push(2)
 * push(3)
 *
 * pop() // 3
 * pop() // 2
 * peek() // 1
 * ```
 */
export function useStack<T>(config: { maxSize?: number } = {}): {
  items: Ref<T[]>;
  size: ComputedRef<number>;
  isEmpty: ComputedRef<boolean>;
  isFull: ComputedRef<boolean>;
  push: (item: T) => void;
  pop: () => T | undefined;
  peek: () => T | undefined;
  clear: () => void;
} {
  const { maxSize = Infinity } = config;

  const items = ref<T[]>([]) as Ref<T[]>;

  const size = computed(() => items.value.length);
  const isEmpty = computed(() => items.value.length === 0);
  const isFull = computed(() => items.value.length >= maxSize);

  const push = (item: T) => {
    if (isFull.value) {
      throw new Error("Stack is full");
    }
    items.value = [...items.value, item];
  };

  const pop = (): T | undefined => {
    if (isEmpty.value) return undefined;
    const item = items.value[items.value.length - 1];
    items.value = items.value.slice(0, -1);
    return item;
  };

  const peek = (): T | undefined => {
    return items.value[items.value.length - 1];
  };

  const clear = () => {
    items.value = [];
  };

  return {
    items,
    size,
    isEmpty,
    isFull,
    push,
    pop,
    peek,
    clear,
  };
}

/**
 * 使用任务队列
 *
 * @description 创建可控的异步任务队列
 *
 * @example
 * ```ts
 * const {
 *   add,
 *   start,
 *   pause,
 *   tasks,
 *   isRunning,
 *   pendingCount,
 *   completedCount
 * } = useTaskQueue<string>({
 *   concurrency: 2,
 *   retries: 3
 * })
 *
 * // 添加任务
 * add(async () => {
 *   await fetch('/api/data')
 *   return 'done'
 * })
 *
 * // 开始处理
 * start()
 *
 * // 暂停
 * pause()
 * ```
 */
export function useTaskQueue<T = unknown>(config: TaskQueueConfig = {}): UseTaskQueueReturn<T> {
  const { concurrency = 1, interval = 0, retries = 0, retryDelay = 1000, autoStart = false } = config;

  const tasks = ref<Task<T>[]>([]) as Ref<Task<T>[]>;
  const isRunning = ref(false);
  const isPaused = ref(false);

  const pendingCount = computed(() => tasks.value.filter((t) => t.status === "pending").length);
  const runningCount = computed(() => tasks.value.filter((t) => t.status === "running").length);
  const completedCount = computed(() => tasks.value.filter((t) => t.status === "completed").length);
  const failedCount = computed(() => tasks.value.filter((t) => t.status === "failed").length);

  let processing = false;

  const processNext = async () => {
    if (isPaused.value || !isRunning.value || processing) return;

    const runningTasks = tasks.value.filter((t) => t.status === "running");
    if (runningTasks.length >= concurrency) return;

    const pendingTask = tasks.value.find((t) => t.status === "pending");
    if (!pendingTask) {
      if (runningTasks.length === 0) {
        isRunning.value = false;
      }
      return;
    }

    processing = true;
    pendingTask.status = "running";

    try {
      pendingTask.result = await pendingTask.fn();
      pendingTask.status = "completed";
      pendingTask.completedAt = Date.now();
    } catch (e) {
      if (pendingTask.retryCount < retries) {
        pendingTask.retryCount++;
        pendingTask.status = "pending";

        if (retryDelay > 0) {
          await new Promise((resolve) => setTimeout(resolve, retryDelay));
        }
      } else {
        pendingTask.status = "failed";
        pendingTask.error = e instanceof Error ? e : new Error(String(e));
        pendingTask.completedAt = Date.now();
      }
    }

    processing = false;

    if (interval > 0) {
      await new Promise((resolve) => setTimeout(resolve, interval));
    }

    // 继续处理下一个
    processNext();
  };

  const add = (fn: () => Promise<T>): string | number => {
    const id = generateId();
    const task: Task<T> = {
      id,
      fn,
      status: "pending",
      retryCount: 0,
      createdAt: Date.now(),
    };

    tasks.value = [...tasks.value, task];

    if (autoStart && !isRunning.value) {
      start();
    } else if (isRunning.value && !isPaused.value) {
      processNext();
    }

    return id;
  };

  const start = () => {
    if (isRunning.value) return;
    isRunning.value = true;
    isPaused.value = false;

    // 启动多个并发处理
    for (let i = 0; i < concurrency; i++) {
      processNext();
    }
  };

  const pause = () => {
    isPaused.value = true;
  };

  const resume = () => {
    isPaused.value = false;
    for (let i = 0; i < concurrency; i++) {
      processNext();
    }
  };

  const clear = () => {
    tasks.value = tasks.value.filter((t) => t.status === "running");
  };

  const retryFailed = () => {
    tasks.value = tasks.value.map((t) => {
      if (t.status === "failed") {
        return { ...t, status: "pending" as TaskStatus, retryCount: 0, error: undefined };
      }
      return t;
    });

    if (isRunning.value && !isPaused.value) {
      processNext();
    }
  };

  const remove = (id: string | number): boolean => {
    const index = tasks.value.findIndex((t) => t.id === id);
    if (index === -1) return false;

    const task = tasks.value[index];
    if (task.status === "running") return false;

    tasks.value = [...tasks.value.slice(0, index), ...tasks.value.slice(index + 1)];
    return true;
  };

  const waitAll = (): Promise<void> => {
    return new Promise((resolve) => {
      const checkComplete = () => {
        const hasActive = tasks.value.some((t) => t.status === "pending" || t.status === "running");
        if (!hasActive) {
          resolve();
        } else {
          setTimeout(checkComplete, 100);
        }
      };
      checkComplete();
    });
  };

  return {
    tasks,
    pendingCount,
    runningCount,
    completedCount,
    failedCount,
    isRunning,
    isPaused,
    add,
    start,
    pause,
    resume,
    clear,
    retryFailed,
    remove,
    waitAll,
  };
}

/**
 * 使用通知队列
 *
 * @description 创建通知消息队列
 *
 * @example
 * ```ts
 * const { notifications, add, remove, info, success, error } = useNotificationQueue({
 *   maxVisible: 5,
 *   defaultDuration: 3000
 * })
 *
 * info('操作成功')
 * error('发生错误')
 *
 * // 自定义通知
 * add({
 *   type: 'warning',
 *   title: '警告',
 *   message: '请注意...',
 *   duration: 5000
 * })
 * ```
 */
export function useNotificationQueue(config: NotificationQueueConfig = {}): UseNotificationQueueReturn {
  const { maxVisible = 5, defaultDuration = 3000, position = "top" } = config;

  const notifications = ref<Notification[]>([]);
  const timers = new Map<string | number, ReturnType<typeof setTimeout>>();

  const add = (notification: Omit<Notification, "id" | "createdAt">): string | number => {
    const id = generateId();
    const item: Notification = {
      ...notification,
      id,
      createdAt: Date.now(),
      duration: notification.duration ?? defaultDuration,
      closable: notification.closable ?? true,
    };

    // 根据位置添加
    if (position === "top") {
      notifications.value = [item, ...notifications.value];
    } else {
      notifications.value = [...notifications.value, item];
    }

    // 限制可见数量
    if (notifications.value.length > maxVisible) {
      const oldest = position === "top" ? notifications.value[notifications.value.length - 1] : notifications.value[0];
      remove(oldest.id);
    }

    // 设置自动移除
    if (item.duration && item.duration > 0) {
      const timer = setTimeout(() => {
        remove(id);
      }, item.duration);
      timers.set(id, timer);
    }

    return id;
  };

  const remove = (id: string | number) => {
    const timer = timers.get(id);
    if (timer) {
      clearTimeout(timer);
      timers.delete(id);
    }

    notifications.value = notifications.value.filter((n) => n.id !== id);
  };

  const clear = () => {
    timers.forEach((timer) => clearTimeout(timer));
    timers.clear();
    notifications.value = [];
  };

  const info = (message: string, title?: string): string | number => {
    return add({ type: "info", message, title });
  };

  const success = (message: string, title?: string): string | number => {
    return add({ type: "success", message, title });
  };

  const warning = (message: string, title?: string): string | number => {
    return add({ type: "warning", message, title });
  };

  const error = (message: string, title?: string): string | number => {
    return add({ type: "error", message, title });
  };

  onUnmounted(() => {
    clear();
  });

  return {
    notifications,
    add,
    remove,
    clear,
    info,
    success,
    warning,
    error,
  };
}

/**
 * 使用历史队列
 *
 * @description 创建固定大小的历史记录队列
 *
 * @example
 * ```ts
 * const { items, add, undo, redo, canUndo, canRedo, current } = useHistoryQueue<string>(10)
 *
 * add('state1')
 * add('state2')
 * add('state3')
 *
 * undo() // 回到 state2
 * undo() // 回到 state1
 * redo() // 前进到 state2
 * ```
 */
export function useHistoryQueue<T>(maxSize = 50): {
  items: Ref<T[]>;
  current: ComputedRef<T | undefined>;
  currentIndex: Ref<number>;
  add: (item: T) => void;
  undo: () => T | undefined;
  redo: () => T | undefined;
  canUndo: ComputedRef<boolean>;
  canRedo: ComputedRef<boolean>;
  clear: () => void;
  goto: (index: number) => T | undefined;
} {
  const items = ref<T[]>([]) as Ref<T[]>;
  const currentIndex = ref(-1);

  const current = computed(() => items.value[currentIndex.value]);
  const canUndo = computed(() => currentIndex.value > 0);
  const canRedo = computed(() => currentIndex.value < items.value.length - 1);

  const add = (item: T) => {
    // 移除当前位置之后的所有项
    items.value = items.value.slice(0, currentIndex.value + 1);

    // 添加新项
    items.value = [...items.value, item];
    currentIndex.value = items.value.length - 1;

    // 限制大小
    if (items.value.length > maxSize) {
      items.value = items.value.slice(-maxSize);
      currentIndex.value = items.value.length - 1;
    }
  };

  const undo = (): T | undefined => {
    if (!canUndo.value) return undefined;
    currentIndex.value--;
    return items.value[currentIndex.value];
  };

  const redo = (): T | undefined => {
    if (!canRedo.value) return undefined;
    currentIndex.value++;
    return items.value[currentIndex.value];
  };

  const clear = () => {
    items.value = [];
    currentIndex.value = -1;
  };

  const goto = (index: number): T | undefined => {
    if (index < 0 || index >= items.value.length) return undefined;
    currentIndex.value = index;
    return items.value[index];
  };

  return {
    items,
    current,
    currentIndex,
    add,
    undo,
    redo,
    canUndo,
    canRedo,
    clear,
    goto,
  };
}

/**
 * 使用环形缓冲区
 *
 * @description 创建固定大小的环形缓冲区
 *
 * @example
 * ```ts
 * const { items, push, shift, isFull, size } = useRingBuffer<number>(5)
 *
 * push(1, 2, 3, 4, 5) // 缓冲区已满
 * push(6) // 1 被移除，6 被添加
 *
 * shift() // 2
 * ```
 */
export function useRingBuffer<T>(capacity: number): {
  items: Ref<T[]>;
  size: ComputedRef<number>;
  isFull: ComputedRef<boolean>;
  isEmpty: ComputedRef<boolean>;
  push: (...items: T[]) => void;
  shift: () => T | undefined;
  peek: () => T | undefined;
  peekLast: () => T | undefined;
  clear: () => void;
  toArray: () => T[];
} {
  const items = ref<T[]>([]) as Ref<T[]>;

  const size = computed(() => items.value.length);
  const isFull = computed(() => items.value.length >= capacity);
  const isEmpty = computed(() => items.value.length === 0);

  const push = (...newItems: T[]) => {
    for (const item of newItems) {
      if (items.value.length >= capacity) {
        items.value = items.value.slice(1);
      }
      items.value = [...items.value, item];
    }
  };

  const shift = (): T | undefined => {
    if (isEmpty.value) return undefined;
    const first = items.value[0];
    items.value = items.value.slice(1);
    return first;
  };

  const peek = (): T | undefined => items.value[0];
  const peekLast = (): T | undefined => items.value[items.value.length - 1];

  const clear = () => {
    items.value = [];
  };

  const toArray = (): T[] => [...items.value];

  return {
    items,
    size,
    isFull,
    isEmpty,
    push,
    shift,
    peek,
    peekLast,
    clear,
    toArray,
  };
}
