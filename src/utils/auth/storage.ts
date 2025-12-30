/**
 * 认证存储工具 - 使用 VueUse useLocalStorage
 *
 * 提供响应式的 localStorage 存储，支持在 Vue 组件和普通 JS 中使用。
 * 在 Vue 组件中可以直接 watch 这些 ref，实现状态同步。
 */
import { useLocalStorage } from "@vueuse/core";
import type { AuthUserBriefDTO } from "@models";

// ============================================================================
// 响应式存储 Refs
// ============================================================================

/** 访问令牌（响应式） */
export const accessToken = useLocalStorage<string | null>("auth_access_token", null);

/** 刷新令牌（响应式） */
export const refreshToken = useLocalStorage<string | null>("auth_refresh_token", null);

/** 令牌过期时间（响应式） */
export const tokenExpiry = useLocalStorage<number | null>("auth_token_expiry", null);

/** 当前用户信息（响应式持久化，使用 JSON 序列化） */
export const storedUser = useLocalStorage<AuthUserBriefDTO | null>("auth_current_user", null, {
  serializer: {
    read: (v) => (v ? JSON.parse(v) : null),
    write: (v) => JSON.stringify(v),
  },
});

// ============================================================================
// 便捷函数
// ============================================================================

/**
 * 清除所有认证令牌和用户信息
 */
export const clearAuthTokens = (): void => {
  accessToken.value = null;
  refreshToken.value = null;
  tokenExpiry.value = null;
  storedUser.value = null;
};

/**
 * 检查是否有访问令牌
 */
export const hasAccessToken = (): boolean => !!accessToken.value;
