/**
 * Admin 系统概览 API
 */
import { apiClient } from "../auth/client";
import type { ApiResponse } from "@/types/response";
import type { SystemStats } from "@/types/admin";

/**
 * 获取系统统计信息
 */
export const getSystemStats = async (): Promise<SystemStats> => {
  const { data } = await apiClient.get<ApiResponse<SystemStats>>("/api/admin/overview/stats");

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "获取统计信息失败");
};
