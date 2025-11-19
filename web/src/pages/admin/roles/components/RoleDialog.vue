<script setup lang="ts">
import { ref, watch, computed } from "vue";
import type { Role, CreateRoleRequest, UpdateRoleRequest } from "@/types/admin";

interface Props {
  modelValue: boolean;
  role?: Role | null;
  mode: "create" | "edit";
}

interface Emits {
  (e: "update:modelValue", value: boolean): void;
  (e: "save", data: CreateRoleRequest | UpdateRoleRequest): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const formData = ref<CreateRoleRequest & UpdateRoleRequest>({
  name: "",
  display_name: "",
  description: "",
});

const valid = ref(false);
const form = ref();

const rules = {
  name: [(v: string) => !!v || "角色标识不能为空", (v: string) => /^[a-z_]+$/.test(v) || "只能包含小写字母和下划线"],
  display_name: [(v: string) => !!v || "显示名称不能为空"],
};

const dialogTitle = computed(() => (props.mode === "create" ? "新建角色" : "编辑角色"));

watch(
  () => props.role,
  (newRole) => {
    if (newRole && props.mode === "edit") {
      formData.value = {
        name: newRole.name,
        display_name: newRole.display_name,
        description: newRole.description || "",
      };
    } else {
      resetForm();
    }
  },
  { immediate: true },
);

const resetForm = () => {
  formData.value = {
    name: "",
    display_name: "",
    description: "",
  };
  form.value?.resetValidation();
};

const closeDialog = () => {
  emit("update:modelValue", false);
  resetForm();
};

const handleSave = async () => {
  const { valid: isValid } = await form.value.validate();
  if (!isValid) return;

  if (props.mode === "create") {
    emit("save", formData.value as CreateRoleRequest);
  } else {
    const updateData: UpdateRoleRequest = {
      display_name: formData.value.display_name,
      description: formData.value.description,
    };
    emit("save", updateData);
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
          <v-text-field v-model="formData.name" label="角色标识" :rules="rules.name" :disabled="mode === 'edit'" variant="outlined" density="comfortable" class="mb-2" hint="如: admin, editor, viewer" persistent-hint></v-text-field>

          <v-text-field v-model="formData.display_name" label="显示名称" :rules="rules.display_name" variant="outlined" density="comfortable" class="mb-2" hint="如: 管理员, 编辑者"></v-text-field>

          <v-textarea v-model="formData.description" label="描述（可选）" variant="outlined" density="comfortable" rows="3"></v-textarea>
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
