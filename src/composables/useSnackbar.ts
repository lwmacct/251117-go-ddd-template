/**
 * useSnackbar - 全局消息提示 Composable
 *
 * 提供便捷的消息提示方法，封装 Snackbar Store
 *
 * @example
 * ```ts
 * const { success, error, warning, info } = useSnackbar()
 *
 * // 显示成功消息
 * success('操作成功')
 *
 * // 显示错误消息（自动延长显示时间）
 * error('操作失败：' + err.message)
 *
 * // 显示警告消息
 * warning('请注意：此操作不可撤销')
 * ```
 */
import { useSnackbarStore } from "@/stores/snackbar";

export function useSnackbar() {
  const store = useSnackbarStore();

  return {
    /** 显示成功消息（绿色，3秒） */
    success: store.success,
    /** 显示错误消息（红色，5秒） */
    error: store.error,
    /** 显示警告消息（橙色，4秒） */
    warning: store.warning,
    /** 显示信息消息（蓝色，3秒） */
    info: store.info,
    /** 通用显示方法 */
    show: store.show,
  };
}
