<!--
  TaskDialog.vue - 创建/编辑任务对话框
-->
<script setup lang="ts">
import { ref, watch, computed } from "vue";
import type { TaskTaskDTO, TaskCreateTaskDTO, TaskUpdateTaskDTO } from "@models";

interface Props {
  modelValue: boolean;
  task?: TaskTaskDTO | null;
  members: Array<{ title: string; value: number }>;
  loading?: boolean;
}

interface Emits {
  (e: "update:modelValue", value: boolean): void;
  (e: "save", data: TaskCreateTaskDTO | TaskUpdateTaskDTO): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

// 表单状态
const form = ref<TaskCreateTaskDTO | TaskUpdateTaskDTO>({
  title: "",
  description: "",
  assignee_id: undefined,
});

const isEditing = computed(() => !!props.task?.id);

// 对话框标题
const dialogTitle = computed(() => (isEditing.value ? "编辑任务" : "新建任务"));

// 重置表单
const resetForm = () => {
  form.value = {
    title: "",
    description: "",
    assignee_id: undefined,
  };
};

// 监听对话框打开，初始化表单
watch(
  () => props.modelValue,
  (open) => {
    if (open && props.task) {
      // 编辑模式：填充表单
      form.value = {
        title: props.task.title || "",
        description: props.task.description || "",
        assignee_id: props.task.assignee_id,
      };
    } else if (!open) {
      // 关闭时重置
      resetForm();
    }
  },
);

// 保存处理
const handleSave = () => {
  if (!form.value.title?.trim()) {
    return;
  }
  emit("save", form.value);
};

// 取消处理
const handleCancel = () => {
  emit("update:modelValue", false);
};
</script>

<template>
  <v-dialog :model-value="modelValue" max-width="600" @update:model-value="handleCancel">
    <v-card>
      <v-card-title class="d-flex align-center py-4">
        <v-icon icon="mdi-pencil" start />
        {{ dialogTitle }}
      </v-card-title>

      <v-card-text class="pt-4">
        <v-form @submit.prevent="handleSave">
          <!-- 标题 -->
          <v-text-field
            v-model="form.title"
            label="任务标题"
            placeholder="输入任务标题..."
            variant="outlined"
            counter="200"
            :rules="[(v) => !!v?.trim() || '标题不能为空']"
            required
            autofocus
          />

          <!-- 描述 -->
          <v-textarea
            v-model="form.description"
            label="任务描述"
            placeholder="添加详细描述..."
            variant="outlined"
            counter="2000"
            rows="4"
            auto-grow
            class="mt-4"
          />

          <!-- 指派成员 -->
          <v-select
            v-model="form.assignee_id"
            :items="members"
            label="指派给"
            placeholder="未指派"
            variant="outlined"
            clearable
            hide-details
            class="mt-4"
            prepend-inner-icon="mdi-account"
          >
            <template #no-data>
              <v-list-item>
                <v-list-item-title class="text-medium-emphasis"> 暂无成员可选 </v-list-item-title>
              </v-list-item>
            </template>
          </v-select>
        </v-form>
      </v-card-text>

      <v-card-actions class="pa-4 pt-0">
        <v-spacer />
        <v-btn variant="text" @click="handleCancel">取消</v-btn>
        <v-btn color="primary" variant="elevated" :loading="loading" :disabled="!form.title?.trim()" @click="handleSave">
          保存
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
