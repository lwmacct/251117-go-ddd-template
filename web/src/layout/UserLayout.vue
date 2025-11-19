<script setup lang="ts">
import AppBars from "@/views/AppBars/index.vue";
import Navigation from "@/views/Navigation/index.vue";
import type { MenuItem } from "@/views/Navigation/types";

/**
 * 用户中心菜单配置
 */
const menuItems: MenuItem[] = [
  {
    title: "个人资料",
    path: "/user/profile",
    icon: "mdi-account-circle",
  },
  {
    title: "安全设置",
    path: "/user/security",
    icon: "mdi-shield-lock",
  },
  {
    title: "访问令牌",
    path: "/user/tokens",
    icon: "mdi-key-variant",
  },
];
</script>

<template>
  <!-- 顶部导航栏 -->
  <AppBars />
  <!-- 主布局：左侧菜单 + 内容区域 -->
  <v-main class="kso-main">
    <!-- 左侧菜单 -->
    <Navigation :items="menuItems" :show-icon="true" :width="200" />
    <!-- 主内容区域 -->
    <v-container fluid class="kso-content pa-6">
      <router-view />
    </v-container>
  </v-main>
</template>

<style scoped></style>
