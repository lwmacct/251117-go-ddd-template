/**
 * Lifecycle Composable
 * 提供增强的生命周期钩子工具函数
 */

import {
  onMounted,
  onUnmounted,
  onBeforeMount,
  onBeforeUnmount,
  onUpdated,
  onBeforeUpdate,
  onActivated,
  onDeactivated,
  onErrorCaptured,
  onRenderTracked,
  onRenderTriggered,
  ref,
  watch,
  nextTick,
  getCurrentInstance,
  type Ref,
  type ComponentInternalInstance,
} from "vue";

// ============================================================================
// 类型定义
// ============================================================================

/**
 * 挂载状态
 */
export interface MountedState {
  /** 是否已挂载 */
  isMounted: Ref<boolean>;
  /** 是否正在挂载 */
  isMounting: Ref<boolean>;
  /** 是否已卸载 */
  isUnmounted: Ref<boolean>;
}

/**
 * 生命周期追踪选项
 */
export interface LifecycleTrackerOptions {
  /** 是否输出日志 */
  log?: boolean;
  /** 自定义日志前缀 */
  prefix?: string;
}

/**
 * 生命周期事件
 */
export type LifecycleEvent =
  | "beforeMount"
  | "mounted"
  | "beforeUpdate"
  | "updated"
  | "beforeUnmount"
  | "unmounted"
  | "activated"
  | "deactivated"
  | "errorCaptured";

/**
 * 生命周期追踪返回值
 */
export interface LifecycleTrackerReturn {
  /** 生命周期事件历史 */
  history: Ref<Array<{ event: LifecycleEvent; timestamp: number }>>;
  /** 当前生命周期阶段 */
  currentPhase: Ref<LifecycleEvent | null>;
  /** 组件是否活跃 */
  isActive: Ref<boolean>;
}

/**
 * 异步挂载选项
 */
export interface AsyncMountedOptions {
  /** 超时时间 */
  timeout?: number;
  /** 超时回调 */
  onTimeout?: () => void;
}

/**
 * 清理函数
 */
export type CleanupFn = () => void;

// ============================================================================
// 核心函数
// ============================================================================

/**
 * 使用挂载状态
 *
 * @description 追踪组件的挂载状态
 *
 * @example
 * ```ts
 * const { isMounted, isMounting, isUnmounted } = useMountedState()
 *
 * // 在异步操作中检查
 * async function fetchData() {
 *   const data = await api.getData()
 *   if (isMounted.value) {
 *     // 安全更新状态
 *   }
 * }
 * ```
 */
export function useMountedState(): MountedState {
  const isMounted = ref(false);
  const isMounting = ref(true);
  const isUnmounted = ref(false);

  onMounted(() => {
    isMounted.value = true;
    isMounting.value = false;
  });

  onUnmounted(() => {
    isMounted.value = false;
    isUnmounted.value = true;
  });

  return {
    isMounted,
    isMounting,
    isUnmounted,
  };
}

/**
 * 安全的 onMounted
 *
 * @description 返回一个 ref 表示是否已挂载，用于异步操作安全检查
 *
 * @example
 * ```ts
 * const isMounted = useSafeMounted()
 *
 * async function loadData() {
 *   const data = await fetchData()
 *   if (!isMounted.value) return // 组件已卸载，取消更新
 *   state.value = data
 * }
 * ```
 */
export function useSafeMounted(): Readonly<Ref<boolean>> {
  const isMounted = ref(false);

  onMounted(() => {
    isMounted.value = true;
  });

  onUnmounted(() => {
    isMounted.value = false;
  });

  return isMounted;
}

/**
 * 挂载时执行（带清理）
 *
 * @description 在挂载时执行回调，支持返回清理函数
 *
 * @example
 * ```ts
 * useMounted(() => {
 *   const handler = () => console.log('click')
 *   window.addEventListener('click', handler)
 *
 *   // 返回清理函数
 *   return () => {
 *     window.removeEventListener('click', handler)
 *   }
 * })
 * ```
 */
export function useMounted(callback: () => CleanupFn | void): void {
  let cleanup: CleanupFn | void;

  onMounted(() => {
    cleanup = callback();
  });

  onUnmounted(() => {
    cleanup?.();
  });
}

/**
 * 异步挂载
 *
 * @description 等待组件挂载完成
 *
 * @example
 * ```ts
 * await useAsyncMounted()
 * // 此时组件已挂载
 * console.log('mounted')
 * ```
 */
export function useAsyncMounted(options: AsyncMountedOptions = {}): Promise<void> {
  const { timeout, onTimeout } = options;

  return new Promise((resolve, reject) => {
    let timer: ReturnType<typeof setTimeout> | null = null;

    if (timeout) {
      timer = setTimeout(() => {
        onTimeout?.();
        reject(new Error("Mount timeout"));
      }, timeout);
    }

    onMounted(() => {
      if (timer) clearTimeout(timer);
      resolve();
    });
  });
}

/**
 * 卸载时执行
 *
 * @description 在组件卸载时执行回调
 *
 * @example
 * ```ts
 * useUnmounted(() => {
 *   console.log('Component unmounted')
 *   cleanup()
 * })
 * ```
 */
export function useUnmounted(callback: () => void): void {
  onUnmounted(callback);
}

/**
 * 更新时执行
 *
 * @description 在组件更新时执行回调
 *
 * @example
 * ```ts
 * useUpdated(() => {
 *   console.log('Component updated')
 * })
 * ```
 */
export function useUpdated(callback: () => void): void {
  onUpdated(callback);
}

/**
 * 激活时执行（keep-alive）
 *
 * @description 在组件激活时执行回调
 *
 * @example
 * ```ts
 * useActivated(() => {
 *   console.log('Component activated')
 *   refreshData()
 * })
 * ```
 */
export function useActivated(callback: () => void): void {
  onActivated(callback);
}

/**
 * 停用时执行（keep-alive）
 *
 * @description 在组件停用时执行回调
 *
 * @example
 * ```ts
 * useDeactivated(() => {
 *   console.log('Component deactivated')
 *   pauseTimer()
 * })
 * ```
 */
export function useDeactivated(callback: () => void): void {
  onDeactivated(callback);
}

/**
 * 生命周期追踪器
 *
 * @description 追踪组件的所有生命周期事件
 *
 * @example
 * ```ts
 * const { history, currentPhase, isActive } = useLifecycleTracker({
 *   log: true,
 *   prefix: 'MyComponent'
 * })
 *
 * // 查看历史
 * console.log(history.value)
 * // [{ event: 'beforeMount', timestamp: 1234567890 }, ...]
 * ```
 */
export function useLifecycleTracker(options: LifecycleTrackerOptions = {}): LifecycleTrackerReturn {
  const { log = false, prefix = "Component" } = options;

  const history = ref<Array<{ event: LifecycleEvent; timestamp: number }>>([]);
  const currentPhase = ref<LifecycleEvent | null>(null);
  const isActive = ref(false);

  const record = (event: LifecycleEvent) => {
    const entry = { event, timestamp: Date.now() };
    history.value = [...history.value, entry];
    currentPhase.value = event;

    if (log) {
      console.log(`[${prefix}] ${event}`, entry.timestamp);
    }
  };

  onBeforeMount(() => record("beforeMount"));
  onMounted(() => {
    record("mounted");
    isActive.value = true;
  });
  onBeforeUpdate(() => record("beforeUpdate"));
  onUpdated(() => record("updated"));
  onBeforeUnmount(() => record("beforeUnmount"));
  onUnmounted(() => {
    record("unmounted");
    isActive.value = false;
  });
  onActivated(() => {
    record("activated");
    isActive.value = true;
  });
  onDeactivated(() => {
    record("deactivated");
    isActive.value = false;
  });
  onErrorCaptured(() => {
    record("errorCaptured");
    return true;
  });

  return {
    history,
    currentPhase,
    isActive,
  };
}

/**
 * 使用清理函数
 *
 * @description 注册多个清理函数，在卸载时统一执行
 *
 * @example
 * ```ts
 * const { onCleanup, cleanup } = useCleanup()
 *
 * onCleanup(() => {
 *   clearInterval(timer1)
 * })
 *
 * onCleanup(() => {
 *   subscription.unsubscribe()
 * })
 *
 * // 也可以手动触发清理
 * cleanup()
 * ```
 */
export function useCleanup(): {
  onCleanup: (fn: CleanupFn) => void;
  cleanup: () => void;
} {
  const cleanupFns: CleanupFn[] = [];

  const onCleanup = (fn: CleanupFn) => {
    cleanupFns.push(fn);
  };

  const cleanup = () => {
    cleanupFns.forEach((fn) => fn());
    cleanupFns.length = 0;
  };

  onUnmounted(cleanup);

  return {
    onCleanup,
    cleanup,
  };
}

/**
 * 延迟挂载
 *
 * @description 延迟执行挂载回调
 *
 * @example
 * ```ts
 * useMountedDelay(() => {
 *   // 延迟 100ms 执行
 *   initializeChart()
 * }, 100)
 * ```
 */
export function useMountedDelay(callback: () => void, delay: number): void {
  let timer: ReturnType<typeof setTimeout> | null = null;

  onMounted(() => {
    timer = setTimeout(callback, delay);
  });

  onUnmounted(() => {
    if (timer) clearTimeout(timer);
  });
}

/**
 * 挂载后下一帧执行
 *
 * @description 在挂载后的下一个动画帧执行
 *
 * @example
 * ```ts
 * useMountedNextFrame(() => {
 *   // DOM 已完全渲染
 *   measureElement()
 * })
 * ```
 */
export function useMountedNextFrame(callback: () => void): void {
  let frameId: number | null = null;

  onMounted(() => {
    frameId = requestAnimationFrame(callback);
  });

  onUnmounted(() => {
    if (frameId !== null) cancelAnimationFrame(frameId);
  });
}

/**
 * 挂载后下一个 tick 执行
 *
 * @description 在挂载后等待 Vue 更新 DOM
 *
 * @example
 * ```ts
 * useMountedNextTick(() => {
 *   // Vue 已完成 DOM 更新
 *   scrollToBottom()
 * })
 * ```
 */
export function useMountedNextTick(callback: () => void): void {
  onMounted(() => {
    nextTick(callback);
  });
}

/**
 * 条件挂载
 *
 * @description 当条件为真时执行挂载回调
 *
 * @example
 * ```ts
 * const isReady = ref(false)
 *
 * useMountedWhen(isReady, () => {
 *   console.log('Ready and mounted!')
 * })
 *
 * // 稍后
 * isReady.value = true // 触发回调
 * ```
 */
export function useMountedWhen(condition: Ref<boolean>, callback: () => CleanupFn | void): void {
  const { isMounted } = useMountedState();
  let cleanup: CleanupFn | void;

  watch(
    [isMounted, condition],
    ([mounted, ready]) => {
      if (mounted && ready) {
        cleanup = callback();
      } else {
        cleanup?.();
        cleanup = undefined;
      }
    },
    { immediate: true }
  );

  onUnmounted(() => {
    cleanup?.();
  });
}

/**
 * 使用组件实例
 *
 * @description 获取当前组件实例
 *
 * @example
 * ```ts
 * const { instance, uid, proxy } = useInstance()
 *
 * console.log(uid) // 组件唯一 ID
 * ```
 */
export function useInstance(): {
  instance: ComponentInternalInstance | null;
  uid: number;
  proxy: ComponentInternalInstance["proxy"];
} {
  const instance = getCurrentInstance();

  return {
    instance,
    uid: instance?.uid ?? 0,
    proxy: instance?.proxy ?? null,
  };
}

/**
 * 渲染计数
 *
 * @description 追踪组件渲染次数
 *
 * @example
 * ```ts
 * const { count, reset } = useRenderCount()
 *
 * console.log(count.value) // 渲染次数
 * ```
 */
export function useRenderCount(): {
  count: Ref<number>;
  reset: () => void;
} {
  const count = ref(0);

  onUpdated(() => {
    count.value++;
  });

  onMounted(() => {
    count.value = 1;
  });

  const reset = () => {
    count.value = 0;
  };

  return {
    count,
    reset,
  };
}

/**
 * 错误捕获
 *
 * @description 捕获组件及其子组件的错误
 *
 * @example
 * ```ts
 * const { error, clearError } = useErrorCapture((err, instance, info) => {
 *   console.error('Error captured:', err)
 *   reportError(err)
 * })
 *
 * if (error.value) {
 *   // 显示错误 UI
 * }
 * ```
 */
export function useErrorCapture(
  handler?: (err: Error, instance: ComponentInternalInstance | null, info: string) => boolean | void
): {
  error: Ref<Error | null>;
  clearError: () => void;
} {
  const error = ref<Error | null>(null);

  onErrorCaptured((err, instance, info) => {
    error.value = err;

    if (handler) {
      return handler(err, instance, info);
    }

    return true; // 阻止错误向上传播
  });

  const clearError = () => {
    error.value = null;
  };

  return {
    error,
    clearError,
  };
}

/**
 * 渲染追踪（开发模式）
 *
 * @description 追踪组件的响应式依赖
 *
 * @example
 * ```ts
 * useRenderTracking({
 *   onTracked: (e) => console.log('Tracked:', e),
 *   onTriggered: (e) => console.log('Triggered:', e)
 * })
 * ```
 */
export function useRenderTracking(options: {
  onTracked?: (event: unknown) => void;
  onTriggered?: (event: unknown) => void;
}): void {
  const { onTracked, onTriggered } = options;

  if (onTracked) {
    onRenderTracked((e) => onTracked(e));
  }

  if (onTriggered) {
    onRenderTriggered((e) => onTriggered(e));
  }
}

/**
 * 一次性挂载
 *
 * @description 确保回调只在第一次挂载时执行
 *
 * @example
 * ```ts
 * useMountedOnce(() => {
 *   // 即使组件被 keep-alive 重新激活，也只执行一次
 *   initializeOnce()
 * })
 * ```
 */
export function useMountedOnce(callback: () => void): void {
  let called = false;

  onMounted(() => {
    if (!called) {
      called = true;
      callback();
    }
  });
}

/**
 * 挂载/激活时执行
 *
 * @description 在挂载或激活时执行，适用于 keep-alive 组件
 *
 * @example
 * ```ts
 * useMountedOrActivated(() => {
 *   // 每次显示时都执行
 *   refreshData()
 * })
 * ```
 */
export function useMountedOrActivated(callback: () => CleanupFn | void): void {
  let cleanup: CleanupFn | void;

  const execute = () => {
    cleanup?.();
    cleanup = callback();
  };

  onMounted(execute);
  onActivated(execute);

  onDeactivated(() => {
    cleanup?.();
    cleanup = undefined;
  });

  onUnmounted(() => {
    cleanup?.();
  });
}

/**
 * 组件显示状态
 *
 * @description 追踪组件是否可见（考虑 keep-alive）
 *
 * @example
 * ```ts
 * const isVisible = useComponentVisible()
 *
 * watch(isVisible, (visible) => {
 *   if (visible) {
 *     startAnimation()
 *   } else {
 *     stopAnimation()
 *   }
 * })
 * ```
 */
export function useComponentVisible(): Readonly<Ref<boolean>> {
  const isVisible = ref(false);

  onMounted(() => {
    isVisible.value = true;
  });

  onActivated(() => {
    isVisible.value = true;
  });

  onDeactivated(() => {
    isVisible.value = false;
  });

  onUnmounted(() => {
    isVisible.value = false;
  });

  return isVisible;
}

/**
 * 组件存活时间
 *
 * @description 追踪组件从挂载到现在的时间
 *
 * @example
 * ```ts
 * const { aliveTime, mountedAt } = useComponentAliveTime()
 *
 * console.log(`组件已存活 ${aliveTime.value}ms`)
 * ```
 */
export function useComponentAliveTime(): {
  aliveTime: Ref<number>;
  mountedAt: Ref<number | null>;
} {
  const mountedAt = ref<number | null>(null);
  const aliveTime = ref(0);
  let timer: ReturnType<typeof setInterval> | null = null;

  onMounted(() => {
    mountedAt.value = Date.now();
    timer = setInterval(() => {
      if (mountedAt.value) {
        aliveTime.value = Date.now() - mountedAt.value;
      }
    }, 1000);
  });

  onUnmounted(() => {
    if (timer) clearInterval(timer);
  });

  return {
    aliveTime,
    mountedAt,
  };
}

/**
 * 组合生命周期钩子
 *
 * @description 将多个生命周期回调组合在一起
 *
 * @example
 * ```ts
 * useLifecycle({
 *   onMounted: () => console.log('mounted'),
 *   onUnmounted: () => console.log('unmounted'),
 *   onUpdated: () => console.log('updated'),
 *   onActivated: () => console.log('activated'),
 *   onDeactivated: () => console.log('deactivated')
 * })
 * ```
 */
export function useLifecycle(hooks: {
  onBeforeMount?: () => void;
  onMounted?: () => void;
  onBeforeUpdate?: () => void;
  onUpdated?: () => void;
  onBeforeUnmount?: () => void;
  onUnmounted?: () => void;
  onActivated?: () => void;
  onDeactivated?: () => void;
}): void {
  if (hooks.onBeforeMount) onBeforeMount(hooks.onBeforeMount);
  if (hooks.onMounted) onMounted(hooks.onMounted);
  if (hooks.onBeforeUpdate) onBeforeUpdate(hooks.onBeforeUpdate);
  if (hooks.onUpdated) onUpdated(hooks.onUpdated);
  if (hooks.onBeforeUnmount) onBeforeUnmount(hooks.onBeforeUnmount);
  if (hooks.onUnmounted) onUnmounted(hooks.onUnmounted);
  if (hooks.onActivated) onActivated(hooks.onActivated);
  if (hooks.onDeactivated) onDeactivated(hooks.onDeactivated);
}
