<script setup lang="ts">
import { ref } from "vue";
import { useAdminOrgs } from "./composables/useAdminOrgs";
import { ITEMS_PER_PAGE_OPTIONS } from "@/composables";
import OrgDialog from "./components/OrgDialog.vue";
import OrgMemberPanel from "./components/OrgMemberPanel.vue";
import OrgTeamPanel from "./components/OrgTeamPanel.vue";
import CopyButton from "@/components/CopyButton.vue";
import type { OrgOrgDTO, OrgCreateOrgDTO, OrgUpdateOrgDTO } from "@models";

/**
 * 组织管理页面
 * 用于查看和管理系统组织、成员和团队
 */

// 使用 composable
const {
  orgs,
  loading,
  searchQuery,
  statusFilter,
  pagination,
  fetchOrgs: _fetchOrgs,
  createOrg,
  updateOrg,
  deleteOrg,
  onTableOptionsUpdate,
  exportOrgs,
} = useAdminOrgs();

// 对话框状态
const orgDialog = ref(false);
const deleteDialog = ref(false);

// 编辑状态
const dialogMode = ref<"create" | "edit">("create");
const selectedOrg = ref<OrgOrgDTO | null>(null);
const orgToDelete = ref<OrgOrgDTO | null>(null);

// 标签页状态
const activeTab = ref<"list" | "members" | "teams">("list");
const selectedOrgForManage = ref<OrgOrgDTO | null>(null);

// 状态选项
const statusOptions = [
  { title: "全部", value: "" },
  { title: "启用", value: "active" },
  { title: "禁用", value: "suspended" },
];

// 表头配置
const headers = [
  { title: "ID", key: "id", sortable: true },
  { title: "组织标识", key: "name", sortable: true },
  { title: "组织名称", key: "display_name", sortable: true },
  { title: "描述", key: "description" },
  { title: "状态", key: "status", sortable: true },
  { title: "创建时间", key: "created_at", sortable: true },
  { title: "操作", key: "actions", sortable: false },
];

// 打开创建对话框
const openCreateDialog = () => {
  dialogMode.value = "create";
  selectedOrg.value = null;
  orgDialog.value = true;
};

// 打开编辑对话框
const openEditDialog = (org: OrgOrgDTO) => {
  dialogMode.value = "edit";
  selectedOrg.value = org;
  orgDialog.value = true;
};

// 打开删除确认对话框
const openDeleteDialog = (org: OrgOrgDTO) => {
  orgToDelete.value = org;
  deleteDialog.value = true;
};

// 管理成员
const manageMembers = (org: OrgOrgDTO) => {
  selectedOrgForManage.value = org;
  activeTab.value = "members";
};

// 管理团队
const manageTeams = (org: OrgOrgDTO) => {
  selectedOrgForManage.value = org;
  activeTab.value = "teams";
};

// 保存组织（创建或编辑）
const handleSaveOrg = async (data: OrgCreateOrgDTO | OrgUpdateOrgDTO) => {
  let success = false;

  if (dialogMode.value === "create") {
    success = await createOrg(data as OrgCreateOrgDTO);
  } else if (selectedOrg.value?.id) {
    success = await updateOrg(selectedOrg.value.id, data as OrgUpdateOrgDTO);
  }

  if (success) {
    orgDialog.value = false;
  }
};

// 确认删除
const confirmDelete = async () => {
  if (!orgToDelete.value?.id) return;

  const success = await deleteOrg(orgToDelete.value.id);
  if (success) {
    deleteDialog.value = false;
    orgToDelete.value = null;
  }
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

// 状态颜色
const getStatusColor = (status?: string) => {
  const colors: Record<string, string> = {
    active: "success",
    inactive: "error",
    suspended: "warning",
  };
  return colors[status ?? ""] || "default";
};

// 状态文本
const getStatusText = (status?: string) => {
  const texts: Record<string, string> = {
    active: "启用",
    inactive: "禁用",
    suspended: "暂停",
  };
  return texts[status ?? ""] || (status ?? "-");
};

// 标签页配置
const tabs = computed(() => {
  const baseTabs = [{ value: "list", title: "组织列表" }];
  if (selectedOrgForManage.value) {
    baseTabs.push(
      { value: "members", title: `成员管理 (${selectedOrgForManage.value.display_name})` },
      { value: "teams", title: `团队管理 (${selectedOrgForManage.value.display_name})` },
    );
  }
  return baseTabs;
});
</script>

<template>
  <div class="orgs-page">
    <!-- 标题 -->
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-6">组织管理</h1>
      </v-col>
    </v-row>

    <!-- 标签页 -->
    <v-tabs v-model="activeTab" color="primary">
      <v-tab v-for="tab in tabs" :key="tab.value" :value="tab.value">
        {{ tab.title }}
      </v-tab>
      <v-btn
        v-if="selectedOrgForManage && activeTab !== 'list'"
        variant="text"
        size="small"
        @click="
          selectedOrgForManage = null;
          activeTab = 'list';
        "
      >
        <v-icon start>mdi-close</v-icon>
        清除选择
      </v-btn>
    </v-tabs>

    <v-window v-model="activeTab" class="mt-4">
      <!-- 组织列表 -->
      <v-window-item value="list">
        <v-row>
          <v-col cols="12">
            <v-card>
              <v-card-title>
                <v-row align="center">
                  <v-col cols="12" md="4">
                    <v-text-field
                      v-model="searchQuery"
                      prepend-inner-icon="mdi-magnify"
                      label="搜索组织"
                      single-line
                      hide-details
                      clearable
                      variant="outlined"
                      density="compact"
                      placeholder="输入后自动搜索..."
                    ></v-text-field>
                  </v-col>
                  <v-col cols="12" md="2">
                    <v-select
                      v-model="statusFilter"
                      label="状态筛选"
                      :items="statusOptions"
                      item-title="title"
                      item-value="value"
                      hide-details
                      clearable
                      variant="outlined"
                      density="compact"
                    ></v-select>
                  </v-col>
                  <v-col cols="12" md="6" class="text-right">
                    <v-btn variant="outlined" class="mr-2" :loading="loading" @click="exportOrgs">
                      <v-icon start>mdi-download</v-icon>
                      导出
                    </v-btn>
                    <v-btn color="primary" @click="openCreateDialog">
                      <v-icon start>mdi-plus</v-icon>
                      新建组织
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
                  :items="orgs"
                  :items-length="pagination.total"
                  :loading="loading"
                  loading-text="加载中..."
                  no-data-text="暂无组织数据"
                  @update:options="onTableOptionsUpdate"
                >
                  <!-- ID 列 -->
                  <template #item.id="{ item }">
                    <div class="d-flex align-center">
                      <span>{{ item.id }}</span>
                      <CopyButton :text="String(item.id)" size="x-small" />
                    </div>
                  </template>

                  <!-- 组织标识列 -->
                  <template #item.name="{ item }">
                    <span class="font-weight-medium">{{ item.name }}</span>
                  </template>

                  <!-- 描述列 -->
                  <template #item.description="{ item }">
                    <span class="text-body-2 text-truncate d-inline-block" style="max-width: 200px">
                      {{ item.description || "-" }}
                    </span>
                  </template>

                  <!-- 状态列 -->
                  <template #item.status="{ item }">
                    <v-chip :color="getStatusColor(item.status)" size="small">
                      {{ getStatusText(item.status) }}
                    </v-chip>
                  </template>

                  <!-- 创建时间列 -->
                  <template #item.created_at="{ item }">
                    <span class="text-body-2">{{ formatDate(item.created_at) }}</span>
                  </template>

                  <!-- 操作列 -->
                  <template #item.actions="{ item }">
                    <v-tooltip text="成员管理">
                      <template #activator="{ props }">
                        <v-btn
                          icon="mdi-account-group"
                          size="small"
                          variant="text"
                          color="primary"
                          v-bind="props"
                          @click="manageMembers(item)"
                        ></v-btn>
                      </template>
                    </v-tooltip>

                    <v-tooltip text="团队管理">
                      <template #activator="{ props }">
                        <v-btn
                          icon="mdi-account-multiple"
                          size="small"
                          variant="text"
                          color="info"
                          v-bind="props"
                          @click="manageTeams(item)"
                        ></v-btn>
                      </template>
                    </v-tooltip>

                    <v-tooltip text="编辑">
                      <template #activator="{ props }">
                        <v-btn
                          icon="mdi-pencil"
                          size="small"
                          variant="text"
                          v-bind="props"
                          @click="openEditDialog(item)"
                        ></v-btn>
                      </template>
                    </v-tooltip>

                    <v-tooltip text="删除">
                      <template #activator="{ props }">
                        <v-btn
                          icon="mdi-delete"
                          size="small"
                          variant="text"
                          color="error"
                          v-bind="props"
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
      </v-window-item>

      <!-- 成员管理 -->
      <v-window-item value="members">
        <OrgMemberPanel v-if="selectedOrgForManage?.id" :org-id="selectedOrgForManage.id" />
        <v-alert v-else type="info">请先选择一个组织</v-alert>
      </v-window-item>

      <!-- 团队管理 -->
      <v-window-item value="teams">
        <OrgTeamPanel v-if="selectedOrgForManage?.id" :org-id="selectedOrgForManage.id" />
        <v-alert v-else type="info">请先选择一个组织</v-alert>
      </v-window-item>
    </v-window>

    <!-- 创建/编辑对话框 -->
    <OrgDialog v-model="orgDialog" :org="selectedOrg" :mode="dialogMode" @save="handleSaveOrg" />

    <!-- 删除确认对话框 -->
    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h5">确认删除</v-card-title>
        <v-card-text>
          确定要删除组织 <strong>{{ orgToDelete?.display_name || orgToDelete?.name }}</strong> 吗？此操作不可恢复。
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="deleteDialog = false">取消</v-btn>
          <v-btn color="error" variant="elevated" :loading="loading" @click="confirmDelete">删除</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script lang="ts">
import { computed } from "vue";
</script>

<style scoped>
.orgs-page {
  width: 100%;
}
</style>
