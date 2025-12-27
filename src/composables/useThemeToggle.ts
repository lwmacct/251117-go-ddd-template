/**
 * 主题切换 Composable
 *
 * 提供统一的主题切换逻辑，支持：
 * - 手动切换 light/dark
 * - 持久化到 localStorage（使用 VueUse useLocalStorage）
 * - 响应式状态
 *
 * 使用场景：
 * - AuthLayout 主题切换按钮
 * - AppBars 主题切换按钮
 */
import { computed } from "vue";
import { useTheme } from "vuetify";
import { useLocalStorage } from "@vueuse/core";

// 模块级别的主题存储（响应式，自动同步到 localStorage）
const storedTheme = useLocalStorage<"light" | "dark">("theme", "light");

export function useThemeToggle() {
  const theme = useTheme();

  /** 当前主题名称 */
  const currentTheme = computed(() => theme.global.name.value);

  /** 是否为暗色主题 */
  const isDark = computed(() => currentTheme.value === "dark");

  /** 切换主题图标 */
  const themeIcon = computed(() => (isDark.value ? "mdi-weather-sunny" : "mdi-weather-night"));

  /** 切换主题 */
  function toggleTheme() {
    const newTheme = isDark.value ? "light" : "dark";
    theme.change(newTheme);
    storedTheme.value = newTheme; // 响应式存储，自动同步
  }

  return {
    currentTheme,
    isDark,
    themeIcon,
    toggleTheme,
  };
}
