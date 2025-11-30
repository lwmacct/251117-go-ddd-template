/**
 * 点击外部检测 Composable
 * 检测点击是否发生在指定元素外部
 */

import { ref, onMounted, onUnmounted, type Ref } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UseClickOutsideOptions {
  /** 是否立即激活，默认 true */
  immediate?: boolean;
  /** 事件类型，默认 "pointerdown" */
  event?: "click" | "mousedown" | "mouseup" | "pointerdown" | "pointerup";
  /** 是否检测右键点击，默认 true */
  detectRightClick?: boolean;
  /** 忽略的元素选择器或元素列表 */
  ignore?: (string | Ref<HTMLElement | null | undefined>)[];
  /** 是否在捕获阶段处理，默认 true */
  capture?: boolean;
}

// ============================================================================
// 主函数
// ============================================================================

/**
 * 点击外部检测
 * @example
 * const target = ref<HTMLElement>()
 * const { isOutside } = useClickOutside(target, () => {
 *   console.log('点击了外部')
 * })
 */
export function useClickOutside(
  target: Ref<HTMLElement | null | undefined> | Ref<HTMLElement | null | undefined>[],
  callback?: (event: PointerEvent | MouseEvent) => void,
  options: UseClickOutsideOptions = {}
) {
  const {
    immediate = true,
    event = "pointerdown",
    detectRightClick = true,
    ignore = [],
    capture = true,
  } = options;

  const isOutside = ref(false);
  let isActive = immediate;

  // 获取所有目标元素
  const getTargets = (): HTMLElement[] => {
    const targets = Array.isArray(target) ? target : [target];
    return targets
      .map((t) => t.value)
      .filter((el): el is HTMLElement => el != null);
  };

  // 检查是否应该忽略
  const shouldIgnore = (event: Event): boolean => {
    const path = event.composedPath();

    return ignore.some((item) => {
      if (typeof item === "string") {
        // 选择器
        return path.some((el) => {
          if (el instanceof Element) {
            return el.matches(item);
          }
          return false;
        });
      } else {
        // Ref<HTMLElement>
        return item.value && path.includes(item.value);
      }
    });
  };

  // 处理点击事件
  const handleClick = (e: Event) => {
    if (!isActive) return;

    const event = e as PointerEvent | MouseEvent;

    // 检查右键点击
    if (!detectRightClick && event.button === 2) {
      return;
    }

    // 检查是否应该忽略
    if (shouldIgnore(event)) {
      return;
    }

    const targets = getTargets();
    if (targets.length === 0) return;

    const path = event.composedPath();

    // 检查是否点击在目标元素内部
    const isClickInside = targets.some((target) => path.includes(target));

    isOutside.value = !isClickInside;

    if (!isClickInside) {
      callback?.(event);
    }
  };

  // 激活监听
  const activate = () => {
    isActive = true;
    window.addEventListener(event, handleClick, { capture, passive: true });
  };

  // 停用监听
  const deactivate = () => {
    isActive = false;
    window.removeEventListener(event, handleClick, { capture });
  };

  onMounted(() => {
    if (immediate) {
      activate();
    }
  });

  onUnmounted(() => {
    deactivate();
  });

  return {
    isOutside,
    activate,
    deactivate,
  };
}

// ============================================================================
// 可切换的点击外部检测
// ============================================================================

export interface UseClickOutsideToggleOptions extends UseClickOutsideOptions {
  /** 初始状态 */
  initialValue?: boolean;
}

/**
 * 可切换的点击外部检测
 * 适用于下拉菜单、弹出框等场景
 * @example
 * const menuRef = ref<HTMLElement>()
 * const { isOpen, open, close, toggle } = useClickOutsideToggle(menuRef)
 */
export function useClickOutsideToggle(
  target: Ref<HTMLElement | null | undefined> | Ref<HTMLElement | null | undefined>[],
  options: UseClickOutsideToggleOptions = {}
) {
  const { initialValue = false, ...clickOutsideOptions } = options;

  const isOpen = ref(initialValue);

  const close = () => {
    isOpen.value = false;
  };

  const open = () => {
    isOpen.value = true;
  };

  const toggle = () => {
    isOpen.value = !isOpen.value;
  };

  // 点击外部时关闭
  const { isOutside, activate, deactivate } = useClickOutside(
    target,
    () => {
      if (isOpen.value) {
        close();
      }
    },
    clickOutsideOptions
  );

  return {
    isOpen,
    isOutside,
    open,
    close,
    toggle,
    activate,
    deactivate,
  };
}

// ============================================================================
// Vue 指令形式
// ============================================================================

export interface ClickOutsideBinding {
  handler: (event: PointerEvent | MouseEvent) => void;
  ignore?: string[];
}

/**
 * v-click-outside 指令
 * @example
 * // 在 main.ts 中注册
 * app.directive('click-outside', vClickOutside)
 *
 * // 使用
 * <div v-click-outside="handleClickOutside">...</div>
 * <div v-click-outside="{ handler: handleClickOutside, ignore: ['.ignore-class'] }">...</div>
 */
export const vClickOutside = {
  mounted(
    el: HTMLElement,
    binding: { value: ((e: Event) => void) | ClickOutsideBinding }
  ) {
    const handler =
      typeof binding.value === "function"
        ? binding.value
        : binding.value.handler;
    const ignore =
      typeof binding.value === "function" ? [] : binding.value.ignore || [];

    const handleClick = (event: Event) => {
      const e = event as PointerEvent | MouseEvent;
      const path = e.composedPath();

      // 检查是否应该忽略
      const shouldIgnore = ignore.some((selector) => {
        return path.some((el) => {
          if (el instanceof Element) {
            return el.matches(selector);
          }
          return false;
        });
      });

      if (shouldIgnore) return;

      // 检查是否点击在元素外部
      if (!path.includes(el)) {
        handler(e);
      }
    };

    // 存储处理函数以便卸载
    (el as HTMLElement & { __clickOutside: (e: Event) => void }).__clickOutside =
      handleClick;

    window.addEventListener("pointerdown", handleClick, { capture: true });
  },

  unmounted(el: HTMLElement) {
    const handleClick = (el as HTMLElement & { __clickOutside?: (e: Event) => void })
      .__clickOutside;
    if (handleClick) {
      window.removeEventListener("pointerdown", handleClick, { capture: true });
    }
  },
};
