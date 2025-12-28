<script setup lang="ts">
/**
 * 所有页面面板
 *
 * 功能：
 * - 按分类分组展示所有页面
 * - 收藏星标切换
 * - 滚动容器
 */
import { useNavbarStore } from "@/stores";
import FavoriteButton from "./FavoriteButton.vue";

const emit = defineEmits<{
  navigate: [path: string];
}>();

const navbarStore = useNavbarStore();

/** 处理点击 */
function handleClick(path: string) {
  emit("navigate", path);
}
</script>

<template>
  <div class="products-panel">
    <!-- 头部 -->
    <div class="panel-header">
      <div class="header-title">
        <v-icon size="small" class="mr-2">mdi-apps</v-icon>
        <span>所有页面</span>
        <v-chip size="x-small" variant="tonal" class="ml-2">
          {{ navbarStore.allMenuItems.length }}
        </v-chip>
      </div>
    </div>

    <!-- 分类列表 -->
    <div class="categories-list">
      <div v-for="(items, category) in navbarStore.menuItemsByCategory" :key="category" class="category-group">
        <!-- 分类标题 -->
        <div class="category-header">
          <span>{{ category }}</span>
          <v-chip size="x-small" variant="text">
            {{ items.length }}
          </v-chip>
        </div>

        <!-- 页面列表 -->
        <v-list density="compact" nav>
          <v-list-item
            v-for="item in items"
            :key="item.path"
            :prepend-icon="item.icon"
            :title="item.title"
            :subtitle="item.description"
            class="page-item"
            @click="handleClick(item.path)"
          >
            <template #append>
              <FavoriteButton :path="item.path" :is-favorite="item.isFavorite" @toggle="navbarStore.toggleFavorite($event)" />
            </template>
          </v-list-item>
        </v-list>
      </div>
    </div>
  </div>
</template>

<style scoped>
.products-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-width: 300px;
  max-width: 400px;
}

.panel-header {
  padding: 12px 16px;
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}

.header-title {
  display: flex;
  align-items: center;
  font-weight: 600;
}

.categories-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.category-group {
  margin-bottom: 8px;
}

.category-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px 4px;
  font-size: 0.75rem;
  font-weight: 600;
  opacity: 0.7;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.page-item {
  border-radius: 4px;
  margin: 2px 8px;
  transition: all 0.2s ease;
}

.page-item:hover {
  background-color: rgba(var(--v-theme-primary), 0.08);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

/* 滚动条样式 */
.categories-list::-webkit-scrollbar {
  width: 4px;
}

.categories-list::-webkit-scrollbar-track {
  background: transparent;
}

.categories-list::-webkit-scrollbar-thumb {
  background: rgba(var(--v-border-color), 0.3);
  border-radius: 2px;
}
</style>
