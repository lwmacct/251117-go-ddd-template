<script setup lang="ts">
import { ref, watch } from "vue";
import { orgMemberApi, extractList } from "@/api";
import type { OrgMemberDTO, OrgAddMemberDTO, OrgUpdateMemberRoleDTO } from "@models";
import { useServerPagination, useSnackbar } from "@/composables";
import { ITEMS_PER_PAGE_OPTIONS } from "@/composables";
import type { AppError } from "@/api/errors";

interface Props {
  orgId: number;
}

const props = defineProps<Props>();
const { success, error, warning } = useSnackbar();

const members = ref<OrgMemberDTO[]>([]);
const loading = ref(false);
const addDialog = ref(false);
const deleteDialog = ref(false);

const selectedMember = ref<OrgMemberDTO | null>(null);
const memberToDelete = ref<OrgMemberDTO | null>(null);

// 分页
const { pagination, updateTotal, onTableOptionsUpdate: baseOnTableOptionsUpdate, getParams } = useServerPagination();

// 新增成员表单
const newMember = ref<OrgAddMemberDTO>({
  user_id: 0,
  role: "member",
});

// 角色选项
const roleOptions = [
  { title: "所有者", value: "owner" },
  { title: "管理员", value: "admin" },
  { title: "成员", value: "member" },
];

// 表头配置
const headers = [
  { title: "ID", key: "id" },
  { title: "用户名", key: "username" },
  { title: "邮箱", key: "email" },
  { title: "角色", key: "role" },
  { title: "加入时间", key: "joined_at" },
  { title: "操作", key: "actions" },
];

/**
 * 获取成员列表
 */
const fetchMembers = async () => {
  if (!props.orgId) return;

  loading.value = true;
  try {
    const { limit, page } = getParams();
    const response = await orgMemberApi.apiOrgOrgIdMembersGet(props.orgId, limit, page);
    const result = extractList<OrgMemberDTO>(response.data);
    members.value = result.data;
    updateTotal(result.pagination.total, result.pagination.total_pages);
  } catch (err) {
    error((err as Error).message || "获取成员列表失败");
  } finally {
    loading.value = false;
  }
};

/**
 * 添加成员
 */
const addMember = async () => {
  if (!props.orgId || !newMember.value.user_id) return;

  loading.value = true;
  try {
    await orgMemberApi.apiOrgOrgIdMembersPost(props.orgId, newMember.value);
    success("成员添加成功");
    addDialog.value = false;
    newMember.value = { user_id: 0, role: "member" };
    await fetchMembers();
  } catch (err) {
    error((err as Error).message || "添加成员失败");
  } finally {
    loading.value = false;
  }
};

/**
 * 移除成员
 */
const removeMember = async (member: OrgMemberDTO) => {
  if (!props.orgId) return;

  memberToDelete.value = member;
  deleteDialog.value = true;
};

/**
 * 确认移除成员
 */
const confirmRemove = async () => {
  if (!props.orgId || !memberToDelete.value?.user_id) return;

  loading.value = true;
  try {
    await orgMemberApi.apiOrgOrgIdMembersUserIdDelete(props.orgId, memberToDelete.value.user_id);
    success("成员移除成功");
    deleteDialog.value = false;
    memberToDelete.value = null;
    await fetchMembers();
  } catch (err) {
    const appErr = err as AppError;
    // 400 错误是业务规则违规，显示为警告而非错误
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
 * 更新角色
 */
const updateRole = async (member: OrgMemberDTO) => {
  if (!props.orgId) return;

  selectedMember.value = member;
};

/**
 * 保存角色更新
 */
const saveRole = async () => {
  if (!props.orgId || !selectedMember.value?.user_id) return;

  loading.value = true;
  try {
    const data: OrgUpdateMemberRoleDTO = {
      role: selectedMember.value.role as "owner" | "admin" | "member",
    };
    await orgMemberApi.apiOrgOrgIdMembersUserIdRolePut(props.orgId, selectedMember.value.user_id, data);
    success("角色更新成功");
    selectedMember.value = null;
    await fetchMembers();
  } catch (err) {
    const appErr = err as AppError;
    // 400 错误是业务规则违规，显示为警告而非错误
    if (appErr.status === 400) {
      warning(appErr.message || "更新角色失败");
    } else {
      error(appErr.message || "更新角色失败");
    }
  } finally {
    loading.value = false;
  }
};

/**
 * 表格选项变化处理
 */
const onTableOptionsUpdate = (options: { page: number; itemsPerPage: number }) => {
  baseOnTableOptionsUpdate(options, fetchMembers);
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
    owner: "error",
    admin: "warning",
    member: "default",
  };
  return colors[role ?? ""] || "default";
};

// 角色文本
const getRoleText = (role?: string) => {
  const texts: Record<string, string> = {
    owner: "所有者",
    admin: "管理员",
    member: "成员",
  };
  return texts[role ?? ""] || (role ?? "-");
};

// 监听 orgId 变化，重置分页并重新加载
// 注意：首次加载由 v-data-table-server 的 @update:options 触发
watch(
  () => props.orgId,
  () => {
    if (!props.orgId) return;
    // 重置到第一页
    pagination.page = 1;
    fetchMembers();
  },
);
</script>

<template>
  <div class="org-member-panel">
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title class="d-flex align-center">
            <span>成员管理</span>
            <v-spacer />
            <v-btn color="primary" size="small" @click="addDialog = true">
              <v-icon start>mdi-plus</v-icon>
              添加成员
            </v-btn>
          </v-card-title>

          <v-card-text>
            <v-data-table-server
              density="compact"
              :items-per-page="pagination.limit"
              :items-per-page-options="ITEMS_PER_PAGE_OPTIONS"
              :page="pagination.page"
              :headers="headers"
              :items="members"
              :items-length="pagination.total"
              :loading="loading"
              loading-text="加载中..."
              no-data-text="暂无成员"
              @update:options="onTableOptionsUpdate"
            >
              <!-- ID 列 -->
              <template #item.id="{ item }">
                <span class="text-body-2">{{ item.id }}</span>
              </template>

              <!-- 角色列 -->
              <template #item.role="{ item }">
                <v-chip v-if="!selectedMember || selectedMember.id !== item.id" :color="getRoleColor(item.role)" size="small">
                  {{ getRoleText(item.role) }}
                </v-chip>
                <v-select
                  v-else
                  v-model="item.role"
                  :items="roleOptions"
                  item-title="title"
                  item-value="value"
                  density="compact"
                  variant="outlined"
                  hide-details
                  style="max-width: 120px"
                  @update:model-value="saveRole"
                />
              </template>

              <!-- 加入时间列 -->
              <template #item.joined_at="{ item }">
                <span class="text-body-2">{{ formatDate(item.joined_at) }}</span>
              </template>

              <!-- 操作列 -->
              <template #item.actions="{ item }">
                <v-tooltip text="更新角色">
                  <template #activator="{ props }">
                    <v-btn
                      icon="mdi-shield-account"
                      size="small"
                      variant="text"
                      color="primary"
                      v-bind="props"
                      @click="updateRole(item)"
                    ></v-btn>
                  </template>
                </v-tooltip>

                <v-tooltip text="移除">
                  <template #activator="{ props }">
                    <v-btn
                      icon="mdi-minus"
                      size="small"
                      variant="text"
                      color="error"
                      v-bind="props"
                      @click="removeMember(item)"
                    ></v-btn>
                  </template>
                </v-tooltip>
              </template>
            </v-data-table-server>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 添加成员对话框 -->
    <v-dialog v-model="addDialog" max-width="400">
      <v-card>
        <v-card-title>添加成员</v-card-title>
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
          <v-btn variant="text" @click="addDialog = false">取消</v-btn>
          <v-btn color="primary" :loading="loading" :disabled="!newMember.user_id" @click="addMember"> 添加 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 删除确认对话框 -->
    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h5">确认移除</v-card-title>
        <v-card-text>
          确定要移除成员 <strong>{{ memberToDelete?.username }}</strong> 吗？
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="deleteDialog = false">取消</v-btn>
          <v-btn color="error" :loading="loading" @click="confirmRemove">移除</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.org-member-panel {
  width: 100%;
}
</style>
