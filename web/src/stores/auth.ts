/**
 * 认证状态管理 Store (Pinia)
 */
import { defineStore } from "pinia";
import { ref, computed } from "vue";
import {
  login as basicLogin,
  register as basicRegister,
  logout as apiLogout,
  getCurrentUser,
  AuthAPI,
} from "@/api/auth";
import { getAccessToken, clearAuthTokens, saveAccessToken, saveRefreshToken } from "@/utils/auth";
import type { LoginRequest, BasicLoginRequest, BasicRegisterRequest, User, LoginResult } from "@/types/auth";

export const useAuthStore = defineStore("auth", () => {
  // 状态
  const currentUser = ref<User | null>(null);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  // 计算属性
  const isAuthenticated = computed(() => !!currentUser.value);
  const hasToken = computed(() => !!getAccessToken());

  /**
   * 初始化认证状态
   * 检查 localStorage 中的 token，如果存在则获取用户信息
   */
  async function initAuth() {
    const token = getAccessToken();
    if (!token) {
      currentUser.value = null;
      return;
    }

    try {
      isLoading.value = true;
      error.value = null;
      const user = await getCurrentUser();
      currentUser.value = user;
    } catch (err) {
      // token 无效，清除并重置状态
      clearAuthTokens();
      currentUser.value = null;
      error.value = (err as Error).message || "Failed to initialize auth";
    } finally {
      isLoading.value = false;
    }
  }

  /**
   * 登录
   * 支持两种方式：标准登录（带验证码）和基础登录（已废弃）
   */
  async function login(credentials: LoginRequest | BasicLoginRequest): Promise<LoginResult> {
    try {
      isLoading.value = true;
      error.value = null;

      // 检查是否是标准登录请求 (带验证码)
      if ("captcha_id" in credentials) {
        // 使用标准 API（带验证码）
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
            saveAccessToken(response.data.access_token);
          }
          if (response.data?.refresh_token) {
            saveRefreshToken(response.data.refresh_token);
          }

          // 保存用户信息
          if (response.data?.user) {
            currentUser.value = response.data.user;
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
      } else {
        // 使用基础 API（已废弃，仅向后兼容）
        const response = await basicLogin(credentials);
        currentUser.value = response.user;
        return {
          success: true,
          requiresTwoFactor: false,
        };
      }
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
   * 注册（基础版，已废弃）
   * @deprecated 使用带验证码的注册流程
   */
  async function register(data: BasicRegisterRequest) {
    try {
      isLoading.value = true;
      error.value = null;
      const response = await basicRegister(data);
      currentUser.value = response.user;
      return response;
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (err: any) {
      error.value = err.response?.data?.error || err.message || "Registration failed";
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  /**
   * 登出
   */
  function logout() {
    apiLogout();
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
  function updateUser(user: User) {
    currentUser.value = user;
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
    register,
    logout,
    clearError,
    updateUser,
  };
});
