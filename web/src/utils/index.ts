/**
 * 工具函数统一导出
 * 提供所有工具函数的统一入口
 */

// ============================================================================
// 数组工具
// ============================================================================

export {
  // 分组与分块
  chunk,
  groupBy,
  partition,
  // 去重与交集
  unique,
  intersection,
  difference,
  union,
  // 查找与索引
  findIndex,
  findLastIndex,
  findAllIndices,
  // 排序
  sortBy,
  sortByMultiple,
  // 变换
  flatten,
  shuffle,
  sample,
  // 移动与交换
  move,
  swap,
  // 聚合
  sum,
  average,
  maxBy,
  minBy,
  // 树形结构
  arrayToTree,
  treeToArray,
  findInTree,
  filterTree,
  type TreeNode,
  type ArrayToTreeOptions,
} from "./array";

// ============================================================================
// 对象工具
// ============================================================================

export {
  // 深度操作
  deepClone,
  deepMerge,
  deepEqual,
  // 路径操作
  parsePath,
  get,
  set,
  has,
  unset,
  // 筛选与转换
  pick,
  omit,
  filterObject,
  mapObject,
  invert,
  // 类型检查
  isObject,
  isEmpty,
  isPlainObject,
  // 遍历
  forEachObject,
  deepForEach,
  // 差异比较
  diff,
  // 其他
  compact,
  defaults,
  entries,
  fromEntries,
  type ObjectDiff,
} from "./object";

// ============================================================================
// 字符串工具
// ============================================================================

export {
  // 截断与填充
  truncate,
  truncateMiddle,
  padStart,
  padEnd,
  // 大小写转换
  capitalize,
  capitalizeWords,
  toCamelCase,
  toKebabCase,
  toSnakeCase,
  toPascalCase,
  // URL 与路径
  slugify,
  // 掩码
  maskEmail,
  maskPhone,
  maskIdCard,
  maskBankCard,
  mask,
  // 格式化
  template,
  // 验证
  isEmail,
  isPhone,
  isUrl,
  isIdCard,
  // 转换
  escapeHtml,
  unescapeHtml,
  escapeRegExp,
  // 其他
  removeWhitespace,
  normalizeWhitespace,
  countWords,
  countChars,
  getInitials,
  highlight,
  generateRandomString,
} from "./string";

// ============================================================================
// 数字工具
// ============================================================================

export {
  formatNumber,
  formatCurrency,
  formatPercent,
  formatFileSize,
  formatCompact,
  toFixed,
  clamp,
  random,
  randomInt,
  isInRange,
  round,
  floor,
  ceil,
  formatDuration,
  parseNumber,
  toOrdinal,
  type FormatNumberOptions,
  type FormatCurrencyOptions,
  type FormatFileSizeOptions,
  type FormatDurationOptions,
} from "./number";

// ============================================================================
// 日期工具
// ============================================================================

export {
  formatDate,
  formatDateTime,
  formatTime,
  formatRelative,
  formatDateRange,
  isToday,
  isYesterday,
  isTomorrow,
  isThisWeek,
  isThisMonth,
  isThisYear,
  isSameDay,
  startOfDay,
  endOfDay,
  startOfWeek,
  endOfWeek,
  startOfMonth,
  endOfMonth,
  addDays,
  addMonths,
  addYears,
  differenceInDays,
  differenceInMonths,
  differenceInYears,
  parseDate,
  type FormatDateOptions,
} from "./date";

// ============================================================================
// 验证工具
// ============================================================================

export { rules, createValidator, type ValidationRule } from "./validation";

// ============================================================================
// 平台检测
// ============================================================================

export {
  detectBrowser,
  detectOS,
  detectDevice,
  getPlatformInfo,
  // 快捷常量
  isChrome,
  isFirefox,
  isSafari,
  isEdge,
  isOpera,
  isIE,
  isWindows,
  isMacOS,
  isLinux,
  isIOS,
  isAndroid,
  isMobile,
  isTablet,
  isDesktop,
  isTouchDevice,
  supportsWebGL,
  supportsWebP,
  supportsServiceWorker,
  supportsNotification,
  supportsGeolocation,
  type BrowserInfo,
  type OSInfo,
  type DeviceInfo,
  type PlatformInfo,
} from "./platform";

// ============================================================================
// 限流工具
// ============================================================================

export {
  throttle,
  debounce,
  createRateLimiter,
  retry,
  withTimeout,
  createDeduplicator,
  type ThrottleOptions,
  type ThrottledFunction,
  type RateLimiterOptions,
  type RateLimiter,
  type RetryOptions,
} from "./throttle";

// ============================================================================
// ID 生成
// ============================================================================

export {
  // UUID
  uuid,
  shortUuid,
  isValidUuid,
  // NanoID 风格
  nanoid,
  customId,
  alphanumericId,
  numericId,
  alphabeticId,
  hexId,
  // 时间戳 ID
  timestampId,
  sortableId,
  extractTimestamp,
  // 前缀 ID
  prefixedId,
  createPrefixedIdGenerator,
  // 序列 ID
  createSequence,
  createFormattedSequence,
  // 高级 ID
  createSnowflake,
  ulid,
  extractUlidTimestamp,
  // 实用工具
  uniqueDomId,
  createIdFactory,
  generateIds,
  isUniqueId,
  ensureUniqueId,
  type SnowflakeConfig,
} from "./id";

// ============================================================================
// 导入导出
// ============================================================================

export { exportToCSV, type ExportOptions } from "./export";

export {
  parseCSV,
  validateImportData,
  type ParseCSVOptions,
  type ValidationResult,
  type FieldValidator,
} from "./import";

// ============================================================================
// 颜色工具
// ============================================================================

export {
  // 解析
  parseHex,
  parseRgb,
  parseHsl,
  parseColor,
  // 转换
  rgbToHex,
  rgbToHsl,
  hslToRgb,
  rgbToHsv,
  hsvToRgb,
  // 格式化
  formatRgb,
  formatHsl,
  // 颜色操作
  lighten,
  darken,
  saturate,
  desaturate,
  setAlpha,
  invert as invertColor,
  grayscale,
  mix,
  complement,
  // 颜色分析
  getLuminance,
  getContrast,
  isDark,
  isLight,
  getTextColor,
  // 颜色生成
  randomColor,
  generateGradient,
  generatePalette,
  // 其他
  getNamedColor,
  isValidColor,
  type RGB,
  type RGBA,
  type HSL,
  type HSLA,
  type HSV,
} from "./color";

// ============================================================================
// URL 工具
// ============================================================================

export {
  // URL 解析
  parseURL,
  isValidURL,
  isAbsoluteURL,
  isRelativeURL,
  // 查询参数
  parseQuery,
  buildQuery,
  getQueryParam,
  setQueryParam,
  removeQueryParam,
  mergeQueryParams,
  // URL 操作
  joinURL,
  normalizeURL,
  getBasePath,
  getFileName,
  getExtension,
  setHash,
  removeHash,
  // URL 编码
  encodeURLComponent,
  decodeURLComponent,
  encodeURL,
  decodeURL,
  // URL 构建
  buildURL,
  createURLBuilder,
  // 特殊 URL
  createDataURL,
  parseDataURL,
  createBlobURL,
  revokeBlobURL,
  getCurrentURL,
  getCurrentQueryParams,
  type ParsedURL,
  type QueryParams,
  type URLBuilderOptions,
} from "./url";

// ============================================================================
// DOM 工具
// ============================================================================

export {
  // 元素查询
  getElement,
  getElements,
  elementExists,
  waitForElement,
  // 类名操作
  addClass,
  removeClass,
  toggleClass,
  hasClass,
  replaceClass,
  // 样式操作
  getStyle,
  setStyle,
  removeStyle,
  getCSSVariable,
  setCSSVariable,
  // 属性操作
  getAttribute,
  setAttribute,
  removeAttribute,
  hasAttribute,
  getDataAttribute,
  setDataAttribute,
  // 尺寸和位置
  getRect,
  getSize,
  getOffset,
  getWindowSize,
  getDocumentSize,
  getScrollPosition,
  // 滚动操作
  scrollTo,
  scrollToTop,
  scrollToBottom,
  isInViewport,
  isPartiallyInViewport,
  // 焦点操作
  focus,
  blur,
  getActiveElement,
  hasFocus,
  // 元素创建和操作
  createElement,
  removeElement,
  cloneElement,
  insertBefore,
  insertAfter,
  wrap,
  unwrap,
  // 事件工具
  on,
  off,
  once,
  trigger,
  // 可见性
  show,
  hide,
  toggle,
  isVisible,
  isHidden,
  type ElementTarget,
  type Position,
  type Size,
  type Rect,
} from "./dom";

// ============================================================================
// Cookie 工具
// ============================================================================

export {
  // 基础操作
  getCookie,
  setCookie,
  removeCookie,
  hasCookie,
  // 批量操作
  getAllCookies,
  getCookies,
  setCookies,
  removeCookies,
  clearAllCookies,
  // JSON Cookie
  getJsonCookie,
  setJsonCookie,
  // 解析工具
  parseCookieString,
  serializeCookie,
  // 工具函数
  getCookieCount,
  getCookieNames,
  areCookiesEnabled,
  getCookiesSize,
  getCookiesRemainingSpace,
  // 管理器
  createCookieManager,
  // 预设配置
  SESSION_COOKIE,
  PERSISTENT_COOKIE,
  SECURE_COOKIE,
  CROSS_SITE_COOKIE,
  type CookieOptions,
  type CookieManager,
} from "./cookie";

// ============================================================================
// 类型工具
// ============================================================================

export {
  // 基础类型检查
  isString,
  isNumber,
  isFiniteNumber,
  isInteger,
  isBoolean,
  isNull,
  isUndefined,
  isNullish,
  isDefined,
  hasValue,
  // 复杂类型检查
  isObject as isTypeObject,
  isPlainObject as isTypePlainObject,
  isArray,
  isArrayOf,
  isNonEmptyArray,
  isFunction,
  isAsyncFunction,
  isSymbol,
  isBigInt,
  // 特殊类型检查
  isDate,
  isRegExp,
  isError,
  isPromise,
  isMap,
  isSet,
  isWeakMap,
  isWeakSet,
  // DOM 类型检查
  isElement,
  isHTMLElement,
  isNode,
  isBlob,
  isFile,
  isFormData,
  isArrayBuffer,
  // 字符串检查
  isEmptyString,
  isBlankString,
  isNonEmptyString,
  isNonBlankString,
  // 类型转换
  toString as typeToString,
  toNumber,
  toInteger,
  toBoolean,
  toArray,
  toDate,
  // 断言函数
  assertDefined,
  assert,
  assertNever,
  // 类型守卫组合
  unionGuard,
  intersectionGuard,
  notGuard,
  // 实用工具
  getType,
  isSameType,
  isPrimitive,
  isIterable,
  isJsonSerializable,
} from "./type";
