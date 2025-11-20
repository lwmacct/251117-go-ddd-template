/**
 * 工具函数：将后端 ListResponse 映射为前端使用的 PaginatedResponse
 */
import type { ListApiResponse } from "@/types/response";
import type { PaginatedResponse } from "@/types/common";

interface PaginationFallback {
  page: number;
  limit: number;
  total?: number;
}

/**
 * normalizeListResponse 将后端 ListResponse 转换为前端统一的 PaginatedResponse
 */
export function normalizeListResponse<T>(
  payload: ListApiResponse<T[]>,
  fallback: PaginationFallback
): PaginatedResponse<T> {
  if (!Array.isArray(payload.data)) {
    throw new Error(payload.message || "Invalid list response");
  }

  const meta = payload.meta;
  const limit = meta?.per_page ?? fallback.limit;
  const page = meta?.page ?? fallback.page;
  const total = meta?.total ?? fallback.total ?? payload.data.length;
  const totalPages = meta?.total_pages ?? Math.max(Math.ceil(total / Math.max(limit, 1)), 1);

  return {
    data: payload.data,
    pagination: {
      page,
      limit,
      total,
      total_pages: totalPages,
      has_more: meta?.has_more ?? page < totalPages,
    },
  };
}
