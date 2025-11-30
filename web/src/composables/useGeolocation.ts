/**
 * Geolocation Composable
 * 提供地理位置相关的响应式状态和功能
 */

import { ref, onUnmounted, type Ref, computed } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface GeolocationCoordinates {
  /** 纬度 */
  latitude: number;
  /** 经度 */
  longitude: number;
  /** 海拔高度（米） */
  altitude: number | null;
  /** 精度（米） */
  accuracy: number;
  /** 海拔精度（米） */
  altitudeAccuracy: number | null;
  /** 方向（度，0-360） */
  heading: number | null;
  /** 速度（米/秒） */
  speed: number | null;
}

export interface GeolocationPosition {
  coords: GeolocationCoordinates;
  timestamp: number;
}

export interface UseGeolocationOptions {
  /** 是否立即获取位置 */
  immediate?: boolean;
  /** 是否启用高精度模式 */
  enableHighAccuracy?: boolean;
  /** 位置缓存最大时间（毫秒） */
  maximumAge?: number;
  /** 超时时间（毫秒） */
  timeout?: number;
}

export interface UseGeolocationReturn {
  /** 当前坐标 */
  coords: Ref<GeolocationCoordinates | null>;
  /** 位置时间戳 */
  locatedAt: Ref<number | null>;
  /** 是否正在定位 */
  isLoading: Ref<boolean>;
  /** 错误信息 */
  error: Ref<GeolocationPositionError | null>;
  /** 是否支持地理位置 */
  isSupported: Ref<boolean>;
  /** 获取当前位置 */
  getCurrentPosition: () => Promise<GeolocationPosition | null>;
  /** 暂停位置监听 */
  pause: () => void;
  /** 恢复位置监听 */
  resume: () => void;
}

// ============================================================================
// useGeolocation - 地理位置
// ============================================================================

/**
 * 获取和跟踪地理位置
 * @example
 * const { coords, isLoading, error, getCurrentPosition } = useGeolocation()
 *
 * // 获取当前位置
 * const position = await getCurrentPosition()
 *
 * // 使用响应式坐标
 * console.log(coords.value?.latitude, coords.value?.longitude)
 */
export function useGeolocation(options: UseGeolocationOptions = {}): UseGeolocationReturn {
  const { immediate = true, enableHighAccuracy = true, maximumAge = 30000, timeout = 27000 } = options;

  const isSupported = ref(typeof navigator !== "undefined" && "geolocation" in navigator);
  const coords = ref<GeolocationCoordinates | null>(null);
  const locatedAt = ref<number | null>(null);
  const error = ref<GeolocationPositionError | null>(null);
  const isLoading = ref(false);

  let watchId: number | null = null;

  const positionOptions: PositionOptions = {
    enableHighAccuracy,
    maximumAge,
    timeout,
  };

  const updatePosition = (position: GeolocationPosition) => {
    coords.value = {
      latitude: position.coords.latitude,
      longitude: position.coords.longitude,
      altitude: position.coords.altitude,
      accuracy: position.coords.accuracy,
      altitudeAccuracy: position.coords.altitudeAccuracy,
      heading: position.coords.heading,
      speed: position.coords.speed,
    };
    locatedAt.value = position.timestamp;
    error.value = null;
    isLoading.value = false;
  };

  const handleError = (err: GeolocationPositionError) => {
    error.value = err;
    isLoading.value = false;
  };

  const getCurrentPosition = (): Promise<GeolocationPosition | null> => {
    if (!isSupported.value) {
      return Promise.resolve(null);
    }

    isLoading.value = true;

    return new Promise((resolve) => {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          updatePosition(position);
          resolve({
            coords: coords.value!,
            timestamp: locatedAt.value!,
          });
        },
        (err) => {
          handleError(err);
          resolve(null);
        },
        positionOptions
      );
    });
  };

  const pause = () => {
    if (watchId !== null) {
      navigator.geolocation.clearWatch(watchId);
      watchId = null;
    }
  };

  const resume = () => {
    if (!isSupported.value || watchId !== null) {
      return;
    }

    watchId = navigator.geolocation.watchPosition(updatePosition, handleError, positionOptions);
  };

  if (immediate && isSupported.value) {
    resume();
  }

  onUnmounted(pause);

  return {
    coords,
    locatedAt,
    isLoading,
    error,
    isSupported,
    getCurrentPosition,
    pause,
    resume,
  };
}

// ============================================================================
// useGeolocationDistance - 距离计算
// ============================================================================

export interface DistanceUnit {
  /** 单位名称 */
  name: string;
  /** 转换因子（相对于米） */
  factor: number;
}

export const DISTANCE_UNITS: Record<string, DistanceUnit> = {
  m: { name: "meters", factor: 1 },
  km: { name: "kilometers", factor: 0.001 },
  mi: { name: "miles", factor: 0.000621371 },
  ft: { name: "feet", factor: 3.28084 },
  yd: { name: "yards", factor: 1.09361 },
  nm: { name: "nautical miles", factor: 0.000539957 },
};

/**
 * 计算两点之间的距离（Haversine 公式）
 * @example
 * const distance = calculateDistance(
 *   { latitude: 40.7128, longitude: -74.0060 },  // 纽约
 *   { latitude: 34.0522, longitude: -118.2437 }, // 洛杉矶
 *   'km'
 * )
 * // 约 3936 公里
 */
export function calculateDistance(
  from: { latitude: number; longitude: number },
  to: { latitude: number; longitude: number },
  unit: keyof typeof DISTANCE_UNITS = "m"
): number {
  const R = 6371000; // 地球半径（米）

  const toRad = (deg: number) => (deg * Math.PI) / 180;

  const dLat = toRad(to.latitude - from.latitude);
  const dLon = toRad(to.longitude - from.longitude);

  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(toRad(from.latitude)) * Math.cos(toRad(to.latitude)) * Math.sin(dLon / 2) * Math.sin(dLon / 2);

  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  const distanceInMeters = R * c;

  return distanceInMeters * DISTANCE_UNITS[unit].factor;
}

/**
 * 格式化距离
 * @example
 * formatDistance(1500) // '1.5 km'
 * formatDistance(500) // '500 m'
 */
export function formatDistance(meters: number): string {
  if (meters >= 1000) {
    return `${(meters / 1000).toFixed(1)} km`;
  }
  return `${Math.round(meters)} m`;
}

// ============================================================================
// useGeolocationWatch - 位置监听
// ============================================================================

export interface UseGeolocationWatchOptions extends UseGeolocationOptions {
  /** 位置变化回调 */
  onPositionChange?: (coords: GeolocationCoordinates) => void;
  /** 错误回调 */
  onError?: (error: GeolocationPositionError) => void;
  /** 最小距离变化（米）才触发回调 */
  minDistance?: number;
}

/**
 * 持续监听位置变化
 * @example
 * useGeolocationWatch({
 *   onPositionChange: (coords) => {
 *     updateUserLocation(coords.latitude, coords.longitude)
 *   },
 *   minDistance: 10, // 移动超过 10 米才更新
 * })
 */
export function useGeolocationWatch(options: UseGeolocationWatchOptions = {}): UseGeolocationReturn {
  const { onPositionChange, onError, minDistance = 0, ...restOptions } = options;

  const geolocation = useGeolocation({
    ...restOptions,
    immediate: false,
  });

  const lastCoords = ref<GeolocationCoordinates | null>(null);

  // 监听位置变化
  const checkPositionChange = () => {
    const currentCoords = geolocation.coords.value;
    if (!currentCoords) return;

    if (lastCoords.value && minDistance > 0) {
      const distance = calculateDistance(
        {
          latitude: lastCoords.value.latitude,
          longitude: lastCoords.value.longitude,
        },
        {
          latitude: currentCoords.latitude,
          longitude: currentCoords.longitude,
        }
      );

      if (distance < minDistance) {
        return;
      }
    }

    lastCoords.value = { ...currentCoords };
    onPositionChange?.(currentCoords);
  };

  // 监听错误
  const checkError = () => {
    if (geolocation.error.value) {
      onError?.(geolocation.error.value);
    }
  };

  // 自动开始监听
  geolocation.resume();

  // 设置定期检查
  const checkInterval = setInterval(() => {
    checkPositionChange();
    checkError();
  }, 1000);

  onUnmounted(() => {
    clearInterval(checkInterval);
  });

  return geolocation;
}

// ============================================================================
// useGeolocationBounds - 边界检测
// ============================================================================

export interface GeolocationBounds {
  north: number;
  south: number;
  east: number;
  west: number;
}

export interface UseGeolocationBoundsOptions {
  /** 边界范围 */
  bounds: GeolocationBounds;
  /** 进入边界回调 */
  onEnter?: () => void;
  /** 离开边界回调 */
  onLeave?: () => void;
}

/**
 * 检测是否在指定边界内
 * @example
 * const { isInBounds } = useGeolocationBounds({
 *   bounds: {
 *     north: 40.8,
 *     south: 40.6,
 *     east: -73.9,
 *     west: -74.1,
 *   },
 *   onEnter: () => console.log('进入区域'),
 *   onLeave: () => console.log('离开区域'),
 * })
 */
export function useGeolocationBounds(options: UseGeolocationBoundsOptions): {
  isInBounds: Ref<boolean>;
  geolocation: UseGeolocationReturn;
} {
  const { bounds, onEnter, onLeave } = options;

  const geolocation = useGeolocation({ immediate: true });
  const isInBounds = ref(false);
  const wasInBounds = ref(false);

  const checkBounds = () => {
    const coords = geolocation.coords.value;
    if (!coords) return;

    const inBounds =
      coords.latitude <= bounds.north &&
      coords.latitude >= bounds.south &&
      coords.longitude <= bounds.east &&
      coords.longitude >= bounds.west;

    if (inBounds !== isInBounds.value) {
      isInBounds.value = inBounds;

      if (inBounds && !wasInBounds.value) {
        onEnter?.();
      } else if (!inBounds && wasInBounds.value) {
        onLeave?.();
      }

      wasInBounds.value = inBounds;
    }
  };

  const checkInterval = setInterval(checkBounds, 1000);

  onUnmounted(() => {
    clearInterval(checkInterval);
  });

  return {
    isInBounds,
    geolocation,
  };
}

// ============================================================================
// 坐标转换工具
// ============================================================================

/**
 * 度分秒转十进制度
 * @example
 * dmsToDecimal(40, 42, 46, 'N') // 40.7128
 */
export function dmsToDecimal(
  degrees: number,
  minutes: number,
  seconds: number,
  direction: "N" | "S" | "E" | "W"
): number {
  let decimal = degrees + minutes / 60 + seconds / 3600;
  if (direction === "S" || direction === "W") {
    decimal = -decimal;
  }
  return decimal;
}

/**
 * 十进制度转度分秒
 * @example
 * decimalToDms(40.7128, 'lat')
 * // { degrees: 40, minutes: 42, seconds: 46, direction: 'N' }
 */
export function decimalToDms(
  decimal: number,
  type: "lat" | "lng"
): {
  degrees: number;
  minutes: number;
  seconds: number;
  direction: "N" | "S" | "E" | "W";
} {
  const direction = type === "lat" ? (decimal >= 0 ? "N" : "S") : decimal >= 0 ? "E" : "W";

  const absolute = Math.abs(decimal);
  const degrees = Math.floor(absolute);
  const minutesFloat = (absolute - degrees) * 60;
  const minutes = Math.floor(minutesFloat);
  const seconds = (minutesFloat - minutes) * 60;

  return {
    degrees,
    minutes,
    seconds: Math.round(seconds * 100) / 100,
    direction,
  };
}

/**
 * 格式化坐标为字符串
 * @example
 * formatCoordinates(40.7128, -74.0060)
 * // '40.7128°N, 74.0060°W'
 */
export function formatCoordinates(latitude: number, longitude: number, precision: number = 4): string {
  const latDir = latitude >= 0 ? "N" : "S";
  const lngDir = longitude >= 0 ? "E" : "W";

  return `${Math.abs(latitude).toFixed(precision)}°${latDir}, ${Math.abs(longitude).toFixed(precision)}°${lngDir}`;
}

/**
 * 生成 Google Maps 链接
 * @example
 * getGoogleMapsUrl(40.7128, -74.0060)
 * // 'https://www.google.com/maps?q=40.7128,-74.006'
 */
export function getGoogleMapsUrl(latitude: number, longitude: number, zoom: number = 15): string {
  return `https://www.google.com/maps?q=${latitude},${longitude}&z=${zoom}`;
}

/**
 * 生成 Apple Maps 链接
 * @example
 * getAppleMapsUrl(40.7128, -74.0060)
 */
export function getAppleMapsUrl(latitude: number, longitude: number, zoom: number = 15): string {
  return `https://maps.apple.com/?ll=${latitude},${longitude}&z=${zoom}`;
}

/**
 * 生成导航链接（跨平台）
 * @example
 * getNavigationUrl(40.7128, -74.0060)
 */
export function getNavigationUrl(latitude: number, longitude: number): string {
  // 检测平台
  const isIOS = typeof navigator !== "undefined" && /iPad|iPhone|iPod/.test(navigator.userAgent);

  if (isIOS) {
    return getAppleMapsUrl(latitude, longitude);
  }

  return getGoogleMapsUrl(latitude, longitude);
}
