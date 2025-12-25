/**
 * Admin 系统概览 Composable
 */
import { ref } from "vue";
import { overviewApi, extractData } from "@/api";
import type { StatsStatsDTO } from "@models";

export function useOverview() {
  const stats = ref<StatsStatsDTO | null>(null);
  const loading = ref(false);
  const errorMessage = ref("");

  const fetchStats = async () => {
    loading.value = true;
    errorMessage.value = "";

    try {
      const response = await overviewApi.apiAdminOverviewStatsGet();
      stats.value = extractData<StatsStatsDTO>(response.data) ?? null;
    } catch (error) {
      errorMessage.value = (error as Error).message || "获取统计信息失败";
      console.error("Failed to fetch stats:", error);
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
