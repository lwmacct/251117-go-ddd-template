/**
 * 认证状态管理 Store (Pinia)
 */
import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { AuthAPI } from "@/api";
import { accessToken, refreshToken, clearAuthTokens, storedUser } from "@/utils/auth";
import type { AuthLoginDTO, AuthUserBriefDTO } from "@models";
import type { LoginResult } from "@/api";

export const useAuthStore = defineStore("auth", () => {
  // 状态
  const currentUser = ref<AuthUserBriefDTO | null>(null);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  // 计算属性
  const isAuthenticated = computed(() => !!currentUser.value);
  const hasToken = computed(() => !!accessToken.value);

  /**
   * 初始化认证状态
   * 从 localStorage 恢复用户信息（不请求 profile API）
   */
  function initAuth() {
    const token = accessToken.value;
    if (!token) {
      currentUser.value = null;
      storedUser.value = null;
      return;
    }

    // 从 localStorage 恢复用户信息
    if (storedUser.value) {
      currentUser.value = storedUser.value;
    }
  }

  /**
   * 登录（标准版，带验证码）
   */
  async function login(credentials: AuthLoginDTO): Promise<LoginResult> {
    try {
      isLoading.value = true;
      error.value = null;

      const response = await AuthAPI.login(credentials);

      if (response.code === 200) {
        // 检查是否需要 2FA
        if (response.data?.requires_2fa) {
          return {
            success: false,
            requiresTwoFactor: true,
            sessionToken: response.data.session_token,
            message: response.message,
          };
        }

        // 登录成功 - 保存 token
        if (response.data?.access_token) {
          accessToken.value = response.data.access_token;
        }
        if (response.data?.refresh_token) {
          refreshToken.value = response.data.refresh_token;
        }

        // 保存用户信息（内存 + localStorage）
        if (response.data?.user) {
          currentUser.value = response.data.user;
          storedUser.value = response.data.user;
        }
        return {
          success: true,
          requiresTwoFactor: false,
          message: response.message,
        };
      }

      // 登录失败
      error.value = response.message;
      return {
        success: false,
        requiresTwoFactor: false,
        message: response.message,
      };
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (err: any) {
      const message = err.response?.data?.error || err.message || "Login failed";
      error.value = message;
      return {
        success: false,
        requiresTwoFactor: false,
        message,
      };
    } finally {
      isLoading.value = false;
    }
  }

  /**
   * 登出
   */
  function logout() {
    clearAuthTokens();
    currentUser.value = null;
    error.value = null;
  }

  /**
   * 清除错误
   */
  function clearError() {
    error.value = null;
  }

  /**
   * 更新用户信息
   */
  function updateUser(user: AuthUserBriefDTO) {
    currentUser.value = user;
    storedUser.value = user;
  }

  return {
    // 状态
    currentUser,
    isLoading,
    error,

    // 计算属性
    isAuthenticated,
    hasToken,

    // 方法
    initAuth,
    login,
    logout,
    clearError,
    updateUser,
  };
});
