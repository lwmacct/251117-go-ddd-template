/**
 * API 响应类型定义
 * 与后端 internal/adapters/http/response/response.go 保持一致
 */

/**
 * 统一 API 响应结构
 * 对应后端 UnifiedResponse
 */
export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data?: T;
  error?: string;
}

/**
 * 列表响应结构（带分页）
 * 对应后端 ListResponse
 */
export interface ListApiResponse<T = unknown> {
  code: number;
  message: string;
  data: T;
  meta?: PaginationMeta;
}

/**
 * 分页元数据
 * 对应后端 PaginationMeta
 */
export interface PaginationMeta {
  total: number;
  page: number;
  per_page: number;
  total_pages?: number;
  has_more?: boolean;
}

/**
 * 错误响应结构
 * 对应后端 ErrorResponse
 */
export interface ErrorResponse {
  code: number;
  message: string;
  error?: ErrorDetail;
}

/**
 * 错误详情
 * 对应后端 ErrorDetail
 */
export interface ErrorDetail {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}
