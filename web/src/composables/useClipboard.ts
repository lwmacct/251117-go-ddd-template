/**
 * 剪贴板操作 Composable
 * 提供复制到剪贴板功能，支持成功/失败回调
 */
import { ref } from "vue";

export interface UseClipboardOptions {
  /** 复制成功后显示成功状态的持续时间（毫秒） */
  successDuration?: number;
}

export function useClipboard(options: UseClipboardOptions = {}) {
  const { successDuration = 2000 } = options;

  const copied = ref(false);
  const error = ref<string | null>(null);

  /**
   * 复制文本到剪贴板
   * @param text 要复制的文本
   * @returns 是否成功
   */
  const copy = async (text: string): Promise<boolean> => {
    error.value = null;

    try {
      // 使用现代 Clipboard API
      if (navigator.clipboard && window.isSecureContext) {
        await navigator.clipboard.writeText(text);
      } else {
        // 降级方案：使用 execCommand
        const textArea = document.createElement("textarea");
        textArea.value = text;
        textArea.style.position = "fixed";
        textArea.style.left = "-9999px";
        textArea.style.top = "-9999px";
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();

        try {
          document.execCommand("copy");
        } finally {
          document.body.removeChild(textArea);
        }
      }

      copied.value = true;

      // 自动重置复制状态
      setTimeout(() => {
        copied.value = false;
      }, successDuration);

      return true;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "复制失败";
      copied.value = false;
      return false;
    }
  };

  return {
    copied,
    error,
    copy,
  };
}
