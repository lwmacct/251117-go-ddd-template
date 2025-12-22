/**
 * Personal Access Token 相关类型定义
 */

/** Token 状态 */
export type TokenStatus = "active" | "disabled" | "expired";

/** Personal Access Token */
export interface PersonalAccessToken {
  id: number;
  user_id: number;
  name: string;
  token_prefix: string;
  permissions: string[];
  ip_whitelist?: string[];
  expires_at?: string;
  last_used_at?: string;
  status: TokenStatus;
  created_at: string;
  updated_at?: string;
}

/** 创建 Token 请求 */
export interface CreateTokenRequest {
  name: string;
  permissions?: string[];
  expires_at?: string;
  expires_in?: number;
  ip_whitelist?: string[];
  description?: string;
}

/** 创建 Token 响应（包含完整 token） */
export interface CreateTokenResponse {
  token: PersonalAccessToken;
  plain_token: string; // 明文 token，仅在创建时返回一次
}
