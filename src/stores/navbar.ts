/**
 * Navbar 状态管理 Store
 *
 * 聚合导航相关的全局状态：
 * - 收藏夹（localStorage 存储）
 * - 访问历史（委托给 useAccessHistory，带 localStorage 持久化）
 * - 抽屉状态
 *
 * 注意：菜单项现在由 useMenus composable 动态生成，基于角色过滤
 */

import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { useMenus } from "@/composables";
import { useAccessHistory } from "@/views/AppBars/composables";
import type { MenuItem } from "@/views/AppBars/types";

/** 收藏夹存储 Key */
const FAVORITES_KEY = "navbar_favorites";

/**
 * 从 localStorage 加载收藏夹
 */
function loadFavorites(): string[] {
  try {
    const data = localStorage.getItem(FAVORITES_KEY);
    return data ? JSON.parse(data) : [];
  } catch {
    return [];
  }
}

/**
 * 保存收藏夹到 localStorage
 */
function saveFavorites(paths: string[]): void {
  try {
    localStorage.setItem(FAVORITES_KEY, JSON.stringify(paths));
  } catch {
    // 静默失败
  }
}

export const useNavbarStore = defineStore("navbar", () => {
  // ==================== 动态菜单（基于角色） ====================

  const { adminMenus, userMenus } = useMenus();

  // ==================== 访问历史（委托给 composable）====================

  // 使用 composable 管理访问历史（带 localStorage 持久化）
  // userId 暂时使用 null，后续可从 authStore 获取
  const historyManager = useAccessHistory(null);

  // ==================== 状态 ====================

  /** 收藏的路径列表 */
  const favoritePaths = ref<string[]>(loadFavorites());

  /** 抽屉状态 */
  const drawerOpen = ref(false);

  // ==================== 计算属性 ====================

  /**
   * 访问历史（从 composable 获取）
   */
  const accessHistory = historyManager.recentPages;

  /**
   * 所有菜单项（合并管理员和用户菜单，基于角色动态过滤）
   */
  const allMenuItems = computed<MenuItem[]>(() => {
    const items: MenuItem[] = [];

    // 转换管理员菜单（已经按角色过滤）
    adminMenus.value.forEach((item, index) => {
      items.push({
        id: `admin-${index}`,
        title: item.title,
        icon: item.icon || "mdi-file",
        path: item.path,
        category: "管理后台",
        isFavorite: favoritePaths.value.includes(item.path),
      });
    });

    // 转换用户菜单（已经按角色过滤）
    userMenus.value.forEach((item, index) => {
      items.push({
        id: `user-${index}`,
        title: item.title,
        icon: item.icon || "mdi-file",
        path: item.path,
        category: "用户中心",
        isFavorite: favoritePaths.value.includes(item.path),
      });
    });

    return items;
  });

  /**
   * 收藏的菜单项
   */
  const favoriteItems = computed<MenuItem[]>(() => {
    return allMenuItems.value.filter((item) => favoritePaths.value.includes(item.path));
  });

  /**
   * 最近访问的菜单项
   */
  const recentItems = computed<MenuItem[]>(() => {
    return accessHistory.value.slice(0, 10).map((record) => ({
      id: record.path,
      title: record.title,
      icon: record.icon,
      path: record.path,
      category: record.category,
      isFavorite: favoritePaths.value.includes(record.path),
    }));
  });

  /**
   * 按分类分组的菜单项
   */
  const menuItemsByCategory = computed(() => {
    const groups: Record<string, MenuItem[]> = {};
    allMenuItems.value.forEach((item) => {
      const category = item.category || "其他";
      if (!groups[category]) {
        groups[category] = [];
      }
      groups[category].push(item);
    });
    return groups;
  });

  // ==================== 方法 ====================

  /**
   * 切换收藏状态
   */
  function toggleFavorite(path: string) {
    const index = favoritePaths.value.indexOf(path);
    if (index === -1) {
      favoritePaths.value.push(path);
    } else {
      favoritePaths.value.splice(index, 1);
    }
    saveFavorites(favoritePaths.value);

    // 同步更新访问历史中的收藏状态
    historyManager.updateFavoriteStatus(path, index === -1);
  }

  /**
   * 检查是否已收藏
   */
  function isFavorite(path: string): boolean {
    return favoritePaths.value.includes(path);
  }

  /**
   * 记录访问（委托给 composable）
   */
  function recordAccess(record: Omit<Parameters<typeof historyManager.recordAccess>[0], "isFavorite">) {
    historyManager.recordAccess({
      ...record,
      isFavorite: favoritePaths.value.includes(record.path),
    });
  }

  /**
   * 清空访问历史（委托给 composable）
   */
  function clearHistory() {
    historyManager.clearHistory();
  }

  /**
   * 处理导航（关闭抽屉）
   */
  function handleNavigation() {
    drawerOpen.value = false;
  }

  /**
   * 获取当前路由对应的菜单项
   */
  function getMenuItemByPath(path: string): MenuItem | undefined {
    return allMenuItems.value.find((item) => item.path === path);
  }

  return {
    // 状态
    drawerOpen,
    favoritePaths,
    accessHistory,

    // 计算属性
    allMenuItems,
    favoriteItems,
    recentItems,
    menuItemsByCategory,

    // 方法
    toggleFavorite,
    isFavorite,
    recordAccess,
    clearHistory,
    handleNavigation,
    getMenuItemByPath,
  };
});
