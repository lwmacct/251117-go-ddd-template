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
  type AuditlogAuditLogDTO,
  type AuditAction,
  type AuditStatus,
} from "@/api";
import { exportToCSV, formatDateForExport, type CSVColumn } from "@/utils/export";

interface AuditLogQueryParams {
  page?: number;
  limit?: number;
  user_id?: number;
  action?: AuditAction;
  resource?: string;
  status?: AuditStatus;
  start_date?: string;
  end_date?: string;
}

export function useAuditLogs() {
  // 状态管理
  const logs = ref<AuditlogAuditLogDTO[]>([]);
  const selectedLog = ref<AuditlogAuditLogDTO | null>(null);
  const loading = ref(false);
  const exporting = ref(false);
  const errorMessage = ref("");
  const successMessage = ref("");

  // 过滤条件
  const filters = reactive<AuditLogQueryParams>({
    user_id: undefined,
    action: "",
    resource: "",
    status: "",
    start_date: "",
    end_date: "",
  });

  // 分页状态
  const pagination = reactive({
    page: 1,
    limit: 20,
    total: 0,
  });

  /**
   * 获取审计日志列表
   */
  const fetchLogs = async () => {
    loading.value = true;
    errorMessage.value = "";

    try {
      const response = await adminAuditLogApi.apiAdminAuditlogsGet(
        pagination.page,
        pagination.limit,
        filters.user_id,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (filters.action || undefined) as any,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (filters.resource || undefined) as any,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (filters.status || undefined) as any,
        filters.start_date || undefined,
        filters.end_date || undefined,
      );
      const result = extractList<AuditlogAuditLogDTO>(response.data);
      logs.value = result.data;
      pagination.total = result.pagination.total;
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取审计日志失败";
      console.error("Failed to fetch audit logs:", error);
    } finally {
      loading.value = false;
    }
  };

  /**
   * 获取日志详情
   */
  const fetchLogDetail = async (id: number) => {
    loading.value = true;
    errorMessage.value = "";

    try {
      const response = await adminAuditLogApi.apiAdminAuditlogsIdGet(id);
      selectedLog.value = extractData<AuditlogAuditLogDTO>(response.data) ?? null;
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取日志详情失败";
      console.error("Failed to fetch log detail:", error);
    } finally {
      loading.value = false;
    }
  };

  /**
   * 应用过滤条件
   */
  const applyFilters = () => {
    pagination.page = 1; // 重置到第一页
    fetchLogs();
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
    pagination.page = 1;
    fetchLogs();
  };

  /**
   * 切换页码
   */
  const changePage = (page: number) => {
    pagination.page = page;
    fetchLogs();
  };

  /**
   * 清除消息提示
   */
  const clearMessages = () => {
    errorMessage.value = "";
    successMessage.value = "";
  };

  /**
   * 导出审计日志为 CSV
   * 导出当前过滤条件下的数据（最多 1000 条）
   */
  const exportLogs = async () => {
    exporting.value = true;
    errorMessage.value = "";

    try {
      // 使用当前过滤条件，获取最多 1000 条数据
      const response = await adminAuditLogApi.apiAdminAuditlogsGet(
        1,
        1000,
        filters.user_id,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (filters.action || undefined) as any,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (filters.resource || undefined) as any,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (filters.status || undefined) as any,
        filters.start_date || undefined,
        filters.end_date || undefined,
      );
      const result = extractList<AuditlogAuditLogDTO>(response.data);

      if (result.data.length === 0) {
        errorMessage.value = "没有数据可导出";
        return;
      }

      // 定义 CSV 列
      const columns: CSVColumn<AuditlogAuditLogDTO>[] = [
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

      successMessage.value = `成功导出 ${result.data.length} 条记录`;
    } catch (error) {
      errorMessage.value = (error as Error).message || "导出失败";
      console.error("Failed to export audit logs:", error);
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
    errorMessage,
    successMessage,
    filters,
    pagination,

    // 方法
    fetchLogs,
    fetchLogDetail,
    applyFilters,
    resetFilters,
    changePage,
    clearMessages,
    exportLogs,
  };
}
