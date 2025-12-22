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

  const { data } = await apiClient.get<ListApiResponse<Role[]>>("/api/admin/roles", {
    params: {
      page,
      limit,
      search: params.search,
    },
  });

  return normalizeListResponse<Role>(data, { page, limit });
};

/**
 * 获取角色详情
 */
export const getRole = async (id: number): Promise<Role> => {
  const { data } = await apiClient.get<ApiResponse<Role>>(`/admin/roles/${id}`);

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "获取角色详情失败");
};

/**
 * 创建角色
 */
export const createRole = async (params: CreateRoleRequest): Promise<Role> => {
  const { data } = await apiClient.post<ApiResponse<Role>>("/api/admin/roles", params);

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "创建角色失败");
};

/**
 * 更新角色
 */
export const updateRole = async (id: number, params: UpdateRoleRequest): Promise<Role> => {
  const { data } = await apiClient.put<ApiResponse<Role>>(`/admin/roles/${id}`, params);

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "更新角色失败");
};

/**
 * 删除角色
 */
export const deleteRole = async (id: number): Promise<void> => {
  await apiClient.delete(`/admin/roles/${id}`);
};

/**
 * 设置角色权限
 */
export const setPermissions = async (id: number, params: SetPermissionsRequest): Promise<void> => {
  await apiClient.put(`/admin/roles/${id}/permissions`, params);
};
