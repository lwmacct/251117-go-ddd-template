<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import type { AdminUser, CreateUserRequest, UpdateUserRequest } from '@/types/admin';
import RoleSelector from './RoleSelector.vue';

/**
 * 用户创建/编辑对话框
 */

interface Props {
  modelValue: boolean;
  user?: AdminUser | null;
  mode: 'create' | 'edit';
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'save', data: CreateUserRequest | UpdateUserRequest): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

// 表单数据
const formData = ref<CreateUserRequest & UpdateUserRequest>({
  username: '',
  email: '',
  password: '',
  full_name: '',
  status: 'active',
});

// 表单验证
const valid = ref(false);
const form = ref();

// 验证规则
const rules = {
  username: [
    (v: string) => !!v || '用户名不能为空',
    (v: string) => (v && v.length >= 3) || '用户名至少3个字符',
    (v: string) => /^[a-zA-Z0-9_]+$/.test(v) || '用户名只能包含字母、数字和下划线',
  ],
  email: [
    (v: string) => !!v || '邮箱不能为空',
    (v: string) => /.+@.+\..+/.test(v) || '邮箱格式不正确',
  ],
  password: [
    (v: string) => {
      if (props.mode === 'create') {
        return !!v || '密码不能为空';
      }
      return true;
    },
    (v: string) => {
      if (v && v.length > 0) {
        return v.length >= 6 || '密码至少6个字符';
      }
      return true;
    },
  ],
};

// 对话框标题
const dialogTitle = computed(() => {
  return props.mode === 'create' ? '新建用户' : '编辑用户';
});

// 监听用户数据变化，初始化表单
watch(
  () => props.user,
  (newUser) => {
    if (newUser && props.mode === 'edit') {
      formData.value = {
        username: newUser.username,
        email: newUser.email,
        full_name: newUser.full_name || '',
        status: newUser.status,
        password: '',
      };
    } else {
      resetForm();
    }
  },
  { immediate: true }
);

// 重置表单
const resetForm = () => {
  formData.value = {
    username: '',
    email: '',
    password: '',
    full_name: '',
    status: 'active',
  };
  form.value?.resetValidation();
};

// 关闭对话框
const closeDialog = () => {
  emit('update:modelValue', false);
  resetForm();
};

// 保存
const handleSave = async () => {
  const { valid: isValid } = await form.value.validate();
  if (!isValid) return;

  if (props.mode === 'create') {
    const createData: CreateUserRequest = {
      username: formData.value.username,
      email: formData.value.email,
      password: formData.value.password!,
      full_name: formData.value.full_name || undefined,
      status: formData.value.status as 'active' | 'inactive',
    };
    emit('save', createData);
  } else {
    const updateData: UpdateUserRequest = {
      email: formData.value.email,
      full_name: formData.value.full_name || undefined,
      status: formData.value.status,
    };
    emit('save', updateData);
  }

  closeDialog();
};
</script>

<template>
  <v-dialog :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)" max-width="600" persistent>
    <v-card>
      <v-card-title>
        <span class="text-h5">{{ dialogTitle }}</span>
      </v-card-title>

      <v-card-text>
        <v-form ref="form" v-model="valid">
          <v-text-field v-model="formData.username" label="用户名" :rules="rules.username" :disabled="mode === 'edit'" variant="outlined" density="comfortable" class="mb-2" hint="只能包含字母、数字和下划线" persistent-hint></v-text-field>

          <v-text-field v-model="formData.email" label="邮箱" type="email" :rules="rules.email" variant="outlined" density="comfortable" class="mb-2"></v-text-field>

          <v-text-field v-model="formData.password" label="密码" type="password" :rules="rules.password" variant="outlined" density="comfortable" class="mb-2" :hint="mode === 'edit' ? '留空则不修改密码' : '至少6个字符'" persistent-hint></v-text-field>

          <v-text-field v-model="formData.full_name" label="全名（可选）" variant="outlined" density="comfortable" class="mb-2"></v-text-field>

          <v-select v-model="formData.status" label="状态" :items="[
            { title: '启用', value: 'active' },
            { title: '禁用', value: 'inactive' },
          ]" variant="outlined" density="comfortable" class="mb-2"></v-select>
        </v-form>
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn variant="text" @click="closeDialog">取消</v-btn>
        <v-btn color="primary" variant="elevated" @click="handleSave" :disabled="!valid">保存</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
