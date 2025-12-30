<script setup lang="ts">
import { ref, watch, computed } from "vue";
import type { RolePermissionDTO, RolePermissionInputDTO } from "@models";

interface Props {
  modelValue: boolean;
  roleId: number;
  rolePermissions: RolePermissionDTO[];
}

interface Emits {
  (e: "update:modelValue", value: boolean): void;
  (e: "save", permissions: RolePermissionInputDTO[]): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

// 编辑中的权限列表
const editingPermissions = ref<RolePermissionInputDTO[]>([]);
const errorMessage = ref("");

// 预定义的 Operation 模式示例（URN 风格）
const operationExamples = [
  { pattern: "sys:users:*", desc: "用户管理所有操作" },
  { pattern: "sys:roles:*", desc: "角色管理所有操作" },
  { pattern: "sys:settings:*", desc: "设置管理所有操作" },
  { pattern: "self:profile:*", desc: "个人资料所有操作" },
  { pattern: "self:tokens:*", desc: "PAT 令牌所有操作" },
  { pattern: "*:*:*", desc: "超级管理员（所有操作）" },
];

// 预定义的 Resource 模式示例（URN 风格）
const resourceExamples = [
  { pattern: "*:*:*", desc: "所有资源" },
  { pattern: "self:user:@me", desc: "仅限自身" },
  { pattern: "sys:user:*", desc: "所有系统用户" },
];

// 初始化编辑数据
const initPermissions = () => {
  editingPermissions.value = props.rolePermissions.map((p) => ({
    operation_pattern: p.operation_pattern ?? "",
    resource_pattern: p.resource_pattern ?? "*:*:*",
  }));
};

// 添加新权限
const addPermission = () => {
  editingPermissions.value.push({
    operation_pattern: "",
    resource_pattern: "*:*:*",
  });
};

// 删除权限
const removePermission = (index: number) => {
  editingPermissions.value.splice(index, 1);
};

// 验证权限
const isValid = computed(() => {
  return editingPermissions.value.every((p) => p.operation_pattern.trim() !== "");
});

const closeDialog = () => {
  emit("update:modelValue", false);
};

const handleSave = () => {
  if (!isValid.value) {
    errorMessage.value = "所有权限必须填写 Operation 模式";
    return;
  }
  emit("save", editingPermissions.value);
  closeDialog();
};

// 监听对话框打开
watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal) {
      initPermissions();
      errorMessage.value = "";
    }
  },
  { immediate: true },
);
</script>

<template>
  <v-dialog :model-value="modelValue" max-width="800" scrollable @update:model-value="emit('update:modelValue', $event)">
    <v-card>
      <v-card-title>
        <span class="text-h5">设置权限</span>
      </v-card-title>

      <v-card-text style="max-height: 500px">
        <v-alert v-if="errorMessage" type="error" class="mb-4" closable @click:close="errorMessage = ''">
          {{ errorMessage }}
        </v-alert>

        <!-- 说明 -->
        <v-alert type="info" variant="tonal" class="mb-4">
          <div class="text-subtitle-2 mb-2">URN 风格 RBAC 说明</div>
          <div class="text-body-2">
            权限由 <strong>Operation 模式</strong> + <strong>Resource 模式</strong> 组成。
            <br />
            • URN 格式：<code>scope:type:identifier</code>
            <br />
            • Operation 示例：<code>sys:users:create</code>、<code>self:profile:*</code>
            <br />
            • Resource 示例：<code>*:*:*</code>（所有）、<code>self:user:@me</code>（自身）
            <br />
            • Scope 类型：<code>sys</code>（系统）、<code>self</code>（用户自身）、<code>public</code>（公开）
          </div>
        </v-alert>

        <!-- 权限列表 -->
        <div class="mb-4">
          <div v-for="(perm, index) in editingPermissions" :key="index" class="d-flex align-center gap-2 mb-2">
            <v-text-field
              v-model="perm.operation_pattern"
              label="Operation 模式"
              placeholder="如: sys:users:*"
              density="compact"
              variant="outlined"
              style="flex: 2"
              :rules="[(v) => !!v || 'Operation 模式必填']"
            />
            <v-text-field
              v-model="perm.resource_pattern"
              label="Resource 模式"
              placeholder="默认: *:*:*"
              density="compact"
              variant="outlined"
              style="flex: 1"
            />
            <v-btn icon="mdi-delete" variant="text" color="error" size="small" @click="removePermission(index)" />
          </div>

          <v-btn variant="tonal" color="primary" size="small" prepend-icon="mdi-plus" @click="addPermission"> 添加权限 </v-btn>
        </div>

        <!-- 快速添加示例 -->
        <v-expansion-panels variant="accordion" class="mb-4">
          <v-expansion-panel title="常用 Operation 模式">
            <v-expansion-panel-text>
              <v-chip
                v-for="ex in operationExamples"
                :key="ex.pattern"
                class="ma-1"
                size="small"
                @click="
                  editingPermissions.push({
                    operation_pattern: ex.pattern,
                    resource_pattern: '*:*:*',
                  })
                "
              >
                {{ ex.pattern }} - {{ ex.desc }}
              </v-chip>
            </v-expansion-panel-text>
          </v-expansion-panel>
          <v-expansion-panel title="常用 Resource 模式">
            <v-expansion-panel-text>
              <v-chip v-for="ex in resourceExamples" :key="ex.pattern" class="ma-1" size="small" disabled>
                {{ ex.pattern }} - {{ ex.desc }}
              </v-chip>
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn variant="text" @click="closeDialog">取消</v-btn>
        <v-btn color="primary" variant="elevated" :disabled="!isValid" @click="handleSave">保存</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
