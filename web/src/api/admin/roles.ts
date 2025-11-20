/**
 * Admin 角色管理 API
 */
import { apiClient } from "../auth/client";
import { normalizeListResponse } from "../helpers/pagination";
import type { ApiResponse, ListApiResponse } from "@/types/response";
import type { Role, CreateRoleRequest, UpdateRoleRequest, SetPermissionsRequest } from "@/types/admin";
import type { PaginatedResponse, PaginationParams } from "@/types/common";

/**
 * 获取角色列表（分页）
 */
export const listRoles = async (params: Partial<PaginationParams> = {}): Promise<PaginatedResponse<Role>> => {
  const page = params.page ?? 1;
  const limit = params.limit ?? 20;

  try {
    const { data } = await apiClient.get<ListApiResponse<Role[]>>("/api/admin/roles", {
      params: {
        page,
        limit,
        search: params.search,
      },
    });

    return normalizeListResponse<Role>(data, { page, limit });
  } catch (error: any) {
    const serverError = error.response?.data?.error || error.response?.data?.message;
    throw new Error(serverError || error.message || "获取角色列表失败");
  }
};

/**
 * 获取角色详情
 */
export const getRole = async (id: number): Promise<Role> => {
  try {
    const { data } = await apiClient.get<ApiResponse<Role>>(`/admin/roles/${id}`);

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || "获取角色详情失败");
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "获取角色详情失败");
  }
};

/**
 * 创建角色
 */
export const createRole = async (params: CreateRoleRequest): Promise<Role> => {
  try {
    const { data } = await apiClient.post<ApiResponse<Role>>("/api/admin/roles", params);

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || "创建角色失败");
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "创建角色失败");
  }
};

/**
 * 更新角色
 */
export const updateRole = async (id: number, params: UpdateRoleRequest): Promise<Role> => {
  try {
    const { data } = await apiClient.put<ApiResponse<Role>>(`/admin/roles/${id}`, params);

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || "更新角色失败");
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "更新角色失败");
  }
};

/**
 * 删除角色
 */
export const deleteRole = async (id: number): Promise<void> => {
  try {
    await apiClient.delete(`/admin/roles/${id}`);
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "删除角色失败");
  }
};

/**
 * 设置角色权限
 */
export const setPermissions = async (id: number, params: SetPermissionsRequest): Promise<void> => {
  try {
    await apiClient.put(`/admin/roles/${id}/permissions`, params);
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "设置权限失败");
  }
};
