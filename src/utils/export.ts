/**
 * 数据导出工具函数
 */

/**
 * CSV 列配置
 */
export interface CSVColumn<T> {
  /** 列标题 */
  header: string;
  /** 数据字段名或取值函数 */
  key: keyof T | ((item: T) => string | number | undefined);
}

/**
 * 导出配置
 */
export interface ExportOptions {
  /** 文件名（不含扩展名） */
  filename: string;
  /** 是否添加 BOM（用于 Excel 正确显示中文） */
  withBOM?: boolean;
}

/**
 * 转义 CSV 字段值
 * 处理包含逗号、引号、换行符的值
 */
function escapeCSVValue(value: string | number | undefined | null): string {
  if (value === undefined || value === null) {
    return "";
  }

  const str = String(value);

  // 如果包含逗号、引号或换行符，需要用双引号包裹
  if (str.includes(",") || str.includes('"') || str.includes("\n") || str.includes("\r")) {
    // 将引号替换为两个引号（CSV 转义规则）
    return `"${str.replace(/"/g, '""')}"`;
  }

  return str;
}

/**
 * 将数据导出为 CSV 文件
 *
 * @param data - 要导出的数据数组
 * @param columns - 列配置
 * @param options - 导出选项
 *
 * @example
 * ```ts
 * exportToCSV(logs, [
 *   { header: 'ID', key: 'id' },
 *   { header: '用户', key: (item) => item.username || item.user_id },
 *   { header: '时间', key: 'created_at' },
 * ], { filename: 'audit-logs' });
 * ```
 */
export function exportToCSV<T>(data: T[], columns: CSVColumn<T>[], options: ExportOptions): void {
  const { filename, withBOM = true } = options;

  // 构建表头
  const headers = columns.map((col) => escapeCSVValue(col.header));
  const headerRow = headers.join(",");

  // 构建数据行
  const dataRows = data.map((item) => {
    const values = columns.map((col) => {
      let value: string | number | undefined;

      if (typeof col.key === "function") {
        value = col.key(item);
      } else {
        value = item[col.key] as string | number | undefined;
      }

      return escapeCSVValue(value);
    });

    return values.join(",");
  });

  // 组合 CSV 内容
  const csvContent = [headerRow, ...dataRows].join("\n");

  // 添加 BOM 以支持 Excel 正确显示中文
  const bom = withBOM ? "\uFEFF" : "";
  const blob = new Blob([bom + csvContent], { type: "text/csv;charset=utf-8" });

  // 创建下载链接
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = `${filename}.csv`;

  // 触发下载
  document.body.appendChild(link);
  link.click();

  // 清理
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}

/**
 * 格式化日期为本地字符串
 */
export function formatDateForExport(dateString?: string): string {
  if (!dateString) return "";

  return new Date(dateString).toLocaleString("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}
