/**
 * Admin 用户管理 Composable
 */
import { ref, reactive, computed } from 'vue';
import { AdminUsersAPI } from '@/api/admin';
import type { AdminUser, CreateUserRequest, UpdateUserRequest } from '@/types/admin';
import type { PaginationMeta } from '@/types/common';

export function useAdminUsers() {
  // 状态
  const users = ref<AdminUser[]>([]);
  const loading = ref(false);
  const searchQuery = ref('');
  const pagination = reactive<PaginationMeta>({
    page: 1,
    limit: 20,
    total: 0,
    total_pages: 0,
  });

  // 错误和成功消息
  const errorMessage = ref('');
  const successMessage = ref('');

  /**
   * 获取用户列表
   */
  const fetchUsers = async () => {
    loading.value = true;
    errorMessage.value = '';

    try {
      const response = await AdminUsersAPI.listUsers({
        page: pagination.page,
        limit: pagination.limit,
        search: searchQuery.value || undefined,
      });

      users.value = response.data;
      Object.assign(pagination, response.pagination);
    } catch (error: any) {
      errorMessage.value = error.message || '获取用户列表失败';
      console.error('Failed to fetch users:', error);
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
    } catch (error: any) {
      errorMessage.value = error.message || '获取用户详情失败';
      console.error('Failed to fetch user:', error);
      return null;
    }
  };

  /**
   * 创建用户
   */
  const createUser = async (data: CreateUserRequest): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = '';
    successMessage.value = '';

    try {
      await AdminUsersAPI.createUser(data);
      successMessage.value = '用户创建成功';
      await fetchUsers(); // 刷新列表
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || '创建用户失败';
      console.error('Failed to create user:', error);
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
    errorMessage.value = '';
    successMessage.value = '';

    try {
      await AdminUsersAPI.updateUser(id, data);
      successMessage.value = '用户更新成功';
      await fetchUsers(); // 刷新列表
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || '更新用户失败';
      console.error('Failed to update user:', error);
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
    errorMessage.value = '';
    successMessage.value = '';

    try {
      await AdminUsersAPI.deleteUser(id);
      successMessage.value = '用户删除成功';
      await fetchUsers(); // 刷新列表
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || '删除用户失败';
      console.error('Failed to delete user:', error);
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
    errorMessage.value = '';
    successMessage.value = '';

    try {
      await AdminUsersAPI.assignRoles(id, { role_ids: roleIds });
      successMessage.value = '角色分配成功';
      await fetchUsers(); // 刷新列表
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || '角色分配失败';
      console.error('Failed to assign roles:', error);
      return false;
    } finally {
      loading.value = false;
    }
  };

  /**
   * 搜索用户
   */
  const searchUsers = (query: string) => {
    searchQuery.value = query;
    pagination.page = 1; // 重置到第一页
    fetchUsers();
  };

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
    errorMessage.value = '';
    successMessage.value = '';
  };

  return {
    // 状态
    users,
    loading,
    searchQuery,
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
    searchUsers,
    changePage,
    clearMessages,
  };
}
