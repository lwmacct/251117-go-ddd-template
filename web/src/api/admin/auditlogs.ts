/**
 * Admin 审计日志 API
 */
import { apiClient } from "../auth/client";
import { normalizeListResponse } from "../helpers/pagination";
import type { ApiResponse, ListApiResponse } from "@/types/response";
import type { AuditLog, AuditLogQueryParams } from "@/types/admin/audit";
import type { PaginatedResponse } from "@/types/common";

/**
 * 获取审计日志列表（分页 + 过滤）
 */
export const listAuditLogs = async (params: AuditLogQueryParams = {}): Promise<PaginatedResponse<AuditLog>> => {
  const page = params.page ?? 1;
  const limit = params.limit ?? 20;

  const { data } = await apiClient.get<ListApiResponse<AuditLog[]>>("/api/admin/audit-logs", {
    params: {
      page,
      limit,
      user_id: params.user_id,
      action: params.action,
      resource: params.resource,
      status: params.status,
      start_date: params.start_date,
      end_date: params.end_date,
    },
  });

  return normalizeListResponse<AuditLog>(data, { page, limit });
};

/**
 * 获取审计日志详情
 */
export const getAuditLog = async (id: number): Promise<AuditLog> => {
  const { data } = await apiClient.get<ApiResponse<AuditLog>>(`/api/admin/audit-logs/${id}`);

  if (data.data) {
    return data.data;
  }

  throw new Error(data.error || "获取审计日志详情失败");
};
