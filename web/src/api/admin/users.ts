/**
 * Admin 用户管理 API
 */
import { apiClient } from "../auth/client";
import type { ApiResponse } from "@/types/auth";
import type { AdminUser, CreateUserRequest, UpdateUserRequest, AssignRolesRequest } from "@/types/admin";
import type { PaginatedResponse, PaginationParams } from "@/types/common";

/**
 * 获取用户列表（分页）
 */
export const listUsers = async (params: Partial<PaginationParams>): Promise<PaginatedResponse<AdminUser>> => {
  try {
    const { data } = await apiClient.get<ApiResponse<PaginatedResponse<AdminUser>>>("/admin/users", { params });

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || "获取用户列表失败");
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "获取用户列表失败");
  }
};

/**
 * 获取用户详情
 */
export const getUser = async (id: number): Promise<AdminUser> => {
  try {
    const { data } = await apiClient.get<ApiResponse<AdminUser>>(`/admin/users/${id}`);

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || "获取用户详情失败");
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "获取用户详情失败");
  }
};

/**
 * 创建用户
 */
export const createUser = async (params: CreateUserRequest): Promise<AdminUser> => {
  try {
    const { data } = await apiClient.post<ApiResponse<AdminUser>>("/admin/users", params);

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || "创建用户失败");
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "创建用户失败");
  }
};

/**
 * 更新用户
 */
export const updateUser = async (id: number, params: UpdateUserRequest): Promise<AdminUser> => {
  try {
    const { data } = await apiClient.put<ApiResponse<AdminUser>>(`/admin/users/${id}`, params);

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || "更新用户失败");
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "更新用户失败");
  }
};

/**
 * 删除用户
 */
export const deleteUser = async (id: number): Promise<void> => {
  try {
    await apiClient.delete(`/admin/users/${id}`);
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "删除用户失败");
  }
};

/**
 * 分配角色
 */
export const assignRoles = async (id: number, params: AssignRolesRequest): Promise<AdminUser> => {
  try {
    const { data } = await apiClient.put<ApiResponse<AdminUser>>(`/admin/users/${id}/roles`, params);

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || "分配角色失败");
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "分配角色失败");
  }
};
