<script setup lang="ts">
/**
 * 全局搜索框组件
 *
 * 功能：
 * - Cmd+K (Mac) / Ctrl+K (Windows) 快捷键
 * - 搜索结果分组（最近、收藏、匹配）
 * - 键盘导航（↑↓ Enter ESC）
 * - 响应式设计（移动端仅显示图标）
 */
import { ref, computed, watch, onMounted, onUnmounted } from "vue";
import { useRouter } from "vue-router";
import { useNavbarStore } from "@/stores";

const router = useRouter();
const navbarStore = useNavbarStore();
const searchQuery = ref("");
const searchDialog = ref(false);
const selectedIndex = ref(0);

/** 检测操作系统 */
const isMac = computed(() => {
  return navigator.platform.toUpperCase().includes("MAC");
});

/** 快捷键显示文本 */
const shortcutText = computed(() => {
  return isMac.value ? "⌘K" : "Ctrl+K";
});

/** 搜索结果 */
const searchResults = computed(() => {
  if (!searchQuery.value.trim()) {
    // 无搜索词时，显示最近访问和收藏
    return {
      recent: navbarStore.recentItems.slice(0, 3),
      favorites: navbarStore.favoriteItems.slice(0, 3),
      matches: [],
    };
  }

  const query = searchQuery.value.toLowerCase();
  const matches = navbarStore.allMenuItems.filter((item) => {
    return (
      item.title.toLowerCase().includes(query) ||
      item.description?.toLowerCase().includes(query) ||
      item.category?.toLowerCase().includes(query) ||
      item.path.toLowerCase().includes(query)
    );
  });

  return {
    recent: [],
    favorites: [],
    matches: matches.slice(0, 8),
  };
});

/** 所有可选项（用于键盘导航） */
const allItems = computed(() => {
  return [...searchResults.value.recent, ...searchResults.value.favorites, ...searchResults.value.matches];
});

/** 打开搜索框 */
function openSearch() {
  searchDialog.value = true;
  searchQuery.value = "";
  selectedIndex.value = 0;
}

/** 关闭搜索框 */
function closeSearch() {
  searchDialog.value = false;
  searchQuery.value = "";
  selectedIndex.value = 0;
}

/** 导航到指定页面 */
function navigateTo(path: string) {
  navbarStore.handleNavigation();
  router.push(path);
  closeSearch();
}

/** 键盘导航 */
function handleKeydown(event: KeyboardEvent) {
  const items = allItems.value;
  if (!items.length) return;

  switch (event.key) {
    case "ArrowDown":
      event.preventDefault();
      selectedIndex.value = Math.min(selectedIndex.value + 1, items.length - 1);
      break;
    case "ArrowUp":
      event.preventDefault();
      selectedIndex.value = Math.max(selectedIndex.value - 1, 0);
      break;
    case "Enter":
      event.preventDefault();
      const selectedItem = items[selectedIndex.value];
      if (selectedItem) {
        navigateTo(selectedItem.path);
      }
      break;
    case "Escape":
      event.preventDefault();
      closeSearch();
      break;
  }
}

/** 重置选中索引 */
watch(searchQuery, () => {
  selectedIndex.value = 0;
});

/** 全局快捷键监听 */
function handleGlobalKeydown(event: KeyboardEvent) {
  // Cmd+K (Mac) 或 Ctrl+K (Windows/Linux)
  if ((event.metaKey || event.ctrlKey) && event.key === "k") {
    event.preventDefault();
    openSearch();
  }
}

onMounted(() => {
  window.addEventListener("keydown", handleGlobalKeydown);
});

onUnmounted(() => {
  window.removeEventListener("keydown", handleGlobalKeydown);
});
</script>

<template>
  <div class="search-box">
    <!-- 搜索框触发器 - 桌面版 -->
    <div class="search-trigger d-none d-sm-flex" @click="openSearch">
      <v-icon size="small" class="search-icon">mdi-magnify</v-icon>
      <span class="search-text">搜索或跳转至...</span>
      <v-chip size="x-small" variant="outlined" class="shortcut-chip">
        {{ shortcutText }}
      </v-chip>
    </div>

    <!-- 搜索图标按钮 - 移动版 -->
    <v-btn icon size="small" variant="text" class="d-sm-none" @click="openSearch">
      <v-icon>mdi-magnify</v-icon>
    </v-btn>

    <!-- 搜索对话框 -->
    <v-dialog v-model="searchDialog" max-width="600" class="search-dialog" @keydown="handleKeydown">
      <v-card>
        <v-card-text class="pa-0">
          <!-- 搜索输入框 -->
          <v-text-field
            v-model="searchQuery"
            placeholder="搜索或跳转至..."
            prepend-inner-icon="mdi-magnify"
            variant="solo"
            flat
            hide-details
            autofocus
            clearable
            class="search-input"
          />

          <!-- 搜索结果 -->
          <div v-if="allItems.length > 0" class="search-results">
            <!-- 最近访问 -->
            <div v-if="searchResults.recent.length > 0" class="result-section">
              <div class="result-header">
                <v-icon size="small" class="mr-2">mdi-history</v-icon>
                最近访问
              </div>
              <v-list density="compact">
                <v-list-item
                  v-for="(item, index) in searchResults.recent"
                  :key="item.path"
                  :prepend-icon="item.icon"
                  :title="item.title"
                  :subtitle="item.category"
                  :active="selectedIndex === index"
                  class="search-result-item"
                  @click="navigateTo(item.path)"
                >
                  <template #append>
                    <v-icon size="small">mdi-arrow-right</v-icon>
                  </template>
                </v-list-item>
              </v-list>
            </div>

            <!-- 收藏夹 -->
            <div v-if="searchResults.favorites.length > 0" class="result-section">
              <div class="result-header">
                <v-icon size="small" class="mr-2">mdi-star</v-icon>
                收藏夹
              </div>
              <v-list density="compact">
                <v-list-item
                  v-for="(item, index) in searchResults.favorites"
                  :key="item.path"
                  :prepend-icon="item.icon"
                  :title="item.title"
                  :subtitle="item.category"
                  :active="selectedIndex === searchResults.recent.length + index"
                  class="search-result-item"
                  @click="navigateTo(item.path)"
                >
                  <template #append>
                    <v-icon size="small">mdi-arrow-right</v-icon>
                  </template>
                </v-list-item>
              </v-list>
            </div>

            <!-- 搜索匹配 -->
            <div v-if="searchResults.matches.length > 0" class="result-section">
              <div class="result-header">
                <v-icon size="small" class="mr-2">mdi-text-search</v-icon>
                搜索结果
              </div>
              <v-list density="compact">
                <v-list-item
                  v-for="(item, index) in searchResults.matches"
                  :key="item.path"
                  :prepend-icon="item.icon"
                  :title="item.title"
                  :subtitle="item.description || item.category"
                  :active="selectedIndex === searchResults.recent.length + searchResults.favorites.length + index"
                  class="search-result-item"
                  @click="navigateTo(item.path)"
                >
                  <template #append>
                    <v-chip size="x-small" variant="outlined">
                      {{ item.category }}
                    </v-chip>
                  </template>
                </v-list-item>
              </v-list>
            </div>
          </div>

          <!-- 空状态 -->
          <div v-else class="empty-state">
            <v-icon size="64" color="grey-lighten-1">mdi-magnify</v-icon>
            <div class="text-h6 text-grey mt-4">未找到匹配结果</div>
            <div class="text-caption text-grey">请尝试其他搜索词</div>
          </div>

          <!-- 提示信息 -->
          <div class="search-tips">
            <v-chip size="x-small" variant="text" prepend-icon="mdi-keyboard"> ↑↓ 导航 </v-chip>
            <v-chip size="x-small" variant="text" prepend-icon="mdi-keyboard-return"> 回车跳转 </v-chip>
            <v-chip size="x-small" variant="text" prepend-icon="mdi-keyboard-esc"> ESC 关闭 </v-chip>
          </div>
        </v-card-text>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.search-box {
  display: flex;
  align-items: center;
  margin-right: 8px;
}

.search-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 6px;
  background-color: rgba(var(--v-theme-surface), 1);
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 200px;
}

.search-trigger:hover {
  background-color: rgba(var(--v-theme-primary), 0.08);
  border-color: rgb(var(--v-theme-primary));
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.search-icon {
  opacity: 0.7;
  flex-shrink: 0;
}

.search-text {
  opacity: 0.8;
  font-size: 0.875rem;
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.shortcut-chip {
  opacity: 0.6;
  font-size: 0.7rem;
  font-weight: 600;
  height: 20px;
  flex-shrink: 0;
}

.search-input {
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}

.search-results {
  max-height: 400px;
  overflow-y: auto;
}

.result-section {
  padding: 8px 0;
}

.result-header {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  font-size: 0.75rem;
  font-weight: 600;
  opacity: 0.7;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.search-result-item {
  cursor: pointer;
  transition: all 0.2s ease;
}

.search-result-item:hover {
  background-color: rgba(var(--v-theme-primary), 0.08);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 16px;
}

.search-tips {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  background-color: rgba(var(--v-theme-surface-variant), 0.3);
}

/* 滚动条样式 */
.search-results::-webkit-scrollbar {
  width: 6px;
}

.search-results::-webkit-scrollbar-track {
  background: transparent;
}

.search-results::-webkit-scrollbar-thumb {
  background: rgba(var(--v-border-color), 0.3);
  border-radius: 3px;
}

.search-results::-webkit-scrollbar-thumb:hover {
  background: rgba(var(--v-border-color), 0.5);
}
</style>
