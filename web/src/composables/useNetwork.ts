/**
 * 网络状态检测 Composable
 * 检测在线/离线状态和网络连接信息
 */

import { ref, computed, onMounted, onUnmounted } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface NetworkState {
  /** 是否在线 */
  isOnline: boolean;
  /** 自上次状态变化以来的时间 */
  since?: Date;
  /** 下行速度 (Mbps) */
  downlink?: number;
  /** 下行速度上限 (Mbps) */
  downlinkMax?: number;
  /** 有效类型: slow-2g, 2g, 3g, 4g */
  effectiveType?: "slow-2g" | "2g" | "3g" | "4g";
  /** 往返时间 (ms) */
  rtt?: number;
  /** 是否启用数据节省模式 */
  saveData?: boolean;
  /** 连接类型 */
  type?:
    | "bluetooth"
    | "cellular"
    | "ethernet"
    | "none"
    | "wifi"
    | "wimax"
    | "other"
    | "unknown";
}

export interface UseNetworkOptions {
  /** 是否立即激活，默认 true */
  immediate?: boolean;
}

// ============================================================================
// 扩展 Navigator 类型
// ============================================================================

interface NetworkInformation extends EventTarget {
  downlink?: number;
  downlinkMax?: number;
  effectiveType?: "slow-2g" | "2g" | "3g" | "4g";
  rtt?: number;
  saveData?: boolean;
  type?:
    | "bluetooth"
    | "cellular"
    | "ethernet"
    | "none"
    | "wifi"
    | "wimax"
    | "other"
    | "unknown";
  onchange?: EventListener;
}

declare global {
  interface Navigator {
    connection?: NetworkInformation;
    mozConnection?: NetworkInformation;
    webkitConnection?: NetworkInformation;
  }
}

// ============================================================================
// 主函数
// ============================================================================

/**
 * 网络状态检测
 * @example
 * const { isOnline, effectiveType, downlink } = useNetwork()
 *
 * // 监听离线状态
 * watch(isOnline, (online) => {
 *   if (!online) showOfflineNotification()
 * })
 */
export function useNetwork(options: UseNetworkOptions = {}) {
  const { immediate = true } = options;

  // 状态
  const isOnline = ref(navigator.onLine);
  const since = ref<Date | undefined>();
  const downlink = ref<number | undefined>();
  const downlinkMax = ref<number | undefined>();
  const effectiveType = ref<NetworkState["effectiveType"]>();
  const rtt = ref<number | undefined>();
  const saveData = ref<boolean | undefined>();
  const type = ref<NetworkState["type"]>();

  // 获取网络连接对象
  const getConnection = (): NetworkInformation | undefined => {
    return (
      navigator.connection ||
      navigator.mozConnection ||
      navigator.webkitConnection
    );
  };

  // 更新网络信息
  const updateNetworkInfo = () => {
    const connection = getConnection();

    if (connection) {
      downlink.value = connection.downlink;
      downlinkMax.value = connection.downlinkMax;
      effectiveType.value = connection.effectiveType;
      rtt.value = connection.rtt;
      saveData.value = connection.saveData;
      type.value = connection.type;
    }
  };

  // 在线事件处理
  const handleOnline = () => {
    isOnline.value = true;
    since.value = new Date();
  };

  // 离线事件处理
  const handleOffline = () => {
    isOnline.value = false;
    since.value = new Date();
  };

  // 网络变化处理
  const handleChange = () => {
    updateNetworkInfo();
  };

  // 计算属性：是否为慢速网络
  const isSlowConnection = computed(() => {
    return effectiveType.value === "slow-2g" || effectiveType.value === "2g";
  });

  // 计算属性：是否为快速网络
  const isFastConnection = computed(() => {
    return effectiveType.value === "4g";
  });

  // 计算属性：网络状态描述
  const connectionStatus = computed<"good" | "moderate" | "poor" | "offline">(() => {
    if (!isOnline.value) return "offline";
    if (effectiveType.value === "4g") return "good";
    if (effectiveType.value === "3g") return "moderate";
    return "poor";
  });

  // 启动监听
  const activate = () => {
    window.addEventListener("online", handleOnline);
    window.addEventListener("offline", handleOffline);

    const connection = getConnection();
    if (connection) {
      connection.addEventListener("change", handleChange);
    }

    // 初始化
    updateNetworkInfo();
  };

  // 停止监听
  const deactivate = () => {
    window.removeEventListener("online", handleOnline);
    window.removeEventListener("offline", handleOffline);

    const connection = getConnection();
    if (connection) {
      connection.removeEventListener("change", handleChange);
    }
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
    // 基础状态
    isOnline,
    since,

    // Network Information API
    downlink,
    downlinkMax,
    effectiveType,
    rtt,
    saveData,
    type,

    // 计算属性
    isSlowConnection,
    isFastConnection,
    connectionStatus,

    // 方法
    activate,
    deactivate,
  };
}

// ============================================================================
// 离线检测 Hook
// ============================================================================

export interface UseOnlineOptions {
  /** 离线回调 */
  onOffline?: () => void;
  /** 在线回调 */
  onOnline?: () => void;
}

/**
 * 简化的在线/离线检测
 * @example
 * const isOnline = useOnline({
 *   onOffline: () => toast.warning('网络已断开'),
 *   onOnline: () => toast.success('网络已恢复')
 * })
 */
export function useOnline(options: UseOnlineOptions = {}) {
  const { onOffline, onOnline } = options;

  const isOnline = ref(navigator.onLine);

  const handleOnline = () => {
    isOnline.value = true;
    onOnline?.();
  };

  const handleOffline = () => {
    isOnline.value = false;
    onOffline?.();
  };

  onMounted(() => {
    window.addEventListener("online", handleOnline);
    window.addEventListener("offline", handleOffline);
  });

  onUnmounted(() => {
    window.removeEventListener("online", handleOnline);
    window.removeEventListener("offline", handleOffline);
  });

  return isOnline;
}

// ============================================================================
// 网络速度检测
// ============================================================================

export interface UseNetworkSpeedOptions {
  /** 测试文件 URL（应该是一个小文件） */
  testUrl?: string;
  /** 测试文件大小（字节），用于计算速度 */
  testFileSize?: number;
  /** 超时时间（毫秒） */
  timeout?: number;
}

/**
 * 网络速度检测
 * @example
 * const { speed, isLoading, test } = useNetworkSpeed()
 * await test()
 * console.log(`下载速度: ${speed.value} Mbps`)
 */
export function useNetworkSpeed(options: UseNetworkSpeedOptions = {}) {
  const {
    testUrl = "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
    testFileSize = 13504, // Google logo 大约大小
    timeout = 10000,
  } = options;

  const speed = ref<number | null>(null);
  const isLoading = ref(false);
  const error = ref<Error | null>(null);

  const test = async (): Promise<number | null> => {
    isLoading.value = true;
    error.value = null;

    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), timeout);

      const startTime = performance.now();

      // 添加时间戳避免缓存
      const url = `${testUrl}?_=${Date.now()}`;

      const response = await fetch(url, {
        method: "GET",
        cache: "no-store",
        signal: controller.signal,
      });

      if (!response.ok) {
        throw new Error(`HTTP error: ${response.status}`);
      }

      // 读取完整响应
      await response.blob();

      clearTimeout(timeoutId);

      const endTime = performance.now();
      const duration = (endTime - startTime) / 1000; // 秒

      // 计算速度 (Mbps)
      const bitsLoaded = testFileSize * 8;
      const speedMbps = bitsLoaded / duration / 1000000;

      speed.value = Math.round(speedMbps * 100) / 100;
      return speed.value;
    } catch (e) {
      const err = e instanceof Error ? e : new Error(String(e));
      error.value = err;
      return null;
    } finally {
      isLoading.value = false;
    }
  };

  return {
    speed,
    isLoading,
    error,
    test,
  };
}

// ============================================================================
// 网络状态提示组件数据
// ============================================================================

/**
 * 网络状态提示数据
 * @example
 * const { shouldShowBanner, bannerType, bannerMessage } = useNetworkBanner()
 */
export function useNetworkBanner() {
  const { isOnline, isSlowConnection, connectionStatus } = useNetwork();

  const shouldShowBanner = computed(() => {
    return !isOnline.value || isSlowConnection.value;
  });

  const bannerType = computed<"error" | "warning" | "info">(() => {
    if (!isOnline.value) return "error";
    if (isSlowConnection.value) return "warning";
    return "info";
  });

  const bannerMessage = computed(() => {
    if (!isOnline.value) {
      return "您当前处于离线状态，部分功能可能不可用";
    }
    if (isSlowConnection.value) {
      return "网络连接较慢，加载可能需要更长时间";
    }
    return "";
  });

  return {
    isOnline,
    shouldShowBanner,
    bannerType,
    bannerMessage,
    connectionStatus,
  };
}
