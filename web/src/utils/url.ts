/**
 * URL 工具函数
 * 提供 URL 解析、构建和操作
 */

// ============================================================================
// 类型定义
// ============================================================================

export interface ParsedURL {
  href: string;
  protocol: string;
  host: string;
  hostname: string;
  port: string;
  pathname: string;
  search: string;
  hash: string;
  origin: string;
  username: string;
  password: string;
}

export interface QueryParams {
  [key: string]: string | string[] | undefined;
}

// ============================================================================
// URL 解析
// ============================================================================

/**
 * 解析 URL
 * @example
 * parseURL('https://example.com:8080/path?a=1#hash')
 * // { protocol: 'https:', host: 'example.com:8080', ... }
 */
export function parseURL(url: string, base?: string): ParsedURL {
  const parsed = new URL(url, base);

  return {
    href: parsed.href,
    protocol: parsed.protocol,
    host: parsed.host,
    hostname: parsed.hostname,
    port: parsed.port,
    pathname: parsed.pathname,
    search: parsed.search,
    hash: parsed.hash,
    origin: parsed.origin,
    username: parsed.username,
    password: parsed.password,
  };
}

/**
 * 检查是否为有效 URL
 * @example
 * isValidURL('https://example.com') // true
 * isValidURL('not-a-url') // false
 */
export function isValidURL(url: string): boolean {
  try {
    new URL(url);
    return true;
  } catch {
    return false;
  }
}

/**
 * 检查是否为绝对 URL
 * @example
 * isAbsoluteURL('https://example.com') // true
 * isAbsoluteURL('/path/to/page') // false
 */
export function isAbsoluteURL(url: string): boolean {
  return /^[a-z][a-z0-9+.-]*:\/\//i.test(url);
}

/**
 * 检查是否为相对 URL
 * @example
 * isRelativeURL('/path/to/page') // true
 * isRelativeURL('https://example.com') // false
 */
export function isRelativeURL(url: string): boolean {
  return !isAbsoluteURL(url);
}

// ============================================================================
// 查询参数
// ============================================================================

/**
 * 解析查询字符串
 * @example
 * parseQuery('a=1&b=2&c=3') // { a: '1', b: '2', c: '3' }
 * parseQuery('?a=1&b=2') // { a: '1', b: '2' }
 * parseQuery('a=1&a=2') // { a: ['1', '2'] }
 */
export function parseQuery(query: string): QueryParams {
  const search = query.startsWith("?") ? query.slice(1) : query;

  if (!search) {
    return {};
  }

  const params: QueryParams = {};
  const searchParams = new URLSearchParams(search);

  for (const [key, value] of searchParams.entries()) {
    const existing = params[key];

    if (existing === undefined) {
      params[key] = value;
    } else if (Array.isArray(existing)) {
      existing.push(value);
    } else {
      params[key] = [existing, value];
    }
  }

  return params;
}

/**
 * 构建查询字符串
 * @example
 * buildQuery({ a: '1', b: '2' }) // 'a=1&b=2'
 * buildQuery({ a: ['1', '2'] }) // 'a=1&a=2'
 */
export function buildQuery(params: QueryParams, options: { encode?: boolean; prefix?: boolean } = {}): string {
  const { encode = true, prefix = false } = options;

  const searchParams = new URLSearchParams();

  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === null) {
      continue;
    }

    if (Array.isArray(value)) {
      for (const v of value) {
        searchParams.append(key, encode ? encodeURIComponent(v) : v);
      }
    } else {
      searchParams.append(key, encode ? encodeURIComponent(value) : value);
    }
  }

  const query = searchParams.toString();
  return prefix && query ? `?${query}` : query;
}

/**
 * 获取 URL 查询参数
 * @example
 * getQueryParam('https://example.com?a=1&b=2', 'a') // '1'
 * getQueryParam('?a=1&a=2', 'a') // ['1', '2']
 */
export function getQueryParam(url: string, key: string): string | string[] | undefined {
  const search = url.includes("?") ? url.split("?")[1].split("#")[0] : url;
  const params = parseQuery(search);
  return params[key];
}

/**
 * 设置 URL 查询参数
 * @example
 * setQueryParam('https://example.com?a=1', 'b', '2')
 * // 'https://example.com?a=1&b=2'
 */
export function setQueryParam(url: string, key: string, value: string | string[]): string {
  const parsed = new URL(url, "http://localhost");
  const params = parseQuery(parsed.search);

  params[key] = value;
  parsed.search = buildQuery(params, { prefix: true });

  // 如果原 URL 是相对路径，只返回路径部分
  if (isRelativeURL(url)) {
    return parsed.pathname + parsed.search + parsed.hash;
  }

  return parsed.href;
}

/**
 * 删除 URL 查询参数
 * @example
 * removeQueryParam('https://example.com?a=1&b=2', 'a')
 * // 'https://example.com?b=2'
 */
export function removeQueryParam(url: string, key: string): string {
  const parsed = new URL(url, "http://localhost");
  const params = parseQuery(parsed.search);

  delete params[key];
  parsed.search = buildQuery(params, { prefix: true });

  if (isRelativeURL(url)) {
    return parsed.pathname + parsed.search + parsed.hash;
  }

  return parsed.href;
}

/**
 * 合并查询参数到 URL
 * @example
 * mergeQueryParams('https://example.com?a=1', { b: '2', c: '3' })
 * // 'https://example.com?a=1&b=2&c=3'
 */
export function mergeQueryParams(url: string, params: QueryParams): string {
  const parsed = new URL(url, "http://localhost");
  const existingParams = parseQuery(parsed.search);
  const mergedParams = { ...existingParams, ...params };

  parsed.search = buildQuery(mergedParams, { prefix: true });

  if (isRelativeURL(url)) {
    return parsed.pathname + parsed.search + parsed.hash;
  }

  return parsed.href;
}

// ============================================================================
// URL 操作
// ============================================================================

/**
 * 连接 URL 路径
 * @example
 * joinURL('https://example.com', 'path', 'to', 'page')
 * // 'https://example.com/path/to/page'
 * joinURL('/api', '/users/', '/1')
 * // '/api/users/1'
 */
export function joinURL(...parts: string[]): string {
  if (parts.length === 0) return "";

  let result = parts[0];

  for (let i = 1; i < parts.length; i++) {
    const part = parts[i];

    // 移除前导斜杠
    const cleanPart = part.replace(/^\/+/, "");

    // 确保结果以斜杠结尾
    if (!result.endsWith("/") && cleanPart) {
      result += "/";
    }

    result += cleanPart;
  }

  // 移除尾部斜杠（除非是根路径）
  if (result.length > 1 && result.endsWith("/")) {
    result = result.slice(0, -1);
  }

  return result;
}

/**
 * 规范化 URL
 * @example
 * normalizeURL('https://example.com//path///to//page')
 * // 'https://example.com/path/to/page'
 */
export function normalizeURL(url: string): string {
  // 保留协议部分的双斜杠
  const protocolMatch = url.match(/^([a-z][a-z0-9+.-]*:\/\/)/i);
  const protocol = protocolMatch ? protocolMatch[1] : "";
  const rest = protocolMatch ? url.slice(protocol.length) : url;

  // 规范化路径中的多余斜杠
  const normalized = rest.replace(/\/+/g, "/");

  return protocol + normalized;
}

/**
 * 获取 URL 的基础路径
 * @example
 * getBasePath('https://example.com/path/to/page')
 * // 'https://example.com/path/to'
 */
export function getBasePath(url: string): string {
  const parsed = new URL(url, "http://localhost");
  const pathParts = parsed.pathname.split("/").filter(Boolean);
  pathParts.pop();

  const basePath = "/" + pathParts.join("/");

  if (isAbsoluteURL(url)) {
    return parsed.origin + basePath;
  }

  return basePath || "/";
}

/**
 * 获取 URL 的文件名
 * @example
 * getFileName('https://example.com/path/to/file.pdf')
 * // 'file.pdf'
 */
export function getFileName(url: string): string {
  const parsed = new URL(url, "http://localhost");
  const pathParts = parsed.pathname.split("/");
  return pathParts[pathParts.length - 1] || "";
}

/**
 * 获取 URL 的扩展名
 * @example
 * getExtension('https://example.com/path/to/file.pdf')
 * // 'pdf'
 */
export function getExtension(url: string): string {
  const fileName = getFileName(url);
  const dotIndex = fileName.lastIndexOf(".");

  if (dotIndex === -1 || dotIndex === 0) {
    return "";
  }

  return fileName.slice(dotIndex + 1).toLowerCase();
}

/**
 * 添加或更新 URL 的 hash
 * @example
 * setHash('https://example.com/page', 'section')
 * // 'https://example.com/page#section'
 */
export function setHash(url: string, hash: string): string {
  const parsed = new URL(url, "http://localhost");
  parsed.hash = hash.startsWith("#") ? hash : `#${hash}`;

  if (isRelativeURL(url)) {
    return parsed.pathname + parsed.search + parsed.hash;
  }

  return parsed.href;
}

/**
 * 移除 URL 的 hash
 * @example
 * removeHash('https://example.com/page#section')
 * // 'https://example.com/page'
 */
export function removeHash(url: string): string {
  const hashIndex = url.indexOf("#");
  return hashIndex === -1 ? url : url.slice(0, hashIndex);
}

// ============================================================================
// URL 编码
// ============================================================================

/**
 * 编码 URL 组件
 * @example
 * encodeURLComponent('hello world') // 'hello%20world'
 */
export function encodeURLComponent(str: string): string {
  return encodeURIComponent(str);
}

/**
 * 解码 URL 组件
 * @example
 * decodeURLComponent('hello%20world') // 'hello world'
 */
export function decodeURLComponent(str: string): string {
  try {
    return decodeURIComponent(str);
  } catch {
    return str;
  }
}

/**
 * 编码 URL（保留有效字符）
 * @example
 * encodeURL('https://example.com/path?q=hello world')
 * // 'https://example.com/path?q=hello%20world'
 */
export function encodeURL(url: string): string {
  return encodeURI(url);
}

/**
 * 解码 URL
 * @example
 * decodeURL('https://example.com/path?q=hello%20world')
 * // 'https://example.com/path?q=hello world'
 */
export function decodeURL(url: string): string {
  try {
    return decodeURI(url);
  } catch {
    return url;
  }
}

// ============================================================================
// URL 构建器
// ============================================================================

export interface URLBuilderOptions {
  protocol?: string;
  hostname?: string;
  port?: string | number;
  pathname?: string;
  query?: QueryParams;
  hash?: string;
  username?: string;
  password?: string;
}

/**
 * 构建 URL
 * @example
 * buildURL({
 *   protocol: 'https',
 *   hostname: 'example.com',
 *   pathname: '/path',
 *   query: { a: '1', b: '2' }
 * })
 * // 'https://example.com/path?a=1&b=2'
 */
export function buildURL(options: URLBuilderOptions): string {
  const { protocol = "https", hostname = "localhost", port, pathname = "/", query, hash, username, password } = options;

  let url = `${protocol}://${hostname}`;

  if (port) {
    url += `:${port}`;
  }

  url += pathname.startsWith("/") ? pathname : `/${pathname}`;

  if (query && Object.keys(query).length > 0) {
    url += buildQuery(query, { prefix: true });
  }

  if (hash) {
    url += hash.startsWith("#") ? hash : `#${hash}`;
  }

  // 添加认证信息
  if (username) {
    const auth = password ? `${username}:${password}` : username;
    url = url.replace("://", `://${auth}@`);
  }

  return url;
}

/**
 * 创建 URL 构建器（链式调用）
 * @example
 * createURLBuilder('https://example.com')
 *   .setPath('/api/users')
 *   .setQuery({ page: '1' })
 *   .setHash('top')
 *   .toString()
 * // 'https://example.com/api/users?page=1#top'
 */
export function createURLBuilder(baseURL: string = "") {
  let url = baseURL;

  const builder = {
    setProtocol(protocol: string) {
      const parsed = new URL(url || "http://localhost");
      parsed.protocol = protocol;
      url = parsed.href;
      return builder;
    },

    setHost(host: string) {
      const parsed = new URL(url || "http://localhost");
      parsed.host = host;
      url = parsed.href;
      return builder;
    },

    setPath(pathname: string) {
      const parsed = new URL(url || "http://localhost");
      parsed.pathname = pathname;
      url = parsed.href;
      return builder;
    },

    appendPath(path: string) {
      const parsed = new URL(url || "http://localhost");
      parsed.pathname = joinURL(parsed.pathname, path);
      url = parsed.href;
      return builder;
    },

    setQuery(params: QueryParams) {
      url = mergeQueryParams(url || "http://localhost", params);
      return builder;
    },

    addQuery(key: string, value: string | string[]) {
      url = setQueryParam(url || "http://localhost", key, value);
      return builder;
    },

    removeQuery(key: string) {
      url = removeQueryParam(url || "http://localhost", key);
      return builder;
    },

    setHash(hash: string) {
      url = setHash(url || "http://localhost", hash);
      return builder;
    },

    toString() {
      return url;
    },

    toURL() {
      return new URL(url);
    },
  };

  return builder;
}

// ============================================================================
// 特殊 URL 处理
// ============================================================================

/**
 * 创建 Data URL
 * @example
 * createDataURL('Hello', 'text/plain') // 'data:text/plain;base64,SGVsbG8='
 */
export function createDataURL(data: string, mimeType: string = "text/plain"): string {
  const base64 = btoa(unescape(encodeURIComponent(data)));
  return `data:${mimeType};base64,${base64}`;
}

/**
 * 解析 Data URL
 * @example
 * parseDataURL('data:text/plain;base64,SGVsbG8=')
 * // { mimeType: 'text/plain', data: 'Hello' }
 */
export function parseDataURL(dataURL: string): { mimeType: string; data: string } | null {
  const match = dataURL.match(/^data:([^;]+);base64,(.+)$/);

  if (!match) {
    return null;
  }

  try {
    const data = decodeURIComponent(escape(atob(match[2])));
    return {
      mimeType: match[1],
      data,
    };
  } catch {
    return null;
  }
}

/**
 * 创建 Blob URL
 * @example
 * const blob = new Blob(['Hello'], { type: 'text/plain' })
 * createBlobURL(blob) // 'blob:http://localhost/...'
 */
export function createBlobURL(blob: Blob): string {
  return URL.createObjectURL(blob);
}

/**
 * 释放 Blob URL
 * @example
 * revokeBlobURL(blobURL)
 */
export function revokeBlobURL(url: string): void {
  URL.revokeObjectURL(url);
}

/**
 * 获取当前页面 URL
 */
export function getCurrentURL(): string {
  if (typeof window === "undefined") {
    return "";
  }
  return window.location.href;
}

/**
 * 获取当前页面的查询参数
 */
export function getCurrentQueryParams(): QueryParams {
  if (typeof window === "undefined") {
    return {};
  }
  return parseQuery(window.location.search);
}
