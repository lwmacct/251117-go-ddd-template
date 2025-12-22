/**
 * 认证相关类型定义
 */
import type { User } from "./user";

/**
 * 登录请求（标准版，带验证码）
 * 支持多种登录方式：手机号/用户名/邮箱
 */
export interface LoginRequest {
  account: string; // 手机号/用户名/邮箱
  password: string;
  captcha_id: string;
  captcha: string;
}

/**
 * 注册请求（标准版，带验证码）
 */
export interface RegisterRequest {
  email: string;
  password: string;
  captcha_id: string;
  captcha: string;
}

// ============================================================================
// 已废弃的类型（向后兼容）
// ============================================================================

/**
 * @deprecated 使用 LoginRequest 代替
 * 基础登录请求（不带验证码，已废弃）
 */
export interface BasicLoginRequest {
  login: string; // 用户名或邮箱
  password: string;
}

/**
 * @deprecated 使用 RegisterRequest 代替
 * 基础注册请求（不带验证码，已废弃）
 */
export interface BasicRegisterRequest {
  username: string;
  email: string;
  password: string;
  full_name?: string;
}

/** 验证码数据 */
export interface CaptchaData {
  id: string;
  image: string; // Base64 编码的图片
}

/** 认证响应 */
export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
  user: User;
}

/** 登录结果 (支持 2FA)  */
export interface LoginResult {
  success: boolean;
  message?: string;
  requiresTwoFactor: boolean;
  sessionToken?: string;
}
