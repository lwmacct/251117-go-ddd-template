<script setup lang="ts">
import { computed } from "vue";
import { useTheme } from "vuetify";

const theme = useTheme();

/** 获取当前主题名称 (最可靠的状态源)  */
const currentThemeName = computed(() => theme.global.name.value);

/** 判断是否为暗色主题 (基于主题名称)  */
const isDark = computed(() => currentThemeName.value === "dark");

/** 主题切换函数 - 使用 Vuetify 3.9+ 新 API */
const toggleTheme = () => {
  // 基于当前主题名称直接切换 (避免依赖 current.value.dark 的延迟更新)
  const newTheme = currentThemeName.value === "dark" ? "light" : "dark";

  // 使用 Vuetify 3.9+ 推荐的 change API
  theme.change(newTheme);

  // 保存到 localStorage
  localStorage.setItem("theme", newTheme);
};
</script>

<template>
  <v-app-bar color="primary" density="compact" flat>
    <v-app-bar-title class="d-flex align-center">
      <v-icon icon="mdi-cloud-outline" size="24" class="mr-2" />
      <span class="text-h6">DDD Template</span>
    </v-app-bar-title>

    <v-spacer />

    <!-- 主题切换 -->
    <v-btn icon size="small" title="切换主题" @click="toggleTheme">
      <v-icon :icon="isDark ? 'mdi-weather-sunny' : 'mdi-weather-night'" />
    </v-btn>
  </v-app-bar>
</template>
