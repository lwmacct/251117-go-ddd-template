/**
 * Admin 角色管理 Composable
 */
import { ref, watch } from "vue";
import { adminRoleApi, extractList, extractData } from "@/api";
import type { RoleRoleDTO, RoleCreateDTO, RoleUpdateDTO, RolePermissionInputDTO } from "@models";
import { exportToCSV, type CSVColumn } from "@/utils/export";
import { refDebounced } from "@vueuse/core";
import { useServerPagination, useSnackbar } from "@/composables";

export function useRoles() {
  const roles = ref<RoleRoleDTO[]>([]);
  const searchQuery = ref("");
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

  const fetchRoles = async () => {
    loading.value = true;

    try {
      const { limit, page } = getParams();
      const response = await adminRoleApi.apiSystemRolesGet(limit, page);
      const result = extractList<RoleRoleDTO>(response.data);
      roles.value = result.data;
      updateTotal(result.pagination.total, result.pagination.total_pages);
    } catch (err) {
      error((err as Error).message || "获取角色列表失败");
      console.error("Failed to fetch roles:", err);
    } finally {
      loading.value = false;
    }
  };

  const fetchRole = async (id: number): Promise<RoleRoleDTO | null> => {
    try {
      const response = await adminRoleApi.apiSystemRolesIdGet(id);
      return extractData<RoleRoleDTO>(response.data) ?? null;
    } catch (err) {
      error((err as Error).message || "获取角色详情失败");
      return null;
    }
  };

  const createRole = async (data: RoleCreateDTO): Promise<boolean> => {
    loading.value = true;

    try {
      await adminRoleApi.apiSystemRolesPost(data);
      success("角色创建成功");
      await fetchRoles();
      return true;
    } catch (err) {
      error((err as Error).message || "创建角色失败");
      return false;
    } finally {
      loading.value = false;
    }
  };

  const updateRole = async (id: number, data: RoleUpdateDTO): Promise<boolean> => {
    loading.value = true;

    try {
      await adminRoleApi.apiSystemRolesIdPut(id, data);
      success("角色更新成功");
      await fetchRoles();
      return true;
    } catch (err) {
      error((err as Error).message || "更新角色失败");
      return false;
    } finally {
      loading.value = false;
    }
  };

  const deleteRole = async (id: number): Promise<boolean> => {
    loading.value = true;

    try {
      await adminRoleApi.apiSystemRolesIdDelete(id);
      success("角色删除成功");
      await fetchRoles();
      return true;
    } catch (err) {
      error((err as Error).message || "删除角色失败");
      return false;
    } finally {
      loading.value = false;
    }
  };

  const setPermissions = async (id: number, permissions: RolePermissionInputDTO[]): Promise<boolean> => {
    loading.value = true;

    try {
      await adminRoleApi.apiSystemRolesIdPermissionsPut(id, { permissions });
      success("权限设置成功");
      await fetchRoles();
      return true;
    } catch (err) {
      error((err as Error).message || "设置权限失败");
      return false;
    } finally {
      loading.value = false;
    }
  };

  // 监听防抖搜索值变化，自动触发搜索
  watch(debouncedSearchQuery, () => {
    resetAndFetch(fetchRoles);
  });

  /**
   * 表格选项变化处理
   */
  const onTableOptionsUpdate = (options: { page: number; itemsPerPage: number }) => {
    baseOnTableOptionsUpdate(options, fetchRoles);
  };

  /**
   * 导出角色列表为 CSV
   */
  const exportRoles = async () => {
    loading.value = true;

    try {
      // 获取所有角色（最多 1000 条）
      const response = await adminRoleApi.apiSystemRolesGet(1000, 1);
      const result = extractList<RoleRoleDTO>(response.data);

      if (result.data.length === 0) {
        error("没有可导出的数据");
        return;
      }

      // 定义 CSV 列
      const columns: CSVColumn<RoleRoleDTO>[] = [
        { header: "ID", key: "id" },
        { header: "角色标识", key: "name" },
        { header: "显示名称", key: "display_name" },
        { header: "描述", key: (item) => item.description || "-" },
        { header: "系统角色", key: (item) => (item.is_system ? "是" : "否") },
        {
          header: "权限数量",
          key: (item) => item.permissions?.length || 0,
        },
        {
          header: "权限列表",
          key: (item) => item.permissions?.map((p) => p.operation_pattern).join(", ") || "-",
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
      const filename = `角色列表_${timestamp}.csv`;

      // 导出
      exportToCSV(result.data, columns, { filename, withBOM: true });
      success(`成功导出 ${result.data.length} 个角色`);
    } catch (err) {
      error((err as Error).message || "导出失败");
      console.error("Failed to export roles:", err);
    } finally {
      loading.value = false;
    }
  };

  return {
    roles,
    loading,
    searchQuery,
    debouncedSearchQuery,
    pagination,
    fetchRoles,
    fetchRole,
    createRole,
    updateRole,
    deleteRole,
    setPermissions,
    onTableOptionsUpdate,
    exportRoles,
  };
}
