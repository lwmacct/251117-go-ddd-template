/**
 * Admin 系统概览 API
 */
import { apiClient } from '../auth/client';
import type { ApiResponse } from '@/types/auth';
import type { SystemStats } from '@/types/admin';

/**
 * 获取系统统计信息
 */
export const getSystemStats = async (): Promise<SystemStats> => {
  try {
    const { data } = await apiClient.get<ApiResponse<SystemStats>>('/admin/overview/stats');

    if (data.data) {
      return data.data;
    }

    throw new Error(data.error || '获取统计信息失败');
  } catch (error: any) {
    throw new Error(error.response?.data?.error || error.message || '获取统计信息失败');
  }
};
