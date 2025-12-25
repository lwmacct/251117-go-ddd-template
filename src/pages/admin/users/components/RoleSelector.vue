<script setup lang="ts">
import { ref, onMounted } from "vue";
import { adminRoleApi, extractList, type RoleRoleDTO } from "@/api";

/**
 * 角色选择器组件
 * 用于为用户分配角色
 */

interface Props {
  modelValue: boolean;
  userId: number;
  userRoles: RoleRoleDTO[];
}

interface Emits {
  (e: "update:modelValue", value: boolean): void;
  (e: "save", roleIds: number[]): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

// 状态
const loading = ref(false);
const roles = ref<RoleRoleDTO[]>([]);
const selectedRoleIds = ref<number[]>([]);
const errorMessage = ref("");

// 获取所有角色
const fetchRoles = async () => {
  loading.value = true;
  errorMessage.value = "";

  try {
    const response = await adminRoleApi.apiAdminRolesGet(1, 100);
    const result = extractList<RoleRoleDTO>(response.data);
    roles.value = result.data;

    // 初始化已选角色（过滤掉 undefined）
    selectedRoleIds.value = props.userRoles.map((role) => role.id).filter((id): id is number => id !== undefined);
  } catch (error) {
    errorMessage.value = (error as Error).message || "获取角色列表失败";
    console.error("Failed to fetch roles:", error);
  } finally {
    loading.value = false;
  }
};

// 关闭对话框
const closeDialog = () => {
  emit("update:modelValue", false);
};

// 保存角色分配
const handleSave = () => {
  emit("save", selectedRoleIds.value);
  closeDialog();
};

// 初始化
onMounted(() => {
  if (props.modelValue) {
    fetchRoles();
  }
});
</script>

<template>
  <v-dialog :model-value="modelValue" max-width="500" @update:model-value="emit('update:modelValue', $event)">
    <v-card>
      <v-card-title>
        <span class="text-h5">分配角色</span>
      </v-card-title>

      <v-card-text>
        <v-alert v-if="errorMessage" type="error" class="mb-4" closable @click:close="errorMessage = ''">
          {{ errorMessage }}
        </v-alert>

        <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-4"></v-progress-linear>

        <v-list v-if="!loading && roles.length > 0">
          <v-list-item v-for="role in roles" :key="role.id">
            <template #prepend>
              <v-checkbox-btn v-model="selectedRoleIds" :value="role.id"></v-checkbox-btn>
            </template>

            <v-list-item-title>
              {{ role.display_name }}
              <v-chip v-if="role.is_system" size="x-small" class="ml-2" color="primary">系统</v-chip>
            </v-list-item-title>

            <v-list-item-subtitle v-if="role.description">
              {{ role.description }}
            </v-list-item-subtitle>
          </v-list-item>
        </v-list>

        <v-alert v-if="!loading && roles.length === 0" type="info"> 暂无可用角色 </v-alert>
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn variant="text" @click="closeDialog">取消</v-btn>
        <v-btn color="primary" variant="elevated" :disabled="loading" @click="handleSave">保存</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
