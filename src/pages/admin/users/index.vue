<script setup lang="ts">
import { ref } from "vue";
import { useAdminUsers } from "./composables/useAdminUsers";
import UserDialog from "./components/UserDialog.vue";
import UserImportDialog from "./components/UserImportDialog.vue";
import RoleSelector from "./components/RoleSelector.vue";
import CopyButton from "@/components/CopyButton.vue";
import type { UserUserWithRolesDTO, UserCreateUserDTO, UserUpdateUserDTO } from "@models";

/**
 * 用户管理页面
 * 用于查看和管理系统用户
 */

// 使用 composable（移除 searchUsers，搜索现在通过 watch 自动触发）
const {
  users,
  loading,
  searchQuery,
  pagination,
  errorMessage,
  successMessage,
  fetchUsers,
  createUser,
  updateUser,
  deleteUser,
  assignRoles,
  onTableOptionsUpdate,
  clearMessages,
  exportUsers,
} = useAdminUsers();

// 对话框状态
const userDialog = ref(false);
const importDialog = ref(false);
const roleSelectorDialog = ref(false);
const deleteDialog = ref(false);

// 编辑状态
const dialogMode = ref<"create" | "edit">("create");
const selectedUser = ref<UserUserWithRolesDTO | null>(null);
const userToDelete = ref<UserUserWithRolesDTO | null>(null);

// 表头配置
const headers = [
  { title: "ID", key: "id", sortable: true },
  { title: "用户名", key: "username", sortable: true },
  { title: "邮箱", key: "email", sortable: true },
  { title: "全名", key: "full_name" },
  { title: "角色", key: "roles" },
  { title: "状态", key: "status", sortable: true },
  { title: "创建时间", key: "created_at", sortable: true },
  { title: "操作", key: "actions", sortable: false },
];

// 打开创建对话框
const openCreateDialog = () => {
  dialogMode.value = "create";
  selectedUser.value = null;
  userDialog.value = true;
};

// 打开编辑对话框
const openEditDialog = (user: UserUserWithRolesDTO) => {
  dialogMode.value = "edit";
  selectedUser.value = user;
  userDialog.value = true;
};

// 打开角色选择器
const openRoleSelector = (user: UserUserWithRolesDTO) => {
  selectedUser.value = user;
  roleSelectorDialog.value = true;
};

// 打开删除确认对话框
const openDeleteDialog = (user: UserUserWithRolesDTO) => {
  userToDelete.value = user;
  deleteDialog.value = true;
};

// 保存用户（创建或编辑）
const handleSaveUser = async (data: UserCreateUserDTO | UserUpdateUserDTO) => {
  let success = false;

  if (dialogMode.value === "create") {
    success = await createUser(data as UserCreateUserDTO);
  } else if (selectedUser.value?.id) {
    success = await updateUser(selectedUser.value.id, data as UserUpdateUserDTO);
  }

  if (success) {
    userDialog.value = false;
  }
};

// 保存角色分配
const handleSaveRoles = async (roleIds: number[]) => {
  if (!selectedUser.value?.id) return;

  const success = await assignRoles(selectedUser.value.id, roleIds);
  if (success) {
    roleSelectorDialog.value = false;
  }
};

// 确认删除
const confirmDelete = async () => {
  if (!userToDelete.value?.id) return;

  const success = await deleteUser(userToDelete.value.id);
  if (success) {
    deleteDialog.value = false;
    userToDelete.value = null;
  }
};

// 导入完成处理
const handleImported = (result: { success: number; failed: number }) => {
  if (result.success > 0) {
    fetchUsers(); // 刷新用户列表
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

// 格式化角色
const formatRoles = (user: UserUserWithRolesDTO) => {
  if (!user.roles || user.roles.length === 0) return "-";
  return user.roles.map((role) => role.display_name).join(", ");
};

// 状态颜色
const getStatusColor = (status?: string) => {
  const colors: Record<string, string> = {
    active: "success",
    inactive: "warning",
    banned: "error",
  };
  return colors[status ?? ""] || "default";
};

// 状态文本
const getStatusText = (status?: string) => {
  const texts: Record<string, string> = {
    active: "启用",
    inactive: "禁用",
    banned: "封禁",
  };
  return texts[status ?? ""] || (status ?? "-");
};
</script>

<template>
  <div class="users-page">
    <!-- 标题 -->
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-6">用户管理</h1>
      </v-col>
    </v-row>

    <!-- 消息提示 -->
    <v-row v-if="errorMessage || successMessage">
      <v-col cols="12">
        <v-alert v-if="errorMessage" type="error" closable @click:close="clearMessages">
          {{ errorMessage }}
        </v-alert>
        <v-alert v-if="successMessage" type="success" closable @click:close="clearMessages">
          {{ successMessage }}
        </v-alert>
      </v-col>
    </v-row>

    <!-- 用户列表卡片 -->
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-row align="center">
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="searchQuery"
                  prepend-inner-icon="mdi-magnify"
                  label="搜索用户（用户名或邮箱）"
                  single-line
                  hide-details
                  clearable
                  variant="outlined"
                  density="compact"
                  placeholder="输入后自动搜索..."
                ></v-text-field>
              </v-col>
              <v-col cols="12" md="6" class="text-right">
                <v-btn variant="outlined" class="mr-2" :loading="loading" @click="exportUsers">
                  <v-icon start>mdi-download</v-icon>
                  导出
                </v-btn>
                <v-btn variant="outlined" class="mr-2" @click="importDialog = true">
                  <v-icon start>mdi-upload</v-icon>
                  批量导入
                </v-btn>
                <v-btn color="primary" @click="openCreateDialog">
                  <v-icon start>mdi-plus</v-icon>
                  新建用户
                </v-btn>
              </v-col>
            </v-row>
          </v-card-title>

          <v-card-text>
            <v-data-table-server
              :items-per-page="pagination.limit"
              :page="pagination.page"
              :headers="headers"
              :items="users"
              :items-length="pagination.total"
              :loading="loading"
              loading-text="加载中..."
              no-data-text="暂无用户数据"
              @update:options="onTableOptionsUpdate"
            >
              <!-- ID 列 -->
              <template #item.id="{ item }">
                <div class="d-flex align-center">
                  <span>{{ item.id }}</span>
                  <CopyButton :text="String(item.id)" size="x-small" />
                </div>
              </template>

              <!-- 角色列 -->
              <template #item.roles="{ item }">
                <span class="text-body-2">{{ formatRoles(item) }}</span>
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
                <v-tooltip text="编辑">
                  <template #activator="{ props }">
                    <v-btn icon="mdi-pencil" size="small" variant="text" v-bind="props" @click="openEditDialog(item)"></v-btn>
                  </template>
                </v-tooltip>

                <v-tooltip text="分配角色">
                  <template #activator="{ props }">
                    <v-btn
                      icon="mdi-shield-account"
                      size="small"
                      variant="text"
                      color="primary"
                      v-bind="props"
                      @click="openRoleSelector(item)"
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

    <!-- 创建/编辑对话框 -->
    <UserDialog v-model="userDialog" :user="selectedUser" :mode="dialogMode" @save="handleSaveUser" />

    <!-- 批量导入对话框 -->
    <UserImportDialog v-model="importDialog" @imported="handleImported" />

    <!-- 角色选择器 -->
    <RoleSelector
      v-if="selectedUser?.id && roleSelectorDialog"
      v-model="roleSelectorDialog"
      :user-id="selectedUser.id"
      :user-roles="
        (selectedUser.roles || []).filter((r): r is { id: number; name?: string; display_name?: string } => r.id !== undefined)
      "
      @save="handleSaveRoles"
    />

    <!-- 删除确认对话框 -->
    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h5">确认删除</v-card-title>
        <v-card-text>
          确定要删除用户 <strong>{{ userToDelete?.username }}</strong> 吗？此操作不可恢复。
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

<style scoped>
.users-page {
  width: 100%;
}
</style>
