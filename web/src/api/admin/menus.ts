/**
 * Admin 菜单管理 API
 */
import { apiClient } from "../auth/client";
import type { ApiResponse } from "@/types/response";
import type { Menu, CreateMenuRequest, UpdateMenuRequest, ReorderMenusRequest } from "@/types/admin";

/**
 * 获取菜单列表（树形结构）
 */
export const listMenus = async (): Promise<Menu[]> => {
  const { data } = await apiClient.get<ApiResponse<Menu[]>>("/api/admin/menus");

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "获取菜单列表失败");
};

/**
 * 获取菜单详情
 */
export const getMenu = async (id: number): Promise<Menu> => {
  const { data } = await apiClient.get<ApiResponse<Menu>>(`/admin/menus/${id}`);

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "获取菜单详情失败");
};

/**
 * 创建菜单
 */
export const createMenu = async (params: CreateMenuRequest): Promise<Menu> => {
  const { data } = await apiClient.post<ApiResponse<Menu>>("/api/admin/menus", params);

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "创建菜单失败");
};

/**
 * 更新菜单
 */
export const updateMenu = async (id: number, params: UpdateMenuRequest): Promise<Menu> => {
  const { data } = await apiClient.put<ApiResponse<Menu>>(`/admin/menus/${id}`, params);

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "更新菜单失败");
};

/**
 * 删除菜单
 */
export const deleteMenu = async (id: number): Promise<void> => {
  await apiClient.delete(`/admin/menus/${id}`);
};

/**
 * 批量更新菜单排序
 */
export const reorderMenus = async (params: ReorderMenusRequest): Promise<void> => {
  await apiClient.put("/api/admin/menus/reorder", params);
};
