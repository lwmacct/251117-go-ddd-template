/**
 * 焦点陷阱 Composable
 * 将键盘焦点限制在指定容器内，用于模态框等场景的无障碍访问
 */

import { ref, watch, onMounted, onUnmounted, type Ref } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UseFocusTrapOptions {
  /** 是否立即激活，默认 false */
  immediate?: boolean;
  /** 激活时是否自动聚焦第一个元素，默认 true */
  autoFocus?: boolean;
  /** 停用时是否恢复之前的焦点，默认 true */
  restoreFocus?: boolean;
  /** 初始聚焦的元素选择器 */
  initialFocus?: string;
  /** 按 Escape 键时是否停用，默认 false */
  escapeDeactivates?: boolean;
  /** 点击外部时是否停用，默认 false */
  clickOutsideDeactivates?: boolean;
  /** 允许在容器外 Tab，默认 false */
  allowOutsideTab?: boolean;
  /** 停用回调 */
  onDeactivate?: () => void;
}

// ============================================================================
// 常量
// ============================================================================

// 可聚焦元素选择器
const FOCUSABLE_SELECTORS = [
  'a[href]:not([tabindex="-1"])',
  'area[href]:not([tabindex="-1"])',
  'input:not([disabled]):not([tabindex="-1"])',
  'select:not([disabled]):not([tabindex="-1"])',
  'textarea:not([disabled]):not([tabindex="-1"])',
  'button:not([disabled]):not([tabindex="-1"])',
  'iframe:not([tabindex="-1"])',
  'audio[controls]:not([tabindex="-1"])',
  'video[controls]:not([tabindex="-1"])',
  '[contenteditable]:not([tabindex="-1"])',
  '[tabindex]:not([tabindex="-1"])',
].join(",");

// ============================================================================
// 工具函数
// ============================================================================

/**
 * 获取容器内所有可聚焦元素
 */
function getFocusableElements(container: HTMLElement): HTMLElement[] {
  const elements = Array.from(container.querySelectorAll<HTMLElement>(FOCUSABLE_SELECTORS));

  // 过滤不可见元素
  return elements.filter((el) => {
    return el.offsetWidth > 0 && el.offsetHeight > 0 && getComputedStyle(el).visibility !== "hidden";
  });
}

/**
 * 获取第一个可聚焦元素
 */
function getFirstFocusable(container: HTMLElement): HTMLElement | null {
  const elements = getFocusableElements(container);
  return elements[0] || null;
}

/**
 * 获取最后一个可聚焦元素
 */
function getLastFocusable(container: HTMLElement): HTMLElement | null {
  const elements = getFocusableElements(container);
  return elements[elements.length - 1] || null;
}

// ============================================================================
// 主函数
// ============================================================================

/**
 * 焦点陷阱
 * @example
 * const modalRef = ref<HTMLElement>()
 * const { activate, deactivate, isActive } = useFocusTrap(modalRef)
 *
 * // 打开模态框时激活
 * activate()
 *
 * // 关闭时停用
 * deactivate()
 */
export function useFocusTrap(target: Ref<HTMLElement | null | undefined>, options: UseFocusTrapOptions = {}) {
  const {
    immediate = false,
    autoFocus = true,
    restoreFocus = true,
    initialFocus,
    escapeDeactivates = false,
    clickOutsideDeactivates = false,
    allowOutsideTab = false,
    onDeactivate,
  } = options;

  const isActive = ref(false);
  let previouslyFocusedElement: HTMLElement | null = null;

  // 处理 Tab 键
  const handleKeyDown = (event: KeyboardEvent) => {
    if (!isActive.value) return;
    const container = target.value;
    if (!container) return;

    // Escape 键
    if (event.key === "Escape" && escapeDeactivates) {
      event.preventDefault();
      deactivate();
      onDeactivate?.();
      return;
    }

    // Tab 键
    if (event.key !== "Tab") return;

    if (allowOutsideTab) return;

    const focusableElements = getFocusableElements(container);
    if (focusableElements.length === 0) {
      event.preventDefault();
      return;
    }

    const firstElement = focusableElements[0];
    const lastElement = focusableElements[focusableElements.length - 1];
    const activeElement = document.activeElement as HTMLElement;

    // Shift + Tab（向前）
    if (event.shiftKey) {
      if (activeElement === firstElement || !container.contains(activeElement)) {
        event.preventDefault();
        lastElement.focus();
      }
    }
    // Tab（向后）
    else {
      if (activeElement === lastElement || !container.contains(activeElement)) {
        event.preventDefault();
        firstElement.focus();
      }
    }
  };

  // 处理点击外部
  const handleClick = (event: MouseEvent) => {
    if (!isActive.value || !clickOutsideDeactivates) return;
    const container = target.value;
    if (!container) return;

    if (!container.contains(event.target as Node)) {
      deactivate();
      onDeactivate?.();
    }
  };

  // 激活焦点陷阱
  const activate = () => {
    if (isActive.value) return;
    const container = target.value;
    if (!container) return;

    // 保存当前焦点
    if (restoreFocus) {
      previouslyFocusedElement = document.activeElement as HTMLElement;
    }

    isActive.value = true;

    // 添加事件监听
    document.addEventListener("keydown", handleKeyDown);
    if (clickOutsideDeactivates) {
      document.addEventListener("mousedown", handleClick);
    }

    // 自动聚焦
    if (autoFocus) {
      // 使用 nextTick 确保 DOM 已更新
      setTimeout(() => {
        let elementToFocus: HTMLElement | null = null;

        // 优先使用指定的初始焦点
        if (initialFocus) {
          elementToFocus = container.querySelector(initialFocus);
        }

        // 否则聚焦第一个可聚焦元素
        if (!elementToFocus) {
          elementToFocus = getFirstFocusable(container);
        }

        // 如果没有可聚焦元素，聚焦容器本身
        if (!elementToFocus) {
          if (!container.hasAttribute("tabindex")) {
            container.setAttribute("tabindex", "-1");
          }
          elementToFocus = container;
        }

        elementToFocus.focus();
      }, 0);
    }
  };

  // 停用焦点陷阱
  const deactivate = () => {
    if (!isActive.value) return;

    isActive.value = false;

    // 移除事件监听
    document.removeEventListener("keydown", handleKeyDown);
    document.removeEventListener("mousedown", handleClick);

    // 恢复焦点
    if (restoreFocus && previouslyFocusedElement) {
      previouslyFocusedElement.focus();
      previouslyFocusedElement = null;
    }
  };

  // 切换状态
  const toggle = () => {
    if (isActive.value) {
      deactivate();
    } else {
      activate();
    }
  };

  // 聚焦第一个元素
  const focusFirst = () => {
    const container = target.value;
    if (!container) return;

    const first = getFirstFocusable(container);
    first?.focus();
  };

  // 聚焦最后一个元素
  const focusLast = () => {
    const container = target.value;
    if (!container) return;

    const last = getLastFocusable(container);
    last?.focus();
  };

  // 监听目标变化
  watch(target, (newTarget, oldTarget) => {
    if (isActive.value && oldTarget && !newTarget) {
      deactivate();
    }
  });

  // 立即激活
  onMounted(() => {
    if (immediate && target.value) {
      activate();
    }
  });

  // 清理
  onUnmounted(() => {
    if (isActive.value) {
      deactivate();
    }
  });

  return {
    isActive,
    activate,
    deactivate,
    toggle,
    focusFirst,
    focusLast,
  };
}

// ============================================================================
// 响应式焦点陷阱
// ============================================================================

/**
 * 响应式焦点陷阱
 * 根据 ref 值自动激活/停用
 * @example
 * const isOpen = ref(false)
 * const modalRef = ref<HTMLElement>()
 * useFocusTrapWhenTrue(modalRef, isOpen)
 */
export function useFocusTrapWhenTrue(
  target: Ref<HTMLElement | null | undefined>,
  source: Ref<boolean>,
  options?: UseFocusTrapOptions
) {
  const trap = useFocusTrap(target, options);

  watch(
    source,
    (value) => {
      if (value) {
        // 等待 DOM 更新
        setTimeout(() => trap.activate(), 0);
      } else {
        trap.deactivate();
      }
    },
    { immediate: true }
  );

  return trap;
}

// ============================================================================
// 返回焦点管理
// ============================================================================

/**
 * 焦点返回管理
 * 保存当前焦点并在之后恢复
 * @example
 * const { save, restore } = useFocusReturn()
 *
 * // 打开模态框前保存焦点
 * save()
 *
 * // 关闭模态框后恢复焦点
 * restore()
 */
export function useFocusReturn() {
  const savedElement = ref<HTMLElement | null>(null);

  const save = () => {
    savedElement.value = document.activeElement as HTMLElement;
  };

  const restore = () => {
    if (savedElement.value) {
      savedElement.value.focus();
      savedElement.value = null;
    }
  };

  return {
    savedElement,
    save,
    restore,
  };
}

// ============================================================================
// 自动聚焦
// ============================================================================

/**
 * 自动聚焦指令数据
 * @example
 * const autoFocusRef = ref<HTMLElement>()
 * useAutoFocus(autoFocusRef)
 */
export function useAutoFocus(target: Ref<HTMLElement | null | undefined>) {
  onMounted(() => {
    // 等待 DOM 更新
    setTimeout(() => {
      target.value?.focus();
    }, 0);
  });
}
