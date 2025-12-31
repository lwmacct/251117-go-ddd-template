<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { orgTeamApi, orgTeamMemberApi, extractList } from "@/api";
import type { OrgTeamDTO, OrgCreateTeamDTO, OrgUpdateTeamDTO, OrgTeamMemberDTO, OrgAddTeamMemberDTO } from "@models";
import { useServerPagination, useSnackbar } from "@/composables";
import { ITEMS_PER_PAGE_OPTIONS } from "@/composables";
import type { AppError } from "@/api/errors";
import TeamDialog from "./TeamDialog.vue";

interface Props {
  orgId: number;
}

const props = defineProps<Props>();
const { success, error, warning } = useSnackbar();

const teams = ref<OrgTeamDTO[]>([]);
const teamMembers = ref<OrgTeamMemberDTO[]>([]);
const loading = ref(false);
const teamDialog = ref(false);
const deleteDialog = ref(false);
const addMemberDialog = ref(false);

const selectedTeam = ref<OrgTeamDTO | null>(null);
const teamToDelete = ref<OrgTeamDTO | null>(null);
const activeTeamId = ref<number | null>(null);

const dialogMode = ref<"create" | "edit">("create");

// 分页
const { pagination, updateTotal, onTableOptionsUpdate: baseOnTableOptionsUpdate, getParams } = useServerPagination();

// 新增成员表单
const newMember = ref<OrgAddTeamMemberDTO>({
  user_id: 0,
  role: "member",
});

// 角色选项
const roleOptions = [
  { title: "负责人", value: "lead" },
  { title: "成员", value: "member" },
];

// 表头配置 - 团队
const teamHeaders = [
  { title: "ID", key: "id" },
  { title: "团队标识", key: "name" },
  { title: "团队名称", key: "display_name" },
  { title: "描述", key: "description" },
  { title: "操作", key: "actions" },
];

// 表头配置 - 团队成员
const memberHeaders = [
  { title: "ID", key: "id" },
  { title: "用户名", key: "username" },
  { title: "邮箱", key: "email" },
  { title: "角色", key: "role" },
  { title: "加入时间", key: "joined_at" },
  { title: "操作", key: "actions" },
];

/**
 * 获取团队列表
 */
const fetchTeams = async () => {
  if (!props.orgId) return;

  loading.value = true;
  try {
    const { limit, page } = getParams();
    const response = await orgTeamApi.apiOrgOrgIdTeamsGet(props.orgId, limit, page);
    const result = extractList<OrgTeamDTO>(response.data);
    teams.value = result.data;
    updateTotal(result.pagination.total, result.pagination.total_pages);
  } catch (err) {
    const appErr = err as AppError;
    if (appErr.status === 400) {
      warning(appErr.message || "获取团队列表失败");
    } else {
      error(appErr.message || "获取团队列表失败");
    }
  } finally {
    loading.value = false;
  }
};

/**
 * 获取团队成员列表
 */
const fetchTeamMembers = async (teamId: number) => {
  if (!props.orgId) return;

  loading.value = true;
  try {
    const response = await orgTeamMemberApi.apiOrgOrgIdTeamsTeamIdMembersGet(props.orgId, teamId, 100, 1);
    const result = extractList<OrgTeamMemberDTO>(response.data);
    teamMembers.value = result.data;
  } catch (err) {
    const appErr = err as AppError;
    if (appErr.status === 400) {
      warning(appErr.message || "获取团队成员列表失败");
    } else {
      error(appErr.message || "获取团队成员列表失败");
    }
  } finally {
    loading.value = false;
  }
};

/**
 * 创建团队
 */
const createTeam = async (data: OrgCreateTeamDTO) => {
  if (!props.orgId) return;

  loading.value = true;
  try {
    await orgTeamApi.apiOrgOrgIdTeamsPost(props.orgId, data);
    success("团队创建成功");
    teamDialog.value = false;
    await fetchTeams();
  } catch (err) {
    const appErr = err as AppError;
    if (appErr.status === 400) {
      warning(appErr.message || "创建团队失败");
    } else {
      error(appErr.message || "创建团队失败");
    }
  } finally {
    loading.value = false;
  }
};

/**
 * 更新团队
 */
const updateTeam = async (data: OrgUpdateTeamDTO) => {
  if (!props.orgId || !selectedTeam.value?.id) return;

  loading.value = true;
  try {
    await orgTeamApi.apiOrgOrgIdTeamsTeamIdPut(props.orgId, selectedTeam.value.id, data);
    success("团队更新成功");
    teamDialog.value = false;
    await fetchTeams();
  } catch (err) {
    const appErr = err as AppError;
    if (appErr.status === 400) {
      warning(appErr.message || "更新团队失败");
    } else {
      error(appErr.message || "更新团队失败");
    }
  } finally {
    loading.value = false;
  }
};

/**
 * 打开创建团队对话框
 */
const openCreateDialog = () => {
  dialogMode.value = "create";
  selectedTeam.value = null;
  teamDialog.value = true;
};

/**
 * 打开编辑团队对话框
 */
const openEditDialog = (team: OrgTeamDTO) => {
  dialogMode.value = "edit";
  selectedTeam.value = team;
  teamDialog.value = true;
};

/**
 * 删除团队
 */
const deleteTeam = (team: OrgTeamDTO) => {
  teamToDelete.value = team;
  deleteDialog.value = true;
};

/**
 * 确认删除团队
 */
const confirmDelete = async () => {
  if (!props.orgId || !teamToDelete.value?.id) return;

  const deletedTeamId = teamToDelete.value.id;

  loading.value = true;
  try {
    await orgTeamApi.apiOrgOrgIdTeamsTeamIdDelete(props.orgId, deletedTeamId);
    success("团队删除成功");
    deleteDialog.value = false;
    teamToDelete.value = null;
    // 如果删除的是当前选中的团队，清空成员列表
    if (activeTeamId.value === deletedTeamId) {
      activeTeamId.value = null;
      teamMembers.value = [];
    }
    await fetchTeams();
  } catch (err) {
    const appErr = err as AppError;
    if (appErr.status === 400) {
      warning(appErr.message || "删除团队失败");
    } else {
      error(appErr.message || "删除团队失败");
    }
  } finally {
    loading.value = false;
  }
};

/**
 * 选择团队查看成员
 */
const selectTeam = (team: OrgTeamDTO) => {
  if (!team.id) return;
  activeTeamId.value = team.id;
  fetchTeamMembers(team.id);
};

/**
 * 打开添加成员对话框
 */
const openAddMemberDialog = () => {
  newMember.value = { user_id: 0, role: "member" };
  addMemberDialog.value = true;
};

/**
 * 添加团队成员
 */
const addTeamMember = async () => {
  if (!props.orgId || !activeTeamId.value || !newMember.value.user_id) return;

  loading.value = true;
  try {
    await orgTeamMemberApi.apiOrgOrgIdTeamsTeamIdMembersPost(props.orgId, activeTeamId.value, newMember.value);
    success("成员添加成功");
    addMemberDialog.value = false;
    await fetchTeamMembers(activeTeamId.value);
  } catch (err) {
    const appErr = err as AppError;
    if (appErr.status === 400) {
      warning(appErr.message || "添加成员失败");
    } else {
      error(appErr.message || "添加成员失败");
    }
  } finally {
    loading.value = false;
  }
};

/**
 * 移除团队成员
 */
const removeTeamMember = async (member: OrgTeamMemberDTO) => {
  if (!props.orgId || !activeTeamId.value || !member.user_id) return;

  loading.value = true;
  try {
    await orgTeamMemberApi.apiOrgOrgIdTeamsTeamIdMembersUserIdDelete(props.orgId, activeTeamId.value, member.user_id);
    success("成员移除成功");
    await fetchTeamMembers(activeTeamId.value);
  } catch (err) {
    const appErr = err as AppError;
    if (appErr.status === 400) {
      warning(appErr.message || "移除成员失败");
    } else {
      error(appErr.message || "移除成员失败");
    }
  } finally {
    loading.value = false;
  }
};

/**
 * 表格选项变化处理
 */
const onTableOptionsUpdate = (options: { page: number; itemsPerPage: number }) => {
  baseOnTableOptionsUpdate(options, fetchTeams);
};

/**
 * 处理团队对话框保存
 */
const handleTeamDialogSave = (data: OrgCreateTeamDTO | OrgUpdateTeamDTO) => {
  if (dialogMode.value === "create") {
    createTeam(data as OrgCreateTeamDTO);
  } else {
    updateTeam(data as OrgUpdateTeamDTO);
  }
};

// 格式化日期
const formatDate = (dateString?: string) => {
  if (!dateString) return "-";
  return new Date(dateString).toLocaleString("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  });
};

// 角色颜色
const getRoleColor = (role?: string) => {
  const colors: Record<string, string> = {
    lead: "warning",
    member: "default",
  };
  return colors[role ?? ""] || "default";
};

// 角色文本
const getRoleText = (role?: string) => {
  const texts: Record<string, string> = {
    lead: "负责人",
    member: "成员",
  };
  return texts[role ?? ""] || (role ?? "-");
};

// 监听 orgId 变化
watch(
  () => props.orgId,
  () => {
    activeTeamId.value = null;
    teamMembers.value = [];
    fetchTeams();
  },
);

onMounted(fetchTeams);
</script>

<template>
  <div class="org-team-panel">
    <v-row>
      <!-- 团队列表 -->
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title class="d-flex align-center">
            <span>团队列表</span>
            <v-spacer />
            <v-btn color="primary" size="small" @click="openCreateDialog">
              <v-icon start>mdi-plus</v-icon>
              新建团队
            </v-btn>
          </v-card-title>

          <v-card-text>
            <v-data-table-server
              :items-per-page="pagination.limit"
              :items-per-page-options="ITEMS_PER_PAGE_OPTIONS"
              :page="pagination.page"
              :headers="teamHeaders"
              :items="teams"
              :items-length="pagination.total"
              :loading="loading"
              loading-text="加载中..."
              no-data-text="暂无团队"
              density="compact"
              @update:options="onTableOptionsUpdate"
            >
              <!-- ID 列 -->
              <template #item.id="{ item }">
                <span class="text-body-2">{{ item.id }}</span>
              </template>

              <!-- 描述列 -->
              <template #item.description="{ item }">
                <span class="text-body-2 text-truncate d-inline-block" style="max-width: 100px">
                  {{ item.description || "-" }}
                </span>
              </template>

              <!-- 操作列 -->
              <template #item.actions="{ item }">
                <v-btn
                  icon="mdi-account-group"
                  size="x-small"
                  variant="text"
                  :color="activeTeamId === item.id ? 'primary' : 'default'"
                  @click="selectTeam(item)"
                >
                  <v-tooltip activator="parent" text="查看成员" />
                </v-btn>

                <v-btn icon="mdi-pencil" size="x-small" variant="text" @click="openEditDialog(item)">
                  <v-tooltip activator="parent" text="编辑" />
                </v-btn>

                <v-btn icon="mdi-delete" size="x-small" variant="text" color="error" @click="deleteTeam(item)">
                  <v-tooltip activator="parent" text="删除" />
                </v-btn>
              </template>
            </v-data-table-server>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- 团队成员 -->
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title class="d-flex align-center">
            <span>团队成员</span>
            <v-spacer />
            <v-btn v-if="activeTeamId" color="primary" size="small" @click="openAddMemberDialog">
              <v-icon start>mdi-plus</v-icon>
              添加成员
            </v-btn>
          </v-card-title>

          <v-card-text>
            <v-data-table-server
              :headers="memberHeaders"
              :items="teamMembers"
              :items-length="teamMembers.length"
              :loading="loading"
              loading-text="加载中..."
              :no-data-text="activeTeamId ? '暂无成员' : '请先选择团队'"
              hide-default-footer
              density="compact"
            >
              <!-- ID 列 -->
              <template #item.id="{ item }">
                <span class="text-body-2">{{ item.id }}</span>
              </template>

              <!-- 角色列 -->
              <template #item.role="{ item }">
                <v-chip :color="getRoleColor(item.role)" size="small">
                  {{ getRoleText(item.role) }}
                </v-chip>
              </template>

              <!-- 加入时间列 -->
              <template #item.joined_at="{ item }">
                <span class="text-body-2">{{ formatDate(item.joined_at) }}</span>
              </template>

              <!-- 操作列 -->
              <template #item.actions="{ item }">
                <v-btn icon="mdi-minus" size="x-small" variant="text" color="error" @click="removeTeamMember(item)">
                  <v-tooltip activator="parent" text="移除" />
                </v-btn>
              </template>
            </v-data-table-server>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 创建/编辑团队对话框 -->
    <TeamDialog v-model="teamDialog" :team="selectedTeam" :mode="dialogMode" @save="handleTeamDialogSave" />

    <!-- 删除确认对话框 -->
    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h5">确认删除</v-card-title>
        <v-card-text>
          确定要删除团队 <strong>{{ teamToDelete?.display_name || teamToDelete?.name }}</strong> 吗？此操作不可恢复。
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="deleteDialog = false">取消</v-btn>
          <v-btn color="error" :loading="loading" @click="confirmDelete">删除</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 添加成员对话框 -->
    <v-dialog v-model="addMemberDialog" max-width="400">
      <v-card>
        <v-card-title>添加团队成员</v-card-title>
        <v-card-text>
          <v-text-field
            v-model.number="newMember.user_id"
            label="用户 ID"
            type="number"
            variant="outlined"
            required
            class="mb-2"
          />
          <v-select
            v-model="newMember.role"
            label="角色"
            :items="roleOptions"
            item-title="title"
            item-value="value"
            variant="outlined"
            required
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="addMemberDialog = false">取消</v-btn>
          <v-btn color="primary" :loading="loading" :disabled="!newMember.user_id" @click="addTeamMember"> 添加 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.org-team-panel {
  width: 100%;
}
</style>
