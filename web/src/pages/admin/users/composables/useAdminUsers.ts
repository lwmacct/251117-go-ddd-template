/**
 * Admin 用户管理 Composable
 */
import { ref, reactive, watch } from "vue";
import { AdminUsersAPI } from "@/api/admin";
import { exportToCSV, type CSVColumn } from "@/utils/export";
import { refDebounced } from "@vueuse/core";
import type { AdminUser, CreateUserRequest, UpdateUserRequest } from "@/types/admin";
import type { PaginationMeta } from "@/types/common";

export function useAdminUsers() {
  // 状态
  const users = ref<AdminUser[]>([]);
  const loading = ref(false);
  const searchQuery = ref("");
  // 防抖搜索值，300ms 延迟
  const debouncedSearchQuery = refDebounced(searchQuery, 300);
  const pagination = reactive<PaginationMeta>({
    page: 1,
    limit: 20,
    total: 0,
    total_pages: 0,
  });

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
      const response = await AdminUsersAPI.listUsers({
        page: pagination.page,
        limit: pagination.limit,
        search: debouncedSearchQuery.value || undefined,
      });

      users.value = response.data;
      Object.assign(pagination, response.pagination);
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
  const fetchUser = async (id: number): Promise<AdminUser | null> => {
    try {
      return await AdminUsersAPI.getUser(id);
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取用户详情失败";
      console.error("Failed to fetch user:", error);
      return null;
    }
  };

  /**
   * 创建用户
   */
  const createUser = async (data: CreateUserRequest): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await AdminUsersAPI.createUser(data);
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
  const updateUser = async (id: number, data: UpdateUserRequest): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = "";
    successMessage.value = "";

    try {
      await AdminUsersAPI.updateUser(id, data);
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
      await AdminUsersAPI.deleteUser(id);
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
      await AdminUsersAPI.assignRoles(id, { role_ids: roleIds });
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
    pagination.page = 1; // 重置到第一页
    fetchUsers();
  });

  /**
   * 翻页
   */
  const changePage = (page: number) => {
    pagination.page = page;
    fetchUsers();
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
      const response = await AdminUsersAPI.listUsers({
        page: 1,
        limit: 1000,
        search: searchQuery.value || undefined,
      });

      if (response.data.length === 0) {
        errorMessage.value = "没有可导出的数据";
        return;
      }

      // 定义 CSV 列
      const columns: CSVColumn<AdminUser>[] = [
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
      exportToCSV(response.data, columns, { filename, withBOM: true });
      successMessage.value = `成功导出 ${response.data.length} 条用户记录`;
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
    changePage,
    clearMessages,
    exportUsers,
  };
}
