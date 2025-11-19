<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { AdminPermissionsAPI } from "@/api/admin";
import type { Permission, PermissionTreeNode } from "@/types/admin";

interface Props {
  modelValue: boolean;
  roleId: number;
  rolePermissions: Permission[];
}

interface Emits {
  (e: "update:modelValue", value: boolean): void;
  (e: "save", permissionIds: number[]): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const loading = ref(false);
const permissions = ref<Permission[]>([]);
const selectedPermissionIds = ref<number[]>([]);
const errorMessage = ref("");
const expandAll = ref(false);

// 将权限列表转换为树形结构
const permissionTree = computed<PermissionTreeNode[]>(() => {
  const tree: Map<string, PermissionTreeNode> = new Map();

  permissions.value.forEach((perm) => {
    const domain = perm.domain;
    const resource = perm.resource;
    const domainKey = domain;
    const resourceKey = `${domain}:${resource}`;

    // 创建 domain 节点
    if (!tree.has(domainKey)) {
      tree.set(domainKey, {
        key: domainKey,
        label: domain,
        children: [],
      });
    }

    // 创建 resource 节点
    const domainNode = tree.get(domainKey)!;
    let resourceNode = domainNode.children?.find((c) => c.key === resourceKey);

    if (!resourceNode) {
      resourceNode = {
        key: resourceKey,
        label: resource,
        children: [],
      };
      domainNode.children!.push(resourceNode);
    }

    // 添加 action 节点（叶子节点）
    resourceNode.children!.push({
      key: `${perm.code}`,
      label: `${perm.action} ${perm.description ? `(${perm.description})` : ""}`,
      permission: perm,
    });
  });

  return Array.from(tree.values());
});

const fetchPermissions = async () => {
  loading.value = true;
  errorMessage.value = "";

  try {
    const allPermissions = await AdminPermissionsAPI.getAllPermissions();
    permissions.value = allPermissions;
    selectedPermissionIds.value = props.rolePermissions.map((p) => p.id);
  } catch (error: any) {
    errorMessage.value = error.message || "获取权限列表失败";
  } finally {
    loading.value = false;
  }
};

const closeDialog = () => {
  emit("update:modelValue", false);
};

const handleSave = () => {
  emit("save", selectedPermissionIds.value);
  closeDialog();
};

// 切换权限选择
const togglePermission = (permId: number) => {
  const index = selectedPermissionIds.value.indexOf(permId);
  if (index > -1) {
    selectedPermissionIds.value.splice(index, 1);
  } else {
    selectedPermissionIds.value.push(permId);
  }
};

onMounted(() => {
  if (props.modelValue) {
    fetchPermissions();
  }
});
</script>

<template>
  <v-dialog :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)" max-width="800" scrollable>
    <v-card>
      <v-card-title>
        <span class="text-h5">设置权限</span>
      </v-card-title>

      <v-card-text style="max-height: 500px">
        <v-alert v-if="errorMessage" type="error" class="mb-4" closable @click:close="errorMessage = ''">
          {{ errorMessage }}
        </v-alert>

        <div class="d-flex justify-end mb-2">
          <v-btn variant="text" size="small" @click="expandAll = !expandAll">
            {{ expandAll ? "折叠全部" : "展开全部" }}
          </v-btn>
        </div>

        <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-4"></v-progress-linear>

        <!-- 权限树 -->
        <div v-if="!loading && permissionTree.length > 0" class="permission-tree">
          <v-expansion-panels v-model:model-value="expandAll ? permissionTree.map((_, i) => i) : []" multiple>
            <v-expansion-panel v-for="domainNode in permissionTree" :key="domainNode.key">
              <v-expansion-panel-title>
                <v-icon start>mdi-folder</v-icon>
                <strong>{{ domainNode.label }}</strong>
              </v-expansion-panel-title>
              <v-expansion-panel-text>
                <div v-for="resourceNode in domainNode.children" :key="resourceNode.key" class="ml-4 mb-3">
                  <div class="text-subtitle-2 mb-2">
                    <v-icon size="small">mdi-file-tree</v-icon>
                    {{ resourceNode.label }}
                  </div>
                  <div class="ml-6">
                    <v-checkbox v-for="actionNode in resourceNode.children" :key="actionNode.key" v-model="selectedPermissionIds" :value="actionNode.permission!.id" :label="actionNode.label" density="compact" hide-details></v-checkbox>
                  </div>
                </div>
              </v-expansion-panel-text>
            </v-expansion-panel>
          </v-expansion-panels>
        </div>

        <v-alert v-if="!loading && permissions.length === 0" type="info"> 暂无可用权限 </v-alert>
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn variant="text" @click="closeDialog">取消</v-btn>
        <v-btn color="primary" variant="elevated" @click="handleSave" :disabled="loading">保存</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<style scoped>
.permission-tree {
  width: 100%;
}
</style>
