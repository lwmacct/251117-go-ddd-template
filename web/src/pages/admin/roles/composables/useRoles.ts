/**
 * Admin 角色管理 Composable
 */
import { ref, reactive } from 'vue';
import { AdminRolesAPI } from '@/api/admin';
import type { Role, CreateRoleRequest, UpdateRoleRequest } from '@/types/admin';
import type { PaginationMeta } from '@/types/common';

export function useRoles() {
  const roles = ref<Role[]>([]);
  const loading = ref(false);
  const searchQuery = ref('');
  const pagination = reactive<PaginationMeta>({
    page: 1,
    limit: 20,
    total: 0,
    total_pages: 0,
  });

  const errorMessage = ref('');
  const successMessage = ref('');

  const fetchRoles = async () => {
    loading.value = true;
    errorMessage.value = '';

    try {
      const response = await AdminRolesAPI.listRoles({
        page: pagination.page,
        limit: pagination.limit,
        search: searchQuery.value || undefined,
      });

      roles.value = response.data;
      Object.assign(pagination, response.pagination);
    } catch (error: any) {
      errorMessage.value = error.message || '获取角色列表失败';
      console.error('Failed to fetch roles:', error);
    } finally {
      loading.value = false;
    }
  };

  const fetchRole = async (id: number): Promise<Role | null> => {
    try {
      return await AdminRolesAPI.getRole(id);
    } catch (error: any) {
      errorMessage.value = error.message || '获取角色详情失败';
      return null;
    }
  };

  const createRole = async (data: CreateRoleRequest): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = '';
    successMessage.value = '';

    try {
      await AdminRolesAPI.createRole(data);
      successMessage.value = '角色创建成功';
      await fetchRoles();
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || '创建角色失败';
      return false;
    } finally {
      loading.value = false;
    }
  };

  const updateRole = async (id: number, data: UpdateRoleRequest): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = '';
    successMessage.value = '';

    try {
      await AdminRolesAPI.updateRole(id, data);
      successMessage.value = '角色更新成功';
      await fetchRoles();
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || '更新角色失败';
      return false;
    } finally {
      loading.value = false;
    }
  };

  const deleteRole = async (id: number): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = '';
    successMessage.value = '';

    try {
      await AdminRolesAPI.deleteRole(id);
      successMessage.value = '角色删除成功';
      await fetchRoles();
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || '删除角色失败';
      return false;
    } finally {
      loading.value = false;
    }
  };

  const setPermissions = async (id: number, permissionIds: number[]): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = '';
    successMessage.value = '';

    try {
      await AdminRolesAPI.setPermissions(id, { permission_ids: permissionIds });
      successMessage.value = '权限设置成功';
      await fetchRoles();
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || '设置权限失败';
      return false;
    } finally {
      loading.value = false;
    }
  };

  const searchRoles = (query: string) => {
    searchQuery.value = query;
    pagination.page = 1;
    fetchRoles();
  };

  const changePage = (page: number) => {
    pagination.page = page;
    fetchRoles();
  };

  const clearMessages = () => {
    errorMessage.value = '';
    successMessage.value = '';
  };

  return {
    roles,
    loading,
    searchQuery,
    pagination,
    errorMessage,
    successMessage,
    fetchRoles,
    fetchRole,
    createRole,
    updateRole,
    deleteRole,
    setPermissions,
    searchRoles,
    changePage,
    clearMessages,
  };
}
