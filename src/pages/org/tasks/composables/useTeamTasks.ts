/**
 * 团队任务 Composable
 *
 * 提供团队任务的完整 CRUD 功能和状态管理
 */
import { ref, computed } from "vue";
import { orgTeamTaskApi, orgTeamMemberApi, extractList, extractData } from "@/api";
import type { TaskTaskDTO, TaskCreateTaskDTO, TaskUpdateTaskDTO, OrgTeamMemberDTO } from "@models";
import { refDebounced } from "@vueuse/core";
import { useServerPagination, useSnackbar, ITEMS_PER_PAGE_OPTIONS } from "@/composables";

/**
 * 任务状态枚举
 */
export const TaskStatus = {
  Pending: "pending",
  InProgress: "in_progress",
  Completed: "completed",
} as const;

export type TaskStatusValue = (typeof TaskStatus)[keyof typeof TaskStatus];

/**
 * 任务状态配置
 */
export const TaskStatusConfig: Record<TaskStatusValue, { label: string; color: string; icon: string }> = {
  [TaskStatus.Pending]: { label: "待处理", color: "grey", icon: "mdi-clock-outline" },
  [TaskStatus.InProgress]: { label: "进行中", color: "info", icon: "mdi-play-circle-outline" },
  [TaskStatus.Completed]: { label: "已完成", color: "success", icon: "mdi-check-circle-outline" },
};

/**
 * 状态流转规则（key → 允许的下一状态）
 */
export const StatusTransitions: Record<TaskStatusValue, TaskStatusValue[]> = {
  [TaskStatus.Pending]: [TaskStatus.InProgress, TaskStatus.Completed],
  [TaskStatus.InProgress]: [TaskStatus.Pending, TaskStatus.Completed],
  [TaskStatus.Completed]: [TaskStatus.InProgress], // 允许重新打开
};

/**
 * 团队任务管理 Composable
 * @param orgId 组织 ID
 * @param teamId 团队 ID
 */
export function useTeamTasks(orgId: number, teamId: number) {
  // ========== 状态 ==========
  const tasks = ref<TaskTaskDTO[]>([]);
  const searchQuery = ref("");
  const statusFilter = ref<TaskStatusValue | "all">("all");
  const debouncedSearchQuery = refDebounced(searchQuery, 300);

  const members = ref<OrgTeamMemberDTO[]>([]);
  const membersLoading = ref(false);

  // 使用通用分页 composable
  const { pagination, loading, onTableOptionsUpdate, updateTotal, resetAndFetch, getParams } = useServerPagination();

  const { success, error } = useSnackbar();

  // ========== 数据获取 ==========

  /**
   * 获取任务列表
   */
  const fetchTasks = async () => {
    loading.value = true;
    try {
      const { limit, page } = getParams();
      const response = await orgTeamTaskApi.apiOrgOrgIdTeamsTeamIdTasksGet(orgId, teamId, limit, page);
      const result = extractList<TaskTaskDTO>(response.data);
      tasks.value = result.data;
      updateTotal(result.pagination.total, result.pagination.total_pages);
    } catch (err) {
      error((err as Error).message || "获取任务列表失败");
      console.error("Failed to fetch tasks:", err);
    } finally {
      loading.value = false;
    }
  };

  /**
   * 获取团队成员列表（用于指派）
   */
  const fetchMembers = async () => {
    membersLoading.value = true;
    try {
      // 获取所有团队成员（设置一个较大的 limit）
      const response = await orgTeamMemberApi.apiOrgOrgIdTeamsTeamIdMembersGet(orgId, teamId, 1000, 1);
      const result = extractList<OrgTeamMemberDTO>(response.data);
      members.value = result.data;
    } catch (err) {
      error((err as Error).message || "获取团队成员列表失败");
      console.error("Failed to fetch members:", err);
    } finally {
      membersLoading.value = false;
    }
  };

  /**
   * 获取单个任务详情
   */
  const fetchTask = async (id: number): Promise<TaskTaskDTO | null> => {
    try {
      const response = await orgTeamTaskApi.apiOrgOrgIdTeamsTeamIdTasksIdGet(orgId, teamId, id);
      return extractData<TaskTaskDTO>(response.data) ?? null;
    } catch (err) {
      error((err as Error).message || "获取任务详情失败");
      console.error("Failed to fetch task:", err);
      return null;
    }
  };

  // ========== CRUD 操作 ==========

  /**
   * 创建任务
   */
  const createTask = async (data: TaskCreateTaskDTO): Promise<boolean> => {
    loading.value = true;
    try {
      await orgTeamTaskApi.apiOrgOrgIdTeamsTeamIdTasksPost(orgId, teamId, data);
      success("任务创建成功");
      await fetchTasks();
      return true;
    } catch (err) {
      error((err as Error).message || "创建任务失败");
      console.error("Failed to create task:", err);
      return false;
    } finally {
      loading.value = false;
    }
  };

  /**
   * 更新任务
   */
  const updateTask = async (id: number, data: TaskUpdateTaskDTO): Promise<boolean> => {
    loading.value = true;
    try {
      await orgTeamTaskApi.apiOrgOrgIdTeamsTeamIdTasksIdPut(orgId, teamId, id, data);
      success("任务更新成功");
      await fetchTasks();
      return true;
    } catch (err) {
      error((err as Error).message || "更新任务失败");
      console.error("Failed to update task:", err);
      return false;
    } finally {
      loading.value = false;
    }
  };

  /**
   * 删除任务（软删除）
   */
  const deleteTask = async (id: number): Promise<boolean> => {
    loading.value = true;
    try {
      await orgTeamTaskApi.apiOrgOrgIdTeamsTeamIdTasksIdDelete(orgId, teamId, id);
      success("任务删除成功");
      await fetchTasks();
      return true;
    } catch (err) {
      error((err as Error).message || "删除任务失败");
      console.error("Failed to delete task:", err);
      return false;
    } finally {
      loading.value = false;
    }
  };

  /**
   * 更新任务状态
   */
  const updateStatus = async (id: number, status: TaskStatusValue): Promise<boolean> => {
    loading.value = true;
    try {
      await orgTeamTaskApi.apiOrgOrgIdTeamsTeamIdTasksIdPut(orgId, teamId, id, {
        status,
      });
      success("状态更新成功");
      await fetchTasks();
      return true;
    } catch (err) {
      error((err as Error).message || "状态更新失败");
      console.error("Failed to update status:", err);
      return false;
    } finally {
      loading.value = false;
    }
  };

  // ========== 辅助函数 ==========

  /**
   * 获取可用的下一状态选项
   */
  const getNextStatuses = (currentStatus: TaskStatusValue): TaskStatusValue[] => {
    return StatusTransitions[currentStatus] || [];
  };

  /**
   * 获取成员选择选项（用于 v-select）
   */
  const memberOptions = computed(() => {
    return members.value
      .filter((m) => m.user_id !== undefined)
      .map((m) => ({
        title: m.full_name || m.username || `User ${m.user_id}`,
        value: m.user_id as number,
      }));
  });

  /**
   * 获取成员显示名称
   */
  const getMemberName = (memberId: number | undefined | null): string => {
    if (!memberId) return "-";
    const member = members.value.find((m) => m.user_id === memberId);
    return member?.full_name || member?.username || `User ${memberId}`;
  };

  /**
   * 筛选后的任务列表
   */
  const filteredTasks = computed(() => {
    let result = tasks.value;

    // 状态筛选
    if (statusFilter.value !== "all") {
      result = result.filter((t) => t.status === statusFilter.value);
    }

    // 搜索筛选
    if (debouncedSearchQuery.value) {
      const query = debouncedSearchQuery.value.toLowerCase();
      result = result.filter((t) => t.title?.toLowerCase().includes(query) || t.description?.toLowerCase().includes(query));
    }

    return result;
  });

  /**
   * 处理表格选项更新（带搜索和状态筛选变化监听）
   */
  const handleTableOptionsUpdate = (options: { page: number; itemsPerPage: number }) => {
    onTableOptionsUpdate(options, fetchTasks);
  };

  return {
    // 状态
    tasks,
    filteredTasks,
    searchQuery,
    statusFilter,
    members,
    memberOptions,
    loading,
    membersLoading,
    pagination,

    // 方法
    fetchTasks,
    fetchMembers,
    fetchTask,
    createTask,
    updateTask,
    deleteTask,
    updateStatus,
    getNextStatuses,
    getMemberName,
    handleTableOptionsUpdate,
    resetAndFetch: () => resetAndFetch(fetchTasks),

    // 常量
    TaskStatus,
    TaskStatusConfig,
    StatusTransitions,
    ITEMS_PER_PAGE_OPTIONS,
  };
}
