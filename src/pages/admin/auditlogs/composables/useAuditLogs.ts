/**
 * 审计日志管理 Composable
 *
 * 提供审计日志列表查询、过滤、分页、导出等功能
 */
import { ref, reactive } from "vue";
import {
  adminAuditLogApi,
  extractList,
  extractData,
  type AuditAuditDTO,
  type AuditAuditActionDTO,
  type AuditCategoryOptionDTO,
  type AuditAuditActionsResponseDTO,
  type AuditStatus,
} from "@/api";
import { exportToCSV, formatDateForExport, type CSVColumn } from "@/utils/export";
import { useServerPagination, useSnackbar } from "@/composables";

interface AuditLogQueryParams {
  page?: number;
  limit?: number;
  user_id?: number;
  action?: string;
  resource?: string;
  status?: AuditStatus;
  start_date?: string;
  end_date?: string;
}

export function useAuditLogs() {
  // 状态管理
  const logs = ref<AuditAuditDTO[]>([]);
  const selectedLog = ref<AuditAuditDTO | null>(null);
  const exporting = ref(false);

  // 消息提示
  const { success, error } = useSnackbar();

  // 筛选选项（从 API 动态获取）
  const actionOptions = ref<AuditAuditActionDTO[]>([]);
  const categoryOptions = ref<AuditCategoryOptionDTO[]>([]);

  // 使用通用分页 composable
  const {
    pagination,
    loading,
    onTableOptionsUpdate: baseOnTableOptionsUpdate,
    updateTotal,
    resetAndFetch,
    getParams,
  } = useServerPagination();

  // 过滤条件
  const filters = reactive<AuditLogQueryParams>({
    user_id: undefined,
    action: "",
    resource: "",
    status: "",
    start_date: "",
    end_date: "",
  });

  /**
   * 获取筛选选项（操作类型、资源分类）
   */
  const fetchFilterOptions = async () => {
    try {
      const response = await adminAuditLogApi.apiSystemAuditActionsGet();
      const result = extractData<AuditAuditActionsResponseDTO>(response.data);
      actionOptions.value = result?.actions ?? [];
      categoryOptions.value = result?.categories ?? [];
    } catch (error) {
      console.error("Failed to fetch filter options:", error);
    }
  };

  /**
   * 获取审计日志列表
   */
  const fetchLogs = async () => {
    loading.value = true;

    try {
      const { limit, page } = getParams();
      // 参数按字母顺序：action, endDate, limit, page, resource, startDate, status, userId
      const response = await adminAuditLogApi.apiSystemAuditGet(
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (filters.action || undefined) as any,
        filters.end_date || undefined,
        limit,
        page,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (filters.resource || undefined) as any,
        filters.start_date || undefined,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (filters.status || undefined) as any,
        filters.user_id,
      );
      const result = extractList<AuditAuditDTO>(response.data);
      logs.value = result.data;
      updateTotal(result.pagination.total, result.pagination.total_pages);
    } catch (err) {
      error((err as Error).message || "获取审计日志失败");
      console.error("Failed to fetch audit logs:", err);
    } finally {
      loading.value = false;
    }
  };

  /**
   * 获取日志详情
   */
  const fetchLogDetail = async (id: number) => {
    loading.value = true;

    try {
      const response = await adminAuditLogApi.apiSystemAuditIdGet(id);
      selectedLog.value = extractData<AuditAuditDTO>(response.data) ?? null;
    } catch (err) {
      error((err as Error).message || "获取日志详情失败");
      console.error("Failed to fetch log detail:", err);
    } finally {
      loading.value = false;
    }
  };

  /**
   * 应用过滤条件
   */
  const applyFilters = () => {
    resetAndFetch(fetchLogs);
  };

  /**
   * 重置过滤条件
   */
  const resetFilters = () => {
    filters.user_id = undefined;
    filters.action = "";
    filters.resource = "";
    filters.status = "";
    filters.start_date = "";
    filters.end_date = "";
    resetAndFetch(fetchLogs);
  };

  /**
   * 表格选项变化处理（分页、每页条数、排序）
   * 由 v-data-table-server 的 @update:options 触发
   */
  const onTableOptionsUpdate = (options: { page: number; itemsPerPage: number }) => {
    baseOnTableOptionsUpdate(options, fetchLogs);
  };

  /**
   * 导出审计日志为 CSV
   * 导出当前过滤条件下的数据（最多 1000 条）
   */
  const exportLogs = async () => {
    exporting.value = true;

    try {
      // 参数按字母顺序：action, endDate, limit, page, resource, startDate, status, userId
      const response = await adminAuditLogApi.apiSystemAuditGet(
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (filters.action || undefined) as any,
        filters.end_date || undefined,
        1000,
        1,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (filters.resource || undefined) as any,
        filters.start_date || undefined,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (filters.status || undefined) as any,
        filters.user_id,
      );
      const result = extractList<AuditAuditDTO>(response.data);

      if (result.data.length === 0) {
        error("没有数据可导出");
        return;
      }

      // 定义 CSV 列
      const columns: CSVColumn<AuditAuditDTO>[] = [
        { header: "ID", key: "id" },
        { header: "用户ID", key: "user_id" },
        { header: "操作类型", key: "action" },
        { header: "资源类型", key: "resource" },
        { header: "状态", key: "status" },
        { header: "IP地址", key: (item) => item.ip_address || "" },
        { header: "UserUserWithRolesDTO Agent", key: (item) => item.user_agent || "" },
        { header: "详情", key: (item) => item.details || "" },
        { header: "创建时间", key: (item) => formatDateForExport(item.created_at) },
      ];

      // 生成文件名（包含时间戳）
      const timestamp = new Date().toISOString().slice(0, 10);
      const filename = `audit-logs-${timestamp}`;

      // 导出 CSV
      exportToCSV(result.data, columns, { filename });

      success(`成功导出 ${result.data.length} 条记录`);
    } catch (err) {
      error((err as Error).message || "导出失败");
      console.error("Failed to export audit logs:", err);
    } finally {
      exporting.value = false;
    }
  };

  return {
    // 状态
    logs,
    selectedLog,
    loading,
    exporting,
    filters,
    pagination,

    // 筛选选项
    actionOptions,
    categoryOptions,

    // 方法
    fetchLogs,
    fetchLogDetail,
    fetchFilterOptions,
    applyFilters,
    resetFilters,
    onTableOptionsUpdate,
    exportLogs,
  };
}
