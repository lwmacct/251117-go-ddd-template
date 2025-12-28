/**
 * 响应式 Tabs Composable
 *
 * 提供响应式 Tabs 布局逻辑，支持：
 * - 响应式断点切换（垂直/水平布局）
 * - URL 查询参数同步（使用 VueUse useRouteQuery）
 * - Tab 切换回调（用于懒加载场景）
 *
 * 使用场景：
 * - admin/settings 系统设置页面
 * - user/settings 用户设置页面
 * - user/security 安全设置页面
 * - 任何需要响应式 Tabs 的页面
 */
import { computed, watch, type ComputedRef, type Ref } from "vue";
import { useRouteQuery } from "@vueuse/router";
import { useDisplay } from "vuetify";

/** Tab 项定义 */
export interface TabItem {
  /** Tab 唯一标识 */
  value: string;
  /** Tab 显示文本 */
  label: string;
  /** Tab 图标（可选） */
  icon?: string;
}

/** useResponsiveTabs 配置选项 */
export interface UseResponsiveTabsOptions {
  /** 默认 Tab（URL 无参数时使用） */
  defaultTab: string;
  /** URL 查询参数名，默认 'tab' */
  queryParam?: string;
  /** 响应式断点（px），默认 720 */
  breakpoint?: number;
  /** Tab 切换回调（用于懒加载） */
  onTabChange?: (tab: string) => void | Promise<void>;
}

/** useResponsiveTabs 返回值 */
export interface UseResponsiveTabsReturn {
  /** 当前激活的 Tab */
  currentTab: Ref<string>;
  /** 是否使用垂直布局（宽屏） */
  isVertical: ComputedRef<boolean>;
  /** 处理 Tab 切换（更新 URL + 触发回调） */
  handleTabChange: (tab: string) => Promise<void>;
}

/**
 * 响应式 Tabs Composable
 *
 * @example
 * ```ts
 * const { currentTab, isVertical, handleTabChange } = useResponsiveTabs({
 *   defaultTab: 'general',
 *   onTabChange: async (tab) => {
 *     // 懒加载逻辑
 *     if (!isCategoryLoaded(tab)) {
 *       await fetchSchemaByCategory(tab)
 *     }
 *   },
 * })
 * ```
 */
export function useResponsiveTabs(options: UseResponsiveTabsOptions): UseResponsiveTabsReturn {
  const { defaultTab, queryParam = "tab", breakpoint = 720, onTabChange } = options;

  const { width } = useDisplay();

  // VueUse: 自动双向同步 URL 查询参数
  const currentTab = useRouteQuery(queryParam, defaultTab, {
    mode: "replace", // 使用 replace 而非 push，避免历史记录堆积
  });

  /** 是否使用垂直布局（宽度 >= 断点） */
  const isVertical = computed(() => width.value >= breakpoint);

  /** 处理 Tab 切换 */
  const handleTabChange = async (tab: string) => {
    if (tab === currentTab.value) return;
    currentTab.value = tab;
    if (onTabChange) {
      await onTabChange(tab);
    }
  };

  // 监听 URL 变化（浏览器前进/后退）
  watch(
    currentTab,
    async (newTab, oldTab) => {
      if (newTab !== oldTab && onTabChange) {
        await onTabChange(newTab);
      }
    },
    { flush: "post" },
  );

  return {
    currentTab,
    isVertical,
    handleTabChange,
  };
}
