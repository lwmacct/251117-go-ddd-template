/**
 * 快手 API 类型定义（精简版）
 */

/** 通用响应结构 */
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

/** 快手文本响应 */
export type KuaishouTextResponse = string;
