<script setup lang="ts">
import { ref, watch, computed } from "vue";
import type { OrgOrgDTO, OrgCreateOrgDTO, OrgUpdateOrgDTO } from "@models";

interface Props {
  modelValue: boolean;
  org?: OrgOrgDTO | null;
  mode: "create" | "edit";
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
  (e: "save", data: OrgCreateOrgDTO | OrgUpdateOrgDTO): void;
}>();

// 表单数据
const formData = ref<{
  name: string;
  display_name: string;
  description: string;
  avatar: string;
  status: string;
}>({
  name: "",
  display_name: "",
  description: "",
  avatar: "",
  status: "active",
});

// 表单引用
const formRef = ref<HTMLFormElement | null>(null);

// 表单验证规则
const rules = {
  name: [
    (v: string) => !!v || "组织标识不能为空",
    (v: string) => (v && v.length >= 2) || "组织标识至少2个字符",
    (v: string) => /^[a-z0-9-]+$/.test(v) || "组织标识只能包含小写字母、数字和连字符",
  ],
  display_name: [(v: string) => !!v || "组织名称不能为空", (v: string) => (v && v.length >= 2) || "组织名称至少2个字符"],
  avatar: [(v: string) => !v || /^https?:\/\//.test(v) || "头像必须是有效的 URL"],
};

// 监听 org 变化，更新表单数据
watch(
  () => props.org,
  (newOrg) => {
    if (newOrg) {
      formData.value = {
        name: newOrg.name || "",
        display_name: newOrg.display_name || "",
        description: newOrg.description || "",
        avatar: newOrg.avatar || "",
        status: newOrg.status || "active",
      };
    } else {
      formData.value = {
        name: "",
        display_name: "",
        description: "",
        avatar: "",
        status: "active",
      };
    }
  },
  { immediate: true },
);

// 状态选项
const statusOptions = [
  { title: "启用", value: "active" },
  { title: "禁用", value: "suspended" },
];

// 对话框标题
const dialogTitle = computed(() => {
  return props.mode === "create" ? "新建组织" : "编辑组织";
});

/**
 * 保存处理
 */
const handleSave = async () => {
  // 验证表单
  const valid = await formRef.value?.validate?.();
  if (valid?.length > 0) {
    return; // 验证失败
  }

  if (props.mode === "create") {
    emit("save", {
      name: formData.value.name,
      display_name: formData.value.display_name,
      description: formData.value.description || undefined,
      avatar: formData.value.avatar || undefined,
    } as OrgCreateOrgDTO);
  } else {
    emit("save", {
      display_name: formData.value.display_name || undefined,
      description: formData.value.description || undefined,
      avatar: formData.value.avatar || undefined,
      status: formData.value.status || undefined,
    } as OrgUpdateOrgDTO);
  }
  emit("update:modelValue", false);
};

/**
 * 关闭对话框
 */
const handleClose = () => {
  emit("update:modelValue", false);
};
</script>

<template>
  <v-dialog :model-value="modelValue" max-width="600" @click:outside="handleClose">
    <v-card>
      <v-card-title>{{ dialogTitle }}</v-card-title>

      <v-card-text>
        <v-form ref="formRef">
          <v-text-field
            v-model="formData.name"
            label="组织标识"
            :rules="rules.name"
            :disabled="mode === 'edit'"
            variant="outlined"
            required
            class="mb-2"
            hint="创建后不可修改，只能包含小写字母、数字和连字符"
            persistent-hint
          />

          <v-text-field
            v-model="formData.display_name"
            label="组织名称"
            :rules="rules.display_name"
            variant="outlined"
            required
            class="mb-2"
          />

          <v-textarea v-model="formData.description" label="组织描述" variant="outlined" rows="3" auto-grow class="mb-2" />

          <v-text-field
            v-model="formData.avatar"
            label="头像 URL"
            :rules="rules.avatar"
            variant="outlined"
            class="mb-2"
            hint="可选，组织 Logo 图片地址"
            persistent-hint
          />

          <v-select
            v-if="mode === 'edit'"
            v-model="formData.status"
            label="状态"
            :items="statusOptions"
            item-title="title"
            item-value="value"
            variant="outlined"
          />
        </v-form>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="handleClose">取消</v-btn>
        <v-btn color="primary" variant="elevated" @click="handleSave">保存</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
