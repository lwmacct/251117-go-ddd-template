/**
 * Personal Access Token API
 */
import { apiClient } from "../auth/client";
import type { ApiResponse } from "@/types/response";
import type { PersonalAccessToken, CreateTokenRequest, CreateTokenResponse } from "@/types/user";

/**
 * 获取 Token 列表
 */
export const listTokens = async (): Promise<PersonalAccessToken[]> => {
  try {
    const { data } = await apiClient.get<ApiResponse<PersonalAccessToken[]>>("/api/user/tokens");

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || "获取 Token 列表失败");
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "获取 Token 列表失败");
  }
};

/**
 * 获取 Token 详情
 */
export const getToken = async (id: number): Promise<PersonalAccessToken> => {
  try {
    const { data } = await apiClient.get<ApiResponse<PersonalAccessToken>>(`/api/user/tokens/${id}`);

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || "获取 Token 详情失败");
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "获取 Token 详情失败");
  }
};

/**
 * 创建 Token
 */
export const createToken = async (params: CreateTokenRequest): Promise<CreateTokenResponse> => {
  try {
    const { data } = await apiClient.post<ApiResponse<CreateTokenResponse>>("/api/user/tokens", params);

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || "创建 Token 失败");
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "创建 Token 失败");
  }
};

/**
 * 删除 Token
 */
export const deleteToken = async (id: number): Promise<void> => {
  try {
    await apiClient.delete(`/api/user/tokens/${id}`);
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "删除 Token 失败");
  }
};

/**
 * 禁用 Token
 */
export const disableToken = async (id: number): Promise<void> => {
  try {
    await apiClient.patch(`/api/user/tokens/${id}/disable`);
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "禁用 Token 失败");
  }
};

/**
 * 启用 Token
 */
export const enableToken = async (id: number): Promise<void> => {
  try {
    await apiClient.patch(`/api/user/tokens/${id}/enable`);
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || "启用 Token 失败");
  }
};
