/**
 * Admin 系统设置 API
 */
import { apiClient } from '../auth/client';
import type { ApiResponse } from '@/types/auth';
import type { SettingGroup, UpdateSettingsRequest } from '@/types/admin';

/**
 * 获取所有设置（按分组）
 */
export const getSettings = async (): Promise<SettingGroup[]> => {
  try {
    const { data } = await apiClient.get<ApiResponse<SettingGroup[]>>('/admin/settings');

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || '获取设置失败');
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || '获取设置失败');
  }
};

/**
 * 批量更新设置
 */
export const updateSettings = async (params: UpdateSettingsRequest): Promise<void> => {
  try {
    await apiClient.put('/admin/settings', params);
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || '更新设置失败');
  }
};
