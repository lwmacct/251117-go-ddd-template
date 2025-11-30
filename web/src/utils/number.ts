/**
 * 数字格式化工具
 * 提供货币、百分比、大数字等格式化函数
 */

export interface NumberFormatOptions {
  /** 语言 */
  locale?: string;
  /** 小数位数 */
  decimals?: number;
  /** 空值显示文本 */
  fallback?: string;
}

const defaultOptions: NumberFormatOptions = {
  locale: "zh-CN",
  decimals: 2,
  fallback: "-",
};

/**
 * 格式化数字（添加千分位分隔符）
 */
export function formatNumber(value: number | string | null | undefined, options: NumberFormatOptions = {}): string {
  const opts = { ...defaultOptions, ...options };

  if (value === null || value === undefined || value === "") {
    return opts.fallback!;
  }

  const num = typeof value === "string" ? parseFloat(value) : value;

  if (isNaN(num)) {
    return opts.fallback!;
  }

  return new Intl.NumberFormat(opts.locale, {
    minimumFractionDigits: 0,
    maximumFractionDigits: opts.decimals,
  }).format(num);
}

/**
 * 格式化货币
 */
export function formatCurrency(
  value: number | string | null | undefined,
  options: NumberFormatOptions & { currency?: string } = {}
): string {
  const opts = { ...defaultOptions, currency: "CNY", ...options };

  if (value === null || value === undefined || value === "") {
    return opts.fallback!;
  }

  const num = typeof value === "string" ? parseFloat(value) : value;

  if (isNaN(num)) {
    return opts.fallback!;
  }

  return new Intl.NumberFormat(opts.locale, {
    style: "currency",
    currency: opts.currency,
    minimumFractionDigits: opts.decimals,
    maximumFractionDigits: opts.decimals,
  }).format(num);
}

/**
 * 格式化百分比
 */
export function formatPercent(
  value: number | string | null | undefined,
  options: NumberFormatOptions & { multiply?: boolean } = {}
): string {
  const opts = { ...defaultOptions, multiply: false, decimals: 1, ...options };

  if (value === null || value === undefined || value === "") {
    return opts.fallback!;
  }

  let num = typeof value === "string" ? parseFloat(value) : value;

  if (isNaN(num)) {
    return opts.fallback!;
  }

  // 如果传入的是小数（如 0.15），乘以 100
  if (opts.multiply) {
    num = num * 100;
  }

  return `${num.toFixed(opts.decimals)}%`;
}

/**
 * 格式化大数字（如 1.2万、3.5亿）
 */
export function formatCompact(value: number | string | null | undefined, options: NumberFormatOptions = {}): string {
  const opts = { ...defaultOptions, decimals: 1, ...options };

  if (value === null || value === undefined || value === "") {
    return opts.fallback!;
  }

  const num = typeof value === "string" ? parseFloat(value) : value;

  if (isNaN(num)) {
    return opts.fallback!;
  }

  const absNum = Math.abs(num);
  const sign = num < 0 ? "-" : "";

  if (absNum >= 100000000) {
    return `${sign}${(absNum / 100000000).toFixed(opts.decimals)}亿`;
  }
  if (absNum >= 10000) {
    return `${sign}${(absNum / 10000).toFixed(opts.decimals)}万`;
  }

  return formatNumber(num, opts);
}

/**
 * 格式化文件大小
 */
export function formatFileSize(
  bytes: number | string | null | undefined,
  options: Pick<NumberFormatOptions, "decimals" | "fallback"> = {}
): string {
  const opts = { decimals: 2, fallback: "-", ...options };

  if (bytes === null || bytes === undefined || bytes === "") {
    return opts.fallback;
  }

  const num = typeof bytes === "string" ? parseFloat(bytes) : bytes;

  if (isNaN(num) || num < 0) {
    return opts.fallback;
  }

  if (num === 0) {
    return "0 B";
  }

  const units = ["B", "KB", "MB", "GB", "TB", "PB"];
  const exponent = Math.min(Math.floor(Math.log(num) / Math.log(1024)), units.length - 1);
  const value = num / Math.pow(1024, exponent);

  return `${value.toFixed(opts.decimals)} ${units[exponent]}`;
}

/**
 * 格式化持续时间（秒转为 HH:MM:SS 或友好格式）
 */
export function formatDuration(
  seconds: number | string | null | undefined,
  options: { format?: "time" | "friendly"; fallback?: string } = {}
): string {
  const opts = { format: "friendly" as const, fallback: "-", ...options };

  if (seconds === null || seconds === undefined || seconds === "") {
    return opts.fallback;
  }

  const num = typeof seconds === "string" ? parseFloat(seconds) : seconds;

  if (isNaN(num) || num < 0) {
    return opts.fallback;
  }

  const hours = Math.floor(num / 3600);
  const minutes = Math.floor((num % 3600) / 60);
  const secs = Math.floor(num % 60);

  if (opts.format === "time") {
    const pad = (n: number) => n.toString().padStart(2, "0");
    if (hours > 0) {
      return `${pad(hours)}:${pad(minutes)}:${pad(secs)}`;
    }
    return `${pad(minutes)}:${pad(secs)}`;
  }

  // 友好格式
  const parts: string[] = [];
  if (hours > 0) parts.push(`${hours}小时`);
  if (minutes > 0) parts.push(`${minutes}分钟`);
  if (secs > 0 || parts.length === 0) parts.push(`${secs}秒`);

  return parts.join(" ");
}

/**
 * 将数字转换为序数词
 */
export function toOrdinal(value: number): string {
  return `第${value}`;
}

/**
 * 数字范围限制
 */
export function clamp(value: number, min: number, max: number): number {
  return Math.min(Math.max(value, min), max);
}
