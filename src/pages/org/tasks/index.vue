<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useTeamTasks, TaskStatus, type TaskStatusValue } from "./composables/useTeamTasks";
import TaskDialog from "./components/TaskDialog.vue";
import TaskStatusChip from "./components/TaskStatusChip.vue";
import CopyButton from "@/components/CopyButton.vue";
import type { TaskTaskDTO, TaskCreateTaskDTO, TaskUpdateTaskDTO } from "@models";

interface Props {
  orgId: number;
  teamId: number;
}

const props = defineProps<Props>();

// 使用 composable
const {
  filteredTasks,
  searchQuery,
  statusFilter,
  memberOptions,
  loading,
  pagination,
  fetchTasks,
  fetchMembers,
  createTask,
  updateTask,
  deleteTask,
  updateStatus,
  getNextStatuses,
  getMemberName,
  handleTableOptionsUpdate,
  TaskStatusConfig,
  ITEMS_PER_PAGE_OPTIONS,
} = useTeamTasks(props.orgId, props.teamId);

// 对话框状态
const taskDialog = ref(false);
const deleteDialog = ref(false);
const statusMenu = ref<{ [key: number]: boolean }>({});

// 编辑状态
const selectedTask = ref<TaskTaskDTO | null>(null);
const taskToDelete = ref<TaskTaskDTO | null>(null);

// 表头配置
const headers = [
  { title: "ID", key: "id", sortable: true, width: 80 },
  { title: "任务标题", key: "title", sortable: true },
  { title: "描述", key: "description", sortable: false },
  { title: "状态", key: "status", sortable: true, width: 140 },
  { title: "指派给", key: "assignee", sortable: false, width: 140 },
  { title: "创建时间", key: "created_at", sortable: true, width: 170 },
  { title: "操作", key: "actions", sortable: false, width: 160, align: "start" as const },
];

// 状态筛选选项
const statusFilterOptions = [
  { title: "全部", value: "all" },
  { title: TaskStatusConfig[TaskStatus.Pending].label, value: TaskStatus.Pending },
  { title: TaskStatusConfig[TaskStatus.InProgress].label, value: TaskStatus.InProgress },
  { title: TaskStatusConfig[TaskStatus.Completed].label, value: TaskStatus.Completed },
];

// 初始化加载
// 注意：fetchTasks() 由 v-data-table-server 的 @update:options 自动触发
onMounted(() => {
  fetchMembers();
});

// 监听筛选条件变化，重置到第一页
watch([statusFilter], () => {
  fetchTasks();
});

// 打开创建对话框
const openCreateDialog = () => {
  selectedTask.value = null;
  taskDialog.value = true;
};

// 打开编辑对话框
const openEditDialog = (task: TaskTaskDTO) => {
  selectedTask.value = task;
  taskDialog.value = true;
};

// 打开删除确认对话框
const openDeleteDialog = (task: TaskTaskDTO) => {
  taskToDelete.value = task;
  deleteDialog.value = true;
};

// 保存任务（创建或编辑）
const handleSaveTask = async (data: TaskCreateTaskDTO | TaskUpdateTaskDTO) => {
  let success = false;

  if (selectedTask.value?.id) {
    success = await updateTask(selectedTask.value.id, data as TaskUpdateTaskDTO);
  } else {
    success = await createTask(data as TaskCreateTaskDTO);
  }

  if (success) {
    taskDialog.value = false;
  }
};

// 确认删除
const confirmDelete = async () => {
  if (!taskToDelete.value?.id) return;

  const success = await deleteTask(taskToDelete.value.id);
  if (success) {
    deleteDialog.value = false;
    taskToDelete.value = null;
  }
};

// 快速更改状态
const handleStatusChange = async (task: TaskTaskDTO, newStatus: TaskStatusValue) => {
  statusMenu.value[task.id!] = false;
  await updateStatus(task.id!, newStatus);
};

// 格式化日期
const formatDate = (dateString?: string) => {
  if (!dateString) return "-";
  return new Date(dateString).toLocaleString("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
};

// 截断描述文本
const truncateDescription = (text?: string, maxLength = 50) => {
  if (!text) return "-";
  if (text.length <= maxLength) return text;
  return text.slice(0, maxLength) + "...";
};
</script>

<template>
  <div class="tasks-page">
    <!-- 标题 -->
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-2">团队任务</h1>
        <p class="text-body-2 text-medium-emphasis mb-0">管理团队任务，跟踪工作进度</p>
      </v-col>
    </v-row>

    <!-- 任务列表卡片 -->
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-row align="center">
              <!-- 搜索和筛选 -->
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="searchQuery"
                  prepend-inner-icon="mdi-magnify"
                  label="搜索任务（标题或描述）"
                  single-line
                  hide-details
                  clearable
                  variant="outlined"
                  density="compact"
                  placeholder="输入后自动搜索..."
                ></v-text-field>
              </v-col>

              <!-- 状态筛选和新建按钮 -->
              <v-col cols="12" md="6" class="d-flex align-center justify-end">
                <v-select
                  v-model="statusFilter"
                  :items="statusFilterOptions"
                  label="状态筛选"
                  variant="outlined"
                  density="compact"
                  hide-details
                  class="mr-4"
                  style="max-width: 150px"
                ></v-select>
                <v-btn color="primary" @click="openCreateDialog">
                  <v-icon start>mdi-plus</v-icon>
                  新建任务
                </v-btn>
              </v-col>
            </v-row>
          </v-card-title>

          <v-card-text>
            <v-data-table-server
              :items-per-page="pagination.limit"
              :items-per-page-options="ITEMS_PER_PAGE_OPTIONS"
              :page="pagination.page"
              :headers="headers"
              :items="filteredTasks"
              :items-length="pagination.total"
              :loading="loading"
              loading-text="加载中..."
              no-data-text="暂无任务数据"
              @update:options="handleTableOptionsUpdate"
            >
              <!-- ID 列 -->
              <template #item.id="{ item }">
                <div class="d-flex align-center">
                  <span class="text-body-2">{{ item.id }}</span>
                  <CopyButton :text="String(item.id)" size="x-small" />
                </div>
              </template>

              <!-- 标题列 -->
              <template #item.title="{ item }">
                <span class="font-weight-medium">{{ item.title || "-" }}</span>
              </template>

              <!-- 描述列 -->
              <template #item.description="{ item }">
                <span class="text-body-2 text-medium-emphasis">{{ truncateDescription(item.description) }}</span>
              </template>

              <!-- 状态列 -->
              <template #item.status="{ item }">
                <TaskStatusChip :status="item.status" />
              </template>

              <!-- 指派给列 -->
              <template #item.assignee="{ item }">
                <div class="d-flex align-center">
                  <v-avatar size="24" color="surface-light" class="mr-2">
                    <v-icon size="16" icon="mdi-account" />
                  </v-avatar>
                  <span class="text-body-2">{{ getMemberName(item.assignee_id) }}</span>
                </div>
              </template>

              <!-- 创建时间列 -->
              <template #item.created_at="{ item }">
                <span class="text-body-2 text-medium-emphasis">{{ formatDate(item.created_at) }}</span>
              </template>

              <!-- 操作列 -->
              <template #item.actions="{ item }">
                <!-- 编辑 -->
                <v-tooltip text="编辑">
                  <template #activator="{ props: tooltipProps }">
                    <v-btn
                      icon="mdi-pencil"
                      size="small"
                      variant="text"
                      v-bind="tooltipProps"
                      @click="openEditDialog(item)"
                    ></v-btn>
                  </template>
                </v-tooltip>

                <!-- 更改状态 -->
                <v-menu v-model="statusMenu[item.id!]">
                  <template #activator="{ props: menuProps }">
                    <v-tooltip text="更改状态">
                      <template #activator="{ props: tooltipProps }">
                        <v-btn icon="mdi-flag" size="small" variant="text" v-bind="{ ...menuProps, ...tooltipProps }"></v-btn>
                      </template>
                    </v-tooltip>
                  </template>
                  <v-list>
                    <v-list-item
                      v-for="status in getNextStatuses(item.status as TaskStatusValue)"
                      :key="status"
                      @click="handleStatusChange(item, status)"
                    >
                      <template #prepend>
                        <v-icon :icon="TaskStatusConfig[status].icon" :color="TaskStatusConfig[status].color" />
                      </template>
                      <v-list-item-title>{{ TaskStatusConfig[status].label }}</v-list-item-title>
                    </v-list-item>
                  </v-list>
                </v-menu>

                <!-- 删除 -->
                <v-tooltip text="删除">
                  <template #activator="{ props: deleteTooltipProps }">
                    <v-btn
                      icon="mdi-delete"
                      size="small"
                      variant="text"
                      color="error"
                      v-bind="deleteTooltipProps"
                      @click="openDeleteDialog(item)"
                    ></v-btn>
                  </template>
                </v-tooltip>
              </template>
            </v-data-table-server>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 创建/编辑对话框 -->
    <TaskDialog v-model="taskDialog" :task="selectedTask" :members="memberOptions" :loading="loading" @save="handleSaveTask" />

    <!-- 删除确认对话框 -->
    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h5">确认删除</v-card-title>
        <v-card-text>
          确定要删除任务 <strong>{{ taskToDelete?.title }}</strong> 吗？此操作不可恢复。
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="deleteDialog = false">取消</v-btn>
          <v-btn color="error" variant="elevated" :loading="loading" @click="confirmDelete"> 删除 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.tasks-page {
  width: 100%;
}
</style>
