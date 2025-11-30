/**
 * 滚动锁定 Composable
 * 锁定页面或元素滚动，用于模态框、抽屉等场景
 */

import { ref, watch, onUnmounted, type Ref } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UseScrollLockOptions {
  /** 是否立即锁定，默认 false */
  immediate?: boolean;
  /** 锁定时是否保持滚动条占位，防止页面抖动，默认 true */
  reserveScrollBarGap?: boolean;
  /** 允许滚动的元素选择器 */
  allowTouchMove?: string[];
}

// ============================================================================
// 工具函数
// ============================================================================

// 获取滚动条宽度
function getScrollbarWidth(): number {
  const outer = document.createElement("div");
  outer.style.visibility = "hidden";
  outer.style.overflow = "scroll";
  document.body.appendChild(outer);

  const inner = document.createElement("div");
  outer.appendChild(inner);

  const scrollbarWidth = outer.offsetWidth - inner.offsetWidth;
  outer.parentNode?.removeChild(outer);

  return scrollbarWidth;
}

// 检查是否有滚动条
function hasScrollbar(): boolean {
  return document.body.scrollHeight > window.innerHeight;
}

// ============================================================================
// 主函数
// ============================================================================

// 全局锁定计数器（支持嵌套）
let lockCount = 0;
let originalStyle: {
  overflow: string;
  paddingRight: string;
  position: string;
  top: string;
  width: string;
} | null = null;
let scrollPosition = 0;

/**
 * 滚动锁定
 * @example
 * const { isLocked, lock, unlock } = useScrollLock()
 *
 * // 显示模态框时
 * const showModal = () => {
 *   lock()
 *   modalVisible.value = true
 * }
 */
export function useScrollLock(options: UseScrollLockOptions = {}) {
  const { immediate = false, reserveScrollBarGap = true } = options;

  const isLocked = ref(false);

  // 锁定滚动
  const lock = () => {
    if (isLocked.value) return;

    lockCount++;
    isLocked.value = true;

    // 只在第一次锁定时保存原始样式
    if (lockCount === 1) {
      const body = document.body;
      const html = document.documentElement;

      // 保存原始样式
      originalStyle = {
        overflow: body.style.overflow,
        paddingRight: body.style.paddingRight,
        position: body.style.position,
        top: body.style.top,
        width: body.style.width,
      };

      // 保存滚动位置
      scrollPosition = window.scrollY;

      // 计算滚动条宽度
      const scrollbarWidth = reserveScrollBarGap && hasScrollbar() ? getScrollbarWidth() : 0;

      // 应用锁定样式
      body.style.overflow = "hidden";
      body.style.position = "fixed";
      body.style.top = `-${scrollPosition}px`;
      body.style.width = "100%";

      if (scrollbarWidth > 0) {
        body.style.paddingRight = `${scrollbarWidth}px`;
      }

      // 也锁定 html 元素（某些浏览器需要）
      html.style.overflow = "hidden";
    }
  };

  // 解锁滚动
  const unlock = () => {
    if (!isLocked.value) return;

    lockCount--;
    isLocked.value = false;

    // 只在最后一次解锁时恢复原始样式
    if (lockCount === 0 && originalStyle) {
      const body = document.body;
      const html = document.documentElement;

      // 恢复原始样式
      body.style.overflow = originalStyle.overflow;
      body.style.paddingRight = originalStyle.paddingRight;
      body.style.position = originalStyle.position;
      body.style.top = originalStyle.top;
      body.style.width = originalStyle.width;

      // 恢复 html
      html.style.overflow = "";

      // 恢复滚动位置
      window.scrollTo(0, scrollPosition);

      originalStyle = null;
    }
  };

  // 切换锁定状态
  const toggle = () => {
    if (isLocked.value) {
      unlock();
    } else {
      lock();
    }
  };

  // 立即锁定
  if (immediate) {
    lock();
  }

  // 组件卸载时解锁
  onUnmounted(() => {
    if (isLocked.value) {
      unlock();
    }
  });

  return {
    isLocked,
    lock,
    unlock,
    toggle,
  };
}

// ============================================================================
// 响应式滚动锁定
// ============================================================================

/**
 * 响应式滚动锁定
 * 根据 ref 值自动锁定/解锁
 * @example
 * const isOpen = ref(false)
 * useScrollLockWhenTrue(isOpen)
 *
 * // 打开模态框时自动锁定
 * isOpen.value = true
 */
export function useScrollLockWhenTrue(source: Ref<boolean>, options?: UseScrollLockOptions) {
  const { isLocked, lock, unlock } = useScrollLock(options);

  watch(
    source,
    (value) => {
      if (value) {
        lock();
      } else {
        unlock();
      }
    },
    { immediate: true }
  );

  return {
    isLocked,
    lock,
    unlock,
  };
}

// ============================================================================
// 元素滚动锁定
// ============================================================================

/**
 * 元素滚动锁定
 * 锁定指定元素的滚动（而非整个页面）
 * @example
 * const containerRef = ref<HTMLElement>()
 * const { lock, unlock } = useElementScrollLock(containerRef)
 */
export function useElementScrollLock(target: Ref<HTMLElement | null | undefined>) {
  const isLocked = ref(false);
  let originalOverflow = "";

  const lock = () => {
    const element = target.value;
    if (!element || isLocked.value) return;

    originalOverflow = element.style.overflow;
    element.style.overflow = "hidden";
    isLocked.value = true;
  };

  const unlock = () => {
    const element = target.value;
    if (!element || !isLocked.value) return;

    element.style.overflow = originalOverflow;
    isLocked.value = false;
  };

  const toggle = () => {
    if (isLocked.value) {
      unlock();
    } else {
      lock();
    }
  };

  onUnmounted(() => {
    if (isLocked.value) {
      unlock();
    }
  });

  return {
    isLocked,
    lock,
    unlock,
    toggle,
  };
}

// ============================================================================
// 滚动位置保存/恢复
// ============================================================================

/**
 * 滚动位置保存/恢复
 * @example
 * const { save, restore, savedPosition } = useScrollPosition()
 *
 * // 保存位置
 * save()
 *
 * // 切换页面后恢复
 * restore()
 */
export function useScrollPosition() {
  const savedPosition = ref({ x: 0, y: 0 });

  const save = () => {
    savedPosition.value = {
      x: window.scrollX,
      y: window.scrollY,
    };
  };

  const restore = (behavior: ScrollBehavior = "instant") => {
    window.scrollTo({
      left: savedPosition.value.x,
      top: savedPosition.value.y,
      behavior,
    });
  };

  const scrollTo = (position: { x?: number; y?: number }, behavior: ScrollBehavior = "smooth") => {
    window.scrollTo({
      left: position.x ?? window.scrollX,
      top: position.y ?? window.scrollY,
      behavior,
    });
  };

  const scrollToTop = (behavior: ScrollBehavior = "smooth") => {
    window.scrollTo({ top: 0, behavior });
  };

  const scrollToBottom = (behavior: ScrollBehavior = "smooth") => {
    window.scrollTo({
      top: document.documentElement.scrollHeight,
      behavior,
    });
  };

  return {
    savedPosition,
    save,
    restore,
    scrollTo,
    scrollToTop,
    scrollToBottom,
  };
}

// ============================================================================
// 滚动方向检测
// ============================================================================

/**
 * 滚动方向检测
 * @example
 * const { direction, isScrollingUp, isScrollingDown } = useScrollDirection()
 */
export function useScrollDirection() {
  const direction = ref<"up" | "down" | "none">("none");
  const isScrollingUp = ref(false);
  const isScrollingDown = ref(false);

  let lastScrollY = 0;
  let ticking = false;

  const updateDirection = () => {
    const currentScrollY = window.scrollY;

    if (currentScrollY > lastScrollY) {
      direction.value = "down";
      isScrollingUp.value = false;
      isScrollingDown.value = true;
    } else if (currentScrollY < lastScrollY) {
      direction.value = "up";
      isScrollingUp.value = true;
      isScrollingDown.value = false;
    } else {
      direction.value = "none";
      isScrollingUp.value = false;
      isScrollingDown.value = false;
    }

    lastScrollY = currentScrollY;
    ticking = false;
  };

  const handleScroll = () => {
    if (!ticking) {
      requestAnimationFrame(updateDirection);
      ticking = true;
    }
  };

  if (typeof window !== "undefined") {
    lastScrollY = window.scrollY;
    window.addEventListener("scroll", handleScroll, { passive: true });
  }

  onUnmounted(() => {
    window.removeEventListener("scroll", handleScroll);
  });

  return {
    direction,
    isScrollingUp,
    isScrollingDown,
  };
}
