/**
 * Admin 系统概览 Composable
 */
import { ref } from 'vue';
import { AdminOverviewAPI } from '@/api/admin';
import type { SystemStats } from '@/types/admin';

export function useOverview() {
  const stats = ref<SystemStats | null>(null);
  const loading = ref(false);
  const errorMessage = ref('');

  const fetchStats = async () => {
    loading.value = true;
    errorMessage.value = '';

    try {
      stats.value = await AdminOverviewAPI.getSystemStats();
    } catch (error: any) {
      errorMessage.value = error.message || '获取统计信息失败';
      console.error('Failed to fetch stats:', error);
    } finally {
      loading.value = false;
    }
  };

  return {
    stats,
    loading,
    errorMessage,
    fetchStats,
  };
}
