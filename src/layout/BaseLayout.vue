<!--
  BaseLayout.vue - 共享基础布局组件
  用于 AdminLayout 和 UserLayout 的公共结构
-->
<script setup lang="ts">
import AppBars from "@/views/AppBars/index.vue";
import Navigation from "@/views/Navigation/index.vue";
import AppBreadcrumb from "@/components/AppBreadcrumb.vue";
import PageTransition from "@/components/PageTransition.vue";
import type { MenuItem } from "@/views/Navigation/types";

interface Props {
  /** 菜单项列表 */
  menuItems: MenuItem[];
  /** 是否显示面包屑导航 */
  showBreadcrumb?: boolean;
  /** 是否启用页面过渡动画 */
  enableTransition?: boolean;
  /** 过渡动画类型 */
  transitionName?: "fade" | "slide-left" | "slide-right" | "slide-up" | "scale" | "none";
  /** 导航抽屉宽度 */
  navWidth?: number;
  /** 是否显示导航图标 */
  showNavIcon?: boolean;
}

withDefaults(defineProps<Props>(), {
  showBreadcrumb: true,
  enableTransition: true,
  transitionName: "fade",
  navWidth: 200,
  showNavIcon: true,
});
</script>

<template>
  <!-- 顶部导航栏 -->
  <AppBars />

  <!-- 主布局：左侧菜单 + 内容区域 -->
  <v-main>
    <!-- 左侧菜单 -->
    <Navigation :items="menuItems" :show-icon="showNavIcon" :width="navWidth" />

    <!-- 主内容区域 -->
    <v-container fluid class="pa-6">
      <!-- 面包屑导航 (可选) -->
      <AppBreadcrumb v-if="showBreadcrumb" class="mb-4" />

      <!-- 页面内容（带过渡动画） -->
      <router-view v-if="enableTransition" v-slot="{ Component }">
        <PageTransition :name="transitionName">
          <component :is="Component" />
        </PageTransition>
      </router-view>

      <!-- 页面内容（无过渡动画） -->
      <router-view v-else />
    </v-container>
  </v-main>
</template>
