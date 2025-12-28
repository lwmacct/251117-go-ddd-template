<script setup lang="ts">
/**
 * 最近访问面板
 *
 * 功能：
 * - 显示最近 10 条访问记录
 * - 时间格式化（刚刚、X分钟前、HH:MM 等）
 * - 清空历史按钮
 * - 收藏切换
 */
import { computed } from "vue";
import { useNavbarStore } from "@/stores";
import FavoriteButton from "./FavoriteButton.vue";

const emit = defineEmits<{
  navigate: [path: string];
}>();

const navbarStore = useNavbarStore();

/** 最近访问记录 */
const recentRecords = computed(() => {
  return navbarStore.accessHistory.slice(0, 10);
});

/** 格式化时间 */
function formatTime(timestamp: number): string {
  const now = Date.now();
  const diff = now - timestamp;
  const minutes = Math.floor(diff / 60000);
  const hours = Math.floor(diff / 3600000);

  if (minutes < 1) return "刚刚";
  if (minutes < 60) return `${minutes}分钟前`;
  if (hours < 24) {
    const date = new Date(timestamp);
    return `${date.getHours().toString().padStart(2, "0")}:${date.getMinutes().toString().padStart(2, "0")}`;
  }

  const date = new Date(timestamp);
  const today = new Date();
  const yesterday = new Date(today);
  yesterday.setDate(yesterday.getDate() - 1);

  if (date.toDateString() === yesterday.toDateString()) {
    return "昨天";
  }

  return `${(date.getMonth() + 1).toString().padStart(2, "0")}-${date.getDate().toString().padStart(2, "0")} ${date.getHours().toString().padStart(2, "0")}:${date.getMinutes().toString().padStart(2, "0")}`;
}

/** 处理点击 */
function handleClick(path: string) {
  emit("navigate", path);
}

/** 清空历史 */
function handleClearHistory() {
  navbarStore.clearHistory();
}
</script>

<template>
  <div class="recent-pages-panel">
    <!-- 头部 -->
    <div class="panel-header">
      <div class="header-title">
        <v-icon size="small" class="mr-2">mdi-history</v-icon>
        <span>最近访问</span>
        <v-chip v-if="recentRecords.length > 0" size="x-small" variant="tonal" class="ml-2">
          {{ recentRecords.length }}
        </v-chip>
      </div>
    </div>

    <!-- 记录列表 -->
    <div v-if="recentRecords.length > 0" class="records-list">
      <v-list density="compact" nav>
        <v-list-item
          v-for="record in recentRecords"
          :key="record.path"
          :prepend-icon="record.icon"
          class="record-item"
          @click="handleClick(record.path)"
        >
          <template #title>
            <div class="record-title">{{ record.title }}</div>
          </template>
          <template #subtitle>
            <div class="record-meta">
              <v-chip v-if="record.category" size="x-small" variant="text" class="category-chip">
                {{ record.category }}
              </v-chip>
              <span class="time-text">{{ formatTime(record.timestamp) }}</span>
            </div>
          </template>
          <template #append>
            <FavoriteButton :path="record.path" :is-favorite="record.isFavorite" @toggle="navbarStore.toggleFavorite($event)" />
          </template>
        </v-list-item>
      </v-list>
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <v-icon size="48" color="grey-lighten-1">mdi-history</v-icon>
      <div class="text-body-2 text-grey mt-2">暂无访问记录</div>
    </div>

    <!-- 底部操作 -->
    <div v-if="recentRecords.length > 0" class="panel-footer">
      <v-btn size="small" variant="text" prepend-icon="mdi-delete-outline" @click="handleClearHistory"> 清空历史 </v-btn>
    </div>
  </div>
</template>

<style scoped>
.recent-pages-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-width: 280px;
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

.records-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.record-item {
  border-radius: 4px;
  margin: 2px 8px;
  transition: all 0.2s ease;
}

.record-item:hover {
  background-color: rgba(var(--v-theme-primary), 0.08);
}

.record-title {
  font-size: 0.875rem;
}

.record-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.75rem;
  opacity: 0.7;
}

.category-chip {
  height: 18px;
  font-size: 0.65rem;
}

.time-text {
  font-size: 0.7rem;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  padding: 32px;
}

.panel-footer {
  padding: 8px 16px;
  border-top: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  display: flex;
  justify-content: center;
}

/* 滚动条样式 */
.records-list::-webkit-scrollbar {
  width: 4px;
}

.records-list::-webkit-scrollbar-track {
  background: transparent;
}

.records-list::-webkit-scrollbar-thumb {
  background: rgba(var(--v-border-color), 0.3);
  border-radius: 2px;
}
</style>
