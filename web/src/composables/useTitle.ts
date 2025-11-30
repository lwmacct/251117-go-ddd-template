/**
 * Title Composable
 * 提供文档标题的响应式管理
 */

import {
  ref,
  watch,
  computed,
  onMounted,
  onUnmounted,
  type Ref,
  type ComputedRef,
  type MaybeRef,
} from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UseTitleOptions {
  /** 标题模板，%s 会被替换为标题内容 */
  template?: string;
  /** 是否在组件卸载时恢复原标题 */
  restoreOnUnmount?: boolean;
  /** 是否观察标题变化 */
  observe?: boolean;
}

export interface UseTitleReturn {
  /** 当前标题 */
  title: Ref<string>;
  /** 是否支持 */
  isSupported: ComputedRef<boolean>;
}

export interface UseTitleTemplateOptions {
  /** 分隔符 */
  separator?: string;
  /** 网站名称 */
  siteName?: string;
  /** 网站名称位置 */
  siteNamePosition?: "prefix" | "suffix";
}

export interface UseTitleTemplateReturn {
  /** 页面标题 */
  pageTitle: Ref<string>;
  /** 完整标题 */
  fullTitle: ComputedRef<string>;
  /** 设置页面标题 */
  setPageTitle: (title: string) => void;
}

export interface UseFaviconOptions {
  /** 是否在组件卸载时恢复原图标 */
  restoreOnUnmount?: boolean;
}

export interface UseFaviconReturn {
  /** 当前图标 URL */
  favicon: Ref<string | null>;
  /** 是否支持 */
  isSupported: ComputedRef<boolean>;
}

// ============================================================================
// useTitle - 文档标题管理
// ============================================================================

/**
 * 响应式文档标题管理
 * @example
 * const { title } = useTitle('My Page')
 * title.value = 'New Title'
 *
 * // 使用模板
 * const { title } = useTitle('Home', { template: '%s | My App' })
 * // 结果: 'Home | My App'
 */
export function useTitle(
  newTitle?: MaybeRef<string | null | undefined>,
  options: UseTitleOptions = {}
): UseTitleReturn {
  const { template, restoreOnUnmount = false, observe = false } = options;

  const isSupported = computed(
    () => typeof document !== "undefined" && "title" in document
  );

  // 保存原始标题
  const originalTitle =
    typeof document !== "undefined" ? document.title : "";

  // 创建响应式标题
  const title = ref(
    typeof newTitle === "object" && newTitle !== null && "value" in newTitle
      ? (newTitle as Ref<string | null | undefined>).value ?? originalTitle
      : (newTitle as string | null | undefined) ?? originalTitle
  );

  // 应用标题到文档
  const applyTitle = (value: string) => {
    if (!isSupported.value) return;

    const finalTitle = template ? template.replace("%s", value) : value;
    document.title = finalTitle;
  };

  // 监听标题变化
  watch(
    () =>
      typeof newTitle === "object" && newTitle !== null && "value" in newTitle
        ? (newTitle as Ref<string | null | undefined>).value
        : title.value,
    (value) => {
      if (value !== null && value !== undefined) {
        title.value = value;
        applyTitle(value);
      }
    },
    { immediate: true }
  );

  // 监听外部标题变化
  if (observe && isSupported.value) {
    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (
          mutation.type === "childList" &&
          mutation.target === document.querySelector("title")
        ) {
          title.value = document.title;
        }
      });
    });

    onMounted(() => {
      const titleElement = document.querySelector("title");
      if (titleElement) {
        observer.observe(titleElement, { childList: true });
      }
    });

    onUnmounted(() => {
      observer.disconnect();
    });
  }

  // 组件卸载时恢复原标题
  if (restoreOnUnmount) {
    onUnmounted(() => {
      if (isSupported.value) {
        document.title = originalTitle;
      }
    });
  }

  return {
    title,
    isSupported,
  };
}

// ============================================================================
// useTitleTemplate - 带模板的标题管理
// ============================================================================

/**
 * 带模板的标题管理
 * @example
 * const { pageTitle, fullTitle, setPageTitle } = useTitleTemplate({
 *   siteName: 'My App',
 *   separator: ' - ',
 *   siteNamePosition: 'suffix'
 * })
 *
 * setPageTitle('Home')
 * // fullTitle.value: 'Home - My App'
 */
export function useTitleTemplate(
  options: UseTitleTemplateOptions = {}
): UseTitleTemplateReturn {
  const {
    separator = " | ",
    siteName = "",
    siteNamePosition = "suffix",
  } = options;

  const pageTitle = ref("");

  const fullTitle = computed(() => {
    if (!siteName) return pageTitle.value;
    if (!pageTitle.value) return siteName;

    return siteNamePosition === "prefix"
      ? `${siteName}${separator}${pageTitle.value}`
      : `${pageTitle.value}${separator}${siteName}`;
  });

  // 应用到文档
  watch(fullTitle, (value) => {
    if (typeof document !== "undefined") {
      document.title = value;
    }
  });

  const setPageTitle = (title: string) => {
    pageTitle.value = title;
  };

  return {
    pageTitle,
    fullTitle,
    setPageTitle,
  };
}

// ============================================================================
// useDocumentTitle - 简化的文档标题 Hook
// ============================================================================

/**
 * 简化的文档标题设置
 * @example
 * useDocumentTitle('My Page')
 *
 * // 响应式
 * const title = ref('My Page')
 * useDocumentTitle(title)
 */
export function useDocumentTitle(title: MaybeRef<string>): void {
  watch(
    () => (typeof title === "object" && "value" in title ? title.value : title),
    (value) => {
      if (typeof document !== "undefined") {
        document.title = value;
      }
    },
    { immediate: true }
  );
}

// ============================================================================
// useFavicon - Favicon 管理
// ============================================================================

/**
 * 响应式 Favicon 管理
 * @example
 * const { favicon } = useFavicon()
 * favicon.value = '/new-favicon.ico'
 *
 * // 或者直接设置
 * useFavicon('/favicon.ico')
 */
export function useFavicon(
  newIcon?: MaybeRef<string | null | undefined>,
  options: UseFaviconOptions = {}
): UseFaviconReturn {
  const { restoreOnUnmount = false } = options;

  const isSupported = computed(() => typeof document !== "undefined");

  // 获取当前 favicon
  const getCurrentFavicon = (): string | null => {
    if (!isSupported.value) return null;
    const link = document.querySelector<HTMLLinkElement>(
      'link[rel="icon"], link[rel="shortcut icon"]'
    );
    return link?.href ?? null;
  };

  // 保存原始 favicon
  const originalFavicon = getCurrentFavicon();

  // 创建响应式 favicon
  const favicon = ref<string | null>(
    typeof newIcon === "object" && newIcon !== null && "value" in newIcon
      ? (newIcon as Ref<string | null | undefined>).value ?? null
      : (newIcon as string | null | undefined) ?? null
  );

  // 应用 favicon
  const applyFavicon = (url: string | null) => {
    if (!isSupported.value || !url) return;

    let link = document.querySelector<HTMLLinkElement>('link[rel="icon"]');

    if (!link) {
      link = document.createElement("link");
      link.rel = "icon";
      document.head.appendChild(link);
    }

    link.href = url;
  };

  // 监听变化
  watch(
    () =>
      typeof newIcon === "object" && newIcon !== null && "value" in newIcon
        ? (newIcon as Ref<string | null | undefined>).value
        : favicon.value,
    (value) => {
      if (value !== undefined) {
        favicon.value = value;
        applyFavicon(value);
      }
    },
    { immediate: true }
  );

  // 组件卸载时恢复
  if (restoreOnUnmount) {
    onUnmounted(() => {
      if (originalFavicon) {
        applyFavicon(originalFavicon);
      }
    });
  }

  return {
    favicon,
    isSupported,
  };
}

// ============================================================================
// usePageLeave - 页面离开检测
// ============================================================================

export interface UsePageLeaveOptions {
  /** 是否触发回调 */
  enabled?: boolean;
}

export interface UsePageLeaveReturn {
  /** 是否已离开页面 */
  isLeft: Ref<boolean>;
}

/**
 * 检测用户是否离开页面
 * @example
 * const { isLeft } = usePageLeave()
 *
 * watch(isLeft, (left) => {
 *   if (left) {
 *     // 暂停视频、保存草稿等
 *   }
 * })
 */
export function usePageLeave(
  callback?: () => void,
  options: UsePageLeaveOptions = {}
): UsePageLeaveReturn {
  const { enabled = true } = options;

  const isLeft = ref(false);

  if (typeof document === "undefined") {
    return { isLeft };
  }

  const handleVisibilityChange = () => {
    isLeft.value = document.hidden;
    if (enabled && document.hidden && callback) {
      callback();
    }
  };

  const handleMouseLeave = (event: MouseEvent) => {
    // 鼠标离开文档区域
    if (
      event.clientY <= 0 ||
      event.clientX <= 0 ||
      event.clientX >= window.innerWidth ||
      event.clientY >= window.innerHeight
    ) {
      isLeft.value = true;
      if (enabled && callback) {
        callback();
      }
    }
  };

  const handleMouseEnter = () => {
    isLeft.value = false;
  };

  onMounted(() => {
    document.addEventListener("visibilitychange", handleVisibilityChange);
    document.addEventListener("mouseleave", handleMouseLeave);
    document.addEventListener("mouseenter", handleMouseEnter);
  });

  onUnmounted(() => {
    document.removeEventListener("visibilitychange", handleVisibilityChange);
    document.removeEventListener("mouseleave", handleMouseLeave);
    document.removeEventListener("mouseenter", handleMouseEnter);
  });

  return {
    isLeft,
  };
}

// ============================================================================
// useDocumentVisibility - 文档可见性
// ============================================================================

export type DocumentVisibilityState = "visible" | "hidden" | "prerender";

export interface UseDocumentVisibilityReturn {
  /** 文档可见状态 */
  visibility: Ref<DocumentVisibilityState>;
  /** 是否可见 */
  isVisible: ComputedRef<boolean>;
  /** 是否支持 */
  isSupported: ComputedRef<boolean>;
}

/**
 * 检测文档可见性状态
 * @example
 * const { visibility, isVisible } = useDocumentVisibility()
 *
 * watch(isVisible, (visible) => {
 *   if (visible) {
 *     // 恢复动画、重新连接等
 *   } else {
 *     // 暂停动画、断开连接等
 *   }
 * })
 */
export function useDocumentVisibility(): UseDocumentVisibilityReturn {
  const isSupported = computed(
    () => typeof document !== "undefined" && "visibilityState" in document
  );

  const visibility = ref<DocumentVisibilityState>(
    isSupported.value
      ? (document.visibilityState as DocumentVisibilityState)
      : "visible"
  );

  const isVisible = computed(() => visibility.value === "visible");

  if (isSupported.value) {
    const handleVisibilityChange = () => {
      visibility.value = document.visibilityState as DocumentVisibilityState;
    };

    onMounted(() => {
      document.addEventListener("visibilitychange", handleVisibilityChange);
    });

    onUnmounted(() => {
      document.removeEventListener("visibilitychange", handleVisibilityChange);
    });
  }

  return {
    visibility,
    isVisible,
    isSupported,
  };
}

// ============================================================================
// useHead - 简单的 Head 管理
// ============================================================================

export interface HeadConfig {
  /** 页面标题 */
  title?: string;
  /** 元标签 */
  meta?: Array<{
    name?: string;
    property?: string;
    content: string;
  }>;
  /** 链接标签 */
  link?: Array<{
    rel: string;
    href: string;
    type?: string;
  }>;
}

/**
 * 简单的 Head 管理
 * @example
 * useHead({
 *   title: 'My Page',
 *   meta: [
 *     { name: 'description', content: 'Page description' },
 *     { property: 'og:title', content: 'My Page' }
 *   ],
 *   link: [
 *     { rel: 'canonical', href: 'https://example.com/page' }
 *   ]
 * })
 */
export function useHead(config: HeadConfig): void {
  const createdElements: HTMLElement[] = [];

  const applyHead = () => {
    if (typeof document === "undefined") return;

    // 设置标题
    if (config.title) {
      document.title = config.title;
    }

    // 添加 meta 标签
    if (config.meta) {
      config.meta.forEach((meta) => {
        const el = document.createElement("meta");
        if (meta.name) el.name = meta.name;
        if (meta.property) el.setAttribute("property", meta.property);
        el.content = meta.content;
        document.head.appendChild(el);
        createdElements.push(el);
      });
    }

    // 添加 link 标签
    if (config.link) {
      config.link.forEach((link) => {
        const el = document.createElement("link");
        el.rel = link.rel;
        el.href = link.href;
        if (link.type) el.type = link.type;
        document.head.appendChild(el);
        createdElements.push(el);
      });
    }
  };

  const cleanup = () => {
    createdElements.forEach((el) => {
      el.parentNode?.removeChild(el);
    });
    createdElements.length = 0;
  };

  onMounted(applyHead);
  onUnmounted(cleanup);
}

// ============================================================================
// useScript - 动态脚本加载
// ============================================================================

export interface UseScriptOptions {
  /** 是否立即加载 */
  immediate?: boolean;
  /** 是否异步加载 */
  async?: boolean;
  /** 是否延迟加载 */
  defer?: boolean;
  /** 脚本类型 */
  type?: string;
  /** 跨域设置 */
  crossOrigin?: "anonymous" | "use-credentials";
  /** 完整性校验 */
  integrity?: string;
  /** 引用策略 */
  referrerPolicy?: string;
  /** 是否在组件卸载时移除 */
  removeOnUnmount?: boolean;
}

export interface UseScriptReturn {
  /** 是否正在加载 */
  isLoading: Ref<boolean>;
  /** 是否加载完成 */
  isLoaded: Ref<boolean>;
  /** 加载错误 */
  error: Ref<Error | null>;
  /** 手动加载脚本 */
  load: () => Promise<HTMLScriptElement>;
  /** 移除脚本 */
  unload: () => void;
}

/**
 * 动态加载脚本
 * @example
 * const { isLoaded, load } = useScript('https://example.com/script.js', {
 *   immediate: true,
 *   async: true
 * })
 *
 * // 或手动加载
 * await load()
 */
export function useScript(
  src: string,
  options: UseScriptOptions = {}
): UseScriptReturn {
  const {
    immediate = true,
    async = true,
    defer = false,
    type,
    crossOrigin,
    integrity,
    referrerPolicy,
    removeOnUnmount = false,
  } = options;

  const isLoading = ref(false);
  const isLoaded = ref(false);
  const error = ref<Error | null>(null);

  let scriptElement: HTMLScriptElement | null = null;

  const load = (): Promise<HTMLScriptElement> => {
    return new Promise((resolve, reject) => {
      if (typeof document === "undefined") {
        reject(new Error("Document is not available"));
        return;
      }

      // 检查是否已存在
      const existing = document.querySelector<HTMLScriptElement>(
        `script[src="${src}"]`
      );
      if (existing) {
        scriptElement = existing;
        isLoaded.value = true;
        resolve(existing);
        return;
      }

      isLoading.value = true;
      error.value = null;

      const script = document.createElement("script");
      script.src = src;
      script.async = async;
      script.defer = defer;
      if (type) script.type = type;
      if (crossOrigin) script.crossOrigin = crossOrigin;
      if (integrity) script.integrity = integrity;
      if (referrerPolicy) script.referrerPolicy = referrerPolicy;

      script.onload = () => {
        isLoading.value = false;
        isLoaded.value = true;
        scriptElement = script;
        resolve(script);
      };

      script.onerror = (e) => {
        isLoading.value = false;
        const err = new Error(`Failed to load script: ${src}`);
        error.value = err;
        reject(err);
      };

      document.head.appendChild(script);
    });
  };

  const unload = () => {
    if (scriptElement && scriptElement.parentNode) {
      scriptElement.parentNode.removeChild(scriptElement);
      scriptElement = null;
      isLoaded.value = false;
    }
  };

  if (immediate) {
    onMounted(() => {
      load().catch(() => {});
    });
  }

  if (removeOnUnmount) {
    onUnmounted(unload);
  }

  return {
    isLoading,
    isLoaded,
    error,
    load,
    unload,
  };
}

// ============================================================================
// useStylesheet - 动态样式表加载
// ============================================================================

export interface UseStylesheetOptions {
  /** 是否立即加载 */
  immediate?: boolean;
  /** 媒体查询 */
  media?: string;
  /** 是否在组件卸载时移除 */
  removeOnUnmount?: boolean;
}

export interface UseStylesheetReturn {
  /** 是否正在加载 */
  isLoading: Ref<boolean>;
  /** 是否加载完成 */
  isLoaded: Ref<boolean>;
  /** 加载错误 */
  error: Ref<Error | null>;
  /** 手动加载 */
  load: () => Promise<HTMLLinkElement>;
  /** 移除样式表 */
  unload: () => void;
}

/**
 * 动态加载样式表
 * @example
 * const { isLoaded, load } = useStylesheet('/styles/theme.css')
 *
 * // 暗色主题
 * useStylesheet('/styles/dark.css', {
 *   media: '(prefers-color-scheme: dark)'
 * })
 */
export function useStylesheet(
  href: string,
  options: UseStylesheetOptions = {}
): UseStylesheetReturn {
  const { immediate = true, media, removeOnUnmount = false } = options;

  const isLoading = ref(false);
  const isLoaded = ref(false);
  const error = ref<Error | null>(null);

  let linkElement: HTMLLinkElement | null = null;

  const load = (): Promise<HTMLLinkElement> => {
    return new Promise((resolve, reject) => {
      if (typeof document === "undefined") {
        reject(new Error("Document is not available"));
        return;
      }

      // 检查是否已存在
      const existing = document.querySelector<HTMLLinkElement>(
        `link[href="${href}"]`
      );
      if (existing) {
        linkElement = existing;
        isLoaded.value = true;
        resolve(existing);
        return;
      }

      isLoading.value = true;
      error.value = null;

      const link = document.createElement("link");
      link.rel = "stylesheet";
      link.href = href;
      if (media) link.media = media;

      link.onload = () => {
        isLoading.value = false;
        isLoaded.value = true;
        linkElement = link;
        resolve(link);
      };

      link.onerror = () => {
        isLoading.value = false;
        const err = new Error(`Failed to load stylesheet: ${href}`);
        error.value = err;
        reject(err);
      };

      document.head.appendChild(link);
    });
  };

  const unload = () => {
    if (linkElement && linkElement.parentNode) {
      linkElement.parentNode.removeChild(linkElement);
      linkElement = null;
      isLoaded.value = false;
    }
  };

  if (immediate) {
    onMounted(() => {
      load().catch(() => {});
    });
  }

  if (removeOnUnmount) {
    onUnmounted(unload);
  }

  return {
    isLoading,
    isLoaded,
    error,
    load,
    unload,
  };
}
