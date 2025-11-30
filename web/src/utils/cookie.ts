/**
 * Cookie 工具函数
 * 提供 Cookie 的读取、写入和管理
 */

// ============================================================================
// 类型定义
// ============================================================================

export interface CookieOptions {
  /** 过期时间（秒）或 Date 对象 */
  expires?: number | Date;
  /** 有效路径 */
  path?: string;
  /** 有效域名 */
  domain?: string;
  /** 是否仅 HTTPS */
  secure?: boolean;
  /** SameSite 策略 */
  sameSite?: "Strict" | "Lax" | "None";
  /** 最大存活时间（秒） */
  maxAge?: number;
}

export interface ParsedCookie {
  name: string;
  value: string;
  options?: Partial<CookieOptions>;
}

// ============================================================================
// 基础操作
// ============================================================================

/**
 * 获取 Cookie
 * @example
 * getCookie('token') // 'abc123'
 * getCookie('nonexistent') // null
 */
export function getCookie(name: string): string | null {
  if (typeof document === "undefined") {
    return null;
  }

  const cookies = document.cookie.split(";");

  for (const cookie of cookies) {
    const [cookieName, cookieValue] = cookie.trim().split("=");

    if (cookieName === name) {
      try {
        return decodeURIComponent(cookieValue);
      } catch {
        return cookieValue;
      }
    }
  }

  return null;
}

/**
 * 设置 Cookie
 * @example
 * setCookie('token', 'abc123')
 * setCookie('token', 'abc123', { expires: 3600, secure: true })
 * setCookie('token', 'abc123', { expires: new Date('2024-12-31') })
 */
export function setCookie(name: string, value: string, options: CookieOptions = {}): void {
  if (typeof document === "undefined") {
    return;
  }

  const { expires, path = "/", domain, secure, sameSite = "Lax", maxAge } = options;

  let cookieString = `${encodeURIComponent(name)}=${encodeURIComponent(value)}`;

  if (expires !== undefined) {
    let expiresDate: Date;

    if (expires instanceof Date) {
      expiresDate = expires;
    } else {
      expiresDate = new Date();
      expiresDate.setTime(expiresDate.getTime() + expires * 1000);
    }

    cookieString += `; expires=${expiresDate.toUTCString()}`;
  }

  if (maxAge !== undefined) {
    cookieString += `; max-age=${maxAge}`;
  }

  if (path) {
    cookieString += `; path=${path}`;
  }

  if (domain) {
    cookieString += `; domain=${domain}`;
  }

  if (secure) {
    cookieString += "; secure";
  }

  if (sameSite) {
    cookieString += `; samesite=${sameSite}`;
  }

  document.cookie = cookieString;
}

/**
 * 删除 Cookie
 * @example
 * removeCookie('token')
 * removeCookie('token', { path: '/', domain: '.example.com' })
 */
export function removeCookie(name: string, options: Pick<CookieOptions, "path" | "domain"> = {}): void {
  setCookie(name, "", {
    ...options,
    expires: new Date(0),
  });
}

/**
 * 检查 Cookie 是否存在
 * @example
 * hasCookie('token') // true/false
 */
export function hasCookie(name: string): boolean {
  return getCookie(name) !== null;
}

// ============================================================================
// 批量操作
// ============================================================================

/**
 * 获取所有 Cookie
 * @example
 * getAllCookies() // { token: 'abc', user: 'john' }
 */
export function getAllCookies(): Record<string, string> {
  if (typeof document === "undefined") {
    return {};
  }

  const cookies: Record<string, string> = {};
  const cookieString = document.cookie;

  if (!cookieString) {
    return cookies;
  }

  for (const cookie of cookieString.split(";")) {
    const [name, value] = cookie.trim().split("=");

    if (name) {
      try {
        cookies[name] = decodeURIComponent(value || "");
      } catch {
        cookies[name] = value || "";
      }
    }
  }

  return cookies;
}

/**
 * 获取多个 Cookie
 * @example
 * getCookies(['token', 'user']) // { token: 'abc', user: 'john' }
 */
export function getCookies(names: string[]): Record<string, string | null> {
  const result: Record<string, string | null> = {};

  for (const name of names) {
    result[name] = getCookie(name);
  }

  return result;
}

/**
 * 设置多个 Cookie
 * @example
 * setCookies({ token: 'abc', user: 'john' }, { expires: 3600 })
 */
export function setCookies(cookies: Record<string, string>, options: CookieOptions = {}): void {
  for (const [name, value] of Object.entries(cookies)) {
    setCookie(name, value, options);
  }
}

/**
 * 删除多个 Cookie
 * @example
 * removeCookies(['token', 'user'])
 */
export function removeCookies(names: string[], options: Pick<CookieOptions, "path" | "domain"> = {}): void {
  for (const name of names) {
    removeCookie(name, options);
  }
}

/**
 * 清除所有 Cookie
 * @example
 * clearAllCookies()
 */
export function clearAllCookies(options: Pick<CookieOptions, "path" | "domain"> = {}): void {
  const cookies = getAllCookies();

  for (const name of Object.keys(cookies)) {
    removeCookie(name, options);
  }
}

// ============================================================================
// JSON Cookie
// ============================================================================

/**
 * 获取 JSON Cookie
 * @example
 * getJsonCookie<User>('user') // { id: 1, name: 'John' }
 */
export function getJsonCookie<T>(name: string): T | null {
  const value = getCookie(name);

  if (value === null) {
    return null;
  }

  try {
    return JSON.parse(value) as T;
  } catch {
    return null;
  }
}

/**
 * 设置 JSON Cookie
 * @example
 * setJsonCookie('user', { id: 1, name: 'John' })
 */
export function setJsonCookie<T>(name: string, value: T, options: CookieOptions = {}): void {
  setCookie(name, JSON.stringify(value), options);
}

// ============================================================================
// Cookie 解析
// ============================================================================

/**
 * 解析 Cookie 字符串
 * @example
 * parseCookieString('token=abc; user=john')
 * // { token: 'abc', user: 'john' }
 */
export function parseCookieString(cookieString: string): Record<string, string> {
  const cookies: Record<string, string> = {};

  if (!cookieString) {
    return cookies;
  }

  for (const cookie of cookieString.split(";")) {
    const [name, value] = cookie.trim().split("=");

    if (name) {
      try {
        cookies[name] = decodeURIComponent(value || "");
      } catch {
        cookies[name] = value || "";
      }
    }
  }

  return cookies;
}

/**
 * 序列化 Cookie
 * @example
 * serializeCookie('token', 'abc', { expires: 3600, secure: true })
 * // 'token=abc; expires=...; path=/; secure; samesite=Lax'
 */
export function serializeCookie(name: string, value: string, options: CookieOptions = {}): string {
  const { expires, path = "/", domain, secure, sameSite = "Lax", maxAge } = options;

  let cookieString = `${encodeURIComponent(name)}=${encodeURIComponent(value)}`;

  if (expires !== undefined) {
    let expiresDate: Date;

    if (expires instanceof Date) {
      expiresDate = expires;
    } else {
      expiresDate = new Date();
      expiresDate.setTime(expiresDate.getTime() + expires * 1000);
    }

    cookieString += `; expires=${expiresDate.toUTCString()}`;
  }

  if (maxAge !== undefined) {
    cookieString += `; max-age=${maxAge}`;
  }

  if (path) {
    cookieString += `; path=${path}`;
  }

  if (domain) {
    cookieString += `; domain=${domain}`;
  }

  if (secure) {
    cookieString += "; secure";
  }

  if (sameSite) {
    cookieString += `; samesite=${sameSite}`;
  }

  return cookieString;
}

// ============================================================================
// Cookie 工具
// ============================================================================

/**
 * 获取 Cookie 数量
 * @example
 * getCookieCount() // 5
 */
export function getCookieCount(): number {
  return Object.keys(getAllCookies()).length;
}

/**
 * 获取所有 Cookie 名称
 * @example
 * getCookieNames() // ['token', 'user', 'session']
 */
export function getCookieNames(): string[] {
  return Object.keys(getAllCookies());
}

/**
 * 检查 Cookie 是否启用
 * @example
 * areCookiesEnabled() // true/false
 */
export function areCookiesEnabled(): boolean {
  if (typeof document === "undefined") {
    return false;
  }

  try {
    // 尝试设置测试 cookie
    const testKey = "__cookie_test__";
    document.cookie = `${testKey}=1`;
    const enabled = document.cookie.indexOf(testKey) !== -1;

    // 清理测试 cookie
    document.cookie = `${testKey}=; expires=Thu, 01 Jan 1970 00:00:00 GMT`;

    return enabled;
  } catch {
    return false;
  }
}

/**
 * 估算 Cookie 总大小（字节）
 * @example
 * getCookiesSize() // 1024
 */
export function getCookiesSize(): number {
  if (typeof document === "undefined") {
    return 0;
  }

  return new Blob([document.cookie]).size;
}

/**
 * 获取 Cookie 剩余可用空间（估算）
 * 浏览器通常限制每个域名约 4KB
 * @example
 * getCookiesRemainingSpace() // 3072
 */
export function getCookiesRemainingSpace(): number {
  const MAX_COOKIE_SIZE = 4096; // 4KB
  return Math.max(0, MAX_COOKIE_SIZE - getCookiesSize());
}

// ============================================================================
// Cookie 管理器
// ============================================================================

export interface CookieManager {
  get(name: string): string | null;
  getJson<T>(name: string): T | null;
  set(name: string, value: string, options?: CookieOptions): void;
  setJson<T>(name: string, value: T, options?: CookieOptions): void;
  remove(name: string): void;
  has(name: string): boolean;
  clear(): void;
  getAll(): Record<string, string>;
}

/**
 * 创建 Cookie 管理器
 * @example
 * const cookies = createCookieManager({ path: '/', secure: true })
 * cookies.set('token', 'abc')
 * cookies.get('token') // 'abc'
 */
export function createCookieManager(defaultOptions: CookieOptions = {}): CookieManager {
  return {
    get(name: string) {
      return getCookie(name);
    },

    getJson<T>(name: string) {
      return getJsonCookie<T>(name);
    },

    set(name: string, value: string, options: CookieOptions = {}) {
      setCookie(name, value, { ...defaultOptions, ...options });
    },

    setJson<T>(name: string, value: T, options: CookieOptions = {}) {
      setJsonCookie(name, value, { ...defaultOptions, ...options });
    },

    remove(name: string) {
      removeCookie(name, {
        path: defaultOptions.path,
        domain: defaultOptions.domain,
      });
    },

    has(name: string) {
      return hasCookie(name);
    },

    clear() {
      clearAllCookies({
        path: defaultOptions.path,
        domain: defaultOptions.domain,
      });
    },

    getAll() {
      return getAllCookies();
    },
  };
}

// ============================================================================
// 预设 Cookie 配置
// ============================================================================

/** 会话 Cookie 配置（浏览器关闭后过期） */
export const SESSION_COOKIE: CookieOptions = {
  path: "/",
  sameSite: "Lax",
};

/** 持久 Cookie 配置（7 天） */
export const PERSISTENT_COOKIE: CookieOptions = {
  path: "/",
  maxAge: 7 * 24 * 60 * 60, // 7 days
  sameSite: "Lax",
};

/** 安全 Cookie 配置 */
export const SECURE_COOKIE: CookieOptions = {
  path: "/",
  secure: true,
  sameSite: "Strict",
};

/** 跨站 Cookie 配置 */
export const CROSS_SITE_COOKIE: CookieOptions = {
  path: "/",
  secure: true,
  sameSite: "None",
};
