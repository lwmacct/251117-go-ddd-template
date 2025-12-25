/**
 * API 客户端配置
 *
 * 本文件提供：
 * 1. axios 实例（带 token 刷新逻辑）
 * 2. 所有 OpenAPI 生成的 API 类实例
 */
import axios, { type AxiosError } from "axios";
import { accessToken, refreshToken, clearAuthTokens } from "@/utils/auth";
import type { AuthLoginResponseDTO } from "@models";
import type { ApiResponse, ErrorResponse } from "../types";
import { extractErrorFromAxios } from "../errors";
import { Configuration } from "@generated";
import {
  AdminAuditLogApi,
  AdminMenuManagementApi,
  AdminRoleManagementApi,
  AdminSettingsApi,
  AdminUserManagementApi,
  AuthenticationApi,
  Authentication2FAApi,
  OverviewApi,
  SystemApi,
  UserPersonalAccessTokenApi,
  UserProfileApi,
} from "@generated/api";

/** axios 实例 */
export const apiClient = axios.create({
  timeout: 10000,
});

// 请求拦截器 - 添加 token
apiClient.interceptors.request.use(
  (config) => {
    const token = accessToken.value;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
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
        const currentRefreshToken = refreshToken.value;
        if (currentRefreshToken) {
          const { data } = await axios.post<ApiResponse<AuthLoginResponseDTO>>("/api/auth/refresh", {
            refresh_token: currentRefreshToken,
          });

          if (data.data) {
            // 保存新 token
            accessToken.value = data.data.access_token;
            refreshToken.value = data.data.refresh_token;

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
  },
);

// ============== API 实例 ==============
// 创建 API 配置，使用带有 token 刷新功能的 axios 实例
const configuration = new Configuration({
  basePath: "", // 使用 apiClient 已配置的 baseURL
});

// 导出配置好的 API 实例
export const adminAuditLogApi = new AdminAuditLogApi(configuration, "", apiClient);
export const adminMenuApi = new AdminMenuManagementApi(configuration, "", apiClient);
export const adminRoleApi = new AdminRoleManagementApi(configuration, "", apiClient);
export const adminSettingsApi = new AdminSettingsApi(configuration, "", apiClient);
export const adminUserApi = new AdminUserManagementApi(configuration, "", apiClient);
export const authApi = new AuthenticationApi(configuration, "", apiClient);
export const auth2faApi = new Authentication2FAApi(configuration, "", apiClient);
export const overviewApi = new OverviewApi(configuration, "", apiClient);
export const systemApi = new SystemApi(configuration, "", apiClient);
export const userTokensApi = new UserPersonalAccessTokenApi(configuration, "", apiClient);
export const userProfileApi = new UserProfileApi(configuration, "", apiClient);
