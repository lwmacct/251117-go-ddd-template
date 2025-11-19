/**
 * Admin 权限管理 API
 */
import { apiClient } from "../auth/client";
import type { ApiResponse } from "@/types/auth";
import type { Permission } from "@/types/admin";
import type { PaginatedResponse, PaginationParams } from "@/types/common";

/**
 * 获取权限列表（分页）
 */
export const listPermissions = async (params: Partial<PaginationParams> = {}): Promise<PaginatedResponse<Permission>> => {
  try {
    const { data } = await apiClient.get<ApiResponse<PaginatedResponse<Permission>>>("/admin/permissions", { params });

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || "获取权限列表失败");
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "获取权限列表失败");
  }
};

/**
 * 获取所有权限（不分页，用于选择器）
 */
export const getAllPermissions = async (): Promise<Permission[]> => {
  try {
    const { data } = await apiClient.get<ApiResponse<PaginatedResponse<Permission>>>("/admin/permissions", {
      params: { page: 1, limit: 1000 },
    });

    if (data.data) {
      return data.data.data;
    }

    throw new Error(data.error || "获取权限列表失败");
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "获取权限列表失败");
  }
};
