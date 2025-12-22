/**
 * 用户相关 API
 */
import { apiClient } from "./client";
import type { User } from "@/types/auth";
import type { ApiResponse } from "@/types/response";

/**
 * 获取当前用户信息
 */
export const getCurrentUser = async (): Promise<User> => {
  const { data } = await apiClient.get<ApiResponse<User>>("/api/auth/me");

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
  await apiClient.put<ApiResponse<null>>("/api/user/me/password", params);
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
  const { data } = await apiClient.put<ApiResponse<User>>("/api/user/me", params);

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "更新个人资料失败");
};

/**
 * 删除账户请求参数
 */
export interface DeleteAccountRequest {
  password: string;
}

/**
 * 删除账户
 * 需要用户输入密码确认身份
 */
export const deleteAccount = async (params: DeleteAccountRequest): Promise<void> => {
  await apiClient.delete<ApiResponse<null>>("/api/user/account", { data: params });
};
