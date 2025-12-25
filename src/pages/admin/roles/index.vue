<script setup lang="ts">
import { ref } from "vue";
import { useRoles } from "./composables/useRoles";
import RoleDialog from "./components/RoleDialog.vue";
import PermissionSelector from "./components/PermissionSelector.vue";
import type { RoleRoleDTO, RoleCreateRoleDTO, RoleUpdateRoleDTO } from "@models";

// 使用 composable（搜索现在通过 watch 自动触发）
const {
  roles,
  loading,
  searchQuery,
  pagination,
  errorMessage,
  successMessage,
  createRole,
  updateRole,
  deleteRole,
  setPermissions,
  onTableOptionsUpdate,
  clearMessages,
  exportRoles,
} = useRoles();

const roleDialog = ref(false);
const permissionDialog = ref(false);
const deleteDialog = ref(false);
const dialogMode = ref<"create" | "edit">("create");
const selectedRole = ref<RoleRoleDTO | null>(null);
const roleToDelete = ref<RoleRoleDTO | null>(null);

const headers = [
  { title: "ID", key: "id", sortable: true },
  { title: "角色标识", key: "name", sortable: true },
  { title: "显示名称", key: "display_name", sortable: true },
  { title: "描述", key: "description" },
  { title: "权限数量", key: "permissions" },
  { title: "系统角色", key: "is_system" },
  { title: "创建时间", key: "created_at", sortable: true },
  { title: "操作", key: "actions", sortable: false },
];

const openCreateDialog = () => {
  dialogMode.value = "create";
  selectedRole.value = null;
  roleDialog.value = true;
};

const openEditDialog = (role: RoleRoleDTO) => {
  if (role.is_system) {
    errorMessage.value = "系统角色不能编辑";
    return;
  }
  dialogMode.value = "edit";
  selectedRole.value = role;
  roleDialog.value = true;
};

const openPermissionSelector = (role: RoleRoleDTO) => {
  selectedRole.value = role;
  permissionDialog.value = true;
};

const openDeleteDialog = (role: RoleRoleDTO) => {
  if (role.is_system) {
    errorMessage.value = "系统角色不能删除";
    return;
  }
  roleToDelete.value = role;
  deleteDialog.value = true;
};

const handleSaveRole = async (data: RoleCreateRoleDTO | RoleUpdateRoleDTO) => {
  let success = false;

  if (dialogMode.value === "create") {
    success = await createRole(data as RoleCreateRoleDTO);
  } else if (selectedRole.value?.id) {
    success = await updateRole(selectedRole.value.id, data as RoleUpdateRoleDTO);
  }

  if (success) {
    roleDialog.value = false;
  }
};

const handleSavePermissions = async (permissionIds: number[]) => {
  if (!selectedRole.value?.id) return;

  const success = await setPermissions(selectedRole.value.id, permissionIds);
  if (success) {
    permissionDialog.value = false;
  }
};

const confirmDelete = async () => {
  if (!roleToDelete.value?.id) return;

  const success = await deleteRole(roleToDelete.value.id);
  if (success) {
    deleteDialog.value = false;
    roleToDelete.value = null;
  }
};

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

const formatPermissions = (role: RoleRoleDTO) => {
  return role.permissions?.length || 0;
};
</script>

<template>
  <div class="roles-page">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-6">角色管理</h1>
      </v-col>
    </v-row>

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

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-row align="center">
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="searchQuery"
                  prepend-inner-icon="mdi-magnify"
                  label="搜索角色"
                  single-line
                  hide-details
                  clearable
                  variant="outlined"
                  density="compact"
                  placeholder="输入后自动搜索..."
                ></v-text-field>
              </v-col>
              <v-col cols="12" md="6" class="text-right">
                <v-btn variant="outlined" class="mr-2" :loading="loading" @click="exportRoles">
                  <v-icon start>mdi-download</v-icon>
                  导出
                </v-btn>
                <v-btn color="primary" @click="openCreateDialog">
                  <v-icon start>mdi-plus</v-icon>
                  新建角色
                </v-btn>
              </v-col>
            </v-row>
          </v-card-title>

          <v-card-text>
            <v-data-table-server
              :items-per-page="pagination.limit"
              :page="pagination.page"
              :headers="headers"
              :items="roles"
              :items-length="pagination.total"
              :loading="loading"
              loading-text="加载中..."
              no-data-text="暂无角色数据"
              @update:options="onTableOptionsUpdate"
            >
              <template #item.permissions="{ item }">
                <v-chip size="small" color="primary"> {{ formatPermissions(item) }} 个权限 </v-chip>
              </template>

              <template #item.is_system="{ item }">
                <v-chip v-if="item.is_system" size="small" color="info"> 系统 </v-chip>
                <span v-else>-</span>
              </template>

              <template #item.created_at="{ item }">
                <span class="text-body-2">{{ formatDate(item.created_at) }}</span>
              </template>

              <template #item.actions="{ item }">
                <v-tooltip text="编辑">
                  <template #activator="{ props }">
                    <v-btn
                      icon="mdi-pencil"
                      size="small"
                      variant="text"
                      v-bind="props"
                      :disabled="item.is_system"
                      @click="openEditDialog(item)"
                    ></v-btn>
                  </template>
                </v-tooltip>

                <v-tooltip text="设置权限">
                  <template #activator="{ props }">
                    <v-btn
                      icon="mdi-shield-lock"
                      size="small"
                      variant="text"
                      color="primary"
                      v-bind="props"
                      @click="openPermissionSelector(item)"
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
                      :disabled="item.is_system"
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

    <RoleDialog v-model="roleDialog" :role="selectedRole" :mode="dialogMode" @save="handleSaveRole" />

    <PermissionSelector
      v-if="selectedRole?.id && permissionDialog"
      v-model="permissionDialog"
      :role-id="selectedRole.id"
      :role-permissions="selectedRole.permissions || []"
      @save="handleSavePermissions"
    />

    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title class="text-h5">确认删除</v-card-title>
        <v-card-text>
          确定要删除角色 <strong>{{ roleToDelete?.display_name }}</strong> 吗？此操作不可恢复。
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
.roles-page {
  width: 100%;
}
</style>
