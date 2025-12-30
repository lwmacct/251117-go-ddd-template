/**
 * Admin 系统概览 Composable
 */
import { ref } from "vue";
import { overviewApi, extractData } from "@/api";
import type { StatsStatsDTO } from "@models";
import { useSnackbar } from "@/composables";

export function useOverview() {
  const stats = ref<StatsStatsDTO | null>(null);
  const loading = ref(false);

  // 消息提示
  const { error } = useSnackbar();

  const fetchStats = async () => {
    loading.value = true;

    try {
      const response = await overviewApi.apiSystemOverviewStatsGet();
      stats.value = extractData<StatsStatsDTO>(response.data) ?? null;
    } catch (err) {
      error((err as Error).message || "获取统计信息失败");
      console.error("Failed to fetch stats:", err);
    } finally {
      loading.value = false;
    }
  };

  return {
    stats,
    loading,
    fetchStats,
  };
}
