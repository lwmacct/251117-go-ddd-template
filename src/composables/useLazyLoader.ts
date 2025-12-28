/**
 * 通用懒加载 Composable
 *
 * 核心功能：
 * - 按 key 追踪已加载/正在加载状态
 * - 防止并发重复请求（竞态防护）
 * - 支持自定义合并策略
 *
 * 使用场景：
 * - Tab 内容懒加载（按 tabKey）
 * - 分页列表（按 page number）
 * - 树形结构（按 parentId 加载子节点）
 *
 * @example
 * ```typescript
 * const { data, load, isLoaded, loading } = useLazyLoader({
 *   fetcher: (key) => api.fetchByCategory(key),
 *   merge: (existing, incoming) => [...existing, ...incoming],
 * });
 *
 * // 加载数据
 * await load('general');
 *
 * // 检查是否已加载
 * if (!isLoaded('notification')) {
 *   await load('notification');
 * }
 * ```
 */
import { shallowRef, type Ref, type ShallowRef } from "vue";

export interface UseLazyLoaderOptions<K extends string | number, T> {
  /**
   * 数据获取函数
   * @param key 加载的 key（如 category key、page number）
   * @returns 返回数据数组的 Promise
   */
  fetcher: (key: K) => Promise<T[]>;

  /**
   * 合并策略
   * @param existing 现有数据
   * @param incoming 新加载的数据
   * @param key 当前加载的 key
   * @returns 合并后的数据
   * @default 简单追加 [...existing, ...incoming]
   */
  merge?: (existing: T[], incoming: T[], key: K) => T[];

  /**
   * 从数据项提取 key（用于按 key 查找/更新）
   */
  getKey?: (item: T) => K;
}

export interface UseLazyLoaderReturn<K extends string | number, T> {
  /** 已加载的数据（使用 shallowRef 优化大数据集性能） */
  data: ShallowRef<T[]>;

  /** 已加载的 key 集合 */
  loaded: Ref<Set<K>>;

  /** 是否正在加载 */
  loading: Ref<boolean>;

  /** 最近一次加载的错误 */
  error: ShallowRef<Error | null>;

  /**
   * 加载指定 key 的数据
   * - 如果已加载或正在加载，直接跳过（防止重复请求）
   * - 加载成功后数据会合并到 data 中
   */
  load: (key: K) => Promise<void>;

  /** 检查指定 key 是否已加载 */
  isLoaded: (key: K) => boolean;

  /** 重置所有状态（清空数据和加载记录） */
  reset: () => void;
}

/**
 * 通用懒加载 Composable
 *
 * @param options 配置选项
 * @returns 懒加载状态和方法
 */
export function useLazyLoader<K extends string | number, T>(options: UseLazyLoaderOptions<K, T>): UseLazyLoaderReturn<K, T> {
  const { fetcher, merge } = options;

  // 使用 shallowRef 优化大数据集性能
  const data = shallowRef<T[]>([]) as ShallowRef<T[]>;
  const loaded = shallowRef<Set<K>>(new Set());
  const pending = shallowRef<Set<K>>(new Set());
  const loading = shallowRef(false);
  const error = shallowRef<Error | null>(null);

  /**
   * 加载指定 key 的数据
   */
  const load = async (key: K): Promise<void> => {
    // 如果已加载或正在加载，跳过（防止竞态条件）
    if (loaded.value.has(key) || pending.value.has(key)) {
      return;
    }

    // 标记为正在加载
    pending.value = new Set([...pending.value, key]);
    loading.value = true;
    error.value = null;

    try {
      const result = await fetcher(key);

      // 合并数据
      if (merge) {
        data.value = merge([...data.value], result, key);
      } else {
        // 默认：简单追加
        data.value = [...data.value, ...result];
      }

      // 标记为已加载
      loaded.value = new Set([...loaded.value, key]);
    } catch (e) {
      error.value = e as Error;
      console.error(`Failed to load data for key "${key}":`, e);
    } finally {
      // 清除正在加载标记
      const newPending = new Set(pending.value);
      newPending.delete(key);
      pending.value = newPending;

      // 只有当没有其他 pending 请求时才设置 loading = false
      loading.value = pending.value.size > 0;
    }
  };

  /**
   * 检查指定 key 是否已加载
   */
  const isLoaded = (key: K): boolean => {
    return loaded.value.has(key);
  };

  /**
   * 重置所有状态
   */
  const reset = (): void => {
    data.value = [];
    loaded.value = new Set();
    pending.value = new Set();
    loading.value = false;
    error.value = null;
  };

  return {
    data,
    loaded,
    loading,
    error,
    load,
    isLoaded,
    reset,
  };
}
