/**
 * Preferences Composable
 * 提供用户偏好设置的响应式管理
 */

import { ref, watch, computed, type Ref, type ComputedRef } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UsePreferredDarkReturn {
  /** 是否偏好深色模式 */
  isDark: Ref<boolean>;
  /** 是否支持 */
  isSupported: ComputedRef<boolean>;
}

export interface UsePreferredLanguageReturn {
  /** 首选语言 */
  language: Ref<string>;
  /** 语言列表 */
  languages: Ref<readonly string[]>;
  /** 是否支持 */
  isSupported: ComputedRef<boolean>;
}

export interface UsePreferredReducedMotionReturn {
  /** 是否偏好减少动画 */
  isReduced: Ref<boolean>;
  /** 是否支持 */
  isSupported: ComputedRef<boolean>;
}

export interface UsePreferredContrastReturn {
  /** 对比度偏好 */
  contrast: Ref<"more" | "less" | "custom" | "no-preference">;
  /** 是否支持 */
  isSupported: ComputedRef<boolean>;
}

export interface UsePreferredColorSchemeReturn {
  /** 颜色方案偏好 */
  colorScheme: Ref<"light" | "dark" | "no-preference">;
  /** 是否支持 */
  isSupported: ComputedRef<boolean>;
}

// ============================================================================
// usePreferredDark - 深色模式偏好
// ============================================================================

/**
 * 检测用户是否偏好深色模式
 * @example
 * const { isDark } = usePreferredDark()
 * watch(isDark, (dark) => {
 *   document.body.classList.toggle('dark', dark)
 * })
 */
export function usePreferredDark(): UsePreferredDarkReturn {
  const isSupported = computed(() => typeof window !== "undefined" && "matchMedia" in window);

  const isDark = ref(false);
  let mediaQuery: MediaQueryList | null = null;

  const update = () => {
    if (mediaQuery) {
      isDark.value = mediaQuery.matches;
    }
  };

  if (isSupported.value) {
    mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
    isDark.value = mediaQuery.matches;

    if ("addEventListener" in mediaQuery) {
      mediaQuery.addEventListener("change", update);
    } else {
      // 旧版浏览器
      (mediaQuery as MediaQueryList).addListener(update);
    }
  }

  return {
    isDark,
    isSupported,
  };
}

// ============================================================================
// usePreferredLanguage - 语言偏好
// ============================================================================

/**
 * 获取用户的语言偏好
 * @example
 * const { language, languages } = usePreferredLanguage()
 * console.log(language.value) // 'zh-CN'
 */
export function usePreferredLanguage(): UsePreferredLanguageReturn {
  const isSupported = computed(() => typeof navigator !== "undefined" && "languages" in navigator);

  const language = ref(typeof navigator !== "undefined" ? navigator.language : "en");

  const languages = ref<readonly string[]>(typeof navigator !== "undefined" ? navigator.languages : ["en"]);

  if (isSupported.value) {
    const handleLanguageChange = () => {
      language.value = navigator.language;
      languages.value = navigator.languages;
    };

    window.addEventListener("languagechange", handleLanguageChange);
  }

  return {
    language,
    languages,
    isSupported,
  };
}

// ============================================================================
// usePreferredReducedMotion - 减少动画偏好
// ============================================================================

/**
 * 检测用户是否偏好减少动画
 * @example
 * const { isReduced } = usePreferredReducedMotion()
 * const transition = computed(() =>
 *   isReduced.value ? 'none' : 'all 0.3s ease'
 * )
 */
export function usePreferredReducedMotion(): UsePreferredReducedMotionReturn {
  const isSupported = computed(() => typeof window !== "undefined" && "matchMedia" in window);

  const isReduced = ref(false);
  let mediaQuery: MediaQueryList | null = null;

  const update = () => {
    if (mediaQuery) {
      isReduced.value = mediaQuery.matches;
    }
  };

  if (isSupported.value) {
    mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
    isReduced.value = mediaQuery.matches;

    if ("addEventListener" in mediaQuery) {
      mediaQuery.addEventListener("change", update);
    } else {
      (mediaQuery as MediaQueryList).addListener(update);
    }
  }

  return {
    isReduced,
    isSupported,
  };
}

// ============================================================================
// usePreferredContrast - 对比度偏好
// ============================================================================

/**
 * 检测用户的对比度偏好
 * @example
 * const { contrast } = usePreferredContrast()
 * // contrast.value: 'more' | 'less' | 'custom' | 'no-preference'
 */
export function usePreferredContrast(): UsePreferredContrastReturn {
  const isSupported = computed(() => typeof window !== "undefined" && "matchMedia" in window);

  const contrast = ref<"more" | "less" | "custom" | "no-preference">("no-preference");

  const update = () => {
    if (window.matchMedia("(prefers-contrast: more)").matches) {
      contrast.value = "more";
    } else if (window.matchMedia("(prefers-contrast: less)").matches) {
      contrast.value = "less";
    } else if (window.matchMedia("(prefers-contrast: custom)").matches) {
      contrast.value = "custom";
    } else {
      contrast.value = "no-preference";
    }
  };

  if (isSupported.value) {
    update();

    // 监听变化
    const queries = [
      window.matchMedia("(prefers-contrast: more)"),
      window.matchMedia("(prefers-contrast: less)"),
      window.matchMedia("(prefers-contrast: custom)"),
    ];

    queries.forEach((query) => {
      if ("addEventListener" in query) {
        query.addEventListener("change", update);
      } else {
        (query as MediaQueryList).addListener(update);
      }
    });
  }

  return {
    contrast,
    isSupported,
  };
}

// ============================================================================
// usePreferredColorScheme - 颜色方案偏好
// ============================================================================

/**
 * 检测用户的颜色方案偏好
 * @example
 * const { colorScheme } = usePreferredColorScheme()
 * // colorScheme.value: 'light' | 'dark' | 'no-preference'
 */
export function usePreferredColorScheme(): UsePreferredColorSchemeReturn {
  const isSupported = computed(() => typeof window !== "undefined" && "matchMedia" in window);

  const colorScheme = ref<"light" | "dark" | "no-preference">("no-preference");

  const update = () => {
    if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
      colorScheme.value = "dark";
    } else if (window.matchMedia("(prefers-color-scheme: light)").matches) {
      colorScheme.value = "light";
    } else {
      colorScheme.value = "no-preference";
    }
  };

  if (isSupported.value) {
    update();

    const darkQuery = window.matchMedia("(prefers-color-scheme: dark)");
    const lightQuery = window.matchMedia("(prefers-color-scheme: light)");

    if ("addEventListener" in darkQuery) {
      darkQuery.addEventListener("change", update);
      lightQuery.addEventListener("change", update);
    } else {
      (darkQuery as MediaQueryList).addListener(update);
      (lightQuery as MediaQueryList).addListener(update);
    }
  }

  return {
    colorScheme,
    isSupported,
  };
}

// ============================================================================
// useDark - 深色模式控制
// ============================================================================

export interface UseDarkOptions {
  /** 默认值 */
  initialValue?: boolean | "auto";
  /** 存储键名 */
  storageKey?: string;
  /** 属性选择器 */
  selector?: string;
  /** 属性名 */
  attribute?: string;
  /** 深色类名 */
  valueDark?: string;
  /** 浅色类名 */
  valueLight?: string;
  /** 变化回调 */
  onChanged?: (isDark: boolean) => void;
}

export interface UseDarkReturn {
  /** 是否为深色模式 */
  isDark: Ref<boolean>;
  /** 切换模式 */
  toggle: () => void;
  /** 设置为深色 */
  setDark: () => void;
  /** 设置为浅色 */
  setLight: () => void;
  /** 系统偏好 */
  systemIsDark: Ref<boolean>;
}

/**
 * 深色模式控制
 * @example
 * const { isDark, toggle } = useDark()
 * // 自动同步到 document.documentElement
 */
export function useDark(options: UseDarkOptions = {}): UseDarkReturn {
  const {
    initialValue = "auto",
    storageKey = "vue-use-dark",
    selector = "html",
    attribute = "class",
    valueDark = "dark",
    valueLight = "",
    onChanged,
  } = options;

  const { isDark: systemIsDark } = usePreferredDark();

  // 从存储读取
  const getStoredValue = (): boolean | null => {
    if (typeof localStorage === "undefined") return null;
    const stored = localStorage.getItem(storageKey);
    if (stored === "true") return true;
    if (stored === "false") return false;
    return null;
  };

  // 确定初始值
  const getInitialValue = (): boolean => {
    const stored = getStoredValue();
    if (stored !== null) return stored;
    if (initialValue === "auto") return systemIsDark.value;
    return initialValue;
  };

  const isDark = ref(getInitialValue());

  // 应用到 DOM
  const applyDark = (dark: boolean) => {
    if (typeof document === "undefined") return;

    const el = document.querySelector(selector);
    if (!el) return;

    if (attribute === "class") {
      if (dark) {
        if (valueDark) el.classList.add(valueDark);
        if (valueLight) el.classList.remove(valueLight);
      } else {
        if (valueDark) el.classList.remove(valueDark);
        if (valueLight) el.classList.add(valueLight);
      }
    } else {
      el.setAttribute(attribute, dark ? valueDark : valueLight);
    }

    onChanged?.(dark);
  };

  // 保存到存储
  const saveValue = (dark: boolean) => {
    if (typeof localStorage !== "undefined") {
      localStorage.setItem(storageKey, String(dark));
    }
  };

  // 监听变化
  watch(
    isDark,
    (dark) => {
      saveValue(dark);
      applyDark(dark);
    },
    { immediate: true }
  );

  // 监听系统变化（仅当使用 auto 时）
  if (initialValue === "auto") {
    watch(systemIsDark, (dark) => {
      if (getStoredValue() === null) {
        isDark.value = dark;
      }
    });
  }

  const toggle = () => {
    isDark.value = !isDark.value;
  };

  const setDark = () => {
    isDark.value = true;
  };

  const setLight = () => {
    isDark.value = false;
  };

  return {
    isDark,
    toggle,
    setDark,
    setLight,
    systemIsDark,
  };
}

// ============================================================================
// useColorMode - 颜色模式管理
// ============================================================================

export type ColorMode = "light" | "dark" | "auto" | string;

export interface UseColorModeOptions {
  /** 默认模式 */
  initialValue?: ColorMode;
  /** 存储键名 */
  storageKey?: string;
  /** 可用模式 */
  modes?: ColorMode[];
  /** 属性选择器 */
  selector?: string;
  /** 属性名 */
  attribute?: string;
}

export interface UseColorModeReturn {
  /** 当前模式 */
  mode: Ref<ColorMode>;
  /** 实际应用的模式 */
  resolvedMode: ComputedRef<"light" | "dark">;
  /** 设置模式 */
  setMode: (mode: ColorMode) => void;
  /** 循环切换模式 */
  cycle: () => void;
}

/**
 * 颜色模式管理
 * @example
 * const { mode, resolvedMode, cycle } = useColorMode({
 *   modes: ['light', 'dark', 'auto']
 * })
 */
export function useColorMode(options: UseColorModeOptions = {}): UseColorModeReturn {
  const {
    initialValue = "auto",
    storageKey = "vue-use-color-mode",
    modes = ["light", "dark", "auto"],
    selector = "html",
    attribute = "data-theme",
  } = options;

  const { isDark: systemIsDark } = usePreferredDark();

  // 从存储读取
  const getStoredValue = (): ColorMode | null => {
    if (typeof localStorage === "undefined") return null;
    const stored = localStorage.getItem(storageKey);
    return stored && modes.includes(stored) ? stored : null;
  };

  const mode = ref<ColorMode>(getStoredValue() || initialValue);

  // 解析实际模式
  const resolvedMode = computed<"light" | "dark">(() => {
    if (mode.value === "auto") {
      return systemIsDark.value ? "dark" : "light";
    }
    return mode.value === "dark" ? "dark" : "light";
  });

  // 应用到 DOM
  const applyMode = (modeValue: string) => {
    if (typeof document === "undefined") return;

    const el = document.querySelector(selector);
    if (!el) return;

    el.setAttribute(attribute, modeValue);
  };

  // 保存到存储
  const saveValue = (modeValue: ColorMode) => {
    if (typeof localStorage !== "undefined") {
      localStorage.setItem(storageKey, modeValue);
    }
  };

  // 监听变化
  watch(
    mode,
    (value) => {
      saveValue(value);
      applyMode(value === "auto" ? resolvedMode.value : value);
    },
    { immediate: true }
  );

  // 监听系统变化（当模式为 auto 时）
  watch(
    () => mode.value === "auto" && systemIsDark.value,
    () => {
      if (mode.value === "auto") {
        applyMode(resolvedMode.value);
      }
    }
  );

  const setMode = (newMode: ColorMode) => {
    mode.value = newMode;
  };

  const cycle = () => {
    const currentIndex = modes.indexOf(mode.value);
    const nextIndex = (currentIndex + 1) % modes.length;
    mode.value = modes[nextIndex];
  };

  return {
    mode,
    resolvedMode,
    setMode,
    cycle,
  };
}

// ============================================================================
// useTransparency - 透明度偏好
// ============================================================================

export interface UsePreferredTransparencyReturn {
  /** 是否偏好透明度减少 */
  isReduced: Ref<boolean>;
  /** 是否支持 */
  isSupported: ComputedRef<boolean>;
}

/**
 * 检测用户是否偏好减少透明度
 * @example
 * const { isReduced } = usePreferredTransparency()
 */
export function usePreferredTransparency(): UsePreferredTransparencyReturn {
  const isSupported = computed(() => typeof window !== "undefined" && "matchMedia" in window);

  const isReduced = ref(false);

  if (isSupported.value) {
    const mediaQuery = window.matchMedia("(prefers-reduced-transparency: reduce)");
    isReduced.value = mediaQuery.matches;

    if ("addEventListener" in mediaQuery) {
      mediaQuery.addEventListener("change", (e) => {
        isReduced.value = e.matches;
      });
    }
  }

  return {
    isReduced,
    isSupported,
  };
}
