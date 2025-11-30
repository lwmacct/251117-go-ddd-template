/**
 * 存储工具 Composable
 * 提供 localStorage/sessionStorage 的响应式封装
 */
import { ref, watch, type Ref } from "vue";

export interface UseStorageOptions<T> {
  /** 存储类型 */
  storage?: "local" | "session";
  /** 默认值 */
  defaultValue?: T;
  /** 是否监听变化并自动保存 */
  watchChange?: boolean;
  /** 过期时间（毫秒），仅 localStorage 有效 */
  expires?: number;
  /** 序列化函数 */
  serializer?: {
    read: (raw: string) => T;
    write: (value: T) => string;
  };
}

interface StorageItem<T> {
  value: T;
  expires?: number;
}

/**
 * 创建响应式存储
 */
export function useStorage<T>(
  key: string,
  options: UseStorageOptions<T> = {}
): {
  value: Ref<T | null>;
  set: (newValue: T, customExpires?: number) => void;
  get: () => T | null;
  remove: () => void;
  isExpired: () => boolean;
} {
  const {
    storage = "local",
    defaultValue = null,
    watchChange = true,
    expires,
    serializer = {
      read: (raw: string) => JSON.parse(raw) as T,
      write: (value: T) => JSON.stringify(value),
    },
  } = options;

  const storageApi = storage === "local" ? localStorage : sessionStorage;

  // 读取初始值
  const readValue = (): T | null => {
    try {
      const raw = storageApi.getItem(key);
      if (raw === null) {
        return defaultValue as T | null;
      }

      const item = JSON.parse(raw) as StorageItem<T>;

      // 检查过期
      if (item.expires && Date.now() > item.expires) {
        storageApi.removeItem(key);
        return defaultValue as T | null;
      }

      return item.value;
    } catch {
      return defaultValue as T | null;
    }
  };

  const value = ref<T | null>(readValue()) as Ref<T | null>;

  /**
   * 设置值
   */
  const set = (newValue: T, customExpires?: number) => {
    try {
      const expiresAt = customExpires ?? expires;
      const item: StorageItem<T> = {
        value: newValue,
        expires: expiresAt ? Date.now() + expiresAt : undefined,
      };
      storageApi.setItem(key, JSON.stringify(item));
      value.value = newValue;
    } catch (error) {
      console.error(`Failed to save to ${storage}Storage:`, error);
    }
  };

  /**
   * 获取值
   */
  const get = (): T | null => {
    return readValue();
  };

  /**
   * 删除值
   */
  const remove = () => {
    storageApi.removeItem(key);
    value.value = defaultValue as T | null;
  };

  /**
   * 检查是否过期
   */
  const isExpired = (): boolean => {
    try {
      const raw = storageApi.getItem(key);
      if (raw === null) return true;

      const item = JSON.parse(raw) as StorageItem<T>;
      if (item.expires && Date.now() > item.expires) {
        return true;
      }
      return false;
    } catch {
      return true;
    }
  };

  // 监听变化并自动保存
  if (watchChange) {
    watch(
      value,
      (newValue) => {
        if (newValue !== null) {
          set(newValue);
        }
      },
      { deep: true }
    );
  }

  return {
    value,
    set,
    get,
    remove,
    isExpired,
  };
}

/**
 * 快捷方法：localStorage
 */
export function useLocalStorage<T>(
  key: string,
  options: Omit<UseStorageOptions<T>, "storage"> = {}
) {
  return useStorage<T>(key, { ...options, storage: "local" });
}

/**
 * 快捷方法：sessionStorage
 */
export function useSessionStorage<T>(
  key: string,
  options: Omit<UseStorageOptions<T>, "storage"> = {}
) {
  return useStorage<T>(key, { ...options, storage: "session" });
}

/**
 * 清除所有带特定前缀的存储项
 */
export function clearStorageByPrefix(prefix: string, storage: "local" | "session" = "local") {
  const storageApi = storage === "local" ? localStorage : sessionStorage;
  const keysToRemove: string[] = [];

  for (let i = 0; i < storageApi.length; i++) {
    const key = storageApi.key(i);
    if (key && key.startsWith(prefix)) {
      keysToRemove.push(key);
    }
  }

  keysToRemove.forEach((key) => storageApi.removeItem(key));

  return keysToRemove.length;
}
