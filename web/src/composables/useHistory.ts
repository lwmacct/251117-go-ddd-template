/**
 * 历史记录工具 Composable
 * 提供撤销/重做功能
 */

import { ref, computed, watch, type Ref, type ComputedRef } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UseHistoryOptions<T> {
  /** 最大历史记录数量，默认无限制 */
  capacity?: number;
  /** 是否深拷贝，默认 true */
  deep?: boolean;
  /** 是否立即记录初始状态，默认 true */
  immediate?: boolean;
  /** 自定义克隆函数 */
  clone?: (value: T) => T;
  /** 自定义比较函数 */
  equals?: (a: T, b: T) => boolean;
}

export interface UseHistoryReturn<T> {
  /** 当前值 */
  current: Ref<T>;
  /** 是否可以撤销 */
  canUndo: ComputedRef<boolean>;
  /** 是否可以重做 */
  canRedo: ComputedRef<boolean>;
  /** 历史记录数量 */
  historyCount: ComputedRef<number>;
  /** 撤销 */
  undo: () => void;
  /** 重做 */
  redo: () => void;
  /** 清除历史 */
  clear: () => void;
  /** 手动提交 */
  commit: () => void;
  /** 重置到初始状态 */
  reset: () => void;
  /** 跳转到指定历史 */
  go: (delta: number) => void;
  /** 获取历史记录 */
  history: ComputedRef<T[]>;
  /** 当前索引 */
  index: Ref<number>;
}

// ============================================================================
// 工具函数
// ============================================================================

// 默认深拷贝
function defaultClone<T>(value: T): T {
  return JSON.parse(JSON.stringify(value));
}

// 默认比较
function defaultEquals<T>(a: T, b: T): boolean {
  return JSON.stringify(a) === JSON.stringify(b);
}

// ============================================================================
// 主函数
// ============================================================================

/**
 * 历史记录（撤销/重做）
 * @example
 * const state = ref({ count: 0 })
 * const { current, undo, redo, canUndo, canRedo } = useHistory(state)
 *
 * // 修改值会自动记录历史
 * current.value.count++
 *
 * // 撤销
 * undo()
 *
 * // 重做
 * redo()
 */
export function useHistory<T>(source: Ref<T>, options: UseHistoryOptions<T> = {}): UseHistoryReturn<T> {
  const { capacity = Infinity, deep = true, immediate = true, clone = defaultClone, equals = defaultEquals } = options;

  // 历史记录栈
  const history = ref<T[]>([]) as Ref<T[]>;
  const index = ref(-1);

  // 是否正在撤销/重做（避免触发 watch）
  let isUpdating = false;

  // 克隆值
  const cloneValue = (value: T): T => {
    try {
      return clone(value);
    } catch {
      return value;
    }
  };

  // 提交当前状态
  const commit = () => {
    const currentValue = source.value;

    // 如果与最后一条记录相同，不记录
    if (history.value.length > 0 && index.value >= 0) {
      const lastValue = history.value[index.value];
      if (equals(currentValue, lastValue)) {
        return;
      }
    }

    // 如果不在历史末尾，清除后面的记录
    if (index.value < history.value.length - 1) {
      history.value = history.value.slice(0, index.value + 1);
    }

    // 添加新记录
    history.value.push(cloneValue(currentValue));
    index.value = history.value.length - 1;

    // 限制容量
    if (capacity !== Infinity && history.value.length > capacity) {
      const overflow = history.value.length - capacity;
      history.value = history.value.slice(overflow);
      index.value = Math.max(0, index.value - overflow);
    }
  };

  // 是否可以撤销
  const canUndo = computed(() => index.value > 0);

  // 是否可以重做
  const canRedo = computed(() => index.value < history.value.length - 1);

  // 历史记录数量
  const historyCount = computed(() => history.value.length);

  // 撤销
  const undo = () => {
    if (!canUndo.value) return;

    isUpdating = true;
    index.value--;
    source.value = cloneValue(history.value[index.value]);
    isUpdating = false;
  };

  // 重做
  const redo = () => {
    if (!canRedo.value) return;

    isUpdating = true;
    index.value++;
    source.value = cloneValue(history.value[index.value]);
    isUpdating = false;
  };

  // 跳转
  const go = (delta: number) => {
    const newIndex = index.value + delta;
    if (newIndex < 0 || newIndex >= history.value.length) return;

    isUpdating = true;
    index.value = newIndex;
    source.value = cloneValue(history.value[index.value]);
    isUpdating = false;
  };

  // 清除历史
  const clear = () => {
    history.value = [cloneValue(source.value)];
    index.value = 0;
  };

  // 重置到初始状态
  const reset = () => {
    if (history.value.length === 0) return;

    isUpdating = true;
    index.value = 0;
    source.value = cloneValue(history.value[0]);
    isUpdating = false;
  };

  // 监听值变化
  watch(
    source,
    () => {
      if (!isUpdating) {
        commit();
      }
    },
    { deep, immediate }
  );

  return {
    current: source,
    canUndo,
    canRedo,
    historyCount,
    undo,
    redo,
    clear,
    commit,
    reset,
    go,
    history: computed(() => history.value),
    index,
  };
}

// ============================================================================
// 手动历史记录
// ============================================================================

/**
 * 手动控制的历史记录
 * 不自动监听，需要手动调用 commit
 * @example
 * const { current, commit, undo, redo } = useManualHistory({ count: 0 })
 *
 * // 修改值
 * current.value.count++
 *
 * // 手动提交
 * commit()
 */
export function useManualHistory<T>(initialValue: T, options: Omit<UseHistoryOptions<T>, "immediate"> = {}) {
  const source = ref(initialValue) as Ref<T>;
  const history = useHistory(source, { ...options, immediate: false });

  // 初始提交
  history.commit();

  return history;
}

// ============================================================================
// 带时间戳的历史记录
// ============================================================================

export interface HistoryEntry<T> {
  value: T;
  timestamp: number;
}

/**
 * 带时间戳的历史记录
 */
export function useTimestampedHistory<T>(source: Ref<T>, options: UseHistoryOptions<T> = {}) {
  const entries = ref<HistoryEntry<T>[]>([]) as Ref<HistoryEntry<T>[]>;
  const index = ref(-1);

  const history = useHistory(source, {
    ...options,
    immediate: false,
  });

  // 包装 commit
  const commitWithTimestamp = () => {
    history.commit();

    // 同步更新带时间戳的记录
    if (history.index.value >= entries.value.length) {
      entries.value.push({
        value: history.history.value[history.index.value],
        timestamp: Date.now(),
      });
    }

    index.value = history.index.value;
  };

  // 初始提交
  commitWithTimestamp();

  // 监听变化
  watch(source, commitWithTimestamp, { deep: options.deep ?? true });

  return {
    ...history,
    entries: computed(() => entries.value),
    commit: commitWithTimestamp,
  };
}

// ============================================================================
// 状态快照
// ============================================================================

export interface UseSnapshotOptions<T> {
  /** 最大快照数量，默认 10 */
  maxSnapshots?: number;
  /** 自定义克隆函数 */
  clone?: (value: T) => T;
}

/**
 * 状态快照
 * 保存和恢复命名快照
 * @example
 * const { save, restore, list, remove } = useSnapshot(state)
 *
 * // 保存快照
 * save('before-edit')
 *
 * // 恢复快照
 * restore('before-edit')
 */
export function useSnapshot<T>(source: Ref<T>, options: UseSnapshotOptions<T> = {}) {
  const { maxSnapshots = 10, clone = defaultClone } = options;

  const snapshots = ref<Map<string, T>>(new Map()) as Ref<Map<string, T>>;

  // 保存快照
  const save = (name: string) => {
    // 限制数量
    if (snapshots.value.size >= maxSnapshots && !snapshots.value.has(name)) {
      const firstKey = snapshots.value.keys().next().value;
      if (firstKey) {
        snapshots.value.delete(firstKey);
      }
    }

    snapshots.value.set(name, clone(source.value));
  };

  // 恢复快照
  const restore = (name: string): boolean => {
    const snapshot = snapshots.value.get(name);
    if (snapshot === undefined) return false;

    source.value = clone(snapshot);
    return true;
  };

  // 删除快照
  const remove = (name: string): boolean => {
    return snapshots.value.delete(name);
  };

  // 列出所有快照
  const list = computed(() => Array.from(snapshots.value.keys()));

  // 检查快照是否存在
  const has = (name: string): boolean => {
    return snapshots.value.has(name);
  };

  // 清除所有快照
  const clear = () => {
    snapshots.value.clear();
  };

  // 获取快照
  const get = (name: string): T | undefined => {
    const snapshot = snapshots.value.get(name);
    return snapshot ? clone(snapshot) : undefined;
  };

  return {
    save,
    restore,
    remove,
    list,
    has,
    clear,
    get,
    count: computed(() => snapshots.value.size),
  };
}
