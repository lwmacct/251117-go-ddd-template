/**
 * Snackbar Store - 全局消息提示队列管理
 *
 * 配合 Vuetify v-snackbar-queue 组件使用
 * 支持 success/error/warning/info 四种消息类型
 */
import { defineStore } from "pinia";
import { ref } from "vue";

export interface SnackbarMessage {
  text: string;
  color?: "success" | "error" | "warning" | "info";
  timeout?: number;
}

export const useSnackbarStore = defineStore("snackbar", () => {
  const queue = ref<SnackbarMessage[]>([]);

  /**
   * 显示消息（通用方法）
   */
  function show(message: SnackbarMessage | string) {
    const record = typeof message === "string" ? { text: message } : message;
    queue.value.push({
      timeout: 3000,
      color: "success",
      ...record,
    });
  }

  /**
   * 显示成功消息
   */
  function success(text: string, timeout = 3000) {
    show({ text, color: "success", timeout });
  }

  /**
   * 显示错误消息
   */
  function error(text: string, timeout = 5000) {
    show({ text, color: "error", timeout });
  }

  /**
   * 显示警告消息
   */
  function warning(text: string, timeout = 4000) {
    show({ text, color: "warning", timeout });
  }

  /**
   * 显示信息消息
   */
  function info(text: string, timeout = 3000) {
    show({ text, color: "info", timeout });
  }

  return { queue, show, success, error, warning, info };
});
