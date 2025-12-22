/**
 * API 客户端配置
 */
import axios, { type AxiosError } from "axios";
import { getAccessToken, getRefreshToken, saveAccessToken, saveRefreshToken, clearAuthTokens } from "@/utils/auth";
import type { AuthResponse } from "@/types/auth";
import type { ApiResponse, ErrorResponse } from "@/types/response";
import { extractErrorFromAxios } from "../errors";

/** axios 实例 */
export const apiClient = axios.create({
  timeout: 10000,
});

// 请求拦截器 - 添加 token
apiClient.interceptors.request.use(
  (config) => {
    const token = getAccessToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器 - 处理错误
apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<ErrorResponse>) => {
    const originalRequest = error.config as typeof error.config & { _retry?: boolean };

    // 如果是 401 错误且不是刷新 token 请求，尝试刷新 token
    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = getRefreshToken();
        if (refreshToken) {
          const { data } = await axios.post<ApiResponse<AuthResponse>>("/api/auth/refresh", {
            refresh_token: refreshToken,
          });

          if (data.data) {
            // 保存新 token
            saveAccessToken(data.data.access_token);
            saveRefreshToken(data.data.refresh_token);

            // 重试原请求
            originalRequest.headers.Authorization = `Bearer ${data.data.access_token}`;
            return apiClient(originalRequest);
          }
        }
      } catch {
        // 刷新失败，清除 token 并跳转到登录页
        clearAuthTokens();
        window.location.href = "/#/auth/login";
        return Promise.reject(extractErrorFromAxios(error));
      }
    }

    // 统一错误转换
    return Promise.reject(extractErrorFromAxios(error));
  }
);
