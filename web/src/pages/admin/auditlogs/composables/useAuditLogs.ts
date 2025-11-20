/**
 * 审计日志管理 Composable
 *
 * 提供审计日志列表查询、过滤、分页等功能
 */
import { ref, reactive } from "vue";
import { listAuditLogs, getAuditLog } from "@/api/admin/auditlogs";
import type { AuditLog, AuditLogQueryParams } from "@/types/admin/audit";

export function useAuditLogs() {
  // 状态管理
  const logs = ref<AuditLog[]>([]);
  const selectedLog = ref<AuditLog | null>(null);
  const loading = ref(false);
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
      const queryParams: AuditLogQueryParams = {
        page: pagination.page,
        limit: pagination.limit,
        user_id: filters.user_id,
        action: filters.action || undefined,
        resource: filters.resource || undefined,
        status: filters.status || undefined,
        start_date: filters.start_date || undefined,
        end_date: filters.end_date || undefined,
      };

      const response = await listAuditLogs(queryParams);
      logs.value = response.data;
      pagination.total = response.pagination.total;
    } catch (error: any) {
      errorMessage.value = error.message || "获取审计日志失败";
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
      selectedLog.value = await getAuditLog(id);
    } catch (error: any) {
      errorMessage.value = error.message || "获取日志详情失败";
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

  return {
    // 状态
    logs,
    selectedLog,
    loading,
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
  };
}
