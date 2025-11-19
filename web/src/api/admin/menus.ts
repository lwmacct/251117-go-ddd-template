/**
 * Admin 菜单管理 API
 */
import { apiClient } from '../auth/client';
import type { ApiResponse } from '@/types/auth';
import type { Menu, CreateMenuRequest, UpdateMenuRequest, ReorderMenusRequest } from '@/types/admin';

/**
 * 获取菜单列表（树形结构）
 */
export const listMenus = async (): Promise<Menu[]> => {
  try {
    const { data } = await apiClient.get<ApiResponse<Menu[]>>('/admin/menus');

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || '获取菜单列表失败');
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || '获取菜单列表失败');
  }
};

/**
 * 获取菜单详情
 */
export const getMenu = async (id: number): Promise<Menu> => {
  try {
    const { data } = await apiClient.get<ApiResponse<Menu>>(`/admin/menus/${id}`);

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || '获取菜单详情失败');
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || '获取菜单详情失败');
  }
};

/**
 * 创建菜单
 */
export const createMenu = async (params: CreateMenuRequest): Promise<Menu> => {
  try {
    const { data } = await apiClient.post<ApiResponse<Menu>>('/admin/menus', params);

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || '创建菜单失败');
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || '创建菜单失败');
  }
};

/**
 * 更新菜单
 */
export const updateMenu = async (id: number, params: UpdateMenuRequest): Promise<Menu> => {
  try {
    const { data } = await apiClient.put<ApiResponse<Menu>>(`/admin/menus/${id}`, params);

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || '更新菜单失败');
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || '更新菜单失败');
  }
};

/**
 * 删除菜单
 */
export const deleteMenu = async (id: number): Promise<void> => {
  try {
    await apiClient.delete(`/admin/menus/${id}`);
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || '删除菜单失败');
  }
};

/**
 * 批量更新菜单排序
 */
export const reorderMenus = async (params: ReorderMenusRequest): Promise<void> => {
  try {
    await apiClient.put('/admin/menus/reorder', params);
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || '更新排序失败');
  }
};
