<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import type { Menu, CreateMenuRequest, UpdateMenuRequest } from '@/types/admin';

interface Props {
  modelValue: boolean;
  menu?: Menu | null;
  mode: 'create' | 'edit';
  parentMenus: Menu[];
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'save', data: CreateMenuRequest | UpdateMenuRequest): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const formData = ref<CreateMenuRequest & UpdateMenuRequest>({
  title: '',
  path: '',
  icon: '',
  parent_id: undefined,
  order: 0,
  visible: true,
});

const valid = ref(false);
const form = ref();

const rules = {
  title: [(v: string) => !!v || '标题不能为空'],
  path: [(v: string) => !!v || '路径不能为空'],
};

const dialogTitle = computed(() => (props.mode === 'create' ? '新建菜单' : '编辑菜单'));

// 扁平化菜单列表（用于父菜单选择）
const flatMenus = computed(() => {
  const result: Array<{ id: number; title: string; level: number }> = [];
  const flatten = (menus: Menu[], level = 0) => {
    menus.forEach((menu) => {
      if (props.mode === 'edit' && props.menu && menu.id === props.menu.id) {
        return; // 不能选择自己作为父菜单
      }
      result.push({ id: menu.id, title: '  '.repeat(level) + menu.title, level });
      if (menu.children && menu.children.length > 0) {
        flatten(menu.children, level + 1);
      }
    });
  };
  flatten(props.parentMenus);
  return result;
});

watch(
  () => props.menu,
  (newMenu) => {
    if (newMenu && props.mode === 'edit') {
      formData.value = {
        title: newMenu.title,
        path: newMenu.path,
        icon: newMenu.icon || '',
        parent_id: newMenu.parent_id,
        order: newMenu.order,
        visible: newMenu.visible,
      };
    } else {
      resetForm();
    }
  },
  { immediate: true }
);

const resetForm = () => {
  formData.value = {
    title: '',
    path: '',
    icon: '',
    parent_id: undefined,
    order: 0,
    visible: true,
  };
  form.value?.resetValidation();
};

const closeDialog = () => {
  emit('update:modelValue', false);
  resetForm();
};

const handleSave = async () => {
  const { valid: isValid } = await form.value.validate();
  if (!isValid) return;

  emit('save', formData.value);
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
          <v-text-field v-model="formData.title" label="菜单标题" :rules="rules.title" variant="outlined" density="comfortable" class="mb-2"></v-text-field>

          <v-text-field v-model="formData.path" label="路由路径" :rules="rules.path" variant="outlined" density="comfortable" class="mb-2" hint="如: /admin/users"></v-text-field>

          <v-text-field v-model="formData.icon" label="图标（可选）" variant="outlined" density="comfortable" class="mb-2" hint="MDI 图标名称，如: mdi-account"></v-text-field>

          <v-select v-model="formData.parent_id" label="父菜单（可选）" :items="[{ id: undefined, title: '无（顶级菜单）', level: 0 }, ...flatMenus]" item-title="title" item-value="id" variant="outlined" density="comfortable" class="mb-2" clearable></v-select>

          <v-text-field v-model.number="formData.order" label="排序" type="number" variant="outlined" density="comfortable" class="mb-2" hint="数字越小越靠前"></v-text-field>

          <v-switch v-model="formData.visible" label="是否可见" color="primary"></v-switch>
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
