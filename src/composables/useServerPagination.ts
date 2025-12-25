/**
 * 服务端分页 Composable
 *
 * 统一处理 v-data-table-server 的分页逻辑：
 * - 首次请求 vs 后续请求的判断（避免重复请求）
 * - 只更新 total，不覆盖 page/limit（前端是分页状态的唯一真相来源）
 * - 提供正确的 API 参数顺序（limit, page）
 * - 安全处理 -1（Vuetify "All"）等无效值
 */
import { ref, reactive, type Ref } from "vue";
import type { PaginationState } from "@/api";

/** 标准分页选项（不含 -1 "All"，避免服务端分页性能问题） */
export const ITEMS_PER_PAGE_OPTIONS = [30, 50, 100];

export interface UseServerPaginationOptions {
  /** 默认每页条数，默认 30 */
  defaultLimit?: number;
  /** 默认页码，默认 1 */
  defaultPage?: number;
}

export interface UseServerPaginationReturn {
  /** 分页状态（响应式） */
  pagination: PaginationState;
  /** 加载状态 */
  loading: Ref<boolean>;
  /** 是否已初始化（用于区分首次请求和后续请求） */
  initialized: Ref<boolean>;
  /**
   * 处理 v-data-table-server 的 @update:options 事件
   * @param options - 表格选项 { page, itemsPerPage }
   * @param fetchFn - 数据获取函数
   */
  onTableOptionsUpdate: (options: { page: number; itemsPerPage: number }, fetchFn: () => Promise<void>) => void;
  /**
   * 更新 total（从 API 响应中），不覆盖 page/limit
   * @param total - 总记录数
   * @param totalPages - 总页数（可选）
   */
  updateTotal: (total: number, totalPages?: number) => void;
  /**
   * 重置到第一页并执行 fetch
   * @param fetchFn - 数据获取函数
   */
  resetAndFetch: (fetchFn: () => Promise<void>) => void;
  /**
   * 获取分页参数（正确顺序：limit, page）
   * 用于调用 API 时保证参数顺序正确
   */
  getParams: () => { limit: number; page: number };
}

/**
 * 创建服务端分页状态和处理函数
 *
 * @example
 * ```typescript
 * const { pagination, loading, onTableOptionsUpdate, updateTotal, getParams } = useServerPagination();
 *
 * const fetchData = async () => {
 *   loading.value = true;
 *   try {
 *     const { limit, page } = getParams();
 *     const response = await api.getData(limit, page);
 *     data.value = response.data;
 *     updateTotal(response.pagination.total, response.pagination.total_pages);
 *   } finally {
 *     loading.value = false;
 *   }
 * };
 * ```
 */
export function useServerPagination(options: UseServerPaginationOptions = {}): UseServerPaginationReturn {
  const { defaultLimit = 30, defaultPage = 1 } = options;

  const pagination = reactive<PaginationState>({
    page: defaultPage,
    limit: defaultLimit,
    total: 0,
    total_pages: 0,
  });

  const loading = ref(false);
  const initialized = ref(false);

  const onTableOptionsUpdate = (opts: { page: number; itemsPerPage: number }, fetchFn: () => Promise<void>) => {
    // 安全检查：-1 表示 Vuetify "All"，转换为默认值
    const safeLimit = opts.itemsPerPage <= 0 ? defaultLimit : opts.itemsPerPage;
    // 检测值是否真正变化
    const isChanged = opts.page !== pagination.page || safeLimit !== pagination.limit;

    // 首次调用：无条件发起请求
    // 后续调用：只有当值真正变化时才发起请求
    if (!initialized.value || isChanged) {
      pagination.page = opts.page;
      pagination.limit = safeLimit;
      initialized.value = true;
      fetchFn();
    }
  };

  const updateTotal = (total: number, totalPages?: number) => {
    pagination.total = total;
    if (totalPages !== undefined) {
      pagination.total_pages = totalPages;
    }
  };

  const resetAndFetch = (fetchFn: () => Promise<void>) => {
    pagination.page = 1;
    fetchFn();
  };

  const getParams = () => ({
    limit: pagination.limit,
    page: pagination.page,
  });

  return {
    pagination,
    loading,
    initialized,
    onTableOptionsUpdate,
    updateTotal,
    resetAndFetch,
    getParams,
  };
}
