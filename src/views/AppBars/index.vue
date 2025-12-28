<script setup lang="ts">
/**
 * AppBars 主组件
 *
 * 功能：
 * - 顶部导航栏（搜索、主题切换、用户菜单）
 * - 抽屉菜单（最近访问、所有页面、收藏夹）
 * - 悬停面板（HoverPanel）
 *
 * 优化点（相对 SaaS 模板）：
 * - 悬停面板增加延迟防抖，避免意外关闭
 * - 移动端使用点击模式
 */
import { ref, watch, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useNavbarStore } from "@/stores";
import type { AppBarProps, HoverPanelType } from "./types";

// 子组件
import UserMenu from "./components/UserMenu.vue";
import SearchBox from "./components/SearchBox.vue";
import ThemeToggle from "./components/ThemeToggle.vue";
import MenuListSection from "./components/MenuListSection.vue";
import RecentPagesPanel from "./components/RecentPagesPanel.vue";
import ProductsPanel from "./components/ProductsPanel.vue";
import DrawerMenuItem from "./components/DrawerMenuItem.vue";

// Props
const props = withDefaults(defineProps<AppBarProps>(), {
  logo: undefined,
  logoIcon: "mdi-cube-outline",
  appName: undefined,
  showNavIcon: true,
  navIcon: "mdi-menu",
  showDrawer: true,
  drawerWidth: 220,
  elevation: 2,
  color: undefined,
  height: 56,
});

// 禁用属性继承（多根节点）
defineOptions({
  inheritAttrs: false,
});

const router = useRouter();
const navbarStore = useNavbarStore();

// 状态
const drawer = ref(false);
const hoveredItem = ref<HoverPanelType>(null);
const hoverTimeout = ref<ReturnType<typeof setTimeout> | null>(null);

// 初始化
onMounted(() => {
  // 同步抽屉状态
  navbarStore.drawerOpen = drawer.value;
});

// 监听抽屉关闭时隐藏面板
watch(drawer, (newValue) => {
  if (!newValue) {
    hoveredItem.value = null;
  }
  navbarStore.drawerOpen = newValue;
});

/** 导航 */
function navigateTo(path: string) {
  router.push(path);
  drawer.value = false;
  hoveredItem.value = null;
}

/** 跳转首页 */
function navigateToHome() {
  router.push("/");
}

/** 处理导航图标点击 */
function handleNavIconClick() {
  if (props.showDrawer) {
    drawer.value = !drawer.value;
  }
}

/** 处理悬停进入（带延迟防抖） */
function handleMouseEnter(panelType: HoverPanelType) {
  // 清除之前的延迟
  if (hoverTimeout.value) {
    clearTimeout(hoverTimeout.value);
    hoverTimeout.value = null;
  }
  hoveredItem.value = panelType;
}

/** 处理悬停离开（带延迟） */
function handleMouseLeave() {
  // 200ms 延迟，避免意外关闭
  hoverTimeout.value = setTimeout(() => {
    hoveredItem.value = null;
  }, 200);
}

/** 保持面板打开 */
function keepPanelOpen() {
  if (hoverTimeout.value) {
    clearTimeout(hoverTimeout.value);
    hoverTimeout.value = null;
  }
}

/** 切换收藏 */
function handleToggleFavorite(path: string) {
  navbarStore.toggleFavorite(path);
}
</script>

<template>
  <!-- 头部导航栏 -->
  <v-app-bar app :elevation="elevation" :color="color" :height="height" class="app-bar">
    <!-- 导航图标 -->
    <v-app-bar-nav-icon v-if="showNavIcon" class="nav-icon" @click="handleNavIconClick">
      <v-icon>{{ drawer ? "mdi-close" : navIcon }}</v-icon>
    </v-app-bar-nav-icon>

    <!-- 分割线 -->
    <v-divider v-if="showNavIcon" vertical class="mr-3" />

    <!-- Logo 区域 -->
    <div class="logo-area" @click="navigateToHome">
      <img v-if="logo" :src="logo" alt="Logo" class="logo-image" />
      <v-icon v-else size="28">{{ logoIcon }}</v-icon>
      <span v-if="appName" class="app-name">{{ appName }}</span>
    </div>

    <!-- 右侧区域 -->
    <v-spacer />
    <SearchBox />
    <ThemeToggle />
    <UserMenu />
  </v-app-bar>

  <!-- 抽屉菜单 -->
  <v-navigation-drawer v-if="showDrawer" v-model="drawer" app temporary :width="drawerWidth" class="app-drawer">
    <v-list nav>
      <!-- 最近访问 -->
      <DrawerMenuItem
        icon="mdi-history"
        title="最近访问"
        panel-type="recent"
        :hovered-item="hoveredItem"
        @mouse-enter="handleMouseEnter"
      />

      <!-- 所有页面 -->
      <DrawerMenuItem
        icon="mdi-apps"
        title="所有页面"
        panel-type="all"
        :hovered-item="hoveredItem"
        @mouse-enter="handleMouseEnter"
      />

      <v-divider class="my-2" />

      <!-- 收藏列表 -->
      <MenuListSection
        title="收藏页面"
        :items="navbarStore.favoriteItems"
        empty-text="暂无收藏页面"
        @navigate="navigateTo"
        @toggle-favorite="handleToggleFavorite"
      />
    </v-list>

    <!-- 悬停面板 -->
    <div v-if="hoveredItem && drawer" class="hover-panel" @mouseenter="keepPanelOpen" @mouseleave="handleMouseLeave">
      <!-- 连接区域 -->
      <div class="connection-area" />

      <!-- 最近访问面板 -->
      <RecentPagesPanel v-if="hoveredItem === 'recent'" @navigate="navigateTo" />

      <!-- 所有页面面板 -->
      <ProductsPanel v-if="hoveredItem === 'all'" @navigate="navigateTo" />
    </div>
  </v-navigation-drawer>
</template>

<style scoped>
.app-bar {
  z-index: 1000;
}

.nav-icon {
  margin: 4px;
}

/* Logo 区域 */
.logo-area {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.logo-area:hover {
  background-color: rgba(var(--v-theme-on-surface), 0.04);
}

.logo-image {
  height: 32px;
  width: auto;
}

.app-name {
  margin-left: 8px;
  font-size: 1.1rem;
  font-weight: 500;
}

.app-drawer {
  z-index: 999;
}

/* 悬停面板 */
.hover-panel {
  position: absolute;
  left: 100%;
  top: 0;
  min-width: 280px;
  max-width: 400px;
  height: 100%;
  background-color: rgb(var(--v-theme-surface));
  border-left: thin solid rgba(var(--v-border-color), var(--v-border-opacity));
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
  overflow: visible;
  z-index: 1;
}

/* 连接区域：确保鼠标可以从抽屉移动到面板 */
.connection-area {
  position: absolute;
  left: -20px;
  top: 0;
  width: 20px;
  height: 100%;
}
</style>
