/**
 * Permission Composable
 * 提供浏览器权限 API 的封装
 */

import { ref, onUnmounted, type Ref, computed } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

/** 权限名称类型 */
export type PermissionName =
  | "geolocation"
  | "notifications"
  | "push"
  | "midi"
  | "camera"
  | "microphone"
  | "speaker-selection"
  | "device-info"
  | "background-fetch"
  | "background-sync"
  | "bluetooth"
  | "persistent-storage"
  | "ambient-light-sensor"
  | "accelerometer"
  | "gyroscope"
  | "magnetometer"
  | "clipboard-read"
  | "clipboard-write"
  | "display-capture"
  | "screen-wake-lock";

/** 权限状态 */
export type PermissionState = "granted" | "denied" | "prompt";

export interface UsePermissionOptions {
  /** 权限名称 */
  name: PermissionName;
  /** 是否立即查询 */
  immediate?: boolean;
}

export interface UsePermissionReturn {
  /** 权限状态 */
  state: Ref<PermissionState | undefined>;
  /** 是否已授权 */
  isGranted: Ref<boolean>;
  /** 是否被拒绝 */
  isDenied: Ref<boolean>;
  /** 是否等待用户选择 */
  isPrompt: Ref<boolean>;
  /** 是否支持此权限 */
  isSupported: Ref<boolean>;
  /** 刷新权限状态 */
  query: () => Promise<PermissionState | undefined>;
}

// ============================================================================
// usePermission - 权限查询
// ============================================================================

/**
 * 查询浏览器权限状态
 * @example
 * const { state, isGranted, query } = usePermission({ name: 'notifications' })
 * if (isGranted.value) {
 *   sendNotification()
 * }
 */
export function usePermission(options: UsePermissionOptions): UsePermissionReturn {
  const { name, immediate = true } = options;

  const state = ref<PermissionState | undefined>(undefined);
  const isSupported = ref(false);

  let permissionStatus: PermissionStatus | null = null;

  const isGranted = computed(() => state.value === "granted");
  const isDenied = computed(() => state.value === "denied");
  const isPrompt = computed(() => state.value === "prompt");

  const handleChange = () => {
    if (permissionStatus) {
      state.value = permissionStatus.state;
    }
  };

  const query = async (): Promise<PermissionState | undefined> => {
    if (typeof navigator === "undefined" || !navigator.permissions) {
      isSupported.value = false;
      return undefined;
    }

    try {
      isSupported.value = true;
      permissionStatus = await navigator.permissions.query({
        name: name as PermissionName,
      });

      state.value = permissionStatus.state;
      permissionStatus.addEventListener("change", handleChange);

      return state.value;
    } catch {
      isSupported.value = false;
      return undefined;
    }
  };

  if (immediate) {
    query();
  }

  onUnmounted(() => {
    if (permissionStatus) {
      permissionStatus.removeEventListener("change", handleChange);
    }
  });

  return {
    state,
    isGranted,
    isDenied,
    isPrompt,
    isSupported,
    query,
  };
}

// ============================================================================
// useNotificationPermission - 通知权限
// ============================================================================

export interface UseNotificationPermissionReturn {
  /** 权限状态 */
  permission: Ref<NotificationPermission>;
  /** 是否已授权 */
  isGranted: Ref<boolean>;
  /** 是否被拒绝 */
  isDenied: Ref<boolean>;
  /** 是否支持通知 */
  isSupported: Ref<boolean>;
  /** 请求权限 */
  request: () => Promise<NotificationPermission>;
  /** 发送通知 */
  notify: (title: string, options?: NotificationOptions) => Notification | null;
}

/**
 * 通知权限管理
 * @example
 * const { isGranted, request, notify } = useNotificationPermission()
 *
 * async function sendNotification() {
 *   if (!isGranted.value) {
 *     await request()
 *   }
 *   notify('Hello!', { body: 'This is a notification' })
 * }
 */
export function useNotificationPermission(): UseNotificationPermissionReturn {
  const isSupported = ref(typeof window !== "undefined" && "Notification" in window);
  const permission = ref<NotificationPermission>(
    isSupported.value ? Notification.permission : "denied"
  );

  const isGranted = computed(() => permission.value === "granted");
  const isDenied = computed(() => permission.value === "denied");

  const request = async (): Promise<NotificationPermission> => {
    if (!isSupported.value) {
      return "denied";
    }

    try {
      permission.value = await Notification.requestPermission();
      return permission.value;
    } catch {
      return "denied";
    }
  };

  const notify = (
    title: string,
    options?: NotificationOptions
  ): Notification | null => {
    if (!isSupported.value || !isGranted.value) {
      return null;
    }

    return new Notification(title, options);
  };

  return {
    permission,
    isGranted,
    isDenied,
    isSupported,
    request,
    notify,
  };
}

// ============================================================================
// useClipboardPermission - 剪贴板权限
// ============================================================================

export interface UseClipboardPermissionReturn {
  /** 读取权限状态 */
  readState: Ref<PermissionState | undefined>;
  /** 写入权限状态 */
  writeState: Ref<PermissionState | undefined>;
  /** 是否可以读取 */
  canRead: Ref<boolean>;
  /** 是否可以写入 */
  canWrite: Ref<boolean>;
  /** 查询读取权限 */
  queryRead: () => Promise<PermissionState | undefined>;
  /** 查询写入权限 */
  queryWrite: () => Promise<PermissionState | undefined>;
}

/**
 * 剪贴板权限管理
 * @example
 * const { canRead, canWrite } = useClipboardPermission()
 */
export function useClipboardPermission(): UseClipboardPermissionReturn {
  const readPermission = usePermission({ name: "clipboard-read" });
  const writePermission = usePermission({ name: "clipboard-write" });

  return {
    readState: readPermission.state,
    writeState: writePermission.state,
    canRead: readPermission.isGranted,
    canWrite: writePermission.isGranted,
    queryRead: readPermission.query,
    queryWrite: writePermission.query,
  };
}

// ============================================================================
// useCameraPermission - 摄像头权限
// ============================================================================

export interface UseCameraPermissionReturn {
  /** 权限状态 */
  state: Ref<PermissionState | undefined>;
  /** 是否已授权 */
  isGranted: Ref<boolean>;
  /** 是否支持 */
  isSupported: Ref<boolean>;
  /** 请求权限 */
  request: () => Promise<MediaStream | null>;
  /** 停止所有轨道 */
  stop: () => void;
  /** 当前媒体流 */
  stream: Ref<MediaStream | null>;
}

/**
 * 摄像头权限管理
 * @example
 * const { isGranted, request, stream } = useCameraPermission()
 *
 * async function startCamera() {
 *   const mediaStream = await request()
 *   if (mediaStream) {
 *     videoElement.srcObject = mediaStream
 *   }
 * }
 */
export function useCameraPermission(): UseCameraPermissionReturn {
  const permission = usePermission({ name: "camera" });
  const stream = ref<MediaStream | null>(null);

  const request = async (): Promise<MediaStream | null> => {
    if (typeof navigator === "undefined" || !navigator.mediaDevices) {
      return null;
    }

    try {
      stream.value = await navigator.mediaDevices.getUserMedia({ video: true });
      await permission.query();
      return stream.value;
    } catch {
      await permission.query();
      return null;
    }
  };

  const stop = () => {
    if (stream.value) {
      stream.value.getTracks().forEach((track) => track.stop());
      stream.value = null;
    }
  };

  onUnmounted(stop);

  return {
    state: permission.state,
    isGranted: permission.isGranted,
    isSupported: permission.isSupported,
    request,
    stop,
    stream,
  };
}

// ============================================================================
// useMicrophonePermission - 麦克风权限
// ============================================================================

export interface UseMicrophonePermissionReturn {
  /** 权限状态 */
  state: Ref<PermissionState | undefined>;
  /** 是否已授权 */
  isGranted: Ref<boolean>;
  /** 是否支持 */
  isSupported: Ref<boolean>;
  /** 请求权限 */
  request: () => Promise<MediaStream | null>;
  /** 停止所有轨道 */
  stop: () => void;
  /** 当前媒体流 */
  stream: Ref<MediaStream | null>;
}

/**
 * 麦克风权限管理
 * @example
 * const { isGranted, request, stream } = useMicrophonePermission()
 */
export function useMicrophonePermission(): UseMicrophonePermissionReturn {
  const permission = usePermission({ name: "microphone" });
  const stream = ref<MediaStream | null>(null);

  const request = async (): Promise<MediaStream | null> => {
    if (typeof navigator === "undefined" || !navigator.mediaDevices) {
      return null;
    }

    try {
      stream.value = await navigator.mediaDevices.getUserMedia({ audio: true });
      await permission.query();
      return stream.value;
    } catch {
      await permission.query();
      return null;
    }
  };

  const stop = () => {
    if (stream.value) {
      stream.value.getTracks().forEach((track) => track.stop());
      stream.value = null;
    }
  };

  onUnmounted(stop);

  return {
    state: permission.state,
    isGranted: permission.isGranted,
    isSupported: permission.isSupported,
    request,
    stop,
    stream,
  };
}

// ============================================================================
// useGeolocationPermission - 地理位置权限
// ============================================================================

export interface UseGeolocationPermissionReturn {
  /** 权限状态 */
  state: Ref<PermissionState | undefined>;
  /** 是否已授权 */
  isGranted: Ref<boolean>;
  /** 是否被拒绝 */
  isDenied: Ref<boolean>;
  /** 是否支持 */
  isSupported: Ref<boolean>;
  /** 刷新状态 */
  query: () => Promise<PermissionState | undefined>;
}

/**
 * 地理位置权限管理
 * @example
 * const { isGranted, isDenied } = useGeolocationPermission()
 */
export function useGeolocationPermission(): UseGeolocationPermissionReturn {
  return usePermission({ name: "geolocation" });
}

// ============================================================================
// useScreenWakeLock - 屏幕唤醒锁
// ============================================================================

export interface UseScreenWakeLockReturn {
  /** 是否已锁定 */
  isActive: Ref<boolean>;
  /** 是否支持 */
  isSupported: Ref<boolean>;
  /** 请求锁定 */
  request: () => Promise<boolean>;
  /** 释放锁定 */
  release: () => Promise<void>;
}

/**
 * 屏幕唤醒锁
 * 阻止屏幕变暗或锁定
 * @example
 * const { isActive, request, release } = useScreenWakeLock()
 *
 * // 在视频播放时保持屏幕常亮
 * onMounted(() => request())
 * onUnmounted(() => release())
 */
export function useScreenWakeLock(): UseScreenWakeLockReturn {
  const isActive = ref(false);
  const isSupported = ref(typeof navigator !== "undefined" && "wakeLock" in navigator);

  let wakeLock: WakeLockSentinel | null = null;

  const handleVisibilityChange = async () => {
    if (document.visibilityState === "visible" && isActive.value) {
      await request();
    }
  };

  const request = async (): Promise<boolean> => {
    if (!isSupported.value) {
      return false;
    }

    try {
      wakeLock = await (navigator as Navigator).wakeLock.request("screen");
      isActive.value = true;

      wakeLock.addEventListener("release", () => {
        isActive.value = false;
      });

      document.addEventListener("visibilitychange", handleVisibilityChange);

      return true;
    } catch {
      isActive.value = false;
      return false;
    }
  };

  const release = async (): Promise<void> => {
    if (wakeLock) {
      await wakeLock.release();
      wakeLock = null;
    }
    isActive.value = false;
    document.removeEventListener("visibilitychange", handleVisibilityChange);
  };

  onUnmounted(() => {
    release();
  });

  return {
    isActive,
    isSupported,
    request,
    release,
  };
}

// ============================================================================
// 权限检查工具
// ============================================================================

/**
 * 检查是否支持 Permissions API
 */
export function isPermissionsApiSupported(): boolean {
  return typeof navigator !== "undefined" && "permissions" in navigator;
}

/**
 * 批量查询权限
 * @example
 * const results = await queryPermissions(['camera', 'microphone'])
 * // { camera: 'granted', microphone: 'prompt' }
 */
export async function queryPermissions(
  names: PermissionName[]
): Promise<Record<PermissionName, PermissionState | undefined>> {
  const results: Record<string, PermissionState | undefined> = {};

  if (!isPermissionsApiSupported()) {
    for (const name of names) {
      results[name] = undefined;
    }
    return results;
  }

  await Promise.all(
    names.map(async (name) => {
      try {
        const status = await navigator.permissions.query({
          name: name as PermissionName,
        });
        results[name] = status.state;
      } catch {
        results[name] = undefined;
      }
    })
  );

  return results;
}
