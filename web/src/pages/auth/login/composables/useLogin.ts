/**
 * Login 页面状态管理 Composable
 *
 * 管理登录表单状态，包括：
 * - 登录表单数据（账号、密码、验证码）
 * - 验证码获取和显示
 * - 错误消息管理
 * - Session token（用于2FA验证流程）
 */

import { ref, computed } from "vue";
import { useAuthStore } from "@/stores/auth";
import { PlatformAuthAPI } from "@/api";
import type { CaptchaData } from "@/api";
import type { LoginResult } from "@/types";

/**
 * Login 页面状态管理 Composable
 *
 * 管理登录表单状态（支持手机号/用户名/邮箱登录）
 */
export function useLogin() {
  const authStore = useAuthStore();

  // === 状态 ===
  const account = ref(""); // 手机号/用户名/邮箱
  const password = ref("");
  const errorMessage = ref("");
  const showError = ref(false);
  const captchaData = ref<CaptchaData | null>(null);
  const captchaCode = ref("");
  const loadingCaptcha = ref(false);
  const sessionToken = ref<string>(""); // 临时会话token（用于2FA验证）

  // 错误提示自动消失定时器
  let errorTimer: ReturnType<typeof setTimeout> | null = null;

  // === 计算属性 ===
  const isFormValid = computed(() => {
    return account.value.length > 0 && password.value.length > 0 && captchaCode.value.length > 0 && captchaData.value !== null;
  });

  const isLoading = computed(() => authStore.isLoading);

  // 验证码图片（用于模板，避免类型错误）
  const captchaImage = computed(() => captchaData.value?.image || "");

  // === 方法 ===

  /**
   * 显示错误消息（5秒后自动消失）
   */
  const showErrorMessage = (message: string, duration = 5000) => {
    // 清除之前的定时器
    if (errorTimer) {
      clearTimeout(errorTimer);
    }

    errorMessage.value = message;

    // 设置自动消失
    if (duration > 0) {
      errorTimer = setTimeout(() => {
        errorMessage.value = "";
      }, duration);
    }
  };

  /**
   * 清除错误消息
   */
  const clearError = () => {
    if (errorTimer) {
      clearTimeout(errorTimer);
      errorTimer = null;
    }
    errorMessage.value = "";
  };

  /**
   * 获取验证码
   */
  const fetchCaptcha = async () => {
    try {
      loadingCaptcha.value = true;
      const response = await PlatformAuthAPI.getCaptcha();
      if (response.data) {
        captchaData.value = response.data;
        captchaCode.value = ""; // 清空验证码输入
      }
    } catch (error) {
      showErrorMessage(error instanceof Error ? error.message : "获取验证码失败");
    } finally {
      loadingCaptcha.value = false;
    }
  };

  /**
   * 登录
   */
  const login = async (): Promise<LoginResult> => {
    if (!isFormValid.value) {
      showErrorMessage("请填写完整的登录信息");
      return {
        success: false,
        message: "请填写完整的登录信息",
        requiresTwoFactor: false,
      };
    }

    if (!captchaData.value) {
      showErrorMessage("请先获取验证码");
      return {
        success: false,
        message: "请先获取验证码",
        requiresTwoFactor: false,
      };
    }

    // 清除之前的错误消息
    clearError();

    const result = await authStore.login({
      account: account.value,
      password: password.value,
      captcha_id: captchaData.value.id,
      captcha: captchaCode.value,
    });

    // 如果需要2FA验证，保存session token
    if (result.requiresTwoFactor && result.sessionToken) {
      sessionToken.value = result.sessionToken;
    }

    // 如果不需要2FA验证且登录失败，显示错误
    if (!result.success && !result.requiresTwoFactor) {
      showErrorMessage(result.message || "登录失败");
      // 登录失败后刷新验证码
      await fetchCaptcha();
    }
    // 如果需要2FA验证，不显示错误消息，直接切换组件

    return result;
  };

  /**
   * 清空表单
   */
  const clearForm = () => {
    account.value = "";
    password.value = "";
    captchaCode.value = "";
    clearError();
    captchaData.value = null;
  };

  return {
    // 状态
    account,
    password,
    errorMessage,
    showError,
    captchaData,
    captchaCode,
    loadingCaptcha,
    sessionToken,

    // 计算属性
    isFormValid,
    isLoading,
    captchaImage,

    // 方法
    login,
    clearForm,
    fetchCaptcha,
  };
}
