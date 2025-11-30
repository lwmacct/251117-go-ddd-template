/**
 * Mouse Composable
 * 提供鼠标相关的响应式状态和事件处理
 */

import { ref, onMounted, onUnmounted, type Ref, computed, watch } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface MousePosition {
  x: number;
  y: number;
}

export interface UseMouseOptions {
  /** 目标元素，默认为 window */
  target?: Window | HTMLElement | Ref<HTMLElement | null>;
  /** 坐标类型：page（相对文档）、client（相对视口）、screen（相对屏幕） */
  type?: "page" | "client" | "screen";
  /** 是否包含触摸事件 */
  touch?: boolean;
  /** 是否立即启动 */
  immediate?: boolean;
  /** 初始位置 */
  initialValue?: MousePosition;
}

export interface UseMouseReturn {
  /** X 坐标 */
  x: Ref<number>;
  /** Y 坐标 */
  y: Ref<number>;
  /** 位置对象 */
  position: Ref<MousePosition>;
  /** 源类型（mouse 或 touch） */
  sourceType: Ref<"mouse" | "touch" | null>;
}

export interface UseMousePressedOptions {
  /** 目标元素 */
  target?: Window | HTMLElement | Ref<HTMLElement | null>;
  /** 是否包含触摸事件 */
  touch?: boolean;
}

export interface UseMousePressedReturn {
  /** 是否按下 */
  pressed: Ref<boolean>;
  /** 按下的按钮（0=左键，1=中键，2=右键） */
  button: Ref<number>;
}

// ============================================================================
// useMouse - 鼠标位置
// ============================================================================

/**
 * 跟踪鼠标位置
 * @example
 * const { x, y, position } = useMouse()
 * // x.value 和 y.value 会随鼠标移动自动更新
 */
export function useMouse(options: UseMouseOptions = {}): UseMouseReturn {
  const { target, type = "page", touch = true, immediate = true, initialValue = { x: 0, y: 0 } } = options;

  const x = ref(initialValue.x);
  const y = ref(initialValue.y);
  const sourceType = ref<"mouse" | "touch" | null>(null);

  const position = computed(() => ({
    x: x.value,
    y: y.value,
  }));

  const getTarget = (): Window | HTMLElement | null => {
    if (!target) {
      return typeof window !== "undefined" ? window : null;
    }
    if (target instanceof Window || target instanceof HTMLElement) {
      return target;
    }
    return target.value;
  };

  const updatePosition = (event: MouseEvent | TouchEvent) => {
    if (event instanceof MouseEvent) {
      sourceType.value = "mouse";
      if (type === "page") {
        x.value = event.pageX;
        y.value = event.pageY;
      } else if (type === "client") {
        x.value = event.clientX;
        y.value = event.clientY;
      } else if (type === "screen") {
        x.value = event.screenX;
        y.value = event.screenY;
      }
    } else if (event instanceof TouchEvent && event.touches.length > 0) {
      sourceType.value = "touch";
      const touchEvent = event.touches[0];
      if (type === "page") {
        x.value = touchEvent.pageX;
        y.value = touchEvent.pageY;
      } else if (type === "client") {
        x.value = touchEvent.clientX;
        y.value = touchEvent.clientY;
      } else if (type === "screen") {
        x.value = touchEvent.screenX;
        y.value = touchEvent.screenY;
      }
    }
  };

  let cleanup: (() => void) | null = null;

  const start = () => {
    const el = getTarget();
    if (!el) return;

    el.addEventListener("mousemove", updatePosition as EventListener);
    if (touch) {
      el.addEventListener("touchmove", updatePosition as EventListener);
      el.addEventListener("touchstart", updatePosition as EventListener);
    }

    cleanup = () => {
      el.removeEventListener("mousemove", updatePosition as EventListener);
      if (touch) {
        el.removeEventListener("touchmove", updatePosition as EventListener);
        el.removeEventListener("touchstart", updatePosition as EventListener);
      }
    };
  };

  if (immediate) {
    onMounted(start);
  }

  onUnmounted(() => {
    cleanup?.();
  });

  return {
    x,
    y,
    position,
    sourceType,
  };
}

// ============================================================================
// useMousePressed - 鼠标按下状态
// ============================================================================

/**
 * 跟踪鼠标按下状态
 * @example
 * const { pressed, button } = useMousePressed()
 * // pressed.value 为 true 时表示鼠标按下
 */
export function useMousePressed(options: UseMousePressedOptions = {}): UseMousePressedReturn {
  const { target, touch = true } = options;

  const pressed = ref(false);
  const button = ref(-1);

  const getTarget = (): Window | HTMLElement | null => {
    if (!target) {
      return typeof window !== "undefined" ? window : null;
    }
    if (target instanceof Window || target instanceof HTMLElement) {
      return target;
    }
    return target.value;
  };

  const handleMouseDown = (event: MouseEvent) => {
    pressed.value = true;
    button.value = event.button;
  };

  const handleMouseUp = () => {
    pressed.value = false;
    button.value = -1;
  };

  const handleTouchStart = () => {
    pressed.value = true;
    button.value = 0;
  };

  const handleTouchEnd = () => {
    pressed.value = false;
    button.value = -1;
  };

  onMounted(() => {
    const el = getTarget();
    if (!el) return;

    el.addEventListener("mousedown", handleMouseDown as EventListener);
    el.addEventListener("mouseup", handleMouseUp);
    el.addEventListener("mouseleave", handleMouseUp);

    if (touch) {
      el.addEventListener("touchstart", handleTouchStart);
      el.addEventListener("touchend", handleTouchEnd);
      el.addEventListener("touchcancel", handleTouchEnd);
    }
  });

  onUnmounted(() => {
    const el = getTarget();
    if (!el) return;

    el.removeEventListener("mousedown", handleMouseDown as EventListener);
    el.removeEventListener("mouseup", handleMouseUp);
    el.removeEventListener("mouseleave", handleMouseUp);

    if (touch) {
      el.removeEventListener("touchstart", handleTouchStart);
      el.removeEventListener("touchend", handleTouchEnd);
      el.removeEventListener("touchcancel", handleTouchEnd);
    }
  });

  return {
    pressed,
    button,
  };
}

// ============================================================================
// useMouseInElement - 元素内鼠标位置
// ============================================================================

export interface UseMouseInElementOptions {
  /** 是否包含触摸事件 */
  touch?: boolean;
}

export interface UseMouseInElementReturn {
  /** 相对于元素的 X 坐标 */
  x: Ref<number>;
  /** 相对于元素的 Y 坐标 */
  y: Ref<number>;
  /** 是否在元素内部 */
  isOutside: Ref<boolean>;
  /** 元素左边距 */
  elementX: Ref<number>;
  /** 元素上边距 */
  elementY: Ref<number>;
  /** 元素宽度 */
  elementWidth: Ref<number>;
  /** 元素高度 */
  elementHeight: Ref<number>;
}

/**
 * 跟踪鼠标在元素内的相对位置
 * @example
 * const target = ref<HTMLElement | null>(null)
 * const { x, y, isOutside } = useMouseInElement(target)
 */
export function useMouseInElement(
  target: Ref<HTMLElement | null>,
  options: UseMouseInElementOptions = {}
): UseMouseInElementReturn {
  const { touch = true } = options;

  const x = ref(0);
  const y = ref(0);
  const isOutside = ref(true);
  const elementX = ref(0);
  const elementY = ref(0);
  const elementWidth = ref(0);
  const elementHeight = ref(0);

  const updatePosition = (event: MouseEvent | TouchEvent) => {
    const el = target.value;
    if (!el) return;

    const rect = el.getBoundingClientRect();
    elementX.value = rect.left;
    elementY.value = rect.top;
    elementWidth.value = rect.width;
    elementHeight.value = rect.height;

    let clientX: number;
    let clientY: number;

    if (event instanceof MouseEvent) {
      clientX = event.clientX;
      clientY = event.clientY;
    } else if (event instanceof TouchEvent && event.touches.length > 0) {
      clientX = event.touches[0].clientX;
      clientY = event.touches[0].clientY;
    } else {
      return;
    }

    x.value = clientX - rect.left;
    y.value = clientY - rect.top;

    isOutside.value = x.value < 0 || y.value < 0 || x.value > rect.width || y.value > rect.height;
  };

  const handleMouseEnter = () => {
    isOutside.value = false;
  };

  const handleMouseLeave = () => {
    isOutside.value = true;
  };

  watch(
    target,
    (el) => {
      if (!el) return;

      el.addEventListener("mousemove", updatePosition as EventListener);
      el.addEventListener("mouseenter", handleMouseEnter);
      el.addEventListener("mouseleave", handleMouseLeave);

      if (touch) {
        el.addEventListener("touchmove", updatePosition as EventListener);
        el.addEventListener("touchstart", updatePosition as EventListener);
      }
    },
    { immediate: true }
  );

  onUnmounted(() => {
    const el = target.value;
    if (!el) return;

    el.removeEventListener("mousemove", updatePosition as EventListener);
    el.removeEventListener("mouseenter", handleMouseEnter);
    el.removeEventListener("mouseleave", handleMouseLeave);

    if (touch) {
      el.removeEventListener("touchmove", updatePosition as EventListener);
      el.removeEventListener("touchstart", updatePosition as EventListener);
    }
  });

  return {
    x,
    y,
    isOutside,
    elementX,
    elementY,
    elementWidth,
    elementHeight,
  };
}

// ============================================================================
// useHover - 悬停状态
// ============================================================================

export interface UseHoverOptions {
  /** 进入延迟（毫秒） */
  delayEnter?: number;
  /** 离开延迟（毫秒） */
  delayLeave?: number;
}

export interface UseHoverReturn {
  /** 是否悬停 */
  isHovered: Ref<boolean>;
}

/**
 * 跟踪元素悬停状态
 * @example
 * const target = ref<HTMLElement | null>(null)
 * const { isHovered } = useHover(target)
 */
export function useHover(target: Ref<HTMLElement | null>, options: UseHoverOptions = {}): UseHoverReturn {
  const { delayEnter = 0, delayLeave = 0 } = options;

  const isHovered = ref(false);
  let enterTimer: ReturnType<typeof setTimeout> | null = null;
  let leaveTimer: ReturnType<typeof setTimeout> | null = null;

  const clearTimers = () => {
    if (enterTimer) {
      clearTimeout(enterTimer);
      enterTimer = null;
    }
    if (leaveTimer) {
      clearTimeout(leaveTimer);
      leaveTimer = null;
    }
  };

  const handleMouseEnter = () => {
    clearTimers();
    if (delayEnter > 0) {
      enterTimer = setTimeout(() => {
        isHovered.value = true;
      }, delayEnter);
    } else {
      isHovered.value = true;
    }
  };

  const handleMouseLeave = () => {
    clearTimers();
    if (delayLeave > 0) {
      leaveTimer = setTimeout(() => {
        isHovered.value = false;
      }, delayLeave);
    } else {
      isHovered.value = false;
    }
  };

  watch(
    target,
    (el, _, onCleanup) => {
      if (!el) return;

      el.addEventListener("mouseenter", handleMouseEnter);
      el.addEventListener("mouseleave", handleMouseLeave);

      onCleanup(() => {
        el.removeEventListener("mouseenter", handleMouseEnter);
        el.removeEventListener("mouseleave", handleMouseLeave);
        clearTimers();
      });
    },
    { immediate: true }
  );

  onUnmounted(clearTimers);

  return {
    isHovered,
  };
}

// ============================================================================
// useCursor - 光标样式
// ============================================================================

export type CursorType =
  | "auto"
  | "default"
  | "none"
  | "context-menu"
  | "help"
  | "pointer"
  | "progress"
  | "wait"
  | "cell"
  | "crosshair"
  | "text"
  | "vertical-text"
  | "alias"
  | "copy"
  | "move"
  | "no-drop"
  | "not-allowed"
  | "grab"
  | "grabbing"
  | "all-scroll"
  | "col-resize"
  | "row-resize"
  | "n-resize"
  | "e-resize"
  | "s-resize"
  | "w-resize"
  | "ne-resize"
  | "nw-resize"
  | "se-resize"
  | "sw-resize"
  | "ew-resize"
  | "ns-resize"
  | "nesw-resize"
  | "nwse-resize"
  | "zoom-in"
  | "zoom-out";

export interface UseCursorReturn {
  /** 当前光标类型 */
  cursor: Ref<CursorType>;
  /** 设置光标 */
  setCursor: (cursor: CursorType) => void;
  /** 重置光标 */
  resetCursor: () => void;
}

/**
 * 控制光标样式
 * @example
 * const { cursor, setCursor, resetCursor } = useCursor()
 * setCursor('pointer')
 * // 或者直接修改 cursor.value = 'grab'
 */
export function useCursor(initialValue: CursorType = "auto"): UseCursorReturn {
  const cursor = ref<CursorType>(initialValue);
  const originalCursor = ref<string>("");

  const setCursor = (value: CursorType) => {
    cursor.value = value;
  };

  const resetCursor = () => {
    cursor.value = initialValue;
  };

  watch(
    cursor,
    (value) => {
      if (typeof document !== "undefined") {
        if (!originalCursor.value) {
          originalCursor.value = document.body.style.cursor;
        }
        document.body.style.cursor = value;
      }
    },
    { immediate: true }
  );

  onUnmounted(() => {
    if (typeof document !== "undefined" && originalCursor.value !== undefined) {
      document.body.style.cursor = originalCursor.value;
    }
  });

  return {
    cursor,
    setCursor,
    resetCursor,
  };
}

// ============================================================================
// useDropZone - 拖放区域
// ============================================================================

export interface UseDropZoneOptions {
  /** 是否阻止默认行为 */
  preventDefault?: boolean;
  /** 接受的文件类型 */
  accept?: string[];
  /** 拖入回调 */
  onDragEnter?: (event: DragEvent) => void;
  /** 拖出回调 */
  onDragLeave?: (event: DragEvent) => void;
  /** 拖放回调 */
  onDrop?: (files: File[], event: DragEvent) => void;
}

export interface UseDropZoneReturn {
  /** 是否正在拖放 */
  isOverDropZone: Ref<boolean>;
  /** 拖放的文件 */
  files: Ref<File[]>;
}

/**
 * 创建拖放区域
 * @example
 * const target = ref<HTMLElement | null>(null)
 * const { isOverDropZone, files } = useDropZone(target, {
 *   onDrop: (files) => uploadFiles(files)
 * })
 */
export function useDropZone(target: Ref<HTMLElement | null>, options: UseDropZoneOptions = {}): UseDropZoneReturn {
  const { preventDefault = true, accept, onDragEnter, onDragLeave, onDrop } = options;

  const isOverDropZone = ref(false);
  const files = ref<File[]>([]);
  let counter = 0;

  const filterFiles = (fileList: FileList): File[] => {
    const filtered = Array.from(fileList);

    if (!accept || accept.length === 0) {
      return filtered;
    }

    return filtered.filter((file) => {
      return accept.some((type) => {
        if (type.startsWith(".")) {
          return file.name.endsWith(type);
        }
        if (type.endsWith("/*")) {
          return file.type.startsWith(type.slice(0, -1));
        }
        return file.type === type;
      });
    });
  };

  const handleDragEnter = (event: DragEvent) => {
    if (preventDefault) {
      event.preventDefault();
    }
    counter++;
    isOverDropZone.value = true;
    onDragEnter?.(event);
  };

  const handleDragLeave = (event: DragEvent) => {
    if (preventDefault) {
      event.preventDefault();
    }
    counter--;
    if (counter === 0) {
      isOverDropZone.value = false;
      onDragLeave?.(event);
    }
  };

  const handleDragOver = (event: DragEvent) => {
    if (preventDefault) {
      event.preventDefault();
    }
  };

  const handleDrop = (event: DragEvent) => {
    if (preventDefault) {
      event.preventDefault();
    }
    counter = 0;
    isOverDropZone.value = false;

    if (event.dataTransfer?.files) {
      files.value = filterFiles(event.dataTransfer.files);
      onDrop?.(files.value, event);
    }
  };

  watch(
    target,
    (el, _, onCleanup) => {
      if (!el) return;

      el.addEventListener("dragenter", handleDragEnter);
      el.addEventListener("dragleave", handleDragLeave);
      el.addEventListener("dragover", handleDragOver);
      el.addEventListener("drop", handleDrop);

      onCleanup(() => {
        el.removeEventListener("dragenter", handleDragEnter);
        el.removeEventListener("dragleave", handleDragLeave);
        el.removeEventListener("dragover", handleDragOver);
        el.removeEventListener("drop", handleDrop);
      });
    },
    { immediate: true }
  );

  return {
    isOverDropZone,
    files,
  };
}
