/**
 * 认证相关 API
 */

import { apiClient } from "./client";
import { saveAccessToken, saveRefreshToken, clearAuthTokens } from "@/utils/auth";
import type { LoginRequest, RegisterRequest, AuthResponse, ApiResponse } from "@/types/auth";

/**
 * 用户登录
 */
export const login = async (req: LoginRequest): Promise<AuthResponse> => {
  const { data } = await apiClient.post<ApiResponse<AuthResponse>>("/login", req);

  if (data.data) {
    // 保存 token
    saveAccessToken(data.data.access_token);
    saveRefreshToken(data.data.refresh_token);
    return data.data;
  }

  throw new Error(data.error || "Login failed");
};

/**
 * 用户注册
 */
export const register = async (req: RegisterRequest): Promise<AuthResponse> => {
  const { data } = await apiClient.post<ApiResponse<AuthResponse>>("/register", req);

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
  const { data } = await apiClient.post<ApiResponse<AuthResponse>>("/refresh", {
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
