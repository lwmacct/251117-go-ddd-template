/**
 * Admin 用户管理 Composable
 */
import { ref, watch } from "vue";
import { adminUserApi, extractList, extractData } from "@/api";
import type { UserUserWithRolesDTO, UserCreateUserDTO, UserUpdateUserDTO, UserAssignRolesDTO } from "@models";
import { exportToCSV, type CSVColumn } from "@/utils/export";
import { refDebounced } from "@vueuse/core";
import { useServerPagination } from "@/composables";

export function useAdminUsers() {
  // 状态
  const users = ref<UserUserWithRolesDTO[]>([]);
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

  // 错误和成功消息
  const errorMessage = ref("");
  const successMessage = ref("");

  /**
   * 获取用户列表
   */
  const fetchUsers = async () => {
    loading.value = true;
    errorMessage.value = "";

    try {
      const { limit, page } = getParams();
      // 注意：API 参数顺序是 (limit, page, search)
      const response = await adminUserApi.apiAdminUsersGet(limit, page, debouncedSearchQuery.value || undefined);
      const result = extractList<UserUserWithRolesDTO>(response.data);
      users.value = result.data;
      updateTotal(result.pagination.total, result.pagination.total_pages);
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取用户列表失败";
      console.error("Failed to fetch users:", error);
    } finally {
      loading.value = false;
    }
  };

  /**
   * 获取单个用户详情
   */
  const fetchUser = async (id: number): Promise<UserUserWithRolesDTO | null> => {
    try {
      const response = await adminUserApi.apiAdminUsersIdGet(id);
      return extractData<UserUserWithRolesDTO>(response.data) ?? null;
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取用户详情失败";
      console.error("Failed to fetch user:", error);
      return null;
    }
  };

  /**
   * 创建用户
   */
  const createUser = async (data: UserCreateUserDTO): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await adminUserApi.apiAdminUsersPost(data);
      successMessage.value = "用户创建成功";
      await fetchUsers(); // 刷新列表
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "创建用户失败";
      console.error("Failed to create user:", error);
      return false;
    } finally {
      loading.value = false;
    }
  };

  /**
   * 更新用户
   */
  const updateUser = async (id: number, data: UserUpdateUserDTO): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await adminUserApi.apiAdminUsersIdPut(id, data);
      successMessage.value = "用户更新成功";
      await fetchUsers(); // 刷新列表
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "更新用户失败";
      console.error("Failed to update user:", error);
      return false;
    } finally {
      loading.value = false;
    }
  };

  /**
   * 删除用户
   */
  const deleteUser = async (id: number): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await adminUserApi.apiAdminUsersIdDelete(id);
      successMessage.value = "用户删除成功";
      await fetchUsers(); // 刷新列表
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "删除用户失败";
      console.error("Failed to delete user:", error);
      return false;
    } finally {
      loading.value = false;
    }
  };

  /**
   * 分配角色
   */
  const assignRoles = async (id: number, roleIds: number[]): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      const data: UserAssignRolesDTO = { role_ids: roleIds };
      await adminUserApi.apiAdminUsersIdRolesPut(id, data);
      successMessage.value = "角色分配成功";
      await fetchUsers(); // 刷新列表
      return true;
    } catch (error) {
      errorMessage.value = (error as Error).message || "角色分配失败";
      console.error("Failed to assign roles:", error);
      return false;
    } finally {
      loading.value = false;
    }
  };

  // 监听防抖搜索值变化，自动触发搜索
  watch(debouncedSearchQuery, () => {
    resetAndFetch(fetchUsers);
  });

  /**
   * 表格选项变化处理（分页、每页条数、排序）
   * 由 v-data-table-server 的 @update:options 触发
   */
  const onTableOptionsUpdate = (options: { page: number; itemsPerPage: number }) => {
    baseOnTableOptionsUpdate(options, fetchUsers);
  };

  /**
   * 清除消息
   */
  const clearMessages = () => {
    errorMessage.value = "";
    successMessage.value = "";
  };

  /**
   * 导出用户列表为 CSV
   */
  const exportUsers = async () => {
    loading.value = true;
    errorMessage.value = "";

    try {
      // 获取所有用户（最多 1000 条）
      // 注意：API 参数顺序是 (limit, page, search)
      const response = await adminUserApi.apiAdminUsersGet(1000, 1, searchQuery.value || undefined);
      const result = extractList<UserUserWithRolesDTO>(response.data);

      if (result.data.length === 0) {
        errorMessage.value = "没有可导出的数据";
        return;
      }

      // 定义 CSV 列
      const columns: CSVColumn<UserUserWithRolesDTO>[] = [
        { header: "ID", key: "id" },
        { header: "用户名", key: "username" },
        { header: "邮箱", key: "email" },
        { header: "姓名", key: "full_name" },
        {
          header: "状态",
          key: (item) => (item.status === "active" ? "启用" : item.status === "inactive" ? "禁用" : "封禁"),
        },
        {
          header: "角色",
          key: (item) => item.roles?.map((r) => r.display_name).join(", ") || "-",
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
        {
          header: "更新时间",
          key: (item) =>
            item.updated_at
              ? new Date(item.updated_at).toLocaleString("zh-CN", {
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
      const filename = `用户列表_${timestamp}.csv`;

      // 导出
      exportToCSV(result.data, columns, { filename, withBOM: true });
      successMessage.value = `成功导出 ${result.data.length} 条用户记录`;
    } catch (error) {
      errorMessage.value = (error as Error).message || "导出失败";
      console.error("Failed to export users:", error);
    } finally {
      loading.value = false;
    }
  };

  return {
    // 状态
    users,
    loading,
    searchQuery,
    debouncedSearchQuery,
    pagination,
    errorMessage,
    successMessage,

    // 方法
    fetchUsers,
    fetchUser,
    createUser,
    updateUser,
    deleteUser,
    assignRoles,
    onTableOptionsUpdate,
    clearMessages,
    exportUsers,
  };
}
