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

/**
 * 修改密码请求参数
 */
export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

/**
 * 修改密码
 */
export const changePassword = async (params: ChangePasswordRequest): Promise<void> => {
  try {
    await apiClient.put<ApiResponse<null>>("/user/me/password", params);
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "修改密码失败");
  }
};

/**
 * 更新个人资料请求参数
 */
export interface UpdateProfileRequest {
  full_name?: string;
  avatar?: string;
  bio?: string;
}

/**
 * 更新个人资料
 */
export const updateProfile = async (params: UpdateProfileRequest): Promise<User> => {
  try {
    const { data } = await apiClient.put<ApiResponse<User>>("/user/me", params);

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || "更新个人资料失败");
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "更新个人资料失败");
  }
};
