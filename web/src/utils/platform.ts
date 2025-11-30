/**
 * 平台检测工具
 * 检测浏览器、操作系统、设备类型等
 */

// ============================================================================
// 类型定义
// ============================================================================

export interface BrowserInfo {
  name: string;
  version: string;
  major: number;
}

export interface OSInfo {
  name: string;
  version: string;
}

export interface DeviceInfo {
  type: "mobile" | "tablet" | "desktop";
  vendor: string;
  model: string;
}

export interface PlatformInfo {
  browser: BrowserInfo;
  os: OSInfo;
  device: DeviceInfo;
  isTouch: boolean;
  isStandalone: boolean;
  language: string;
  languages: string[];
  cookieEnabled: boolean;
  doNotTrack: boolean;
}

// ============================================================================
// 浏览器检测
// ============================================================================

/**
 * 检测浏览器信息
 */
export function detectBrowser(): BrowserInfo {
  const ua = navigator.userAgent;

  // Chrome
  if (/Chrome/.test(ua) && !/Chromium|Edge|OPR|Edg/.test(ua)) {
    const match = ua.match(/Chrome\/(\d+)\.(\d+)/);
    return {
      name: "Chrome",
      version: match ? `${match[1]}.${match[2]}` : "",
      major: match ? parseInt(match[1], 10) : 0,
    };
  }

  // Edge (Chromium)
  if (/Edg/.test(ua)) {
    const match = ua.match(/Edg\/(\d+)\.(\d+)/);
    return {
      name: "Edge",
      version: match ? `${match[1]}.${match[2]}` : "",
      major: match ? parseInt(match[1], 10) : 0,
    };
  }

  // Firefox
  if (/Firefox/.test(ua)) {
    const match = ua.match(/Firefox\/(\d+)\.(\d+)/);
    return {
      name: "Firefox",
      version: match ? `${match[1]}.${match[2]}` : "",
      major: match ? parseInt(match[1], 10) : 0,
    };
  }

  // Safari
  if (/Safari/.test(ua) && !/Chrome|Chromium/.test(ua)) {
    const match = ua.match(/Version\/(\d+)\.(\d+)/);
    return {
      name: "Safari",
      version: match ? `${match[1]}.${match[2]}` : "",
      major: match ? parseInt(match[1], 10) : 0,
    };
  }

  // Opera
  if (/OPR/.test(ua)) {
    const match = ua.match(/OPR\/(\d+)\.(\d+)/);
    return {
      name: "Opera",
      version: match ? `${match[1]}.${match[2]}` : "",
      major: match ? parseInt(match[1], 10) : 0,
    };
  }

  // IE
  if (/MSIE|Trident/.test(ua)) {
    const match = ua.match(/(?:MSIE |rv:)(\d+)\.(\d+)/);
    return {
      name: "IE",
      version: match ? `${match[1]}.${match[2]}` : "",
      major: match ? parseInt(match[1], 10) : 0,
    };
  }

  return { name: "Unknown", version: "", major: 0 };
}

// ============================================================================
// 操作系统检测
// ============================================================================

/**
 * 检测操作系统信息
 */
export function detectOS(): OSInfo {
  const ua = navigator.userAgent;
  const platform = navigator.platform;

  // Windows
  if (/Windows/.test(ua)) {
    const match = ua.match(/Windows NT (\d+\.\d+)/);
    const versionMap: Record<string, string> = {
      "10.0": "10/11",
      "6.3": "8.1",
      "6.2": "8",
      "6.1": "7",
      "6.0": "Vista",
      "5.1": "XP",
    };
    return {
      name: "Windows",
      version: match ? versionMap[match[1]] || match[1] : "",
    };
  }

  // macOS
  if (/Mac OS X/.test(ua)) {
    const match = ua.match(/Mac OS X (\d+[._]\d+)/);
    return {
      name: "macOS",
      version: match ? match[1].replace("_", ".") : "",
    };
  }

  // iOS
  if (/iPhone|iPad|iPod/.test(ua)) {
    const match = ua.match(/OS (\d+_\d+)/);
    return {
      name: "iOS",
      version: match ? match[1].replace("_", ".") : "",
    };
  }

  // Android
  if (/Android/.test(ua)) {
    const match = ua.match(/Android (\d+\.?\d*)/);
    return {
      name: "Android",
      version: match ? match[1] : "",
    };
  }

  // Linux
  if (/Linux/.test(platform)) {
    return { name: "Linux", version: "" };
  }

  return { name: "Unknown", version: "" };
}

// ============================================================================
// 设备检测
// ============================================================================

/**
 * 检测设备信息
 */
export function detectDevice(): DeviceInfo {
  const ua = navigator.userAgent;

  // 移动设备
  const isMobile = /iPhone|iPod|Android.*Mobile|webOS|BlackBerry|IEMobile|Opera Mini/i.test(ua);

  // 平板
  const isTablet = /iPad|Android(?!.*Mobile)|Tablet/i.test(ua);

  // 设备类型
  let type: DeviceInfo["type"] = "desktop";
  if (isMobile) type = "mobile";
  else if (isTablet) type = "tablet";

  // 设备厂商和型号
  let vendor = "";
  let model = "";

  if (/iPhone/.test(ua)) {
    vendor = "Apple";
    model = "iPhone";
  } else if (/iPad/.test(ua)) {
    vendor = "Apple";
    model = "iPad";
  } else if (/Samsung/.test(ua)) {
    vendor = "Samsung";
    const match = ua.match(/SM-\w+/);
    model = match ? match[0] : "";
  } else if (/Huawei|HUAWEI/.test(ua)) {
    vendor = "Huawei";
  } else if (/Xiaomi|Mi /.test(ua)) {
    vendor = "Xiaomi";
  }

  return { type, vendor, model };
}

// ============================================================================
// 特性检测
// ============================================================================

/**
 * 检测是否支持触摸
 */
export function isTouchDevice(): boolean {
  return (
    "ontouchstart" in window ||
    navigator.maxTouchPoints > 0 ||
    (window.matchMedia && window.matchMedia("(pointer: coarse)").matches)
  );
}

/**
 * 检测是否在 PWA 独立模式
 */
export function isStandalone(): boolean {
  return (
    window.matchMedia("(display-mode: standalone)").matches ||
    (navigator as Navigator & { standalone?: boolean }).standalone === true
  );
}

/**
 * 检测是否支持 WebGL
 */
export function supportsWebGL(): boolean {
  try {
    const canvas = document.createElement("canvas");
    return !!(window.WebGLRenderingContext && (canvas.getContext("webgl") || canvas.getContext("experimental-webgl")));
  } catch {
    return false;
  }
}

/**
 * 检测是否支持 WebP
 */
export async function supportsWebP(): Promise<boolean> {
  return new Promise((resolve) => {
    const img = new Image();
    img.onload = () => resolve(img.width === 1);
    img.onerror = () => resolve(false);
    img.src = "data:image/webp;base64,UklGRiQAAABXRUJQVlA4IBgAAAAwAQCdASoBAAEAAwA0JaQAA3AA/vuUAAA=";
  });
}

/**
 * 检测是否支持 AVIF
 */
export async function supportsAVIF(): Promise<boolean> {
  return new Promise((resolve) => {
    const img = new Image();
    img.onload = () => resolve(img.width === 1);
    img.onerror = () => resolve(false);
    img.src =
      "data:image/avif;base64,AAAAIGZ0eXBhdmlmAAAAAGF2aWZtaWYxbWlhZk1BMUIAAADybWV0YQAAAAAAAAAoaGRscgAAAAAAAAAAcGljdAAAAAAAAAAAAAAAAGxpYmF2aWYAAAAADnBpdG0AAAAAAAEAAAAeaWxvYwAAAABEAAABAAEAAAABAAABGgAAAB0AAAAoaWluZgAAAAAAAQAAABppbmZlAgAAAAABAABhdjAxQ29sb3IAAAAAamlwcnAAAABLaXBjbwAAABRpc3BlAAAAAAAAAAIAAAACAAAAEHBpeGkAAAAAAwgICAAAAAxhdjFDgQ0MAAAAABNjb2xybmNseAACAAIAAYAAAAAXaXBtYQAAAAAAAAABAAEEAQKDBAAAACVtZGF0EgAKBzgABpAQ0AIUDxgADAAxFQAAP/H///sBAAA=";
  });
}

/**
 * 检测是否为暗色模式偏好
 */
export function prefersDarkMode(): boolean {
  return window.matchMedia("(prefers-color-scheme: dark)").matches;
}

/**
 * 检测是否偏好减少动画
 */
export function prefersReducedMotion(): boolean {
  return window.matchMedia("(prefers-reduced-motion: reduce)").matches;
}

// ============================================================================
// 综合平台信息
// ============================================================================

/**
 * 获取完整平台信息
 */
export function getPlatformInfo(): PlatformInfo {
  return {
    browser: detectBrowser(),
    os: detectOS(),
    device: detectDevice(),
    isTouch: isTouchDevice(),
    isStandalone: isStandalone(),
    language: navigator.language,
    languages: [...navigator.languages],
    cookieEnabled: navigator.cookieEnabled,
    doNotTrack: navigator.doNotTrack === "1",
  };
}

// ============================================================================
// 便捷判断函数
// ============================================================================

/** 是否为 Chrome 浏览器 */
export const isChrome = (): boolean => detectBrowser().name === "Chrome";

/** 是否为 Firefox 浏览器 */
export const isFirefox = (): boolean => detectBrowser().name === "Firefox";

/** 是否为 Safari 浏览器 */
export const isSafari = (): boolean => detectBrowser().name === "Safari";

/** 是否为 Edge 浏览器 */
export const isEdge = (): boolean => detectBrowser().name === "Edge";

/** 是否为 IE 浏览器 */
export const isIE = (): boolean => detectBrowser().name === "IE";

/** 是否为 Windows 系统 */
export const isWindows = (): boolean => detectOS().name === "Windows";

/** 是否为 macOS 系统 */
export const isMacOS = (): boolean => detectOS().name === "macOS";

/** 是否为 iOS 系统 */
export const isIOS = (): boolean => detectOS().name === "iOS";

/** 是否为 Android 系统 */
export const isAndroid = (): boolean => detectOS().name === "Android";

/** 是否为 Linux 系统 */
export const isLinux = (): boolean => detectOS().name === "Linux";

/** 是否为移动设备 */
export const isMobile = (): boolean => detectDevice().type === "mobile";

/** 是否为平板设备 */
export const isTablet = (): boolean => detectDevice().type === "tablet";

/** 是否为桌面设备 */
export const isDesktop = (): boolean => detectDevice().type === "desktop";

/** 是否为苹果设备 */
export const isAppleDevice = (): boolean => {
  const os = detectOS().name;
  return os === "macOS" || os === "iOS";
};
