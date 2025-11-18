/**
 * 用户相关 API
 */
import { apiClient } from "./client";
import type { User, ApiResponse } from "@/types/auth";

/**
 * 获取当前用户信息
 */
export const getCurrentUser = async (): Promise<User> => {
  const { data } = await apiClient.get<ApiResponse<User>>("/me");

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "Failed to get user info");
};
