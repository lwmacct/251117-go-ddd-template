/**
 * Admin 组织管理 Composable
 */
import { ref, watch } from "vue";
import { adminOrgApi, extractList, extractData } from "@/api";
import type { OrgOrgDTO, OrgCreateOrgDTO, OrgUpdateOrgDTO } from "@models";
import { exportToCSV, type CSVColumn } from "@/utils/export";
import { refDebounced } from "@vueuse/core";
import { useServerPagination, useSnackbar } from "@/composables";

export function useAdminOrgs() {
  // 状态
  const orgs = ref<OrgOrgDTO[]>([]);
  const searchQuery = ref("");
  const statusFilter = ref<"active" | "inactive" | "">("");
  // 防抖搜索值，300ms 延迟
  const debouncedSearchQuery = refDebounced(searchQuery, 300);

  // 使用通用分页 composable
  const {
    pagination,
    loading,
    onTableOptionsUpdate: baseOnTableOptionsUpdate,
    updateTotal,
    resetAndFetch,
    getParams,
  } = useServerPagination();

  // 消息提示
  const { success, error } = useSnackbar();

  /**
   * 获取组织列表
   */
  const fetchOrgs = async () => {
    loading.value = true;

    try {
      const { limit, page } = getParams();
      const response = await adminOrgApi.apiSystemOrgsGet(limit, page, statusFilter.value || undefined);
      const result = extractList<OrgOrgDTO>(response.data);
      orgs.value = result.data;
      updateTotal(result.pagination.total, result.pagination.total_pages);
    } catch (err) {
      error((err as Error).message || "获取组织列表失败");
      console.error("Failed to fetch orgs:", err);
    } finally {
      loading.value = false;
    }
  };

  /**
   * 获取单个组织详情
   */
  const fetchOrg = async (id: number): Promise<OrgOrgDTO | null> => {
    try {
      const response = await adminOrgApi.apiSystemOrgsIdGet(id);
      return extractData<OrgOrgDTO>(response.data) ?? null;
    } catch (err) {
      error((err as Error).message || "获取组织详情失败");
      console.error("Failed to fetch org:", err);
      return null;
    }
  };

  /**
   * 创建组织
   */
  const createOrg = async (data: OrgCreateOrgDTO): Promise<boolean> => {
    loading.value = true;

    try {
      await adminOrgApi.apiSystemOrgsPost(data);
      success("组织创建成功");
      await fetchOrgs(); // 刷新列表
      return true;
    } catch (err) {
      error((err as Error).message || "创建组织失败");
      console.error("Failed to create org:", err);
      return false;
    } finally {
      loading.value = false;
    }
  };

  /**
   * 更新组织
   */
  const updateOrg = async (id: number, data: OrgUpdateOrgDTO): Promise<boolean> => {
    loading.value = true;

    try {
      await adminOrgApi.apiSystemOrgsIdPut(id, data);
      success("组织更新成功");
      await fetchOrgs(); // 刷新列表
      return true;
    } catch (err) {
      error((err as Error).message || "更新组织失败");
      console.error("Failed to update org:", err);
      return false;
    } finally {
      loading.value = false;
    }
  };

  /**
   * 删除组织
   */
  const deleteOrg = async (id: number): Promise<boolean> => {
    loading.value = true;

    try {
      await adminOrgApi.apiSystemOrgsIdDelete(id);
      success("组织删除成功");
      await fetchOrgs(); // 刷新列表
      return true;
    } catch (err) {
      error((err as Error).message || "删除组织失败");
      console.error("Failed to delete org:", err);
      return false;
    } finally {
      loading.value = false;
    }
  };

  // 监听搜索值变化，自动触发搜索
  watch(debouncedSearchQuery, () => {
    resetAndFetch(fetchOrgs);
  });

  // 监听状态筛选变化
  watch(statusFilter, () => {
    resetAndFetch(fetchOrgs);
  });

  /**
   * 表格选项变化处理（分页、每页条数、排序）
   * 由 v-data-table-server 的 @update:options 触发
   */
  const onTableOptionsUpdate = (options: { page: number; itemsPerPage: number }) => {
    baseOnTableOptionsUpdate(options, fetchOrgs);
  };

  /**
   * 导出组织列表为 CSV
   */
  const exportOrgs = async () => {
    loading.value = true;

    try {
      // 获取所有组织（最多 1000 条）
      const response = await adminOrgApi.apiSystemOrgsGet(1000, 1);
      const result = extractList<OrgOrgDTO>(response.data);

      if (result.data.length === 0) {
        error("没有可导出的数据");
        return;
      }

      // 过滤数据
      let exportData = result.data;
      if (statusFilter.value) {
        exportData = exportData.filter((o) => o.status === statusFilter.value);
      }

      // 定义 CSV 列
      const columns: CSVColumn<OrgOrgDTO>[] = [
        { header: "ID", key: "id" },
        { header: "组织标识", key: "name" },
        { header: "组织名称", key: "display_name" },
        { header: "描述", key: "description" },
        {
          header: "状态",
          key: (item) => (item.status === "active" ? "启用" : "禁用"),
        },
        {
          header: "创建时间",
          key: (item) =>
            item.created_at
              ? new Date(item.created_at).toLocaleString("zh-CN", {
                  year: "numeric",
                  month: "2-digit",
                  day: "2-digit",
                  hour: "2-digit",
                  minute: "2-digit",
                })
              : "-",
        },
      ];

      // 生成文件名
      const timestamp = new Date().toISOString().slice(0, 10);
      const filename = `组织列表_${timestamp}.csv`;

      // 导出
      exportToCSV(exportData, columns, { filename, withBOM: true });
      success(`成功导出 ${exportData.length} 条组织记录`);
    } catch (err) {
      error((err as Error).message || "导出失败");
      console.error("Failed to export orgs:", err);
    } finally {
      loading.value = false;
    }
  };

  return {
    // 状态
    orgs,
    loading,
    searchQuery,
    debouncedSearchQuery,
    statusFilter,
    pagination,

    // 方法
    fetchOrgs,
    fetchOrg,
    createOrg,
    updateOrg,
    deleteOrg,
    onTableOptionsUpdate,
    exportOrgs,
  };
}
