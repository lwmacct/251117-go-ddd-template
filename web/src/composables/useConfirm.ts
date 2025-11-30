/**
 * 确认对话框 Composable
 * 提供程序化调用确认对话框的能力
 */
import { ref, readonly } from "vue";

export interface ConfirmOptions {
  /** 标题 */
  title?: string;
  /** 消息内容 */
  message: string;
  /** 对话框类型 */
  type?: "delete" | "warning" | "info";
  /** 确认按钮文本 */
  confirmText?: string;
  /** 取消按钮文本 */
  cancelText?: string;
}

export function useConfirm() {
  const visible = ref(false);
  const loading = ref(false);
  const options = ref<ConfirmOptions>({
    message: "",
  });

  let resolvePromise: ((value: boolean) => void) | null = null;

  /**
   * 显示确认对话框
   * @param opts 配置选项
   * @returns Promise<boolean> - 用户点击确认返回 true，取消返回 false
   */
  const confirm = (opts: ConfirmOptions): Promise<boolean> => {
    options.value = opts;
    visible.value = true;

    return new Promise((resolve) => {
      resolvePromise = resolve;
    });
  };

  /**
   * 处理确认
   */
  const handleConfirm = () => {
    if (resolvePromise) {
      resolvePromise(true);
      resolvePromise = null;
    }
    visible.value = false;
  };

  /**
   * 处理取消
   */
  const handleCancel = () => {
    if (resolvePromise) {
      resolvePromise(false);
      resolvePromise = null;
    }
    visible.value = false;
  };

  /**
   * 设置加载状态
   */
  const setLoading = (value: boolean) => {
    loading.value = value;
  };

  /**
   * 快捷方法：删除确认
   */
  const confirmDelete = (itemName: string): Promise<boolean> => {
    return confirm({
      type: "delete",
      title: "确认删除",
      message: `确定要删除 <strong>${itemName}</strong> 吗？此操作不可恢复。`,
      confirmText: "删除",
    });
  };

  /**
   * 快捷方法：危险操作确认
   */
  const confirmDanger = (message: string, title = "警告"): Promise<boolean> => {
    return confirm({
      type: "warning",
      title,
      message,
      confirmText: "继续",
    });
  };

  return {
    // 状态（只读）
    visible: readonly(visible),
    loading: readonly(loading),
    options: readonly(options),

    // 方法
    confirm,
    confirmDelete,
    confirmDanger,
    handleConfirm,
    handleCancel,
    setLoading,
  };
}
