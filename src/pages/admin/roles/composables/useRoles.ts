/**
 * Admin 角色管理 Composable
 */
import { ref, reactive, watch } from "vue";
import { adminRoleApi, extractList, extractData } from "@/api";
import type { PaginationState } from "@/api";
import type { RoleRoleDTO, RolePermissionDTO, RoleCreateRoleDTO, RoleUpdateRoleDTO, RoleSetPermissionsDTO } from "@models";
import { exportToCSV, type CSVColumn } from "@/utils/export";
import { refDebounced } from "@vueuse/core";

export function useRoles() {
  const roles = ref<RoleRoleDTO[]>([]);
  const loading = ref(false);
  const searchQuery = ref("");
  // 防抖搜索值，300ms 延迟
  const debouncedSearchQuery = refDebounced(searchQuery, 300);
  const pagination = reactive<PaginationState>({
    page: 1,
    limit: 20,
    total: 0,
    total_pages: 0,
  });

  const errorMessage = ref("");
  const successMessage = ref("");

  const fetchRoles = async () => {
    loading.value = true;
    errorMessage.value = "";

    try {
      const response = await adminRoleApi.apiAdminRolesGet(pagination.page, pagination.limit);
      const result = extractList<RoleRoleDTO>(response.data);
      roles.value = result.data;
      Object.assign(pagination, result.pagination);
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取角色列表失败";
      console.error("Failed to fetch roles:", error);
    } finally {
      loading.value = false;
    }
  };

  const fetchRole = async (id: number): Promise<RoleRoleDTO | null> => {
    try {
      const response = await adminRoleApi.apiAdminRolesIdGet(id);
      return extractData<RoleRoleDTO>(response.data) ?? null;
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取角色详情失败";
      return null;
    }
  };

  const createRole = async (data: RoleCreateRoleDTO): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await adminRoleApi.apiAdminRolesPost(data);
      successMessage.value = "角色创建成功";
      await fetchRoles();
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "创建角色失败";
      return false;
    } finally {
      loading.value = false;
    }
  };

  const updateRole = async (id: number, data: RoleUpdateRoleDTO): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await adminRoleApi.apiAdminRolesIdPut(id, data);
      successMessage.value = "角色更新成功";
      await fetchRoles();
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "更新角色失败";
      return false;
    } finally {
      loading.value = false;
    }
  };

  const deleteRole = async (id: number): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await adminRoleApi.apiAdminRolesIdDelete(id);
      successMessage.value = "角色删除成功";
      await fetchRoles();
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "删除角色失败";
      return false;
    } finally {
      loading.value = false;
    }
  };

  const setPermissions = async (id: number, permissionIds: number[]): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      const data: RoleSetPermissionsDTO = { permission_ids: permissionIds };
      await adminRoleApi.apiAdminRolesIdPermissionsPut(id, data);
      successMessage.value = "权限设置成功";
      await fetchRoles();
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "设置权限失败";
      return false;
    } finally {
      loading.value = false;
    }
  };

  /**
   * 获取权限列表
   */
  const listPermissions = async (params?: { page?: number; limit?: number }) => {
    const response = await adminRoleApi.apiAdminPermissionsGet(params?.page, params?.limit);
    return extractList<RolePermissionDTO>(response.data);
  };

  // 监听防抖搜索值变化，自动触发搜索
  watch(debouncedSearchQuery, () => {
    pagination.page = 1;
    fetchRoles();
  });

  const changePage = (page: number) => {
    pagination.page = page;
    fetchRoles();
  };

  const clearMessages = () => {
    errorMessage.value = "";
    successMessage.value = "";
  };

  /**
   * 导出角色列表为 CSV
   */
  const exportRoles = async () => {
    loading.value = true;
    errorMessage.value = "";

    try {
      // 获取所有角色（最多 1000 条）
      const response = await adminRoleApi.apiAdminRolesGet(1, 1000);
      const result = extractList<RoleRoleDTO>(response.data);

      if (result.data.length === 0) {
        errorMessage.value = "没有可导出的数据";
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
          key: (item) => item.permissions?.map((p) => p.code).join(", ") || "-",
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
      successMessage.value = `成功导出 ${result.data.length} 个角色`;
    } catch (error) {
      errorMessage.value = (error as Error).message || "导出失败";
      console.error("Failed to export roles:", error);
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
    errorMessage,
    successMessage,
    fetchRoles,
    fetchRole,
    createRole,
    updateRole,
    deleteRole,
    setPermissions,
    listPermissions,
    changePage,
    clearMessages,
    exportRoles,
  };
}
