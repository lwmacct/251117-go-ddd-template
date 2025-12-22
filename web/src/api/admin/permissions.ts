/**
 * Admin 权限管理 API
 */
import { apiClient } from "../auth/client";
import { normalizeListResponse } from "../helpers/pagination";
import type { ListApiResponse } from "@/types/response";
import type { Permission } from "@/types/admin";
import type { PaginatedResponse, PaginationParams } from "@/types/common";

/**
 * 获取权限列表（分页）
 */
export const listPermissions = async (
  params: Partial<PaginationParams> = {}
): Promise<PaginatedResponse<Permission>> => {
  const page = params.page ?? 1;
  const limit = params.limit ?? 50;

  const { data } = await apiClient.get<ListApiResponse<Permission[]>>("/api/admin/permissions", {
    params: {
      page,
      limit,
      search: params.search,
    },
  });

  return normalizeListResponse<Permission>(data, { page, limit });
};

/**
 * 获取所有权限（不分页，用于选择器）
 */
export const getAllPermissions = async (): Promise<Permission[]> => {
  const response = await listPermissions({ page: 1, limit: 1000 });
  return response.data;
};
