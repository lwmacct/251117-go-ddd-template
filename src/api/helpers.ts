/**
 * API 响应辅助函数
 *
 * 提供从 OpenAPI 生成的响应中提取数据的通用函数，
 * 替代原有的适配层，使组件可以直接调用生成的 API。
 */

import type { ResponsePaginationMeta } from "@models";

/**
 * 分页响应结果
 */
export interface ListResult<T> {
  data: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

/**
 * 从列表 API 响应中提取数据和分页信息
 *
 * @example
 * const response = await adminUserApi.apiAdminUsersGet(page, limit);
 * const result = extractList(response.data);
 * users.value = result.data;
 * Object.assign(pagination, result.pagination);
 */
export function extractList<T>(response: { data?: T[]; meta?: ResponsePaginationMeta }): ListResult<T> {
  return {
    data: response.data ?? [],
    pagination: {
      page: response.meta?.page ?? 1,
      limit: response.meta?.per_page ?? 20,
      total: response.meta?.total ?? 0,
      total_pages: response.meta?.total_pages ?? 0,
    },
  };
}

/**
 * 从单条数据 API 响应中提取数据
 *
 * @example
 * const response = await adminUserApi.apiAdminUsersIdGet(id);
 * const user = extractData(response.data);
 */
export function extractData<T>(response: { data?: T }): T | undefined {
  return response.data;
}

/**
 * 从单条数据 API 响应中提取数据（断言非空）
 *
 * @example
 * const response = await authApi.apiAuthLoginPost(credentials);
 * const authResult = extractDataRequired(response.data);
 */
export function extractDataRequired<T>(response: { data?: T }): T {
  if (response.data === undefined) {
    throw new Error("Expected data in response but got undefined");
  }
  return response.data;
}
