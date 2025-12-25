<script setup lang="ts">
import { ref, watch, computed } from "vue";
import type { UserUserWithRolesDTO, UserCreateUserDTO, UserUpdateUserDTO } from "@models";
// RoleSelector is reserved for future role assignment feature
import _RoleSelector from "./RoleSelector.vue";
import PasswordStrengthIndicator from "@/components/PasswordStrengthIndicator.vue";

/**
 * 用户创建/编辑对话框
 */

interface Props {
  modelValue: boolean;
  user?: UserUserWithRolesDTO | null;
  mode: "create" | "edit";
}

interface Emits {
  (e: "update:modelValue", value: boolean): void;
  (e: "save", data: UserCreateUserDTO | UserUpdateUserDTO): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

// 表单数据
const formData = ref<{
  username: string;
  email: string;
  password: string;
  full_name: string;
  status: "active" | "inactive" | "banned";
}>({
  username: "",
  email: "",
  password: "",
  full_name: "",
  status: "active",
});

// 表单验证
const valid = ref(false);
const form = ref();

// 验证规则
const rules = {
  username: [
    (v: string) => !!v || "用户名不能为空",
    (v: string) => (v && v.length >= 3) || "用户名至少3个字符",
    (v: string) => /^[a-zA-Z0-9_]+$/.test(v) || "用户名只能包含字母、数字和下划线",
  ],
  email: [(v: string) => !!v || "邮箱不能为空", (v: string) => /.+@.+\..+/.test(v) || "邮箱格式不正确"],
  password: [
    (v: string) => {
      if (props.mode === "create") {
        return !!v || "密码不能为空";
      }
      return true;
    },
    (v: string) => {
      if (v && v.length > 0) {
        return v.length >= 6 || "密码至少6个字符";
      }
      return true;
    },
  ],
};

// 对话框标题
const dialogTitle = computed(() => {
  return props.mode === "create" ? "新建用户" : "编辑用户";
});

// 重置表单（需要在 watch 之前定义）
const resetForm = () => {
  formData.value = {
    username: "",
    email: "",
    password: "",
    full_name: "",
    status: "active",
  };
  form.value?.resetValidation();
};

// 监听用户数据变化，初始化表单
watch(
  () => props.user,
  (newUser) => {
    if (newUser && props.mode === "edit") {
      formData.value = {
        username: newUser.username ?? "",
        email: newUser.email ?? "",
        full_name: newUser.full_name ?? "",
        status: (newUser.status as "active" | "inactive" | "banned") ?? "active",
        password: "",
      };
    } else {
      resetForm();
    }
  },
  { immediate: true },
);

// 关闭对话框
const closeDialog = () => {
  emit("update:modelValue", false);
  resetForm();
};

// 保存
const handleSave = async () => {
  const { valid: isValid } = await form.value.validate();
  if (!isValid) return;

  if (props.mode === "create") {
    const createData: UserCreateUserDTO = {
      username: formData.value.username,
      email: formData.value.email,
      password: formData.value.password!,
      full_name: formData.value.full_name || undefined,
      status: formData.value.status as "active" | "inactive",
    };
    emit("save", createData);
  } else {
    const updateData: UserUpdateUserDTO = {
      email: formData.value.email,
      full_name: formData.value.full_name || undefined,
      status: formData.value.status,
    };
    emit("save", updateData);
  }

  closeDialog();
};
</script>

<template>
  <v-dialog :model-value="modelValue" max-width="600" persistent @update:model-value="emit('update:modelValue', $event)">
    <v-card>
      <v-card-title>
        <span class="text-h5">{{ dialogTitle }}</span>
      </v-card-title>

      <v-card-text>
        <v-form ref="form" v-model="valid">
          <v-text-field
            v-model="formData.username"
            label="用户名"
            :rules="rules.username"
            :disabled="mode === 'edit'"
            variant="outlined"
            density="comfortable"
            class="mb-2"
            hint="只能包含字母、数字和下划线"
            persistent-hint
          ></v-text-field>

          <v-text-field
            v-model="formData.email"
            label="邮箱"
            type="email"
            :rules="rules.email"
            variant="outlined"
            density="comfortable"
            class="mb-2"
          ></v-text-field>

          <v-text-field
            v-model="formData.password"
            label="密码"
            type="password"
            :rules="rules.password"
            variant="outlined"
            density="comfortable"
            class="mb-2"
            :hint="mode === 'edit' ? '留空则不修改密码' : '至少6个字符'"
            persistent-hint
          ></v-text-field>

          <!-- 密码强度指示器（仅创建模式显示） -->
          <PasswordStrengthIndicator v-if="mode === 'create'" :password="formData.password" :show-hints="false" class="mb-4" />

          <v-text-field
            v-model="formData.full_name"
            label="全名（可选）"
            variant="outlined"
            density="comfortable"
            class="mb-2"
          ></v-text-field>

          <v-select
            v-model="formData.status"
            label="状态"
            :items="
              mode === 'create'
                ? [
                    { title: '启用', value: 'active' },
                    { title: '禁用', value: 'inactive' },
                  ]
                : [
                    { title: '启用', value: 'active' },
                    { title: '禁用', value: 'inactive' },
                    { title: '封禁', value: 'banned' },
                  ]
            "
            variant="outlined"
            density="comfortable"
            class="mb-2"
          ></v-select>
        </v-form>
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn variant="text" @click="closeDialog">取消</v-btn>
        <v-btn color="primary" variant="elevated" :disabled="!valid" @click="handleSave">保存</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
