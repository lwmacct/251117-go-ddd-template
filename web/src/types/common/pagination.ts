/**
 * 通用分页类型定义
 */

/** 分页请求参数 */
export interface PaginationParams {
  page: number;
  limit: number;
  search?: string;
  sort_by?: string;
  sort_order?: "asc" | "desc";
}

/** 分页元数据 */
export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

/** 分页响应 */
export interface PaginatedResponse<T> {
  data: T[];
  pagination: PaginationMeta;
}
