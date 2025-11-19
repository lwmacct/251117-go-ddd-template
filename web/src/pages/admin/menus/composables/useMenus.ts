/**
 * Admin 菜单管理 Composable
 */
import { ref } from 'vue';
import { AdminMenusAPI } from '@/api/admin';
import type { Menu, CreateMenuRequest, UpdateMenuRequest, ReorderMenusRequest } from '@/types/admin';

export function useMenus() {
  const menus = ref<Menu[]>([]);
  const loading = ref(false);
  const errorMessage = ref('');
  const successMessage = ref('');

  const fetchMenus = async () => {
    loading.value = true;
    errorMessage.value = '';

    try {
      menus.value = await AdminMenusAPI.listMenus();
    } catch (error: any) {
      errorMessage.value = error.message || '获取菜单列表失败';
      console.error('Failed to fetch menus:', error);
    } finally {
      loading.value = false;
    }
  };

  const fetchMenu = async (id: number): Promise<Menu | null> => {
    try {
      return await AdminMenusAPI.getMenu(id);
    } catch (error: any) {
      errorMessage.value = error.message || '获取菜单详情失败';
      return null;
    }
  };

  const createMenu = async (data: CreateMenuRequest): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = '';
    successMessage.value = '';

    try {
      await AdminMenusAPI.createMenu(data);
      successMessage.value = '菜单创建成功';
      await fetchMenus();
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || '创建菜单失败';
      return false;
    } finally {
      loading.value = false;
    }
  };

  const updateMenu = async (id: number, data: UpdateMenuRequest): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = '';
    successMessage.value = '';

    try {
      await AdminMenusAPI.updateMenu(id, data);
      successMessage.value = '菜单更新成功';
      await fetchMenus();
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || '更新菜单失败';
      return false;
    } finally {
      loading.value = false;
    }
  };

  const deleteMenu = async (id: number): Promise<boolean> => {
    loading.value = true;
    errorMessage.value = '';
    successMessage.value = '';

    try {
      await AdminMenusAPI.deleteMenu(id);
      successMessage.value = '菜单删除成功';
      await fetchMenus();
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || '删除菜单失败';
      return false;
    } finally {
      loading.value = false;
    }
  };

  const reorderMenus = async (data: ReorderMenusRequest): Promise<boolean> => {
    errorMessage.value = '';
    successMessage.value = '';

    try {
      await AdminMenusAPI.reorderMenus(data);
      successMessage.value = '菜单排序已更新';
      return true;
    } catch (error: any) {
      errorMessage.value = error.message || '更新排序失败';
      return false;
    }
  };

  const clearMessages = () => {
    errorMessage.value = '';
    successMessage.value = '';
  };

  return {
    menus,
    loading,
    errorMessage,
    successMessage,
    fetchMenus,
    fetchMenu,
    createMenu,
    updateMenu,
    deleteMenu,
    reorderMenus,
    clearMessages,
  };
}
