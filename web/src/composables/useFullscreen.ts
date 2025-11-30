/**
 * 全屏 API Composable
 * 控制元素或文档进入/退出全屏模式
 */

import { ref, computed, onMounted, onUnmounted, type Ref } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UseFullscreenOptions {
  /** 是否在组件卸载时自动退出全屏，默认 true */
  autoExitOnUnmount?: boolean;
}

export interface UseFullscreenReturn {
  /** 是否处于全屏模式 */
  isFullscreen: Ref<boolean>;
  /** 是否支持全屏 */
  isSupported: Ref<boolean>;
  /** 当前全屏的元素 */
  fullscreenElement: Ref<Element | null>;
  /** 进入全屏 */
  enter: () => Promise<void>;
  /** 退出全屏 */
  exit: () => Promise<void>;
  /** 切换全屏 */
  toggle: () => Promise<void>;
}

// ============================================================================
// 浏览器兼容性处理
// ============================================================================

interface FullscreenDocument extends Document {
  mozFullScreenElement?: Element;
  webkitFullscreenElement?: Element;
  msFullscreenElement?: Element;
  mozCancelFullScreen?: () => Promise<void>;
  webkitExitFullscreen?: () => Promise<void>;
  msExitFullscreen?: () => Promise<void>;
}

interface FullscreenElement extends Element {
  mozRequestFullScreen?: () => Promise<void>;
  webkitRequestFullscreen?: () => Promise<void>;
  msRequestFullscreen?: () => Promise<void>;
}

// 获取全屏元素
function getFullscreenElement(): Element | null {
  const doc = document as FullscreenDocument;
  return (
    doc.fullscreenElement ||
    doc.mozFullScreenElement ||
    doc.webkitFullscreenElement ||
    doc.msFullscreenElement ||
    null
  );
}

// 请求全屏
async function requestFullscreen(element: Element): Promise<void> {
  const el = element as FullscreenElement;

  if (el.requestFullscreen) {
    return el.requestFullscreen();
  } else if (el.mozRequestFullScreen) {
    return el.mozRequestFullScreen();
  } else if (el.webkitRequestFullscreen) {
    return el.webkitRequestFullscreen();
  } else if (el.msRequestFullscreen) {
    return el.msRequestFullscreen();
  }

  throw new Error("Fullscreen API is not supported");
}

// 退出全屏
async function exitFullscreen(): Promise<void> {
  const doc = document as FullscreenDocument;

  if (doc.exitFullscreen) {
    return doc.exitFullscreen();
  } else if (doc.mozCancelFullScreen) {
    return doc.mozCancelFullScreen();
  } else if (doc.webkitExitFullscreen) {
    return doc.webkitExitFullscreen();
  } else if (doc.msExitFullscreen) {
    return doc.msExitFullscreen();
  }

  throw new Error("Fullscreen API is not supported");
}

// 检查是否支持全屏
function isFullscreenSupported(): boolean {
  const doc = document as FullscreenDocument;
  return !!(
    doc.fullscreenEnabled ||
    (doc as FullscreenDocument & { mozFullScreenEnabled?: boolean }).mozFullScreenEnabled ||
    (doc as FullscreenDocument & { webkitFullscreenEnabled?: boolean }).webkitFullscreenEnabled ||
    (doc as FullscreenDocument & { msFullscreenEnabled?: boolean }).msFullscreenEnabled
  );
}

// 全屏变化事件名
function getFullscreenChangeEvent(): string {
  if ("onfullscreenchange" in document) {
    return "fullscreenchange";
  } else if ("onmozfullscreenchange" in document) {
    return "mozfullscreenchange";
  } else if ("onwebkitfullscreenchange" in document) {
    return "webkitfullscreenchange";
  } else if ("onmsfullscreenchange" in document) {
    return "MSFullscreenChange";
  }
  return "fullscreenchange";
}

// ============================================================================
// 主函数
// ============================================================================

/**
 * 全屏控制
 * @example
 * // 元素全屏
 * const videoRef = ref<HTMLVideoElement>()
 * const { isFullscreen, toggle } = useFullscreen(videoRef)
 *
 * // 文档全屏
 * const { isFullscreen, toggle } = useFullscreen()
 */
export function useFullscreen(
  target?: Ref<Element | null | undefined>,
  options: UseFullscreenOptions = {}
): UseFullscreenReturn {
  const { autoExitOnUnmount = true } = options;

  const isFullscreen = ref(false);
  const fullscreenElement = ref<Element | null>(null);
  const isSupported = ref(isFullscreenSupported());

  // 获取目标元素
  const getTarget = (): Element => {
    return target?.value || document.documentElement;
  };

  // 更新状态
  const updateState = () => {
    const element = getFullscreenElement();
    isFullscreen.value = element !== null;
    fullscreenElement.value = element;
  };

  // 进入全屏
  const enter = async (): Promise<void> => {
    if (!isSupported.value) {
      console.warn("Fullscreen API is not supported");
      return;
    }

    const element = getTarget();
    await requestFullscreen(element);
  };

  // 退出全屏
  const exit = async (): Promise<void> => {
    if (!isSupported.value) {
      console.warn("Fullscreen API is not supported");
      return;
    }

    if (isFullscreen.value) {
      await exitFullscreen();
    }
  };

  // 切换全屏
  const toggle = async (): Promise<void> => {
    if (isFullscreen.value) {
      await exit();
    } else {
      await enter();
    }
  };

  // 全屏变化事件
  const eventName = getFullscreenChangeEvent();

  onMounted(() => {
    updateState();
    document.addEventListener(eventName, updateState);
  });

  onUnmounted(() => {
    document.removeEventListener(eventName, updateState);

    // 自动退出全屏
    if (autoExitOnUnmount && isFullscreen.value) {
      exit();
    }
  });

  return {
    isFullscreen,
    isSupported,
    fullscreenElement,
    enter,
    exit,
    toggle,
  };
}

// ============================================================================
// 文档全屏简化版
// ============================================================================

/**
 * 文档全屏（简化版）
 * @example
 * const { isFullscreen, toggle } = useDocumentFullscreen()
 */
export function useDocumentFullscreen(options?: UseFullscreenOptions) {
  return useFullscreen(undefined, options);
}

// ============================================================================
// 全屏按钮组件数据
// ============================================================================

/**
 * 全屏按钮数据
 * @example
 * const { icon, tooltip, toggle } = useFullscreenButton()
 */
export function useFullscreenButton(
  target?: Ref<Element | null | undefined>,
  options?: UseFullscreenOptions
) {
  const { isFullscreen, isSupported, toggle, enter, exit } = useFullscreen(
    target,
    options
  );

  const icon = computed(() =>
    isFullscreen.value ? "mdi-fullscreen-exit" : "mdi-fullscreen"
  );

  const tooltip = computed(() =>
    isFullscreen.value ? "退出全屏" : "全屏显示"
  );

  return {
    isFullscreen,
    isSupported,
    icon,
    tooltip,
    toggle,
    enter,
    exit,
  };
}
