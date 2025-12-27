<script setup lang="ts">
import { useThemeToggle } from "@/composables/useThemeToggle";

const { isDark, themeIcon, toggleTheme } = useThemeToggle();
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
      <v-icon :icon="themeIcon" />
    </v-btn>
  </v-app-bar>
</template>
