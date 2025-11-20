/**
 * 认证相关 API（基础版，已废弃）
 * @deprecated 使用 AuthAPI 代替，支持验证码和2FA
 */

import { apiClient } from "./client";
import { saveAccessToken, saveRefreshToken, clearAuthTokens } from "@/utils/auth";
import type { BasicLoginRequest, BasicRegisterRequest, AuthResponse } from "@/types/auth";
import type { ApiResponse } from "@/types/response";

/**
 * 用户登录（基础版）
 * @deprecated 使用 AuthAPI.login() 代替
 */
export const login = async (req: BasicLoginRequest): Promise<AuthResponse> => {
  const { data } = await apiClient.post<ApiResponse<AuthResponse>>("/api/auth/login", req);

  if (data.data) {
    // 保存 token
    saveAccessToken(data.data.access_token);
    saveRefreshToken(data.data.refresh_token);
    return data.data;
  }

  throw new Error(data.error || "Login failed");
};

/**
 * 用户注册（基础版）
 * @deprecated 使用 AuthAPI.register() 代替
 */
export const register = async (req: BasicRegisterRequest): Promise<AuthResponse> => {
  const { data } = await apiClient.post<ApiResponse<AuthResponse>>("/api/auth/register", req);

  if (data.data) {
    // 保存 token
    saveAccessToken(data.data.access_token);
    saveRefreshToken(data.data.refresh_token);
    return data.data;
  }

  throw new Error(data.error || "Registration failed");
};

/**
 * 刷新访问令牌
 */
export const refreshToken = async (refreshToken: string): Promise<AuthResponse> => {
  const { data } = await apiClient.post<ApiResponse<AuthResponse>>("/api/auth/refresh", {
    refresh_token: refreshToken,
  });

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "Token refresh failed");
};

/**
 * 用户登出
 */
export const logout = () => {
  clearAuthTokens();
};
