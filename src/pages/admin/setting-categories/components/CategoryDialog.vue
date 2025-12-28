<script setup lang="ts">
import { ref, watch, computed } from "vue";
import type { SettingCategoryDTO, HandlerCreateCategoryRequest, HandlerUpdateCategoryRequest } from "@models";

interface Props {
  modelValue: boolean;
  category?: SettingCategoryDTO | null;
  mode: "create" | "edit";
}

interface Emits {
  (e: "update:modelValue", value: boolean): void;
  (e: "save", data: HandlerCreateCategoryRequest | HandlerUpdateCategoryRequest): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const formData = ref<HandlerCreateCategoryRequest & HandlerUpdateCategoryRequest>({
  key: "",
  label: "",
  icon: "mdi-cog",
  order: 0,
});

const valid = ref(false);
const form = ref();

const rules = {
  key: [(v: string) => !!v || "分类标识不能为空", (v: string) => /^[a-z_]+$/.test(v) || "只能包含小写字母和下划线"],
  label: [(v: string) => !!v || "显示名称不能为空"],
};

// 常用图标选项
const iconOptions = [
  { value: "mdi-cog", title: "mdi-cog (设置)" },
  { value: "mdi-shield-lock", title: "mdi-shield-lock (安全)" },
  { value: "mdi-bell", title: "mdi-bell (通知)" },
  { value: "mdi-cloud-upload", title: "mdi-cloud-upload (备份)" },
  { value: "mdi-palette", title: "mdi-palette (主题)" },
  { value: "mdi-account", title: "mdi-account (用户)" },
  { value: "mdi-email", title: "mdi-email (邮件)" },
  { value: "mdi-database", title: "mdi-database (数据)" },
  { value: "mdi-application", title: "mdi-application (应用)" },
  { value: "mdi-tune", title: "mdi-tune (调整)" },
];

const dialogTitle = computed(() => (props.mode === "create" ? "新建分类" : "编辑分类"));

watch(
  () => props.category,
  (newCategory) => {
    if (newCategory && props.mode === "edit") {
      formData.value = {
        key: newCategory.key ?? "",
        label: newCategory.label ?? "",
        icon: newCategory.icon ?? "mdi-cog",
        order: newCategory.order ?? 0,
      };
    } else {
      resetForm();
    }
  },
  { immediate: true },
);

function resetForm() {
  formData.value = {
    key: "",
    label: "",
    icon: "mdi-cog",
    order: 0,
  };
  form.value?.resetValidation();
}

const closeDialog = () => {
  emit("update:modelValue", false);
  resetForm();
};

const handleSave = async () => {
  const { valid: isValid } = await form.value.validate();
  if (!isValid) return;

  if (props.mode === "create") {
    emit("save", formData.value as HandlerCreateCategoryRequest);
  } else {
    const updateData: HandlerUpdateCategoryRequest = {
      label: formData.value.label,
      icon: formData.value.icon,
      order: formData.value.order,
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
            v-model="formData.key"
            label="分类标识"
            :rules="rules.key"
            :disabled="mode === 'edit'"
            variant="outlined"
            density="comfortable"
            class="mb-2"
            hint="如: general, security, notifications"
            persistent-hint
          ></v-text-field>

          <v-text-field
            v-model="formData.label"
            label="显示名称"
            :rules="rules.label"
            variant="outlined"
            density="comfortable"
            class="mb-2"
            hint="如: 常规设置, 安全设置"
          ></v-text-field>

          <v-select
            v-model="formData.icon"
            label="图标"
            :items="iconOptions"
            item-title="title"
            item-value="value"
            variant="outlined"
            density="comfortable"
            class="mb-2"
          >
            <template #prepend-inner>
              <v-icon>{{ formData.icon }}</v-icon>
            </template>
            <template #item="{ props: itemProps, item }">
              <v-list-item v-bind="itemProps">
                <template #prepend>
                  <v-icon>{{ item.value }}</v-icon>
                </template>
              </v-list-item>
            </template>
          </v-select>

          <v-text-field
            v-model.number="formData.order"
            label="排序权重"
            type="number"
            variant="outlined"
            density="comfortable"
            hint="数值越小排序越靠前"
            persistent-hint
          ></v-text-field>
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
