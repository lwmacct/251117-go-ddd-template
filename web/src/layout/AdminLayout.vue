<script setup lang="ts">
import AppBars from "@/views/AppBars/index.vue";
import Navigation from "@/views/Navigation/index.vue";
import type { MenuItem } from "@/views/Navigation/types";

/**
 * 管理后台菜单配置
 */
const menuItems: MenuItem[] = [
  {
    title: "数据概览",
    path: "/admin/overview",
    icon: "mdi-speedometer",
  },
  {
    title: "角色管理",
    path: "/admin/roles",
    icon: "mdi-account-group",
  },
  {
    title: "用户管理",
    path: "/admin/users",
    icon: "mdi-account",
  },
  {
    title: "菜单管理",
    path: "/admin/menus",
    icon: "mdi-menu",
  },
  {
    title: "系统设置",
    path: "/admin/settings",
    icon: "mdi-cog",
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
