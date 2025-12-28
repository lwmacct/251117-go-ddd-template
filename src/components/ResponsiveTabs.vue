<!--
  ResponsiveTabs - 响应式 Tabs 布局组件

  根据屏幕宽度自动切换布局：
  - 宽屏 (vertical=true)：左侧垂直 Tabs + 右侧内容
  - 窄屏 (vertical=false)：顶部水平 Tabs + 下方内容

  使用示例：
  <ResponsiveTabs
    :model-value="currentTab"
    :tabs="tabs"
    :vertical="isVertical"
    @update:model-value="handleTabChange"
  >
    <template #tab1>内容 1</template>
    <template #tab2>内容 2</template>
  </ResponsiveTabs>
-->
<script setup lang="ts">
import type { TabItem } from "@/composables/useResponsiveTabs";

interface Props {
  /** 当前激活的 Tab (v-model) */
  modelValue: string;
  /** Tab 项列表 */
  tabs: TabItem[];
  /** 是否垂直布局（宽屏） */
  vertical?: boolean;
  /** Tab 栏背景色 */
  bgColor?: string;
  /** 垂直布局时左侧 Tab 栏宽度 */
  railWidth?: string;
}

withDefaults(defineProps<Props>(), {
  vertical: false,
  bgColor: undefined,
  railWidth: "200px",
});

const emit = defineEmits<{
  "update:modelValue": [value: string];
}>();
</script>

<template>
  <!-- 宽屏：左侧垂直 Tabs -->
  <v-row v-if="vertical" no-gutters>
    <v-col :style="{ minWidth: railWidth, maxWidth: railWidth }">
      <v-tabs
        :model-value="modelValue"
        direction="vertical"
        :bg-color="bgColor"
        show-arrows
        next-icon="mdi-chevron-down"
        prev-icon="mdi-chevron-up"
        @update:model-value="emit('update:modelValue', $event as string)"
      >
        <v-tab v-for="tab in tabs" :key="tab.value" :value="tab.value" :prepend-icon="tab.icon">
          {{ tab.label }}
        </v-tab>
      </v-tabs>
    </v-col>
    <v-divider vertical />
    <v-col>
      <v-tabs-window :model-value="modelValue" class="pa-4">
        <v-tabs-window-item v-for="tab in tabs" :key="tab.value" :value="tab.value">
          <slot :name="tab.value" :tab="tab" />
        </v-tabs-window-item>
      </v-tabs-window>
    </v-col>
  </v-row>

  <!-- 窄屏：顶部水平 Tabs -->
  <div v-else>
    <v-tabs
      :model-value="modelValue"
      :bg-color="bgColor"
      show-arrows
      @update:model-value="emit('update:modelValue', $event as string)"
    >
      <v-tab v-for="tab in tabs" :key="tab.value" :value="tab.value" :prepend-icon="tab.icon">
        {{ tab.label }}
      </v-tab>
    </v-tabs>
    <v-tabs-window :model-value="modelValue" class="pa-4">
      <v-tabs-window-item v-for="tab in tabs" :key="tab.value" :value="tab.value">
        <slot :name="tab.value" :tab="tab" />
      </v-tabs-window-item>
    </v-tabs-window>
  </div>
</template>
