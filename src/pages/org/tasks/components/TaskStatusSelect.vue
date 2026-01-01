<!--
  TaskStatusSelect.vue - 任务状态选择器
  根据当前状态显示可流转的下一状态选项
-->
<script setup lang="ts">
import { computed } from "vue";
import type { TaskStatusValue } from "../composables/useTeamTasks";
import { TaskStatusConfig, StatusTransitions } from "../composables/useTeamTasks";

interface Props {
  modelValue: TaskStatusValue;
  disabled?: boolean;
}

interface Emits {
  (e: "update:modelValue", value: TaskStatusValue): void;
}

interface StatusOption {
  title: string;
  value: TaskStatusValue;
  color: string;
  icon: string;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

/**
 * 获取可用的状态选项（基于当前状态的流转规则）
 */
const availableStatuses = computed<StatusOption[]>(() => {
  const transitions = StatusTransitions[props.modelValue] || [];
  return transitions.map((status) => ({
    title: TaskStatusConfig[status].label,
    value: status,
    color: TaskStatusConfig[status].color,
    icon: TaskStatusConfig[status].icon,
  }));
});

const internalValue = computed({
  get: () => props.modelValue,
  set: (value: TaskStatusValue) => emit("update:modelValue", value),
});
</script>

<template>
  <v-select
    v-model="internalValue"
    :items="availableStatuses"
    :disabled="disabled"
    variant="outlined"
    density="compact"
    hide-details
    prepend-inner-icon="mdi-flag"
  >
    <template #selection="{ item }">
      <v-chip :color="(item as unknown as StatusOption).color" size="small" density="comfortable">
        <v-icon :icon="(item as unknown as StatusOption).icon" start size="small" />
        {{ (item as unknown as StatusOption).title }}
      </v-chip>
    </template>
    <template #item="{ item, props: itemProps }">
      <v-list-item v-bind="itemProps">
        <template #prepend>
          <v-icon :icon="item.raw.icon" :color="item.raw.color" />
        </template>
        <v-list-item-title>{{ item.raw.title }}</v-list-item-title>
      </v-list-item>
    </template>
  </v-select>
</template>
