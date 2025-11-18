/**
 * API 响应类型定义
 */

/** API 响应 */
export interface ApiResponse<T = any> {
  message?: string;
  data?: T;
  error?: string;
}
