/**
 * Admin 用户管理 API
 */
import { apiClient } from "../auth/client";
import { normalizeListResponse } from "../helpers/pagination";
import type { ApiResponse, ListApiResponse } from "@/types/response";
import type { AdminUser, CreateUserRequest, UpdateUserRequest, AssignRolesRequest } from "@/types/admin";
import type { PaginatedResponse, PaginationParams } from "@/types/common";

/**
 * 获取用户列表（分页）
 */
export const listUsers = async (params: Partial<PaginationParams> = {}): Promise<PaginatedResponse<AdminUser>> => {
  const page = params.page ?? 1;
  const limit = params.limit ?? 20;

  const { data } = await apiClient.get<ListApiResponse<AdminUser[]>>("/api/admin/users", {
    params: {
      page,
      limit,
      search: params.search,
    },
  });

  return normalizeListResponse<AdminUser>(data, { page, limit });
};

/**
 * 获取用户详情
 */
export const getUser = async (id: number): Promise<AdminUser> => {
  const { data } = await apiClient.get<ApiResponse<AdminUser>>(`/admin/users/${id}`);

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "获取用户详情失败");
};

/**
 * 创建用户
 */
export const createUser = async (params: CreateUserRequest): Promise<AdminUser> => {
  const { data } = await apiClient.post<ApiResponse<AdminUser>>("/api/admin/users", params);

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "创建用户失败");
};

/**
 * 更新用户
 */
export const updateUser = async (id: number, params: UpdateUserRequest): Promise<AdminUser> => {
  const { data } = await apiClient.put<ApiResponse<AdminUser>>(`/admin/users/${id}`, params);

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "更新用户失败");
};

/**
 * 删除用户
 */
export const deleteUser = async (id: number): Promise<void> => {
  await apiClient.delete(`/admin/users/${id}`);
};

/**
 * 分配角色
 */
export const assignRoles = async (id: number, params: AssignRolesRequest): Promise<AdminUser> => {
  const { data } = await apiClient.put<ApiResponse<AdminUser>>(`/admin/users/${id}/roles`, params);

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "分配角色失败");
};

/**
 * 批量创建用户请求
 */
export interface BatchCreateUserRequest {
  username: string;
  email: string;
  password: string;
  full_name?: string;
  status?: string;
}

/**
 * 批量创建用户响应
 */
export interface BatchCreateResult {
  success: number;
  failed: number;
  errors: Array<{
    index: number;
    username: string;
    error: string;
  }>;
}

/**
 * 批量创建用户
 * TODO: 后端 API 未实现，需要实现 POST /api/admin/users/batch
 */
export const batchCreateUsers = async (users: BatchCreateUserRequest[]): Promise<BatchCreateResult> => {
  const { data } = await apiClient.post<ApiResponse<BatchCreateResult>>("/api/admin/users/batch", { users });

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "批量创建用户失败");
};
