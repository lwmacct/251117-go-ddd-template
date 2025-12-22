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
  const { data } = await apiClient.get<ApiResponse<PersonalAccessToken[]>>("/api/user/tokens");

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "获取 Token 列表失败");
};

/**
 * 获取 Token 详情
 */
export const getToken = async (id: number): Promise<PersonalAccessToken> => {
  const { data } = await apiClient.get<ApiResponse<PersonalAccessToken>>(`/api/user/tokens/${id}`);

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "获取 Token 详情失败");
};

/**
 * 创建 Token
 */
export const createToken = async (params: CreateTokenRequest): Promise<CreateTokenResponse> => {
  const { data } = await apiClient.post<ApiResponse<CreateTokenResponse>>("/api/user/tokens", params);

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "创建 Token 失败");
};

/**
 * 删除 Token
 */
export const deleteToken = async (id: number): Promise<void> => {
  await apiClient.delete(`/api/user/tokens/${id}`);
};

/**
 * 禁用 Token
 */
export const disableToken = async (id: number): Promise<void> => {
  await apiClient.patch(`/api/user/tokens/${id}/disable`);
};

/**
 * 启用 Token
 */
export const enableToken = async (id: number): Promise<void> => {
  await apiClient.patch(`/api/user/tokens/${id}/enable`);
};
