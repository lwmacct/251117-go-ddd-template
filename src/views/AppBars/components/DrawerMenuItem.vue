<script setup lang="ts">
/**
 * 抽屉菜单项组件（通用）
 *
 * 合并自 RecentPagesMenuItem + AllPagesMenuItem
 * 用于抽屉菜单中的入口项，悬停时显示对应面板
 */
import type { HoverPanelType } from "../types";

interface Props {
  /** 图标 */
  icon: string;
  /** 标题 */
  title: string;
  /** 面板类型 */
  panelType: Exclude<HoverPanelType, null>;
  /** 数量（可选，不传则不显示徽标） */
  count?: number;
  /** 当前悬停的面板 */
  hoveredItem: HoverPanelType;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  mouseEnter: [panelType: HoverPanelType];
}>();

function handleMouseEnter() {
  emit("mouseEnter", props.panelType);
}
</script>

<template>
  <v-list-item
    :prepend-icon="icon"
    :title="title"
    :class="{ 'menu-item-active': hoveredItem === panelType }"
    class="drawer-menu-item"
    @mouseenter="handleMouseEnter"
  >
    <template #append>
      <v-chip v-if="count && count > 0" size="x-small" variant="tonal">
        {{ count }}
      </v-chip>
      <v-icon size="small" class="chevron-icon">mdi-chevron-right</v-icon>
    </template>
  </v-list-item>
</template>

<style scoped>
.drawer-menu-item {
  border-radius: 4px;
  margin: 2px 8px;
  transition: all 0.2s ease;
}

.drawer-menu-item:hover,
.menu-item-active {
  background-color: rgba(var(--v-theme-primary), 0.08);
}

.chevron-icon {
  opacity: 0.5;
  transition: transform 0.2s ease;
}

.drawer-menu-item:hover .chevron-icon,
.menu-item-active .chevron-icon {
  opacity: 1;
  transform: translateX(2px);
}
</style>
