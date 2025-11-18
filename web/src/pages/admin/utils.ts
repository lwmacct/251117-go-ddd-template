/**
 * 快手页面通用工具函数
 */

/** 格式化日期为 YYYYMMDD */
export const formatDateToYYYYMMDD = (date: Date): string => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}${month}${day}`;
};

/** 解析 YYYYMMDD 为 Date */
export const parseDateFromYYYYMMDD = (dateStr: string): Date => {
  const year = parseInt(dateStr.substring(0, 4));
  const month = parseInt(dateStr.substring(4, 6)) - 1;
  const day = parseInt(dateStr.substring(6, 8));
  return new Date(year, month, day);
};

/** 获取今天的日期字符串 (YYYYMMDD) */
export const getTodayString = (): string => {
  return formatDateToYYYYMMDD(new Date());
};

/** 获取昨天的日期字符串 (YYYYMMDD) */
export const getYesterdayString = (): string => {
  const yesterday = new Date();
  yesterday.setDate(yesterday.getDate() - 1);
  return formatDateToYYYYMMDD(yesterday);
};
