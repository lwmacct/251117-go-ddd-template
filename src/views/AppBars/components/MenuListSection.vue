<script setup lang="ts">
/**
 * 可复用菜单列表组件
 *
 * 功能：
 * - 渲染菜单项列表
 * - 支持拖拽排序（收藏夹）
 * - 空状态提示
 * - 收藏星标切换
 */
import { computed } from "vue";
import type { MenuItem } from "../types";
import FavoriteButton from "./FavoriteButton.vue";

interface Props {
  /** 标题 */
  title: string;
  /** 菜单项列表 */
  items: MenuItem[];
  /** 空状态提示文本 */
  emptyText?: string;
  /** 是否支持拖拽 */
  draggable?: boolean;
  /** 最大显示数量 */
  maxItems?: number;
}

const props = withDefaults(defineProps<Props>(), {
  emptyText: "暂无数据",
  draggable: false,
  maxItems: 10,
});

const emit = defineEmits<{
  navigate: [path: string];
  toggleFavorite: [path: string];
  reorder: [fromIndex: number, toIndex: number];
}>();

/** 显示的菜单项（限制数量） */
const displayItems = computed(() => {
  return props.items.slice(0, props.maxItems);
});

/** 处理点击 */
function handleClick(path: string) {
  emit("navigate", path);
}
</script>

<template>
  <div class="menu-list-section">
    <!-- 标题 -->
    <div class="section-header">
      <span class="section-title">{{ title }}</span>
      <v-chip v-if="items.length > 0" size="x-small" variant="text">
        {{ items.length }}
      </v-chip>
    </div>

    <!-- 菜单列表 -->
    <v-list v-if="displayItems.length > 0" density="compact" nav>
      <v-list-item
        v-for="item in displayItems"
        :key="item.path"
        :prepend-icon="item.icon"
        :title="item.title"
        class="menu-item"
        @click="handleClick(item.path)"
      >
        <template #append>
          <FavoriteButton :path="item.path" :is-favorite="item.isFavorite" @toggle="emit('toggleFavorite', $event)" />
        </template>
      </v-list-item>
    </v-list>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <v-chip size="small" variant="text" prepend-icon="mdi-information-outline">
        {{ emptyText }}
      </v-chip>
    </div>
  </div>
</template>

<style scoped>
.menu-list-section {
  padding: 4px 0;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px 4px;
}

.section-title {
  font-size: 0.75rem;
  font-weight: 600;
  opacity: 0.7;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.menu-item {
  border-radius: 4px;
  margin: 2px 8px;
  transition: all 0.2s ease;
}

.menu-item:hover {
  background-color: rgba(var(--v-theme-primary), 0.08);
}

.empty-state {
  display: flex;
  justify-content: center;
  padding: 16px;
}
</style>
