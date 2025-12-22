/**
 * 日期/时间格式化工具
 * 提供统一的日期时间格式化函数
 */

export interface DateFormatOptions {
  /** 是否显示时间 */
  showTime?: boolean;
  /** 是否显示秒 */
  showSeconds?: boolean;
  /** 空值显示文本 */
  fallback?: string;
  /** 语言 */
  locale?: string;
}

const defaultOptions: DateFormatOptions = {
  showTime: true,
  showSeconds: false,
  fallback: "-",
  locale: "zh-CN",
};

/**
 * 格式化日期时间
 * @param date 日期字符串、Date 对象或时间戳
 * @param options 格式化选项
 */
export function formatDateTime(
  date: string | Date | number | null | undefined,
  options: DateFormatOptions = {}
): string {
  const opts = { ...defaultOptions, ...options };

  if (!date) return opts.fallback!;

  const d = typeof date === "string" || typeof date === "number" ? new Date(date) : date;

  if (isNaN(d.getTime())) return opts.fallback!;

  const formatOptions: Intl.DateTimeFormatOptions = {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  };

  if (opts.showTime) {
    formatOptions.hour = "2-digit";
    formatOptions.minute = "2-digit";
    if (opts.showSeconds) {
      formatOptions.second = "2-digit";
    }
  }

  return d.toLocaleString(opts.locale, formatOptions);
}

/**
 * 格式化日期（不含时间）
 */
export function formatDate(
  date: string | Date | number | null | undefined,
  options: Omit<DateFormatOptions, "showTime" | "showSeconds"> = {}
): string {
  return formatDateTime(date, { ...options, showTime: false });
}

/**
 * 格式化时间（不含日期）
 */
export function formatTime(
  date: string | Date | number | null | undefined,
  options: Pick<DateFormatOptions, "showSeconds" | "fallback" | "locale"> = {}
): string {
  const opts = { ...defaultOptions, ...options };

  if (!date) return opts.fallback!;

  const d = typeof date === "string" || typeof date === "number" ? new Date(date) : date;

  if (isNaN(d.getTime())) return opts.fallback!;

  const formatOptions: Intl.DateTimeFormatOptions = {
    hour: "2-digit",
    minute: "2-digit",
  };

  if (opts.showSeconds) {
    formatOptions.second = "2-digit";
  }

  return d.toLocaleTimeString(opts.locale, formatOptions);
}

/**
 * 相对时间格式化
 * 如："刚刚"、"5 分钟前"、"2 小时前"、"昨天"、"3 天前"
 */
export function formatRelativeTime(
  date: string | Date | number | null | undefined,
  options: Pick<DateFormatOptions, "fallback"> = {}
): string {
  const fallback = options.fallback ?? "-";

  if (!date) return fallback;

  const d = typeof date === "string" || typeof date === "number" ? new Date(date) : date;

  if (isNaN(d.getTime())) return fallback;

  const now = new Date();
  const diff = now.getTime() - d.getTime();
  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);
  const months = Math.floor(days / 30);
  const years = Math.floor(days / 365);

  if (seconds < 60) {
    return "刚刚";
  } else if (minutes < 60) {
    return `${minutes} 分钟前`;
  } else if (hours < 24) {
    return `${hours} 小时前`;
  } else if (days === 1) {
    return "昨天";
  } else if (days < 7) {
    return `${days} 天前`;
  } else if (days < 30) {
    const weeks = Math.floor(days / 7);
    return `${weeks} 周前`;
  } else if (months < 12) {
    return `${months} 个月前`;
  } else {
    return `${years} 年前`;
  }
}

/**
 * 判断是否是今天
 */
export function isToday(date: string | Date | number): boolean {
  const d = typeof date === "string" || typeof date === "number" ? new Date(date) : date;
  const today = new Date();

  return (
    d.getDate() === today.getDate() && d.getMonth() === today.getMonth() && d.getFullYear() === today.getFullYear()
  );
}

/**
 * 判断是否是昨天
 */
export function isYesterday(date: string | Date | number): boolean {
  const d = typeof date === "string" || typeof date === "number" ? new Date(date) : date;
  const yesterday = new Date();
  yesterday.setDate(yesterday.getDate() - 1);

  return (
    d.getDate() === yesterday.getDate() &&
    d.getMonth() === yesterday.getMonth() &&
    d.getFullYear() === yesterday.getFullYear()
  );
}

/**
 * 智能格式化（根据时间远近选择格式）
 * - 今天：显示时间
 * - 昨天：显示"昨天 HH:mm"
 * - 今年：显示"MM-DD HH:mm"
 * - 更早：显示"YYYY-MM-DD"
 */
export function formatSmart(
  date: string | Date | number | null | undefined,
  options: Pick<DateFormatOptions, "fallback" | "locale"> = {}
): string {
  const opts = { ...defaultOptions, ...options };

  if (!date) return opts.fallback!;

  const d = typeof date === "string" || typeof date === "number" ? new Date(date) : date;

  if (isNaN(d.getTime())) return opts.fallback!;

  const now = new Date();

  if (isToday(d)) {
    return formatTime(d, { locale: opts.locale });
  }

  if (isYesterday(d)) {
    return `昨天 ${formatTime(d, { locale: opts.locale })}`;
  }

  if (d.getFullYear() === now.getFullYear()) {
    return d.toLocaleString(opts.locale, {
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
    });
  }

  return formatDate(d, { locale: opts.locale });
}
