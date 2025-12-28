<script setup lang="ts">
/**
 * 主题切换组件
 *
 * 功能：
 * - 切换亮色/暗色主题
 * - 图标动画效果
 * - 状态持久化（通过 useThemeToggle）
 */
import { useThemeToggle } from "@/composables/useThemeToggle";

const { themeIcon, toggleTheme } = useThemeToggle();
</script>

<template>
  <v-btn icon size="small" variant="text" title="切换主题" class="theme-toggle-btn" @click="toggleTheme">
    <v-icon :icon="themeIcon" class="theme-icon" />
  </v-btn>
</template>

<style scoped>
.theme-toggle-btn {
  transition: transform 0.2s ease;
}

.theme-toggle-btn:hover {
  transform: scale(1.1);
}

.theme-icon {
  transition: transform 0.3s ease;
}

.theme-toggle-btn:hover .theme-icon {
  transform: rotate(30deg);
}
</style>
