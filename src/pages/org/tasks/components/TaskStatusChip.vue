<!--
  TaskStatusChip.vue - 任务状态显示组件
  根据状态显示不同颜色的徽章
-->
<script setup lang="ts">
import { computed } from "vue";
import type { TaskStatusValue } from "../composables/useTeamTasks";
import { TaskStatusConfig } from "../composables/useTeamTasks";

interface Props {
  status?: TaskStatusValue | string;
}

const props = defineProps<Props>();

const config = computed(() => {
  const s = props.status as TaskStatusValue;
  return TaskStatusConfig[s] || { label: "未知", color: "grey", icon: "mdi-help" };
});
</script>

<template>
  <v-chip :color="config.color" size="small">
    <v-icon :icon="config.icon" start size="small" />
    {{ config.label }}
  </v-chip>
</template>
