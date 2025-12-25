/**
 * 认证存储工具 - 使用 VueUse useLocalStorage
 *
 * 提供响应式的 localStorage 存储，支持在 Vue 组件和普通 JS 中使用。
 * 在 Vue 组件中可以直接 watch 这些 ref，实现状态同步。
 */
import { useLocalStorage } from "@vueuse/core";

// ============================================================================
// 响应式存储 Refs
// ============================================================================

/** 访问令牌（响应式） */
export const accessToken = useLocalStorage<string | null>("access_token", null);

/** 刷新令牌（响应式） */
export const refreshToken = useLocalStorage<string | null>("refresh_token", null);

/** 令牌过期时间（响应式） */
export const tokenExpiry = useLocalStorage<number | null>("token_expiry", null);

// ============================================================================
// 便捷函数
// ============================================================================

/**
 * 清除所有认证令牌
 */
export const clearAuthTokens = (): void => {
  accessToken.value = null;
  refreshToken.value = null;
  tokenExpiry.value = null;
};

/**
 * 检查是否有访问令牌
 */
export const hasAccessToken = (): boolean => !!accessToken.value;
