<script setup lang="ts">
import { ref, watch, computed } from "vue";
import type { OrgTeamDTO, OrgCreateTeamDTO, OrgUpdateTeamDTO } from "@models";

interface Props {
  modelValue: boolean;
  team?: OrgTeamDTO | null;
  mode: "create" | "edit";
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
  (e: "save", data: OrgCreateTeamDTO | OrgUpdateTeamDTO): void;
}>();

// 表单数据
const formData = ref<{
  name: string;
  display_name: string;
  description: string;
}>({
  name: "",
  display_name: "",
  description: "",
});

// 表单引用
const formRef = ref<HTMLFormElement | null>(null);

// 表单验证规则
const rules = {
  name: [
    (v: string) => !!v || "团队标识不能为空",
    (v: string) => (v && v.length >= 2) || "团队标识至少2个字符",
    (v: string) => /^[a-z0-9-]+$/.test(v) || "团队标识只能包含小写字母、数字和连字符",
  ],
  display_name: [(v: string) => !!v || "团队名称不能为空", (v: string) => (v && v.length >= 2) || "团队名称至少2个字符"],
};

// 监听 team 变化，更新表单数据
watch(
  () => props.team,
  (newTeam) => {
    if (newTeam) {
      formData.value = {
        name: newTeam.name || "",
        display_name: newTeam.display_name || "",
        description: newTeam.description || "",
      };
    } else {
      formData.value = {
        name: "",
        display_name: "",
        description: "",
      };
    }
  },
  { immediate: true },
);

// 对话框标题
const dialogTitle = computed(() => {
  return props.mode === "create" ? "新建团队" : "编辑团队";
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
    } as OrgCreateTeamDTO);
  } else {
    emit("save", {
      display_name: formData.value.display_name || undefined,
      description: formData.value.description || undefined,
    } as OrgUpdateTeamDTO);
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
            label="团队标识"
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
            label="团队名称"
            :rules="rules.display_name"
            variant="outlined"
            required
            class="mb-2"
          />

          <v-textarea v-model="formData.description" label="团队描述" variant="outlined" rows="3" auto-grow class="mb-2" />
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
