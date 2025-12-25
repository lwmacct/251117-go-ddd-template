/**
 * 统一错误处理模块
 * 提供类型安全的错误类和提取函数
 */
import type { AxiosError } from "axios";
import type { ErrorResponse } from "./types";

/**
 * 应用统一错误类
 * 所有 API 错误都会被转换为此类型
 */
export class AppError extends Error {
  constructor(
    message: string,
    public readonly code: string = "UNKNOWN",
    public readonly status: number = 500,
    public readonly details?: Record<string, unknown>,
  ) {
    super(message);
    this.name = "AppError";
  }

  /**
   * 类型守卫 - 用于安全地判断是否为 AppError
   */
  static isAppError(error: unknown): error is AppError {
    return error instanceof AppError;
  }
}

/**
 * 从任意错误中提取 AppError
 * 在 axios 拦截器中使用，确保所有错误都被统一处理
 */
export function extractErrorFromAxios(error: unknown): AppError {
  // 1. 已经是 AppError，直接返回
  if (AppError.isAppError(error)) {
    return error;
  }

  // 2. AxiosError - 最常见的情况
  if (isAxiosError(error)) {
    const response = error.response?.data as ErrorResponse | undefined;
    const status = error.response?.status ?? 500;

    // 优先使用结构化错误详情
    if (response?.error) {
      return new AppError(
        response.error.message || response.message,
        response.error.code || String(status),
        status,
        response.error.details as Record<string, unknown> | undefined,
      );
    }

    // 回退到顶层 message
    if (response?.message) {
      return new AppError(response.message, String(status), status);
    }

    // 网络错误
    if (error.code === "ECONNABORTED") {
      return new AppError("请求超时", "TIMEOUT", 408);
    }
    if (!error.response) {
      return new AppError("网络连接失败", "NETWORK_ERROR", 0);
    }

    return new AppError(error.message, String(status), status);
  }

  // 3. 标准 Error
  if (error instanceof Error) {
    return new AppError(error.message);
  }

  // 4. 未知类型
  return new AppError(String(error));
}

/**
 * AxiosError 类型守卫
 */
function isAxiosError(error: unknown): error is AxiosError {
  return typeof error === "object" && error !== null && "isAxiosError" in error && (error as AxiosError).isAxiosError === true;
}
