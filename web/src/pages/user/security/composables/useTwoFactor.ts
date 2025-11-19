import { ref } from "vue";
import { PlatformAuthAPI } from "@/api";

/**
 * 2FA 管理 Composable
 *
 * 职责：
 * - 封装 2FA 相关的业务逻辑
 * - 管理 2FA 状态和数据
 * - 提供可复用的 2FA 操作方法
 */
export function useTwoFactor() {
  // ========== 状态管理 ==========
  const loading = ref(false);
  const enabled = ref(false);
  const recoveryCodesCount = ref(0);
  const showSetupDialog = ref(false);
  const showDisableDialog = ref(false);

  // 设置 2FA 相关状态
  const qrcodeImage = ref("");
  const secret = ref("");
  const verifyCode = ref("");
  const recoveryCodes = ref<string[]>([]);
  const setupStep = ref<"qrcode" | "verify" | "codes">("qrcode");

  // 消息状态
  const errorMessage = ref("");
  const successMessage = ref("");

  // ========== API 调用方法 ==========

  /**
   * 获取 2FA 状态
   */
  async function fetchStatus() {
    try {
      loading.value = true;
      const response = await PlatformAuthAPI.get2FAStatus();
      if (response.code === 200 && response.data) {
        enabled.value = response.data.enabled;
        recoveryCodesCount.value = response.data.recovery_codes_count;
      }
    } catch (error: any) {
      console.error("获取 2FA 状态失败:", error);
      errorMessage.value = error.message || "获取 2FA 状态失败";
    } finally {
      loading.value = false;
    }
  }

  /**
   * 开始设置 2FA
   */
  async function startSetup() {
    try {
      loading.value = true;
      errorMessage.value = "";

      const response = await PlatformAuthAPI.setup2FA();
      if (response.code === 200 && response.data) {
        qrcodeImage.value = response.data.qrcode_img;
        secret.value = response.data.secret;
        setupStep.value = "qrcode";
        showSetupDialog.value = true;
      }
    } catch (error: any) {
      errorMessage.value = error.message || "设置 2FA 失败";
    } finally {
      loading.value = false;
    }
  }

  /**
   * 验证并启用 2FA
   */
  async function verifyAndEnable() {
    if (verifyCode.value.length !== 6) {
      errorMessage.value = "请输入 6 位验证码";
      return;
    }

    try {
      loading.value = true;
      errorMessage.value = "";

      const response = await PlatformAuthAPI.enable2FA(verifyCode.value);
      if (response.code === 200 && response.data) {
        recoveryCodes.value = response.data.recovery_codes;
        setupStep.value = "codes";
        successMessage.value = "2FA 已成功启用！";

        // 更新状态
        await fetchStatus();
      }
    } catch (error: any) {
      errorMessage.value = error.message || "验证失败";
    } finally {
      loading.value = false;
    }
  }

  /**
   * 禁用 2FA
   */
  async function disable2FA() {
    try {
      loading.value = true;
      errorMessage.value = "";

      const response = await PlatformAuthAPI.disable2FA();
      if (response.code === 200) {
        successMessage.value = "2FA 已成功禁用";
        showDisableDialog.value = false;

        // 更新状态
        await fetchStatus();
      }
    } catch (error: any) {
      errorMessage.value = error.message || "禁用 2FA 失败";
    } finally {
      loading.value = false;
    }
  }

  // ========== 工具方法 ==========

  /**
   * 关闭设置对话框
   */
  function closeSetupDialog() {
    showSetupDialog.value = false;
    setupStep.value = "qrcode";
    verifyCode.value = "";
    qrcodeImage.value = "";
    secret.value = "";
    recoveryCodes.value = [];
    errorMessage.value = "";
  }

  /**
   * 复制文本到剪贴板
   */
  async function copyToClipboard(text: string) {
    try {
      await navigator.clipboard.writeText(text);
      successMessage.value = "已复制到剪贴板";
      setTimeout(() => {
        successMessage.value = "";
      }, 2000);
    } catch (error) {
      errorMessage.value = "复制失败，请手动复制";
    }
  }

  /**
   * 下载恢复码
   */
  function downloadRecoveryCodes() {
    const text = recoveryCodes.value.join("\n");
    const blob = new Blob([text], { type: "text/plain" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = "2fa-recovery-codes.txt";
    link.click();
    URL.revokeObjectURL(url);
  }

  // ========== 返回公共接口 ==========
  return {
    // 状态
    loading,
    enabled,
    recoveryCodesCount,
    showSetupDialog,
    showDisableDialog,
    qrcodeImage,
    secret,
    verifyCode,
    recoveryCodes,
    setupStep,
    errorMessage,
    successMessage,

    // 方法
    fetchStatus,
    startSetup,
    verifyAndEnable,
    disable2FA,
    closeSetupDialog,
    copyToClipboard,
    downloadRecoveryCodes,
  };
}
