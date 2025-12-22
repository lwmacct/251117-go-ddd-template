/**
 * Admin 系统设置 API
 */
import { apiClient } from "../auth/client";
import type { ApiResponse } from "@/types/response";
import type { Setting, CreateSettingRequest, UpdateSettingRequest, UpdateSettingsRequest } from "@/types/admin";

/**
 * 获取所有设置（按分组）
 */
export const getSettings = async (category?: string): Promise<Setting[]> => {
  const { data } = await apiClient.get<ApiResponse<Setting[]>>("/api/admin/settings", {
    params: category ? { category } : undefined,
  });

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "获取设置失败");
};

/**
 * 按分类获取设置
 */
export const getSettingsByCategory = async (category: string): Promise<Setting[]> => {
  return getSettings(category);
};

/**
 * 获取单个设置
 */
export const getSetting = async (key: string): Promise<Setting | null> => {
  const { data } = await apiClient.get<ApiResponse<Setting>>(`/api/admin/settings/${key}`);
  return data.data ?? null;
};

/**
 * 创建设置
 */
export const createSetting = async (payload: CreateSettingRequest): Promise<Setting> => {
  const { data } = await apiClient.post<ApiResponse<Setting>>("/api/admin/settings", payload);

  if (data.data) {
    return data.data;
  }

  throw new Error(data.message || "创建设置失败");
};

/**
 * 更新设置
 */
export const updateSetting = async (key: string, payload: UpdateSettingRequest): Promise<Setting> => {
  const { data } = await apiClient.put<ApiResponse<Setting>>(`/api/admin/settings/${key}`, payload);

  if (data.data) {
    return data.data;
  }

  throw new Error(data.message || "更新设置失败");
};

/**
 * 删除设置
 */
export const deleteSetting = async (key: string): Promise<void> => {
  await apiClient.delete(`/api/admin/settings/${key}`);
};

/**
 * 批量更新设置
 */
export const updateSettings = async (params: UpdateSettingsRequest): Promise<void> => {
  await apiClient.put("/api/admin/settings/batch", params);
};

/**
 * 批量更新设置（别名，向后兼容）
 */
export const batchUpdateSettings = async (params: UpdateSettingsRequest): Promise<void> => {
  await updateSettings(params);
};
